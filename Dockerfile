# Build stage
FROM golang:1.19 AS builder
ENV CGO_ENABLED 0
WORKDIR /work
COPY . .
RUN go get -d -v .
RUN go build -o /dist/dapr-ipfs-binding -v .

# Final stage
FROM gcr.io/distroless/static-debian11
COPY --from=builder /dist/dapr-ipfs-binding /
CMD ["/dapr-ipfs-binding"]
