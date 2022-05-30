FROM golang:1.18-alpine as base

RUN apk add --update --no-cache bash musl-dev make git && adduser -D -g '' appuser

WORKDIR /build/
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

FROM base as builder

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
COPY --from=builder /build/bin /app
COPY --from=builder /build/backend/cmd/todo-server/config.example.json /app/config.json
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/passwd
USER appuser
CMD ["/app/todo-server"]