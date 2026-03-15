#!/usr/bin/env bash
# build-android.sh — Compile the Go backend as an Android .aar using gomobile.
#
# Prerequisites:
#   go install golang.org/x/mobile/cmd/gomobile@latest
#   gomobile init
#   Android SDK/NDK installed (ANDROID_HOME set)
#   JDK installed (JAVA_HOME set, or Android Studio JBR will be auto-detected)
#
# Usage:
#   ./scripts/build-android.sh
#
# Output:
#   android/app/libs/otis.aar
#   android/app/libs/otis-sources.jar

set -euo pipefail

# Auto-detect JAVA_HOME from Android Studio JBR if not set
if [ -z "${JAVA_HOME:-}" ]; then
  if [ -d "/c/Program Files/Android/Android Studio/jbr" ]; then
    export JAVA_HOME="/c/Program Files/Android/Android Studio/jbr"
  elif [ -d "$HOME/Library/Java/JavaVirtualMachines" ]; then
    JAVA_HOME="$(/usr/libexec/java_home 2>/dev/null || true)"
  fi
fi

if [ -n "${JAVA_HOME:-}" ]; then
  export PATH="${JAVA_HOME}/bin:${PATH}"
fi

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
OUTPUT_DIR="${PROJECT_ROOT}/android/app/libs"

mkdir -p "$OUTPUT_DIR"

echo "==> Building Android .aar with gomobile..."
cd "$PROJECT_ROOT"

gomobile bind \
  -target=android \
  -androidapi=21 \
  -o "${OUTPUT_DIR}/otis.aar" \
  ./app/mobile

AAR_SIZE=$(du -h "${OUTPUT_DIR}/otis.aar" | cut -f1)
echo "==> Build complete: ${OUTPUT_DIR}/otis.aar (${AAR_SIZE})"
echo "==> Exported functions: Start(port int, dataDir string) error, Stop() error"
