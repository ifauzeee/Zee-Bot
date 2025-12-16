import asyncio
from pyrogram import Client, filters
from pyrogram.enums import ChatMemberStatus

@Client.on_message(filters.command("ban", prefixes=".") & filters.me)
async def ban_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a user to ban.**")
        return
    
    user = message.reply_to_message.from_user
    try:
        await client.ban_chat_member(message.chat.id, user.id)
        await message.edit(f"**🚫 Banned:** {user.mention}")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("unban", prefixes=".") & filters.me)
async def unban_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a user to unban.**")
        return
    
    user = message.reply_to_message.from_user
    try:
        await client.unban_chat_member(message.chat.id, user.id)
        await message.edit(f"**✅ Unbanned:** {user.mention}")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("kick", prefixes=".") & filters.me)
async def kick_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a user to kick.**")
        return
    
    user = message.reply_to_message.from_user
    try:
        await client.ban_chat_member(message.chat.id, user.id)
        await asyncio.sleep(1)
        await client.unban_chat_member(message.chat.id, user.id)
        await message.edit(f"**👢 Kicked:** {user.mention}")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("mute", prefixes=".") & filters.me)
async def mute_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a user to mute.**")
        return
    
    user = message.reply_to_message.from_user
    try:
        await client.restrict_chat_member(
            message.chat.id, 
            user.id,
            permissions={}
        )
        await message.edit(f"**🔇 Muted:** {user.mention}")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("unmute", prefixes=".") & filters.me)
async def unmute_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a user to unmute.**")
        return
    
    user = message.reply_to_message.from_user
    try:
        await client.restrict_chat_member(
            message.chat.id,
            user.id,
            permissions=message.chat.permissions
        )
        await message.edit(f"**🔊 Unmuted:** {user.mention}")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("pin", prefixes=".") & filters.me)
async def pin_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a message to pin.**")
        return
    
    try:
        await message.reply_to_message.pin()
        await message.edit("**📌 Message pinned.**")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("unpin", prefixes=".") & filters.me)
async def unpin_cmd(client, message):
    try:
        await client.unpin_chat_message(message.chat.id)
        await message.edit("**📌 Message unpinned.**")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")
