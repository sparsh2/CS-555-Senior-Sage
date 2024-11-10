import json
import os
import re
from datetime import datetime, timedelta
from voice_interactions import stt_whisper, tts_whisper
from chat_completion import openai_complete

# Base directory is the folder containing this script file (assuming it's in "llm" folder)
BASE_DIR = os.path.dirname(os.path.abspath(__file__))

# Paths using os.path.join() for OS independence
USER_INFO_FILE = os.path.join(BASE_DIR, 'user_info.json')
LOGS_DIR = os.path.join(BASE_DIR, 'logs')
REMINDER_DIR = os.path.join(BASE_DIR, 'reminder')

# Ensure the logs directory exists
os.makedirs(LOGS_DIR, exist_ok=True)
os.makedirs(REMINDER_DIR, exist_ok=True)

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

def summarize_conversation(conversation):
    """Generate a summary of a conversation session."""
    summary = []
    for msg in conversation['messages']:
        timestamp = msg['timestamp']
        user_message = msg['user_message']
        bot_response = msg['bot_response']
        summary.append(f"Timestamp: {timestamp} - You: {user_message}\nBot: {bot_response}")
    return '\n'.join(summary)

def main_func():
    user_info = load_user_info()
    
    print("Login:\n")
    login_message = "Please enter your name to login and continue your conversation:"
    tts_whisper(login_message)
    name = input("Enter your name to continue your conversation:\t").strip()
    
    if not name:
        name_empty_message = "Name cannot be empty. Exiting."
        tts_whisper(name_empty_message)
        return
    
    if name not in user_info:
        welcome_message = f"Welcome to the VA, {name}!"
        tts_whisper(welcome_message)
        voice = select_voice()
        user_info[name] = voice
        save_user_info(user_info)
    else:
        voice = user_info[name]
        welcome_back_message = f"Welcome back, {name}! Your selected voice is: {voice}"
        tts_whisper(welcome_back_message)
    
    # Load previous logs to build the initial context
    past_logs = load_user_logs(name)
    context = []
    for session in past_logs:
        for entry in session.get('messages', []):
            context.append((entry['timestamp'], entry['user_message'], entry['bot_response']))
    
    print("\nYou can start your conversation. Say 'exit' to end.")
    
    # Initialize current conversation session
    cur_time = datetime.now().isoformat()
    current_conversation = {
        'timestamp': cur_time,
        'messages': []
    } 
    
    while True:
        user_message = stt_whisper().strip()
        
        if user_message.lower() == 'exit':
            print("Ending conversation session.")
            break
        
        print(f"You: {user_message}")
        bot_response = openai_complete(user_message, context, voice)

        json_match = re.search(r"\{[^}]+\}+", bot_response.replace("\n", "").replace(" ", ""), re.DOTALL)

        if json_match:
            json_str = json_match.group(0).strip() 

            try:
                bot_response_dict = json.loads(json_str)  # Parse it to a dictionary
                reminder_file_path = os.path.join("reminder", f"{name}_reminder.json")
                print(f"Debug: Saving reminder to {reminder_file_path}")
                add_reminder(name, bot_response_dict)
    
            except json.JSONDecodeError:
                print("Debug: bot_response is not a valid JSON string")
        
        # Update context for the current session
        cur_time = datetime.now().isoformat()
        context.append((cur_time, user_message, bot_response))

        if "Alright then have a great" in re.sub(r'[^\w\s]', '', bot_response):
            current_conversation['messages'].append({
            'timestamp': datetime.now().isoformat(),
            'user_message': user_message,
            'bot_response': bot_response
            })
            break
        
        # Add to current conversation
        current_conversation['messages'].append({
            'timestamp': datetime.now().isoformat(),
            'user_message': user_message,
            'bot_response': bot_response
        })
    
    # Append the completed conversation to the user's log
    append_conversation(name, current_conversation)
    print(f"Conversation session saved for user '{name}'.")

if __name__ == "__main__":
    main_func()
