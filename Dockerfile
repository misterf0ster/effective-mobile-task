FROM golang:latest

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN go build -o effective-mobile-app

CMD ["/app/cmd/effective-mobile-app"]