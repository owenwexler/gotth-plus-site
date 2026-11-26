# syntax=docker/dockerfile:1
# Production image. The CSS and the server are built in the Go image, then copied into a minimal
# distroless image that runs as a non-root user.
#
#   make docker            # docker build -t gotth-plus-site .
#   docker run --rm --entrypoint /app/migrate -v gotth-plus-site-data:/data gotth-plus-site   # schema, before each start
#   docker run --rm -p 8080:8080 -v gotth-plus-site-data:/data gotth-plus-site                # the server
#
# DEPLOYMENT.md has the full procedure.

FROM golang:1.27.4-trixie AS build
WORKDIR /src

# The Tailwind CSS standalone CLI, pinned and checked against its release checksum.
RUN set -eu; \
    case "$(uname -m)" in \
      x86_64) arch=x64; sum=149a8ece5b855587a9c10de6da9dcbcfad0e47c4ed5dca219f721b3eae39d817 ;; \
      aarch64) arch=arm64; sum=f6d7baae1d350ed3a2defe98d0b5cc58e37400c0173b2e38d42f690ccc4bbfe1 ;; \
      *) echo "unsupported architecture: $(uname -m)" >&2; exit 1 ;; \
    esac; \
    curl -fsSL -o /usr/local/bin/tailwindcss \
      "https://github.com/tailwindlabs/tailwindcss/releases/download/v4.3.4/tailwindcss-linux-${arch}"; \
    echo "${sum}  /usr/local/bin/tailwindcss" | sha256sum -c -; \
    chmod +x /usr/local/bin/tailwindcss

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Production serves the committed, reviewed copies of the browser libraries and never downloads
# them (SECURITY.md, "Frontend assets").
RUN for f in htmx.min.js htmx-ext-preload.min.js alpine-csp.min.js hx-optimistic.min.js; do \
      test -f "static/js/$f" || { echo "static/js/$f is missing: run make js, review and commit static/js" >&2; exit 1; }; \
    done

RUN tailwindcss -i static/input.css -o static/output.css --minify
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate
# the SQLite database lives on a volume, so the rest of the filesystem can be read-only
RUN mkdir -p /out/data

FROM gcr.io/distroless/static-debian13:nonroot
WORKDIR /app
COPY --from=build /out/server /app/server
COPY --from=build /out/migrate /app/migrate
COPY --from=build /src/static /app/static
COPY --from=build --chown=65532:65532 /out/data /data
ENV GoEnv=production Port=8080 DatabasePath=/data/app.db
VOLUME /data
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app/server"]
