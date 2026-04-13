# Stage 1: Build frontend
FROM node:22-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.25-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /waken .

# Stage 3: Minimal runtime
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
COPY --from=backend /waken /usr/local/bin/waken
EXPOSE 19527
VOLUME /app/waken/config
ENTRYPOINT ["waken"]
