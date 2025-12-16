import asyncio
import os
from dotenv import load_dotenv
from pyrogram import Client

load_dotenv()

env_api_id = os.getenv("API_ID")
env_api_hash = os.getenv("API_HASH")

if env_api_id and env_api_hash:
    print(f"✅ Found API Config in .env: {env_api_id}")
    use_env = input("Gunakan config ini? (y/n): ").lower()
    if use_env == 'y':
        API_ID = int(env_api_id)
        API_HASH = env_api_hash
    else:
        API_ID = int(input("Masukkan API ID: "))
        API_HASH = input("Masukkan API HASH: ")
else:
    API_ID = int(input("Masukkan API ID: "))
    API_HASH = input("Masukkan API HASH: ")

async def main():
    app = Client(
        "gen_session",
        api_id=API_ID,
        api_hash=API_HASH,
        in_memory=True
    )
    
    try:
        await app.start()
        print("\n✅ BERHASIL LOGIN!")
        session_str = await app.export_session_string()
        
        env_path = ".env"
        if os.path.exists(env_path):
            import re
            
            with open(env_path, "r", encoding="utf-8") as f:
                content = f.read()
            
            if "SESSION_STRING=" in content:
                content = re.sub(r"^SESSION_STRING=.*", f"SESSION_STRING={session_str}", content, flags=re.MULTILINE)
            else:
                if not content.endswith("\n"):
                    content += "\n"
                content += f"SESSION_STRING={session_str}\n"
            
            with open(env_path, "w", encoding="utf-8") as f:
                f.write(content)
                
            print(f"\n✨ SUKSES! Session String telah otomatis disimpan ke {env_path}")
            print("Silakan jalankan ulang bot anda: docker compose up -d --build")
        else:
             print("\n⚠️ File .env tidak ditemukan. Silakan copy manual:")
             print(session_str)

        await app.stop()
    except Exception as e:
        print(f"\n❌ Error: {e}")

if __name__ == "__main__":
    asyncio.run(main())
