from voice_interactions import stt_whisper
from chat_completion import openai_complete

user_info = {}

def main_func():
    print("Login:\n")
    name = input("Enter your name to continue your conversation:\t")

    if name not in user_info:
        print(f"Welcome to the VA, {name}! Please choose a voice to converse with:")
        ip = input("1. Alloy\n2. Echo\n3.Fable\n4. Onyx\n5. Nova\n6. Shimmer\n")
        dic = {"1": "Alloy", "2": "Echo", "3":"Fable", "4": "Onyx", "5": "Nova", "6": "Shimmer"}
        voice = dic[ip]
        user_info[name] = voice
    else:
        print(f"Welcome back {name}, your selected voice is: {user_info[name]}")
        voice = user_info[name]

    context = []
    user_message = stt_whisper()
    i = 0
    while (user_message.lower() != 'exit' or i < 5):
        print(f"You: {user_message.lower()}")
        chat_ans = openai_complete(user_message, context, voice)
        context.append((user_message, chat_ans))
        print(f"\n\n{context}\n\n")
        user_message = stt_whisper()
        i+=1

if __name__ == "__main__":
    main_func()