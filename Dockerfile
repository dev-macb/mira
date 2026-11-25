FROM node:24-bookworm-slim AS web
WORKDIR /src
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:1.26.11-bookworm AS server
WORKDIR /src
COPY server/go.mod server/go.sum ./
RUN go mod download
COPY server/ ./
RUN CGO_ENABLED=0 go build -trimpath -o /out/ ./cmd/...

FROM alpine:3.23 AS docker-inventory
COPY --from=server /out/docker /usr/local/bin/mira-docker
USER 10001:10001
ENTRYPOINT ["mira-docker"]

FROM alpine:3.23 AS http
RUN install -d -m 0750 -o 10001 -g 10001 /srv/data/mira
COPY --from=server /out/api /usr/local/bin/mira-api
COPY --from=web /src/dist /app/web
ENV MIRA_HTTP_ADDR=0.0.0.0:6002 \
    MIRA_WEB_DIR=/app/web \
    MIRA_DB_PATH=/srv/data/mira/mira.db
USER 10001:10001
ENTRYPOINT ["mira-api"]
