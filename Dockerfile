FROM alpine:3.21.3
RUN apk add --no-cache ca-certificates
WORKDIR /root/
# 빌드된 바이너리만 복사
COPY build/linux/alertmanager-webhook-proxy .
EXPOSE 8080
ENTRYPOINT ["./alertmanager-webhook-proxy"]