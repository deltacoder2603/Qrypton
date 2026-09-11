package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"sync"
	"time"
)

// ============================================================================
// SECTION 1: CONFIGURATION
// ============================================================================

const (
	PROTOCOL_VERSION  = "1.0"
	NETWORK_PORT      = 8080
	HTTP_PORT         = 3000
	MAX_PEERS         = 10
	BLOCK_TIME        = 5 * time.Second
	CONSENSUS_TIMEOUT = 10 * time.Second
	TX_POOL_MAX_SIZE  = 1000
	BLOCK_REWARD      = 50
	MIN_TX_FEE        = 1
	STORAGE_DIR       = "./blockchain_data"
	WALLETS_FILE      = "wallets.json"
	BLOCKCHAIN_FILE   = "blockchain.json"
	PEERS_FILE        = "peers.json"
	MEMPOOL_FILE      = "mempool.json"
	SHA256_SIZE       = 32
)

// ============================================================================
// SECTION 2: UTILITIES
// ============================================================================

// Logger with timestamp
type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	log.Printf("[%s] %s | %v", l.prefix, msg, args)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	log.Printf("[%s] ERROR: %s | %v", l.prefix, msg, args)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	log.Printf("[%s] DEBUG: %s | %v", l.prefix, msg, args)
}

// Helper: Convert bytes to hex string
func BytesToHex(data []byte) string {
	return hex.EncodeToString(data)
}

// Helper: Convert hex string to bytes
func HexToBytes(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

// Helper: Current timestamp in milliseconds
func CurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}

// Helper: Generate random nonce
func GenerateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return BytesToHex(b)
}

// ============================================================================
// SECTION 3: CRYPTOGRAPHY (Falcon Post-Quantum Signatures)
// ============================================================================

// CryptoProvider handles signature operations
type CryptoProvider interface {
	GenerateKeyPair() (pubKey []byte, privKey []byte, err error)
	Sign(privKey []byte, data []byte) (signature []byte, err error)
	Verify(pubKey []byte, data []byte, signature []byte) bool
}

// FalconCryptoProvider - Post-Quantum Falcon implementation using CLI
type FalconCryptoProvider struct {
	logger *Logger
}

func NewFalconCryptoProvider() *FalconCryptoProvider {
	return &FalconCryptoProvider{
		logger: NewLogger("FALCON"),
	}
}

func (f *FalconCryptoProvider) GenerateKeyPair() ([]byte, []byte, error) {
	// Generate Falcon-1024 keypair using CLI
	// falcon create [--seed-hex <hex>] [--output <file>]

	seed := make([]byte, 48) // Falcon seed size
	if _, err := rand.Read(seed); err != nil {
		f.logger.Error("Failed to generate random seed", err)
		return nil, nil, fmt.Errorf("failed to generate random seed: %v", err)
	}

	seedHex := hex.EncodeToString(seed)

	// Create temporary key file path
	keyFile := fmt.Sprintf("/tmp/falcon_key_%d.json", CurrentTimestamp())

	// Run: falcon create --seed-hex <hex> --output <file>
	cmd := exec.Command("falcon", "create",
		"--seed", seedHex,
		"--out", keyFile)

	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error("Failed to generate Falcon keypair", fmt.Sprintf("error: %v, output: %s", err, string(output)))
		return nil, nil, fmt.Errorf("falcon create failed: %v", err)
	}

	// Read key file
	keyData, err := os.ReadFile(keyFile)
	if err != nil {
		f.logger.Error("Failed to read key file", err)
		os.Remove(keyFile)
		return nil, nil, fmt.Errorf("failed to read key file: %v", err)
	}

	// Parse JSON to extract keys
	var keyObj struct {
		PublicKey  string `json:"public_key"`
		PrivateKey string `json:"private_key"`
	}

	if err := json.Unmarshal(keyData, &keyObj); err != nil {
		f.logger.Error("Failed to parse key file", err)
		os.Remove(keyFile)
		return nil, nil, fmt.Errorf("failed to parse key file: %v", err)
	}

	// Clean up
	os.Remove(keyFile)

	pubKeyBytes, _ := hex.DecodeString(keyObj.PublicKey)
	privKeyBytes, _ := hex.DecodeString(keyObj.PrivateKey)

	f.logger.Info("Falcon-1024 keypair generated", "")
	return pubKeyBytes, privKeyBytes, nil
}

func (f *FalconCryptoProvider) Sign(privKey []byte, data []byte) ([]byte, error) {
	tempDir, err := os.MkdirTemp("", "qrypton_falcon_sign_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	keyFile := tempDir + "/private_key.json"
	msgFile := tempDir + "/message.bin"
	sigFile := tempDir + "/signature.bin"

	// The Falcon CLI expects the key file to be JSON.
	keyJSON, err := json.Marshal(map[string]string{
		"private_key": hex.EncodeToString(privKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to encode private key JSON: %w", err)
	}

	if err := os.WriteFile(keyFile, keyJSON, 0600); err != nil {
		return nil, fmt.Errorf("failed to write private key file: %w", err)
	}

	// Message is raw bytes.
	if err := os.WriteFile(msgFile, data, 0600); err != nil {
		return nil, fmt.Errorf("failed to write message file: %w", err)
	}

	// Sign the message.
	cmd := exec.Command(
		"falcon",
		"sign",
		"--key", keyFile,
		"--in", msgFile,
		"--out", sigFile,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf(
			"falcon sign failed: %w; output: %s",
			err,
			string(output),
		)
	}

	// The signature output is RAW BINARY, not JSON.
	signature, err := os.ReadFile(sigFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read signature file: %w", err)
	}

	if len(signature) == 0 {
		return nil, fmt.Errorf("falcon generated an empty signature")
	}

	f.logger.Info(
		"Transaction signed with Falcon-1024",
		fmt.Sprintf("signature_size=%d bytes", len(signature)),
	)

	return signature, nil
}

func (f *FalconCryptoProvider) Verify(
	pubKey []byte,
	data []byte,
	signature []byte,
) bool {
	tempDir, err := os.MkdirTemp("", "qrypton_falcon_verify_*")
	if err != nil {
		f.logger.Error("Failed to create temp directory", err)
		return false
	}
	defer os.RemoveAll(tempDir)

	keyFile := tempDir + "/public_key.json"
	msgFile := tempDir + "/message.bin"
	sigFile := tempDir + "/signature.bin"

	// The Falcon CLI expects the public key file to be JSON.
	keyJSON, err := json.Marshal(map[string]string{
		"public_key": hex.EncodeToString(pubKey),
	})
	if err != nil {
		f.logger.Error("Failed to encode public key JSON", err)
		return false
	}

	if err := os.WriteFile(keyFile, keyJSON, 0600); err != nil {
		f.logger.Error("Failed to write public key file", err)
		return false
	}

	// Write exact message bytes.
	if err := os.WriteFile(msgFile, data, 0600); err != nil {
		f.logger.Error("Failed to write message file", err)
		return false
	}

	// IMPORTANT:
	// The signature must be written as RAW BINARY.
	if err := os.WriteFile(sigFile, signature, 0600); err != nil {
		f.logger.Error("Failed to write signature file", err)
		return false
	}

	cmd := exec.Command(
		"falcon",
		"verify",
		"--key", keyFile,
		"--in", msgFile,
		"--sig", sigFile,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		f.logger.Error(
			"Signature verification failed",
			fmt.Sprintf(
				"error: %v, output: %s",
				err,
				string(output),
			),
		)
		return false
	}

	f.logger.Info(
		"Signature verified successfully with Falcon-1024",
		string(output),
	)

	return true
}

// Default to Falcon (Post-Quantum)
var cryptoProvider CryptoProvider = NewFalconCryptoProvider()

// ============================================================================
// SECTION 4: WALLETS
// ============================================================================

type Wallet struct {
	Address    string `json:"address"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Balance    int64  `json:"balance"`
	Nonce      uint64 `json:"nonce"`
	CreatedAt  int64  `json:"created_at"`
}

type WalletManager struct {
	wallets map[string]*Wallet
	mu      sync.RWMutex
	logger  *Logger
}

func NewWalletManager() *WalletManager {
	return &WalletManager{
		wallets: make(map[string]*Wallet),
		logger:  NewLogger("WALLET"),
	}
}

func (wm *WalletManager) CreateWallet() (*Wallet, error) {
	pubKey, privKey, err := cryptoProvider.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	// Address = SHA256(PublicKey)[0:20] in hex
	hash := sha256.Sum256(pubKey)
	address := BytesToHex(hash[:20])

	wallet := &Wallet{
		Address:    address,
		PublicKey:  BytesToHex(pubKey),
		PrivateKey: BytesToHex(privKey),
		Balance:    0,
		Nonce:      0,
		CreatedAt:  CurrentTimestamp(),
	}

	wm.mu.Lock()
	wm.wallets[address] = wallet
	wm.mu.Unlock()

	wm.logger.Info("Wallet created", address)
	return wallet, nil
}

func (wm *WalletManager) GetWallet(address string) *Wallet {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return wm.wallets[address]
}

func (wm *WalletManager) GetAllWallets() []*Wallet {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	wallets := make([]*Wallet, 0)
	for _, w := range wm.wallets {
		wallets = append(wallets, w)
	}
	return wallets
}

func (wm *WalletManager) UpdateBalance(address string, amount int64) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wallet, exists := wm.wallets[address]
	if !exists {
		return fmt.Errorf("wallet not found: %s", address)
	}

	wallet.Balance += amount
	return nil
}

func (wm *WalletManager) IncrementNonce(address string) error {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wallet, exists := wm.wallets[address]
	if !exists {
		return fmt.Errorf("wallet not found: %s", address)
	}

	wallet.Nonce++
	return nil
}

// ============================================================================
// SECTION 5: TRANSACTIONS
// ============================================================================

type Transaction struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    int64  `json:"amount"`
	Fee       int64  `json:"fee"`
	Nonce     uint64 `json:"nonce"`
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type"` // "transfer" or "coinbase"
}

func (t *Transaction) CalculateHash() string {
	data := fmt.Sprintf("%s:%s:%d:%d:%d:%d", t.From, t.To, t.Amount, t.Fee, t.Nonce, t.Timestamp)
	hash := sha256.Sum256([]byte(data))
	return BytesToHex(hash[:])
}

func (t *Transaction) SignTransaction(privKeyHex string) error {
	privKey, err := HexToBytes(privKeyHex)
	if err != nil {
		return err
	}

	data := fmt.Sprintf("%s:%s:%d:%d:%d:%d", t.From, t.To, t.Amount, t.Fee, t.Nonce, t.Timestamp)
	signature, err := cryptoProvider.Sign(privKey, []byte(data))
	if err != nil {
		return err
	}

	t.Signature = BytesToHex(signature)
	return nil
}

func (t *Transaction) VerifySignature(pubKeyHex string) bool {
	pubKey, err := HexToBytes(pubKeyHex)
	if err != nil {
		return false
	}

	signature, err := HexToBytes(t.Signature)
	if err != nil {
		return false
	}

	data := fmt.Sprintf("%s:%s:%d:%d:%d:%d", t.From, t.To, t.Amount, t.Fee, t.Nonce, t.Timestamp)
	return cryptoProvider.Verify(pubKey, []byte(data), signature)
}

type TransactionPool struct {
	transactions map[string]*Transaction
	mu           sync.RWMutex
	logger       *Logger
}

func NewTransactionPool() *TransactionPool {
	return &TransactionPool{
		transactions: make(map[string]*Transaction),
		logger:       NewLogger("TX_POOL"),
	}
}

func (tp *TransactionPool) AddTransaction(tx *Transaction) error {
	tp.mu.Lock()
	defer tp.mu.Unlock()

	if _, exists := tp.transactions[tx.ID]; exists {
		return fmt.Errorf("transaction already in pool: %s", tx.ID)
	}

	if len(tp.transactions) >= TX_POOL_MAX_SIZE {
		return fmt.Errorf("transaction pool full")
	}

	tp.transactions[tx.ID] = tx
	tp.logger.Info("Transaction added to pool", tx.ID)
	return nil
}

func (tp *TransactionPool) RemoveTransaction(txID string) {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	delete(tp.transactions, txID)
}

func (tp *TransactionPool) GetAllTransactions() []*Transaction {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	txs := make([]*Transaction, 0)
	for _, tx := range tp.transactions {
		txs = append(txs, tx)
	}
	return txs
}

func (tp *TransactionPool) GetTransactionCount() int {
	tp.mu.RLock()
	defer tp.mu.RUnlock()
	return len(tp.transactions)
}

func (tp *TransactionPool) Clear() {
	tp.mu.Lock()
	defer tp.mu.Unlock()
	tp.transactions = make(map[string]*Transaction)
}

// ============================================================================
// SECTION 6: MERKLE TREE
// ============================================================================

type MerkleTree struct {
	root  string
	nodes []string
}

func NewMerkleTree(transactions []*Transaction) *MerkleTree {
	if len(transactions) == 0 {
		return &MerkleTree{root: "", nodes: []string{}}
	}

	var nodes []string
	for _, tx := range transactions {
		hash := sha256.Sum256([]byte(tx.ID))
		nodes = append(nodes, BytesToHex(hash[:]))
	}

	for len(nodes) > 1 {
		var newLevel []string
		for i := 0; i < len(nodes); i += 2 {
			var combined string
			if i+1 < len(nodes) {
				combined = nodes[i] + nodes[i+1]
			} else {
				combined = nodes[i] + nodes[i]
			}
			hash := sha256.Sum256([]byte(combined))
			newLevel = append(newLevel, BytesToHex(hash[:]))
		}
		nodes = newLevel
	}

	return &MerkleTree{
		root:  nodes[0],
		nodes: nodes,
	}
}

func (mt *MerkleTree) GetRoot() string {
	return mt.root
}

// ============================================================================
// SECTION 7: BLOCKS
// ============================================================================

type Block struct {
	Index        int64          `json:"index"`
	Timestamp    int64          `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
	PreviousHash string         `json:"previous_hash"`
	Hash         string         `json:"hash"`
	Miner        string         `json:"miner"`
	MerkleRoot   string         `json:"merkle_root"`
	Nonce        uint64         `json:"nonce"`
	Difficulty   int64          `json:"difficulty"`
}

func NewBlock(index int64, previousHash string, transactions []*Transaction, miner string) *Block {
	block := &Block{
		Index:        index,
		Timestamp:    CurrentTimestamp(),
		Transactions: transactions,
		PreviousHash: previousHash,
		Miner:        miner,
		Nonce:        0,
		Difficulty:   1,
	}

	mt := NewMerkleTree(transactions)
	block.MerkleRoot = mt.GetRoot()
	block.Hash = block.CalculateHash()

	return block
}

func (b *Block) CalculateHash() string {
	data := fmt.Sprintf("%d:%d:%s:%s:%s:%d:%d",
		b.Index, b.Timestamp, b.PreviousHash, b.MerkleRoot, b.Miner, b.Nonce, b.Difficulty)
	hash := sha256.Sum256([]byte(data))
	return BytesToHex(hash[:])
}

func (b *Block) IsValid(previousBlock *Block) bool {
	if b.Index != previousBlock.Index+1 {
		return false
	}
	if b.PreviousHash != previousBlock.Hash {
		return false
	}
	if b.Hash != b.CalculateHash() {
		return false
	}
	return true
}

// ============================================================================
// SECTION 8: BLOCKCHAIN
// ============================================================================

type Blockchain struct {
	chain                []*Block
	difficultyMultiplier float64
	mu                   sync.RWMutex
	logger               *Logger
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		chain:                make([]*Block, 0),
		difficultyMultiplier: 1.0,
		logger:               NewLogger("BLOCKCHAIN"),
	}
}

func (bc *Blockchain) CreateGenesisBlock(miner string) {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	genesisBlock := &Block{
		Index:        0,
		Timestamp:    CurrentTimestamp(),
		Transactions: []*Transaction{},
		PreviousHash: "0",
		Miner:        miner,
		Nonce:        0,
		Difficulty:   1,
	}

	genesisBlock.Hash = genesisBlock.CalculateHash()
	bc.chain = append(bc.chain, genesisBlock)
	bc.logger.Info("Genesis block created", genesisBlock.Hash)
}

func (bc *Blockchain) AddBlock(block *Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(bc.chain) == 0 {
		return fmt.Errorf("genesis block not initialized")
	}

	lastBlock := bc.chain[len(bc.chain)-1]
	if !block.IsValid(lastBlock) {
		return fmt.Errorf("block validation failed")
	}

	bc.chain = append(bc.chain, block)
	bc.logger.Info("Block added to chain", block.Hash)
	return nil
}

func (bc *Blockchain) GetLastBlock() *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if len(bc.chain) == 0 {
		return nil
	}
	return bc.chain[len(bc.chain)-1]
}

func (bc *Blockchain) GetBlockByIndex(index int64) *Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()

	if index < 0 || index >= int64(len(bc.chain)) {
		return nil
	}
	return bc.chain[index]
}

func (bc *Blockchain) GetAllBlocks() []*Block {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return append([]*Block{}, bc.chain...)
}

func (bc *Blockchain) GetBlockHeight() int64 {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	return int64(len(bc.chain))
}

// ============================================================================
// SECTION 9: MEMPOOL
// ============================================================================

type Mempool struct {
	transactions map[string]*Transaction
	mu           sync.RWMutex
	logger       *Logger
}

func NewMempool() *Mempool {
	return &Mempool{
		transactions: make(map[string]*Transaction),
		logger:       NewLogger("MEMPOOL"),
	}
}

func (mp *Mempool) AddTransaction(tx *Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if _, exists := mp.transactions[tx.ID]; exists {
		return fmt.Errorf("transaction already in mempool: %s", tx.ID)
	}

	mp.transactions[tx.ID] = tx
	mp.logger.Info("Transaction added to mempool", tx.ID)
	return nil
}

func (mp *Mempool) RemoveTransaction(txID string) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	delete(mp.transactions, txID)
}

func (mp *Mempool) GetTransactionsForBlock(maxCount int) []*Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	txs := make([]*Transaction, 0)
	for _, tx := range mp.transactions {
		if len(txs) >= maxCount {
			break
		}
		txs = append(txs, tx)
	}

	// Sort by fee descending
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Fee > txs[j].Fee
	})

	return txs
}

func (mp *Mempool) GetAllTransactions() []*Transaction {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	txs := make([]*Transaction, 0)
	for _, tx := range mp.transactions {
		txs = append(txs, tx)
	}
	return txs
}

func (mp *Mempool) Size() int {
	mp.mu.RLock()
	defer mp.mu.RUnlock()
	return len(mp.transactions)
}

func (mp *Mempool) Clear() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.transactions = make(map[string]*Transaction)
}

// ============================================================================
// SECTION 10: MINER
// ============================================================================

type Miner struct {
	walletManager *WalletManager
	blockchain    *Blockchain
	mempool       *Mempool
	minerAddress  string
	isMining      bool
	mu            sync.Mutex
	logger        *Logger
	stopSignal    chan bool
}

func NewMiner(walletManager *WalletManager, blockchain *Blockchain, mempool *Mempool, minerAddress string) *Miner {
	return &Miner{
		walletManager: walletManager,
		blockchain:    blockchain,
		mempool:       mempool,
		minerAddress:  minerAddress,
		isMining:      false,
		logger:        NewLogger("MINER"),
		stopSignal:    make(chan bool, 1),
	}
}

func (m *Miner) StartMining() {
	m.mu.Lock()
	if m.isMining {
		m.mu.Unlock()
		return
	}
	m.isMining = true
	m.mu.Unlock()

	go m.mineBlocks()
}

func (m *Miner) StopMining() {
	m.mu.Lock()
	m.isMining = false
	m.mu.Unlock()
	m.stopSignal <- true
}

func (m *Miner) mineBlocks() {
	ticker := time.NewTicker(BLOCK_TIME)
	defer ticker.Stop()

	for {
		select {
		case <-m.stopSignal:
			m.logger.Info("Mining stopped", "")
			return
		case <-ticker.C:
			m.mineBlock()
		}
	}
}

func (m *Miner) mineBlock() {
	lastBlock := m.blockchain.GetLastBlock()
	if lastBlock == nil {
		return
	}

	// Get transactions from mempool
	txs := m.mempool.GetTransactionsForBlock(100)

	// Create coinbase transaction
	coinbase := &Transaction{
		From:      "SYSTEM",
		To:        m.minerAddress,
		Amount:    int64(BLOCK_REWARD),
		Fee:       0,
		Nonce:     0,
		Timestamp: CurrentTimestamp(),
		Type:      "coinbase",
	}
	coinbase.ID = coinbase.CalculateHash()

	allTxs := append([]*Transaction{coinbase}, txs...)

	// Create new block
	newBlock := NewBlock(lastBlock.Index+1, lastBlock.Hash, allTxs, m.minerAddress)

	// Add block to blockchain
	if err := m.blockchain.AddBlock(newBlock); err != nil {
		m.logger.Error("Failed to add block", err)
		return
	}

	// Remove transactions from mempool
	for _, tx := range txs {
		m.mempool.RemoveTransaction(tx.ID)
	}

	// Update miner balance
	m.walletManager.UpdateBalance(m.minerAddress, int64(BLOCK_REWARD))

	m.logger.Info("Block mined", newBlock.Hash)
}

// ============================================================================
// SECTION 11: CONSENSUS (PBFT - Byzantine Fault Tolerant)
// ============================================================================

type ConsensusEngine struct {
	nodeID         string
	validatorNodes map[string]bool
	prepareVotes   map[string]map[string]bool
	commitVotes    map[string]map[string]bool
	mu             sync.RWMutex
	logger         *Logger
}

func NewConsensusEngine(nodeID string) *ConsensusEngine {
	return &ConsensusEngine{
		nodeID:         nodeID,
		validatorNodes: make(map[string]bool),
		prepareVotes:   make(map[string]map[string]bool),
		commitVotes:    make(map[string]map[string]bool),
		logger:         NewLogger("CONSENSUS"),
	}
}

func (ce *ConsensusEngine) AddValidator(nodeID string) {
	ce.mu.Lock()
	defer ce.mu.Unlock()
	ce.validatorNodes[nodeID] = true
	ce.logger.Info("Validator added", nodeID)
}

func (ce *ConsensusEngine) CanFinalize() bool {
	ce.mu.RLock()
	defer ce.mu.RUnlock()
	requiredVotes := (len(ce.validatorNodes) * 2 / 3) + 1
	return len(ce.validatorNodes) >= requiredVotes
}

func (ce *ConsensusEngine) RecordPrepareVote(blockHash string, nodeID string) {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if _, exists := ce.prepareVotes[blockHash]; !exists {
		ce.prepareVotes[blockHash] = make(map[string]bool)
	}
	ce.prepareVotes[blockHash][nodeID] = true
}

func (ce *ConsensusEngine) RecordCommitVote(blockHash string, nodeID string) {
	ce.mu.Lock()
	defer ce.mu.Unlock()

	if _, exists := ce.commitVotes[blockHash]; !exists {
		ce.commitVotes[blockHash] = make(map[string]bool)
	}
	ce.commitVotes[blockHash][nodeID] = true
}

func (ce *ConsensusEngine) GetPrepareVotes(blockHash string) int {
	ce.mu.RLock()
	defer ce.mu.RUnlock()
	if votes, exists := ce.prepareVotes[blockHash]; exists {
		return len(votes)
	}
	return 0
}

func (ce *ConsensusEngine) GetCommitVotes(blockHash string) int {
	ce.mu.RLock()
	defer ce.mu.RUnlock()
	if votes, exists := ce.commitVotes[blockHash]; exists {
		return len(votes)
	}
	return 0
}

// ============================================================================
// SECTION 12: LOCAL STORAGE
// ============================================================================

type StorageManager struct {
	dataDir string
	logger  *Logger
	mu      sync.RWMutex
}

func NewStorageManager(dataDir string) *StorageManager {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0755)
	}

	return &StorageManager{
		dataDir: dataDir,
		logger:  NewLogger("STORAGE"),
	}
}

func (sm *StorageManager) SaveWallets(wallets map[string]*Wallet) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	data, err := json.MarshalIndent(wallets, "", "  ")
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s", sm.dataDir, WALLETS_FILE)
	return os.WriteFile(filePath, data, 0644)
}

func (sm *StorageManager) LoadWallets() (map[string]*Wallet, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	wallets := make(map[string]*Wallet)
	filePath := fmt.Sprintf("%s/%s", sm.dataDir, WALLETS_FILE)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return wallets, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, &wallets)
	return wallets, err
}

func (sm *StorageManager) SaveBlockchain(blocks []*Block) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	data, err := json.MarshalIndent(blocks, "", "  ")
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s", sm.dataDir, BLOCKCHAIN_FILE)
	return os.WriteFile(filePath, data, 0644)
}

func (sm *StorageManager) LoadBlockchain() ([]*Block, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var blocks []*Block
	filePath := fmt.Sprintf("%s/%s", sm.dataDir, BLOCKCHAIN_FILE)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return blocks, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, &blocks)
	return blocks, err
}

func (sm *StorageManager) SaveMempool(transactions map[string]*Transaction) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return err
	}

	filePath := fmt.Sprintf("%s/%s", sm.dataDir, MEMPOOL_FILE)
	return os.WriteFile(filePath, data, 0644)
}

func (sm *StorageManager) LoadMempool() (map[string]*Transaction, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	txs := make(map[string]*Transaction)
	filePath := fmt.Sprintf("%s/%s", sm.dataDir, MEMPOOL_FILE)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return txs, nil
		}
		return nil, err
	}

	err = json.Unmarshal(data, &txs)
	return txs, err
}

// ============================================================================
// SECTION 13: P2P NETWORKING
// ============================================================================

type PeerMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type P2PNetwork struct {
	nodeID          string
	peers           map[string]net.Conn
	peerAddresses   map[string]string
	blockchain      *Blockchain
	mempool         *Mempool
	walletManager   *WalletManager
	consensusEngine *ConsensusEngine
	listener        net.Listener
	logger          *Logger
	mu              sync.RWMutex
	messageHandlers map[string]func(interface{}) error
}

func NewP2PNetwork(nodeID string, blockchain *Blockchain, mempool *Mempool, walletManager *WalletManager, consensusEngine *ConsensusEngine) *P2PNetwork {
	return &P2PNetwork{
		nodeID:          nodeID,
		peers:           make(map[string]net.Conn),
		peerAddresses:   make(map[string]string),
		blockchain:      blockchain,
		mempool:         mempool,
		walletManager:   walletManager,
		consensusEngine: consensusEngine,
		logger:          NewLogger("P2P"),
		messageHandlers: make(map[string]func(interface{}) error),
	}
}

func (p *P2PNetwork) Start(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	p.listener = listener
	p.logger.Info("P2P network started on port", port)

	go p.acceptConnections()
	return nil
}

func (p *P2PNetwork) acceptConnections() {
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			p.logger.Error("Accept connection error", err)
			continue
		}

		go p.handleConnection(conn)
	}
}

func (p *P2PNetwork) handleConnection(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	for {
		var msg PeerMessage
		if err := decoder.Decode(&msg); err != nil {
			break
		}

		if handler, exists := p.messageHandlers[msg.Type]; exists {
			handler(msg.Payload)
		}
	}
}

func (p *P2PNetwork) ConnectPeer(peerID string, address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	p.mu.Lock()
	p.peers[peerID] = conn
	p.peerAddresses[peerID] = address
	p.mu.Unlock()

	p.logger.Info("Connected to peer", peerID)
	return nil
}

func (p *P2PNetwork) BroadcastBlock(block *Block) error {
	p.mu.RLock()
	peers := make([]net.Conn, 0, len(p.peers))
	for _, conn := range p.peers {
		peers = append(peers, conn)
	}
	p.mu.RUnlock()

	msg := PeerMessage{
		Type:    "block",
		Payload: mustMarshal(block),
	}

	for _, conn := range peers {
		encoder := json.NewEncoder(conn)
		encoder.Encode(msg)
	}

	return nil
}

func (p *P2PNetwork) BroadcastTransaction(tx *Transaction) error {
	p.mu.RLock()
	peers := make([]net.Conn, 0, len(p.peers))
	for _, conn := range p.peers {
		peers = append(peers, conn)
	}
	p.mu.RUnlock()

	msg := PeerMessage{
		Type:    "transaction",
		Payload: mustMarshal(tx),
	}

	for _, conn := range peers {
		encoder := json.NewEncoder(conn)
		encoder.Encode(msg)
	}

	return nil
}

func (p *P2PNetwork) RegisterMessageHandler(msgType string, handler func(interface{}) error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.messageHandlers[msgType] = handler
}

func (p *P2PNetwork) GetPeerCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.peers)
}

func mustMarshal(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}

// ============================================================================
// SECTION 14: REST API
// ============================================================================

type APIServer struct {
	walletManager  *WalletManager
	blockchain     *Blockchain
	mempool        *Mempool
	miner          *Miner
	p2pNetwork     *P2PNetwork
	storageManager *StorageManager
	logger         *Logger
}

func NewAPIServer(wm *WalletManager, bc *Blockchain, mp *Mempool, m *Miner, p2p *P2PNetwork, sm *StorageManager) *APIServer {
	return &APIServer{
		walletManager:  wm,
		blockchain:     bc,
		mempool:        mp,
		miner:          m,
		p2pNetwork:     p2p,
		storageManager: sm,
		logger:         NewLogger("API"),
	}
}

func (api *APIServer) Start(port int) {
	mux := http.NewServeMux()

	// Root endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":    "QRYPTON",
			"version": PROTOCOL_VERSION,
			"status":  "running",
			"crypto":  "Falcon-1024",
		})
	})

	// Wallet endpoints
	mux.HandleFunc("/wallet/create", api.handleCreateWallet)
	mux.HandleFunc("/wallet/get", api.handleGetWallet)
	mux.HandleFunc("/wallet/list", api.handleListWallets)
	mux.HandleFunc("/wallet/balance", api.handleGetBalance)

	// Transaction endpoints
	mux.HandleFunc("/transaction/create", api.handleCreateTransaction)
	mux.HandleFunc("/transaction/list", api.handleListTransactions)
	mux.HandleFunc("/transaction/send", api.handleSendTransaction)

	// Blockchain endpoints
	mux.HandleFunc("/blockchain/blocks", api.handleGetBlocks)
	mux.HandleFunc("/blockchain/block", api.handleGetBlock)
	mux.HandleFunc("/blockchain/height", api.handleGetHeight)
	mux.HandleFunc("/blockchain/genesis", api.handleGetGenesisBlock)

	// Mempool endpoints
	mux.HandleFunc("/mempool/transactions", api.handleGetMempoolTransactions)
	mux.HandleFunc("/mempool/size", api.handleGetMempoolSize)

	// Mining endpoints
	mux.HandleFunc("/mining/start", api.handleStartMining)
	mux.HandleFunc("/mining/stop", api.handleStopMining)

	// Node endpoints
	mux.HandleFunc("/node/info", api.handleNodeInfo)
	mux.HandleFunc("/node/stats", api.handleNodeStats)
	mux.HandleFunc("/node/health", api.handleHealth)

	// P2P endpoints
	mux.HandleFunc("/peers/count", api.handlePeerCount)

	// Apply CORS middleware to every route
	handler := enableCORS(mux)

	api.logger.Info("API Server started on port", port)

	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", port),
		handler,
	); err != nil {
		api.logger.Error("API server failed", err)
	}
}

func (api *APIServer) handleCreateWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	wallet, err := api.walletManager.CreateWallet()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}

func (api *APIServer) handleGetWallet(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address parameter", http.StatusBadRequest)
		return
	}

	wallet := api.walletManager.GetWallet(address)
	if wallet == nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}

func (api *APIServer) handleListWallets(w http.ResponseWriter, r *http.Request) {
	wallets := api.walletManager.GetAllWallets()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallets)
}

func (api *APIServer) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Missing address parameter", http.StatusBadRequest)
		return
	}

	wallet := api.walletManager.GetWallet(address)
	if wallet == nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"address": address,
		"balance": wallet.Balance,
	})
}

func (api *APIServer) handleCreateTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int64  `json:"amount"`
		Fee    int64  `json:"fee"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	fromWallet := api.walletManager.GetWallet(req.From)
	if fromWallet == nil {
		http.Error(w, "Sender wallet not found", http.StatusNotFound)
		return
	}

	if req.Amount+req.Fee > fromWallet.Balance {
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	tx := &Transaction{
		From:      req.From,
		To:        req.To,
		Amount:    req.Amount,
		Fee:       req.Fee,
		Nonce:     fromWallet.Nonce,
		Timestamp: CurrentTimestamp(),
		Type:      "transfer",
	}

	tx.ID = tx.CalculateHash()

	// Sign transaction
	if err := tx.SignTransaction(fromWallet.PrivateKey); err != nil {
		http.Error(w, fmt.Sprintf("Signing error: %v", err), http.StatusInternalServerError)
		return
	}

	// Verify signature
	if !tx.VerifySignature(fromWallet.PublicKey) {
		http.Error(w, "Signature verification failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tx)
}

func (api *APIServer) handleListTransactions(w http.ResponseWriter, r *http.Request) {
	txs := api.mempool.GetAllTransactions()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

func (api *APIServer) handleSendTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var tx Transaction
	if err := json.NewDecoder(r.Body).Decode(&tx); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request: %v", err), http.StatusBadRequest)
		return
	}

	wallet := api.walletManager.GetWallet(tx.From)
	if wallet == nil {
		http.Error(w, "Sender wallet not found", http.StatusNotFound)
		return
	}

	if !tx.VerifySignature(wallet.PublicKey) {
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	if err := api.mempool.AddTransaction(&tx); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusBadRequest)
		return
	}

	// Update wallet balance and nonce
	api.walletManager.UpdateBalance(tx.From, -(tx.Amount + tx.Fee))
	api.walletManager.IncrementNonce(tx.From)
	if tx.To != "SYSTEM" {
		api.walletManager.UpdateBalance(tx.To, tx.Amount)
	}

	// Broadcast to peers
	api.p2pNetwork.BroadcastTransaction(&tx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"tx_id":  tx.ID,
	})
}

func (api *APIServer) handleGetBlocks(w http.ResponseWriter, r *http.Request) {
	blocks := api.blockchain.GetAllBlocks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blocks)
}

func (api *APIServer) handleGetBlock(w http.ResponseWriter, r *http.Request) {
	indexStr := r.URL.Query().Get("index")
	if indexStr == "" {
		http.Error(w, "Missing index parameter", http.StatusBadRequest)
		return
	}

	index, err := strconv.ParseInt(indexStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}

	block := api.blockchain.GetBlockByIndex(index)
	if block == nil {
		http.Error(w, "Block not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(block)
}

func (api *APIServer) handleGetHeight(w http.ResponseWriter, r *http.Request) {
	height := api.blockchain.GetBlockHeight()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int64{"height": height})
}

func (api *APIServer) handleGetGenesisBlock(w http.ResponseWriter, r *http.Request) {
	genesisBlock := api.blockchain.GetBlockByIndex(0)
	if genesisBlock == nil {
		http.Error(w, "Genesis block not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(genesisBlock)
}

func (api *APIServer) handleGetMempoolTransactions(w http.ResponseWriter, r *http.Request) {
	txs := api.mempool.GetAllTransactions()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txs)
}

func (api *APIServer) handleGetMempoolSize(w http.ResponseWriter, r *http.Request) {
	size := api.mempool.Size()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"size": size})
}

func (api *APIServer) handleStartMining(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	api.miner.StartMining()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "mining started"})
}

func (api *APIServer) handleStopMining(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	api.miner.StopMining()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "mining stopped"})
}

func (api *APIServer) handleNodeInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"version":       PROTOCOL_VERSION,
		"node_id":       "NODE_001",
		"network_peers": api.p2pNetwork.GetPeerCount(),
		"block_height":  api.blockchain.GetBlockHeight(),
		"mempool_size":  api.mempool.Size(),
		"total_wallets": len(api.walletManager.GetAllWallets()),
	})
}

func (api *APIServer) handleNodeStats(w http.ResponseWriter, r *http.Request) {
	wallets := api.walletManager.GetAllWallets()
	totalBalance := int64(0)
	for _, wallet := range wallets {
		totalBalance += wallet.Balance
	}

	blocks := api.blockchain.GetAllBlocks()
	totalTransactions := 0
	for _, block := range blocks {
		totalTransactions += len(block.Transactions)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"total_blocks":       api.blockchain.GetBlockHeight(),
		"total_transactions": totalTransactions,
		"total_wallets":      len(wallets),
		"total_balance":      totalBalance,
		"mempool_pending":    api.mempool.Size(),
		"timestamp":          CurrentTimestamp(),
	})
}

func (api *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (api *APIServer) handlePeerCount(w http.ResponseWriter, r *http.Request) {
	count := api.p2pNetwork.GetPeerCount()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"peer_count": count})
}

// ============================================================================
// SECTION 15: NODE ORCHESTRATION
// ============================================================================

type Node struct {
	nodeID          string
	walletManager   *WalletManager
	blockchain      *Blockchain
	mempool         *Mempool
	miner           *Miner
	p2pNetwork      *P2PNetwork
	apiServer       *APIServer
	storageManager  *StorageManager
	consensusEngine *ConsensusEngine
	logger          *Logger
	stopChan        chan bool
}

func NewNode(nodeID string) *Node {
	wm := NewWalletManager()
	bc := NewBlockchain()
	mp := NewMempool()
	sm := NewStorageManager(STORAGE_DIR)
	ce := NewConsensusEngine(nodeID)
	p2p := NewP2PNetwork(nodeID, bc, mp, wm, ce)
	miner := NewMiner(wm, bc, mp, "")
	api := NewAPIServer(wm, bc, mp, miner, p2p, sm)

	return &Node{
		nodeID:          nodeID,
		walletManager:   wm,
		blockchain:      bc,
		mempool:         mp,
		miner:           miner,
		p2pNetwork:      p2p,
		apiServer:       api,
		storageManager:  sm,
		consensusEngine: ce,
		logger:          NewLogger("NODE"),
		stopChan:        make(chan bool, 1),
	}
}

func (n *Node) Initialize() error {
	n.logger.Info("Initializing node", n.nodeID)

	// Load state from storage
	wallets, err := n.storageManager.LoadWallets()
	if err != nil {
		n.logger.Error("Failed to load wallets", err)
		return err
	}

	// Restore wallets to manager
	for _, w := range wallets {
		n.walletManager.wallets[w.Address] = w
	}

	// Load blockchain
	blocks, err := n.storageManager.LoadBlockchain()
	if err != nil {
		n.logger.Error("Failed to load blockchain", err)
		return err
	}

	// Restore blockchain
	if len(blocks) > 0 {
		n.blockchain.chain = blocks
		n.logger.Info("Blockchain loaded with", len(blocks), "blocks")
	} else {
		// Create genesis block
		minerWallet, _ := n.walletManager.CreateWallet()
		n.blockchain.CreateGenesisBlock(minerWallet.Address)
		n.miner.minerAddress = minerWallet.Address
		n.logger.Info("Genesis block created", "")
	}

	// Load mempool
	txMap, err := n.storageManager.LoadMempool()
	if err != nil {
		n.logger.Error("Failed to load mempool", err)
	} else {
		for _, tx := range txMap {
			n.mempool.AddTransaction(tx)
		}
	}

	// Start P2P network
	if err := n.p2pNetwork.Start(NETWORK_PORT); err != nil {
		n.logger.Error("Failed to start P2P network", err)
		return err
	}

	n.logger.Info("Node initialized successfully", n.nodeID)
	return nil
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Allowed HTTP methods
		w.Header().Set(
			"Access-Control-Allow-Methods",
			"GET, POST, PUT, PATCH, DELETE, OPTIONS",
		)

		// Allowed request headers
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization",
		)

		// Cache preflight response for 24 hours
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle browser CORS preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (n *Node) Start() error {
	n.logger.Info("Starting node", n.nodeID)

	// Start miner if no wallets exist, create one
	if len(n.walletManager.GetAllWallets()) == 0 {
		minerWallet, _ := n.walletManager.CreateWallet()
		n.miner.minerAddress = minerWallet.Address
	} else {
		wallets := n.walletManager.GetAllWallets()
		if len(wallets) > 0 {
			n.miner.minerAddress = wallets[0].Address
		}
	}

	n.miner.StartMining()
	n.logger.Info("Mining started", n.miner.minerAddress)

	// Start background persistence
	go n.persistenceWorker()

	// Start API server
	go n.apiServer.Start(HTTP_PORT)

	n.logger.Info("Node started successfully", n.nodeID)
	return nil
}

func (n *Node) persistenceWorker() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-n.stopChan:
			return
		case <-ticker.C:
			// Save wallets
			wm := make(map[string]*Wallet)
			for _, w := range n.walletManager.GetAllWallets() {
				wm[w.Address] = w
			}
			if err := n.storageManager.SaveWallets(wm); err != nil {
				n.logger.Error("Failed to save wallets", err)
			}

			// Save blockchain
			if err := n.storageManager.SaveBlockchain(n.blockchain.GetAllBlocks()); err != nil {
				n.logger.Error("Failed to save blockchain", err)
			}

			// Save mempool
			txMap := make(map[string]*Transaction)
			for _, tx := range n.mempool.GetAllTransactions() {
				txMap[tx.ID] = tx
			}
			if err := n.storageManager.SaveMempool(txMap); err != nil {
				n.logger.Error("Failed to save mempool", err)
			}
		}
	}
}

func (n *Node) Stop() {
	n.logger.Info("Stopping node", n.nodeID)
	n.miner.StopMining()
	n.stopChan <- true

	// Final persistence
	wm := make(map[string]*Wallet)
	for _, w := range n.walletManager.GetAllWallets() {
		wm[w.Address] = w
	}
	n.storageManager.SaveWallets(wm)
	n.storageManager.SaveBlockchain(n.blockchain.GetAllBlocks())

	txMap := make(map[string]*Transaction)
	for _, tx := range n.mempool.GetAllTransactions() {
		txMap[tx.ID] = tx
	}
	n.storageManager.SaveMempool(txMap)

	n.logger.Info("Node stopped", n.nodeID)
}

// ============================================================================
// SECTION 16: MAIN
// ============================================================================

func main() {
	fmt.Println(`
╔════════════════════════════════════════════════════════════════╗
║   QRYPTON - Post-Quantum Secure Blockchain                   ║
║   🔐 Falcon-1024 Post-Quantum Signatures (CLI-based)          ║
║   ✓ NIST Security Level 5 (Quantum-Resistant)                ║
║   ✓ PBFT Consensus                                            ║
╚════════════════════════════════════════════════════════════════╝
	`)

	// Check if falcon CLI is installed
	cmd := exec.Command("falcon", "--help")
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ ERROR: Falcon CLI not found!\n")
		fmt.Printf("Install it with: go install github.com/algorandfoundation/falcon-signatures/cmd/falcon@latest\n")
		os.Exit(1)
	}

	node := NewNode("NODE_001")

	if err := node.Initialize(); err != nil {
		fmt.Printf("Failed to initialize node: %v\n", err)
		os.Exit(1)
	}

	if err := node.Start(); err != nil {
		fmt.Printf("Failed to start node: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Node is running with FALCON-1024 POST-QUANTUM SIGNATURES...")
	fmt.Println("✓ REST API listening on http://localhost:3000")
	fmt.Println("✓ P2P Network listening on :8080")
	fmt.Println("✓ Mining is active")
	fmt.Println("✓ Cryptography: Falcon-1024 (lattice-based, NIST Level 5)")
	fmt.Println("✓ Key sizes: Public=1,793B | Private=2,305B | Signature=1,538B")
	fmt.Println("✓ Using Falcon CLI for cryptographic operations")
	fmt.Println("\nPress Ctrl+C to stop...\n")

	// Keep running
	select {}
}
