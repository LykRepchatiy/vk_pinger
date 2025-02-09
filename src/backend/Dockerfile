FROM golang:1.22.2

WORKDIR /app

COPY . .

RUN go build -o backend .

EXPOSE 8080

CMD ["./backend"]