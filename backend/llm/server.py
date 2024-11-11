from flask import Flask, render_template
from flask_socketio import SocketIO, emit

app = Flask(__name__)
app.config['SECRET_KEY'] = 'your_secret_key'
socketio = SocketIO(app)


def authenticate(token):
    # Implement your authentication logic here
    # For example, check the token against a database or other storage
    return token == "your_secret_token"

@socketio.on('connect')
def handle_connect():
    print('Client connected')
    emit('connected', {'data': 'Connected'})

@socketio.on('disconnect')
def handle_disconnect():
    print('Client disconnected')

@socketio.on('message')
def handle_message(data):
    print('Received message:', data)


    # Authenticate the client
    if not authenticate(data['token']):
        emit('error', {'message': 'Unauthorized'})
        return

    # Process the message and emit a response
    emit('response', {'data': 'Message received'})

if __name__ == '__main__':
    socketio.run(app)