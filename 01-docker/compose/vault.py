from flask import Flask

app = Flask(__name__)
key = "O6AMGY6yC9xKi85QxjZTplX3OZm64j88c7Aq4iXqQGA="


@app.route("/secret-key")
def secret_key():
    return key


if __name__ == '__main__':
    app.run("0.0.0.0", 8080)
