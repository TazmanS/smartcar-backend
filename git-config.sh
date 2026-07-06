#!/bin/sh

set -e

mkdir -p .githooks

cat > .githooks/pre-commit <<'EOF'
#!/bin/sh

echo "Generating Swagger..."

swag init -g main.go

git add docs/

echo "Done."
EOF

chmod +x .githooks/pre-commit

git config core.hooksPath .githooks