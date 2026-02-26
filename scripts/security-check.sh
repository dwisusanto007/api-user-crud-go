#!/bin/bash

# Security Check Script
# Checks for potential security issues before committing

set -e

echo "🔒 Running Security Checks..."
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ERRORS=0
WARNINGS=0

# Check 1: .env files
echo "📋 Checking for .env files..."
if git ls-files | grep -q "\.env$"; then
    echo -e "${RED}❌ FAIL: .env file is tracked by git!${NC}"
    echo "   Run: git rm --cached .env"
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✅ PASS: No .env files in git${NC}"
fi
echo ""

# Check 2: Database files
echo "📋 Checking for database files..."
if git ls-files | grep -q "\.db$"; then
    echo -e "${RED}❌ FAIL: Database files are tracked by git!${NC}"
    echo "   Run: git rm --cached *.db"
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✅ PASS: No database files in git${NC}"
fi
echo ""

# Check 3: Hardcoded secrets (in staged files)
echo "📋 Checking for potential hardcoded secrets..."
if git diff --cached | grep -iE "secret.*=.*['\"][^'\"]{20,}['\"]|password.*=.*['\"][^'\"]{8,}['\"]" | grep -v "example\|default-secret\|password123" > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  WARNING: Potential hardcoded secrets found in staged files${NC}"
    echo "   Review your changes carefully"
    WARNINGS=$((WARNINGS + 1))
else
    echo -e "${GREEN}✅ PASS: No obvious hardcoded secrets${NC}"
fi
echo ""

# Check 4: JWT_SECRET in code
echo "📋 Checking for JWT_SECRET usage..."
if git diff --cached | grep -E "JWT_SECRET.*=.*['\"][^'\"]{20,}['\"]" | grep -v "getEnv\|os.Getenv\|example\|default-secret" > /dev/null 2>&1; then
    echo -e "${RED}❌ FAIL: JWT_SECRET appears to be hardcoded!${NC}"
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✅ PASS: JWT_SECRET properly handled${NC}"
fi
echo ""

# Check 5: Binary files
echo "📋 Checking for binary files..."
if git ls-files | grep -qE "api-user-crud-go$|\.exe$|\.dll$|\.so$"; then
    echo -e "${RED}❌ FAIL: Binary files are tracked by git!${NC}"
    echo "   Run: git rm --cached <binary-file>"
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✅ PASS: No binary files in git${NC}"
fi
echo ""

# Check 6: Log files
echo "📋 Checking for log files..."
if git ls-files | grep -q "\.log$"; then
    echo -e "${RED}❌ FAIL: Log files are tracked by git!${NC}"
    echo "   Run: git rm --cached *.log"
    ERRORS=$((ERRORS + 1))
else
    echo -e "${GREEN}✅ PASS: No log files in git${NC}"
fi
echo ""

# Check 7: .gitignore exists and has required entries
echo "📋 Checking .gitignore..."
if [ ! -f .gitignore ]; then
    echo -e "${RED}❌ FAIL: .gitignore file not found!${NC}"
    ERRORS=$((ERRORS + 1))
else
    REQUIRED_ENTRIES=(".env" "*.db" "*.log" "data/")
    MISSING=0
    for entry in "${REQUIRED_ENTRIES[@]}"; do
        if ! grep -q "^${entry}$" .gitignore; then
            echo -e "${YELLOW}⚠️  WARNING: .gitignore missing entry: ${entry}${NC}"
            MISSING=$((MISSING + 1))
        fi
    done
    
    if [ $MISSING -eq 0 ]; then
        echo -e "${GREEN}✅ PASS: .gitignore properly configured${NC}"
    else
        WARNINGS=$((WARNINGS + 1))
    fi
fi
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Security Check Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "Errors:   ${RED}${ERRORS}${NC}"
echo -e "Warnings: ${YELLOW}${WARNINGS}${NC}"
echo ""

if [ $ERRORS -gt 0 ]; then
    echo -e "${RED}❌ Security check FAILED!${NC}"
    echo "   Please fix the errors above before committing."
    exit 1
elif [ $WARNINGS -gt 0 ]; then
    echo -e "${YELLOW}⚠️  Security check passed with warnings${NC}"
    echo "   Review the warnings above."
    exit 0
else
    echo -e "${GREEN}✅ All security checks passed!${NC}"
    echo "   Safe to commit."
    exit 0
fi
