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
    echo -e "Please install Go first (https://go.dev/doc/install) and try again."
    exit 1
fi

if ! command -v git &> /dev/null; then
    echo -e "${RED}[!] Error: Git is not installed.${NC}"
    exit 1
fi

echo -e "${BLUE}[*] Cloning repository...${NC}"
TMP_DIR=$(mktemp -d)
cd "$TMP_DIR" || exit
git clone -q https://github.com/openverselabs/weavers.git
cd weavers || exit

echo -e "${BLUE}[*] Resolving dependencies and building Weavers...${NC}"
go mod tidy &> /dev/null
go build -ldflags="-s -w" -o weavers main.go

echo -e "${BLUE}[*] Moving binary to /usr/local/bin (may require sudo password)...${NC}"
if sudo mv weavers /usr/local/bin/; then
    echo -e "\n${GREEN}[✔] Installation successful!${NC}"
else
    echo -e "${RED}[!] Failed to move binary to /usr/local/bin.${NC}"
    exit 1
fi

cd / || exit
rm -rf "$TMP_DIR"

echo -e "You can now run ${CYAN}weavers${NC} from anywhere in your terminal."
echo -e "Try it: ${CYAN}weavers -d google.com${NC}"
