<script setup>
import { ref } from "vue"
import { gls } from "@/stores/global.js"
import { useRouter } from 'vue-router'

const router = useRouter()

let dj = ref(0)
let lj = ref(0)
let domain = ref('')
let nboffer = ref(0)
let phone  = ref('')
let email = ref('')

let passwordt = ref('password')
let passwordt2 = ref('password')

let newusername = ref('')
let newphone = ref('')
let newemail = ref('')
let newoffer = ref(0)

let passf1 = ref('')
let passf2 = ref('')

let usernameR = ref(false)
let emailR = ref(false)
let phonenumberR = ref(false)
let loading = ref(true)
let saving = ref(false)

// Modal pour suppression
let showDeleteModal = ref(false)
let deleteConfirm = ref('')
let deleting = ref(false)

// 2FA/TOTP variables
let twoFactorEnabled = ref(false)
let loading2FA = ref(true)
let showEnableModal = ref(false)
let showBackupCodesModal = ref(false)
let showDisableModal = ref(false)
let qrCodeImage = ref('')
let totpSecret = ref('')
let backupCodes = ref([])
let verifyCode = ref('')
let disablePassword = ref('')
let enabling2FA = ref(false)
let verifying2FA = ref(false)
let disabling2FA = ref(false)
let regeneratingCodes = ref(false)

// Load 2FA status
fetch(`${import.meta.env.VITE_APP_API}/api/2fa/status`, {
  method: 'POST',
  mode: 'cors',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: gls().username, token: gls().sessionT })
})
  .then(r => r.json())
  .then(data => {
    twoFactorEnabled.value = data.enabled || false
    loading2FA.value = false
  })
  .catch(() => {
    loading2FA.value = false
  })

// Récupère les infos utilisateur
fetch(import.meta.env.VITE_APP_API_INFO_USER, {
  method: "POST",
  mode: "cors",
  body: JSON.stringify({ "username": gls().username, "token": gls().sessionT })
})
  .then(r => r.json())
  .then(a => {
    dj.value = a['DateJoined']
    domain.value = a['Domain']
    nboffer.value = a['Nboffer']
    email.value = a['Email']
    phone.value = a['PhoneNumber']
    lj.value = a["LastLogin"]
    loading.value = false
  })
  .catch(() => {
    loading.value = false
  })

function send() {
  usernameR.value = false
  emailR.value = false
  phonenumberR.value = false

  if (passf1.value !== passf2.value) {
    return
  }

  saving.value = true

  fetch(import.meta.env.VITE_APP_API_INFO_CHANGE, {
    method: "POST",
    mode: "cors",
    body: JSON.stringify({
      "lastusername": gls().username,
      "username": newusername.value,
      "phonenumber": newphone.value,
      "email": newemail.value,
      "nboffer": newoffer.value,
      "password": passf2.value,
      "token": gls().sessionT
    })
  })
    .then(r => r.json())
    .then(a => {
      saving.value = false
      
      // Si le serveur retourne des erreurs de validation (vInfo)
      if (a["username"] || a["email"] || a["phone"]) {
        if (a["username"]) usernameR.value = true
        if (a["email"]) emailR.value = true
        if (a["phone"]) phonenumberR.value = true
        return
      }
      
      // Si changement de username (sessionSend avec Token)
      if (a["Token"]) {
        // Mettre à jour les cookies avec la nouvelle session
        document.cookie = `name=${gls().username}; expires=Fri, 31 Dec 1900 23:59:59 GMT; Secure`
        document.cookie = `sessionToken=${gls().sessionT}; expires=Fri, 31 Dec 1900 23:59:59 GMT; Secure`
        gls().username = a["Username"]
        gls().sessionT = a["Token"]
        document.cookie = `name=${a["Username"]}; expires=${a["Expiry"]}; Secure`
        document.cookie = `sessionToken=${a["Token"]}; expires=${a["Expiry"]}; Secure`
        window.location.href = "/account"
      } 
      // Sinon, succès simple (pas de changement username)
      else if (a["success"]) {
        window.location.href = "/account"
      }
    })
    .catch(() => {
      saving.value = false
    })
}

function openDeleteModal() {
  showDeleteModal.value = true
  deleteConfirm.value = ''
}

function closeDeleteModal() {
  showDeleteModal.value = false
  deleteConfirm.value = ''
}

function deleteAccount() {
  if (deleteConfirm.value !== gls().username) {
    return
  }

  deleting.value = true

  fetch(import.meta.env.VITE_APP_API_DELETE_ACCOUNT, {
    method: "POST",
    mode: "cors",
    body: JSON.stringify({
      "username": gls().username,
      "token": gls().sessionT
    })
  })
    .then(r => r.json())
    .then(a => {
      deleting.value = false
      if (a.success) {
        // Clear cookies and redirect to login
        document.cookie = `name=${gls().username}; expires=Fri, 31 Dec 1900 23:59:59 GMT; Secure`
        document.cookie = `sessionToken=${gls().sessionT}; expires=Fri, 31 Dec 1900 23:59:59 GMT; Secure`
        gls().username = ''
        gls().sessionT = ''
        window.location.href = '/login'
      }
    })
    .catch(() => {
      deleting.value = false
    })
}

function togglePassword1() {
  passwordt.value = (passwordt.value === "password") ? "text" : "password"
}
function togglePassword2() {
  passwordt2.value = (passwordt2.value === "password") ? "text" : "password"
}

// 2FA Functions
async function enable2FA() {
  enabling2FA.value = true
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/2fa/enable`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: gls().username, token: gls().sessionT })
    })
    const data = await response.json()
    if (response.ok) {
      qrCodeImage.value = data.qr_code
      totpSecret.value = data.secret
      backupCodes.value = data.backup_codes
      showEnableModal.value = true
    }
  } catch (error) {
    console.error('Erreur activation 2FA:', error)
  } finally {
    enabling2FA.value = false
  }
}

async function verify2FA() {
  if (!verifyCode.value || verifyCode.value.length !== 6) {
    return
  }
  
  verifying2FA.value = true
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/2fa/verify`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        username: gls().username, 
        token: gls().sessionT,
        code: verifyCode.value 
      })
    })
    
    if (response.ok) {
      twoFactorEnabled.value = true
      showEnableModal.value = false
      showBackupCodesModal.value = true
      verifyCode.value = ''
    } else {
      alert('Code invalide. Vérifiez votre authenticator.')
    }
  } catch (error) {
    console.error('Erreur vérification 2FA:', error)
  } finally {
    verifying2FA.value = false
  }
}

async function disable2FA() {
  if (!disablePassword.value) {
    return
  }
  
  disabling2FA.value = true
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/2fa/disable`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ 
        username: gls().username, 
        token: gls().sessionT,
        password: disablePassword.value 
      })
    })
    
    if (response.ok) {
      twoFactorEnabled.value = false
      showDisableModal.value = false
      disablePassword.value = ''
    } else {
      alert('Mot de passe incorrect')
    }
  } catch (error) {
    console.error('Erreur désactivation 2FA:', error)
  } finally {
    disabling2FA.value = false
  }
}

async function regenerateBackupCodes() {
  regeneratingCodes.value = true
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/2fa/backup-codes`, {
      method: 'POST',
      mode: 'cors',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: gls().username, token: gls().sessionT })
    })
    const data = await response.json()
    if (response.ok) {
      backupCodes.value = data.backup_codes
      showBackupCodesModal.value = true
    }
  } catch (error) {
    console.error('Erreur génération codes de secours:', error)
  } finally {
    regeneratingCodes.value = false
  }
}

function copyBackupCodes() {
  const text = backupCodes.value.join('\n')
  navigator.clipboard.writeText(text).then(() => {
    alert('Codes de secours copiés!')
  })
}

function closeEnableModal() {
  showEnableModal.value = false
  verifyCode.value = ''
  qrCodeImage.value = ''
  totpSecret.value = ''
}
</script>

<template>
  <div class="container">
    <header class="header">
      <RouterLink class="back-btn" to="/account">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"></line>
          <polyline points="12 19 5 12 12 5"></polyline>
        </svg>
        {{ $t('back') || 'Back' }}
      </RouterLink>
      <h1 class="title">
        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"></path>
          <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"></path>
        </svg>
        {{ $t('edit') || 'Edit Profile' }}
      </h1>
    </header>

    <div v-if="loading" class="loading">
      <div class="spinner"></div>
    </div>

    <section v-else class="card">
      <div class="avatar-section">
        <img src="@/assets/napo.png" alt="avatar" class="avatar" />
        <button class="change-avatar-btn">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M23 19a2 2 0 0 1-2 2H3a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h4l2-3h6l2 3h4a2 2 0 0 1 2 2z"></path>
            <circle cx="12" cy="13" r="4"></circle>
          </svg>
          {{ $t('picturep') || 'Change' }}
        </button>
      </div>

      <form class="form-grid" @submit.prevent="send" novalidate>
        <!-- Username -->
        <div class="form-group">
          <label class="label">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
              <circle cx="12" cy="7" r="4"></circle>
            </svg>
            {{ $t('username') }}
          </label>
          <input v-model="newusername" :placeholder="gls().username" class="input" />
          <p v-if="usernameR" class="error">{{ $t('dejaUP') }}</p>
        </div>

        <!-- Password -->
        <div class="form-group">
          <label class="label">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
              <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
            </svg>
            {{ $t('password') }}
          </label>
          <div class="password-input">
            <input :type="passwordt" v-model="passf1" :placeholder="$t('passwordN') || 'New password'" class="input" />
            <button type="button" class="icon-btn" @click="togglePassword1">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                <circle cx="12" cy="12" r="3"></circle>
              </svg>
            </button>
          </div>
          <div class="password-input">
            <input :type="passwordt2" v-model="passf2" :placeholder="$t('repassword') || 'Confirm password'" class="input" />
            <button type="button" class="icon-btn" @click="togglePassword2">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"></path>
                <circle cx="12" cy="12" r="3"></circle>
              </svg>
            </button>
          </div>
          <p v-if="passf1 !== passf2 && passf2 !== ''" class="error">{{ $t('passwordd') }}</p>
        </div>

        <!-- Email -->
        <div class="form-group">
          <label class="label">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
              <polyline points="22,6 12,13 2,6"></polyline>
            </svg>
            {{ $t('emaily') }}
          </label>
          <input v-model="newemail" :placeholder="email" class="input" type="email" />
          <p v-if="emailR" class="error">{{ $t('dejaEP') }}</p>
        </div>

        <!-- Phone -->
        <div class="form-group">
          <label class="label">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
            </svg>
            {{ $t('phoney') }}
          </label>
          <input v-model="newphone" :placeholder="phone" class="input" type="tel" />
          <p v-if="phonenumberR" class="error">{{ $t('dejaPP') }}</p>
        </div>

        <!-- Domain (read-only) -->
        <div class="form-group">
          <label class="label">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="2" y1="12" x2="22" y2="12"></line>
              <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
            </svg>
            {{ $t('domainy') }}
          </label>
          <input readonly :value="domain || '-'" class="input readonly" />
        </div>

        <!-- Offer -->
        <div class="form-group">
          <label class="label">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <polygon points="10,8 16,12 10,16"></polygon>
            </svg>
            {{ $t('offery') }}
          </label>
          <select v-model="newoffer" class="input">
            <option value="0">Free</option>
            <option value="1">Standard</option>
            <option value="2">Premium</option>
          </select>
          <RouterLink class="info-link" to="/about">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"></circle>
              <line x1="12" y1="16" x2="12" y2="12"></line>
              <line x1="12" y1="8" x2="12.01" y2="8"></line>
            </svg>
            {{ $t('About') || 'Learn more' }}
          </RouterLink>
        </div>
      </form>

      <!-- Actions -->
      <div class="actions">
        <button type="button" class="btn primary" @click="send" :disabled="saving">
          <svg v-if="!saving" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <div v-else class="btn-spinner"></div>
          {{ saving ? ($t('saving') || 'Saving...') : ($t('save') || 'Save Changes') }}
        </button>
        <RouterLink class="btn ghost" to="/account">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
          {{ $t('cancel') || 'Cancel' }}
        </RouterLink>
      </div>

      <!-- Authentication Section -->
      <div class="auth-zone">
        <h3 class="auth-title">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
          </svg>
          Authentification
        </h3>
        <p class="auth-subtitle">Sécurisez votre compte avec l'authentification à deux facteurs</p>

        <div v-if="loading2FA" style="text-align: center; padding: 20px;">
          <div class="spinner"></div>
        </div>

        <div v-else class="two-factor-content">
          <div v-if="twoFactorEnabled" class="status-enabled">
            <div class="status-badge success">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
              2FA Activé
            </div>
            <p class="status-text">Votre compte est protégé par l'authentification à deux facteurs.</p>
            <div class="actions">
              <button type="button" class="btn secondary" @click="regenerateBackupCodes" :disabled="regeneratingCodes">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="23 4 23 10 17 10"></polyline>
                  <polyline points="1 20 1 14 7 14"></polyline>
                  <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"></path>
                </svg>
                Régénérer codes de secours
              </button>
              <button type="button" class="btn danger" @click="showDisableModal = true">
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"></circle>
                  <line x1="15" y1="9" x2="9" y2="15"></line>
                  <line x1="9" y1="9" x2="15" y2="15"></line>
                </svg>
                Désactiver 2FA
              </button>
            </div>
          </div>

          <div v-else class="status-disabled">
            <div class="status-badge warning">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
                <line x1="12" y1="9" x2="12" y2="13"></line>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
              2FA Désactivé
            </div>
            <p class="status-text">Protégez votre compte avec l'authentification à deux facteurs.</p>
            <button type="button" class="btn primary" @click="enable2FA" :disabled="enabling2FA">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
              {{ enabling2FA ? 'Activation...' : 'Activer 2FA' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Danger Zone -->
      <div class="danger-zone">
        <h3 class="danger-title">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
            <line x1="12" y1="9" x2="12" y2="13"></line>
            <line x1="12" y1="17" x2="12.01" y2="17"></line>
          </svg>
          {{ $t('dangerZone') || 'Danger Zone' }}
        </h3>
        <div class="danger-content">
          <div>
            <p class="danger-text">{{ $t('deleteAccountWarning') || 'Once you delete your account, there is no going back. Please be certain.' }}</p>
          </div>
          <button type="button" class="btn danger" @click="openDeleteModal">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="3 6 5 6 21 6"></polyline>
              <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
            </svg>
            {{ $t('deleteAccount') || 'Delete Account' }}
          </button>
        </div>
      </div>
    </section>

    <!-- Enable 2FA Modal -->
    <Teleport to="body">
      <div v-if="showEnableModal" class="modal-overlay" @click="closeEnableModal">
        <div class="modal large" @click.stop>
          <div class="modal-header">
            <h2 class="modal-title">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
              Activer l'authentification à deux facteurs
            </h2>
            <button class="modal-close" @click="closeEnableModal">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="setup-steps">
              <div class="step">
                <h3>Étape 1 : Scannez le QR Code</h3>
                <p>Utilisez une application d'authentification comme Google Authenticator, Authy, ou Microsoft Authenticator.</p>
                <div class="qr-container">
                  <img :src="'data:image/png;base64,' + qrCodeImage" alt="QR Code" class="qr-code" />
                </div>
                <p class="secret-label">Ou entrez manuellement ce secret :</p>
                <code class="secret-code">{{ totpSecret }}</code>
              </div>
              
              <div class="step">
                <h3>Étape 2 : Entrez le code de vérification</h3>
                <p>Entrez le code à 6 chiffres affiché dans votre application.</p>
                <input 
                  v-model="verifyCode" 
                  type="text" 
                  maxlength="6" 
                  placeholder="000000" 
                  class="input code-input"
                  @input="verifyCode = verifyCode.replace(/[^0-9]/g, '')"
                />
              </div>
            </div>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn ghost" @click="closeEnableModal" :disabled="verifying2FA">
              Annuler
            </button>
            <button 
              type="button" 
              class="btn primary" 
              @click="verify2FA" 
              :disabled="verifyCode.length !== 6 || verifying2FA"
            >
              <div v-if="verifying2FA" class="btn-spinner"></div>
              {{ verifying2FA ? 'Vérification...' : 'Vérifier et activer' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Backup Codes Modal -->
    <Teleport to="body">
      <div v-if="showBackupCodesModal" class="modal-overlay" @click="showBackupCodesModal = false">
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h2 class="modal-title">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20 6 9 17 4 12"></polyline>
              </svg>
              Codes de secours
            </h2>
            <button class="modal-close" @click="showBackupCodesModal = false">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="warning-box">
              <strong>⚠️ Important :</strong>
              <p>Conservez ces codes dans un endroit sûr. Chaque code ne peut être utilisé qu'une seule fois.</p>
            </div>
            <div class="backup-codes-grid">
              <code v-for="(code, index) in backupCodes" :key="index" class="backup-code">
                {{ code }}
              </code>
            </div>
          </div>
          <div class="modal-actions">
            <button type="button" class="btn secondary" @click="copyBackupCodes">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
              </svg>
              Copier les codes
            </button>
            <button type="button" class="btn primary" @click="showBackupCodesModal = false">
              J'ai sauvegardé les codes
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Disable 2FA Modal -->
    <Teleport to="body">
      <div v-if="showDisableModal" class="modal-overlay" @click="showDisableModal = false">
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h2 class="modal-title danger">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
                <line x1="12" y1="9" x2="12" y2="13"></line>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
              Désactiver la 2FA
            </h2>
            <button class="modal-close" @click="showDisableModal = false">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="warning-box">
              <strong>⚠️ Attention :</strong>
              <p>Désactiver la 2FA rendra votre compte moins sécurisé.</p>
            </div>
            <label class="label">Entrez votre mot de passe pour confirmer :</label>
            <input 
              v-model="disablePassword" 
              type="password" 
              placeholder="Mot de passe" 
              class="input"
            />
          </div>
          <div class="modal-actions">
            <button type="button" class="btn ghost" @click="showDisableModal = false" :disabled="disabling2FA">
              Annuler
            </button>
            <button 
              type="button" 
              class="btn danger" 
              @click="disable2FA" 
              :disabled="!disablePassword || disabling2FA"
            >
              <div v-if="disabling2FA" class="btn-spinner"></div>
              {{ disabling2FA ? 'Désactivation...' : 'Désactiver 2FA' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="showDeleteModal" class="modal-overlay" @click="closeDeleteModal">
        <div class="modal" @click.stop>
          <div class="modal-header">
            <h2 class="modal-title">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
                <line x1="12" y1="9" x2="12" y2="13"></line>
                <line x1="12" y1="17" x2="12.01" y2="17"></line>
              </svg>
              {{ $t('confirmDelete') || 'Confirm Account Deletion' }}
            </h2>
            <button class="modal-close" @click="closeDeleteModal">
              <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"></line>
                <line x1="6" y1="6" x2="18" y2="18"></line>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <p class="modal-text">{{ $t('deleteConfirmText') || 'This action cannot be undone. This will permanently delete your account and remove all your data from our servers.' }}</p>
            <p class="modal-instruction">{{ $t('typeUsername') || 'Please type your username' }} <strong>{{ gls().username }}</strong> {{ $t('toConfirm') || 'to confirm' }}:</p>
            <input v-model="deleteConfirm" :placeholder="gls().username" class="input" />
          </div>
          <div class="modal-actions">
            <button type="button" class="btn ghost" @click="closeDeleteModal" :disabled="deleting">
              {{ $t('cancel') || 'Cancel' }}
            </button>
            <button 
              type="button" 
              class="btn danger" 
              @click="deleteAccount" 
              :disabled="deleteConfirm !== gls().username || deleting"
            >
              <div v-if="deleting" class="btn-spinner"></div>
              {{ deleting ? ($t('deleting') || 'Deleting...') : ($t('deleteAccount') || 'Delete Account') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
* { box-sizing: border-box; font-family: 'Roboto', sans-serif; }

.container {
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

.dark .container { background: transparent; }

.header {
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

.dark .header { background: #1C1C1E; }

.back-btn {
  position: absolute;
  top: 24px;
  left: 32px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 16px;
  color: #495057;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s ease;
  border: 2px solid rgba(0, 48, 143, 0.2);
}
.dark .back-btn {
  background: rgba(30, 30, 40, 0.95);
  border-color: rgba(255, 255, 255, 0.2);
  color: #e5e7eb;
}
.back-btn:hover {
  background: white;
  transform: translateX(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
.dark .back-btn:hover {
  background: rgba(40, 40, 50, 0.95);
}

.title {
  font-size: 3rem;
  font-weight: 700;
  color: #222;
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 0;
  letter-spacing: 2px;
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

.card {
  width: 100%;
  background: transparent;
  border-radius: 0;
  padding: 0;
  box-shadow: none;
}
.dark .card {
  background: transparent;
  color: #eee;
  box-shadow: none;
}

.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;
  padding: 40px 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 32px;
}
.dark .avatar-section { 
  background: #1C1C1E;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 8px 24px rgba(0,0,0,0.1);
}
.change-avatar-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s ease;
  color: #495057;
}
.dark .change-avatar-btn {
  background: #2a2d3a;
  border-color: #3a3d4a;
  color: #e5e7eb;
}
.change-avatar-btn:hover {
  background: #e9ecef;
  transform: translateY(-2px);
}
.dark .change-avatar-btn:hover {
  background: #3a3d4a;
}

.form-grid {
  display: grid;
  gap: 32px;
  padding: 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 32px;
}

.dark .form-grid {
  background: #1C1C1E;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
  color: #333;
  font-size: 0.95rem;
  text-transform: uppercase;
  letter-spacing: 1px;
}
.dark .label { color: #eee; }

.input {
  width: 100%;
  padding: 14px 18px;
  border-radius: 12px;
  border: 2px solid rgba(0, 48, 143, 0.2);
  background: white;
  font-size: 1rem;
  transition: all 0.3s ease;
  font-family: 'Roboto', sans-serif;
}
.input:focus {
  outline: none;
  border-color: blue;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(0, 48, 143, 0.1);
}
.dark .input {
  background: rgba(30, 30, 40, 0.95);
  border-color: rgba(255, 255, 255, 0.2);
  color: white;
}
.dark .input:focus {
  background: rgba(30, 30, 40, 0.95);
  border-color: rgba(255, 255, 255, 0.5);
  box-shadow: 0 0 0 3px rgba(255, 255, 255, 0.1);
}
.input.readonly {
  cursor: not-allowed;
  opacity: 0.6;
}

.password-input {
  position: relative;
  display: flex;
  align-items: center;
  gap: 8px;
}

.icon-btn {
  position: absolute;
  right: 12px;
  border: none;
  background: transparent;
  cursor: pointer;
  padding: 4px;
  color: #6c757d;
  transition: color 0.3s ease;
}
.icon-btn:hover { color: #495057; }
.dark .icon-btn { color: #9ca3af; }
.dark .icon-btn:hover { color: #e5e7eb; }

.info-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #667eea;
  text-decoration: none;
  font-size: 13px;
  font-weight: 500;
  margin-top: 4px;
}
.info-link:hover { text-decoration: underline; }

.error {
  color: #dc3545;
  font-size: 13px;
  margin-top: 4px;
}

.actions {
  display: flex;
  gap: 16px;
  justify-content: center;
  padding: 24px 32px;
  background: rgba(245, 245, 247, 0.5);
  border-radius: 24px;
  margin-bottom: 32px;
}
.dark .actions { 
  background: rgba(30, 30, 40, 0.5);
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 14px 32px;
  border-radius: 16px;
  border: none;
  cursor: pointer;
  font-weight: 600;
  font-size: 1.05rem;
  transition: all 0.3s ease;
  text-decoration: none;
  font-family: 'Roboto', sans-serif;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}
.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.btn.primary {
  background: -webkit-linear-gradient(30deg, blue, red);
  color: #fff;
}
.btn.primary:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}
.btn.secondary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  box-shadow: 0 4px 12px rgba(102,126,234,0.3);
}
.btn.secondary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102,126,234,0.4);
}
.btn.ghost {
  background: transparent;
  color: #6c757d;
  border: 1px solid #dee2e6;
}
.btn.ghost:hover {
  background: #f8f9fa;
  border-color: #adb5bd;
}
.dark .btn.ghost {
  color: #9ca3af;
  border-color: #3a3d4a;
}
.dark .btn.ghost:hover {
  background: #2a2d3a;
}
.btn.danger {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: #fff;
  box-shadow: 0 4px 12px rgba(245,87,108,0.3);
}
.btn.danger:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(245,87,108,0.4);
}

.btn-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

.danger-zone {
  margin-top: 48px;
  padding: 24px;
  background: linear-gradient(135deg, rgba(245,87,108,0.05) 0%, rgba(240,147,251,0.05) 100%);
  border: 2px solid rgba(245,87,108,0.2);
  border-radius: 16px;
}
.dark .danger-zone {
  background: rgba(245,87,108,0.08);
  border-color: rgba(245,87,108,0.3);
}

/* Authentication Zone */
.auth-zone {
  margin-top: 32px;
  padding: 24px;
  background: linear-gradient(135deg, rgba(59,130,246,0.05) 0%, rgba(99,102,241,0.05) 100%);
  border: 2px solid rgba(59,130,246,0.2);
  border-radius: 16px;
}
.dark .auth-zone {
  background: rgba(59,130,246,0.08);
  border-color: rgba(59,130,246,0.3);
}

.auth-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 700;
  color: #3b82f6;
  margin: 0 0 8px;
}
.dark .auth-title { color: #60a5fa; }

.auth-subtitle {
  font-size: 14px;
  color: #6c757d;
  margin: 0 0 20px;
}
.dark .auth-subtitle { color: #9ca3af; }

.danger-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 700;
  color: #dc3545;
  margin: 0 0 16px;
}
.dark .danger-title { color: #f5576c; }

.danger-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
}

.danger-text {
  font-size: 14px;
  color: #6c757d;
  margin: 0;
}
.dark .danger-text { color: #9ca3af; }

/* Modal */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 16px;
}

.modal {
  background: #fff;
  border-radius: 20px;
  max-width: 500px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 60px rgba(0,0,0,0.3);
  animation: modalSlide 0.3s ease;
}
.dark .modal { background: #1f2230; }

@keyframes modalSlide {
  from { opacity: 0; transform: translateY(-20px); }
  to { opacity: 1; transform: translateY(0); }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-bottom: 1px solid #e9ecef;
}
.dark .modal-header { border-bottom-color: #3a3d4a; }

.modal-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 20px;
  font-weight: 700;
  color: #dc3545;
  margin: 0;
}
.dark .modal-title { color: #f5576c; }

.modal-close {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 4px;
  color: #6c757d;
  transition: color 0.3s ease;
}
.modal-close:hover { color: #dc3545; }

.modal-body {
  padding: 24px;
}

.modal-text {
  font-size: 15px;
  color: #495057;
  margin: 0 0 20px;
  line-height: 1.6;
}
.dark .modal-text { color: #e5e7eb; }

.modal-instruction {
  font-size: 14px;
  color: #6c757d;
  margin: 0 0 12px;
}
.dark .modal-instruction { color: #9ca3af; }

.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 20px 24px;
  border-top: 1px solid #e9ecef;
}
.dark .modal-actions { border-top-color: #3a3d4a; }

/* Responsive */
@media (max-width: 768px) {
  .title { font-size: 22px; }
  .card { padding: 24px; }
  .danger-content { flex-direction: column; align-items: stretch; }
  .actions, .modal-actions {
    flex-direction: column;
  }
  .btn { justify-content: center; }
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f3f3;
  border-top: 4px solid #007bff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.btn-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid white;
  border-top: 2px solid transparent;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-right: 8px;
}

/* 2FA Styles */
.two-factor-content {
  padding: 20px 0;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 10px 16px;
  border-radius: 12px;
  font-weight: 600;
  font-size: 14px;
  margin-bottom: 12px;
}

.status-badge.success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.status-badge.warning {
  background: #fff3cd;
  color: #856404;
  border: 1px solid #ffeaa7;
}

.dark .status-badge.success {
  background: #1e3a2e;
  color: #4ade80;
  border-color: #2a5a3e;
}

.dark .status-badge.warning {
  background: #3a3020;
  color: #fbbf24;
  border-color: #4a4030;
}

.status-text {
  color: #6c757d;
  margin: 16px 0;
  font-size: 15px;
}

.dark .status-text {
  color: #9ca3af;
}

.setup-steps {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

.step h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #212529;
}

.dark .step h3 {
  color: #f8f9fa;
}

.step p {
  color: #6c757d;
  font-size: 14px;
  margin-bottom: 16px;
}

.dark .step p {
  color: #9ca3af;
}

.qr-container {
  display: flex;
  justify-content: center;
  padding: 20px;
  background: white;
  border-radius: 12px;
  border: 2px solid #e9ecef;
  margin: 16px 0;
}

.dark .qr-container {
  background: #1a1d2e;
  border-color: #3a3d4a;
}

.qr-code {
  max-width: 256px;
  width: 100%;
  height: auto;
}

.secret-label {
  font-size: 13px;
  color: #6c757d;
  margin: 16px 0 8px;
}

.dark .secret-label {
  color: #9ca3af;
}

.secret-code {
  display: block;
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 12px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  color: #495057;
  text-align: center;
  word-break: break-all;
}

.dark .secret-code {
  background: #1a1d2e;
  border-color: #3a3d4a;
  color: #e5e7eb;
}

.code-input {
  text-align: center;
  font-size: 24px;
  letter-spacing: 8px;
  font-family: 'Courier New', monospace;
  font-weight: bold;
}

.warning-box {
  background: #fff3cd;
  border: 1px solid #ffc107;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 20px;
}

.warning-box strong {
  color: #856404;
  display: block;
  margin-bottom: 8px;
}

.warning-box p {
  color: #856404;
  margin: 0;
  font-size: 14px;
}

.dark .warning-box {
  background: #3a3020;
  border-color: #4a4030;
}

.dark .warning-box strong,
.dark .warning-box p {
  color: #fbbf24;
}

.backup-codes-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-top: 16px;
}

.backup-code {
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 8px;
  padding: 12px;
  text-align: center;
  font-family: 'Courier New', monospace;
  font-size: 16px;
  font-weight: bold;
  color: #495057;
}

.dark .backup-code {
  background: #1a1d2e;
  border-color: #3a3d4a;
  color: #e5e7eb;
}

.modal.large {
  max-width: 600px;
}

@media (max-width: 768px) {
  .backup-codes-grid {
    grid-template-columns: 1fr;
  }
}
</style>

