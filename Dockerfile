FROM node:18-alpine AS frontend-builder

WORKDIR /build/

COPY frontend/package.json .
COPY frontend/package-lock.json .

RUN npm ci

COPY frontend/public/ public/
COPY frontend/src/ src/

# frontend files will be located in /build/build/...
RUN npm run build

FROM golang:1.18-alpine AS backend-base

RUN apk add --update --no-cache bash musl-dev make git && adduser -D -g '' appuser

WORKDIR /build/
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

FROM backend-base AS backend-builder

ARG VERSION=NA
ARG COMMIT=NA
ARG DATE=NA

WORKDIR /build/
COPY go.mod go.mod
COPY go.sum go.sum
COPY backend/ /build/backend

RUN CGO_ENABLED=0 go build \
    -ldflags="-X 'main.version=$VERSION' -X 'main.commit=$COMMIT' -X 'main.date=$DATE' -s -w" \
    -o bin/todo-server ./backend/cmd/todo-server

FROM scratch
VOLUME /data
WORKDIR /app/
COPY --from=backend-builder /build/bin /app
COPY --from=backend-builder /build/backend/cmd/todo-server/config.example.json /app/config.json
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=backend-builder /etc/passwd /etc/passwd
COPY --from=frontend-builder /build/build /app/assets
USER appuser
CMD ["/app/todo-server", "--assets", "/app/assets"]