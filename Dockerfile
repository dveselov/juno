FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -tags netgo -a -v -o api cmd/api
FROM scratch
COPY --from=builder /build/api /app/
ADD ./migrations /app
WORKDIR /app
CMD ["./api"]
