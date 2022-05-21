FROM golang:1.14-alpine
COPY website-controller /app
COPY template /app/template/
WORKDIR /app
CMD ["/app/website-controller"]