import time
from pyrogram import Client, filters

@Client.on_message(filters.command("alive", prefixes=".") & filters.me)
async def alive_check(client, message):
    start = time.time()
    msg = await message.edit("...")
    end = time.time()
    ping = round((end - start) * 1000, 3)
    
    text = (
        f"**SYSTEM ONLINE**\n"
        f"Ping: `{ping}ms`\n"
        f"User: {message.from_user.mention}\n"
        f"Python: `3.10`\n"
        f"Mode: `Docker/Async`"
    )
    await msg.edit(text)
