# ---- Build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/homepage-service ./cmd/server

# ---- Runtime stage ----
FROM alpine:3.20
RUN adduser -D -u 10001 appuser
COPY --from=build /out/homepage-service /usr/local/bin/homepage-service
USER appuser
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/homepage-service"]
