FROM uhub.service.ucloud.cn/entropypool/alpine:3.16

WORKDIR /app

COPY notif-gateway /app

ENTRYPOINT ["/app/notif-gateway"]

EXPOSE 50700 50701 50702 50703

CMD ["run"]
