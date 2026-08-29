FROM ghcr.io/getzola/zola:v0.22.1 AS build
WORKDIR /app
COPY . .
RUN ["zola", "build"]

FROM nginx:alpine@sha256:db35bfc6b2951e7f8a72db5db120288c127ffaeeb4a6d4b95a26fead017d5913
COPY --from=build /app/public /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
