FROM mcr.microsoft.com/vscode/devcontainers/base:ubuntu

COPY --from=mikefarah/yq:4.30.2 /usr/bin/yq /usr/local/bin/yq
COPY --from=koalaman/shellcheck:v0.8.0 /bin/shellcheck /usr/local/bin/shellcheck

USER vscode
RUN echo "source <(lane completion bash)" >> ~/.bashrc
RUN echo "source <(lane completion zsh)" >> ~/.zshrc
