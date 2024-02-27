# ChichaTeleBot
Chicha telegram bot

1. **Build and install:**

```bash
cd /usr/src
git clone https://github.com/matveynator/ChichaTeleBot.git
cd ChichaTeleBot
docker build -t chichatelebot .
rm -rf /usr/src/ChichaTeleBot
```

Replace `your_telegram_bot_token` with your actual Telegram bot token. After that, use the following command to launch the container:
```bash
docker run -d --restart unless-stopped -e TELEGRAM_BOT_TOKEN=your_telegram_bot_token -e DEBUG="false" --name ChichaTeleBot chichatelebot
```

Now you should have a Docker container with Whisper installed, running in the background as the CHICHA telebot. It will transcribe any received voice messages into text using the graphics card installed on your server.
