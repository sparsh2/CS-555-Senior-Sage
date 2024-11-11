import json
import os
import re
from datetime import datetime, timedelta
from voice_interactions import stt_whisper, tts_whisper
from chat_completion import openai_complete
from update_health_question_counter_data import update_health_question_counter, save_user_health_question_counter, load_health_questions


# Base directory is the folder containing this script file (assuming it's in "llm" folder)
BASE_DIR = os.path.dirname(os.path.abspath(__file__))

# Paths using os.path.join() for OS independence
USER_INFO_FILE = os.path.join(BASE_DIR, 'user_info.json')
LOGS_DIR = os.path.join(BASE_DIR, 'logs')
REMINDER_DIR = os.path.join(BASE_DIR, 'reminder')
# USER_HEALTH_LOG_DIR = os.path.join(BASE_DIR, 'user_health_question_counter_logs')  # New directory


# Ensure the logs directory exists
os.makedirs(LOGS_DIR, exist_ok=True)
os.makedirs(REMINDER_DIR, exist_ok=True)
# os.makedirs(USER_HEALTH_LOG_DIR, exist_ok=True)  # Ensure the new folder is created



def select_voice():
    """Prompt the user to select a voice and return the selected voice."""
    voice_select_message = "Please select the voice of your choice"
    tts_whisper(voice_select_message)
    print("Please choose a voice to converse with:")
    print("1. Alloy\n2. Echo\n3. Fable\n4. Onyx\n5. Nova\n6. Shimmer")
    dic = {"1": "Alloy", "2": "Echo", "3": "Fable", "4": "Onyx", "5": "Nova", "6": "Shimmer"}
    while True:
        ip = input("Select an option (1-6):\t").strip()
        if ip in dic:
            return dic[ip]
        else:
            print("Invalid selection. Please choose a number between 1 and 6.")

def load_user_info():
    """Load user information from a JSON file."""
    if os.path.exists(USER_INFO_FILE):
        try:
            with open(USER_INFO_FILE, 'r') as f:
                return json.load(f)
        except json.JSONDecodeError:
            print("Error: user_info.json is corrupted. Starting with empty user info.")
    return {}

def save_user_info(user_info):
    """Save user information to a JSON file."""
    try:
        with open(USER_INFO_FILE, 'w') as f:
            json.dump(user_info, f, indent=4)
    except IOError as e:
        print(f"Error saving user info: {e}")

def load_user_reminders(username):
    """Load reminders for a specific user."""
    reminder_file = os.path.join(REMINDER_DIR, f"{username}_reminders.json")
    if os.path.exists(reminder_file):
        try:
            with open(reminder_file, 'r') as f:
                return json.load(f)
        except json.JSONDecodeError:
            print(f"Warning: {reminder_file} is corrupted. Starting with empty reminders.")
    return []

def save_user_reminders(username, reminders):
    """Save reminders for a specific user."""
    reminder_file = os.path.join(REMINDER_DIR, f"{username}_reminders.json")
    try:
        with open(reminder_file, 'w') as f:
            json.dump(reminders, f, indent=4)
    except IOError as e:
        print(f"Error saving reminders for {username}: {e}")

def add_reminder(username, reminder_details):
    """Add a new reminder to the user's reminders list."""
    reminders = load_user_reminders(username)
    reminders.append(reminder_details)
    save_user_reminders(username, reminders)

def load_user_logs(username):
    """Load conversation logs for a specific user."""
    log_file = os.path.join(LOGS_DIR, f"{username}.json")
    if os.path.exists(log_file):
        try:
            with open(log_file, 'r') as f:
                return json.load(f)
        except json.JSONDecodeError:
            print(f"Warning: {log_file} is corrupted. Starting with empty logs.")
    return []

def save_user_logs(username, logs):
    """Save conversation logs for a specific user."""
    log_file = os.path.join(LOGS_DIR, f"{username}.json")
    try:
        with open(log_file, 'w') as f:
            json.dump(logs, f, indent=4)
    except IOError as e:
        print(f"Error saving logs for {username}: {e}")

def append_conversation(username, conversation):
    """Append a completed conversation session to the user's log."""
    logs = load_user_logs(username)
    logs.append(conversation)
    save_user_logs(username, logs)
