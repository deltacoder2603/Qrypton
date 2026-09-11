# Blockchain Dashboard Configuration Guide

## Quick Start

1. **Update the API URL**
   - Edit `.env.local` and set your blockchain backend URL:
   ```env
   NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://txhm5mlj-3000.inc1.devtunnels.ms
   ```

2. **Start the dashboard**
   ```bash
   pnpm dev
   ```

3. **Access the dashboard**
   - Open `http://localhost:3000` in your browser

## Changing the Blockchain Backend URL

### Option 1: Via Environment File (Recommended for Development)

1. Open `.env.local` in the project root
2. Update the `NEXT_PUBLIC_BLOCKCHAIN_API_URL`:
   ```env
   # Development
   NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://txhm5mlj-3000.inc1.devtunnels.ms
   
   # Or for local development
   NEXT_PUBLIC_BLOCKCHAIN_API_URL=http://localhost:3000
   ```
3. Restart the dev server: `pnpm dev`

### Option 2: Via Settings Panel (Runtime)

1. Click the **⚙️ Settings** button (bottom-right corner)
2. Enter your blockchain API URL in the input field
3. Click **Save**
4. **Reload the page** for changes to take effect
5. The URL is saved to browser localStorage

**Note**: Changes via the settings panel are stored locally and will persist across browser sessions. However, refreshing the page will reload the environment variable value unless you've saved it via the settings panel.

### Option 3: For Production Deployment

If deploying to Vercel:

1. Go to your Vercel project settings
2. Navigate to **Environment Variables**
3. Add: `NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://your-production-blockchain-url.com`
4. Redeploy your application

## Environment Variable Reference

### `NEXT_PUBLIC_BLOCKCHAIN_API_URL` (Required)

The blockchain backend API endpoint that the dashboard communicates with.

- **Type**: URL string
- **Example**: `https://txhm5mlj-3000.inc1.devtunnels.ms`
- **Default**: `http://localhost:3000` (if not set)
- **Public**: Yes (prefixed with `NEXT_PUBLIC_`)

#### Valid URL Formats
```env
# HTTPS (recommended for production)
NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://example.com
NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://example.com:3000
NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://example.com/api

# HTTP (development only)
NEXT_PUBLIC_BLOCKCHAIN_API_URL=http://localhost:3000
NEXT_PUBLIC_BLOCKCHAIN_API_URL=http://192.168.1.100:3000
```

## Required Backend API Endpoints

Ensure your blockchain backend implements these endpoints:

### Wallet Endpoints
- **`GET /api/wallet/all`**
  - Response: `Wallet[]`
  - Get all wallets in the network

- **`POST /api/wallet/create`**
  - Response: `Wallet`
  - Create a new wallet with generated keypair

- **`GET /api/wallet/{address}`**
  - Parameters: `address` (string)
  - Response: `Wallet`
  - Get details for a specific wallet

### Transaction Endpoints
- **`POST /api/transaction/submit`**
  - Body: `{ from: string, to: string, amount: number, fee: number }`
  - Response: `Transaction`
  - Submit a new transaction

- **`GET /api/mempool`**
  - Response: `Transaction[]`
  - Get all pending transactions

### Blockchain Endpoints
- **`GET /api/blockchain`**
  - Response: `Block[]`
  - Get all blocks in the chain

- **`GET /api/blockchain/{index}`**
  - Parameters: `index` (number)
  - Response: `Block`
  - Get a specific block by index

- **`GET /api/blockchain/height`**
  - Response: `{ height: number }`
  - Get current blockchain height

### Mining Endpoints
- **`POST /api/miner/start`**
  - Start the mining process

- **`POST /api/miner/stop`**
  - Stop the mining process

- **`GET /api/miner/status`**
  - Response: `{ is_mining: boolean }`
  - Get current mining status

## Data Types

### Wallet
```typescript
interface Wallet {
  address: string      // Wallet address (derived from public key)
  public_key: string   // Public key (hex encoded)
  balance: number      // Current balance in units
  nonce: number        // Transaction counter
}
```

### Transaction
```typescript
interface Transaction {
  id: string           // Transaction ID
  from: string         // Sender wallet address
  to: string           // Recipient wallet address
  amount: number       // Amount to transfer
  fee: number          // Transaction fee
  nonce: number        // Sender's nonce
  signature: string    // Cryptographic signature
  timestamp: number    // Unix timestamp in milliseconds
  type: 'transfer' | 'coinbase'  // Transaction type
}
```

### Block
```typescript
interface Block {
  index: number        // Block number in chain
  timestamp: number    // Creation time (Unix ms)
  transactions: Transaction[]  // Transactions in block
  previous_hash: string        // Hash of previous block
  hash: string                 // Block's own hash
  miner: string                // Miner's wallet address
  merkle_root: string          // Merkle tree root of transactions
  nonce: number                // Proof of work nonce
  difficulty: number           // Mining difficulty
}
```

## Troubleshooting Configuration Issues

### Dashboard shows "Failed to load wallets"

**Possible causes:**
1. Blockchain backend is not running
2. API URL is incorrect
3. Backend API is not accessible from your location
4. CORS is not properly configured on the backend

**Solutions:**
1. Verify the backend is running and accessible:
   ```bash
   curl https://txhm5mlj-3000.inc1.devtunnels.ms/api/wallet/all
   ```

2. Check the configured URL in the dashboard settings panel (⚙️)

3. Verify `.env.local` has the correct URL and restart the dev server

4. Check browser console (F12 → Console tab) for detailed error messages

5. Ensure backend allows CORS requests from your dashboard URL

### API URL changes not taking effect

**Solutions:**
1. Restart the dev server after updating `.env.local`
2. Clear browser localStorage: Press F12 → Application → Local Storage → Clear All
3. Hard refresh the page: Ctrl+Shift+R (Windows) or Cmd+Shift+R (Mac)
4. Verify the new URL using the settings panel

### Transactions not processing

**Check:**
1. Mining is started (look for mining indicator in header)
2. Sender wallet has sufficient balance (amount + fee)
3. Transaction fees are positive numbers
4. Backend has implemented the `/api/transaction/submit` endpoint

### CORS Errors

**Fix on backend:**
```go
// Add CORS headers to all responses
func cors(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
}
```

## Advanced Configuration

### Using Environment Variables in Custom Scripts

To use the API URL in your own scripts:

```typescript
// Access in client-side code
const apiUrl = process.env.NEXT_PUBLIC_BLOCKCHAIN_API_URL
```

### Runtime Environment Override

The dashboard automatically uses localStorage-saved settings over environment variables:

```typescript
// In lib/api.ts
const API_URL = process.env.NEXT_PUBLIC_BLOCKCHAIN_API_URL || 'http://localhost:3000'
```

To use localStorage value, the settings panel will override on next page load.

## Multi-Environment Setup

### Development
```env
# .env.local
NEXT_PUBLIC_BLOCKCHAIN_API_URL=http://localhost:3000
```

### Staging
```env
# .env.staging
NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://staging-blockchain.example.com
```

### Production
```env
# Set in Vercel dashboard
NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://blockchain.example.com
```

Build for specific environment:
```bash
# Using Vercel environments
vercel env pull .env.production.local
pnpm run build
```

## Security Considerations

⚠️ **Important Security Notes:**

1. **HTTPS Required for Production**
   - Always use HTTPS URLs in production
   - Never send sensitive data over HTTP

2. **API URL is Public**
   - `NEXT_PUBLIC_*` variables are visible in browser
   - Do not put secrets or credentials in this variable

3. **CORS Configuration**
   - Ensure your backend properly validates origins
   - Don't use wildcard `*` for sensitive APIs in production

4. **Authentication**
   - Implement proper authentication on backend
   - Validate all requests server-side

## Next Steps

1. Verify your blockchain backend is running
2. Set the correct API URL in `.env.local`
3. Restart the dev server with `pnpm dev`
4. Create your first wallet
5. Start mining and submit transactions!

For detailed API information, see [README.md](./README.md)
