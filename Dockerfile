FROM golang:1.14-alpine

RUN apk update && apk add curl

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o website-controller .
CMD ["./website-controller"]