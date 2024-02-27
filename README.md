# ChichaTeleBot
Chicha telegram bot

1. **Docker Image Build:**

Open the terminal in the directory where the `Dockerfile` is located and execute the following command:

```bash
docker build -t ChichaTeleBot .
```

2. **Container Launch:**

Replace `your_telegram_bot_token` with your actual Telegram bot token. After that, use the following command to launch the container:

```bash
docker run -d --restart unless-stopped -e TELEGRAM_BOT_TOKEN=your_telegram_bot_token ChichaTeleBot
```

Now you should have a Docker container with Whisper installed, running in the background as the CHICHA telebot. It will transcribe any received voice messages into text using the graphics card installed on your server.
