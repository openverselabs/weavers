#!/bin/bash

CYAN='\033[0;36m'
GREEN='\033[1;32m'
BLUE='\033[1;34m'
RED='\033[1;31m'
NC='\033[0m'

echo -e "${CYAN}"
echo "        Weavers v0.1.0 Installation Script        "
echo -e "${NC}"
echo -e "${GREEN}[+] Starting installation of Weavers...${NC}\n"

if ! command -v go &> /dev/null; then
    echo -e "${RED}[!] Error: Golang is not installed.${NC}"
    exit 1
fi

echo -e "${BLUE}[*] Cloning repository...${NC}"
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR" || exit
git clone -q https://github.com/openverselabs/weavers.git . || exit

echo -e "${BLUE}[*] Initializing Go modules & fetching dependencies...${NC}"
if [ ! -f "go.mod" ]; then
    go mod init github.com/openverselabs/weavers &> /dev/null
fi
go get github.com/likexian/whois &> /dev/null
go mod tidy &> /dev/null

echo -e "${BLUE}[*] Building Weavers...${NC}"
go build -ldflags="-s -w" -o weavers main.go

if [ -f "weavers" ]; then
    echo -e "${BLUE}[*] Moving binary to /usr/local/bin...${NC}"
    sudo mv weavers /usr/local/bin/
    echo -e "\n${GREEN}[✔] Installation successful!${NC}"
else
    echo -e "${RED}[!] Build failed. 'weavers' binary not found.${NC}"
    exit 1
fi

cd / || exit
rm -rf "$TMP_DIR"
