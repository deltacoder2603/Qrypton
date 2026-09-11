#!/bin/bash

BASE_URL="http://localhost:3000"
DELAY=1

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

line() {
echo -e "${BLUE}════════════════════════════════════════════════════════════${NC}"
}

request() {

METHOD=$1
URL=$2
BODY=$3
TITLE=$4

echo -e "${BLUE}$TITLE${NC}"
echo "$METHOD $URL"

if [ -z "$BODY" ]; then
    RESPONSE=$(curl -s -X "$METHOD" "$BASE_URL$URL")
else
    RESPONSE=$(curl -s \
        -X "$METHOD" \
        "$BASE_URL$URL" \
        -H "Content-Type: application/json" \
        -d "$BODY")
fi

if echo "$RESPONSE" | jq . >/dev/null 2>&1; then
    echo "$RESPONSE" | jq .
else
    echo "$RESPONSE"
fi

echo
sleep $DELAY

}

echo
echo "╔══════════════════════════════════════════════════════╗"
echo "║            QRYPTON END TO END TEST                  ║"
echo "╚══════════════════════════════════════════════════════╝"
echo

STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/node/health")

if [ "$STATUS" != "200" ]; then
    echo -e "${RED}Server not running on port 3000${NC}"
    exit 1
fi

line
echo "SECTION 1 - NODE"
line

request GET /node/health "" "Health"
request GET /node/info "" "Node Info"
request GET /node/stats "" "Node Stats"

line
echo "SECTION 2 - CREATE WALLETS"
line

RAW1=$(curl -s -X POST "$BASE_URL/wallet/create")
echo "$RAW1" | jq .

ADDR1=$(echo "$RAW1" | jq -r '.address')

RAW2=$(curl -s -X POST "$BASE_URL/wallet/create")
echo "$RAW2" | jq .

ADDR2=$(echo "$RAW2" | jq -r '.address')

echo
echo "Wallet 1 : $ADDR1"
echo "Wallet 2 : $ADDR2"
echo

sleep 2

request GET /wallet/list "" "Wallet List"
request GET "/wallet/get?address=$ADDR1" "" "Wallet 1"
request GET "/wallet/get?address=$ADDR2" "" "Wallet 2"

request GET "/wallet/balance?address=$ADDR1" "" "Wallet1 Balance"
request GET "/wallet/balance?address=$ADDR2" "" "Wallet2 Balance"

line
echo "SECTION 3 - BLOCKCHAIN"
line

request GET /blockchain/height "" "Height"
request GET /blockchain/genesis "" "Genesis"
request GET /blockchain/blocks "" "All Blocks"

line
echo "SECTION 4 - MINING"
line

request POST /mining/start "" "Start Mining"

echo
echo "Waiting 20 seconds..."
echo
sleep 20

request GET /blockchain/height "" "Height After Mining"
request GET /blockchain/blocks "" "Blocks"

line
echo "SECTION 5 - MEMPOOL"
line

request GET /mempool/transactions "" "Transactions"
request GET /mempool/size "" "Size"

line
echo "SECTION 6 - PEERS"
line

request GET /peers/count "" "Peer Count"

line
echo "SECTION 7 - NODE"
line

request GET /node/stats "" "Statistics"

line
echo "SECTION 8 - STOP MINING"
line

request POST /mining/stop "" "Stop Mining"

line
echo "SECTION 9 - FALCON CLI"
line

echo "Creating Falcon key..."

falcon create --out test_key.json

echo
echo "Signing..."

falcon sign \
    --key test_key.json \
    --msg "Hello QRYPTON" \
    --out signature.bin

echo
echo "Verifying..."

falcon verify \
    --key test_key.json \
    --msg "Hello QRYPTON" \
    --sig signature.bin

echo
echo "Algorand Address"

falcon algorand address \
    --key test_key.json

echo
line
echo -e "${GREEN}ALL TESTS FINISHED${NC}"
line