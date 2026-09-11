'use client'

import { useState, useEffect, useCallback } from 'react'
import { getStatus, Wallet, Block, Transaction } from '@/lib/api'
import { WalletManager } from '@/components/wallet-manager'
import { TransactionForm } from '@/components/transaction-form'
import { BlockchainVisualizer } from '@/components/blockchain-visualizer'
import { MiningControls } from '@/components/mining-controls'
import { SettingsPanel } from '@/components/settings-panel'
import { RefreshCw, Zap } from 'lucide-react'
import { Button } from '@/components/ui/button'

export default function Dashboard() {
  const [blocks, setBlocks] = useState<Block[]>([])
  const [mempool, setMempool] = useState<Transaction[]>([])
  const [wallets, setWallets] = useState<Wallet[]>([])
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null)
  const [loading, setLoading] = useState(false)
  const [isMining, setIsMining] = useState(false)
  const [lastUpdate, setLastUpdate] = useState<Date | null>(null)

  const loadData = useCallback(async () => {
    try {
      setLoading(true)
      const status = await getStatus()
      setBlocks(status.blocks)
      setMempool(status.mempool)
      setWallets(status.wallets)
      setLastUpdate(new Date())
    } catch (error) {
      console.error('Failed to load data:', error)
    } finally {
      setLoading(false)
    }
  }, [])

  // Initial load
  useEffect(() => {
    loadData()
  }, [loadData])

  // Auto-refresh every 3 seconds
  useEffect(() => {
    const interval = setInterval(loadData, 3000)
    return () => clearInterval(interval)
  }, [loadData])

  // Handle wallet selection
  const handleWalletSelect = (wallet: Wallet) => {
    setSelectedWallet(wallet)
  }

  const handleWalletsUpdate = (updated: Wallet[]) => {
    setWallets(updated)
  }

  return (
    <div className="min-h-screen bg-background text-foreground">
      {/* Header */}
      <header className="border-b border-border sticky top-0 z-30 bg-background/95 backdrop-blur">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold tracking-tight">
              Blockchain Dashboard
            </h1>
            <p className="text-sm text-muted-foreground mt-1">
              Real-time blockchain visualization and management
            </p>
          </div>

          <div className="flex items-center gap-3">
            {lastUpdate && (
              <div className="text-xs text-muted-foreground text-right">
                <p>Last updated</p>
                <p className="font-mono">
                  {lastUpdate.toLocaleTimeString()}
                </p>
              </div>
            )}

            <Button
              onClick={loadData}
              disabled={loading}
              variant="outline"
              size="icon"
              className="w-10 h-10"
            >
              <RefreshCw className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
            </Button>

            {isMining && (
              <div className="flex items-center gap-2 bg-green-50/10 border border-green-600/30 px-3 py-2 rounded-lg">
                <Zap className="w-4 h-4 text-green-600 animate-pulse" />
                <span className="text-xs font-bold text-green-600">Mining</span>
              </div>
            )}
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Left Sidebar */}
          <div className="lg:col-span-1 space-y-6 order-2 lg:order-1">
            {/* Wallet Manager */}
            <div className="border border-border rounded-lg p-6 bg-card">
              <WalletManager
                onWalletsUpdate={handleWalletsUpdate}
                selectedWallet={selectedWallet}
                onWalletSelect={handleWalletSelect}
              />
            </div>

            {/* Mining Controls */}
            <MiningControls onStatusChange={setIsMining} />

            {/* Stats */}
            <div className="border border-border rounded-lg p-6 bg-card space-y-3">
              <h3 className="font-bold text-foreground">Network Stats</h3>
              <div className="grid grid-cols-2 gap-2 text-sm">
                <div className="bg-muted p-3 rounded">
                  <p className="text-muted-foreground text-xs">Total Wallets</p>
                  <p className="font-bold text-lg text-foreground">
                    {wallets.length}
                  </p>
                </div>
                <div className="bg-muted p-3 rounded">
                  <p className="text-muted-foreground text-xs">Block Height</p>
                  <p className="font-bold text-lg text-foreground">
                    {blocks.length}
                  </p>
                </div>
                <div className="bg-muted p-3 rounded col-span-2">
                  <p className="text-muted-foreground text-xs">Pending Tx</p>
                  <p className="font-bold text-lg text-foreground">
                    {mempool.length}
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Main Content Area */}
          <div className="lg:col-span-2 space-y-6 order-1 lg:order-2">
            {/* Transaction Form */}
            {wallets.length > 0 ? (
              <TransactionForm
                wallets={wallets}
                selectedWallet={selectedWallet}
                onTransactionSubmitted={loadData}
              />
            ) : (
              <div className="border border-dashed border-border rounded-lg p-8 bg-card text-center">
                <p className="text-muted-foreground font-medium">
                  Create a wallet to get started
                </p>
              </div>
            )}

            {/* Blockchain Visualizer */}
            <BlockchainVisualizer blocks={blocks} mempool={mempool} />
          </div>
        </div>
      </main>

      {/* Settings Panel */}
      <SettingsPanel />
    </div>
  )
}
