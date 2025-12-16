import time
import platform
import pyrogram
from pyrogram import Client, filters
from helpers.utils import human_time

START_TIME = time.time()

@Client.on_message(filters.command("alive", prefixes=".") & filters.me)
async def alive_cmd(client, message):
    uptime = human_time(time.time() - START_TIME)
    start = time.time()
    msg = await message.edit("...")
    ping = round((time.time() - start) * 1000, 2)
    me = await client.get_me()
    
    text = f"""
**🤖 Zee-Bot is Alive!**

**User:** {me.first_name}
**Username:** @{me.username or 'None'}
**Ping:** `{ping}ms`
**Uptime:** `{uptime}`
**Python:** `{platform.python_version()}`
**Pyrogram:** `{pyrogram.__version__}`
**Platform:** `{platform.system()}`
"""
    await msg.edit(text)

@Client.on_message(filters.command("ping", prefixes=".") & filters.me)
async def ping_cmd(client, message):
    start = time.time()
    msg = await message.edit("**Pinging...**")
    ping = round((time.time() - start) * 1000, 2)
    await msg.edit(f"**🏓 Pong!** `{ping}ms`")
