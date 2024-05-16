from flask import Flask

app = Flask(__name__)

@app.route('/')
def home():
    return 'Это главная страница.'

@app.route('/about')
def about():
    return 'Здесь будет информация об авторе сайта.'

@app.route('/blog')
def blog():
    return 'Это блог с заметками о работе и увлечениях.'

if __name__ == '__main__':
    app.run()