# Zee-Ubot Setup Script
# This script automates the login and session encoding process.

Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan
Write-Host "        ZEE-UBOT AUTO SETUP        " -ForegroundColor Cyan
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Cyan

# 1. Ensure .env exists
if (-not (Test-Path ".env")) {
    if (Test-Path ".env.example") {
        Write-Host "[*] Creating .env from .env.example..."
        Copy-Item ".env.example" ".env"
    } else {
        Write-Host "[!] Error: .env.example not found!" -ForegroundColor Red
        exit 1
    }
}

# 2. Run the dedicated login tool
Write-Host "[*] Launching login tool..." -ForegroundColor Yellow
go run ./cmd/login

if ($LASTEXITCODE -ne 0) {
    Write-Host "[!] Login failed or cancelled." -ForegroundColor Red
    exit 1
}

# 3. Process the session file
$sessionPath = "session/session.json"
if (Test-Path $sessionPath) {
    Write-Host "[+] Processing session file..." -ForegroundColor Green
    
    try {
        $bytes = [IO.File]::ReadAllBytes($sessionPath)
        $b64 = [Convert]::ToBase64String($bytes)
        
        # Determine strict content for .env
        $envContent = Get-Content ".env"
        $newContent = @()
        $found = $false
        
        foreach ($line in $envContent) {
            # Trim whitespace
            $trimmed = $line.Trim()
            
            if ($trimmed -match "^SESSION_STRING=") {
                $newContent += "SESSION_STRING=$b64"
                $found = $true
            } elseif ($trimmed -ne "") {
                $newContent += $line
            }
        }
        
        if (-not $found) {
            $newContent += "SESSION_STRING=$b64"
        }
        
        $newContent | Set-Content ".env"
        Write-Host "[+] SUCCESSFULLY updated .env with SESSION_STRING!" -ForegroundColor Green
        Write-Host "[*] Setup complete. You can now use 'task up' to start the bot." -ForegroundColor Cyan
        
        # Cleanup
        if (Test-Path "session") {
            Remove-Item -Recurse -Force "session"
            Write-Host "[*] Temporary session files cleaned up."
        }
        
    } catch {
        Write-Host "[!] Error processing session: $($_.Exception.Message)" -ForegroundColor Red
    }
} else {
    Write-Host "[!] Session file not found. Something went wrong." -ForegroundColor Red
}
