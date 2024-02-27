# ChichaTeleBot
Chicha telegram bot

1. **Build and Install:**

```bash
cd /usr/src
git clone https://github.com/matveynator/ChichaTeleBot.git
cd ChichaTeleBot
docker build -t chichatelebot .
rm -rf /usr/src/ChichaTeleBot
```

2. **Run:**
Replace `your_telegram_bot_token` with your actual Telegram bot token. After that, use the following command to launch the container:

```bash
docker run -d --restart unless-stopped -e TELEGRAM_BOT_TOKEN=your_telegram_bot_token -e DEBUG="false" --name ChichaTeleBot chichatelebot
```

**Important Privacy and Security Measures:**
- 🚨 **CRITICAL:** Ensure that you DO NOT store user messages. It is of utmost importance to prioritize user privacy and comply with data protection regulations.
- 🔐 Set the `DEBUG` environment variable to "false" to disable debugging. This helps prevent unnecessary information exposure, further safeguarding user data.

Now, with these privacy and security measures in place, you have a Docker container with Whisper installed. The CHICHA telebot will transcribe voice messages into text using the graphics card installed on your server.
