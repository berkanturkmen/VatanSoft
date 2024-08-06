FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY .env .
COPY --from=builder /app/myapp .

CMD ["./myapp"]