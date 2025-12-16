from pyrogram import Client, filters
from strings.en import HELP_TEXT, ABOUT_TEXT

@Client.on_message(filters.command("help", prefixes=".") & filters.me)
async def help_cmd(client, message):
    await message.edit(HELP_TEXT, disable_web_page_preview=True)

@Client.on_message(filters.command("about", prefixes=".") & filters.me)
async def about_cmd(client, message):
    await message.edit(ABOUT_TEXT, disable_web_page_preview=True)
