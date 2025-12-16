from pyrogram import Client, filters
from helpers.utils import get_text
from helpers.mongo import db

@Client.on_message(filters.command("save", prefixes=".") & filters.me)
async def save_note(client, message):
    args = message.text.split(None, 2)
    if len(args) < 2:
        await message.edit("**Usage:** `.save [name] [content]` or reply to a message")
        return
    
    name = args[1].lower()
    
    if message.reply_to_message:
        content = {
            "text": message.reply_to_message.text or message.reply_to_message.caption,
            "file_id": None
        }
        if message.reply_to_message.media:
            media_type = message.reply_to_message.media.value
            content["file_id"] = getattr(message.reply_to_message, media_type).file_id
            content["media_type"] = media_type
    elif len(args) > 2:
        content = {"text": args[2], "file_id": None}
    else:
        await message.edit("**Provide content to save.**")
        return
    
    await db.add_note(message.chat.id, name, content)
    await message.edit(f"**✅ Note `{name}` saved to Database.**")

@Client.on_message(filters.command("get", prefixes=".") & filters.me)
async def get_note(client, message):
    args = message.text.split(None, 1)
    if len(args) < 2:
        await message.edit("**Usage:** `.get [name]`")
        return
    
    name = args[1].lower()
    
    note_doc = await db.get_note(message.chat.id, name)
    
    if not note_doc:
        await message.edit(f"**Note `{name}` not found.**")
        return
    
    note = note_doc["content"]
    
    if note.get("file_id"):
        media_type = note.get("media_type", "document")
        try:
            if media_type == "photo":
                await client.send_photo(message.chat.id, note["file_id"], caption=note.get("text"))
            elif media_type == "video":
                await client.send_video(message.chat.id, note["file_id"], caption=note.get("text"))
            elif media_type == "audio":
                await client.send_audio(message.chat.id, note["file_id"], caption=note.get("text"))
            elif media_type == "voice":
                await client.send_voice(message.chat.id, note["file_id"], caption=note.get("text"))
            elif media_type == "sticker":
                await client.send_sticker(message.chat.id, note["file_id"])
            else:
                await client.send_document(message.chat.id, note["file_id"], caption=note.get("text"))
            
            await message.delete()
        except Exception as e:
            await message.edit(f"**Error sending media:** {e}")
    else:
        await message.edit(note["text"])

@Client.on_message(filters.command("notes", prefixes=".") & filters.me)
async def list_notes(client, message):
    notes = await db.get_all_notes(message.chat.id)
    
    if not notes:
        await message.edit("**No notes saved in this chat.**")
        return
    
    note_list = "\n".join([f"• `{n['name']}`" for n in notes])
    await message.edit(f"**📝 Saved Notes (DB):**\n{note_list}")

@Client.on_message(filters.command("clear", prefixes=".") & filters.me)
async def clear_note(client, message):
    args = message.text.split(None, 1)
    if len(args) < 2:
        await message.edit("**Usage:** `.clear [name]`")
        return
    
    name = args[1].lower()
    
    await db.delete_note(message.chat.id, name)
    await message.edit(f"**✅ Note `{name}` deleted.**")

@Client.on_message(filters.command("clearall", prefixes=".") & filters.me)
async def clear_all_notes(client, message):
    await db.delete_all_notes(message.chat.id)
    await message.edit("**✅ All notes cleared from Database.**")
