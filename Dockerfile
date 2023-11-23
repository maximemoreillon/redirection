FROM golang:1.21.4
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping
EXPOSE 7070
CMD ["/docker-gs-ping"]