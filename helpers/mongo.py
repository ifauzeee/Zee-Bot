from motor.motor_asyncio import AsyncIOMotorClient
from config import Config

class Database:
    def __init__(self):
        self.client = AsyncIOMotorClient(Config.MONGO_URL)
        self.db = self.client["ZeeBot"]
        
        self.notes = self.db.notes
        self.pmpermit = self.db.pmpermit

    async def add_note(self, chat_id, name, content):
        await self.notes.update_one(
            {"chat_id": str(chat_id), "name": name},
            {"$set": {"content": content}},
            upsert=True
        )

    async def get_note(self, chat_id, name):
        return await self.notes.find_one({"chat_id": str(chat_id), "name": name})

    async def delete_note(self, chat_id, name):
        await self.notes.delete_one({"chat_id": str(chat_id), "name": name})

    async def get_all_notes(self, chat_id):
        cursor = self.notes.find({"chat_id": str(chat_id)})
        return await cursor.to_list(length=None)
    
    async def delete_all_notes(self, chat_id):
        await self.notes.delete_many({"chat_id": str(chat_id)})

    async def is_approved(self, user_id):
        doc = await self.pmpermit.find_one({"user_id": int(user_id)})
        return bool(doc)

    async def approve_user(self, user_id):
        await self.pmpermit.update_one(
            {"user_id": int(user_id)},
            {"$set": {"approved": True}},
            upsert=True
        )

    async def disapprove_user(self, user_id):
        await self.pmpermit.delete_one({"user_id": int(user_id)})

db = Database()