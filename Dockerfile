FROM golang as builder

WORKDIR /app
COPY . .
RUN go build -o server .

FROM debian
COPY --from=builder /app/server /usr/local/bin/

EXPOSE 8080
CMD [ "server" ]
