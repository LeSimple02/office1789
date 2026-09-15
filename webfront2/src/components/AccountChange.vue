<script setup>
import { ref, computed } from "vue"
import { gls } from "@/stores/global.js"
import { useRouter, RouterLink } from 'vue-router'
import { validatePassword, isValidEmail, isValidPhone, getStrengthColor } from '@/utils/validation'

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

// Vérification email (code via /api/verification)
let emailCode = ref('')
let emailVerificationSending = ref(false)
let emailVerificationChecking = ref(false)
let emailVerificationSent = ref(false)
let emailVerificationSuccess = ref(false)
let emailVerificationError = ref('')

// Validation states
let passwordStrength = ref({ valid: true, errors: [], strength: 0, strengthLabel: '' })
let emailValid = ref(true)
let phoneValid = ref(true)

// Computed password strength color
const strengthColor = computed(() => getStrengthColor(passwordStrength.value.strength))

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

// Subscription/Payment variables
let showSubscriptionModal = ref(false)
let selectedPlan = ref(0)
let processingPayment = ref(false)
let paymentMessage = ref('')

const plans = [
  { 
    id: 0, 
    name: 'Free', 
    price: 0, 
    storage: '1GB', 
    features: [
      'Email basique', 
      'Chat Matrix', 
      'Calendrier', 
      '1GB de stockage', 
      'Pas de partage de fichiers'
    ] 
  },
  { 
    id: 1, 
    name: 'Standard', 
    price: 5, 
    storage: '50GB', 
    features: [
      'Email professionnel', 
      'Support prioritaire ⭐', 
      '50GB de stockage', 
      'Pas de partage de fichiers',
      'Sauvegardes quotidiennes'
    ] 
  },
  { 
    id: 2, 
    name: 'Professional', 
    price: 12, 
    storage: '200GB', 
    features: [
      'Collaboration d\'équipe', 
      'Support prioritaire ⭐', 
      '200GB de stockage', 
      'Partage avec 3 membres',
      'Historique 90 jours',
      'API Access'
    ] 
  },
  { 
    id: 3, 
    name: 'Enterprise', 
    price: 49, 
    storage: 'Illimité', 
    features: [
      'Sécurité avancée', 
      'Support prioritaire 24/7 ⭐⭐', 
      'Stockage illimité', 
      'Partage avec 20 membres', 
      'SSO & SAML',
      'SLA 99.9%',
      'Gestion avancée'
    ] 
  }
]

// Load 2FA status
fetch(`${import.meta.env.VITE_APP_API}/api/2fa/status`, {
  method: 'POST',
  mode: 'cors',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: gls().username, token: gls().sessionT })
})
  .then(async r => {
    const text = await r.text()
    try {
      return JSON.parse(text)
    } catch (e) {
      console.error('2FA status JSON parse error:', e, 'Response:', text)
      return { enabled: false }
    }
  })
  .then(data => {
    twoFactorEnabled.value = data.enabled || false
    loading2FA.value = false
  })
  .catch((err) => {
    console.error('2FA status fetch error:', err)
    loading2FA.value = false
  })

// Récupère les infos utilisateur
fetch(import.meta.env.VITE_APP_API_GETINFOP, {
  method: "POST",
  mode: "cors",
  body: JSON.stringify({ "username": gls().username, "token": gls().sessionT })
})
  .then(async r => {
    const text = await r.text()
    try {
      return JSON.parse(text)
    } catch (e) {
      console.error('User info JSON parse error:', e, 'Response:', text)
      throw new Error('Invalid JSON response')
    }
  })
  .then(a => {
    dj.value = a['DateJoined']
    domain.value = a['Domain']
    nboffer.value = a['Nboffer']
    email.value = a['Email']
    phone.value = a['PhoneNumber']
    lj.value = a["LastLogin"]
    loading.value = false
  })
  .catch((err) => {
    console.error('User info fetch error:', err)
    loading.value = false
  })

// Check password strength
function checkPasswordStrength() {
  if (passf1.value) {
    passwordStrength.value = validatePassword(passf1.value)
  }
}

// Validate email format
function checkEmail() {
  if (newemail.value) {
    emailValid.value = isValidEmail(newemail.value)
  } else {
    emailValid.value = true
  }
}

// Validate phone format
function checkPhone() {
  if (newphone.value && newphone.value.trim() !== '') {
    phoneValid.value = isValidPhone(newphone.value)
  } else {
    phoneValid.value = true // Optional field
  }
}

// Envoi du code de vérification pour le nouvel email
async function sendEmailVerification() {
  if (!newemail.value || !emailValid.value) return

  emailVerificationSending.value = true
  emailVerificationMessage.value = ''

  try {
    const res = await fetch(import.meta.env.VITE_APP_API_VERIFICATION_SEND, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        contact: newemail.value,
        type: 'email',
      }),
    })

    const data = await res.json().catch(() => ({}))
    if (data.success) {
      emailVerificationMessage.value = 'Code envoyé à votre adresse email.'
    } else {
      emailVerificationMessage.value = data.message || "Erreur lors de l'envoi du code."
    }
  } catch (e) {
    console.error('sendEmailVerification (account) error:', e)
    emailVerificationMessage.value = "Erreur réseau lors de l'envoi du code."
  } finally {
    emailVerificationSending.value = false
  }
}

// Vérification du code saisi pour le nouvel email
async function verifyEmailCode() {
  if (!newemail.value || !emailCode.value) return

  emailVerificationChecking.value = true
  emailVerificationMessage.value = ''

  try {
    const res = await fetch(import.meta.env.VITE_APP_API_VERIFICATION_VERIFY, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        contact: newemail.value,
        code: emailCode.value,
        type: 'email',
      }),
    })

    const data = await res.json().catch(() => ({}))
    if (data.success) {
      emailVerificationMessage.value = 'Email vérifié.'
    } else {
      emailVerificationMessage.value = data.message || 'Code invalide ou expiré.'
    }
  } catch (e) {
    console.error('verifyEmailCode (account) error:', e)
    emailVerificationMessage.value = 'Erreur réseau lors de la vérification.'
  } finally {
    emailVerificationChecking.value = false
  }
}

function send() {
  usernameR.value = false
  emailR.value = false
  phonenumberR.value = false

  // Valider email si modifié et non vide
  if (newemail.value && newemail.value.trim() !== '') {
    checkEmail()
    if (!emailValid.value) {
      return
    }
  }

  // Valider téléphone si modifié et non vide
  if (newphone.value && newphone.value.trim() !== '') {
    checkPhone()
    if (!phoneValid.value) {
      return
    }
  }

  // Valider mot de passe si modifié
  if (passf1.value) {
    checkPasswordStrength()
    if (!passwordStrength.value.valid) {
      return
    }
  }

  if (passf1.value !== passf2.value) {
    return
  }

  saving.value = true

  fetch(import.meta.env.VITE_APP_API_CHANGEINFO, {
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

      // Erreurs explicites renvoyées par l'API (email_not_verified, phone_not_verified, etc.)
      if (a.error) {
        if (a.error === 'email_not_verified' || a.error === 'email_taken') {
          emailR.value = true
        }
        if (a.error === 'phone_not_verified' || a.error === 'phone_taken') {
          phonenumberR.value = true
        }
        return
      }
      
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

// Subscription Management Functions
function openSubscriptionModal(planId) {
  selectedPlan.value = planId
  showSubscriptionModal.value = true
  paymentMessage.value = ''
}

function closeSubscriptionModal() {
  showSubscriptionModal.value = false
  selectedPlan.value = 0
  paymentMessage.value = ''
}

async function processSubscription() {
  const plan = plans.find(p => p.id === selectedPlan.value)
  
  if (!plan) {
    paymentMessage.value = 'Plan invalide'
    return
  }

  // Si c'est le plan gratuit, pas besoin de paiement
  if (plan.id === 0) {
    processingPayment.value = true
    paymentMessage.value = ''
    
    try {
      await updateSubscription(0)
      paymentMessage.value = 'Changement vers le plan gratuit effectué avec succès!'
      
      setTimeout(() => {
        closeSubscriptionModal()
        window.location.reload()
      }, 1500)
    } catch (error) {
      paymentMessage.value = error.message || 'Erreur lors du changement de plan'
    } finally {
      processingPayment.value = false
    }
    return
  }

  // Pour les plans payants, rediriger vers Stripe
  processingPayment.value = true
  paymentMessage.value = ''

  try {
    // Créer une session de paiement Stripe
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/stripe/checkout`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        username: gls().username,
        token: gls().sessionT,
        plan_id: selectedPlan.value
      })
    })

    // Lire la réponse brute pour debug
    const responseText = await response.text()
    console.log('Stripe checkout raw response:', responseText)
    console.log('Stripe checkout status:', response.status)
    console.log('Stripe checkout content-type:', response.headers.get('content-type'))

    // Parser le JSON
    let data
    try {
      data = JSON.parse(responseText)
    } catch (parseError) {
      console.error('Stripe checkout JSON parse error:', parseError)
      console.error('Trying to parse:', responseText)
      throw new Error('Invalid JSON response from Stripe endpoint')
    }

    if (!response.ok || !data.success) {
      throw new Error(data.message || 'Failed to create checkout session')
    }

    // Rediriger vers Stripe Checkout
    window.location.href = data.url
    
  } catch (error) {
    console.error('Error creating checkout session:', error)
    paymentMessage.value = error.message || 'Erreur lors de la création de la session de paiement'
    processingPayment.value = false
  }
}

async function updateSubscription(planId) {
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/subscription/change`, {
      method: "POST",
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        "username": gls().username,
        "nboffer": planId,
        "token": gls().sessionT
      })
    })
    
    // Lire la réponse brute pour debug
    const responseText = await response.text()
    console.log('Raw response:', responseText)
    console.log('Response status:', response.status)
    console.log('Content-Type:', response.headers.get('content-type'))
    
    // Vérifier si la réponse est bien du JSON
    const contentType = response.headers.get('content-type')
    if (!contentType || !contentType.includes('application/json')) {
      console.error('Response is not JSON:', responseText)
      throw new Error('Invalid server response format')
    }
    
    // Parser le JSON
    let data
    try {
      data = JSON.parse(responseText)
    } catch (parseError) {
      console.error('JSON parse error:', parseError)
      console.error('Trying to parse:', responseText)
      throw new Error('Invalid JSON response from server')
    }
    
    if (!response.ok || !data.success) {
      throw new Error(data.message || 'Failed to update subscription')
    }
    
    // Mettre à jour l'offre localement
    nboffer.value = planId
    paymentMessage.value = data.message || 'Abonnement mis à jour avec succès !'
    
    return data
    
  } catch (error) {
    console.error('Error updating subscription:', error)
    paymentMessage.value = error.message || 'Erreur lors de la mise à jour de l\'abonnement'
    throw error
  }
}

function getPlanName(id) {
  const names = ['Free', 'Standard', 'Professional', 'Enterprise']
  return names[id] || 'Unknown'
}
</script>

<template>
  <div class="edit-page">
    <!-- Loading -->
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <span>Chargement...</span>
    </div>

    <template v-else>
      <div class="head">
        <h1>Modifier le compte</h1>
        <p>Mettez à jour vos informations.</p>
      </div>

      <!-- Profile info -->
      <div class="card">
        <h3>Informations</h3>
        <div class="info-row">
          <label>Identifiant actuel</label>
          <span class="info-val">{{ gls().username }}</span>
        </div>
        <div class="info-row">
          <label>Offre actuelle</label>
          <span class="info-val">{{ getPlanName(nboffer) }}</span>
        </div>
        <div class="info-row">
          <label>Email actuel</label>
          <span class="info-val">{{ email || '—' }}</span>
        </div>
        <div class="info-row">
          <label>Téléphone actuel</label>
          <span class="info-val">{{ phone || '—' }}</span>
        </div>
      </div>

      <!-- Edit form -->
      <div class="card">
        <h3>Modifications</h3>
        <div class="form">
          <div class="field">
            <label>Nouvel identifiant</label>
            <input v-model="newusername" type="text" placeholder="Nouvel identifiant" />
          </div>
          <div class="field">
            <label>Nouvel email</label>
            <input v-model="newemail" type="email" placeholder="Nouvel email" @blur="checkEmail" />
            <div v-if="email && !emailValid" class="err">Format email invalide</div>
            <div v-if="email && emailValid" class="verify-row">
              <button type="button" @click="sendEmailVerification" class="btn-sec" :disabled="emailVerificationSending">
                <span v-if="emailVerificationSending" class="spinner"></span>
                Envoyer le code
              </button>
            </div>
            <div v-if="email" class="verify-row">
              <input v-model="emailCode" type="text" placeholder="Code de vérification" />
              <button type="button" @click="verifyEmailCode" class="btn-sec" :disabled="emailVerificationChecking">Vérifier</button>
            </div>
            <p v-if="emailVerificationMessage" class="verify-msg">{{ emailVerificationMessage }}</p>
            <p v-if="emailR" class="err">Email déjà utilisé</p>
          </div>
          <div class="field">
            <label>Nouveau téléphone</label>
            <input v-model="newphone" type="tel" placeholder="Nouveau téléphone" @blur="checkPhone" />
            <p v-if="phonenumberR" class="err">Téléphone déjà utilisé</p>
          </div>
          <div class="field">
            <label>Nouveau mot de passe</label>
            <div class="pwd-wrap">
              <input v-model="passf1" :type="passwordt" placeholder="Mot de passe" @input="checkPasswordStrength" />
              <button type="button" @click="togglePassword1" class="pwd-toggle">{{ passwordt === 'password' ? '👁' : '🙈' }}</button>
            </div>
            <div v-if="passf1" class="strength">
              <div class="bar-bg"><div class="bar-fill" :style="{ width: (passwordStrength.strength / 5 * 100) + '%', background: strengthColor }"></div></div>
              <span class="strength-text" :style="{ color: strengthColor }">{{ passwordStrength.strengthLabel }}</span>
            </div>
            <ul v-if="passf1 && passwordStrength.errors.length > 0" class="checklist">
              <li v-for="e in passwordStrength.errors" :key="e">{{ e }}</li>
            </ul>
          </div>
          <div class="field">
            <label>Confirmer mot de passe</label>
            <div class="pwd-wrap">
              <input v-model="passf2" :type="passwordt2" placeholder="Confirmer" />
              <button type="button" @click="togglePassword2" class="pwd-toggle">{{ passwordt2 === 'password' ? '👁' : '🙈' }}</button>
            </div>
            <p v-if="passf1 !== passf2 && passf2" class="err">Les mots de passe ne correspondent pas</p>
          </div>
          <div class="field">
            <label>Offre</label>
            <select v-model="newoffer">
              <option :value="0">Free</option>
              <option :value="1">Standard (5€)</option>
              <option :value="2">Professional (12€)</option>
              <option :value="3">Enterprise (49€)</option>
            </select>
          </div>
          <button @click="send" class="btn-save" :disabled="saving">
            <span v-if="saving" class="spinner"></span>
            {{ saving ? 'Sauvegarde...' : 'Sauvegarder' }}
          </button>
        </div>
      </div>

      <!-- 2FA -->
      <div class="card">
        <h3>Sécurité 2FA</h3>
        <div v-if="loading2FA" class="loading-inline"><div class="spinner-sm"></div> Vérification...</div>
        <div v-else>
          <p class="status-text">2FA {{ twoFactorEnabled ? 'activé' : 'désactivé' }}</p>
          <div class="2fa-actions">
            <button v-if="!twoFactorEnabled" @click="enable2FA" class="btn-2fa" :disabled="enabling2FA">Activer 2FA</button>
            <button v-else @click="showDisableModal = true" class="btn-2fa">Désactiver 2FA</button>
            <button v-if="twoFactorEnabled" @click="regenerateBackupCodes" class="btn-2fa sec" :disabled="regeneratingCodes">Codes de secours</button>
          </div>
        </div>
      </div>

      <!-- Plans -->
      <div class="card">
        <h3>Abonnements</h3>
        <div class="plans-grid">
          <div v-for="plan in plans" :key="plan.id" class="plan-card" :class="{ current: plan.id === nboffer }">
            <h4>{{ plan.name }}</h4>
            <span class="plan-price">{{ plan.price === 0 ? 'Gratuit' : plan.price + '€/mois' }}</span>
            <ul class="plan-features">
              <li v-for="f in plan.features" :key="f">{{ f }}</li>
            </ul>
            <button @click="openSubscriptionModal(plan.id)" class="plan-btn" :disabled="plan.id === nboffer">
              {{ plan.id === nboffer ? 'Plan actuel' : 'Choisir' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Danger zone -->
      <div class="card danger-card">
        <h3>Zone de danger</h3>
        <p>Supprimer votre compte est irréversible.</p>
        <button @click="openDeleteModal" class="btn-danger">Supprimer le compte</button>
      </div>
    </template>

    <!-- 2FA Enable Modal -->
    <div v-if="showEnableModal" class="overlay" @click.self="closeEnableModal">
      <div class="dialog">
        <div class="dialog-head"><h3>Activer 2FA</h3><button @click="closeEnableModal" class="close">×</button></div>
        <div class="dialog-body">
          <p>Scannez ce QR code avec votre application d'authentification :</p>
          <div v-if="qrCodeImage" class="qr-code"><img :src="qrCodeImage" alt="QR Code" /></div>
          <p class="secret-label">Ou saisissez ce secret : <code>{{ totpSecret }}</code></p>
          <p>Entrez le code à 6 chiffres :</p>
          <input v-model="verifyCode" type="text" maxlength="6" placeholder="000000" class="totp-input" @keyup.enter="verify2FA" />
        </div>
        <div class="dialog-foot"><button @click="closeEnableModal" class="btn-cancel">Annuler</button><button @click="verify2FA" class="btn-confirm" :disabled="verifying2FA">{{ verifying2FA ? '...' : 'Vérifier' }}</button></div>
      </div>
    </div>

    <!-- 2FA Disable Modal -->
    <div v-if="showDisableModal" class="overlay" @click.self="showDisableModal = false">
      <div class="dialog">
        <div class="dialog-head"><h3>Désactiver 2FA</h3><button @click="showDisableModal = false" class="close">×</button></div>
        <div class="dialog-body"><p>Entrez votre mot de passe pour confirmer :</p><input v-model="disablePassword" type="password" placeholder="Mot de passe" /></div>
        <div class="dialog-foot"><button @click="showDisableModal = false" class="btn-cancel">Annuler</button><button @click="disable2FA" class="btn-confirm" :disabled="disabling2FA">{{ disabling2FA ? '...' : 'Désactiver' }}</button></div>
      </div>
    </div>

    <!-- Backup Codes Modal -->
    <div v-if="showBackupCodesModal" class="overlay" @click.self="showBackupCodesModal = false">
      <div class="dialog">
        <div class="dialog-head"><h3>Codes de secours</h3><button @click="showBackupCodesModal = false" class="close">×</button></div>
        <div class="dialog-body">
          <p>Conservez ces codes en lieu sûr :</p>
          <div class="codes-grid">
            <code v-for="code in backupCodes" :key="code">{{ code }}</code>
          </div>
          <button @click="copyBackupCodes" class="btn-sec mt">Copier</button>
        </div>
        <div class="dialog-foot"><button @click="showBackupCodesModal = false" class="btn-confirm">Fermer</button></div>
      </div>
    </div>

    <!-- Subscription Modal -->
    <div v-if="showSubscriptionModal" class="overlay" @click.self="closeSubscriptionModal">
      <div class="dialog">
        <div class="dialog-head"><h3>Changement d'abonnement</h3><button @click="closeSubscriptionModal" class="close">×</button></div>
        <div class="dialog-body"><p>Plan sélectionné : {{ plans.find(p => p.id === selectedPlan)?.name }}</p><p v-if="paymentMessage" class="pay-msg">{{ paymentMessage }}</p></div>
        <div class="dialog-foot"><button @click="closeSubscriptionModal" class="btn-cancel">Annuler</button><button @click="processSubscription" class="btn-confirm" :disabled="processingPayment">{{ processingPayment ? 'Traitement...' : 'Confirmer' }}</button></div>
      </div>
    </div>

    <!-- Delete Modal -->
    <div v-if="showDeleteModal" class="overlay" @click.self="closeDeleteModal">
      <div class="dialog">
        <div class="dialog-head"><h3>Supprimer le compte</h3><button @click="closeDeleteModal" class="close">×</button></div>
        <div class="dialog-body">
          <p class="del-warn">Cette action est irréversible. Tapez votre identifiant pour confirmer :</p>
          <input v-model="deleteConfirm" type="text" :placeholder="gls().username" />
        </div>
        <div class="dialog-foot"><button @click="closeDeleteModal" class="btn-cancel">Annuler</button><button @click="deleteAccount" class="btn-confirm danger" :disabled="deleting || deleteConfirm !== gls().username">{{ deleting ? '...' : 'Supprimer' }}</button></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.edit-page { max-width: 560px; margin: 0 auto; padding: 2rem 1.25rem 5rem; animation: fadeInUp .5s var(--ease) both; }
.loading { display: flex; align-items: center; justify-content: center; gap: .625rem; padding: 4rem 2rem; color: var(--text-muted); font-size: .875rem; }
.spinner { width: 24px; height: 24px; border: 3px solid var(--border); border-top-color: var(--ink); border-radius: 50%; animation: spin .7s linear infinite; }
.spinner-sm { width: 16px; height: 16px; border: 2px solid var(--border); border-top-color: var(--ink); border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
.loading-inline { display: flex; align-items: center; gap: .5rem; font-size: .8125rem; color: var(--text-muted); }
.head { margin-bottom: 1.5rem; animation: fadeInDown .4s var(--ease) both; }
.head h1 { font-size: 1.5rem; font-weight: 800; letter-spacing: -.025em; margin-bottom: .375rem; }
.head p { font-size: .875rem; color: var(--text-secondary); }
.card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1.75rem; margin-bottom: 1rem; animation: fadeInUp .4s var(--ease) both; }
.card h3 { font-size: 1.0625rem; font-weight: 700; margin-bottom: 1.125rem; padding-bottom: .625rem; border-bottom: 1px solid var(--border); }
.info-row { display: flex; justify-content: space-between; align-items: center; padding: .5rem 0; border-bottom: 1px solid var(--border); }
.info-row:last-child { border-bottom: none; }
.info-row label { font-size: .8125rem; font-weight: 600; color: var(--text-secondary); }
.info-val { font-size: .875rem; font-weight: 500; color: var(--text-primary); }
.form { display: flex; flex-direction: column; gap: 1.125rem; }
.field { display: flex; flex-direction: column; }
.field label { font-size: .8125rem; font-weight: 600; color: var(--text-secondary); margin-bottom: .375rem; }
.field input, .field select { width: 100%; padding: .625rem .875rem; font-size: .9375rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); transition: border-color var(--tb), box-shadow var(--tb); }
.field input:focus, .field select:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.field input::placeholder { color: var(--text-muted); }
.pwd-wrap { position: relative; }
.pwd-wrap input { padding-right: 38px; }
.pwd-toggle { position: absolute; right: 6px; top: 50%; transform: translateY(-50%); background: none; border: none; font-size: 16px; cursor: pointer; opacity: .6; padding: 4px; border-radius: var(--r-sm); transition: opacity var(--tb); }
.pwd-toggle:hover { opacity: 1; }
.err { font-size: .75rem; color: var(--text-primary); margin-top: .375rem; font-weight: 500; }
.verify-row { display: flex; gap: .5rem; margin-top: .5rem; }
.verify-row input { flex: 1; }
.btn-sec { padding: .5625rem .875rem; font-size: .75rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; transition: all var(--tb); white-space: nowrap; display: flex; align-items: center; gap: .25rem; }
.btn-sec:hover:not(:disabled) { border-color: var(--ink); background: var(--bg-tertiary); }
.btn-sec:disabled { opacity: .5; }
.spinner { width: 14px; height: 14px; border: 2px solid var(--border); border-top-color: var(--ink); border-radius: 50%; animation: spin .6s linear infinite; display: inline-block; }
.verify-msg { font-size: .75rem; color: var(--text-secondary); margin-top: .375rem; animation: fadeIn .2s var(--ease) both; }
.strength { display: flex; align-items: center; gap: .5rem; margin-top: .375rem; }
.bar-bg { flex: 1; height: 4px; background: var(--bg-tertiary); border-radius: var(--r-full); overflow: hidden; }
.bar-fill { height: 100%; border-radius: var(--r-full); transition: width var(--tb); animation: growBar .3s var(--ease) both; }
.strength-text { font-size: .75rem; font-weight: 500; white-space: nowrap; }
.checklist { list-style: none; padding: 0; margin: .375rem 0 0; }
.checklist li { font-size: .75rem; color: var(--text-muted); margin-top: .25rem; animation: fadeIn .2s var(--ease) both; }
.btn-save { padding: .75rem; font-size: .9375rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); display: flex; align-items: center; justify-content: center; gap: .375rem; }
.btn-save:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-md); transform: translateY(-2px); }
.btn-save:disabled { opacity: .5; cursor: not-allowed; }
.status-text { font-size: .875rem; font-weight: 500; margin-bottom: .75rem; }
.2fa-actions { display: flex; gap: .5rem; flex-wrap: wrap; }
.btn-2fa { padding: .5625rem 1.125rem; font-size: .8125rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-2fa:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); }
.btn-2fa.sec { color: var(--text-primary); background: var(--bg-primary); border-color: var(--border-strong); }
.btn-2fa.sec:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.btn-2fa:disabled { opacity: .5; }
.plans-grid { display: grid; grid-template-columns: 1fr 1fr; gap: .875rem; }
.plan-card { padding: 1.25rem; background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--r); transition: all var(--tb); }
.plan-card.current { border-color: var(--ink); border-width: 2px; }
.plan-card:hover { border-color: var(--border-strong); box-shadow: var(--sh-sm); }
.plan-card h4 { font-size: .9375rem; font-weight: 700; margin-bottom: .25rem; }
.plan-price { font-size: 1.125rem; font-weight: 800; letter-spacing: -.02em; margin-bottom: .75rem; display: block; }
.plan-features { list-style: none; padding: 0; margin: 0 0 .875rem; }
.plan-features li { font-size: .6875rem; color: var(--text-secondary); margin-top: .25rem; padding-left: 10px; position: relative; }
.plan-features li::before { content: ''; position: absolute; left: 0; top: 5px; width: 4px; height: 4px; border-radius: 50%; background: var(--ink); }
.plan-btn { width: 100%; padding: .5rem; font-size: .75rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r-sm); cursor: pointer; transition: all var(--tb); }
.plan-btn:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); }
.plan-btn:disabled { opacity: .5; cursor: not-allowed; }
.danger-card { border-color: var(--border-strong); }
.danger-card h3 { color: var(--text-primary); }
.danger-card p { font-size: .8125rem; color: var(--text-secondary); margin-bottom: .75rem; }
.btn-danger { padding: .5625rem 1.125rem; font-size: .8125rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-danger:hover { border-color: var(--ink); background: var(--bg-tertiary); }

/* Modals */
.overlay { position: fixed; inset: 0; z-index: 999; background: var(--bg-overlay); backdrop-filter: blur(6px); display: flex; align-items: center; justify-content: center; animation: fadeIn .2s var(--ease) both; }
.dialog { width: calc(100% - 32px); max-width: 420px; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); box-shadow: var(--sh-xl); animation: popIn .25s var(--ease-spring) both; max-height: 90vh; display: flex; flex-direction: column; }
.dialog-head { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.25rem; border-bottom: 1px solid var(--border); }
.dialog-head h3 { font-size: 1rem; font-weight: 700; }
.close { width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; font-size: 18px; line-height: 1; background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r-sm); cursor: pointer; color: var(--text-secondary); transition: all var(--tb); }
.close:hover { border-color: var(--ink); background: var(--bg-tertiary); transform: rotate(90deg); color: var(--text-primary); }
.dialog-body { padding: 1.25rem; overflow-y: auto; }
.dialog-body p { font-size: .8125rem; color: var(--text-secondary); margin-bottom: .75rem; }
.dialog-body input { width: 100%; padding: .5625rem .75rem; font-size: .875rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); margin-bottom: .5rem; }
.dialog-body input:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.totp-input { font-family: var(--font-mono); letter-spacing: .3em; font-size: 1.125rem; text-align: center; }
.qr-code { text-align: center; margin: .75rem 0; }
.qr-code img { max-width: 200px; border-radius: var(--r); border: 1px solid var(--border); }
.secret-label { font-size: .8125rem; }
.secret-label code { font-family: var(--font-mono); font-size: .75rem; }
.codes-grid { display: grid; grid-template-columns: 1fr 1fr; gap: .375rem; margin: .75rem 0; }
.codes-grid code { padding: .375rem .5rem; background: var(--bg-tertiary); border: 1px solid var(--border); border-radius: var(--r-sm); font-family: var(--font-mono); font-size: .75rem; text-align: center; }
.del-warn { font-size: .8125rem; color: var(--text-primary); font-weight: 500; }
.pay-msg { font-size: .8125rem; color: var(--text-primary); padding: .5rem .75rem; background: var(--bg-tertiary); border: 1px solid var(--border); border-radius: var(--r); animation: fadeIn .2s var(--ease) both; }
.mt { margin-top: .75rem; }
.dialog-foot { display: flex; gap: .5rem; padding: 1rem 1.25rem; border-top: 1px solid var(--border); background: var(--bg-secondary); }
.btn-cancel { flex: 1; padding: .5625rem; font-size: .8125rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-cancel:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.btn-confirm { flex: 1; padding: .5625rem; font-size: .8125rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-confirm:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); }
.btn-confirm:disabled { opacity: .5; }
.btn-confirm.danger { background: var(--text-primary); color: var(--bg-primary); border-color: var(--text-primary); }

@media (max-width: 480px) { .plans-grid { grid-template-columns: 1fr; } }
</style>
