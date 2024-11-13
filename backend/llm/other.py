import requests

llm_token = ""

def llm_authenticate(cfg):
    global llm_token
    username = cfg['llmUsername']
    password = cfg['llmPassword']
    try:
        host = cfg['authzService']['host']
        port = cfg['authzService']['port']
        response = requests.post(f"http://{host}:{port}/auth/login", json={"email": username, "password": password})
        response.raise_for_status()
        token = response.json()['token']
    except Exception as e:
        print(f"Error authenticating with LLM: {e}")
        return False
    llm_token = token
    return True

all_user_data = {}

def pull_user_data(cfg, user_id):
    global llm_token, all_user_data
    try:
        host = cfg['storageService']['host']
        port = cfg['storageService']['port']
        response = requests.get(f"http://{host}:{port}/data", json={"requester_token": llm_token, "user_id": user_id})
        response.raise_for_status()
        user_data = response.json()
        all_user_data[user_id] = {
            'user_data': user_data, 
            'current_session': {}
        }
    except Exception as e:
        print(f"Error pulling user data: {e}")

def del_user_data(user_id):
    global all_user_data
    if user_id in all_user_data:
        del all_user_data[user_id]

def get_response_data_from_llm(user_id, voice_data):
    current_user_data = all_user_data.get(user_id, {})

    pass