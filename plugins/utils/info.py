from pyrogram import Client, filters
from pyrogram.enums import ChatType
from datetime import datetime

@Client.on_message(filters.command("info", prefixes=".") & filters.me)
async def info_cmd(client, message):
    if message.reply_to_message:
        user = message.reply_to_message.from_user
    else:
        args = message.text.split(None, 1)
        if len(args) < 2:
            user = message.from_user
        else:
            try:
                user = await client.get_users(args[1])
            except Exception:
                await message.edit("**User not found.**")
                return
    
    if not user:
        await message.edit("**User not found.**")
        return
    
    common = await client.get_common_chats(user.id)
    
    text = f"""
**👤 User Information**

**ID:** `{user.id}`
**First Name:** {user.first_name}
**Last Name:** {user.last_name or 'None'}
**Username:** @{user.username or 'None'}
**Is Bot:** {user.is_bot}
**Is Premium:** {user.is_premium or False}
**Is Verified:** {user.is_verified or False}
**Is Scam:** {user.is_scam or False}
**DC ID:** {user.dc_id or 'Unknown'}
**Common Groups:** {len(common)}
**Profile Link:** [Click Here](tg://user?id={user.id})
"""
    await message.edit(text)

@Client.on_message(filters.command("id", prefixes=".") & filters.me)
async def id_cmd(client, message):
    if message.reply_to_message:
        user = message.reply_to_message.from_user
        if user:
            await message.edit(f"**User ID:** `{user.id}`\n**Chat ID:** `{message.chat.id}`")
        else:
            await message.edit(f"**Chat ID:** `{message.chat.id}`")
    else:
        await message.edit(f"**Your ID:** `{message.from_user.id}`\n**Chat ID:** `{message.chat.id}`")

@Client.on_message(filters.command("chatinfo", prefixes=".") & filters.me)
async def chatinfo_cmd(client, message):
    chat = message.chat
    
    if chat.type == ChatType.PRIVATE:
        await message.edit("**This is a private chat.**")
        return
    
    try:
        full_chat = await client.get_chat(chat.id)
    except Exception:
        full_chat = chat
    
    text = f"""
**💬 Chat Information**

**ID:** `{chat.id}`
**Title:** {chat.title}
**Type:** {chat.type.name}
**Username:** @{chat.username or 'None'}
**Members:** {getattr(full_chat, 'members_count', 'Unknown')}
**Description:** {getattr(full_chat, 'description', 'None') or 'None'}
**DC ID:** {getattr(full_chat, 'dc_id', 'Unknown')}
"""
    await message.edit(text)

@Client.on_message(filters.command("json", prefixes=".") & filters.me)
async def json_cmd(client, message):
    if message.reply_to_message:
        target = message.reply_to_message
    else:
        target = message
    
    json_text = str(target)
    if len(json_text) > 4000:
        json_text = json_text[:4000] + "..."
    
    await message.edit(f"```json\n{json_text}\n```")
