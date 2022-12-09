### build stage
FROM golang:1.19.4-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main /app/main.go

### run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY config.json.example config.json
COPY config/cube.glb .
COPY config/cube.gltf .
EXPOSE 8888
CMD [ "/app/main" ]
