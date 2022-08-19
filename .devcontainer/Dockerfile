FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu

RUN wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 && chmod +x /usr/local/bin/yq

RUN apt-get update

RUN apt-get install -y shellcheck

RUN apt-get install -y openssl
