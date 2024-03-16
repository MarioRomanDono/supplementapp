FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o supplementapp ./cmd/http-server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/supplementapp .
EXPOSE 8080
RUN adduser -D noroot
USER noroot:noroot
ENTRYPOINT ["./supplementapp"]