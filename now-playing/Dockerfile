FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o now-playing .

FROM alpine:3.20
RUN adduser -D -H app
USER app
COPY --from=build /app/now-playing /now-playing
EXPOSE 8080
ENTRYPOINT ["/now-playing"]
