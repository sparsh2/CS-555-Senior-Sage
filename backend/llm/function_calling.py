import json
import os
from datetime import datetime
from update_health_question_counter_data import *

# Base directory is the folder containing this script file
BASE_DIR = os.path.dirname(os.path.abspath(__file__))

# Paths using os.path.join() for OS independence
RESPONSES_FILE = os.path.join(BASE_DIR, 'curr_response.json')
HEALTH_QUESTIONS_FILE = os.path.join(BASE_DIR, 'health_questions.json')

def responses(q_idx: int, user: str, user_answer: str) -> None:
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
            "user": user,
            "q_idx": q_idx,
            "Questions": health_questions[q_idx_str],
            "Answer": user_answer,
            "Date": datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        })

        # Save the updated responses
        with open(RESPONSES_FILE, 'w', encoding='utf-8') as file:
            json.dump(curr_response, file, indent=4, ensure_ascii=False)
            
        counter_data = load_user_health_question_counter(user)
        update_health_question_counter(user, q_idx_str, counter_data)
        save_user_health_question_counter(user, counter_data)

        print(f"Successfully saved response for question {q_idx}")

    except Exception as e:
        print(f"Error in responses function: {e}")
        raise

    return 