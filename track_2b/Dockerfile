FROM node:24.15.0-bookworm-slim@sha256:4e6b70dd6cbfc88c8157ba19aa3d9f9cce6ba4703576d55459e45efcbc9c5f5d AS frontend
WORKDIR /build/fe
RUN npm install --global pnpm@11.1.3
COPY src/laporsiap_FE/package.json src/laporsiap_FE/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY src/laporsiap_FE/ ./
RUN pnpm build

FROM golang:1.26.9-alpine@sha256:3082400e369fa24d5fc60bca20edab3f6d604e0c5a690ec66b295eff4dd87ade AS backend
WORKDIR /build/be
COPY src/laporsiap_BE/go.mod src/laporsiap_BE/go.sum ./
RUN go mod download
COPY src/laporsiap_BE/ ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/api ./cmd/api
RUN CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/maintenance ./cmd/maintenance

FROM alpine:3.23@sha256:85fe1e81d6758c208f3e1eed4338a1997e19d4be002d4dd32d3100c9a8c010a0
RUN apk add --no-cache ca-certificates && addgroup -S -g 65532 app && adduser -S -D -H -u 65532 -G app app
WORKDIR /app
COPY --from=backend /out/api /out/maintenance /app/
COPY --from=frontend /build/fe/dist /app/frontend
USER 65532:65532
ENV LISTEN_ADDR=0.0.0.0:8787 ALLOW_CONTAINER_BIND=true FRONTEND_DIR=/app/frontend
EXPOSE 8787
HEALTHCHECK --interval=10s --timeout=3s --start-period=20s --retries=6 CMD wget -q -O /dev/null http://127.0.0.1:8787/api/health || exit 1
ENTRYPOINT ["/app/api"]
