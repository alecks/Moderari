FROM golang:alpine

WORKDIR /go/src/moderari
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["moderari"]