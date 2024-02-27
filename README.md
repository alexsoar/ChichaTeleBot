# ChichaTeleBot

ChichaTeleBot is a powerful voice bot for Telegram designed to be your helpful companion. It seamlessly converts spoken words to text, performs language translations based on your speech, and prioritizes the protection of your privacy. The bot is aptly named after Chicha, a delightful companion pincher dog.

## Quick Installation via Docker

Install ChichaTeleBot effortlessly using Docker with the following command. Replace `your_telegram_bot_token` with your actual Telegram bot token. Set the `MODEL` variable to "small," "medium," or "large," and adjust the `DEBUG` variable to "true" or "false" for debugging preferences. Ensure to name the container (--name your_telegram_bot_name) for easy differentiation if you have multiple bots:

```bash

docker run -d --restart unless-stopped  -e DEBUG="false" -e MODEL=medium -e TELEGRAM_BOT_TOKEN="your_telegram_bot_token" --name your_telegram_bot_name matveynator/chichatelebot:latest

```

## Privacy and Security Measures

### 🚨 CRITICAL:
Ensure that user messages are NOT stored to prioritize user privacy and comply with data protection regulations.

### 🔐 Security:
Set the `DEBUG` environment variable to "false" to disable debugging, preventing unnecessary information exposure and enhancing user data protection.

## System Requirements

- 🖥️ CUDA-enabled graphics cards are required.
- 🚀 Recommended NVIDIA graphics card: RTX 4090.
- ✅ Tested and compatible NVIDIA graphics card: RTX 2070.

## Building from Source

Build ChichaTeleBot from the source code with the following commands:

```bash
rm -rf /usr/src/ChichaTeleBot
mkdir -p /usr/src
cd /usr/src
git clone https://github.com/matveynator/ChichaTeleBot.git
cd ChichaTeleBot
docker build -t chichatelebot .
docker run -d --restart unless-stopped -e TELEGRAM_BOT_TOKEN=your_telegram_bot_token -e MODEL=medium -e DEBUG="false" --name your_telegram_bot_name chichatelebot
```

Now, you have a fully functional ChichaTeleBot running on your server, providing a seamless voice-to-text experience while ensuring the utmost privacy and security for your users.
