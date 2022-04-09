FROM golang:1.18.0 as builder
WORKDIR /go/src/github.com/scor2k/go-sigbro-mail-sender/
COPY go.mod go.sum .
RUN go mod download 
COPY *.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o go-sigbro-mail-sender .


FROM alpine:3.15.0
RUN apk --no-cache add ca-certificates
WORKDIR /opt/app
COPY --from=builder /go/src/github.com/scor2k/go-sigbro-mail-sender/go-sigbro-mail-sender .
CMD ["./go-sigbro-mail-sender"]