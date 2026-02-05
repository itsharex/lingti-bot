#!/bin/bash
set -e

# Release script for lingti-bot
# Usage: ./scripts/release.sh [--bot] <version>
# Example: ./scripts/release.sh 1.2.2
#          ./scripts/release.sh --bot 1.2.2

PROJECTNAME="lingti-bot"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# Parse args: skip flags, take the first positional arg as version
VERSION=""
for arg in "$@"; do
    case "$arg" in
        --*|-*) ;;  # ignore flags
        *) VERSION="$arg" ;;
    esac
done

if [ -z "$VERSION" ]; then
    echo "Usage: $0 [--bot] <version>"
    echo "Example: $0 1.2.2"
    exit 1
fi

if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
    echo "Error: Version must be in semver format (e.g., 1.2.0, 1.2.2-beta)"
    exit 1
fi

# Helper to extract Go string value
extract_go_str() {
    grep "$1" "$2" 2>/dev/null | head -1 | sed 's/.*"\([^"]*\)".*/\1/'
}

# Gather and display current versions
v_makefile=$(grep -E "^VERSION\s*:=" "$PROJECT_DIR/Makefile" | sed 's/VERSION[[:space:]]*:=[[:space:]]*//' | tr -d '[:space:]')
v_server=$(extract_go_str 'ServerVersion' "$PROJECT_DIR/internal/mcp/server.go")
v_client=$(extract_go_str 'ClientVersion' "$PROJECT_DIR/internal/platforms/relay/relay.go")

echo ""
echo "  Current versions:"
echo "    Makefile         VERSION := $v_makefile"
echo "    mcp/server.go    ServerVersion = $v_server"
echo "    relay/relay.go   ClientVersion = $v_client"
echo ""
echo "  New version: $VERSION"
echo ""

read -p "  Proceed with release? [y/N] " -n 1 -r
echo
[[ $REPLY =~ ^[Yy]$ ]] || { echo "Cancelled."; exit 0; }

echo "==> Updating version strings to $VERSION"

sed -i '' "s/^VERSION[[:space:]]*:=.*/VERSION := $VERSION/" "$PROJECT_DIR/Makefile"
sed -i '' "s/ServerVersion[[:space:]]*=[[:space:]]*\"[^\"]*\"/ServerVersion = \"$VERSION\"/" "$PROJECT_DIR/internal/mcp/server.go"
sed -i '' "s/ClientVersion[[:space:]]*=[[:space:]]*\"[^\"]*\"/ClientVersion     = \"$VERSION\"/" "$PROJECT_DIR/internal/platforms/relay/relay.go"

echo "==> Building release v$VERSION"

rm -rf dist
mkdir -p dist

make all

# Create archives
echo "==> Creating archives..."
cd dist

for f in "${PROJECTNAME}-${VERSION}"-*; do
    case "$f" in
        *.exe) zip -q "$f.zip" "$f" ;;
        *)     tar -czf "$f.tar.gz" "$f" ;;
    esac
done

echo "==> Generating checksums..."
shasum -a 256 *.tar.gz *.zip > checksums.txt

cd ..

# Commit version bump
git add Makefile internal/mcp/server.go internal/platforms/relay/relay.go
git commit -m "chore: bump version to $VERSION" || true

# Create git tag and push
echo "==> Creating git tag v$VERSION..."
git tag -a "v$VERSION" -m "Release v$VERSION"

echo "==> Pushing tag to remote..."
git push origin "v$VERSION"

# Create GitHub release
echo "==> Creating GitHub release..."
gh release create "v$VERSION" \
    --title "v$VERSION" \
    --generate-notes \
    dist/*.tar.gz \
    dist/*.zip \
    dist/checksums.txt

echo "==> Release v$VERSION complete!"
echo "View at: https://github.com/ruilisi/lingti-bot/releases/tag/v$VERSION"
