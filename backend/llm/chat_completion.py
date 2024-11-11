from voice_interactions import tts_whisper, stt_whisper
import openai
import json
from dotenv import load_dotenv
from function_calling import responses
import os
from update_health_question_counter_data import create_questions_to_ask_stack  

load_dotenv()
openai.api_key = os.getenv('OPENAI_API_KEY')
client = openai.OpenAI()
BASE_DIR = os.path.dirname(os.path.abspath(__file__))

USER_HEALTH_LOG_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'user_health_question_counter_logs')
os.makedirs(USER_HEALTH_LOG_DIR, exist_ok=True)

def load_health_questions_to_ask(username):
    questions_file = os.path.join(USER_HEALTH_LOG_DIR, f"{username}_questions_to_ask_stack.json")
    try:
        with open(questions_file, 'r', encoding='utf-8') as f:
            return json.load(f)
    except (FileNotFoundError, json.JSONDecodeError) as e:
        print(f"Error loading health questions: {e}")
        return {}

def openai_complete(username, user_ip, context, voice):
    tools = [
        {
            "type": "function",
            "function": {
                "name": "responses",
                "description": "Save user responses to health-related questions",
                "parameters": {
                    "type": "object",
                    "properties": {
                        "q_idx": {
                            "type": "integer",
                            "description": "Index of the question from health_questions.json"
                        },
                        "user": {
                            "type": "string",
                            "description": "The user's name"
                        },
                        "user_answer": {
                            "type": "string",
                            "description": "User's response to the health question"
                        }
                    },
                    "required": ["q_idx", "user", "user_answer"]
                }
            }
        }
    ]

    # Load health questions
    health_questions = load_health_questions_to_ask(username)
    
    system_prompt = f'''
    You are a helpful assistant for the elderly, try and have conversations with {username}. Your primary goal is to get answers to these questions provided in the QUESSIONAIRE while being natural at it and not posing questions one after another.
    Make sure to have sense of what questions and contexts you have discussed with the user in the day and try to continue your conversation from there instead of starting from the beginning. Ex: A user mentions their hobby or plans for the day once, try to refer the chat history and documentation to understand the user preference and how you can continue the conversations from there.
    You can make use of the "QUESSIONAIRE" which includes some day to day questions that you should ask the user in a POLITE WAY, making sure it looks NATURAL and not like they are giving a medical form. When the user answers these question, call the ``responses`` function. Also randomness in asking these questions is required.
    Keep short and consize answers not to bother them too much. Do not output long answers as it may be too long for them to read. Keep it short and simple and human-like. 
    Utilize the conext given below to keep track of user queries and your answers. Try to bring up events from past conversations using CHAT HISTORY to make it a more personalised experience for the user.
    If you don't understand a request, ask for clarification rather than making assumptions. Always prioritize user safety by never providing medical diagnosis, treatment recommendations, 
    or interpreting medical results. When in doubt, encourage consulting a healthcare professional.
    The conversation should be done in English (it can include numbers), if the user responds or asks you a question in any other language, return 
    "Pardon, I didn't quite get that,could you try again in English"
    If the user asks you to set a reminder, please extract the relevant information in a json like : {{'reminder for medicines' : {{'time': '8:30 am', 'frequency' : 'daily'}}}}. If the user hasnt provided information regarding time and frequency, ask them gently. Return the json
    When the user tries to end the conversation using "exit","bye" or "see you soon" or anything simlilar, return 'Alright then, have a great day ahead!'

    CHAT HISTORY: {context}
    QUESTIONAIRE: {json.dumps(health_questions, indent=2)}'''

    try:
        completion = client.chat.completions.create(
            model="gpt-4o-mini",
            messages=[
                {"role": "system", "content": system_prompt},
                {"role": "user", "content": user_ip}
            ],
            tools=tools,
            tool_choice="auto"
        )
        
        response = completion.choices[0].message
        print(response)
        
        # Check if the model wants to call a function
        if hasattr(response, 'tool_calls') and response.tool_calls:
            tool_call = response.tool_calls[0]
            if tool_call.function.name == "responses":
                try:
                    # Parse and validate function arguments
                    function_args = json.loads(tool_call.function.arguments)
                    
                    # Call the responses function
                    responses(
                        q_idx=function_args["q_idx"],
                        user=function_args["user"],
                        user_answer=function_args["user_answer"]
                    )
                    
                except Exception as e:
                    print("err in inner except")
                    error_msg = f"Error processing health response: {str(e)}"
                    print(error_msg)
                    raise
        
        # If no function call, return the regular response
        regular_response = response.content
        if regular_response == None:
            follow_up_completion = client.chat.completions.create(
                        model="gpt-4-turbo-preview",
                        messages=[
                            {"role": "system", "content": system_prompt},
                            {"role": "user", "content": user_ip},
                            {"role": "assistant", "content": "I've noted your response. Let me follow up on that."},
                            {"role": "user", "content": "Please acknowledge my previous response and continue our conversation naturally."}
                        ]
                    )
                    
            regular_response = follow_up_completion.choices[0].message.content
        print(f"Chatbot: {regular_response}")
        tts_whisper(regular_response, voice)
        return regular_response

    except Exception as e:
        print("error in outer except")
        error_msg = f"An error occurred: {str(e)}"
        print(error_msg)
        tts_whisper("I'm sorry, I encountered an error while processing your request.", voice)
        return error_msg