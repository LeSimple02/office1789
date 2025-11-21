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

fetch(import.meta.env.VITE_APP_API_INFO_USER, {
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
  <div class="profile-container">
    <header class="profile-header">
      <h1 class="title">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
          <circle cx="12" cy="7" r="4"></circle>
        </svg>
        {{ $t('infop') }}
      </h1>
      <RouterLink to="/account/edit" class="edit-btn">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
        </svg>
        {{ $t('edit') }}
      </RouterLink>
      <button @click="regenerateMatrix" :disabled="regeneratingMatrix" class="matrix-btn">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"/>
        </svg>
        {{ regeneratingMatrix ? 'Régénération...' : 'Régénérer Matrix' }}
      </button>
      <div v-if="matrixMessage && !showMatrixModal" class="matrix-message" :class="{ success: matrixMessage.includes('succès') }">
        {{ matrixMessage }}
      </div>
    </header>

    <!-- Modal de régénération Matrix -->
    <Teleport to="body">
      <div v-if="showMatrixModal" class="modal-overlay" @click="cancelRegenerateMatrix">
        <div class="modal-content" @click.stop>
          <div class="modal-header">
            <h2>Régénérer le compte Matrix</h2>
            <button @click="cancelRegenerateMatrix" class="close-btn">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>

          <div class="modal-body">
            <p class="modal-description">
              Pour régénérer votre compte Matrix, veuillez confirmer votre identité en entrant votre mot de passe.
            </p>

            <div class="form-group">
              <label for="matrix-password">Mot de passe</label>
              <input 
                id="matrix-password"
                type="password" 
                v-model="matrixPassword"
                @keyup.enter="confirmRegenerateMatrix"
                placeholder="Entrez votre mot de passe"
                :disabled="regeneratingMatrix"
              />
            </div>

            <div v-if="matrixMessage" class="modal-message" :class="{ success: matrixMessage.includes('succès'), error: !matrixMessage.includes('succès') }">
              {{ matrixMessage }}
            </div>
          </div>

          <div class="modal-footer">
            <button @click="cancelRegenerateMatrix" class="btn-cancel" :disabled="regeneratingMatrix">
              Annuler
            </button>
            <button @click="confirmRegenerateMatrix" class="btn-confirm" :disabled="regeneratingMatrix || !matrixPassword">
              <span v-if="regeneratingMatrix">Régénération...</span>
              <span v-else>Confirmer</span>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>

    <section v-else class="profile-card">
      <div class="profile-banner">
        <div class="avatar-wrapper">
          <img src="@/assets/napo.png" alt="profile" class="avatar" />
          <div class="avatar-badge">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
              <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
            </svg>
          </div>
        </div>
        <div class="profile-info">
          <h2 class="username">{{ gls().username }}</h2>
          <p class="domain-tag">{{ domain }}</p>
        </div>
      </div>

      <div class="info-grid">
        <div class="info-card">
          <div class="info-icon email">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
              <polyline points="22,6 12,13 2,6"></polyline>
            </svg>
          </div>
          <div class="info-content">
            <span class="info-label">{{ $t('emaily') }}</span>
            <span class="info-value">{{ email }}</span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon phone">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
            </svg>
          </div>
          <div class="info-content">
            <span class="info-label">{{ $t('phoney') }}</span>
            <span class="info-value">{{ phone || '—' }}</span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon offer">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <polygon points="10,8 16,12 10,16"></polygon>
            </svg>
          </div>
          <div class="info-content">
            <span class="info-label">{{ $t('offery') }}</span>
            <span class="info-value">{{ offerName(nboffer) }}</span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon storage">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
            </svg>
          </div>
          <div class="info-content">
            <span class="info-label">Stockage</span>
            <span class="info-value">{{ storageSize(nboffer) }}</span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon security">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
            </svg>
          </div>
          <div class="info-content">
            <span class="info-label">{{ $t('password') }}</span>
            <span class="info-value">●●●●●●●●</span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon calendar">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="16" y1="2" x2="16" y2="6"></line>
              <line x1="8" y1="2" x2="8" y2="6"></line>
              <line x1="3" y1="10" x2="21" y2="10"></line>
            </svg>
          </div>
          <div class="info-content">
            <span class="info-label">{{ $t('datej') }}</span>
            <span class="info-value">{{ new Date(dj).toLocaleDateString() }}</span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-icon clock">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <polyline points="12,6 12,12 16,14"></polyline>
            </svg>
          </div>
          <div class="info-content">
            <span class="info-label">{{ $t('lastj') }}</span>
            <span class="info-value">{{ new Date(lj).toLocaleDateString() }}</span>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
* { box-sizing: border-box; font-family: 'Roboto', sans-serif; }

.profile-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 0;
  width: 100%;
  animation: fadeIn 0.5s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.profile-header {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  padding: 48px 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 32px;
  gap: 24px;
  position: relative;
}

.dark .profile-container { background: transparent; }
.dark .profile-header { background: #1C1C1E; }

.title {
  font-size: 3rem;
  font-weight: 700;
  color: #222;
  display: flex;
  align-items: center;
  gap: 12px;
  letter-spacing: 2px;
  margin: 0;
}
.dark .title { color: white; }

.title svg {
  width: 48px;
  height: 48px;
  padding: 12px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 50%;
  color: white;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.edit-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 32px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 24px;
  color: #fff;
  text-decoration: none;
  font-weight: 600;
  font-size: 1.1rem;
  transition: all 0.3s ease;
  border: none;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  cursor: pointer;
}
.edit-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.matrix-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 14px 32px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 24px;
  color: #fff;
  text-decoration: none;
  font-weight: 600;
  font-size: 1.1rem;
  transition: all 0.3s ease;
  border: none;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  cursor: pointer;
}
.matrix-btn:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}
.matrix-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.matrix-message {
  padding: 12px 24px;
  border-radius: 12px;
  font-size: 0.95rem;
  font-weight: 600;
  background: #fee;
  color: #c33;
  animation: slideDown 0.3s ease;
}
.matrix-message.success {
  background: #efe;
  color: #3c3;
}

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}
.spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}

.profile-card {
  width: 100%;
  background: transparent;
  border-radius: 0;
  padding: 0;
  box-shadow: none;
  overflow: visible;
}
.dark .profile-card {
  background: transparent;
  color: #eee;
  box-shadow: none;
}

.profile-banner {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  padding: 40px 32px;
  display: flex;
  align-items: center;
  gap: 32px;
  position: relative;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 32px;
}
.dark .profile-banner {
  background: #1C1C1E;
}

.avatar-wrapper {
  position: relative;
  flex-shrink: 0;
}
.avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  border: 5px solid rgba(0, 48, 143, 0.2);
  box-shadow: 0 8px 24px rgba(0,0,0,0.15);
  object-fit: cover;
}
.avatar-badge {
  position: absolute;
  bottom: 5px;
  right: 5px;
  width: 32px;
  height: 32px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 50%;
  border: 3px solid #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.profile-info {
  flex: 1;
}
.username {
  font-size: 2rem;
  font-weight: 700;
  color: #222;
  margin: 0 0 8px;
  letter-spacing: 1px;
}
.dark .username { color: white; }

.domain-tag {
  display: inline-block;
  padding: 6px 14px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 20px;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  margin: 0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
  padding: 0;
  margin-bottom: 32px;
}

.info-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 20px;
  transition: all 0.3s ease;
  border: 2px solid rgba(0, 48, 143, 0.1);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}
.dark .info-card {
  background: #1C1C1E;
  border-color: rgba(255, 255, 255, 0.1);
}
.info-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px rgba(0,0,0,0.12);
}
.dark .info-card:hover {
  box-shadow: 0 12px 24px rgba(0,0,0,0.4);
}

.info-icon {
  width: 50px;
  height: 50px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: -webkit-linear-gradient(30deg, blue, red);
  color: #fff;
}

.info-icon.storage {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.info-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.info-label {
  font-size: 0.85rem;
  color: #6c757d;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
}
.dark .info-label { color: #9ca3af; }
.info-value {
  font-size: 1.1rem;
  font-weight: 600;
  color: #212529;
}
.dark .info-value { color: #e5e7eb; }

.security-section {
  padding: 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}
.dark .security-section { 
  background: #1C1C1E;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 16px;
  color: #212529;
}
.dark .section-title { color: #e5e7eb; }

.security-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: #f8f9fa;
  border-radius: 12px;
}
.dark .security-item { background: #2a2d3a; }

.status-badge {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 600;
}
.status-badge.inactive {
  background: #fef3c7;
  color: #92400e;
}
.dark .status-badge.inactive {
  background: #78350f;
  color: #fef3c7;
}

.activity-section {
  padding: 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}
.dark .activity-section { 
  background: #1C1C1E;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}

.activity-section .section-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.activity-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
}

.activity-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #fff;
  border-radius: 16px;
  border: 1px solid rgba(0, 48, 143, 0.1);
  transition: all 0.3s ease;
}
.dark .activity-card {
  background: #2a2d3a;
  border-color: rgba(255, 255, 255, 0.1);
}

.activity-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.1);
}

.activity-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.activity-icon.success {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.activity-icon.info {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.activity-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
}

.activity-label {
  font-size: 0.85rem;
  color: #6c757d;
  font-weight: 600;
}
.dark .activity-label { color: #9ca3af; }

.activity-value {
  font-size: 1rem;
  font-weight: 600;
  color: #212529;
}
.dark .activity-value { color: #e5e7eb; }

.activity-link {
  font-size: 0.95rem;
  font-weight: 600;
  color: #667eea;
  text-decoration: none;
  transition: color 0.3s ease;
}
.activity-link:hover {
  color: #764ba2;
}

.btn-link {
  padding: 8px 16px;
  background: transparent;
  border: 1px solid #667eea;
  color: #667eea;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
}
.btn-link:hover {
  background: #667eea;
  color: #fff;
}

@media (max-width: 768px) {
  .title { font-size: 24px; }
  .profile-banner { flex-direction: column; text-align: center; padding: 32px 24px; }
  .username { font-size: 22px; }
  .avatar { width: 100px; height: 100px; }
  .info-grid { grid-template-columns: 1fr; padding: 24px; }
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
  animation: fadeIn 0.2s ease;
}

.modal-content {
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 90%;
  max-width: 480px;
  max-height: 90vh;
  overflow-y: auto;
  animation: slideUp 0.3s ease;
}
.dark .modal-content {
  background: #1C1C1E;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 28px;
  border-bottom: 1px solid #e5e7eb;
}
.dark .modal-header {
  border-bottom-color: #374151;
}

.modal-header h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
}
.dark .modal-header h2 {
  color: #f9fafb;
}

.close-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: background 0.2s ease;
  color: #6b7280;
}
.close-btn:hover {
  background: #f3f4f6;
}
.dark .close-btn:hover {
  background: #374151;
}

.modal-body {
  padding: 24px 28px;
}

.modal-description {
  font-size: 0.95rem;
  color: #6b7280;
  margin-bottom: 24px;
  line-height: 1.6;
}
.dark .modal-description {
  color: #9ca3af;
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
.dark .form-group label {
  color: #d1d5db;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  font-size: 1rem;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  background: #f9fafb;
  color: #1f2937;
  transition: all 0.2s ease;
  box-sizing: border-box;
}
.dark .form-group input {
  background: #111827;
  border-color: #374151;
  color: #f9fafb;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
  background: #fff;
}
.dark .form-group input:focus {
  background: #1f2937;
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
  animation: slideDown 0.3s ease;
}

.modal-message.success {
  background: #d1fae5;
  color: #065f46;
  border: 1px solid #6ee7b7;
}
.dark .modal-message.success {
  background: #064e3b;
  color: #6ee7b7;
}

.modal-message.error {
  background: #fee2e2;
  color: #991b1b;
  border: 1px solid #fca5a5;
}
.dark .modal-message.error {
  background: #7f1d1d;
  color: #fca5a5;
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 20px 28px;
  border-top: 1px solid #e5e7eb;
}
.dark .modal-footer {
  border-top-color: #374151;
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
  border-color: #d1d5db;
}
.dark .btn-cancel {
  border-color: #374151;
  color: #d1d5db;
}
.dark .btn-cancel:hover:not(:disabled) {
  background: #374151;
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
  transform: none;
}
</style>
