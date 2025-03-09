FROM golang:1.24

WORKDIR /urlShortener

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY cmd/urlShortener/*.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /my_app

CMD ["/my_app"]