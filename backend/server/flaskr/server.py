from flask import Flask, jsonify, make_response, request
import jwt
import service.service as service
from middlewares.auth_middleware import token_required
from login import validate_email_and_password
import os

app = Flask(__name__)
SECRET_KEY = os.environ.get('SECRET_KEY') or 'this is a secret'
print(SECRET_KEY)
app.config['SECRET_KEY'] = SECRET_KEY

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

@app.route("/api/login", methods=["POST"])
def login():
    try:
        data = request.json
        if not data:
            return {
                "message": "Please provide user details",
                "data": None,
                "error": "Bad request"
            }, 400
        # validate input
        print(data.get('email'))
        print(data.get('password'))
        is_validated = validate_email_and_password(data.get('email'), data.get('password'))
        if is_validated is not True:
            return dict(message='Invalid data', data=None, error=is_validated), 400
        user = {
            "email": data["email"],
            "password": data["password"],
            "_id": "fae98fadsfh98"
        }
        if user:
            try:
                # token should expire after 24 hrs
                user["token"] = jwt.encode(
                    {"user_id": user["_id"]},
                    app.config["SECRET_KEY"],
                    algorithm="HS256"
                )
                return {
                    "message": "Successfully fetched auth token",
                    "data": user
                }, 200
            except Exception as e:
                return {
                    "error": "Something went wrong",
                    "message": str(e)
                }, 500
        return {
            "message": "Error fetching auth token!, invalid email or password",
            "data": None,
            "error": "Unauthorized"
        }, 404
    except Exception as e:
        return {
                "message": "Something went wrong!",
                "error": str(e),
                "data": None
        }, 500

@app.route("/api/test", methods=['GET'])
@token_required
def test(current_user):
    print(f'got current_user: {current_user}')
    return {
        "message": "done",
    }, 200

if __name__ == '__main__':
    app.run(host='0.0.0.0')