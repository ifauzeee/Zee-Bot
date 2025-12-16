import asyncio
import logging
import uvloop
from pyrogram import Client, idle
from config import Config

uvloop.install()

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger("Zee-Bot")

PLUGIN_FOLDERS = [
    "plugins.core",
    "plugins.admin",
    "plugins.utils",
    "plugins.media"
]

app = None

async def main():
    global app
    app = Client(
        "ZeeBot",
        api_id=Config.API_ID,
        api_hash=Config.API_HASH,
        session_string=Config.SESSION_STRING,
        plugins=dict(root="plugins")
    )
    
    await app.start()
    me = await app.get_me()
    
    print("=" * 50)
    print(f"✅ Zee-Bot Started Successfully!")
    print(f"👤 User: {me.first_name}")
    print(f"🆔 ID: {me.id}")
    print(f"📛 Username: @{me.username or 'None'}")
    print("=" * 50)
    print("Type .help in any chat to see commands")
    print("=" * 50)
    
    await idle()
    await app.stop()

if __name__ == "__main__":
    asyncio.run(main())
