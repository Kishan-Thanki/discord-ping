#!/bin/bash
# scripts/install-hooks.sh

HOOK_DIR=".git/hooks"
HOOK_FILE="$HOOK_DIR/pre-push"

echo "Installing Git pre-push hook..."

mkdir -p "$HOOK_DIR"

cat << 'EOF' > "$HOOK_FILE"
#!/bin/bash
# Pre-push hook to ensure code is vetted and tested before pushing.

echo "Running Pre-Check (Lint, Vet, and Test)..."
make pre-check

if [ $? -ne 0 ]; then
  echo "Pre-check failed! Push aborted."
  echo "Please fix the errors above and try pushing again."
  exit 1
fi

echo "Pre-check passed! Pushing to GitHub..."
exit 0
EOF

chmod +x "$HOOK_FILE"

echo "Git pre-push hook installed successfully!"
