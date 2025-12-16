import os
from dotenv import load_dotenv

load_dotenv()

class Config:
    API_ID = os.getenv("API_ID")
    API_HASH = os.getenv("API_HASH")
    SESSION_STRING = os.getenv("SESSION_STRING")
    MONGO_URL = os.getenv("MONGO_URL")

    if not API_ID or not API_HASH or not SESSION_STRING or not MONGO_URL:
        print("❌ CRITICAL ERROR: Missing .env configuration!")
        print("Please ensure API_ID, API_HASH, SESSION_STRING, and MONGO_URL are set.")
        
    if API_ID:
        try:
            API_ID = int(API_ID)
        except ValueError:
             print("❌ ERROR: API_ID must be an integer!")
