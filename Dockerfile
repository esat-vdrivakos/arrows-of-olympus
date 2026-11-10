# Build the game (WASM) and the static server.
FROM golang:1.27 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=js GOARCH=wasm go build -trimpath -ldflags="-s -w" -o /out/public/game.wasm ./cmd/game \
 && gzip -9 -k /out/public/game.wasm \
 && cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" /out/public/ \
 && cp public/index.html /out/public/ \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/ /app/
ENV PUBLIC_DIR=/app/public
EXPOSE 8080
ENTRYPOINT ["/app/server"]
