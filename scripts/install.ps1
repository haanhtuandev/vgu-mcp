#Requires -Version 5.1
<#
.SYNOPSIS
    Installs vgu-mcp (VGU Moodle MCP server) on Windows.

.DESCRIPTION
    Downloads the latest release ZIP, verifies its SHA-256 checksum,
    extracts vgu-mcp.exe to %LOCALAPPDATA%\Programs\vgu-mcp, and adds that
    directory to your user PATH. No admin rights required.

.EXAMPLE
    irm https://raw.githubusercontent.com/haanhtuandev/vgu-mcp/main/scripts/install.ps1 | iex

    To install a specific version:
    $env:VGU_MCP_VERSION = '0.1.2'; irm .../install.ps1 | iex

.NOTES
    Re-running updates to the latest release. To uninstall, delete
    %LOCALAPPDATA%\Programs\vgu-mcp and remove it from your user PATH.
#>

$ErrorActionPreference = 'Stop'
$ProgressPreference    = 'SilentlyContinue'
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$Repo = if ($env:VGU_MCP_REPO) { $env:VGU_MCP_REPO } else { 'haanhtuandev/vgu-mcp' }
$Base = "https://github.com/$Repo/releases/download"
$LocalAppData = if ($env:LOCALAPPDATA) { $env:LOCALAPPDATA } else { Join-Path $env:USERPROFILE 'AppData\Local' }

function Write-Step([string]$Message) { Write-Host "==> $Message" -ForegroundColor Green }

if ($env:VGU_MCP_VERSION) {
    $Tag = if ($env:VGU_MCP_VERSION -like 'v*') { $env:VGU_MCP_VERSION } else { "v$($env:VGU_MCP_VERSION)" }
} else {
    Write-Step "Resolving latest release for $Repo ..."
    $Release = Invoke-RestMethod -Headers @{ 'User-Agent' = 'vgu-mcp-installer' } `
        -Uri "https://api.github.com/repos/$Repo/releases/latest"
    $Tag = $Release.tag_name
}
$Version = $Tag -replace '^v', ''
$Asset    = "vgu-mcp_${Version}_windows_amd64.zip"
$Tmp      = Join-Path $env:TEMP 'vgu-mcp-installer'
New-Item -ItemType Directory -Force -Path $Tmp | Out-Null

Write-Step "Downloading $Asset ..."
$ZipPath = Join-Path $Tmp $Asset
Invoke-WebRequest -UseBasicParsing -Uri "$Base/$Tag/$Asset" -OutFile $ZipPath

Write-Step 'Verifying SHA-256 checksum ...'
$ChecksumsPath = Join-Path $Tmp 'checksums.txt'
Invoke-WebRequest -UseBasicParsing -Uri "$Base/$Tag/checksums.txt" -OutFile $ChecksumsPath
$Expected = Get-Content -LiteralPath $ChecksumsPath |
    Where-Object { $_ -match "  $([regex]::Escape($Asset))$" } |
    ForEach-Object { ($_ -split '\s+')[0] }
if (-not $Expected) { throw "No checksum found for $Asset in checksums.txt" }
$Actual = (Get-FileHash -Algorithm SHA256 -LiteralPath $ZipPath).Hash.ToLowerInvariant()
if ($Actual -ne $Expected) { throw "Checksum mismatch for $Asset (expected $Expected, got $Actual)" }

$InstallRoot = Join-Path $LocalAppData 'Programs\vgu-mcp'
$Exe = Join-Path $InstallRoot 'vgu-mcp.exe'

Write-Step "Installing to $InstallRoot ..."
if (Test-Path -LiteralPath $InstallRoot) {
    Remove-Item -LiteralPath $InstallRoot -Recurse -Force
}
Expand-Archive -LiteralPath $ZipPath -DestinationPath $InstallRoot -Force
if (-not (Test-Path -LiteralPath $Exe)) { throw "vgu-mcp.exe was not found in the downloaded archive" }
# Clear the Mark-of-the-Web so Windows SmartScreen does not block first launch.
Unblock-File -LiteralPath $Exe -ErrorAction SilentlyContinue

Write-Step "Adding $InstallRoot to your user PATH ..."
$UserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($UserPath -notlike "*$InstallRoot*") {
    $NewPath = if ([string]::IsNullOrEmpty($UserPath)) { $InstallRoot } else { "$InstallRoot;$UserPath" }
    [Environment]::SetEnvironmentVariable('Path', $NewPath, 'User')
}

Write-Host ''
Write-Host "vgu-mcp $Version installed to $InstallRoot" -ForegroundColor Green
Write-Host ''
Write-Host 'Next steps:'
Write-Host '  1. Open a NEW terminal so PATH picks up vgu-mcp.'
Write-Host '  2. Authenticate once:'
Write-Host '       vgu-mcp setup'
Write-Host '  3. Add it to your AI client (see the project README).'
Write-Host ''
Write-Host 'To upgrade later, just re-run this command.'
Write-Host "To uninstall: remove $InstallRoot and drop it from your user PATH."
