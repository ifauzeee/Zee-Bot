import aiohttp
from pyrogram import Client, filters
from helpers.utils import get_text

@Client.on_message(filters.command("tr", prefixes=".") & filters.me)
async def translate_cmd(client, message):
    args = message.text.split(None, 2)
    
    if len(args) < 2:
        await message.edit("**Usage:** `.tr [lang] [text]` or reply to a message with `.tr [lang]`")
        return
    
    lang = args[1]
    
    if message.reply_to_message:
        text = message.reply_to_message.text or message.reply_to_message.caption
    elif len(args) > 2:
        text = args[2]
    else:
        await message.edit("**Provide text to translate.**")
        return
    
    if not text:
        await message.edit("**No text found.**")
        return
    
    await message.edit("**⏳ Translating...**")
    
    try:
        url = f"https://api.mymemory.translated.net/get?q={text}&langpair=auto|{lang}"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as resp:
                data = await resp.json()
        
        translated = data["responseData"]["translatedText"]
        detected = data["responseData"].get("detectedLanguage", "auto")
        
        await message.edit(
            f"**🌐 Translated to {lang.upper()}**\n\n"
            f"**Original:** `{text[:100]}{'...' if len(text) > 100 else ''}`\n\n"
            f"**Translation:**\n{translated}"
        )
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("tts", prefixes=".") & filters.me)
async def tts_cmd(client, message):
    args = message.text.split(None, 1)
    
    if message.reply_to_message:
        text = message.reply_to_message.text or message.reply_to_message.caption
    elif len(args) > 1:
        text = args[1]
    else:
        await message.edit("**Usage:** `.tts [text]` or reply to a message")
        return
    
    if not text:
        await message.edit("**No text found.**")
        return
    
    if len(text) > 200:
        text = text[:200]
    
    await message.edit("**⏳ Generating audio...**")
    
    try:
        from gtts import gTTS
        import os
        
        tts = gTTS(text=text, lang='en')
        file_path = f"tts_{message.id}.mp3"
        tts.save(file_path)
        
        await message.delete()
        await client.send_voice(message.chat.id, file_path)
        os.remove(file_path)
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")
