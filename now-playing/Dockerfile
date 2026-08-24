FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o now-playing .

FROM cgr.dev/chainguard/static:latest
COPY --from=build /app/now-playing /now-playing
EXPOSE 8080
ENTRYPOINT ["/now-playing"]
