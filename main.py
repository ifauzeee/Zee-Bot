import asyncio
import logging
import uvloop
from pyrogram import Client, idle
from config import Config

uvloop.install()

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(message)s')
logger = logging.getLogger("Zee-Bot")

app = None

async def main():
    global app
    app = Client(
        "ZeeUbot",
        api_id=Config.API_ID,
        api_hash=Config.API_HASH,
        session_string=Config.SESSION_STRING,
        plugins=dict(root="plugins")
    )
    await app.start()
    me = await app.get_me()
    print(f"✅ Zee-Bot Started as {me.first_name} (@{me.username})")
    await idle()
    await app.stop()

if __name__ == "__main__":
    asyncio.run(main())
