# Quick Start Guide - Blockchain Dashboard

Get your blockchain dashboard up and running in 5 minutes! ⚡

## Step 1: Install Dependencies (30 seconds)

```bash
pnpm install
```

## Step 2: Configure Blockchain API URL (1 minute)

Edit `.env.local` in the project root:

```env
NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://txhm5mlj-3000.inc1.devtunnels.ms
```

**Options**:
- **Local dev**: `http://localhost:3000`
- **Remote server**: `https://your-blockchain-api.com`
- **Update later**: Use settings panel (⚙️) in dashboard

## Step 3: Start Development Server (1 minute)

```bash
pnpm dev
```

The terminal will show:
```
- Local:         http://localhost:3000
- Network:       http://100.64.89.251:3000
- Ready in 632ms
```

## Step 4: Open Dashboard (1 minute)

Open your browser to: **http://localhost:3000**

You should see:
- ✅ Header with "Blockchain Dashboard" title
- ✅ Wallet management panel on the left
- ✅ Mining controls
- ✅ Transaction form (if backend is running)
- ✅ Blockchain chain visualization
- ✅ Settings button in bottom-right (⚙️)

## Step 5: Create & Test (2 minutes)

### Create a Wallet
1. Click **"New Wallet"** button
2. Wait for success
3. Wallet appears in the list

### Start Mining
1. Click **"Start Mining"** button
2. Mining status changes to **"Mining Active"**
3. Blocks start forming
4. Click **"Stop Mining"** to stop

### Submit Transaction
1. **Sender**: Select from dropdown
2. **Recipient**: Enter wallet address
3. **Amount**: Enter number
4. **Fee**: (usually 1)
5. Click **"Submit Transaction"**
6. Transaction appears in Mempool
7. Wait for mining → transaction in block

## 🎯 Common Tasks

### Change Blockchain API URL

**While Developing** (fastest):
1. Edit `.env.local`
2. Update `NEXT_PUBLIC_BLOCKCHAIN_API_URL`
3. Restart dev server: `pnpm dev`

**At Runtime** (no restart):
1. Click ⚙️ (Settings button, bottom-right)
2. Enter new URL
3. Click Save
4. Reload page

**For Production** (Vercel):
1. Go to Vercel project dashboard
2. Settings → Environment Variables
3. Add `NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://your-url.com`
4. Redeploy

### View Block Details
1. Scroll to "Blockchain Chain" section
2. Click any block number
3. See all block details including transactions
4. Copy hashes with copy button

### Check Wallet Balance
1. Look at "Wallets" section on left
2. Each wallet shows balance and nonce
3. Total balance shown at bottom

### Monitor Network
1. Check "Network Stats" box (left sidebar)
2. See:
   - Total Wallets
   - Block Height
   - Pending Transactions

## ⚠️ Troubleshooting

### "Failed to load wallets"
**Solution**: Blockchain backend not running or URL is wrong
```bash
# Check if backend is accessible:
curl https://your-blockchain-url/api/wallet/all
```

### Transactions don't appear
**Solution**: Mining not started
- Click "Start Mining" in Mining section
- Wait a few seconds

### API URL changes don't work
**Solution**: Clear cache and restart
```bash
# Clear browser cache:
1. Press F12 to open DevTools
2. Go to Application tab
3. Click "Clear site data"
4. Refresh page
```

Or restart dev server:
```bash
# Ctrl+C to stop
# Then restart:
pnpm dev
```

### Port 3000 already in use
**Solution**: Use different port
```bash
PORT=3001 pnpm dev
```

## 📚 Documentation

| Document | Purpose |
|----------|---------|
| [README.md](./README.md) | Main features & overview |
| [CONFIGURATION.md](./CONFIGURATION.md) | API setup & endpoints |
| [COMPONENTS.md](./COMPONENTS.md) | Code structure |
| [PROJECT_SUMMARY.md](./PROJECT_SUMMARY.md) | Complete overview |
| [FILES_CREATED.md](./FILES_CREATED.md) | What was built |
| [QUICKSTART.md](./QUICKSTART.md) | This file |

## 🚀 Next Steps

### After Setup Works

1. **Explore**: Create wallets, send transactions, watch mining
2. **Read Code**: Review `app/page.tsx` and `lib/api.ts`
3. **Customize**: Update colors in `app/globals.css`
4. **Deploy**: Follow deployment guide in [README.md](./README.md)

### For Production

```bash
# Build for production
pnpm build

# Start production server
pnpm start

# Or deploy to Vercel
vercel deploy
```

## 💡 Pro Tips

### 1. Keyboard Shortcuts
- `F12` - Open DevTools (see console errors)
- `Ctrl+Shift+R` - Hard refresh page
- `Ctrl+Shift+Delete` - Clear browser cache

### 2. Developer Tools
- **DevTools Console**: Click F12 → Console tab
- **Network Tab**: See API calls in Network tab
- **React DevTools**: Install extension for component debugging

### 3. Wallet Tips
- Always note your wallet **address** (used as recipient)
- **Nonce** = number of transactions sent
- **Balance** = how many units you have
- Copy address to send money FROM

### 4. Transaction Tips
- Fee is **additional cost** on top of amount
- Sender must have: **amount + fee**
- Transactions appear in Mempool first
- Appear in Blockchain after mining

### 5. Mining Tips
- Mining creates **new blocks**
- Miner gets **block rewards**
- Transactions move from Mempool → Block → Blockchain
- Can view all transactions in each block

## 🔧 Quick Settings

### Dark Mode
The dashboard respects system dark mode preference. To toggle:
1. System settings (Windows/Mac)
2. Or browser dark mode extension

### API Endpoint
```env
# Development (local)
NEXT_PUBLIC_BLOCKCHAIN_API_URL=http://localhost:3000

# Staging
NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://staging-blockchain.example.com

# Production
NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://blockchain.example.com
```

### Port Number
```bash
# Use different port if 3000 is taken
PORT=3001 pnpm dev
```

## 📊 What the Dashboard Does

### Real-Time Updates
- **Auto-refresh**: Every 3 seconds
- **Mining status**: Every 2 seconds
- **Manual refresh**: Click ↻ button

### Displays
- 📋 All wallets with balances
- 🔗 Complete blockchain with all blocks
- ⏳ Pending transactions (mempool)
- 🎯 Detailed block information
- 📊 Network statistics

### Features
- ➕ Create unlimited wallets
- 💰 Send transactions
- ⛏️ Start/stop mining
- 🔍 Explore blocks
- 💾 Copy addresses/hashes
- ⚙️ Configure API URL

## 🎓 Learning Resources

### Blockchain Concepts
- **Wallet**: Account with address, balance, public/private keys
- **Transaction**: Transfer of units from one wallet to another
- **Block**: Container of transactions, linked by hashes
- **Mining**: Process of creating new blocks
- **Mempool**: Pool of pending transactions

### Dashboard Components
- **WalletManager**: Create and view wallets
- **TransactionForm**: Send transactions
- **BlockchainVisualizer**: View blockchain
- **MiningControls**: Control mining
- **SettingsPanel**: Configure API

## ✅ Checklist

- [ ] Run `pnpm install`
- [ ] Create `.env.local` with API URL
- [ ] Run `pnpm dev`
- [ ] Open `http://localhost:3000`
- [ ] See dashboard loaded
- [ ] (Optional) Create wallet to test backend connection
- [ ] Read [README.md](./README.md) for more details
- [ ] Deploy when ready!

## 🆘 Still Having Issues?

1. **Check DevTools**: Press F12, go to Console tab
2. **Verify Backend**: Test API with curl
3. **Check URL**: Verify `.env.local` has correct URL
4. **Clear Cache**: Ctrl+Shift+Delete → Clear all
5. **Restart Server**: Ctrl+C then `pnpm dev`
6. **Read Docs**: See [CONFIGURATION.md](./CONFIGURATION.md)

## 🎉 You're Ready!

Your blockchain dashboard is ready to use. Start by:

1. Creating wallets
2. Sending transactions
3. Mining blocks
4. Exploring the blockchain

Enjoy! 🚀

---

**Need help?** Check the full documentation:
- [README.md](./README.md) - Complete guide
- [CONFIGURATION.md](./CONFIGURATION.md) - Configuration details
- [COMPONENTS.md](./COMPONENTS.md) - Code architecture

**Happy blockchain exploring!** ⛓️
