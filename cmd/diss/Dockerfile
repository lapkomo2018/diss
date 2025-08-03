FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/diss

FROM gcr.io/distroless/static
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]
