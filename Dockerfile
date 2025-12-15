FROM alpine:latest

RUN addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY tg-questions-bot /app/tg-questions-bot
RUN chmod +x /app/tg-questions-bot && chown app:app /app/tg-questions-bot

EXPOSE 1443

USER app
ENTRYPOINT ["/app/tg-questions-bot"]
