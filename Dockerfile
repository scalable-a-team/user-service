FROM golang:1.17 AS builder
WORKDIR /build/userservice
RUN mkdir -p temp/images
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o userservice .

FROM alpine:3.15
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /build/userservice .
CMD ["./userservice"]