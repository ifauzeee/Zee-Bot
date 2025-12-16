import aiohttp
from pyrogram import Client, filters

@Client.on_message(filters.command("cur", prefixes=".") & filters.me)
async def currency_cmd(client, message):
    args = message.text.split()
    
    if len(args) < 4:
        await message.edit("**Usage:** `.cur [amount] [from] [to]`\n**Example:** `.cur 100 USD IDR`")
        return
    
    try:
        amount = float(args[1])
        from_cur = args[2].upper()
        to_cur = args[3].upper()
    except ValueError:
        await message.edit("**Invalid amount.**")
        return
    
    await message.edit("**⏳ Converting...**")
    
    try:
        url = f"https://api.exchangerate-api.com/v4/latest/{from_cur}"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as resp:
                data = await resp.json()
        
        if to_cur not in data["rates"]:
            await message.edit(f"**Currency `{to_cur}` not found.**")
            return
        
        rate = data["rates"][to_cur]
        result = amount * rate
        
        await message.edit(
            f"**💱 Currency Conversion**\n\n"
            f"**{amount:,.2f} {from_cur}** = **{result:,.2f} {to_cur}**\n"
            f"**Rate:** 1 {from_cur} = {rate:,.4f} {to_cur}"
        )
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("weather", prefixes=".") & filters.me)
async def weather_cmd(client, message):
    args = message.text.split(None, 1)
    
    if len(args) < 2:
        await message.edit("**Usage:** `.weather [city]`")
        return
    
    city = args[1]
    await message.edit("**⏳ Getting weather...**")
    
    try:
        url = f"https://wttr.in/{city}?format=j1"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as resp:
                data = await resp.json()
        
        current = data["current_condition"][0]
        area = data["nearest_area"][0]
        
        await message.edit(
            f"**🌤️ Weather in {area['areaName'][0]['value']}, {area['country'][0]['value']}**\n\n"
            f"**Temperature:** {current['temp_C']}°C / {current['temp_F']}°F\n"
            f"**Feels Like:** {current['FeelsLikeC']}°C\n"
            f"**Condition:** {current['weatherDesc'][0]['value']}\n"
            f"**Humidity:** {current['humidity']}%\n"
            f"**Wind:** {current['windspeedKmph']} km/h\n"
            f"**UV Index:** {current['uvIndex']}"
        )
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")

@Client.on_message(filters.command("define", prefixes=".") & filters.me)
async def define_cmd(client, message):
    args = message.text.split(None, 1)
    
    if len(args) < 2:
        await message.edit("**Usage:** `.define [word]`")
        return
    
    word = args[1]
    await message.edit("**⏳ Looking up...**")
    
    try:
        url = f"https://api.dictionaryapi.dev/api/v2/entries/en/{word}"
        async with aiohttp.ClientSession() as session:
            async with session.get(url) as resp:
                if resp.status != 200:
                    await message.edit(f"**Word `{word}` not found.**")
                    return
                data = await resp.json()
        
        entry = data[0]
        meanings = entry["meanings"][0]
        definition = meanings["definitions"][0]["definition"]
        example = meanings["definitions"][0].get("example", "No example")
        
        await message.edit(
            f"**📖 Definition of `{word}`**\n\n"
            f"**Part of Speech:** {meanings['partOfSpeech']}\n"
            f"**Definition:** {definition}\n"
            f"**Example:** _{example}_"
        )
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")
