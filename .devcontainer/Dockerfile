FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu

RUN apt update && apt install -y xz-utils openssl

RUN wget -qO /usr/local/bin/yq https://github.com/mikefarah/yq/releases/download/v4.27.2/yq_linux_amd64 && chmod +x /usr/local/bin/yq

RUN wget -qO - https://github.com/koalaman/shellcheck/releases/download/v0.8.0/shellcheck-v0.8.0.linux.x86_64.tar.xz | tar -xJv --strip-components=1 -C /tmp && mv /tmp/shellcheck /usr/local/bin/ && rm -rf /tmp/* && chmod +x /usr/local/bin/shellcheck
