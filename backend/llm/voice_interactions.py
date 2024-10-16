import openai
from dotenv import load_dotenv
import os
import time
import sounddevice as sd
from pynput import keyboard
import wave
import io
from pydub import AudioSegment
from pydub.playback import play
from pocketsphinx import LiveSpeech
import speech_recognition as sr

#Loading the environment variables
load_dotenv()
openai.api_key = os.getenv('OPENAI_API_KEY')
client = openai.OpenAI()

stream = None

def stt_whisper():
    # This function uses OpenAI's whisper voice to text model to convert your voice input to text.
    record_audio()
    audio_file = open("user_response.wav", "rb")
    transcript = client.audio.transcriptions.create(
        model="whisper-1",
        file=audio_file
    )
    return transcript.text

def record_audio(duration=None):
    CHUNK = 1024
    FORMAT = 'int16'
    CHANNELS = 1
    RATE = 10000
    #WAVE_OUTPUT_FILENAME = "user_response.wav"
    WAVE_OUTPUT_FILENAME = os.path.join(r"D:\Projects\Senior Sage (Stevens Agile Project)\CS-555-Senior-Sage\backend\llm", f"user_response.wav")

    frames = []
    stream = None
    is_recording = False
    recording_stopped = False

    def start_audio_recording():
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
    
    

    #detect_wake_word()

    listener = keyboard.Listener(on_press=lambda key: stop_recording() if key == keyboard.Key.esc else None)
    listener.start()

    start_time = time.time()
    while listener.running:
        if recording_stopped:
            listener.stop()
        elif duration and (time.time() - start_time) > duration:
            listener.stop()
        time.sleep(0.01)

def detect_wake_word():
    create_keyword_file()  # Create a keyword file for the wake word
    speech = LiveSpeech(
        lm=False,
        keyphrase='hey sage',
        kws_threshold=1e-20,
        verbose=False
    )
    
    print("Listening for wake word...")
    for phrase in speech:
        print("Wake word detected!")
        return True  # Wake word detected, return True to indicate the detection

def create_keyword_file(keyword="hey sage", threshold=1e-20):
    with open("keyword.list", "w") as f:
        f.write(f"{keyword} /{threshold}/\n")
    """def on_key(key):
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
        time.sleep(0.01)"""


# To Do: Aadyant - Wake word (continous listening, how?), Prompt Engineer (concise, keep on asking if not understood), Data extraction;
# To Do: Prasoon - Memory [chat (current session), log_files (history of sessions)], TTS, Data extraction
# 
temp_dir = os.path.join(r"D:\Projects\Senior Sage (Stevens Agile Project)\CS-555-Senior-Sage\backend\llm\temp")
os.environ['TEMP'] = temp_dir
os.environ['TMP'] = temp_dir
def tts_whisper(input, voice):
    response = client.audio.speech.create(
        model="tts-1",
        voice=f"{voice.lower()}",
        input=f"{input}",
    )
    audio_data = response.content
    """temp_audio_path = os.path.join(r"D:\Projects\Senior Sage (Stevens Agile Project)\CS-555-Senior-Sage\backend\llm", "response.mp3")
    with open(temp_audio_path, "wb") as audio_file:
        audio_file.write(audio_data)

    # Load and play the audio file
    audio_segment = AudioSegment.from_file(temp_audio_path, format="mp3")
    play(audio_segment)"""
    audio_segment = AudioSegment.from_file(io.BytesIO(audio_data), format="mp3")
    try:
        play(audio_segment)
    except:
        print("teri maa ki choot")
