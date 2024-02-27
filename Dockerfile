# Используем официальный образ Debian 12
FROM debian:12

# Установка необходимых пакетов
RUN apt-get update && \
    apt-get install -y python3-pip git golang-go

# Установка pipx
RUN python3 -m pip install --upgrade pip pipx

# Установка Whisper через pipx
RUN pipx install git+https://github.com/openai/whisper.git

# Добавление исполняемого файла Whisper в системный путь
RUN ln -s /root/.local/bin/whisper /usr/local/bin/whisper

# Клонирование и компиляция ChichaTeleBot.go
WORKDIR /app
ADD https://raw.githubusercontent.com/matveynator/ChichaTeleBot/main/ChichaTeleBot.go /app/
RUN go build ChichaTeleBot.go

# Запуск telebot как демона
CMD ["/app/ChichaTeleBot"]
