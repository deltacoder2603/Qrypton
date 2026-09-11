'use client'

import { Wallet, createWallet, getAllWallets } from '@/lib/api'
import { useState, useEffect } from 'react'
import { Copy, CheckCircle2, Plus, RefreshCw } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface WalletManagerProps {
  onWalletsUpdate?: (wallets: Wallet[]) => void
  selectedWallet?: Wallet | null
  onWalletSelect?: (wallet: Wallet) => void
}

export function WalletManager({
  onWalletsUpdate,
  selectedWallet,
  onWalletSelect,
}: WalletManagerProps) {
  const [wallets, setWallets] = useState<Wallet[]>([])
  const [loading, setLoading] = useState(false)
  const [creating, setCreating] = useState(false)
  const [copiedAddress, setCopiedAddress] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    loadWallets()
  }, [])

  const loadWallets = async () => {
    try {
      setLoading(true)
      setError(null)
      const fetchedWallets = await getAllWallets()
      setWallets(fetchedWallets)
      onWalletsUpdate?.(fetchedWallets)
    } catch (err) {
      setError('Failed to load wallets')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleCreateWallet = async () => {
    try {
      setCreating(true)
      setError(null)
      const newWallet = await createWallet()
      setWallets([...wallets, newWallet])
      onWalletsUpdate?.([...wallets, newWallet])
    } catch (err) {
      setError('Failed to create wallet')
      console.error(err)
    } finally {
      setCreating(false)
    }
  }

  const copyToClipboard = (address: string) => {
    navigator.clipboard.writeText(address)
    setCopiedAddress(address)
    setTimeout(() => setCopiedAddress(null), 2000)
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold text-foreground">Wallets</h2>
        <div className="flex gap-2">
          <Button
            onClick={loadWallets}
            disabled={loading}
            variant="outline"
            size="sm"
            className="gap-2"
          >
            <RefreshCw className="w-4 h-4" />
            Refresh
          </Button>
          <Button
            onClick={handleCreateWallet}
            disabled={creating || loading}
            className="gap-2"
            size="sm"
          >
            <Plus className="w-4 h-4" />
            New Wallet
          </Button>
        </div>
      </div>

      {error && (
        <div className="bg-red-50/10 border border-red-600/30 rounded-lg p-3 text-sm text-red-600 space-y-1">
          <p className="font-bold">{error}</p>
          <p className="text-xs opacity-80">Make sure your blockchain backend is running at the URL configured in settings.</p>
        </div>
      )}

      {wallets.length === 0 ? (
        <div className="border border-dashed border-border rounded-lg p-8 text-center">
          <p className="text-muted-foreground">No wallets yet</p>
          <p className="text-sm text-muted-foreground mt-1">
            Create your first wallet to get started
          </p>
        </div>
      ) : (
        <div className="grid grid-cols-1 gap-3">
          {wallets.map((wallet) => (
            <div
              key={wallet.address}
              onClick={() => onWalletSelect?.(wallet)}
              className={`border rounded-lg p-4 cursor-pointer transition-all ${
                selectedWallet?.address === wallet.address
                  ? 'border-primary bg-primary/5'
                  : 'border-border hover:border-primary hover:bg-muted/50'
              }`}
            >
              <div className="flex items-start justify-between gap-3">
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2 mb-2">
                    <p className="font-mono font-bold text-sm text-foreground truncate">
                      {wallet.address}
                    </p>
                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        copyToClipboard(wallet.address)
                      }}
                      className="p-1 hover:bg-background rounded transition-colors flex-shrink-0"
                    >
                      {copiedAddress === wallet.address ? (
                        <CheckCircle2 className="w-4 h-4 text-green-600" />
                      ) : (
                        <Copy className="w-4 h-4 text-muted-foreground" />
                      )}
                    </button>
                  </div>

                  <div className="grid grid-cols-2 gap-2 text-xs">
                    <div>
                      <p className="text-muted-foreground">Balance</p>
                      <p className="font-bold text-foreground">
                        {wallet.balance} units
                      </p>
                    </div>
                    <div>
                      <p className="text-muted-foreground">Nonce</p>
                      <p className="font-bold text-foreground">{wallet.nonce}</p>
                    </div>
                  </div>
                </div>

                <div className="flex-shrink-0">
                  <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-primary/20 to-primary/5 flex items-center justify-center">
                    <span className="text-xs font-bold text-primary">
                      {wallet.address.substring(0, 2).toUpperCase()}
                    </span>
                  </div>
                </div>
              </div>

              {/* Public Key Preview */}
              <div className="mt-3 pt-3 border-t border-border/50">
                <p className="text-xs text-muted-foreground mb-1">Public Key</p>
                <p className="font-mono text-xs text-foreground truncate">
                  {wallet.public_key}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}

      {wallets.length > 0 && (
        <div className="bg-muted rounded-lg p-3 text-xs text-muted-foreground">
          <p>Total Wallets: {wallets.length}</p>
          <p className="mt-1">
            Total Balance:{' '}
            <span className="font-bold text-foreground">
              {wallets.reduce((sum, w) => sum + w.balance, 0)} units
            </span>
          </p>
        </div>
      )}
    </div>
  )
}
