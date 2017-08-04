FROM golang:1.8-alpine
WORKDIR /go/src/github.com/Hendra-Huang/circuit-breaker-sample/B
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o b .

FROM alpine:3.6
COPY --from=0 /go/src/github.com/Hendra-Huang/circuit-breaker-sample/B/b /bin/b
RUN chmod +x /bin/b
EXPOSE 2222
ENTRYPOINT [ "/bin/b" ]
