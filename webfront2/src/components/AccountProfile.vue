<script setup>
import { ref } from "vue"
import { gls } from "@/stores/global.js"

let dj = ref(0)
let lj = ref(0)
let domain = ref(0)
let nboffer = ref(0)
let phone  = ref(0)
let email = ref(0)
let loading = ref(true)
let regeneratingMatrix = ref(false)
let matrixMessage = ref('')
let showMatrixModal = ref(false)
let matrixPassword = ref('')

fetch(import.meta.env.VITE_APP_API_GETINFOP, {
  method: "POST",
  mode: "cors",
  body: JSON.stringify({
    "username": gls().username,
    "token": gls().sessionT
  })
})
.then(res => res.json())
.then(data => {
  dj.value = data['DateJoined']
  domain.value = data['Domain']
  nboffer.value = data['Nboffer']
  email.value = data['Email']
  phone.value = data['PhoneNumber']
  lj.value = data["LastLogin"]
  loading.value = false
})
.catch(() => {
  loading.value = false
})

const offerName = (num) => {
  const offers = { 0: 'Free', 1: 'Standard', 2: 'Professional', 3: 'Enterprise' }
  return offers[num] || 'Unknown'
}

const storageSize = (num) => {
  const storage = { 0: '1GB', 1: '50GB', 2: '200GB', 3: 'Illimité' }
  return storage[num] || '1GB'
}

const regenerateMatrix = async () => {
  showMatrixModal.value = true
}

const confirmRegenerateMatrix = async () => {
  if (!matrixPassword.value) {
    matrixMessage.value = 'Veuillez entrer votre mot de passe'
    return
  }
  regeneratingMatrix.value = true
  matrixMessage.value = ''
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/matrix/regenerate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: gls().username,
        token: gls().sessionT,
        password: matrixPassword.value
      })
    })
    const data = await response.json()
    if (response.ok) {
      matrixMessage.value = 'Compte Matrix régénéré avec succès !'
      setTimeout(() => {
        matrixMessage.value = ''
        showMatrixModal.value = false
        matrixPassword.value = ''
      }, 3000)
    } else {
      matrixMessage.value = data.message || 'Erreur lors de la régénération'
    }
  } catch (error) {
    matrixMessage.value = 'Erreur de connexion au serveur'
  } finally {
    regeneratingMatrix.value = false
  }
}

const cancelRegenerateMatrix = () => {
  showMatrixModal.value = false
  matrixPassword.value = ''
  matrixMessage.value = ''
}
</script>

<template>
  <div class="profile">
    <!-- Loading -->
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <span>Chargement du profil...</span>
    </div>

    <!-- Profile content -->
    <div v-else class="content">
      <!-- Header card -->
      <div class="header-card">
        <div class="avatar">{{ gls().username ? gls().username.charAt(0).toUpperCase() : '?' }}</div>
        <div class="user-info">
          <h1>{{ gls().username }}</h1>
          <span class="badge">{{ offerName(nboffer) }}</span>
        </div>
        <RouterLink to="/account/edit" class="edit-link">Modifier</RouterLink>
      </div>

      <!-- Info grid -->
      <div class="info-grid">
        <div class="info-card">
          <span class="info-label">Email</span>
          <span class="info-value">{{ email || '—' }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">Téléphone</span>
          <span class="info-value">{{ phone || '—' }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">Domaine</span>
          <span class="info-value">{{ domain || '—' }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">Stockage</span>
          <span class="info-value">{{ storageSize(nboffer) }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">Inscrit le</span>
          <span class="info-value">{{ dj || '—' }}</span>
        </div>
        <div class="info-card">
          <span class="info-label">Dernière connexion</span>
          <span class="info-value">{{ lj || '—' }}</span>
        </div>
      </div>

      <!-- Matrix section -->
      <div class="matrix-section">
        <h3>Matrix</h3>
        <p>Régénérez votre compte Matrix si nécessaire.</p>
        <button @click="regenerateMatrix" class="btn-matrix">Régénérer le compte Matrix</button>
      </div>
    </div>

    <!-- Matrix modal -->
    <div v-if="showMatrixModal" class="modal-overlay" @click.self="cancelRegenerateMatrix">
      <div class="modal-card">
        <h3>Régénérer Matrix</h3>
        <p>Confirmez avec votre mot de passe.</p>
        <input v-model="matrixPassword" type="password" placeholder="Mot de passe" />
        <p v-if="matrixMessage" class="modal-msg">{{ matrixMessage }}</p>
        <div class="modal-actions">
          <button @click="cancelRegenerateMatrix" class="btn-cancel">Annuler</button>
          <button @click="confirmRegenerateMatrix" class="btn-confirm" :disabled="regeneratingMatrix">
            <span v-if="regeneratingMatrix" class="spinner"></span>
            {{ regeneratingMatrix ? '...' : 'Confirmer' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile { max-width: 640px; margin: 0 auto; padding: 2.5rem 1.25rem 5rem; animation: fadeInUp .5s var(--ease) both; }
.loading { display: flex; align-items: center; justify-content: center; gap: .75rem; padding: 4rem 2rem; color: var(--text-muted); font-size: .9375rem; }
.spinner { width: 24px; height: 24px; border: 3px solid var(--border); border-top-color: var(--ink); border-radius: 50%; animation: spin .7s linear infinite; }

/* Header card */
.header-card { display: flex; align-items: center; gap: 1.25rem; padding: 1.75rem; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); margin-bottom: 1.5rem; animation: fadeInUp .4s var(--ease) both; }
.avatar { width: 56px; height: 56px; border-radius: 50%; background: var(--ink); color: var(--paper); display: flex; align-items: center; justify-content: center; font-size: 1.5rem; font-weight: 800; flex-shrink: 0; transition: transform var(--tb); }
.avatar:hover { transform: scale(1.05); }
.user-info { flex: 1; }
.user-info h1 { font-size: 1.375rem; font-weight: 800; letter-spacing: -.02em; margin-bottom: .25rem; }
.badge { display: inline-block; padding: .2rem .65rem; font-size: .75rem; font-weight: 600; border-radius: var(--r-full); background: var(--bg-tertiary); color: var(--text-secondary); }
.edit-link { padding: .5rem 1.125rem; font-size: .8125rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); transition: all var(--tb); }
.edit-link:hover { border-color: var(--ink); background: var(--bg-tertiary); transform: translateY(-1px); }

/* Info grid */
.info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: .875rem; margin-bottom: 1.5rem; }
.info-card { padding: 1rem 1.25rem; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r); display: flex; flex-direction: column; gap: .25rem; animation: fadeInUp .4s var(--ease) both; transition: border-color var(--tb); }
.info-card:nth-child(2) { animation-delay: 40ms; }
.info-card:nth-child(3) { animation-delay: 80ms; }
.info-card:nth-child(4) { animation-delay: 120ms; }
.info-card:nth-child(5) { animation-delay: 160ms; }
.info-card:nth-child(6) { animation-delay: 200ms; }
.info-card:hover { border-color: var(--border-strong); }
.info-label { font-size: .75rem; font-weight: 600; color: var(--text-muted); }
.info-value { font-size: .9375rem; font-weight: 500; color: var(--text-primary); }

/* Matrix */
.matrix-section { padding: 1.75rem; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); animation: fadeInUp .4s var(--ease) .2s both; }
.matrix-section h3 { font-size: 1.125rem; font-weight: 700; margin-bottom: .375rem; }
.matrix-section p { font-size: .875rem; color: var(--text-secondary); margin-bottom: 1rem; }
.btn-matrix { padding: .625rem 1.25rem; font-size: .875rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-matrix:hover { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-sm); transform: translateY(-1px); }

/* Modal */
.modal-overlay { position: fixed; inset: 0; z-index: 999; background: var(--bg-overlay); backdrop-filter: blur(6px); display: flex; align-items: center; justify-content: center; animation: fadeIn .2s var(--ease) both; }
.modal-card { width: calc(100% - 32px); max-width: 380px; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); box-shadow: var(--sh-xl); padding: 1.75rem; animation: popIn .25s var(--ease-spring) both; }
.modal-card h3 { font-size: 1.125rem; font-weight: 700; margin-bottom: .375rem; }
.modal-card p { font-size: .875rem; color: var(--text-secondary); margin-bottom: 1rem; }
.modal-card input { width: 100%; padding: .625rem .875rem; font-size: .9375rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); margin-bottom: .75rem; }
.modal-card input:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.modal-msg { font-size: .8125rem; color: var(--text-primary); margin-bottom: .75rem; padding: .5rem .75rem; background: var(--bg-tertiary); border: 1px solid var(--border); border-radius: var(--r); animation: fadeIn .2s var(--ease) both; }
.modal-actions { display: flex; gap: .625rem; }
.btn-cancel { flex: 1; padding: .625rem; font-size: .875rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-cancel:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.btn-confirm { flex: 1; padding: .625rem; font-size: .875rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); display: flex; align-items: center; justify-content: center; gap: .375rem; }
.btn-confirm:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); }
.btn-confirm:disabled { opacity: .5; cursor: not-allowed; }

@media (max-width: 480px) { .info-grid { grid-template-columns: 1fr; } }
</style>
