FROM alpine:3

RUN apk --update --no-cache add busybox-static \
    && apk --no-cache add --virtual build-dependencies curl build-base tzdata

ENV TZ=Asia/Tokyo

COPY crontab /var/spool/cron/crontabs/root

CMD ["busybox", "crond", "-f", "-L", "/dev/stderr"]
