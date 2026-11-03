# Blockchain Dashboard

A professional, visually appealing real-time blockchain visualization and management dashboard built with Next.js, TypeScript, and shadcn/ui. Features a modern black/white/gray theme with comprehensive blockchain interaction capabilities.

## Features

### 💼 Wallet Management
- Create new wallets with cryptographic key generation
- View all wallets with balance and nonce information
- Copy wallet addresses and public keys
- Select wallets for transactions
- Real-time balance updates

### 📊 Blockchain Visualization
- **Interactive Block Chain**: Click on blocks to view detailed information
- **Block Details**: Hash, previous hash, merkle root, difficulty, and nonce
- **Transaction Display**: View all transactions in each block
- **Mempool Status**: See pending transactions waiting to be mined
- **Network Stats**: Total wallets, block height, and pending transactions

### 💰 Transaction Management
- Submit transactions from one wallet to another
- Specify transaction amounts and fees
- Real-time balance validation
- Transaction status confirmation
- Automatic total cost calculation (amount + fee)

### ⛏️ Mining Controls
- Start/stop the mining process
- Real-time mining status indicator
- Mining status polling (every 2 seconds)
- Block reward tracking

### ⚙️ Configuration
- Changeable blockchain API endpoint
- Environment variable support
- Persistent settings in localStorage
- Easy URL configuration via settings panel

### 🎨 Professional UI
- Clean black/white/gray color scheme
- Responsive design (mobile to desktop)
- Dark mode support
- Smooth animations and transitions
- Accessibility-first components (shadcn/ui)
- Professional typography and spacing

## Setup & Installation

### Prerequisites
- Node.js 18+ (pnpm recommended)
- A running blockchain backend (see Backend Configuration)

### Installation

1. **Clone and install dependencies**
```bash
pnpm install
```

2. **Configure the blockchain API endpoint**
Create a `.env.local` file:
```env
NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://txhm5mlj-3000.inc1.devtunnels.ms
```

Or update it later via the settings panel (⚙️ button in the bottom-right)

3. **Start the development server**
```bash
pnpm dev
```

The dashboard will be available at `http://localhost:3000`

## Backend Configuration

The dashboard connects to a blockchain backend API. Update the API URL using one of these methods:

### Method 1: Environment Variable (Recommended)
Edit `.env.local`:
```env
NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://your-blockchain-api.com:3000
```

### Method 2: Settings Panel
1. Click the ⚙️ (Settings) button in the bottom-right corner
2. Enter your blockchain API URL
3. Click "Save"
4. Reload the page

### Method 3: Local Storage
The last configured URL is saved in browser localStorage and will persist across sessions.

## Architecture

### Components

- **`app/page.tsx`**: Main dashboard layout and state management
- **`components/wallet-manager.tsx`**: Wallet creation and management
- **`components/transaction-form.tsx`**: Transaction submission interface
- **`components/blockchain-visualizer.tsx`**: Block chain visualization and exploration
- **`components/mining-controls.tsx`**: Mining start/stop controls
- **`components/settings-panel.tsx`**: API URL configuration

### API Client (`lib/api.ts`)

Type-safe API client with endpoints for:
- Wallet operations (create, get, list)
- Transaction operations (submit, view mempool)
- Blockchain operations (get blocks, view chain)
- Mining operations (start, stop, check status)

### Styling

- **Tailwind CSS v4**: Utility-first CSS framework
- **shadcn/ui**: Accessible component library
- **Design Tokens**: Black/white/gray color system with semantic tokens
- **Responsive**: Mobile-first design with Tailwind breakpoints

## Usage Guide

### Creating a Wallet

1. Click the **"New Wallet"** button in the Wallets section
2. A new wallet will be created with a public/private key pair
3. Your wallet address and balance will be displayed

### Sending a Transaction

1. Select a sender wallet from the dropdown in the Transaction form
2. Enter the recipient's wallet address
3. Specify the transaction amount
4. (Optional) Adjust the transaction fee (minimum 1)
5. The dashboard shows total cost (amount + fee)
6. Click **"Submit Transaction"**
7. Successful transactions appear in the mempool and then in blocks once mined

### Mining Blocks

1. Click **"Start Mining"** in the Mining section
2. The miner will begin creating blocks from pending transactions
3. Mining status is displayed with a green indicator
4. Click **"Stop Mining"** to halt the process
5. Mined blocks appear in the blockchain visualization

### Viewing Block Details

1. Scroll to the **Blockchain Chain** section
2. Click on any block number to view its details
3. View all transaction details for that block
4. Copy hashes using the copy button

## API Endpoints

The dashboard expects the following endpoints from your blockchain backend:

### Wallets
- `GET /api/wallet/all` - List all wallets
- `POST /api/wallet/create` - Create a new wallet
- `GET /api/wallet/{address}` - Get wallet details

### Transactions
- `POST /api/transaction/submit` - Submit a transaction
- `GET /api/mempool` - Get pending transactions

### Blockchain
- `GET /api/blockchain` - Get all blocks
- `GET /api/blockchain/{index}` - Get a specific block
- `GET /api/blockchain/height` - Get current block height

### Mining
- `POST /api/miner/start` - Start mining
- `POST /api/miner/stop` - Stop mining
- `GET /api/miner/status` - Get mining status

## Design System

### Color Palette
- **White**: `oklch(1 0 0)` - Backgrounds, cards
- **Black**: `oklch(0.145 0 0)` - Text, foreground
- **Gray Shades**: Various oklch values for borders, muted text, and accents

### Typography
- **Sans Font**: Default for all text
- **Mono Font**: For addresses, hashes, technical information

### Spacing & Borders
- **Border Radius**: 0.625rem (10px)
- **Gaps**: Standard Tailwind spacing scale (4px, 8px, 16px, etc.)

## Deployment

### Deploy to Vercel

```bash
vercel deploy
```

Ensure you set the `NEXT_PUBLIC_BLOCKCHAIN_API_URL` environment variable in Vercel project settings.

## Development

### Project Structure
```
/app          - Next.js app directory
/components   - React components
/lib          - Utilities and API client
/public       - Static assets
```

### Tech Stack
- **Framework**: Next.js 16 with App Router
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4
- **UI Components**: shadcn/ui
- **Icons**: Lucide React
- **HTTP**: Fetch API (no external HTTP library)

### Code Quality
- TypeScript for type safety
- ESLint for code linting
- Responsive design tested on multiple viewports
- Accessibility-first component usage

## Troubleshooting

### "Failed to load wallets"
- Verify the blockchain backend is running
- Check that the API URL is correct in `.env.local`
- Use the settings panel (⚙️) to verify the API URL
- Check browser console for detailed error messages

### Transactions not appearing
- Ensure the mining process is running
- Check that enough time has passed for the transaction to be mined
- Verify the transaction fee is sufficient
- Check that the sender wallet has enough balance

### Mining not starting
- Verify the blockchain backend supports the mining API
- Check if mining is already running (status indicator)
- Wait 2 seconds for status to update
- Check browser console for errors

## Performance

- **Auto-refresh**: Updates every 3 seconds
- **Mining status**: Checked every 2 seconds
- **Responsive**: Optimized for 60fps animations
- **Lazy loading**: Components load only when needed

## Security Notes

⚠️ **This is a frontend dashboard for blockchain interaction.**
- Never expose private keys in browser storage
- Always use HTTPS for blockchain API connections
- Validate all user inputs on both client and server
- Backend must implement proper authentication and authorization

## Future Enhancements

- [ ] Real-time WebSocket updates (instead of polling)
- [ ] Transaction history with filtering
- [ ] Advanced block explorer features
- [ ] Network statistics graphs
- [ ] Wallet import/export functionality
- [ ] Multi-signature transaction support
- [ ] Transaction search and filtering

## License

MIT License - Feel free to use and modify this dashboard for your blockchain projects.

## Support

For issues or questions:
1. Check the troubleshooting section above
2. Verify the blockchain backend is running and accessible
3. Check browser console (F12) for detailed error messages
4. Ensure all required API endpoints are implemented

---

**Built with ❤️ using Next.js and shadcn/ui**
