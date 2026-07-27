# --- build ---
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/prepbot ./cmd/prepbot

# --- run ---
FROM alpine:3.20
# ca-certificates is required for TLS to api.telegram.org / api.anthropic.com.
RUN apk add --no-cache ca-certificates && adduser -D -H prepbot
WORKDIR /app
COPY --from=build /out/prepbot /usr/local/bin/prepbot
COPY data ./data
COPY prompts ./prompts
USER prepbot
ENTRYPOINT ["prepbot"]
