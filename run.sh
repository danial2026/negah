#!/bin/bash

# The Watchman - Builder & Runner

# Colors for flair
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"
echo "      The Watchman Environment Check        "
echo "~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~"

# 1. Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}[ERROR] Go is not installed.${NC} Please install Go 1.21 or higher."
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${GREEN}[OK] Go found: ${GO_VERSION}${NC}"

# 2. Check for Nmap
if ! command -v nmap &> /dev/null; then
    echo -e "${RED}[ERROR] Nmap is not installed.${NC}"
    echo "Install it via:"
    echo "  macOS: brew install nmap"
    echo "  Arch:  sudo pacman -S nmap"
    exit 1
fi

NMAP_VERSION=$(nmap --version | head -n 1 | awk '{print $3}')
echo -e "${GREEN}[OK] Nmap found: ${NMAP_VERSION}${NC}"

# 3. Clean and Build
echo -e "\nBuilding The Watchman..."
go mod tidy
go build -o nscanner .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}[SUCCESS] Build complete.${NC}\n"
    # 4. Run it
    ./nscanner
else
    echo -e "${RED}[ERROR] Build failed.${NC}"
    exit 1
fi
