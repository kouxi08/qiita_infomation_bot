FROM golang:1.18-alpine3.15 AS go
WORKDIR /app
ADD go.mod go.sum main.go ./
ADD config/local.env ./
RUN go mod download
RUN go build -o main /app/main.go
CMD /app/main