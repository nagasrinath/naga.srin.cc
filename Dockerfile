FROM ghcr.io/getzola/zola:v0.21.0 AS build
WORKDIR /app
COPY . .
RUN ["zola", "build"]

FROM nginx:alpine
COPY --from=build /app/public /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
