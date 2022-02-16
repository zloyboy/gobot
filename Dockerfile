FROM golang:latest

WORKDIR /home

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build github.com/zloyboy/gobot/cmd/gobot

CMD [ "./gobot" ]