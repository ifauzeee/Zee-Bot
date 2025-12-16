import asyncio
from pyrogram import Client

API_ID = int(input("Masukkan API ID: "))
API_HASH = input("Masukkan API HASH: ")

async def main():
    app = Client(
        "gen_session",
        api_id=API_ID,
        api_hash=API_HASH,
        in_memory=True
    )
    
    async with app:
        print("\nSedang memproses...")
        session_str = await app.export_session_string()
        print("\n\n✅ BERHASIL! COPY KODE DI BAWAH INI KE .ENV ANDA:\n")
        print("=" * 50)
        print(session_str)
        print("=" * 50)
        print("\n")

if __name__ == "__main__":
    asyncio.run(main())
