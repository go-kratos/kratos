FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY ./app/message/bin /app
COPY ./app/user/bin /app
COPY ./service.sh /app

EXPOSE 8000
EXPOSE 9000

ENTRYPOINT ["sh","/app/service.sh"]