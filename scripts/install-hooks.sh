#!/bin/bash

# Install Git Hooks Script
# This script installs pre-commit hooks for security checks

set -e

echo "🔧 Installing Git Hooks..."
echo ""

# Check if .git directory exists
if [ ! -d .git ]; then
    echo "❌ Error: Not a git repository"
    echo "   Run this script from the root of your git repository"
    exit 1
fi

# Create pre-commit hook
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash

# Pre-commit hook for security checks
echo "🔒 Running pre-commit security checks..."
echo ""

# Run security check script
if [ -f scripts/security-check.sh ]; then
    ./scripts/security-check.sh
    EXIT_CODE=$?
    
    if [ $EXIT_CODE -ne 0 ]; then
        echo ""
        echo "❌ Commit aborted due to security check failures"
        echo "   Fix the issues above or use 'git commit --no-verify' to skip (not recommended)"
        exit 1
    fi
else
    echo "⚠️  Warning: security-check.sh not found, skipping security checks"
fi

echo ""
echo "✅ Pre-commit checks passed, proceeding with commit..."
exit 0
EOF

# Make hook executable
chmod +x .git/hooks/pre-commit

echo "✅ Git hooks installed successfully!"
echo ""
echo "Pre-commit hook will now run security checks before each commit."
echo "To skip the hook (not recommended), use: git commit --no-verify"
echo ""
