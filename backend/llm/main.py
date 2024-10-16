import json
import os
from datetime import datetime, timedelta
from voice_interactions import stt_whisper, detect_wake_word
from chat_completion import openai_complete

# File and directory configurations
USER_INFO_FILE = 'user_info.json'
LOGS_DIR = 'logs'

# Ensure the logs directory exists
os.makedirs(LOGS_DIR, exist_ok=True)

def select_voice():
    """Prompt the user to select a voice and return the selected voice."""
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

def search_logs_by_date(username, date):
    """Search the user's logs for conversations on a specific date."""
    logs = load_user_logs(username)
    results = []
    for session in logs:
        session_date = datetime.fromisoformat(session['timestamp']).date()
        if session_date == date:
            results.append(session)
    return results

def summarize_conversation(conversation):
    """Generate a summary of a conversation session."""
    summary = []
    for msg in conversation['messages']:
        user_message = msg['user_message']
        bot_response = msg['bot_response']
        summary.append(f"You: {user_message}\nBot: {bot_response}")
    return '\n'.join(summary)

def answer_past_conversation_query(username, query):
    """Handle questions related to past conversations."""
    today = datetime.now().date()
    if 'yesterday' in query:
        date_to_search = today - timedelta(days=1)
    elif 'last conversation' in query or 'previous conversation' in query:
        date_to_search = None
    else:
        # If the query doesn't specify a time, handle differently or return "I don't understand"
        return "I can only answer questions about past conversations like 'yesterday' or 'last session'."

    if date_to_search:
        # Search for conversations from yesterday
        conversations = search_logs_by_date(username, date_to_search)
    else:
        # Get the most recent conversation if 'last conversation' is asked
        conversations = load_user_logs(username)[-1:]  # Get only the last conversation

    if not conversations:
        return "I couldn't find any conversations from that time."

    # Summarize the conversations found
    summaries = [summarize_conversation(conv) for conv in conversations]
    return '\n\n'.join(summaries)

def main_func():
    user_info = load_user_info()
    
    print("Login:\n")
    name = input("Enter your name to continue your conversation:\t").strip()
    
    if not name:
        print("Name cannot be empty. Exiting.")
        return
    
    if name not in user_info:
        print(f"Welcome to the VA, {name}!")
        voice = select_voice()
        user_info[name] = voice
        save_user_info(user_info)
    else:
        voice = user_info[name]
        print(f"Welcome back, {name}! Your selected voice is: {voice}")
    
    # Load previous logs to build the initial context
    past_logs = load_user_logs(name)
    context = []
    for session in past_logs:
        for entry in session.get('messages', []):
            context.append((entry['user_message'], entry['bot_response']))
    
    #print("\nYou can start your conversation. Say 'exit' to end.")
    # Wait for wake word before starting the conversation
    if detect_wake_word():
        print("Wake word detected. You can start your conversation. Say 'exit' to end.")
    else:
        print("Wake word not detected. Exiting.")
        return
    
    # Initialize current conversation session
    current_conversation = {
        'timestamp': datetime.now().isoformat(),
        'messages': []
    }
    
    MAX_CONVERSATIONS = 100  # Adjust as needed to limit log size
    
    while True:
        user_message = stt_whisper().strip()
        
        if user_message.lower() == 'exit':
            print("Ending conversation session.")
            break
        
        # If the user asks about past conversations, handle that query
        # Otherwise, proceed with normal chatbot interaction
        print(f"You: {user_message}")
        bot_response = openai_complete(user_message, context, voice)
        print(f"Bot: {bot_response}\n")
        
        # Update context for the current session
        context.append((user_message, bot_response))
        
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
