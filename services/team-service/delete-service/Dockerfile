FROM golang:latest

WORKDIR /app

COPY . .
COPY go.mod go.sum /app/
RUN go mod download
COPY vendor /app/vendor


EXPOSE 8126

CMD [ "go","run","main.go" ]
