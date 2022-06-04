FROM node:18-alpine AS builder

WORKDIR /build/

COPY frontend/package.json .
COPY frontend/package-lock.json .

RUN npm ci

COPY frontend/public/ public/
COPY frontend/src/ src/

RUN npm run build