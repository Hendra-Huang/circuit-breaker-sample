FROM golang:1.8-alpine
WORKDIR /go/src/github.com/Hendra-Huang/circuit-breaker-sample/C
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o c .

FROM alpine:3.6
COPY --from=0 /go/src/github.com/Hendra-Huang/circuit-breaker-sample/C/c /bin/c
RUN chmod +x /bin/c
EXPOSE 3333
ENTRYPOINT [ "/bin/c" ]
CMD ["-delay", "500"]
