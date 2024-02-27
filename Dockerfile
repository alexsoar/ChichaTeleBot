# Используем официальный образ Debian 12
FROM debian:12

# Установка необходимых пакетов
RUN apt-get update && apt-get install -y python3-pip python3-venv git golang-go ffmpeg

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Создание виртуального окружения
RUN python3 -m venv /venv

# Активация виртуального окружения
ENV PATH="/venv/bin:$PATH"

# Установка pipx в виртуальном окружении
RUN /venv/bin/python3 -m pip install --upgrade pip pipx && \
    /venv/bin/python3 -m pipx ensurepath

# Добавление /root/.local/bin в переменную среды PATH
ENV PATH="/root/.local/bin:$PATH"

# Установка Whisper через pipx
RUN /venv/bin/pipx install git+https://github.com/openai/whisper.git

# Добавление /root/.local/bin в переменную среды PATH
ENV PATH="/root/.local/bin:$PATH"

# Клонирование и компиляция ChichaTeleBot.go
WORKDIR /app
RUN git clone https://github.com/matveynator/ChichaTeleBot && \
    cd ChichaTeleBot && \
    rm -f go.mod && \
    rm -f go.sum  && \
    go mod init ChichaTeleBot && \
    go mod tidy && \
    go build -o /usr/local/bin/ChichaTeleBot ChichaTeleBot.go

RUN /root/.local/bin/whisper --model medium  /app/ChichaTeleBot/test.ogg

# Добавление разрешений на выполнение
RUN chmod +x /usr/local/bin/ChichaTeleBot

# Запуск ChichaTeleBot как демона
CMD ["/usr/local/bin/ChichaTeleBot"]
