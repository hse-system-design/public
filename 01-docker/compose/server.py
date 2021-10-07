import base64

from flask import Flask, request
import requests
from cryptography.fernet import Fernet, MultiFernet

app = Flask(__name__)  # Основной объект приложения Flask


def read_key():
    with open('/secret/key.txt', 'br') as f:
        key1 = Fernet(f.read())

    r = requests.get("http://vault:8080/secret-key")
    key2 = Fernet(r.content)

    return MultiFernet([key1, key2])


@app.route('/')
def hello():
    return "Hello, from Flask"


@app.route('/encrypt', methods=["POST"])
def encrypt():
    fnet = read_key()
    data = request.form['data'].encode('utf-8')
    return base64.urlsafe_b64encode(fnet.encrypt(data))


@app.route('/decrypt', methods=["POST"])
def decrypt():
    fnet = read_key()
    data = base64.urlsafe_b64decode(request.form['data'])
    return fnet.decrypt(data).decode('utf-8')


if __name__ == '__main__':
    app.run("0.0.0.0", 9091)  # Запускаем сервер на 9091 порту
