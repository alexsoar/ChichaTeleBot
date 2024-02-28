# Dockerfile of github.com/matveynator/ChichaTeleBot TALK-to-TEXT telegram bot.
FROM ubuntu:22.04

# Set the default runtime to NVIDIA
ENV NVIDIA_VISIBLE_DEVICES all
ENV NVIDIA_DRIVER_CAPABILITIES compute,utility

# Place where models will be stored.
RUN mkdir -p /root/.cache/whisper

# tmpfs (8G)
CMD ["mount", "-t", "tmpfs", "-o", "size=8G", "tmpfs", "/root/.cache/whisper"]

# Install necessary dependencies
RUN apt-get update && apt-get -y install curl gnupg lsb-release

# Add NVIDIA Container Toolkit repository
RUN distribution=$(. /etc/os-release; echo $ID$VERSION_ID) && \
    curl -fsSL https://nvidia.github.io/libnvidia-container/gpgkey | gpg --dearmor -o /usr/share/keyrings/nvidia-container-toolkit-keyring.gpg && \
    curl -s -L https://nvidia.github.io/libnvidia-container/$distribution/nvidia-container-toolkit.list | \
    sed 's#deb https://#deb [signed-by=/usr/share/keyrings/nvidia-container-toolkit-keyring.gpg] https://#g' | \
    tee /etc/apt/sources.list.d/nvidia-container-toolkit.list;

# Install NVIDIA Container Toolkit and CUDA Toolkit
RUN apt-get update; apt-get install -y nvidia-container-toolkit nvidia-cuda-toolkit nvidia-container-runtime;

# Install necessary packages
RUN apt-get install -y python3 python3-pip python3-venv git golang-go ffmpeg;

# Clean up the package manager cache to reduce image size
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* ;

# Create a virtual environment
RUN python3 -m venv /venv

# Activate the virtual environment
ENV PATH="/venv/bin:/root/.local/bin:$PATH"

# Install pipx within the virtual environment
RUN /venv/bin/pip3 install --upgrade pip 

# Install Whisper using pipx
RUN /venv/bin/pip3 install -U openai-whisper

# Clone and compile ChichaTeleBot.go
WORKDIR /app
RUN git clone https://github.com/matveynator/ChichaTeleBot && \
    cd ChichaTeleBot && \
    rm -f go.mod && \
    rm -f go.sum  && \
    go mod init ChichaTeleBot && \
    go mod tidy && \
    go build -o /usr/local/bin/ChichaTeleBot ChichaTeleBot.go

# Run Whisper to generate a summary for the provided audio file
RUN whisper --model medium  /app/ChichaTeleBot/test.ogg

# Add execution permissions
RUN chmod +x /usr/local/bin/ChichaTeleBot

# Run ChichaTeleBot as a daemon
CMD ["/usr/local/bin/ChichaTeleBot"]
