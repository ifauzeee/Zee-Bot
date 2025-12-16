from pyrogram import filters
from functools import wraps

PREFIX = "."

def zee_cmd(command):
    def decorator(func):
        @wraps(func)
        async def wrapper(client, message):
            try:
                return await func(client, message)
            except Exception as e:
                await message.edit(f"**Error:** `{str(e)}`")
        return wrapper
    return decorator

def command(cmd):
    return filters.command(cmd, prefixes=PREFIX) & filters.me
