FROM golang:1.21.4

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o github.com/paulmuenzner/powerplantmanager

EXPOSE 8080

CMD ["./github.com/paulmuenzner/powerplantmanager"]