# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/azure-keyvault-provider \
    ./cmd/provider


FROM alpine:3.22

LABEL org.opencontainers.image.title="Azure Key Vault Secrets Provider" \
      org.opencontainers.image.description="Docker Secrets Engine provider for Azure Key Vault" \
      org.opencontainers.image.vendor="Kryptonian" \
      com.docker.desktop.extension.api.version=">= 0.2.0" \
      com.docker.desktop.extension.icon="https://raw.githubusercontent.com/kryptonian001/docker-azure-keyvault-provider/refs/heads/main/extension-icon.svg" \
      com.docker.extension.screenshots="[]" \
      com.docker.extension.detailed-description="Azure Key Vault provider for Docker Secrets Engine." \
      com.docker.extension.publisher-url="https://github.com/kryptonian001" \
      com.docker.extension.categories="security" \
      com.docker.extension.changelog="Initial development release."

COPY --from=builder /out/azure-keyvault-provider /usr/local/bin/azure-keyvault-provider

COPY metadata.json /
COPY compose.yaml /
COPY extension-icon.svg /

ENTRYPOINT ["/usr/local/bin/azure-keyvault-provider"]