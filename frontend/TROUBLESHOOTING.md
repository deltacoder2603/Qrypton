# Troubleshooting Guide

## "Failed to load wallets" Error

This error appears when the dashboard cannot connect to your blockchain backend API. Follow these steps to resolve it.

### Step 1: Verify Your Backend is Running

The blockchain backend must be running and accessible at the URL configured in the Settings panel.

```bash
# Check if your backend is running
curl https://txhm5mlj-3000.inc1.devtunnels.ms/api/blockchain/height

# You should get a response like:
# {"height": 1}
```

### Step 2: Configure the API URL

1. Click the **Settings** button (gear icon) in the bottom right corner
2. You'll see the "Blockchain API URL" field
3. Verify it matches your backend URL:
   - Default: `http://localhost:3000`
   - Your backend: `https://txhm5mlj-3000.inc1.devtunnels.ms`

### Step 3: Test the Connection

1. In the Settings panel, click the **Test Connection** button
2. If successful, you'll see "✓ Connection successful!"
3. If it fails, check:
   - The URL is correct and accessible
   - Your firewall/network allows the connection
   - The backend API is running

### Step 4: Save and Reload

1. Click **Save** to persist your settings
2. Close the Settings panel
3. **Reload the page** (Ctrl+R or Cmd+R)
4. The dashboard should now load your wallet data

## Common Issues and Solutions

### Issue: "Connection failed" in Test Connection

**Problem:** The backend is not running or the URL is incorrect

**Solution:**
1. Verify your backend is running
2. Check that the URL is exactly correct (including `https://` vs `http://`)
3. Make sure there are no trailing slashes in the URL
4. If using a tunnel/dev tunnel, verify it's still active

### Issue: Wallets Load But Mining Fails

**Problem:** Mining endpoint isn't available

**Solution:**
```bash
# Verify mining endpoint is available
curl https://txhm5mlj-3000.inc1.devtunnels.ms/api/miner/status

# Should return: {"is_mining": false}
```

### Issue: Transactions Not Submitting

**Problem:** The `/api/transaction/submit` endpoint is failing

**Solution:**
1. Check browser console for detailed error message
2. Verify sender wallet has sufficient balance
3. Ensure the transaction fee is reasonable (typically 1 unit)
4. Confirm both sender and receiver addresses are valid

### Issue: Settings Don't Persist After Reload

**Problem:** Browser localStorage is disabled or blocked

**Solution:**
1. Check browser privacy settings
2. Make sure you're not in private/incognito mode
3. Try a different browser
4. The URL can also be set via `.env.local` file:
   ```
   NEXT_PUBLIC_BLOCKCHAIN_API_URL=https://txhm5mlj-3000.inc1.devtunnels.ms
   ```

## Backend API Endpoints Required

Your blockchain backend must provide these endpoints:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/wallet/all` | GET | Get all wallets |
| `/api/wallet/create` | POST | Create a new wallet |
| `/api/wallet/{address}` | GET | Get specific wallet |
| `/api/transaction/submit` | POST | Submit a transaction |
| `/api/mempool` | GET | Get pending transactions |
| `/api/blockchain` | GET | Get all blocks |
| `/api/blockchain/{index}` | GET | Get specific block |
| `/api/blockchain/height` | GET | Get block height |
| `/api/miner/start` | POST | Start mining |
| `/api/miner/stop` | POST | Stop mining |
| `/api/miner/status` | GET | Get mining status |

## Debugging Tips

### Check Browser Console

Open your browser's Developer Tools (F12) and look at the Console tab for error messages:

```javascript
// You'll see messages like:
[v0] Failed to get status: ...
[v0] API URL: https://txhm5mlj-3000.inc1.devtunnels.ms
```

### Test API Endpoints Manually

```bash
# Test wallet creation
curl -X POST https://txhm5mlj-3000.inc1.devtunnels.ms/api/wallet/create \
  -H "Content-Type: application/json"

# Test getting wallets
curl https://txhm5mlj-3000.inc1.devtunnels.ms/api/wallet/all

# Test blockchain height
curl https://txhm5mlj-3000.inc1.devtunnels.ms/api/blockchain/height
```

### Enable Debug Logging

The dashboard automatically logs debugging information. Check:
1. Browser Console (F12)
2. Look for `[v0]` prefixed messages
3. These include the API URL being used and error details

## Network Issues

### CORS (Cross-Origin Resource Sharing) Errors

If you see CORS errors in the console, your backend needs to:

1. Allow requests from your dashboard URL
2. Set proper CORS headers:
   ```
   Access-Control-Allow-Origin: *
   Access-Control-Allow-Methods: GET, POST, OPTIONS
   Access-Control-Allow-Headers: Content-Type
   ```

### Timeout Issues

If requests timeout:
1. The backend might be slow or unavailable
2. Check network connectivity
3. Verify the URL is correct and accessible
4. Increase backend timeout if needed

## Getting Help

If you've tried all these steps and still have issues:

1. Check the browser console for detailed error messages
2. Verify all backend API endpoints are implemented
3. Test endpoints with `curl` or Postman
4. Check that your blockchain backend is running correctly

## Quick Start Checklist

- [ ] Backend is running
- [ ] Backend URL is correct in settings
- [ ] Test Connection shows "✓ Connection successful!"
- [ ] Page has been reloaded after saving settings
- [ ] Wallet endpoints are responding with data
- [ ] No CORS errors in browser console
- [ ] Browser localStorage is enabled
