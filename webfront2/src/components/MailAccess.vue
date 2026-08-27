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
    
    <div v-if="error && !showPasswordModal" class="error">
      {{ error }}
    </div>

    <!-- Modal de confirmation mot de passe -->
    <Teleport to="body">
      <div v-if="showPasswordModal" class="modal-overlay" @click="cancelOpenMail">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Accès à Roundcube</h2>
            <button @click="cancelOpenMail" class="close-btn">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <div class="modal-body">
            <p class="modal-description">
              Pour accéder à votre boîte mail, veuillez confirmer votre identité.
            </p>

            <div class="form-group">
              <label for="mail-password">Mot de passe</label>
              <input 
                id="mail-password"
                type="password" 
                v-model="password"
                @keyup.enter="confirmOpenMail"
                placeholder="Entrez votre mot de passe"
                :disabled="isLoadingMail"
                autofocus
              />
            </div>

            <div v-if="error" class="modal-message error">
              {{ error }}
            </div>
          </div>

          <div class="modal-footer">
            <button @click="cancelOpenMail" class="btn-cancel" :disabled="isLoadingMail">
              Annuler
            </button>
            <button @click="confirmOpenMail" class="btn-confirm" :disabled="isLoadingMail || !password">
              <span v-if="isLoadingMail">Connexion...</span>
              <span v-else>Ouvrir Roundcube</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>
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
    const showPasswordModal = ref(false)
    const password = ref('')

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
          error.value = 'Votre session a expiré, veuillez vous reconnecter'
        }
      } catch (err) {
        console.error('Erreur de vérification de session:', err)
        isConnected.value = false
      }
    }

    // Ouvrir le modal de mot de passe
    const openMail = () => {
      if (!isConnected.value) {
        error.value = 'Vous devez être connecté pour accéder à votre boîte mail'
        return
      }
      showPasswordModal.value = true
      error.value = ''
    }

    // Confirmer et ouvrir Roundcube
    const confirmOpenMail = async () => {
      if (!password.value) {
        error.value = 'Veuillez entrer votre mot de passe'
        return
      }

      isLoadingMail.value = true
      error.value = ''

      try {
        const response = await fetch(`${import.meta.env.VITE_APP_API}/api/mail/auth`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            username: store.username,
            token: store.sessionT,
            password: password.value
          })
        })

        if (response.ok) {
          const data = await response.json()
          // Preferir la configuration d'environnement en production
          const target = (import.meta.env.VITE_ROUNDCUBE_URL && import.meta.env.VITE_ROUNDCUBE_URL.trim()) || data.url
          // Ouvrir Roundcube dans un nouvel onglet avec authentification automatique
          window.open(target, '_blank', 'noopener,noreferrer')
          showPasswordModal.value = false
          password.value = ''
        } else {
          const data = await response.json()
          error.value = data.error || 'Erreur lors de l\'accès à la boîte mail'
        }
      } catch (err) {
        console.error('Erreur:', err)
        error.value = 'Impossible de se connecter au serveur mail'
      } finally {
        isLoadingMail.value = false
      }
    }

    const cancelOpenMail = () => {
      showPasswordModal.value = false
      password.value = ''
      error.value = ''
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
      showPasswordModal,
      password,
      openMail,
      confirmOpenMail,
      cancelOpenMail
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

/* Modal styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.modal-content {
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 90%;
  max-width: 480px;
  animation: slideUp 0.3s ease;
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 28px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}

.close-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  transition: background 0.2s ease;
  color: #6b7280;
}

.close-btn:hover {
  background: #f3f4f6;
}

.modal-body {
  padding: 24px 28px;
}

.modal-description {
  font-size: 0.95rem;
  color: #6b7280;
  margin-bottom: 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 0.9rem;
  font-weight: 600;
  color: #374151;
  margin-bottom: 8px;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  font-size: 1rem;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  background: #f9fafb;
  transition: all 0.2s ease;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
  background: #fff;
}

.form-group input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.modal-message {
  padding: 12px 16px;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 600;
  margin-top: 16px;
}

.modal-message.error {
  background: #fee2e2;
  color: #991b1b;
  border: 1px solid #fca5a5;
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 20px 28px;
  border-top: 1px solid #e5e7eb;
}

.btn-cancel {
  flex: 1;
  padding: 12px 24px;
  font-size: 1rem;
  font-weight: 600;
  border: 2px solid #e5e7eb;
  background: transparent;
  color: #6b7280;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-cancel:hover:not(:disabled) {
  background: #f3f4f6;
}

.btn-cancel:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-confirm {
  flex: 1;
  padding: 12px 24px;
  font-size: 1rem;
  font-weight: 600;
  border: none;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-confirm:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

.btn-confirm:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
