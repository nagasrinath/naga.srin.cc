FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -o now-playing .

FROM cgr.dev/chainguard/static:latest@sha256:f68e3a8244c7d0f4cd56635aaff8e6a533cf6cc3850d8fb339567a5782d6a0b0
COPY --from=build /app/now-playing /now-playing
EXPOSE 8080
ENTRYPOINT ["/now-playing"]
