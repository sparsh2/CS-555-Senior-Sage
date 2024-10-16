from voice_interactions import tts_whisper, stt_whisper
import openai
import langchain
from dotenv import load_dotenv
import os

load_dotenv()
openai.api_key = os.getenv('OPENAI_API_KEY')
client = openai.OpenAI()

def openai_complete(user_ip, context, voice):
    system_prompt = f'''
    You are a helpful assistant for the elderly, try and have conversations with them.
    Be polite with them and ask them questions about their day. 
    Keep short and consize answers not to bother them too much.
    Do not output long answers as it may be too long for them to read. Keep it short and simple and human-like. 
    Utilize the conext given below to keep track of user queries and your answers. Try to remember conversations from previous CHAT HISTORY as provided below.
    Do not ask repeated questions
    If you don't understand a request, ask for clarification rather than making assumptions. Always prioritize user safety by never providing medical diagnosis, treatment recommendations, or interpreting medical results. When in doubt, encourage consulting a healthcare professional.
    The conversation should be done in English, if the user responds or asks you a question in any other language, return "Pardon, I didn't quite get that,could you try again in English"
    When the user tries to end the conversation using "exit","bye" or "see you soon" or anything simlilar, return 'Alright then, have a great day ahead!'

    CHAT HISTORY : {context}'''
    
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