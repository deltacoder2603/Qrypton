'use client'

import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Settings, Check, AlertCircle } from 'lucide-react'

export function SettingsPanel() {
  const [isOpen, setIsOpen] = useState(false)
  const [apiUrl, setApiUrl] = useState(
    process.env.NEXT_PUBLIC_BLOCKCHAIN_API_URL || 'http://localhost:3000'
  )
  const [saved, setSaved] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [testing, setTesting] = useState(false)
  const [connectionStatus, setConnectionStatus] = useState<string | null>(null)

  useEffect(() => {
    // Load from localStorage
    const saved = localStorage.getItem('blockchain_api_url')
    if (saved) {
      setApiUrl(saved)
    }
  }, [])

  const handleSave = () => {
    try {
      setError(null)

      // Basic URL validation
      new URL(apiUrl)

      // Save to localStorage
      localStorage.setItem('blockchain_api_url', apiUrl)
      setSaved(true)

      // Clear saved message after 2 seconds
      setTimeout(() => setSaved(false), 2000)
    } catch (err) {
      setError('Invalid URL format')
    }
  }

  const handleTestConnection = async () => {
    try {
      setError(null)
      setConnectionStatus(null)
      setTesting(true)

      // Test the connection
      const testUrl = new URL('/api/blockchain/height', apiUrl).toString()
      const response = await fetch(testUrl, {
        method: 'GET',
      })

      if (response.ok) {
        setConnectionStatus('✓ Connection successful!')
        setTimeout(() => setConnectionStatus(null), 3000)
      } else {
        setError(`Server returned status ${response.status}`)
      }
    } catch (err) {
      setError(
        'Connection failed. Make sure the URL is correct and the server is running.'
      )
    } finally {
      setTesting(false)
    }
  }

  return (
    <>
      {/* Settings Button */}
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="fixed bottom-6 right-6 w-12 h-12 rounded-full bg-primary text-primary-foreground shadow-lg hover:bg-primary/90 flex items-center justify-center transition-all hover:scale-110 z-40"
        title="Settings"
      >
        <Settings className="w-5 h-5" />
      </button>

      {/* Settings Panel */}
      {isOpen && (
        <div className="fixed bottom-20 right-6 w-96 max-w-[calc(100vw-2rem)] bg-card border border-border rounded-lg shadow-lg p-6 z-40 space-y-4">
          <div>
            <h3 className="text-lg font-bold text-foreground mb-2">Settings</h3>
            <p className="text-sm text-muted-foreground">
              Configure your blockchain backend connection
            </p>
          </div>

          {/* API URL Input */}
          <div>
            <label className="block text-sm font-medium text-foreground mb-2">
              Blockchain API URL
            </label>
            <input
              type="text"
              value={apiUrl}
              onChange={(e) => setApiUrl(e.target.value)}
              placeholder="https://example.com:3000"
              className="w-full px-3 py-2 border border-input bg-background text-foreground rounded-lg focus:outline-none focus:ring-2 focus:ring-primary placeholder:text-muted-foreground text-sm"
            />
            <p className="text-xs text-muted-foreground mt-1">
              Example: https://txhm5mlj-3000.inc1.devtunnels.ms
            </p>
          </div>

          {/* Current Value */}
          <div className="bg-muted rounded-lg p-3">
            <p className="text-xs text-muted-foreground">Current value</p>
            <p className="font-mono text-xs text-foreground break-all">
              {apiUrl}
            </p>
          </div>

          {/* Error Message */}
          {error && (
            <div className="bg-red-50/10 border border-red-600/30 rounded-lg p-3 flex items-start gap-2">
              <AlertCircle className="w-4 h-4 text-red-600 mt-0.5 flex-shrink-0" />
              <p className="text-xs text-red-600">{error}</p>
            </div>
          )}

          {/* Success Message */}
          {saved && (
            <div className="bg-green-50/10 border border-green-600/30 rounded-lg p-3 flex items-start gap-2">
              <Check className="w-4 h-4 text-green-600 mt-0.5 flex-shrink-0" />
              <p className="text-xs text-green-600">Settings saved!</p>
            </div>
          )}

          {/* Connection Status */}
          {connectionStatus && (
            <div className="bg-green-50/10 border border-green-600/30 rounded-lg p-3 flex items-start gap-2">
              <Check className="w-4 h-4 text-green-600 mt-0.5 flex-shrink-0" />
              <p className="text-xs text-green-600">{connectionStatus}</p>
            </div>
          )}

          {/* Note */}
          <div className="bg-primary/10 border border-primary/20 rounded-lg p-3">
            <p className="text-xs text-primary">
              ⚠️ Changes will only take effect after page reload
            </p>
          </div>

          {/* Buttons */}
          <div className="flex gap-2">
            <Button
              onClick={handleTestConnection}
              disabled={testing}
              variant="outline"
              className="flex-1"
            >
              {testing ? 'Testing...' : 'Test Connection'}
            </Button>
            <Button onClick={handleSave} className="flex-1">
              Save
            </Button>
          </div>

          <Button
            onClick={() => setIsOpen(false)}
            variant="outline"
            className="w-full"
          >
            Close
          </Button>
        </div>
      )}

      {/* Overlay */}
      {isOpen && (
        <div
          className="fixed inset-0 z-30"
          onClick={() => setIsOpen(false)}
        />
      )}
    </>
  )
}
