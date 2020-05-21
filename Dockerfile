FROM golang:alpine

WORKDIR /go/src/moderari
COPY . .

RUN go install -v ./cmd/bot

CMD ["moderari"]
