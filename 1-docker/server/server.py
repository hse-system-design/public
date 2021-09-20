from flask import Flask, request

app = Flask(__name__)  # Основной объект приложения Flask


@app.route('/')
def hello():
    return "Hello, from Flask"


@app.route('/biba')
def biba():
    return "Kuka"


if __name__ == '__main__':
    app.run("0.0.0.0", 9091)  # Запускаем сервер на 9091 порту
