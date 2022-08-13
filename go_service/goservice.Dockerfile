
FROM fedora
RUN dnf install git -y

FROM golang:1.16-alpine

RUN apk add --no-cache git

ENV GO111MODULE=on
ENV GOPATH /go
WORKDIR /go_service/go_app

COPY go.mod .

RUN go mod download

COPY . .

COPY .env .

RUN  go build -o main .


EXPOSE 3001

# ENTRYPOINT [ "main" ]

CMD ["go", "run", "main.go"]
