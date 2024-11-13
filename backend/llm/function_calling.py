import json
import os
from helper import *

BASE_DIR = os.path.dirname(os.path.abspath(__file__))
RESPONSES_FILE = os.path.join(BASE_DIR, 'curr_response.json')
HEALTH_QUESTIONS_FILE = os.path.join(BASE_DIR, 'health_questions.json')

def reminders(username:str, remind:dict) -> None:
    add_reminder(username, remind)
    return

def preferences(username: str, preference_type: str, preference_detail: str, sentiment: str):
    """
    Store user preferences in a JSON file. Each user has their own preferences file that maintains
    a history of their stated preferences with timestamps.
    
    Args:
        username (str): The user's name
        preference_type (str): Category of the preference (food, hobby, etc.)
        preference_detail (str): Detailed description of the preference
        sentiment (str): Whether it's a like, dislike, or neutral preference
    """
    # Create preferences directory if it doesn't exist
    preferences_dir = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'user_preferences')
    os.makedirs(preferences_dir, exist_ok=True)
    
    # Path to user's preferences file
    preferences_file = os.path.join(preferences_dir, f"{username.lower()}_preferences.json")
    
    try:
        # Load existing preferences if file exists
        if os.path.exists(preferences_file):
            with open(preferences_file, 'r', encoding='utf-8') as f:
                preferences = json.load(f)
        else:
            preferences = {
                "user": username,
                "last_updated": None,
                "preferences": []
            }
        
        # Create new preference entry
        new_preference = {
            "type": preference_type,
            "detail": preference_detail,
            "sentiment": sentiment,
            "timestamp": datetime.now().isoformat(),
        }
        
        # Add new preference to the list
        preferences["preferences"].append(new_preference)
        preferences["last_updated"] = datetime.now().isoformat()
        
        # Write updated preferences back to file
        with open(preferences_file, 'w', encoding='utf-8') as f:
            json.dump(preferences, f, indent=2, ensure_ascii=False)
            
        print(f"Successfully stored preference for {username}: {preference_detail}")
        return True
        
    except Exception as e:
        print(f"Error storing preference for {username}: {str(e)}")
        raise

def get_user_preferences(username: str) -> dict:
    """
    Retrieve all stored preferences for a given user.
    
    Args:
        username (str): The user's name
        
    Returns:
        dict: Dictionary containing all user preferences
    """
    preferences_dir = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'user_preferences')
    preferences_file = os.path.join(preferences_dir, f"{username.lower()}_preferences.json")
    
    try:
        if os.path.exists(preferences_file):
            with open(preferences_file, 'r', encoding='utf-8') as f:
                return json.load(f)
        return None
    except Exception as e:
        print(f"Error retrieving preferences for {username}: {str(e)}")
        return None

def get_preferences_by_type(username: str, preference_type: str) -> list:
    """
    Retrieve all preferences of a specific type for a given user.
    
    Args:
        username (str): The user's name
        preference_type (str): Type of preferences to retrieve (e.g., 'food', 'hobby')
        
    Returns:
        list: List of preferences of the specified type
    """
    preferences = get_user_preferences(username)
    if preferences and 'preferences' in preferences:
        return [
            pref for pref in preferences['preferences']
            if pref['type'] == preference_type
        ]
    return []

def responses(q_idx: int, username: str, user_answer: str) -> None:
    """
    Save user responses to health-related questions.
    
    Args:
        q_idx (int): Index of the question being answered
        user (str): Name of the user
        user_answer (str): User's response to the question
    """
    print("Processing health response...")
    
    try:
        try:
            with open(RESPONSES_FILE, 'r', encoding='utf-8') as file:
                curr_response = json.load(file)
        except (FileNotFoundError, json.JSONDecodeError):
            curr_response = {
                "user": "",
                "q_idx": 0,
                "Questions": "",
                "Answer": "",
                "Date": ""
            }
        try:
            with open(HEALTH_QUESTIONS_FILE, 'r', encoding='utf-8') as file:
                health_questions = json.load(file)
        except (FileNotFoundError, json.JSONDecodeError) as e:
            print(f"Error loading health questions: {e}")
            return 

        # Convert q_idx to string for dictionary lookup
        q_idx_str = str(q_idx)
        
        # Verify the question index exists
        if q_idx_str not in health_questions:
            print(f"Question index {q_idx} not found in health questions")
            return 

        # Update the response dictionary
        curr_response.update({
            "user": username,
            "q_idx": q_idx,
            "Questions": health_questions[q_idx_str],
            "Answer": user_answer,
            "Date": datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        })

        # Save the updated responses
        with open(RESPONSES_FILE, 'w', encoding='utf-8') as file:
            json.dump(curr_response, file, indent=4, ensure_ascii=False)
            
        counter_data = load_user_health_question_counter(username)
        update_health_question_counter(username, q_idx_str, counter_data)
        save_user_health_question_counter(username, counter_data)

        print(f"Successfully saved response for question {q_idx}")

    except Exception as e:
        print(f"Error in responses function: {e}")
        raise

    return 