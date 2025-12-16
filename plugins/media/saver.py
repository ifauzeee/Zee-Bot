from pyrogram import Client, filters
import os
import asyncio

@Client.on_message(filters.command(["fsave", "savemedia"], prefixes=".") & filters.me)
async def save_media(client, message):
    if not message.reply_to_message:
        await message.edit("**❌ Please reply to a media message (View Once / Timer / Protected Content).**")
        return

    original_msg = message.reply_to_message
    await message.edit("**🔄 Processing media retrieval...**")
    
    try:
        target_msg = await client.get_messages(original_msg.chat.id, original_msg.id)
    except Exception:
        target_msg = original_msg

    if not target_msg.media:
        await message.edit(
            "**❌ Retrieval Failed:** The server returned an empty message payload.\n\n"
            "**Possible Causes:**\n"
            "1. The message was sent to **Self/Saved Messages** (automatically marked as viewed).\n"
            "2. The media has already been **viewed/opened** prior to saving.\n\n"
            "**Recommendation:** Ensure the media remains unopened before attempting to save."
        )
        return

    await message.edit(f"**📥 Fetching {target_msg.media.value.lower().replace('messagemediatype.', '')}...**")

    file_path = None
    try:
        try:
            file_path = await client.download_media(target_msg)
        except ValueError:
            file_id = None
            if target_msg.photo:
                file_id = target_msg.photo.file_id
            elif target_msg.video:
                file_id = target_msg.video.file_id
            elif target_msg.voice:
                file_id = target_msg.voice.file_id
            elif target_msg.video_note:
                file_id = target_msg.video_note.file_id
                
            if file_id:
                file_path = await client.download_media(file_id)
            else:
                raise Exception("Media type invalid or expired.")

        if not file_path:
            raise Exception("Download returned empty path.")

        await message.edit("**📤 Archiving to Saved Messages...**")
        
        caption = f"**Saved from:** {target_msg.chat.title or 'Private Chat'}\n"
        if target_msg.from_user:
            caption += f"**User:** {target_msg.from_user.mention}"

        if target_msg.photo:
            await client.send_photo("me", file_path, caption=caption)
        elif target_msg.video:
            await client.send_video("me", file_path, caption=caption)
        elif target_msg.voice:
            await client.send_voice("me", file_path)
        elif target_msg.video_note:
            await client.send_video_note("me", file_path)
        elif target_msg.audio:
            await client.send_audio("me", file_path)
        else:
            await client.send_document("me", file_path, caption=caption)
            
        await message.edit("**✅ Media successfully archived.**")

    except Exception as e:
        await message.edit(f"**❌ Error:** `{str(e)}`")
        
    finally:
        if file_path and os.path.exists(file_path):
            os.remove(file_path)
