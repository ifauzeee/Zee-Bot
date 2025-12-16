from pyrogram import Client, filters

@Client.on_message(filters.command("spam", prefixes=".") & filters.me)
async def spam_cmd(client, message):
    args = message.text.split(None, 2)
    
    if len(args) < 3:
        await message.edit("**Usage:** `.spam [count] [text]`")
        return
    
    try:
        count = int(args[1])
    except ValueError:
        await message.edit("**Invalid count.**")
        return
    
    if count > 50:
        await message.edit("**Maximum 50 messages.**")
        return
    
    text = args[2]
    await message.delete()
    
    for _ in range(count):
        await client.send_message(message.chat.id, text)

@Client.on_message(filters.command("copy", prefixes=".") & filters.me)
async def copy_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a message to copy.**")
        return
    
    await message.delete()
    await message.reply_to_message.copy(message.chat.id)

@Client.on_message(filters.command("forward", prefixes=".") & filters.me)
async def forward_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a message to forward.**")
        return
    
    args = message.text.split(None, 1)
    if len(args) < 2:
        await message.edit("**Usage:** `.forward [chat_id/username]`")
        return
    
    try:
        await message.reply_to_message.forward(args[1])
        await message.edit(f"**✅ Forwarded to {args[1]}**")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("send", prefixes=".") & filters.me)
async def send_cmd(client, message):
    args = message.text.split(None, 2)
    
    if len(args) < 3:
        await message.edit("**Usage:** `.send [chat_id/username] [text]`")
        return
    
    target = args[1]
    text = args[2]
    
    try:
        await client.send_message(target, text)
        await message.edit(f"**✅ Sent to {target}**")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("read", prefixes=".") & filters.me)
async def read_cmd(client, message):
    try:
        await client.read_chat_history(message.chat.id)
        await message.edit("**✅ Marked as read.**")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("leave", prefixes=".") & filters.me)
async def leave_cmd(client, message):
    await message.edit("**👋 Leaving chat...**")
    try:
        await client.leave_chat(message.chat.id)
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")
