# ChichaTeleBot is a voice bot for Telegram designed to be a helpful companion. 
It can convert spoken words to text, perform language translations based on your speech, and ensures the protection of your privacy. The bot is named after the dog Chicha, which is a companion pincher dog.

1. **Build and Install:**

```bash
rm -rf /usr/src/ChichaTeleBot
mkdir -p /usr/src
cd /usr/src
git clone https://github.com/matveynator/ChichaTeleBot.git
cd ChichaTeleBot
docker build -t chichatelebot .
```

2. **Run:**

Replace `your_telegram_bot_token` with your actual Telegram bot token. Additionally, set the `MODEL` variable to either "small," "medium," or "large" based on your preferred model size. Also, adjust the `DEBUG` variable to "true" or "false" to enable or disable debugging. After making these adjustments, use the following command to launch the container, and **ensure to name it (--name your_telegram_bot_name) according to your Telegram bot's name for easy differentiation if you have multiple bots:**

```bash
docker run -d --restart unless-stopped -e TELEGRAM_BOT_TOKEN=your_telegram_bot_token -e MODEL=medium -e DEBUG="false" --name your_telegram_bot_name --gpus all chichatelebot
```

**Important Privacy and Security Measures:**

- 🚨 **CRITICAL:** Ensure that you DO NOT store user messages. It is of utmost importance to prioritize user privacy and comply with data protection regulations.
- 🔐 Set the `DEBUG` environment variable to "false" to disable debugging. This helps prevent unnecessary information exposure, further safeguarding user data.

**System Requirements:**

- 🖥️ This bot requires the use of NVIDIA CUDA-enabled graphics cards.
- 🚀 Recommended NVIDIA graphics card: RTX 4090.
- ✅ Tested and compatible NVIDIA graphics card: RTX 2070.

Now, with these privacy and security measures in place, you have a Docker container with Whisper installed. The CHICHA telebot will transcribe voice messages into text using the specified model size and the graphics card installed on your server. Naming the container according to your Telegram bot's name will help in distinguishing between multiple bots.
