import os
from pyrogram import Client, filters
from pyrogram.types import InputMediaPhoto

@Client.on_message(filters.command("sticker", prefixes=".") & filters.me)
async def sticker_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to an image to convert to sticker.**")
        return
    
    reply = message.reply_to_message
    
    if not reply.photo and not reply.document:
        await message.edit("**Reply to an image.**")
        return
    
    await message.edit("**⏳ Converting...**")
    
    try:
        file_path = await reply.download()
        await message.delete()
        await client.send_sticker(message.chat.id, file_path)
        os.remove(file_path)
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("img", prefixes=".") & filters.me)
async def img_cmd(client, message):
    if not message.reply_to_message:
        await message.edit("**Reply to a sticker to convert to image.**")
        return
    
    reply = message.reply_to_message
    
    if not reply.sticker:
        await message.edit("**Reply to a sticker.**")
        return
    
    if reply.sticker.is_animated or reply.sticker.is_video:
        await message.edit("**Animated/Video stickers not supported.**")
        return
    
    await message.edit("**⏳ Converting...**")
    
    try:
        file_path = await reply.download()
        await message.delete()
        await client.send_photo(message.chat.id, file_path)
        os.remove(file_path)
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("resize", prefixes=".") & filters.me)
async def resize_cmd(client, message):
    if not message.reply_to_message or not message.reply_to_message.photo:
        await message.edit("**Reply to an image to resize for sticker.**")
        return
    
    await message.edit("**⏳ Resizing...**")
    
    try:
        file_path = await message.reply_to_message.download()
        
        from PIL import Image
        img = Image.open(file_path)
        
        max_size = 512
        ratio = min(max_size / img.width, max_size / img.height)
        new_size = (int(img.width * ratio), int(img.height * ratio))
        img = img.resize(new_size, Image.LANCZOS)
        
        output_path = file_path.rsplit(".", 1)[0] + "_resized.webp"
        img.save(output_path, "WEBP")
        
        await message.delete()
        await client.send_sticker(message.chat.id, output_path)
        
        os.remove(file_path)
        os.remove(output_path)
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")
