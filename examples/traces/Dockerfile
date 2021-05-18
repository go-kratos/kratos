FROM golang:1.15 AS builder

COPY . /src

WORKDIR /src/app/blog

RUN GOPROXY=https://goproxy.cn go build -o ./bin/ ./...

WORKDIR /src/app/user

RUN GOPROXY=https://goproxy.cn go build -o ./bin/ ./...

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/app/blog/bin /app/blog
COPY --from=builder /src/app/user/bin /app/user
COPY --from=builder /src/start.sh /app
COPY --from=builder /src/app/blog/configs /app/blog/configs
COPY --from=builder /src/app/user/configs /app/user/configs

VOLUME /data/conf

WORKDIR /

EXPOSE 8011

ENTRYPOINT ["sh","/app/start.sh"]