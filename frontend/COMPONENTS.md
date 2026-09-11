# Component Documentation

## Overview

This document describes all React components in the blockchain dashboard and how they work together.

## Directory Structure

```
components/
├── blockchain-visualizer.tsx    # Block chain visualization & explorer
├── wallet-manager.tsx           # Wallet creation & management
├── transaction-form.tsx         # Transaction submission form
├── mining-controls.tsx          # Mining start/stop controls
├── settings-panel.tsx           # API URL configuration
└── ui/
    └── button.tsx               # shadcn Button component
```

## Component Details

### 1. BlockchainVisualizer

**Location**: `components/blockchain-visualizer.tsx`

**Purpose**: Displays the entire blockchain and allows interaction with blocks and transactions.

**Props**:
```typescript
interface BlockVisualizerProps {
  blocks: Block[]
  mempool: Transaction[]
}
```

**Features**:
- Interactive block chain visualization with clickable blocks
- Block details panel showing hash, merkle root, difficulty
- Transaction list for each block
- Mempool status showing pending transactions
- Copy-to-clipboard functionality with success feedback
- Responsive scrolling for wide block chains

**State Management**:
- `selectedBlock`: Currently viewed block
- `copiedHash`: Tracks which hash was copied (for feedback)

**Key Methods**:
- `copyToClipboard()`: Copies hash to clipboard with visual feedback

### 2. WalletManager

**Location**: `components/wallet-manager.tsx`

**Purpose**: Creates and displays all wallets with their balances and details.

**Props**:
```typescript
interface WalletManagerProps {
  onWalletsUpdate?: (wallets: Wallet[]) => void
  selectedWallet?: Wallet | null
  onWalletSelect?: (wallet: Wallet) => void
}
```

**Features**:
- Create new wallets with cryptographic key generation
- Display all wallets with address, balance, and nonce
- Wallet selection for transactions
- Copy wallet addresses to clipboard
- Refresh wallet list from backend
- Total balance calculation across all wallets
- Visual wallet avatars with initials

**State Management**:
- `wallets`: All wallets from backend
- `loading`: Loading state during fetch/create
- `creating`: Loading state during wallet creation
- `copiedAddress`: Tracks copied wallet addresses
- `error`: Error messages

**API Calls**:
- `getAllWallets()`: Fetch all wallets on mount
- `createWallet()`: Create new wallet

### 3. TransactionForm

**Location**: `components/transaction-form.tsx`

**Purpose**: Form for submitting transactions between wallets.

**Props**:
```typescript
interface TransactionFormProps {
  wallets: Wallet[]
  selectedWallet?: Wallet | null
  onTransactionSubmitted?: () => void
}
```

**Features**:
- Sender selection dropdown
- Recipient address input
- Amount and fee specification
- Balance validation (amount + fee ≤ balance)
- Total cost calculation display
- Success/error message display
- Form reset on successful submission
- Loading state during submission

**State Management**:
- `fromAddress`: Selected sender
- `toAddress`: Recipient address
- `amount`: Transaction amount
- `fee`: Transaction fee
- `loading`: Submission loading state
- `error`: Error messages
- `success`: Success message with transaction ID

**Validation**:
- Non-empty sender wallet
- Valid recipient address
- Positive amount and fee
- Sufficient balance for total cost

**API Calls**:
- `submitTransaction()`: Submit transaction to backend

### 4. MiningControls

**Location**: `components/mining-controls.tsx`

**Purpose**: Controls mining process with real-time status updates.

**Props**:
```typescript
interface MiningControlsProps {
  onStatusChange?: (isMining: boolean) => void
}
```

**Features**:
- Start/stop mining button
- Real-time mining status indicator
- Animated status pulse when mining active
- Status polling every 2 seconds
- Info panel explaining mining
- Error handling with messages

**State Management**:
- `isMining`: Current mining status
- `loading`: Loading state for start/stop
- `error`: Error messages

**Side Effects**:
- Sets up interval to poll mining status every 2 seconds
- Cleans up interval on unmount

**API Calls**:
- `startMining()`: Start the mining process
- `stopMining()`: Stop the mining process
- `isMining()`: Check current mining status (polling)

### 5. SettingsPanel

**Location**: `components/settings-panel.tsx`

**Purpose**: Allows users to configure the blockchain API URL.

**Features**:
- Floating settings button (bottom-right corner)
- Modal panel with API URL input
- URL validation
- localStorage persistence
- Visual feedback on save
- Instructions for deployment
- Overlay click to close

**State Management**:
- `isOpen`: Panel open/closed state
- `apiUrl`: Current API URL input value
- `saved`: Success feedback display
- `error`: URL validation errors

**Storage**:
- Saves URL to localStorage key: `blockchain_api_url`
- Loads from localStorage on component mount

**Note**: Changes take effect after page reload

### 6. Main Dashboard (app/page.tsx)

**Location**: `app/page.tsx`

**Purpose**: Main container orchestrating all components and data flow.

**State Management**:
- `blocks`: All blockchain blocks
- `mempool`: Pending transactions
- `wallets`: All wallets
- `selectedWallet`: Currently selected wallet
- `loading`: Global loading state
- `isMining`: Current mining status
- `lastUpdate`: Timestamp of last data refresh

**Data Flow**:
1. Component mounts → `loadData()`
2. `getStatus()` fetches blocks, mempool, wallets, block height
3. Sets state with fetched data
4. Components re-render with new data
5. Auto-refresh every 3 seconds via `setInterval`

**Side Effects**:
- Initial data load on mount
- 3-second auto-refresh interval
- Mining status polling (handled by MiningControls)

## API Client (lib/api.ts)

**Location**: `lib/api.ts`

**Purpose**: Type-safe wrapper around blockchain backend API.

**Key Exports**:
- `createWallet()`: Create new wallet
- `getAllWallets()`: Fetch all wallets
- `getWallet(address)`: Fetch single wallet
- `submitTransaction(from, to, amount, fee)`: Submit transaction
- `getMempool()`: Fetch pending transactions
- `getBlockchain()`: Fetch all blocks
- `getBlock(index)`: Fetch single block
- `getBlockHeight()`: Fetch blockchain height
- `startMining()`: Start mining
- `stopMining()`: Stop mining
- `isMining()`: Check mining status

**Base URL**:
Constructed from `NEXT_PUBLIC_BLOCKCHAIN_API_URL` environment variable, defaults to `http://localhost:3000`

**Error Handling**:
- Throws errors on failed requests
- Errors are caught and logged in components

## Component Hierarchy

```
Dashboard (app/page.tsx)
├── Header
│   ├── Title
│   ├── Last Update Time
│   ├── Refresh Button
│   └── Mining Indicator
├── Main Layout (Grid)
│   ├── Left Sidebar (1/3)
│   │   ├── WalletManager
│   │   ├── MiningControls
│   │   └── NetworkStats
│   └── Main Content (2/3)
│       ├── TransactionForm
│       └── BlockchainVisualizer
└── SettingsPanel (Portal)
```

## Data Flow Diagram

```
Backend API
    ↓
lib/api.ts (Type-safe client)
    ↓
app/page.tsx (Dashboard, state management)
    ├── → WalletManager (display wallets)
    ├── → TransactionForm (submit transactions)
    ├── → BlockchainVisualizer (display blocks)
    ├── → MiningControls (control mining)
    └── → SettingsPanel (configure API)
    ↑
Auto-refresh (3 seconds)
```

## State Management Pattern

The dashboard uses React hooks for state management:

1. **useState**: Local component state
2. **useEffect**: Side effects (API calls, intervals)
3. **useCallback**: Memoized functions to prevent unnecessary re-renders
4. **Props**: Pass data and callbacks to child components

**No external state management** (Redux, Zustand, etc.) - hooks are sufficient for this app.

## Styling

### Design System
- **Framework**: Tailwind CSS v4
- **Components**: shadcn/ui (Button, etc.)
- **Colors**: Black/white/gray theme with semantic tokens
- **Responsive**: Mobile-first with breakpoints

### Tailwind Patterns Used
- `flex`, `grid`: Layouts
- `gap-*`: Spacing
- `bg-*`, `text-*`, `border-*`: Colors using design tokens
- `rounded-lg`: Borders
- `hover:`, `focus:`: Interactive states
- `animate-*`: Animations (spin, pulse)

### CSS Custom Properties
Define semantic color tokens in `globals.css`:
- `--background`: Page background
- `--foreground`: Text color
- `--card`: Card background
- `--primary`: Primary brand color
- `--muted`: Muted elements
- etc.

## Performance Considerations

1. **Auto-refresh**: 3-second interval (balance between freshness and API load)
2. **Mining status polling**: 2-second interval (separate from main refresh)
3. **useCallback**: Prevents child components from unnecessary re-renders
4. **Lazy loading**: Components render only when mounted
5. **Copy feedback**: Clears after 2 seconds to reduce memory usage

## Accessibility

- Semantic HTML elements (main, header, button)
- ARIA labels and roles from shadcn/ui
- Keyboard navigation support
- Color contrast meets WCAG standards
- Screen reader friendly text

## Error Handling Strategy

1. **API errors**: Caught and displayed to user via alert components
2. **Form validation**: Client-side validation before submission
3. **Empty states**: Helpful messages when no data available
4. **Loading states**: Buttons disabled, spinners shown
5. **Success feedback**: Temporary messages (auto-clear after 2-5 seconds)

## Testing

To test components in isolation:

```typescript
// Example: Test WalletManager
import { WalletManager } from '@/components/wallet-manager'

export default function TestPage() {
  const [wallets, setWallets] = useState([
    {
      address: '0x123...',
      public_key: 'abc...',
      balance: 100,
      nonce: 0
    }
  ])

  return (
    <WalletManager
      onWalletsUpdate={setWallets}
      onWalletSelect={(w) => console.log('Selected:', w)}
    />
  )
}
```

## Future Enhancements

- [ ] Add TypeScript strict mode
- [ ] Add unit tests (Vitest)
- [ ] Add E2E tests (Playwright)
- [ ] WebSocket support for real-time updates
- [ ] Redux/Zustand for complex state
- [ ] Transaction history/search
- [ ] Dark mode toggle
- [ ] Export/import wallet functionality

## Debugging Tips

1. **Check Console**: Open DevTools (F12) → Console for errors
2. **Check Network**: DevTools → Network to see API calls
3. **React DevTools**: Install React DevTools extension to inspect component tree
4. **API Status**: Verify backend is running: `curl https://your-api/api/wallet/all`
5. **Environment**: Check `.env.local` for correct API URL

---

For more information, see [README.md](./README.md) and [CONFIGURATION.md](./CONFIGURATION.md)
