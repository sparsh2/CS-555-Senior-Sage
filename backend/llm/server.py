from flask import Flask, render_template
from flask_socketio import SocketIO, emit
from functools import wraps
from flask import request
from flask_socketio import disconnect, emit
import requests
from other import llm_authenticate, pull_user_data, del_user_data, get_response_data_from_llm
import time

import yaml

app = Flask(__name__)
app.config['SECRET_KEY'] = 'your_secret_key'
socketio = SocketIO(app, cors_allowed_origins='*')

@app.route('/ping', methods=['GET'])
def ping():
    return 'pong', 200



with open('/app/config/conf.yaml', 'r') as file:
    cfg = yaml.safe_load(file)
# cfg = {'llmUsername': 'llmuser', 'llmPassword': 'llmpassword', 'openaiApiKey': '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', 'storageService': {'host': 'storage-service-svc.senior-sage', 'port': 8080}, 'authzService': {'host': 'authz-svc.senior-sage', 'port': 8080}}
print(cfg)
ok = llm_authenticate(cfg)
if not ok:
    print("Error authenticating with LLM")
    time.sleep(5)
    exit(1)

def token_required(f):
    @wraps(f)
    def decorated(*args, **kwargs):
        global cfg
        token = None
        if "Authorization" in request.headers:
            token = request.headers["Authorization"]
        print(token)
        if not token:
            return {
                "message": "Authentication Token is missing!",
                "data": None,
                "error": "Unauthorized"
            }, 401
        try:
            authzHost = cfg['authzService']['host']
            authzPort = cfg['authzService']['port']
            response = requests.get(f"http://{authzHost}:{authzPort}/auth/verify", json={"jwt_token": token})
            response.raise_for_status()
            current_user = response.json()
            if current_user['valid'] == False:
                return {
                    "message": "Invalid Token",
                    "data": None,
                    "error": "Unauthorized"
                }, 401
            current_user = current_user['user_id']
        except Exception as e:
            return {
                "message": "Something went wrong",
                "data": None,
                "error": str(e)
            }, 500

        return f(current_user, *args, **kwargs)

    return decorated

state_data = {}


@socketio.on('connect', namespace='/llm')
@token_required
def handle_connect(current_user):
    state_data[request.sid] = current_user
    pull_user_data(cfg, current_user)
    print(current_user)
    print('Client connected')
    emit('connected', {'data': 'Connected'})

@socketio.on('disconnect', namespace='/llm')
def handle_disconnect():
    del_user_data(state_data[request.sid])
    del state_data[request.sid]
    print('Client disconnected')


@socketio.on('voice_input', namespace='/llm')
def handle_voice_capture(raw_voice_data):
    user_id = state_data[request.sid]
    voice_response, disconnect = get_response_data_from_llm(user_id, raw_voice_data)
    emit('voice_response', {'data': voice_response, 'disconnect': disconnect})

if __name__ == '__main__':
    print('Starting server...')
    socketio.run(app, allow_unsafe_werkzeug=True, host='0.0.0.0')