import os
from pyrogram import Client, filters
from pyrogram.types import InputMediaPhoto
from helpers import fast_status, UI
from PIL import Image
from pyrogram.raw.functions.stickers import CreateStickerSet, AddStickerToSet
from pyrogram.raw.types import InputStickerSetItem, InputStickerSetShortName, InputDocument
from pyrogram.file_id import FileId

@Client.on_message(filters.command("sticker", prefixes=".") & filters.me)
async def sticker_cmd(client, message):
    if not message.reply_to_message:
        await fast_status(message, UI.WARN, "Missing Argument", "Reply to an image to convert.")
        return
    
    reply = message.reply_to_message
    
    if not reply.photo and not reply.document:
        await fast_status(message, UI.WARN, "Invalid Media", "Reply to a photo or document.")
        return
    
    await fast_status(message, UI.LOAD, "Converting Sticker...")
    
    try:
        file_path = await reply.download()
        await message.delete()
        await client.send_sticker(message.chat.id, file_path)
        os.remove(file_path)
    except Exception as e:
        await fast_status(message, UI.FAIL, "Conversion Error", str(e))

@Client.on_message(filters.command("img", prefixes=".") & filters.me)
async def img_cmd(client, message):
    if not message.reply_to_message:
        await fast_status(message, UI.WARN, "Missing Argument", "Reply to a sticker.")
        return
    
    reply = message.reply_to_message
    
    if not reply.sticker:
        await fast_status(message, UI.WARN, "Invalid Media", "Reply to a sticker.")
        return
    
    if reply.sticker.is_animated or reply.sticker.is_video:
        await fast_status(message, UI.WARN, "Unsupported Format", "Animated/Video stickers not supported.")
        return
    
    await fast_status(message, UI.LOAD, "Converting Image...")
    
    try:
        file_path = await reply.download()
        await message.delete()
        await client.send_photo(message.chat.id, file_path)
        os.remove(file_path)
    except Exception as e:
        await fast_status(message, UI.FAIL, "Conversion Error", str(e))

@Client.on_message(filters.command("resize", prefixes=".") & filters.me)
async def resize_cmd(client, message):
    if not message.reply_to_message or not message.reply_to_message.photo:
        await fast_status(message, UI.WARN, "Missing Argument", "Reply to an image to resize.")
        return
    
    await fast_status(message, UI.LOAD, "Resizing Image...")
    
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
        await fast_status(message, UI.FAIL, "Resize Error", str(e))

@Client.on_message(filters.command("kang", prefixes=".") & filters.me)
async def kang_cmd(client, message):
    reply = message.reply_to_message
    if not reply:
        await fast_status(message, UI.WARN, "Missing Argument", "Reply to a sticker or image to kang.")
        return

    await fast_status(message, UI.LOAD, "Kanging Sticker...")

    try:
        if reply.sticker:
            if reply.sticker.is_animated or reply.sticker.is_video:
                 await fast_status(message, UI.WARN, "Unsupported", "Video/Animated stickers not supported yet.")
                 return
            emoji = reply.sticker.emoji or "🤔"
        elif reply.photo:
            emoji = "🤔"
        else:
             await fast_status(message, UI.WARN, "Invalid Media", "Reply to a photo or sticker.")
             return

        if len(message.command) > 1:
            emoji = message.command[1]

        file_path = await reply.download()
        
        from PIL import Image
        img = Image.open(file_path)
        
        max_size = 512
        if img.width != 512 and img.height != 512:
             ratio = min(max_size / img.width, max_size / img.height)
             new_size = (int(img.width * ratio), int(img.height * ratio))
             img = img.resize(new_size, Image.LANCZOS)
        
        kang_path = file_path.rsplit(".", 1)[0] + "_kang.webp"
        img.save(kang_path, "WEBP")

        sent = await client.send_sticker("me", kang_path)
        sticker_file = sent.sticker
        
        media = FileId.decode(sticker_file.file_id)
        input_doc = InputDocument(
            id=media.media_id,
            access_hash=media.access_hash,
            file_reference=media.file_reference
        )
        
        item = InputStickerSetItem(document=input_doc, emoji=emoji)

        user = await client.get_me()
        pack_num = 1
        pack_name = f"zeebot_{user.id}_{pack_num}"
        pack_title = f"ZeeBot Pack {pack_num} (@{user.username})"

        try:
            await client.invoke(
                AddStickerToSet(
                    stickerset=InputStickerSetShortName(short_name=pack_name),
                    sticker=item
                )
            )
            await fast_status(message, UI.DONE, "Sticker Kanged!", f"Added to [Pack {pack_num}](t.me/addstickers/{pack_name})")
        except Exception as e:
             if "STICKERSET_INVALID" in str(e):
                 try:
                     await client.invoke(
                         CreateStickerSet(
                             user_id=await client.resolve_peer(user.id),
                             title=pack_title,
                             short_name=pack_name,
                             stickers=[item],
                         )
                     )
                     await fast_status(message, UI.DONE, "Pack Created!", f"Created [Pack {pack_num}](t.me/addstickers/{pack_name})")
                 except Exception as create_e:
                     await fast_status(message, UI.FAIL, "Kang Error", f"Create failed: {create_e}")
             else:
                 await fast_status(message, UI.FAIL, "Kang Error", f"Add failed: {e}")

        if os.path.exists(file_path): os.remove(file_path)
        if os.path.exists(kang_path): os.remove(kang_path)

    except Exception as e:
        await fast_status(message, UI.FAIL, "Kang Error", str(e))