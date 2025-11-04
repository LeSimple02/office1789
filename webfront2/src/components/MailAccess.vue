<template>
  <div class="mail-access">
    <button 
      v-if="!isLoadingMail" 
      @click="openMail" 
      class="mail-button"
      :disabled="!isConnected"
    >
      <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
        <polyline points="22,6 12,13 2,6"/>
      </svg>
      Ouvrir ma boîte mail
    </button>
    
    <div v-if="isLoadingMail" class="loading">
      Chargement de votre boîte mail...
    </div>
    
    <div v-if="error" class="error">
      {{ error }}
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { gls } from '@/stores/global'

export default {
  name: 'MailAccess',
  setup() {
    const store = gls()
    const isLoadingMail = ref(false)
    const isConnected = ref(false)
    const error = ref('')

    // Vérifier si l'utilisateur est connecté
    const checkConnection = async () => {
      try {
        const response = await fetch('http://localhost:8080/api/session/check', {
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
          error.value = 'Votre session a expiré, veuillez vous reconnecter'
        }
      } catch (err) {
        console.error('Erreur de vérification de session:', err)
        isConnected.value = false
      }
    }

    // Ouvrir la boîte mail avec SSO
    const openMail = async () => {
      if (!isConnected.value) {
        error.value = 'Vous devez être connecté pour accéder à votre boîte mail'
        return
      }

      isLoadingMail.value = true
      error.value = ''

      try {
        const response = await fetch('http://localhost:8080/api/mail/auth', {
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
          // Ouvrir Roundcube dans un nouvel onglet avec authentification automatique
          window.open(data.url, '_blank', 'noopener,noreferrer')
        } else {
          error.value = 'Erreur lors de l\'accès à la boîte mail'
        }
      } catch (err) {
        console.error('Erreur:', err)
        error.value = 'Impossible de se connecter au serveur mail'
      } finally {
        isLoadingMail.value = false
      }
    }

    onMounted(() => {
      checkConnection()
      // Vérifier la connexion toutes les 5 minutes
      setInterval(checkConnection, 5 * 60 * 1000)
    })

    return {
      isLoadingMail,
      isConnected,
      error,
      openMail
    }
  }
}
</script>

<style scoped>
.mail-access {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.mail-button {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.mail-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.mail-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: #ccc;
}

.loading {
  padding: 10px;
  color: #667eea;
  font-style: italic;
}

.error {
  padding: 10px;
  background: #fee;
  color: #c33;
  border-radius: 4px;
  border-left: 4px solid #c33;
}
</style>
