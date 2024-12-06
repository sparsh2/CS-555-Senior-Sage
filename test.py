import openai
import io
client = openai.OpenAI(api_key='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')

transcript = client.audio.transcriptions.create(
        model="whisper-1",
        # file=io.BufferedReader(io.BytesIO(audio_file))
        file=open("./backend/llm/user_response.wav", 'rb')
        # content=audio_file
    )