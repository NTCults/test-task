FROM golang:1.10.0 as builder
WORKDIR /go/src/test-task
COPY . .
RUN go get .
RUN CGO_ENABLED=0 GOOS=linux go build  -o test-task -a -installsuffix cgo .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/test-task .
CMD ["./test-task"]