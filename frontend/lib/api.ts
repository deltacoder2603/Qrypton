const API_URL = process.env.NEXT_PUBLIC_BLOCKCHAIN_API_URL || 'http://localhost:3000/'

export interface Wallet {
  address: string
  public_key: string
  private_key?: string
  balance: number
  nonce: number
  created_at?: number
}

export interface Transaction {
  id: string
  from: string
  to: string
  amount: number
  fee: number
  nonce: number
  signature: string
  timestamp: number
  type: 'transfer' | 'coinbase'
}

export interface Block {
  index: number
  timestamp: number
  transactions: Transaction[]
  previous_hash: string
  hash: string
  miner: string
  merkle_root: string
  nonce: number
  difficulty: number
}

export interface BlockchainStatus {
  blocks: Block[]
  mempool: Transaction[]
  wallets: Wallet[]
  block_height: number
}

// Wallet API - Create new wallet
export async function createWallet(): Promise<Wallet> {
  const response = await fetch(`${API_URL}/wallet/create`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  })
  if (!response.ok) throw new Error('Failed to create wallet')
  return response.json()
}

// Wallet API - Get wallet by address
export async function getWallet(address: string): Promise<Wallet> {
  const response = await fetch(`${API_URL}/wallet/get?address=${encodeURIComponent(address)}`)
  if (!response.ok) throw new Error('Failed to get wallet')
  return response.json()
}

// Wallet API - Get all wallets
export async function getAllWallets(): Promise<Wallet[]> {
  const response = await fetch(`${API_URL}/wallet/list`)
  if (!response.ok) throw new Error('Failed to get wallets')
  return response.json()
}

// Wallet API - Get wallet balance
export async function getWalletBalance(address: string): Promise<{ address: string; balance: number }> {
  const response = await fetch(`${API_URL}/wallet/balance?address=${encodeURIComponent(address)}`)
  if (!response.ok) throw new Error('Failed to get wallet balance')
  return response.json()
}

// Transaction API - Create and sign transaction
export async function createTransaction(
  from: string,
  to: string,
  amount: number,
  fee: number = 1
): Promise<Transaction> {
  const response = await fetch(`${API_URL}/transaction/create`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ from, to, amount, fee }),
  })
  if (!response.ok) {
    const error = await response.json().catch(() => ({}))
    throw new Error(error.error || 'Failed to create transaction')
  }
  return response.json()
}

// Transaction API - Send signed transaction
export async function sendTransaction(transaction: Transaction): Promise<{ status: string; tx_id: string }> {
  const response = await fetch(`${API_URL}/transaction/send`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(transaction),
  })
  if (!response.ok) {
    const error = await response.json().catch(() => ({}))
    throw new Error(error.error || 'Failed to send transaction')
  }
  return response.json()
}

// Transaction API - Submit transaction (create + send in one call)
export async function submitTransaction(
  from: string,
  to: string,
  amount: number,
  fee: number = 1
): Promise<{ status: string; tx_id: string }> {
  const tx = await createTransaction(from, to, amount, fee)
  return sendTransaction(tx)
}

// Transaction API - Get mempool transactions
export async function getMempool(): Promise<Transaction[]> {
  const response = await fetch(`${API_URL}/mempool/transactions`)
  if (!response.ok) throw new Error('Failed to get mempool')
  return response.json()
}

// Transaction API - Get transaction list (alias for mempool)
export async function getTransactionList(): Promise<Transaction[]> {
  const response = await fetch(`${API_URL}/transaction/list`)
  if (!response.ok) throw new Error('Failed to get transaction list')
  return response.json()
}

// Blockchain API - Get all blocks
export async function getBlockchain(): Promise<Block[]> {
  const response = await fetch(`${API_URL}/blockchain/blocks`)
  if (!response.ok) throw new Error('Failed to get blockchain')
  return response.json()
}

// Blockchain API - Get single block by index
export async function getBlock(index: number): Promise<Block> {
  const response = await fetch(`${API_URL}/blockchain/block?index=${index}`)
  if (!response.ok) throw new Error('Failed to get block')
  return response.json()
}

// Blockchain API - Get blockchain height
export async function getBlockHeight(): Promise<number> {
  const response = await fetch(`${API_URL}/blockchain/height`)
  if (!response.ok) throw new Error('Failed to get block height')
  const data = await response.json()
  return data.height
}

// Blockchain API - Get genesis block
export async function getGenesisBlock(): Promise<Block> {
  const response = await fetch(`${API_URL}/blockchain/genesis`)
  if (!response.ok) throw new Error('Failed to get genesis block')
  return response.json()
}

// Mempool API - Get mempool size
export async function getMempoolSize(): Promise<number> {
  const response = await fetch(`${API_URL}/mempool/size`)
  if (!response.ok) throw new Error('Failed to get mempool size')
  const data = await response.json()
  return data.size
}

// Mining API - Start mining
export async function startMining(): Promise<{ status: string }> {
  const response = await fetch(`${API_URL}/mining/start`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  })
  if (!response.ok) throw new Error('Failed to start mining')
  return response.json()
}

// Mining API - Stop mining
export async function stopMining(): Promise<{ status: string }> {
  const response = await fetch(`${API_URL}/mining/stop`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
  })
  if (!response.ok) throw new Error('Failed to stop mining')
  return response.json()
}

// Node API - Get node info
export async function getNodeInfo(): Promise<{
  version: string
  node_id: string
  network_peers: number
  block_height: number
  mempool_size: number
  total_wallets: number
}> {
  const response = await fetch(`${API_URL}/node/info`)
  if (!response.ok) throw new Error('Failed to get node info')
  return response.json()
}

// Node API - Get node stats
export async function getNodeStats(): Promise<{
  total_blocks: number
  total_transactions: number
  total_wallets: number
  total_balance: number
  mempool_pending: number
  timestamp: number
}> {
  const response = await fetch(`${API_URL}/node/stats`)
  if (!response.ok) throw new Error('Failed to get node stats')
  return response.json()
}

// Node API - Health check
export async function getHealth(): Promise<{ status: string }> {
  const response = await fetch(`${API_URL}/node/health`)
  if (!response.ok) throw new Error('Failed to get health status')
  return response.json()
}

// P2P API - Get peer count
export async function getPeerCount(): Promise<number> {
  const response = await fetch(`${API_URL}/peers/count`)
  if (!response.ok) throw new Error('Failed to get peer count')
  const data = await response.json()
  return data.peer_count
}

// Combined status function
export async function getStatus(): Promise<BlockchainStatus> {
  try {
    const [blocks, mempool, wallets, blockHeight] = await Promise.all([
      getBlockchain(),
      getMempool(),
      getAllWallets(),
      getBlockHeight(),
    ])
    return { blocks, mempool, wallets, block_height: blockHeight }
  } catch (error) {
    const errorMsg = error instanceof Error ? error.message : 'Unknown error'
    console.error('[v0] Failed to get status:', errorMsg)
    console.error('[v0] API URL:', API_URL)
    console.error('[v0] Make sure blockchain backend is running at:', API_URL)
    return { blocks: [], mempool: [], wallets: [], block_height: 0 }
  }
}
