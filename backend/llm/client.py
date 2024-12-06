import socketio
import sounddevice as sd
from pynput import keyboard
import wave
import time
from pydub.playback import play
import pydub
from pydub import AudioSegment
import io
import logging

sio = socketio.Client()


@sio.on('connect', namespace='/llm')
def connect():
    print('connected to server')
    sio.emit('auth', namespace='/llm')

@sio.on('connected', namespace='/llm')
def connected(data):
    # print('connected to server')
    record_audio()
    sio.emit('voice_input', open("user_response.wav", 'rb').read(), namespace='/llm')

@sio.on('voice_response', namespace='/llm')
def handle_voice_response(data):
    # handle_voice_response(data)
    aud_seg = AudioSegment.from_file(io.BytesIO(data['data']), format="mp3")
    play(aud_seg)
    if not data['disconnect']:
        record_audio()
        sio.emit('voice_input', open("user_response.wav", 'rb').read(), namespace='/llm')
    else:
        global close
        close = True

    

@sio.on('connection_denied', namespace='/llm')
def connection_denied(data):
    print('auth failed', data)

# def handle_voice_response(data):
#     # Decode audio data (if necessary)
#     decoded_audio = data['data']

#     # Create pydub.AudioSegment object (if not already an AudioSegment)
#     audio_segment = pydub.AudioSegment.from_mp3(decoded_audio)  # Adjust format if needed

#     try:
#         # Attempt playback with ffplay (fallback if missing)
#         audio_segment.export(f.name, "wav")
#     except (OSError, AttributeError):  # OSError for missing ffplay
#         print("Failed to play audio using ffplay. Falling back to alternative method...")
#         # Implement alternative playback logic (e.g., pyaudio)
#     return audio_segment
close = False
def record_audio(duration=None):
    CHUNK = 1024
    FORMAT = 'int16'
    CHANNELS = 1
    RATE = 10000
    WAVE_OUTPUT_FILENAME = "user_response.wav"

    frames = []
    stream = None
    is_recording = False
    recording_stopped = False
    global close
    print('Press UP to start recording, UP to stop recording')

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
        global close

        if key == keyboard.Key.up:
            if not is_recording:
                print('recording audio')
                record_audio()
                is_recording = True
            else:
                print('recording captured, voice assistant is thinking')
                stop_recording()
                is_recording = False
        elif key == keyboard.Key.down:
            print('closing')
            close = True

    listener = keyboard.Listener(on_press=on_key)
    listener.start()

    start_time = time.time()
    while listener.running:
        if close:
            listener.stop()
            break
        if recording_stopped:
            listener.stop()
        elif duration and (time.time() - start_time) > duration:
            listener.stop()
        
        time.sleep(0.01)

if __name__ == '__main__':
    # logging.getLogger('socketio').setLevel(logging.WARNING)
    # logging.getLogger('engineio').setLevel(logging.WARNING)
    while True:
        print('Select user: press 1 for Alice, 2 for Bob')
        user = input()
        if user == '1':
            print('logging in as Alice')
            time.sleep(0.01)
            user_id = '91c6e005d9fcd20a6cbf1fe1dfdd319d'
            token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzYWdlLXNlcnZlciIsImV4cCI6MTc2OTQ0NzQ5MSwibmJmIjoxNzMzNDUxMDkxLCJpYXQiOjE3MzM0NTEwOTEsInVzZXJfaWQiOiI5MWM2ZTAwNWQ5ZmNkMjBhNmNiZjFmZTFkZmRkMzE5ZCJ9.VM2HcCgZQPu_dGgtsTBbGlvHD3QrRnxTGJwsWJokM9E'
            break
        elif user == '2':
            print('logging in as Bob')
            time.sleep(0.01)
            user_id = 'ee82c18d386dd4f0b6e7fae47e5fc1a6'
            token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzYWdlLXNlcnZlciIsImV4cCI6MTc2OTQ2OTc2NCwibmJmIjoxNzMzNDczMzY0LCJpYXQiOjE3MzM0NzMzNjQsInVzZXJfaWQiOiJlZTgyYzE4ZDM4NmRkNGYwYjZlN2ZhZTQ3ZTVmYzFhNiJ9.2JYxW9M3OiDTc6Lv2uGJlRlDEFxckX9_y0EKM6o93FA'
            break
        else:
            print('Invalid user')
            
    sio.logger.setLevel(logging.WARNING)
    sio.connect('http://127.0.0.1:52714', headers={
        'Authorization': token
    }, namespaces='/llm')
    # sio.wait()
    while close == False:
        time.sleep(0.01)
    print('closing connection')
    sio.disconnect()
    # record_audio()

