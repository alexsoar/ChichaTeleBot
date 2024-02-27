# Используем официальный образ Debian 12
FROM debian:12

# Установка необходимых пакетов
RUN apt-get update && apt-get install -y python3-pip python3-venv git golang-go

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Создание виртуального окружения
RUN python3 -m venv /venv

# Активация виртуального окружения
ENV PATH="/venv/bin:$PATH"

# Установка pipx в виртуальном окружении
RUN /venv/bin/python3 -m pip install --upgrade pip pipx && \
    /venv/bin/python3 -m pipx ensurepath

# Установка Whisper через pipx
RUN /venv/bin/pipx install git+https://github.com/openai/whisper.git

# Клонирование и компиляция ChichaTeleBot.go
WORKDIR /app
RUN git clone https://github.com/matveynator/ChichaTeleBot && \
    cd ChichaTeleBot && \
    go mod init ChichaTeleBot && \
    go mod tidy && \
    go build -o ChichaTeleBot ChichaTeleBot.go && \
    chmod +x ./ChichaTeleBot

# Запуск ChichaTeleBot как демона
CMD ["./ChichaTeleBot"]
