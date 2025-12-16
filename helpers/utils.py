import time
from datetime import datetime

def get_text(message):
    text = message.text.split(None, 1)
    if len(text) > 1:
        return text[1]
    if message.reply_to_message:
        return message.reply_to_message.text or message.reply_to_message.caption
    return None

async def edit_or_reply(message, text, **kwargs):
    if message.from_user.is_self:
        return await message.edit(text, **kwargs)
    return await message.reply(text, **kwargs)

def human_time(seconds):
    seconds = int(seconds)
    if seconds < 60:
        return f"{seconds}s"
    minutes = seconds // 60
    seconds = seconds % 60
    if minutes < 60:
        return f"{minutes}m {seconds}s"
    hours = minutes // 60
    minutes = minutes % 60
    if hours < 24:
        return f"{hours}h {minutes}m"
    days = hours // 24
    hours = hours % 24
    return f"{days}d {hours}h"

def get_readable_size(size):
    units = ["B", "KB", "MB", "GB", "TB"]
    unit_index = 0
    while size >= 1024 and unit_index < len(units) - 1:
        size /= 1024
        unit_index += 1
    return f"{size:.2f} {units[unit_index]}"
