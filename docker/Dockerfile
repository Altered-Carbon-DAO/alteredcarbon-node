FROM alpine:3.14
RUN apk add -U --no-cache ca-certificates

WORKDIR /data
ENV HOME=/data
COPY ./bin/acarbd /usr/bin/acarbd
COPY ./docker/entry-point.sh ./entry-point.sh
EXPOSE 26657

CMD ["acarbd", "start", "--pruning", "nothing", "--log_format", "json"]
