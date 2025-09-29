# --- build web
FROM node:20 AS web
WORKDIR /repo
COPY apps/web ./apps/web
WORKDIR /repo/apps/web
RUN npm ci && npm run build

# --- build server
FROM golang:1.22 AS server
WORKDIR /repo
COPY go.work ./
COPY apps/server ./apps/server
WORKDIR /repo/apps/server
# copy web build into server context for embedding
COPY --from=web /repo/apps/web/dist ../web/dist
RUN go build -o /netscan-server .

# --- runtime
FROM gcr.io/distroless/base-debian12
COPY --from=server /netscan-server /netscan-server
EXPOSE 8080
ENTRYPOINT ["/netscan-server"]