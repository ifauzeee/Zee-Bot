import sys
import traceback
from io import StringIO
from pyrogram import Client, filters
import asyncio

@Client.on_message(filters.command("eval", prefixes=".") & filters.me)
async def eval_cmd(client, message):
    args = message.text.split(None, 1)
    if len(args) < 2:
        await message.edit("**Usage:** `.eval [code]`")
        return
    
    code = args[1]
    old_stdout = sys.stdout
    redirected_output = sys.stdout = StringIO()
    
    try:
        exec_globals = {
            "client": client,
            "message": message,
            "chat": message.chat,
            "asyncio": asyncio,
        }
        
        if "await " in code:
            exec(
                f"async def __aexec__(client, message):\n" +
                "    " + "\n    ".join(code.split("\n")),
                exec_globals
            )
            output = await exec_globals["__aexec__"](client, message)
        else:
            output = eval(code, exec_globals)
    except Exception:
        output = traceback.format_exc()
    finally:
        sys.stdout = old_stdout
    
    printed = redirected_output.getvalue()
    
    if output is not None:
        result = str(output)
    elif printed:
        result = printed
    else:
        result = "No output"
    
    if len(result) > 4000:
        result = result[:4000] + "..."
    
    await message.edit(f"**📤 Output:**\n```python\n{result}\n```")

@Client.on_message(filters.command("sh", prefixes=".") & filters.me)
async def shell_cmd(client, message):
    args = message.text.split(None, 1)
    if len(args) < 2:
        await message.edit("**Usage:** `.sh [command]`")
        return
    
    cmd = args[1]
    await message.edit("**⏳ Executing...**")
    
    try:
        proc = await asyncio.create_subprocess_shell(
            cmd,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE
        )
        stdout, stderr = await proc.communicate()
        
        output = stdout.decode() + stderr.decode()
        if not output:
            output = "No output"
        
        if len(output) > 4000:
            output = output[:4000] + "..."
        
        await message.edit(f"**📤 Output:**\n```\n{output}\n```")
    except Exception as e:
        await message.edit(f"**Error:** `{e}`")
