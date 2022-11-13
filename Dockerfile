FROM golang:1.17

COPY ./ ./
ENV GOPATH=/

RUN go mod download
RUN go build -o balance-service ./cmd/main.go

EXPOSE 7000

CMD ["./balance-service"]