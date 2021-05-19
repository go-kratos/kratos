FROM golang:1.15 AS builder

COPY . /src

COPY ../ /src

WORKDIR /src/app/message

RUN GOPROXY=https://goproxy.cn go build -o ./bin/message ./...

WORKDIR /src/app/user

RUN GOPROXY=https://goproxy.cn go build -o ./bin/user ./...

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/app/message/bin /app
COPY --from=builder /src/app/user/bin /app
COPY --from=builder /src/start.sh /app

EXPOSE 8000
EXPOSE 9000

ENTRYPOINT ["sh","/app/start.sh"]