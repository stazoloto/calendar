FROM golang:1.26-alpine AS builder

WORKDIR /calendar

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o calendar ./cmd/server

FROM scratch

COPY --from=builder /calendar/calendar /calendar

EXPOSE 3333

ENTRYPOINT [ "/calendar" ]