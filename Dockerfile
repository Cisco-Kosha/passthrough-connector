FROM docker.io/golang:alpine as builder
RUN apk add git
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN rm /build/go.sum
RUN go mod tidy
RUN go build -o passthrough-connector .
FROM docker.io/alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/passthrough-connector /app/
WORKDIR /app
CMD ["./passthrough-connector"]