FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o api ./cmd/api
FROM scratch
COPY --from=builder /build/api /app/
COPY ./migrations /app/migrations
WORKDIR /app
CMD ["./api"]
