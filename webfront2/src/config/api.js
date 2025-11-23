export const API_CONFIG = {
  // Configuration multi-environnement
  getBaseURL() {
    // 1. Vérifier les variables d'environnement Vite en priorité
    if (import.meta.env.VITE_APP_API) {
      console.log('[API] Using VITE_APP_API:', import.meta.env.VITE_APP_API)
      return import.meta.env.VITE_APP_API
    }
    
    // 2. Détection de l'environnement
    const protocol = window.location.protocol
    const hostname = window.location.hostname
    
    console.log('[API] No VITE_APP_API, detecting from hostname:', hostname)
    
    // Cordova (file://)
    if (protocol === 'file:' || window.cordova) {
      return localStorage.getItem('api_url') || localStorage.getItem('auto_api_url') || 'http://192.168.56.1:8080'
    }
    
    // Electron (app://)
    if (window.electron) {
      return localStorage.getItem('api_url') || 'http://localhost:8080'
    }
    
    // Web development local
    if (hostname === 'localhost' || hostname === '127.0.0.1') {
      return 'http://localhost:8080'
    }
    
    // Production - si on est sur office1789.com
    if (hostname.includes('office1789.com')) {
      // En production on ignore toute valeur locale forcée dans localStorage
      const forced = localStorage.getItem('api_url') || ''
      if (forced.startsWith('http://localhost')) {
        localStorage.removeItem('api_url')
        localStorage.removeItem('auto_api_url')
        console.log('[API] Cleared localhost api_url in production')
      }
      return 'https://backend.office1789.com'
    }
    
    // Fallback production
    return 'https://backend.office1789.com'
  },
  
  // Permet de changer l'URL API à la volée
  setBaseURL(url) {
    localStorage.setItem('api_url', url)
    window.location.reload()
  },
  
  // Présets d'URLs
  presets: {
    local: 'http://localhost:8080',
    wifi: 'http://192.168.56.1:8080', // sera mis à jour par setup-ip.ps1
    ngrok: 'https://your-ngrok-url.ngrok.io',
    production: 'https://backend.office1789.com'
  }
}

// Fonction utilitaire pour construire une URL complète
export function getApiUrl(endpoint) {
  const base = API_CONFIG.getBaseURL()
  // S'assurer qu'on n'a pas de double slash
  const cleanEndpoint = endpoint.startsWith('/') ? endpoint : '/' + endpoint
  console.log('[API] URL:', base + cleanEndpoint)
  return base + cleanEndpoint
}

// Fonction utilitaire pour tester la connexion
export async function testConnection(url = null) {
  const testUrl = url || API_CONFIG.getBaseURL()
  try {
    const response = await fetch(`${testUrl}/api/health`, { 
      method: 'GET',
      timeout: 5000 
    })
    return response.ok
  } catch (error) {
    console.error('Connection test failed:', error)
    return false
  }
}

// Auto-détection du backend (test des IP candidates)
export async function autoDetectBackend() {
  const candidates = []
  // IP déjà configurée
  const existing = localStorage.getItem('api_url')
  if (existing) candidates.push(existing)
  // Ajout des presets connus
  candidates.push(API_CONFIG.presets.wifi, API_CONFIG.presets.local)
  // Génération IP locale 192.168.x.1-50
  for (let i=1;i<=50;i++) { candidates.push(`http://192.168.0.${i}:8080`) }
  for (let i=1;i<=50;i++) { candidates.push(`http://192.168.1.${i}:8080`) }
  // Nettoyage doublons
  const unique = [...new Set(candidates)].filter(Boolean)
  console.log('[AutoDetect] Testing', unique.length, 'candidates')
  for (const base of unique) {
    try {
      const controller = new AbortController()
      const t = setTimeout(()=>controller.abort(), 1200)
      const resp = await fetch(base + '/api/health', {signal: controller.signal})
      clearTimeout(t)
      if (resp.ok) {
        console.log('[AutoDetect] Backend found:', base)
        localStorage.setItem('auto_api_url', base)
        return base
      }
    } catch(e) {
      // ignore
    }
  }
  console.warn('[AutoDetect] No backend detected')
  return null
}
