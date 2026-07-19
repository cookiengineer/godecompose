#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
OUT_DIR="$SCRIPT_DIR/bin"

VERSION="$(git describe --tags --always --dirty 2>/dev/null || echo "dev")"
COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")"

TARGETS=(
	"linux/amd64"
	"linux/arm64"
	"darwin/amd64"
	"darwin/arm64"
)

rm -rf "$OUT_DIR"
mkdir -p "$OUT_DIR"

for target in "${TARGETS[@]}"; do
	GOOS="${target%/*}"
	GOARCH="${target#*/}"
	OUTPUT="$OUT_DIR/godecompose_${GOOS}_${GOARCH}"

	echo "building $GOOS/$GOARCH -> $OUTPUT"

	CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
		go build \
		-ldflags="-s -w -X main.version=$VERSION -X main.commit=$COMMIT" \
		-trimpath \
		-o "$OUTPUT" \
		"$SCRIPT_DIR/cmd/godecompose"
done

echo ""
echo "done: $(ls "$OUT_DIR" | wc -l) binaries in $OUT_DIR"
echo ""
ls -lh "$OUT_DIR/"
