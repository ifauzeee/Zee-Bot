import time
from pyrogram import Client, filters
from helpers.utils import human_time

AFK_STATUS = {"is_afk": False, "reason": "", "time": 0}

@Client.on_message(filters.command("afk", prefixes=".") & filters.me)
async def afk_cmd(client, message):
    args = message.text.split(None, 1)
    reason = args[1] if len(args) > 1 else "No reason provided."
    
    AFK_STATUS["is_afk"] = True
    AFK_STATUS["reason"] = reason
    AFK_STATUS["time"] = time.time()
    
    await message.edit(f"**💤 AFK Mode Enabled**\n**Reason:** {reason}")

@Client.on_message(filters.command("unafk", prefixes=".") & filters.me)
async def unafk_cmd(client, message):
    if not AFK_STATUS["is_afk"]:
        await message.edit("**You are not AFK.**")
        return
    
    afk_time = human_time(time.time() - AFK_STATUS["time"])
    AFK_STATUS["is_afk"] = False
    AFK_STATUS["reason"] = ""
    AFK_STATUS["time"] = 0
    
    await message.edit(f"**👋 AFK Mode Disabled**\n**Was AFK for:** {afk_time}")

@Client.on_message(filters.private & filters.incoming & ~filters.me & ~filters.bot)
async def afk_reply_handler(client, message):
    if not AFK_STATUS["is_afk"]:
        return
    
    afk_time = human_time(time.time() - AFK_STATUS["time"])
    await message.reply(
        f"**💤 I'm currently AFK**\n"
        f"**Reason:** {AFK_STATUS['reason']}\n"
        f"**Since:** {afk_time} ago"
    )

@Client.on_message(filters.mentioned & ~filters.me)
async def afk_mention_handler(client, message):
    if not AFK_STATUS["is_afk"]:
        return
    
    afk_time = human_time(time.time() - AFK_STATUS["time"])
    await message.reply(
        f"**💤 I'm currently AFK**\n"
        f"**Reason:** {AFK_STATUS['reason']}\n"
        f"**Since:** {afk_time} ago"
    )

@Client.on_message(filters.me & ~filters.command("afk", prefixes="."))
async def auto_unafk(client, message):
    if AFK_STATUS["is_afk"]:
        afk_time = human_time(time.time() - AFK_STATUS["time"])
        AFK_STATUS["is_afk"] = False
        AFK_STATUS["reason"] = ""
        AFK_STATUS["time"] = 0
        
        await message.reply(f"**👋 Welcome back!**\nYou were AFK for {afk_time}")
