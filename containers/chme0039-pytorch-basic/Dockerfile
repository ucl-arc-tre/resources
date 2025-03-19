FROM pytorch/pytorch:2.5.1-cuda12.4-cudnn9-runtime

# Init repos
RUN apt-get update
# Install dev dependencies
RUN apt-get install -y -q curl git vim tmux unzip nano

WORKDIR /workspace

# For those who are fond of requirements.txt
COPY requirements.txt /workspace/requirements.txt
RUN pip install -r /workspace/requirements.txt

# streamlit port
EXPOSE 8501
