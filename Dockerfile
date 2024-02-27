# Use the official Debian 12 image as the base
FROM debian:12

# Set the default runtime to NVIDIA
ENV NVIDIA_VISIBLE_DEVICES all
ENV NVIDIA_DRIVER_CAPABILITIES compute,utility

Install NVIDIA CUDA Toolkit
RUN apt-get install -y gnupg2 wget lsb-release && \
    wget https://developer.download.nvidia.com/compute/cuda/repos/debian/$(lsb_release -cs)/x86_64/cuda-$(dpkg --print-architecture)/cuda-repo-$(lsb_release -cs)_$(dpkg --print-architecture).deb && \
    dpkg -i cuda-repo-$(lsb_release -cs)_$(dpkg --print-architecture).deb && \
    apt-key adv --fetch-keys http://developer.download.nvidia.com/compute/cuda/repos/debian/$(lsb_release -cs)/x86_64/7fa2af80.pub && \
    apt-get update && \
    apt-get install -y cuda nvidia-cuda-toolkit

# Install necessary packages
RUN apt-get install -y python3-pip python3-venv git golang-go ffmpeg

# Clean up the package manager cache to reduce image size
RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Create a virtual environment
RUN python3 -m venv /venv

# Activate the virtual environment
ENV PATH="/venv/bin:$PATH"

# Install pipx within the virtual environment
RUN /venv/bin/python3 -m pip install --upgrade pip pipx && \
    /venv/bin/python3 -m pipx ensurepath

# Add /root/.local/bin to the PATH environment variable
ENV PATH="/root/.local/bin:$PATH"

# Install Whisper using pipx
RUN /venv/bin/pipx install git+https://github.com/openai/whisper.git

# Add /root/.local/bin to the PATH environment variable
ENV PATH="/root/.local/bin:$PATH"

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
RUN /root/.local/bin/whisper --model medium  /app/ChichaTeleBot/test.ogg

# Add execution permissions
RUN chmod +x /usr/local/bin/ChichaTeleBot

# Run ChichaTeleBot as a daemon
CMD ["/usr/local/bin/ChichaTeleBot"]
