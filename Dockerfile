FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o /app/server ./src

FROM --platform=linux/amd64 ubuntu:22.04 AS runner

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        wget ca-certificates fonts-liberation && \
    wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    apt-get install -y ./google-chrome-stable_current_amd64.deb && \
    rm google-chrome-stable_current_amd64.deb && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

ENV CHROME_PATH=/usr/bin/google-chrome-stable

WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
ENTRYPOINT ["./server"]
