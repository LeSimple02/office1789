param(
    [string]$CN = "localhost",
    [int]$Days = 365,
    [switch]$Force
)

$ErrorActionPreference = 'Stop'

# Resolve SSL output directory
$sslDir = Join-Path $PSScriptRoot 'config\ssl'
if (-not (Test-Path $sslDir)) {
    New-Item -ItemType Directory -Path $sslDir | Out-Null
}

$certPath = Join-Path $sslDir 'cert.pem'
$keyPath  = Join-Path $sslDir 'key.pem'

if ((Test-Path $certPath) -and (Test-Path $keyPath) -and -not $Force) {
    Write-Host "Existing certs found at $sslDir. Use -Force to regenerate." -ForegroundColor Yellow
    return
}

# Ensure Docker is available
if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    throw "Docker CLI is not available in PATH. Please install Docker Desktop."
}

# Windows path to bind mount
$mountPath = $sslDir
Write-Host "Generating self-signed cert for CN=$CN into $mountPath ..." -ForegroundColor Cyan

# Use a temporary Alpine container to install openssl and generate certs
$linuxCmd = @(
    'apk add --no-cache openssl >/dev/null 2>&1',
    "openssl req -x509 -newkey rsa:4096 -keyout /out/key.pem -out /out/cert.pem -sha256 -days $Days -nodes -subj '/CN=$CN'"
) -join ' && '

$proc = Start-Process -FilePath 'docker' -ArgumentList @('run','--rm','-v',"$mountPath:/out",'alpine:3','sh','-lc', $linuxCmd) -NoNewWindow -PassThru -Wait
if ($proc.ExitCode -ne 0) {
    throw "Failed to generate certs using Docker (exit code $($proc.ExitCode))."
}

# Basic sanity check
if (-not (Test-Path $certPath)) { throw "Missing $certPath" }
if (-not (Test-Path $keyPath)) { throw "Missing $keyPath" }

Write-Host "Self-signed certs generated:" -ForegroundColor Green
Write-Host "  $certPath" -ForegroundColor Green
Write-Host "  $keyPath" -ForegroundColor Green

Write-Host "You can now (re)start the stack: `ndocker compose -f $($PSScriptRoot)\docker-compose.yml up -d`n" -ForegroundColor Gray
