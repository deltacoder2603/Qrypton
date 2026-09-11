'use client'

import { Wallet, submitTransaction } from '@/lib/api'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Send, AlertCircle, CheckCircle2, Loader } from 'lucide-react'

interface TransactionFormProps {
  wallets: Wallet[]
  selectedWallet?: Wallet | null
  onTransactionSubmitted?: () => void
}

export function TransactionForm({
  wallets,
  selectedWallet,
  onTransactionSubmitted,
}: TransactionFormProps) {
  const [fromAddress, setFromAddress] = useState<string>(
    selectedWallet?.address || ''
  )
  const [toAddress, setToAddress] = useState<string>('')
  const [amount, setAmount] = useState<string>('')
  const [fee, setFee] = useState<string>('1')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [success, setSuccess] = useState<string | null>(null)

  const sender = wallets.find((w) => w.address === fromAddress)
  const balance = sender?.balance || 0

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setSuccess(null)

    // Validation
    if (!fromAddress) {
      setError('Please select a sender wallet')
      return
    }
    if (!toAddress) {
      setError('Please enter a recipient address')
      return
    }
    if (!amount || isNaN(Number(amount)) || Number(amount) <= 0) {
      setError('Please enter a valid amount')
      return
    }
    if (!fee || isNaN(Number(fee)) || Number(fee) < 0) {
      setError('Please enter a valid fee')
      return
    }

    const totalRequired = Number(amount) + Number(fee)
    if (totalRequired > balance) {
      setError(
        `Insufficient balance. Required: ${totalRequired}, Available: ${balance}`
      )
      return
    }

    try {
      setLoading(true)
      const result = await submitTransaction(
        fromAddress,
        toAddress,
        Number(amount),
        Number(fee)
      )
      setSuccess(`Transaction submitted successfully! ID: ${result.tx_id.substring(0, 16)}...`)
      setToAddress('')
      setAmount('')
      setFee('1')
      onTransactionSubmitted?.()

      // Clear success message after 5 seconds
      setTimeout(() => setSuccess(null), 5000)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to submit transaction')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="border border-border rounded-lg p-6 bg-card">
      <h2 className="text-2xl font-bold mb-6 text-foreground">
        Create Transaction
      </h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* From Address */}
        <div>
          <label className="block text-sm font-medium text-foreground mb-2">
            From Address
          </label>
          <select
            value={fromAddress}
            onChange={(e) => setFromAddress(e.target.value)}
            className="w-full px-3 py-2 border border-input bg-background text-foreground rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
          >
            <option value="">Select a wallet</option>
            {wallets.map((wallet) => (
              <option key={wallet.address} value={wallet.address}>
                {wallet.address} (Balance: {wallet.balance})
              </option>
            ))}
          </select>
        </div>

        {/* Sender Info */}
        {sender && (
          <div className="bg-muted p-3 rounded-lg text-sm">
            <p className="text-muted-foreground">Current Balance</p>
            <p className="font-bold text-lg text-foreground">{sender.balance} units</p>
          </div>
        )}

        {/* To Address */}
        <div>
          <label className="block text-sm font-medium text-foreground mb-2">
            To Address
          </label>
          <input
            type="text"
            value={toAddress}
            onChange={(e) => setToAddress(e.target.value)}
            placeholder="Enter recipient address"
            className="w-full px-3 py-2 border border-input bg-background text-foreground rounded-lg focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-muted-foreground"
          />
        </div>

        {/* Amount */}
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-foreground mb-2">
              Amount
            </label>
            <input
              type="number"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              placeholder="0"
              min="1"
              className="w-full px-3 py-2 border border-input bg-background text-foreground rounded-lg focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-muted-foreground"
            />
          </div>

          {/* Fee */}
          <div>
            <label className="block text-sm font-medium text-foreground mb-2">
              Fee
            </label>
            <input
              type="number"
              value={fee}
              onChange={(e) => setFee(e.target.value)}
              placeholder="1"
              min="0"
              className="w-full px-3 py-2 border border-input bg-background text-foreground rounded-lg focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-muted-foreground"
            />
          </div>
        </div>

        {/* Total Cost */}
        {amount && (
          <div className="bg-primary/10 border border-primary/20 rounded-lg p-3">
            <div className="flex justify-between text-sm">
              <span className="text-muted-foreground">Amount:</span>
              <span className="font-bold text-foreground">{amount} units</span>
            </div>
            <div className="flex justify-between text-sm mt-1">
              <span className="text-muted-foreground">Fee:</span>
              <span className="font-bold text-foreground">{fee} units</span>
            </div>
            <div className="border-t border-primary/20 mt-2 pt-2 flex justify-between text-sm">
              <span className="font-bold text-foreground">Total:</span>
              <span className="font-bold text-primary">
                {Number(amount) + Number(fee)} units
              </span>
            </div>
          </div>
        )}

        {/* Error Message */}
        {error && (
          <div className="bg-red-50/10 border border-red-600/30 rounded-lg p-3 flex items-start gap-2">
            <AlertCircle className="w-5 h-5 text-red-600 mt-0.5 flex-shrink-0" />
            <p className="text-sm text-red-600">{error}</p>
          </div>
        )}

        {/* Success Message */}
        {success && (
          <div className="bg-green-50/10 border border-green-600/30 rounded-lg p-3 flex items-start gap-2">
            <CheckCircle2 className="w-5 h-5 text-green-600 mt-0.5 flex-shrink-0" />
            <p className="text-sm text-green-600">{success}</p>
          </div>
        )}

        {/* Submit Button */}
        <Button
          type="submit"
          disabled={loading || !fromAddress}
          className="w-full gap-2"
        >
          {loading ? (
            <>
              <Loader className="w-4 h-4 animate-spin" />
              Submitting...
            </>
          ) : (
            <>
              <Send className="w-4 h-4" />
              Submit Transaction
            </>
          )}
        </Button>
      </form>
    </div>
  )
}
