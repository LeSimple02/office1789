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
          error.value = 'Votre session a expiré, veuillez vous reconnecter'
        }
      } catch (err) {
        console.error('Erreur de vérification de session:', err)
        isConnected.value = false
      }
    }

    const openMail = () => {
      if (!isConnected.value) {
        error.value = 'Vous devez être connecté pour accéder à votre boîte mail'
        return
      }
      showPasswordModal.value = true
      error.value = ''
    }

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
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username: store.username, token: store.sessionT, password: password.value })
        })
        if (response.ok) {
          const data = await response.json()
          const target = (import.meta.env.VITE_ROUNDCUBE_URL && import.meta.env.VITE_ROUNDCUBE_URL.trim()) || data.url
          window.open(target, '_blank', 'noopener,noreferrer')
          showPasswordModal.value = false
          password.value = ''
        } else {
          const data = await response.json()
          error.value = data.error || "Erreur lors de l'accès à la boîte mail"
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
      setInterval(checkConnection, 5 * 60 * 1000)
    })

    return { isLoadingMail, isConnected, error, showPasswordModal, password, openMail, confirmOpenMail, cancelOpenMail }
  }
}
</script>

<template>
  <div class="mail-access">
    <button @click="openMail" class="open-btn">Accéder à la messagerie</button>
    <p v-if="error" class="err-msg">{{ error }}</p>

    <div v-if="showPasswordModal" class="modal-overlay" @click.self="cancelOpenMail">
      <div class="modal-card">
        <h3>Mot de passe mail</h3>
        <p>Entrez votre mot de passe pour accéder à Roundcube.</p>
        <input v-model="password" type="password" placeholder="Mot de passe" @keyup.enter="confirmOpenMail" />
        <div class="modal-actions">
          <button @click="cancelOpenMail" class="btn-cancel">Annuler</button>
          <button @click="confirmOpenMail" class="btn-confirm" :disabled="isLoadingMail">
            <span v-if="isLoadingMail" class="spinner"></span>
            {{ isLoadingMail ? '...' : 'Confirmer' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.mail-access { display: flex; flex-direction: column; align-items: center; gap: .75rem; animation: fadeInUp .4s var(--ease) both; }
.open-btn { padding: .75rem 2rem; font-size: .9375rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.open-btn:hover { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-md); transform: translateY(-2px); }
.err-msg { font-size: .8125rem; color: var(--text-primary); padding: .5rem .875rem; background: var(--bg-tertiary); border: 1px solid var(--border); border-radius: var(--r); animation: fadeIn .2s var(--ease) both; }
.modal-overlay { position: fixed; inset: 0; z-index: 999; background: var(--bg-overlay); backdrop-filter: blur(6px); display: flex; align-items: center; justify-content: center; animation: fadeIn .2s var(--ease) both; }
.modal-card { width: calc(100% - 32px); max-width: 380px; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); box-shadow: var(--sh-xl); padding: 1.75rem; animation: popIn .25s var(--ease-spring) both; }
.modal-card h3 { font-size: 1.125rem; font-weight: 700; margin-bottom: .375rem; }
.modal-card p { font-size: .875rem; color: var(--text-secondary); margin-bottom: 1rem; }
.modal-card input { width: 100%; padding: .625rem .875rem; font-size: .9375rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); margin-bottom: 1rem; }
.modal-card input:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.modal-actions { display: flex; gap: .625rem; }
.btn-cancel { flex: 1; padding: .625rem; font-size: .875rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-cancel:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.btn-confirm { flex: 1; padding: .625rem; font-size: .875rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); display: flex; align-items: center; justify-content: center; gap: .375rem; }
.btn-confirm:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); }
.btn-confirm:disabled { opacity: .5; cursor: not-allowed; }
.spinner { width: 16px; height: 16px; border: 2px solid var(--border); border-top-color: var(--paper); border-radius: 50%; animation: spin .7s linear infinite; }
</style>
