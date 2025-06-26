FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY ../../ .

RUN go build -a -o ./bin/server ./main.go

FROM ubuntu:22.04 AS runner

WORKDIR /app

RUN apt-get update && apt-get install -y wget

RUN wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb

RUN apt-get install -y ./google-chrome-stable_current_amd64.deb

ENV CHROME_PATH="/usr/bin/google-chrome-stable"

COPY --from=builder /app/bin/ .

EXPOSE 8080

ENTRYPOINT ["./server"]