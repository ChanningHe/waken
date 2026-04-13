FROM alpine:3.21
ARG TARGETARCH
RUN apk add --no-cache ca-certificates tzdata
COPY waken-linux-${TARGETARCH} /usr/local/bin/waken
RUN chmod +x /usr/local/bin/waken
EXPOSE 19527
VOLUME /app/waken/config
ENTRYPOINT ["waken"]
