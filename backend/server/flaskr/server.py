from flask import Flask, jsonify, make_response
import service.service as service

app = Flask(__name__)

@app.route("/")
def hello_world():
    return "<p>Hello, World!</p>"

@app.route("/api/delete-preference", methods=['POST'])
def delete_preference():
    done, msg = service.delete_preferences({})
    if not done: 
        return make_response(jsonify({'error': msg, 'message': f'failed to delete preferences: {msg}'}), 500)
    data = {'message': 'preferences deleted'}
    response = make_response(jsonify(data), 200)
    return response

if __name__ == '__main__':
    app.run(host='0.0.0.0')