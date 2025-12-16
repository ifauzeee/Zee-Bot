import asyncio
from pyrogram import Client, filters

@Client.on_message(filters.command("purge", prefixes=".") & filters.me)
async def purge_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a message to start purging.**")
        return
    
    start_id = message.reply_to_message.id
    end_id = message.id
    chat_id = message.chat.id
    
    msg_ids = list(range(start_id, end_id + 1))
    deleted = 0
    
    for i in range(0, len(msg_ids), 100):
        batch = msg_ids[i:i+100]
        try:
            await client.delete_messages(chat_id, batch)
            deleted += len(batch)
        except Exception:
            pass
    
    status = await client.send_message(chat_id, f"**🗑️ Purged {deleted} messages.**")
    await asyncio.sleep(2)
    await status.delete()

@Client.on_message(filters.command("del", prefixes=".") & filters.me)
async def del_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a message to delete.**")
        return
    
    try:
        await message.reply_to_message.delete()
        await message.delete()
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("purgeme", prefixes=".") & filters.me)
async def purgeme_cmd(client, message):
    args = message.text.split()
    if len(args) < 2:
        await message.edit("**Usage:** `.purgeme [number]`")
        return
    
    try:
        count = int(args[1])
    except ValueError:
        await message.edit("**Invalid number.**")
        return
    
    deleted = 0
    async for msg in client.get_chat_history(message.chat.id, limit=count + 1):
        if msg.from_user and msg.from_user.is_self:
            try:
                await msg.delete()
                deleted += 1
            except Exception:
                pass
            if deleted >= count:
                break
    
    status = await client.send_message(message.chat.id, f"**🗑️ Deleted {deleted} of my messages.**")
    await asyncio.sleep(2)
    await status.delete()
