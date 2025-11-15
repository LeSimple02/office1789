<template>
  <div class="mail-access-sso">
    <button 
      @click="openMail" 
      class="mail-button"
      :disabled="isLoadingMail || !isConnected"
    >
      <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
      </svg>
      <span v-if="!isLoadingMail">Ouvrir ma boîte mail</span>
      <span v-else>Connexion en cours...</span>
    </button>

    <p v-if="error" class="error-message">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { gls } from '@/stores/global'

const store = gls()
const isLoadingMail = ref(false)
const isConnected = ref(true)
const error = ref('')

// Vérifier la session au chargement
onMounted(() => {
  checkConnection()
})

// Vérifier si l'utilisateur est connecté
const checkConnection = async () => {
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/session/check`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: store.username,
        token: store.sessionT
      })
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

// Ouvrir Roundcube avec SSO automatique (sans demander le mot de passe)
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
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: store.username,
        token: store.sessionT
      })
    })

    if (response.ok) {
      const data = await response.json()
      // Ouvrir Roundcube avec SSO automatique
      window.open(data.url, '_blank', 'noopener,noreferrer')
    } else {
      const data = await response.json()
      error.value = data.error || 'Erreur lors de l\'accès à la boîte mail'
    }
  } catch (err) {
    console.error('Erreur:', err)
    error.value = 'Impossible de se connecter au serveur'
  } finally {
    isLoadingMail.value = false
  }
}
</script>

<style scoped>
.mail-access-sso {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.mail-button {
  font-family: roboto, sans-serif;
  font-size: 1.2rem;
  padding: 1rem 2.5rem;
  border-radius: 32px;
  border: none;
  cursor: pointer;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  font-weight: 600;
}

.mail-button:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.mail-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-icon {
  width: 24px;
  height: 24px;
}

.error-message {
  color: #e74c3c;
  font-size: 0.9rem;
  text-align: center;
  padding: 12px 20px;
  background: rgba(231, 76, 60, 0.1);
  border-radius: 8px;
  border-left: 4px solid #e74c3c;
}
</style>
