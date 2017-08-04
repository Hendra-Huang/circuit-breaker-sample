FROM golang:1.8-alpine
WORKDIR /go/src/github.com/Hendra-Huang/circuit-breaker-sample/A
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o a .

FROM alpine:3.6
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/Hendra-Huang/circuit-breaker-sample/a /bin/a
RUN chmod +x /bin/a
EXPOSE 1111
ENTRYPOINT [ "/bin/a" ]
