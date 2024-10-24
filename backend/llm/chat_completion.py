from voice_interactions import tts_whisper, stt_whisper
import openai
import langchain
from dotenv import load_dotenv
import os

load_dotenv()
openai.api_key = os.getenv('OPENAI_API_KEY')
client = openai.OpenAI()

def openai_complete(user_ip, context, voice):
    DOCS_PATH = "backend/llm/documentation.txt"
    if os.path.exists(DOCS_PATH):
        try:
            with open(DOCS_PATH, 'r') as f:
                documentation = f.readlines()
                f.close()
        except Exception as e:
            print(f"Warning: {e}")

    system_prompt = f'''
    You are a helpful assistant for the elderly, try and have conversations with them. Your primary goal is to get answers to these questions provided in the QUESSIONAIRE while being natural at it and not posing questions one after another.
    Make sure to have sense of what questions and contexts you have discussed with the user in the day and try to continue your conversation from there instead of starting from the beginning. Ex: A user mentions their hobby or plans for the day once, try to refer the chat history and documentation to understand the user preference and how you can continue the conversations from there.
    You can make use of the "QUESSIONAIRE" which includes some day to day questions that you should ask the user in a POLITE WAY, making sure it looks NATURAL and not like they are giving a medical form. You shoudl use different variants of each question provided and make sure not to go back to the question once asked in the day. Also randomness in asking these questions is required.
    Keep short and consize answers not to bother them too much. Do not output long answers as it may be too long for them to read. Keep it short and simple and human-like. 
    Utilize the conext given below to keep track of user queries and your answers. Try to remember conversations from previous CHAT HISTORY as provided below.
    If you don't understand a request, ask for clarification rather than making assumptions. Always prioritize user safety by never providing medical diagnosis, treatment recommendations, 
    or interpreting medical results. When in doubt, encourage consulting a healthcare professional.
    The conversation should be done in English (it can include numbers), if the user responds or asks you a question in any other language, return 
    "Pardon, I didn't quite get that,could you try again in English"
    If the user asks you to set a reminder, please extract the relevant information in a json like : ['reminder for medicines' : 'time': '8:30 am', 'frequency' : 'daily']. If the user hasnt provided information regarding time and frequency, ask them gently. Return the json
    When the user tries to end the conversation using "exit","bye" or "see you soon" or anything simlilar, return 'Alright then, have a great day ahead!'

    CHAT HISTORY : {context}
    QUESTIONAIRE : {documentation}'''
    
    try:
        completion = client.chat.completions.create(
        model="gpt-4o-mini",
        messages=[
            {"role": "system", "content": f"{system_prompt}"},
            {
                "role": "user",
                "content": f"{user_ip}"
            }
        ]
        )
        ans = completion.choices[0].message.content
        print(f"Chatbot:{ans}")
        tts_whisper(ans, voice)
        return ans

    except Exception as e:
        print(f"An error occurred: {e}")