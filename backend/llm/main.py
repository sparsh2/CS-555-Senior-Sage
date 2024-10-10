import openai
from dotenv import load_dotenv
import os
import time
import sounddevice as sd
from pynput import keyboard
import wave

load_dotenv()
openai.api_key = os.getenv('OPENAI_API_KEY')
client = openai.OpenAI()

def whisper():
    # This function uses OpenAI's whisper voice to text model to convert your voice input to text.
    record_audio()
    audio_file = open("user_response.wav", "rb")
    transcript = client.audio.transcriptions.create(
        model="whisper-1",
        file=audio_file
    )
    return transcript.text

def record_audio(duration=None):
    # This function allows you to record your voice with a press of a button, right now set to 'page down'. You could
    # also bypass the keyboard input logic to consistently talk to the AI without pressing a button.
    CHUNK = 1024
    FORMAT = 'int16'
    CHANNELS = 1
    RATE = 10000
    WAVE_OUTPUT_FILENAME = "user_response.wav"

    frames = []
    stream = None
    is_recording = False
    recording_stopped = False

    def record_audio():
        nonlocal frames, stream
        frames = []

        stream = sd.InputStream(
            samplerate=RATE,
            channels=CHANNELS,
            dtype=FORMAT,
            blocksize=CHUNK,
            callback=callback
        )

        stream.start()

    def callback(indata, frame_count, time, status):
        nonlocal stream
        if is_recording:
            frames.append(indata.copy())

    def stop_recording():
        nonlocal frames, stream, recording_stopped

        stream.stop()
        stream.close()


        wf = wave.open(WAVE_OUTPUT_FILENAME, 'wb')
        wf.setnchannels(CHANNELS)
        wf.setsampwidth(2)
        wf.setframerate(RATE)
        wf.writeframes(b''.join(frames))
        wf.close()
        recording_stopped = True

    def on_key(key):
        nonlocal is_recording

        if key == keyboard.Key.up:
            if not is_recording:
                record_audio()
                is_recording = True
            else:
                stop_recording()
                is_recording = False

    listener = keyboard.Listener(on_press=on_key)
    listener.start()

    start_time = time.time()
    while listener.running:
        if recording_stopped:
            listener.stop()
        elif duration and (time.time() - start_time) > duration:
            listener.stop()
        time.sleep(0.01)

system_prompt = '''You are a helpful assistant for the elderly, try and have conversations with them.
Be polite with them and ask them questions about their day. 
Do not output long answers as it may be too long for them to read. Keep it short and simple and human-like. '''

def openai_complete(user_ip):
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
        print(ans)

    except Exception as e:
        print(f"An error occurred: {e}")

def main_func():
    user_message = whisper()
    i = 0
    while (user_message.lower() != 'exit' or i < 5):
        print(f"You: {user_message.lower()}")
        openai_complete(user_message)
        user_message = whisper()
        i+=1

if __name__ == "__main__":
    main_func()