# ChichaTeleBot - powered by OpenAi WHISPER transformer LM.
ChichaTeleBot is a powerful voice bot for Telegram designed to be your helpful companion. It seamlessly converts spoken words to text, performs language translations based on your speech, and **prioritizes the protection of your privacy**. The bot is named after Chicha, a delightful companion pincher dog. 

Please note: **CHICHA NEVER STORES YOU TEXT OR VOICE MESSAGES.**

🟢 **Live Telegram Demo:** [@ChichaTeleBot](https://t.me/ChichaTeleBot) 

## Quick Installation via Docker:
Install ChichaTeleBot effortlessly using Docker with the following command. Replace `your_telegram_bot_token` with your actual Telegram bot token. Set the `MODEL` variable to "small," "medium," or "large," and adjust the `DEBUG` variable to "true" or "false" for debugging preferences. Ensure to name the container (--name your_telegram_bot_name) for easy differentiation if you have multiple bots:

```bash

docker run -d --restart unless-stopped  -e DEBUG="false" -e MODEL="medium" -e TELEGRAM_BOT_TOKEN="your_telegram_bot_token" --mount type=tmpfs,destination=/root/.cache/whisper --gpus all --name your_telegram_bot_name matveynator/chichatelebot:latest

```

## Privacy and Security Measures:
### 🚨 CRITICAL:
Ensure that user messages are NOT stored to prioritize user privacy and comply with data protection regulations.
### 🔐 Security:
Set the `DEBUG` environment variable to "false" to disable debugging, preventing unnecessary information exposure and enhancing user data protection.

## System Requirements:
- 🖥️ CUDA-enabled graphics cards are required.
- Tested CUDA on UBUNTU Linux ONLY: https://docs.nvidia.com/cuda/cuda-installation-guide-linux/index.html 
- 🚀 Recommended NVIDIA graphics card: RTX 4090.
- ✅ Tested and compatible NVIDIA graphics card: RTX 2070.

## Building from Source:
Build ChichaTeleBot from the source code with the following commands:
```bash
rm -rf /usr/src/ChichaTeleBot
mkdir -p /usr/src
cd /usr/src
git clone https://github.com/matveynator/ChichaTeleBot.git
cd ChichaTeleBot
docker build -t chichatelebot .
```
```
docker run -d --restart unless-stopped -e TELEGRAM_BOT_TOKEN=your_telegram_bot_token -e MODEL=medium -e DEBUG="false" --gpus all --name your_telegram_bot_name chichatelebot
```

Now, you have a fully functional ChichaTeleBot running on your server, providing a seamless voice-to-text experience while ensuring the utmost privacy and security for your users.


## INSTALLING NVIDIA CUDA for Docker on Ubuntu:
```
echo "Installing CUDA Toolkit for Docker onUbuntu..." && distribution=$(. /etc/os-release; echo $ID$VERSION_ID) && curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg && curl -s -L https://nvidia.github.io/libnvidia-container/$distribution/libnvidia-container.list | sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | tee /etc/apt/sources.list.d/nvidia-container-toolkit.list && apt-get update && apt-get -y install --reinstall nvidia-utils-535-server nvidia-dkms-535-server && apt-get install -y nvidia-container-toolkit && systemctl restart docker && echo "CUDA Toolkit installation completed."
```

