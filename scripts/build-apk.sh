#!/usr/bin/env bash
# build-apk.sh — Build a complete Android APK from frontend + Go backend.
#
# Prerequisites:
#   - Node.js 18+ and npm
#   - Go 1.25+ with gomobile installed
#   - Android SDK/NDK (ANDROID_HOME set)
#   - JDK (JAVA_HOME set, or Android Studio JBR will be auto-detected)
#
# Usage:
#   ./scripts/build-apk.sh           # debug APK (default)
#   ./scripts/build-apk.sh release   # release APK (requires signing config)
#
# Output:
#   android/app/build/outputs/apk/debug/app-debug.apk

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BUILD_TYPE="${1:-debug}"

# ---------------------------------------------------------------------------
# Prerequisite checks
# ---------------------------------------------------------------------------
check_prereq() {
  if ! command -v "$1" &>/dev/null; then
    echo "ERROR: '$1' not found. $2"
    exit 1
  fi
}

check_prereq node   "Install Node.js 18+ from https://nodejs.org"
check_prereq npm    "Install Node.js 18+ from https://nodejs.org"
check_prereq go     "Install Go 1.25+ from https://go.dev/dl"
check_prereq gomobile "Run: go install golang.org/x/mobile/cmd/gomobile@latest && gomobile init"

if [ -z "${ANDROID_HOME:-}" ]; then
  # Try common default locations
  if [ -d "$HOME/AppData/Local/Android/Sdk" ]; then
    export ANDROID_HOME="$HOME/AppData/Local/Android/Sdk"
  elif [ -d "$HOME/Library/Android/sdk" ]; then
    export ANDROID_HOME="$HOME/Library/Android/sdk"
  elif [ -d "/usr/local/lib/android/sdk" ]; then
    export ANDROID_HOME="/usr/local/lib/android/sdk"
  else
    echo "ERROR: ANDROID_HOME not set and Android SDK not found in default locations."
    echo "       Install Android Studio or set ANDROID_HOME manually."
    exit 1
  fi
fi
echo "==> ANDROID_HOME: $ANDROID_HOME"

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
  echo "==> JAVA_HOME: $JAVA_HOME"
fi

echo ""
echo "=========================================="
echo "  Otis APK Build ($BUILD_TYPE)"
echo "=========================================="
echo ""

# ---------------------------------------------------------------------------
# Step 1: Build frontend
# ---------------------------------------------------------------------------
echo "==> [1/4] Building frontend..."
cd "$PROJECT_ROOT/frontend"
npm ci --silent 2>/dev/null || npm install --silent
npm run build
echo "==> Frontend built: frontend/dist/"

# ---------------------------------------------------------------------------
# Step 2: Build Go backend as .aar
# ---------------------------------------------------------------------------
echo ""
echo "==> [2/4] Building Go backend (.aar)..."
cd "$PROJECT_ROOT"
bash scripts/build-android.sh

# ---------------------------------------------------------------------------
# Step 3: Capacitor sync
# ---------------------------------------------------------------------------
echo ""
echo "==> [3/4] Syncing Capacitor..."
cd "$PROJECT_ROOT/frontend"
npx cap sync android
echo "==> Capacitor synced."

# ---------------------------------------------------------------------------
# Step 4: Build APK with Gradle
# ---------------------------------------------------------------------------
echo ""
if [ "$BUILD_TYPE" = "release" ]; then
  echo "==> [4/4] Building release APK..."
  cd "$PROJECT_ROOT/android"
  ./gradlew assembleRelease
  APK_PATH="$PROJECT_ROOT/android/app/build/outputs/apk/release/app-release.apk"
else
  echo "==> [4/4] Building debug APK..."
  cd "$PROJECT_ROOT/android"
  ./gradlew assembleDebug
  APK_PATH="$PROJECT_ROOT/android/app/build/outputs/apk/debug/app-debug.apk"
fi

# ---------------------------------------------------------------------------
# Done
# ---------------------------------------------------------------------------
echo ""
echo "=========================================="
echo "  Build complete!"
echo "=========================================="
if [ -f "$APK_PATH" ]; then
  APK_SIZE=$(du -h "$APK_PATH" | cut -f1)
  echo "  APK: $APK_PATH ($APK_SIZE)"
  echo ""
  echo "  Install on device:"
  echo "    adb install $APK_PATH"
else
  echo "  WARNING: APK not found at expected path: $APK_PATH"
  echo "  Check Gradle output above for errors."
  exit 1
fi
