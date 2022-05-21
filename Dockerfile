FROM golang:1.14-alpine
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o website-controller .
CMD ["/app/website-controller"]