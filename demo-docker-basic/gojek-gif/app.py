
from flask import Flask, render_template
import random

app = Flask(__name__)

# list of gojek images
images = [
    "https://media.giphy.com/media/NsBcxCuRcdybSDLbvN/giphy.gif",
    "https://media.giphy.com/media/9LZSJyCzVQhgNiADmS/giphy.gif",
    "https://media.giphy.com/media/vRHCshZDU4tCiLhFYM/giphy.gif",
    "https://media.giphy.com/media/1BGeVoYKFFusTL2Q3j/giphy.gif",
    "https://media.giphy.com/media/229B4icxA8tC1Bo3eH/giphy.gif"
]

@app.route('/')
def index():
    url = random.choice(images)
    return render_template('index.html', url=url)

if __name__ == "__main__":
    app.run(host="0.0.0.0")
