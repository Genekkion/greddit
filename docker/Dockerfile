FROM golang:1.25.4-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ cmd/

COPY internal/ internal/

RUN CGO_ENABLED=0 \
    GOOS=linux \
    go build \
    -ldflags="-w -s" \
    -o greddit \
    cmd/main.go

FROM scratch

WORKDIR /app
COPY --from=builder /app/greddit .

ENTRYPOINT ["/app/greddit"]
