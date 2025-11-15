<template>
  <div class="matrix-access-sso">
    <button 
      @click="openMatrix" 
      class="matrix-button"
      :disabled="isLoadingMatrix || !isConnected"
    >
      <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
      </svg>
      <span v-if="!isLoadingMatrix">Ouvrir Element (Chat)</span>
      <span v-else>Connexion en cours...</span>
    </button>

    <p v-if="error" class="error-message">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { gls } from '@/stores/global'

const store = gls()
const isLoadingMatrix = ref(false)
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

// Ouvrir Element avec SSO automatique (sans demander le mot de passe)
const openMatrix = async () => {
  if (!isConnected.value) {
    error.value = 'Vous devez être connecté pour accéder à Matrix'
    return
  }

  isLoadingMatrix.value = true
  error.value = ''

  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/matrix/sso`, {
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
      // Ouvrir Element avec SSO automatique
      window.open(data.url, '_blank', 'noopener,noreferrer')
    } else {
      const data = await response.json()
      error.value = data.error || 'Erreur lors de l\'accès à Matrix'
    }
  } catch (err) {
    console.error('Erreur:', err)
    error.value = 'Impossible de se connecter au serveur'
  } finally {
    isLoadingMatrix.value = false
  }
}
</script>

<style scoped>
.matrix-access-sso {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.matrix-button {
  font-family: roboto, sans-serif;
  font-size: 1.2rem;
  padding: 1rem 2.5rem;
  border-radius: 32px;
  border: none;
  cursor: pointer;
  background: linear-gradient(135deg, #00b4d8 0%, #0077b6 100%);
  color: white;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  font-weight: 600;
}

.matrix-button:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.matrix-button:disabled {
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
