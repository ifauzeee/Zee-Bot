import aiohttp
import os
from pyrogram import Client, filters

@Client.on_message(filters.command("carbon", prefixes=".") & filters.me)
async def carbon_cmd(client, message):
    args = message.text.split(None, 1)
    
    if message.reply_to_message:
        code = message.reply_to_message.text
    elif len(args) > 1:
        code = args[1]
    else:
        await message.edit("**Usage:** `.carbon [code]` or reply to code")
        return
    
    if not code:
        await message.edit("**No code found.**")
        return
    
    await message.edit("**⏳ Generating carbon image...**")
    
    try:
        url = "https://carbonara.solopov.dev/api/cook"
        payload = {
            "code": code,
            "theme": "dracula",
            "backgroundColor": "#1a1a2e"
        }
        
        async with aiohttp.ClientSession() as session:
            async with session.post(url, json=payload) as resp:
                if resp.status == 200:
                    image_data = await resp.read()
                    file_path = f"carbon_{message.id}.png"
                    with open(file_path, "wb") as f:
                        f.write(image_data)
                    
                    await message.delete()
                    await client.send_photo(message.chat.id, file_path)
                    os.remove(file_path)
                else:
                    await message.edit("**Failed to generate carbon image.**")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("ss", prefixes=".") & filters.me)
async def screenshot_cmd(client, message):
    args = message.text.split(None, 1)
    
    if len(args) < 2:
        await message.edit("**Usage:** `.ss [url]`")
        return
    
    url = args[1]
    if not url.startswith("http"):
        url = "https://" + url
    
    await message.edit("**⏳ Taking screenshot...**")
    
    try:
        api_url = f"https://image.thum.io/get/width/1280/crop/720/noanimate/{url}"
        
        async with aiohttp.ClientSession() as session:
            async with session.get(api_url) as resp:
                if resp.status == 200:
                    image_data = await resp.read()
                    file_path = f"ss_{message.id}.png"
                    with open(file_path, "wb") as f:
                        f.write(image_data)
                    
                    await message.delete()
                    await client.send_photo(message.chat.id, file_path, caption=f"**Screenshot of:** {url}")
                    os.remove(file_path)
                else:
                    await message.edit("**Failed to take screenshot.**")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")
