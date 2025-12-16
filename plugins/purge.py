import asyncio
from pyrogram import Client, filters

@Client.on_message(filters.command("purge", prefixes=".") & filters.me)
async def purge(client, message):
    if not message.reply_to_message:
        await message.edit("Reply to a message.")
        return

    start_id = message.reply_to_message.id
    end_id = message.id
    chat_id = message.chat.id
    
    msgs = list(range(start_id, end_id + 1))
    
    await message.delete()
    
    try:
        await client.delete_messages(chat_id, msgs)
        status = await client.send_message(chat_id, f"Deleted {len(msgs)} messages.")
        await asyncio.sleep(2)
        await status.delete()
    except Exception:
        pass
