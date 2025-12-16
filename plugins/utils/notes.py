from pyrogram import Client, filters
from helpers.utils import get_text

NOTES = {}

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
            content["file_id"] = getattr(message.reply_to_message, message.reply_to_message.media.value).file_id
            content["media_type"] = message.reply_to_message.media.value
    elif len(args) > 2:
        content = {"text": args[2], "file_id": None}
    else:
        await message.edit("**Provide content to save.**")
        return
    
    chat_id = str(message.chat.id)
    if chat_id not in NOTES:
        NOTES[chat_id] = {}
    
    NOTES[chat_id][name] = content
    await message.edit(f"**✅ Note `{name}` saved.**")

@Client.on_message(filters.command("get", prefixes=".") & filters.me)
async def get_note(client, message):
    args = message.text.split(None, 1)
    if len(args) < 2:
        await message.edit("**Usage:** `.get [name]`")
        return
    
    name = args[1].lower()
    chat_id = str(message.chat.id)
    
    if chat_id not in NOTES or name not in NOTES[chat_id]:
        await message.edit(f"**Note `{name}` not found.**")
        return
    
    note = NOTES[chat_id][name]
    
    if note.get("file_id"):
        media_type = note.get("media_type", "document")
        send_func = getattr(client, f"send_{media_type}", client.send_document)
        await message.delete()
        await send_func(message.chat.id, note["file_id"], caption=note.get("text"))
    else:
        await message.edit(note["text"])

@Client.on_message(filters.command("notes", prefixes=".") & filters.me)
async def list_notes(client, message):
    chat_id = str(message.chat.id)
    
    if chat_id not in NOTES or not NOTES[chat_id]:
        await message.edit("**No notes saved in this chat.**")
        return
    
    note_list = "\n".join([f"• `{name}`" for name in NOTES[chat_id].keys()])
    await message.edit(f"**📝 Saved Notes:**\n{note_list}")

@Client.on_message(filters.command("clear", prefixes=".") & filters.me)
async def clear_note(client, message):
    args = message.text.split(None, 1)
    if len(args) < 2:
        await message.edit("**Usage:** `.clear [name]`")
        return
    
    name = args[1].lower()
    chat_id = str(message.chat.id)
    
    if chat_id not in NOTES or name not in NOTES[chat_id]:
        await message.edit(f"**Note `{name}` not found.**")
        return
    
    del NOTES[chat_id][name]
    await message.edit(f"**✅ Note `{name}` deleted.**")

@Client.on_message(filters.command("clearall", prefixes=".") & filters.me)
async def clear_all_notes(client, message):
    chat_id = str(message.chat.id)
    
    if chat_id in NOTES:
        NOTES[chat_id] = {}
    
    await message.edit("**✅ All notes cleared.**")
