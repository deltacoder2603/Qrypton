#!/bin/bash

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║  QRYPTON Blockchain Setup Script      ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}\n"

check_go_installation() {
    echo -e "${YELLOW}Checking Go installation...${NC}"
    if ! command -v go &> /dev/null; then
        echo -e "${RED}✗ Go is not installed${NC}"
        echo "Please install Go from https://golang.org/dl/"
        echo "Minimum version: Go 1.21"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✓ Go ${GO_VERSION} found${NC}\n"
}

check_curl() {
    echo -e "${YELLOW}Checking curl installation...${NC}"
    if ! command -v curl &> /dev/null; then
        echo -e "${RED}✗ curl is not installed${NC}"
        echo "Please install curl"
        exit 1
    fi
    echo -e "${GREEN}✓ curl found${NC}\n"
}

check_jq() {
    echo -e "${YELLOW}Checking jq installation...${NC}"
    if ! command -v jq &> /dev/null; then
        echo -e "${YELLOW}⚠ jq is not installed (optional, recommended for pretty JSON output)${NC}"
        echo "Install with: sudo apt-get install jq (Linux) or brew install jq (macOS)"
    else
        echo -e "${GREEN}✓ jq found${NC}"
    fi
    echo ""
}

install_dependencies() {
    echo -e "${YELLOW}Installing Go dependencies...${NC}"
    go mod download
    go mod tidy
    echo -e "${GREEN}✓ Dependencies installed${NC}\n"
}

build_project() {
    echo -e "${YELLOW}Building QRYPTON Blockchain...${NC}"
    go build -o qrypton-blockchain .
    echo -e "${GREEN}✓ Build complete${NC}\n"
}

setup_database() {
    echo -e "${YELLOW}Initializing database...${NC}"
    rm -f blockchain.db
    echo -e "${GREEN}✓ Database ready${NC}\n"
}

create_directories() {
    echo -e "${YELLOW}Creating directories...${NC}"
    mkdir -p data
    mkdir -p logs
    echo -e "${GREEN}✓ Directories created${NC}\n"
}

show_quick_start() {
    echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  Setup Complete! Quick Start Guide     ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════╝${NC}\n"
    
    echo -e "${GREEN}1. Start the blockchain server:${NC}"
    echo -e "   ${YELLOW}go run main.go${NC}"
    echo -e "   or use: ${YELLOW}make run${NC}\n"
    
    echo -e "${GREEN}2. Open the web dashboard:${NC}"
    echo -e "   Open ${YELLOW}index.html${NC} in your browser\n"
    
    echo -e "${GREEN}3. Run the test suite:${NC}"
    echo -e "   ${YELLOW}./test.sh${NC}\n"
    
    echo -e "${GREEN}4. Check available commands:${NC}"
    echo -e "   ${YELLOW}make help${NC}\n"
    
    echo -e "${GREEN}5. View API documentation:${NC}"
    echo -e "   ${YELLOW}make api-docs${NC}\n"
    
    echo -e "${BLUE}Useful Commands:${NC}"
    echo -e "   ${YELLOW}make run${NC}          - Run the server"
    echo -e "   ${YELLOW}make test${NC}         - Run tests"
    echo -e "   ${YELLOW}make dev${NC}          - Development mode"
    echo -e "   ${YELLOW}make clean${NC}        - Clean build artifacts"
    echo -e "   ${YELLOW}make api-docs${NC}     - Show API documentation"
    echo ""
}

show_api_endpoints() {
    echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  API Endpoints Reference               ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════╝${NC}\n"
    
    echo -e "${GREEN}Base URL: http://localhost:8080${NC}\n"
    
    echo -e "${YELLOW}Transaction Operations:${NC}"
    echo -e "  POST /transaction   - Add a transaction"
    echo -e "  GET  /pending       - View pending transactions\n"
    
    echo -e "${YELLOW}Mining Operations:${NC}"
    echo -e "  POST /mine          - Mine a new block\n"
    
    echo -e "${YELLOW}Blockchain Info:${NC}"
    echo -e "  GET  /chain         - Get entire blockchain"
    echo -e "  GET  /block/:index  - Get specific block"
    echo -e "  GET  /stats         - Get blockchain statistics"
    echo -e "  GET  /validate      - Validate blockchain integrity\n"
    
    echo -e "${YELLOW}Account Operations:${NC}"
    echo -e "  GET  /balance       - Check address balance\n"
    
    echo -e "${YELLOW}Monitoring:${NC}"
    echo -e "  GET  /metrics       - Get performance metrics"
    echo -e "  GET  /              - API info\n"
}

show_example_commands() {
    echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  Example API Calls                     ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════╝${NC}\n"
    
    echo -e "${GREEN}Get blockchain:${NC}"
    echo -e "  ${YELLOW}curl http://localhost:8080/chain${NC}\n"
    
    echo -e "${GREEN}Add transaction:${NC}"
    echo -e "  ${YELLOW}curl -X POST http://localhost:8080/transaction \\${NC}"
    echo -e "  ${YELLOW}-H \"Content-Type: application/json\" \\${NC}"
    echo -e "  ${YELLOW}-d '{\"sender\":\"alice\",\"receiver\":\"bob\",\"amount\":50.0}'${NC}\n"
    
    echo -e "${GREEN}Mine block:${NC}"
    echo -e "  ${YELLOW}curl -X POST http://localhost:8080/mine?miner=miner1${NC}\n"
    
    echo -e "${GREEN}Check balance:${NC}"
    echo -e "  ${YELLOW}curl http://localhost:8080/balance?address=alice${NC}\n"
    
    echo -e "${GREEN}Get stats:${NC}"
    echo -e "  ${YELLOW}curl http://localhost:8080/stats${NC}\n"
    
    echo -e "${GREEN}Validate chain:${NC}"
    echo -e "  ${YELLOW}curl http://localhost:8080/validate${NC}\n"
}

show_troubleshooting() {
    echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  Troubleshooting                       ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════╝${NC}\n"
    
    echo -e "${YELLOW}Port 8080 already in use:${NC}"
    echo -e "  1. Find process: ${YELLOW}lsof -i :8080${NC}"
    echo -e "  2. Kill process: ${YELLOW}kill -9 <PID>${NC}"
    echo -e "  3. Or change port in main.go\n"
    
    echo -e "${YELLOW}Database lock error:${NC}"
    echo -e "  ${YELLOW}make clean${NC}\n"
    
    echo -e "${YELLOW}Build error:${NC}"
    echo -e "  ${YELLOW}go clean -modcache${NC}"
    echo -e "  ${YELLOW}go mod download${NC}"
    echo -e "  ${YELLOW}make build${NC}\n"
}

main() {
    check_go_installation
    check_curl
    check_jq
    install_dependencies
    create_directories
    build_project
    setup_database
    
    echo ""
    show_quick_start
    echo ""
    show_api_endpoints
    echo ""
    show_example_commands
    echo ""
    show_troubleshooting
    
    echo -e "${GREEN}═══════════════════════════════════════${NC}"
    echo -e "${GREEN}✓ Setup complete! Ready to launch.${NC}"
    echo -e "${GREEN}═══════════════════════════════════════${NC}\n"
    
    echo -e "${BLUE}Next step:${NC} ${YELLOW}go run main.go${NC} or ${YELLOW}make run${NC}\n"
}

main