import os
import json
from pyrogram import Client, filters
from helpers import fast_status, UI
from helpers.mongo import db

PMPERMIT_ON = True
WARNS = {}

@Client.on_message(filters.private & ~filters.me & ~filters.bot & ~filters.service, group=99)
async def pm_guard(client, message):
    global PMPERMIT_ON
    
    if not PMPERMIT_ON:
        return
        
    if not message.from_user:
        return

    user_id = message.from_user.id
    
    if await db.is_approved(user_id):
        return

    if user_id not in WARNS:
        WARNS[user_id] = 0
        
    WARNS[user_id] += 1
    count = WARNS[user_id]
    
    if count >= 5:
        await message.reply_text(
            f"**{UI.FAIL} BLOCKED**\n"
            f"{UI.DIVIDER}\n"
            f"Anda terdeteksi melakukan spam.\n"
            f"Sistem otomatis memblokir anda."
        )
        try:
            await client.block_user(user_id)
        except:
            pass
        if user_id in WARNS:
            del WARNS[user_id]
        return

    await message.reply_text(
        f"**{UI.WARN} SECURITY SYSTEM**\n"
        f"{UI.DIVIDER}\n"
        f"Halo **{message.from_user.first_name}**! 👋\n"
        f"Tuanku sedang sibuk/AFK saat ini.\n\n"
        f"⚠️ **Harap Menunggu Konfirmasi!**\n"
        f"Jangan mengirim pesan beruntun atau anda akan diblokir.\n\n"
        f"🛑 **Peringatan:** `{count}/5`"
    )

@Client.on_message(filters.command("pmpermit", prefixes=".") & filters.me)
async def toggle_pm(client, message):
    global PMPERMIT_ON
    args = message.command
    
    if len(args) > 1:
        if args[1].lower() == "on":
            PMPERMIT_ON = True
        elif args[1].lower() == "off":
            PMPERMIT_ON = False
    
    status = "✅ ON" if PMPERMIT_ON else "❌ OFF"
    await fast_status(message, UI.INFO, "PM Security", f"Status: {status}")

@Client.on_message(filters.command("a", prefixes=".") & filters.me)
async def approve_user(client, message):
    if message.chat.type.value != "private":
        await fast_status(message, UI.WARN, "Error", "Command only for Private Chat")
        return
        
    user_id = message.chat.id
    
    await db.approve_user(user_id)
    
    if user_id in WARNS:
        del WARNS[user_id]
    
    await fast_status(message, UI.DONE, "User Approved", "Datanya telah tersimpan di Database.")

@Client.on_message(filters.command("da", prefixes=".") & filters.me)
async def disapprove_user(client, message):
    if message.chat.type.value != "private":
         return
         
    user_id = message.chat.id
    
    await db.disapprove_user(user_id)
    
    await fast_status(message, UI.WARN, "User Disapproved", "Izin chat dicabut dari Database.")