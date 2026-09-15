<script setup>
import { ref, onMounted } from 'vue'
import { gls } from '@/stores/global'

const store = gls()
const isLoadingMail = ref(false)
const isConnected = ref(true)
const error = ref('')

onMounted(() => { checkConnection() })

const checkConnection = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/session/check`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: store.username, token: store.sessionT })
    })
    const data = await response.json()
    isConnected.value = data.connected
    if (!data.connected) {
      error.value = 'Votre session a expiré, veuillez vous reconnecter à Office1789'
    }
  } catch (err) {
    console.error('Erreur de vérification de session:', err)
    isConnected.value = false
    error.value = 'Erreur de connexion au serveur'
  }
}

const openMail = async () => {
  if (!isConnected.value) {
    error.value = 'Vous devez être connecté pour accéder à votre boîte mail'
    return
  }
  isLoadingMail.value = true
  error.value = ''
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/mail/sso`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ username: store.username, token: store.sessionT })
    })
    if (response.ok) {
      const data = await response.json()
      const isElectron = window.navigator.userAgent.toLowerCase().includes('electron')
      if (isElectron && window.electronAPI?.openMailWindow) {
        window.electronAPI.openMailWindow(data.url)
      } else {
        window.open(data.url, '_blank', 'noopener,noreferrer')
      }
    } else {
      const data = await response.json()
      error.value = data.error || "Erreur lors de l'accès à la boîte mail"
    }
  } catch (err) {
    console.error('[Mail SSO] Erreur réseau:', err)
    error.value = 'Impossible de se connecter au serveur'
  } finally {
    isLoadingMail.value = false
  }
}
</script>

<template>
  <div class="sso">
    <button @click="openMail" class="sso-btn" :disabled="isLoadingMail">
      <span v-if="isLoadingMail" class="spinner"></span>
      <span v-if="!isLoadingMail">Ouvrir la messagerie</span>
      <span v-else>Connexion...</span>
    </button>
    <p v-if="error" class="sso-err">{{ error }}</p>
  </div>
</template>

<style scoped>
.sso { display: flex; flex-direction: column; align-items: center; gap: .75rem; animation: fadeInUp .4s var(--ease) .2s both; }
.sso-btn {
  padding: .75rem 2rem; font-size: .9375rem; font-weight: 600;
  color: var(--paper); background: var(--ink); border: 1px solid var(--ink);
  border-radius: var(--r); cursor: pointer; transition: all var(--tb);
  display: flex; align-items: center; gap: .5rem; min-width: 200px; justify-content: center;
}
.sso-btn:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-md); transform: translateY(-2px); }
.sso-btn:disabled { opacity: .5; cursor: not-allowed; }
.spinner { width: 18px; height: 18px; border: 2.5px solid var(--border); border-top-color: var(--paper); border-radius: 50%; animation: spin .7s linear infinite; }
.sso-err { font-size: .8125rem; color: var(--text-primary); padding: .5rem .875rem; background: var(--bg-tertiary); border: 1px solid var(--border); border-radius: var(--r); animation: fadeIn .2s var(--ease) both; }
</style>
