# --- Build do front-end Angular ---
FROM node:24-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json frontend/.npmrc ./
RUN npm ci
COPY frontend/ ./
RUN npx ng build

# --- Build do backend Go ---
FROM golang:1.24-alpine AS backend
WORKDIR /src
COPY backend/go.mod ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /server ./cmd/server

# --- Imagem final ---
FROM alpine:3.20
WORKDIR /app
COPY --from=backend /server /app/server
COPY --from=frontend /app/dist/frontend/browser /app/public
ENV ADDR=:8080 \
    STATIC_DIR=/app/public \
    DATA_FILE=/app/data/tasks.json \
    CORS_ORIGIN=
VOLUME /app/data
EXPOSE 8080
CMD ["/app/server"]
