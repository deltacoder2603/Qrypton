'use client'

import { startMining, stopMining } from '@/lib/api'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Loader, Zap } from 'lucide-react'

interface MiningControlsProps {
  onStatusChange?: (isMining: boolean) => void
}

export function MiningControls({ onStatusChange }: MiningControlsProps) {
  const [isMining, setIsMining] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleStartMining = async () => {
    try {
      setLoading(true)
      setError(null)
      await startMining()
      setIsMining(true)
      onStatusChange?.(true)
    } catch (err) {
      setError('Failed to start mining')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  const handleStopMining = async () => {
    try {
      setLoading(true)
      setError(null)
      await stopMining()
      setIsMining(false)
      onStatusChange?.(false)
    } catch (err) {
      setError('Failed to stop mining')
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="border border-border rounded-lg p-6 bg-card">
      <div className="space-y-4">
        <div>
          <h2 className="text-2xl font-bold text-foreground">Mining</h2>
          <p className="text-sm text-muted-foreground mt-1">
            Control the blockchain mining process
          </p>
        </div>

        {/* Status Indicator */}
        <div
          className={`flex items-center gap-3 p-3 rounded-lg ${
            isMining
              ? 'bg-green-50/10 border border-green-600/30'
              : 'bg-muted border border-border'
          }`}
        >
          <div
            className={`w-3 h-3 rounded-full ${
              isMining ? 'bg-green-600 animate-pulse' : 'bg-muted-foreground'
            }`}
          />
          <span className={isMining ? 'text-green-600 font-bold' : 'text-muted-foreground'}>
            {isMining ? 'Mining Active' : 'Mining Inactive'}
          </span>
        </div>

        {/* Error Message */}
        {error && (
          <div className="bg-red-50/10 border border-red-600/30 rounded-lg p-3 text-sm text-red-600">
            {error}
          </div>
        )}

        {/* Controls */}
        <div className="flex gap-2">
          {!isMining ? (
            <Button
              onClick={handleStartMining}
              disabled={loading}
              className="flex-1 gap-2"
            >
              {loading ? (
                <>
                  <Loader className="w-4 h-4 animate-spin" />
                  Starting...
                </>
              ) : (
                <>
                  <Zap className="w-4 h-4" />
                  Start Mining
                </>
              )}
            </Button>
          ) : (
            <Button
              onClick={handleStopMining}
              disabled={loading}
              variant="destructive"
              className="flex-1 gap-2"
            >
              {loading ? (
                <>
                  <Loader className="w-4 h-4 animate-spin" />
                  Stopping...
                </>
              ) : (
                <>
                  <Zap className="w-4 h-4" />
                  Stop Mining
                </>
              )}
            </Button>
          )}
        </div>

        {/* Info */}
        <div className="bg-muted rounded-lg p-3 text-xs text-muted-foreground space-y-1">
          <p>🔗 Mining adds new blocks to the blockchain</p>
          <p>💰 Block rewards go to the miner&apos;s wallet</p>
          <p>⏱️ Check status every 2 seconds</p>
        </div>
      </div>
    </div>
  )
}
