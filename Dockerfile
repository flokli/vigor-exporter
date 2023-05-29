FROM golang:1.19

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY vigor/* ./vigor/
RUN CGO_ENABLED=0 GOOS=linux go build -o /vigor-exporter

EXPOSE 9103

# Run
ENTRYPOINT ["/vigor-exporter"]