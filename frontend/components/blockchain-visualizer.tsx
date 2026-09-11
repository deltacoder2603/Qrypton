'use client'

import { Block, Transaction } from '@/lib/api'
import { useState, useEffect } from 'react'
import { Copy, CheckCircle2, AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'

interface BlockVisualizerProps {
  blocks: Block[]
  mempool: Transaction[]
}

export function BlockchainVisualizer({ blocks, mempool }: BlockVisualizerProps) {
  const [selectedBlock, setSelectedBlock] = useState<Block | null>(null)
  const [copiedHash, setCopiedHash] = useState<string | null>(null)

  useEffect(() => {
    if (blocks.length > 0) {
      setSelectedBlock(blocks[blocks.length - 1])
    }
  }, [blocks])

  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text)
    setCopiedHash(id)
    setTimeout(() => setCopiedHash(null), 2000)
  }

  return (
    <div className="space-y-6">
      {/* Blockchain Chain Visualization */}
      <div className="border border-border rounded-lg p-6 bg-card">
        <h2 className="text-2xl font-bold mb-6 text-foreground">Blockchain Chain</h2>

        <div className="overflow-x-auto pb-4">
          <div className="flex gap-3 min-w-max px-2">
            {blocks.map((block, index) => (
              <div key={block.index} className="flex items-center">
                <button
                  onClick={() => setSelectedBlock(block)}
                  className={`px-4 py-3 rounded-lg font-mono text-sm font-bold transition-all transform hover:scale-105 ${
                    selectedBlock?.index === block.index
                      ? 'bg-primary text-primary-foreground shadow-lg scale-105'
                      : 'bg-muted text-muted-foreground hover:bg-accent hover:text-accent-foreground'
                  }`}
                >
                  Block #{block.index}
                </button>
                {index < blocks.length - 1 && (
                  <div className="mx-2 text-muted-foreground">→</div>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* Genesis Block Indicator */}
        {blocks.length > 0 && (
          <div className="mt-4 pt-4 border-t border-border">
            <p className="text-xs text-muted-foreground">
              Genesis Block: {blocks[0].hash.substring(0, 16)}...
            </p>
          </div>
        )}
      </div>

      {/* Block Details */}
      {selectedBlock && (
        <div className="border border-border rounded-lg p-6 bg-card space-y-4">
          <div className="flex justify-between items-start">
            <div>
              <h3 className="text-xl font-bold text-foreground">
                Block #{selectedBlock.index}
              </h3>
              <p className="text-sm text-muted-foreground">
                Timestamp: {new Date(selectedBlock.timestamp).toLocaleString()}
              </p>
            </div>
            <div className="text-right">
              <p className="text-sm font-mono text-muted-foreground">
                Difficulty: {selectedBlock.difficulty}
              </p>
              <p className="text-sm font-mono text-muted-foreground">
                Nonce: {selectedBlock.nonce}
              </p>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-3">
            <div className="bg-muted p-3 rounded font-mono text-xs break-all">
              <p className="text-muted-foreground mb-1">Hash:</p>
              <div className="flex items-center gap-2">
                <span className="text-foreground">{selectedBlock.hash}</span>
                <button
                  onClick={() => copyToClipboard(selectedBlock.hash, 'hash')}
                  className="p-1 hover:bg-background rounded transition-colors"
                >
                  {copiedHash === 'hash' ? (
                    <CheckCircle2 className="w-4 h-4 text-green-600" />
                  ) : (
                    <Copy className="w-4 h-4 text-muted-foreground" />
                  )}
                </button>
              </div>
            </div>

            <div className="bg-muted p-3 rounded font-mono text-xs break-all">
              <p className="text-muted-foreground mb-1">Previous Hash:</p>
              <div className="flex items-center gap-2">
                <span className="text-foreground">
                  {selectedBlock.previous_hash}
                </span>
                <button
                  onClick={() =>
                    copyToClipboard(selectedBlock.previous_hash, 'prev')
                  }
                  className="p-1 hover:bg-background rounded transition-colors"
                >
                  {copiedHash === 'prev' ? (
                    <CheckCircle2 className="w-4 h-4 text-green-600" />
                  ) : (
                    <Copy className="w-4 h-4 text-muted-foreground" />
                  )}
                </button>
              </div>
            </div>

            <div className="bg-muted p-3 rounded font-mono text-xs break-all">
              <p className="text-muted-foreground mb-1">Merkle Root:</p>
              <div className="flex items-center gap-2">
                <span className="text-foreground">{selectedBlock.merkle_root}</span>
                <button
                  onClick={() =>
                    copyToClipboard(selectedBlock.merkle_root, 'merkle')
                  }
                  className="p-1 hover:bg-background rounded transition-colors"
                >
                  {copiedHash === 'merkle' ? (
                    <CheckCircle2 className="w-4 h-4 text-green-600" />
                  ) : (
                    <Copy className="w-4 h-4 text-muted-foreground" />
                  )}
                </button>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-2 gap-2 text-sm">
            <div className="bg-muted p-2 rounded">
              <p className="text-muted-foreground">Miner</p>
              <p className="font-mono text-xs truncate text-foreground">
                {selectedBlock.miner}
              </p>
            </div>
            <div className="bg-muted p-2 rounded">
              <p className="text-muted-foreground">Transactions</p>
              <p className="font-bold text-foreground">
                {selectedBlock.transactions.length}
              </p>
            </div>
          </div>

          {/* Transactions in Block */}
          {selectedBlock.transactions.length > 0 && (
            <div className="mt-4 pt-4 border-t border-border">
              <h4 className="font-bold text-foreground mb-3">
                Transactions ({selectedBlock.transactions.length})
              </h4>
              <div className="space-y-2 max-h-64 overflow-y-auto">
                {selectedBlock.transactions.map((tx, idx) => (
                  <div
                    key={tx.id}
                    className="bg-muted p-2 rounded text-xs font-mono"
                  >
                    <div className="flex items-center gap-2 mb-1">
                      <span className="text-muted-foreground">Tx {idx + 1}:</span>
                      <span className="text-foreground truncate flex-1">
                        {tx.id}
                      </span>
                    </div>
                    <div className="text-muted-foreground text-xs space-y-1 ml-4">
                      <p>From: {tx.from.substring(0, 20)}...</p>
                      <p>To: {tx.to.substring(0, 20)}...</p>
                      <p>
                        Amount: {tx.amount} | Fee: {tx.fee}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Mempool Status */}
      {mempool.length > 0 && (
        <div className="border border-yellow-600/30 bg-yellow-50/5 rounded-lg p-6">
          <div className="flex items-start gap-3">
            <AlertCircle className="w-5 h-5 text-yellow-600 mt-0.5 flex-shrink-0" />
            <div className="flex-1">
              <h4 className="font-bold text-foreground">
                Pending Transactions ({mempool.length})
              </h4>
              <p className="text-sm text-muted-foreground mt-1">
                These transactions are waiting to be included in a block.
              </p>
              <div className="mt-3 space-y-2 max-h-48 overflow-y-auto">
                {mempool.map((tx) => (
                  <div
                    key={tx.id}
                    className="bg-muted p-2 rounded text-xs font-mono"
                  >
                    <div className="flex items-center justify-between">
                      <span className="text-foreground truncate flex-1">
                        {tx.id}
                      </span>
                      <span className="text-muted-foreground ml-2">
                        {tx.amount} units
                      </span>
                    </div>
                    <p className="text-muted-foreground text-xs mt-1">
                      {tx.from.substring(0, 10)}... → {tx.to.substring(0, 10)}...
                    </p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
