<template>
  <div class="create-view">
    <!-- Hero Card -->
    <div class="hero-card">
      <div class="icon-wrapper">
        <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z"/>
        </svg>
      </div>
      <h1>✨ {{$t("createac")}}</h1>
      <h2>{{ $t('joinDigitalRevolution') }}</h2>
    </div>

    <!-- Form Card -->
    <div class="form-card">
      <div class="form-section">
        <h3>{{ $t('accountInformation') }}</h3>
        
        <div class="form-group">
          <label>{{ $t('usernameField') }}</label>
          <div class="input-with-suffix">
            <input 
              v-model="username" 
              type="text" 
              class="input-field" 
              minlength="3"
              maxlength="20"
              pattern="[a-zA-Z0-9_]+"
              @blur="checkUsername"
              required 
            />
            <span class="suffix">@office1789.com</span>
          </div>
          <p v-if="!usernameValid.valid" class="error-message">❌ {{ usernameValid.error }}</p>
          <p v-if="usernameR" class="error-message">❌ {{$t('dejaUP')}}</p>
        </div>

        <div class="form-group">
          <label>{{ $t('passwordField') }}</label>
          <div class="input-with-icon">
            <input 
              v-model="passf1" 
              :type="passw" 
              class="input-field" 
              minlength="8"
              @input="checkPasswordStrength"
              required 
            />
            <button @click="show()" class="toggle-password" type="button">👁</button>
          </div>
          <!-- Indicateur de force du mot de passe -->
          <div v-if="passf1" class="password-strength">
            <div class="strength-bar-container">
              <div 
                class="strength-bar" 
                :style="{ width: (passwordStrength.strength / 5 * 100) + '%', backgroundColor: strengthColor }"
              ></div>
            </div>
            <p class="strength-label" :style="{ color: strengthColor }">{{ passwordStrength.strengthLabel }}</p>
          </div>
          <ul v-if="passf1 && passwordStrength.errors.length > 0" class="validation-errors">
            <li v-for="error in passwordStrength.errors" :key="error">{{ error }}</li>
          </ul>
        </div>

        <div class="form-group">
          <label>{{ $t('confirmPasswordField') }}</label>
          <div class="input-with-icon">
            <input v-model="passf2" :type="passw2" class="input-field" required />
            <button @click="show2()" class="toggle-password">👁</button>
          </div>
          <p v-if="passf1!=passf2 && passf2!=''" class="error-message">❌ {{$t('passwordd')}}</p>
        </div>
      </div>

      <div class="form-section">
        <h3>{{ $t('personalInformation') }}</h3>
        
        <div class="form-group">
          <label>{{ $t('emailField') }}</label>
          <input 
            v-model="email" 
            type="email" 
            class="input-field"
            @blur="checkEmail"
          />
          <!-- Vérification email : envoi du code et saisie -->
          <div v-if="email && emailValid" class="verification-actions">
            <button 
              type="button" 
              class="btn-submit secondary" 
              @click="sendEmailVerification" 
              :disabled="emailVerificationSending"
            >
              <span v-if="emailVerificationSending" class="btn-spinner"></span>
              {{ $t('sendVerificationCode') || 'Envoyer un code de vérification' }}
            </button>
          </div>

          <div v-if="email" class="verification-code">
            <input 
              v-model="emailCode" 
              type="text" 
              class="input-field" 
              :placeholder="$t('enterVerificationCode') || 'Code reçu par email'" 
            />
            <button 
              type="button" 
              class="btn-submit secondary" 
              @click="verifyEmailCode" 
              :disabled="emailVerificationChecking || !emailCode"
            >
              <span v-if="emailVerificationChecking" class="btn-spinner"></span>
              {{ $t('verifyCode') || 'Vérifier le code' }}
            </button>
          </div>
          <p v-if="email && !emailValid" class="error-message">❌ {{ $t('invalidEmailFormat') }}</p>
          <p v-if="emailR" class="error-message">❌ {{$t('dejaEP')}}</p>
        </div>

        <div class="form-group" v-if="false">
          <label>{{ $t('phoneNumberField') }}</label>
          <input 
            v-model="phonenumber" 
            type="tel" 
            class="input-field"
            :placeholder="$t('phonePlaceholder')"
            @blur="checkPhone"
          />
          <p v-if="phonenumber && !phoneValid" class="error-message">❌ {{ $t('invalidPhoneFormat') }}</p>
          <p v-if="phonenumberR" class="error-message">❌ {{$t('dejaPP')}}</p>
        </div>
      </div>

      <button @click="verif()" class="btn-submit">
        <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        {{ $t('createMyAccount') }}
      </button>

      <RouterLink to="/login" class="back-link">← {{ $t('alreadyHaveAccount') }}</RouterLink>
    </div>
  </div>
</template>

<script setup>
import {ref, computed, nextTick} from "vue"
import {gls} from "@/stores/global"
import router from "@/router/index"
import { validatePassword, isValidEmail, validateUsername, isValidPhone, getStrengthColor } from '@/utils/validation'

let passw = ref("password")
let passw2 = ref("password")

let username = ref('')
let passf1 = ref('')
let passf2 = ref('')
let email = ref('')
let phonenumber = ref('')

let usernameR = ref(0)
let emailR = ref(0)
let phonenumberR = ref(0)

// Vérification email (code envoyé par /api/verification/send)
let emailCode = ref('')
let emailVerified = ref(false)
let emailVerificationSending = ref(false)
let emailVerificationChecking = ref(false)
let emailVerificationMessage = ref('')

// Validation states
let passwordStrength = ref({ valid: true, errors: [], strength: 0, strengthLabel: '' })
let emailValid = ref(true)
let usernameValid = ref({ valid: true, error: '' })
let phoneValid = ref(true)

// Check password strength
function checkPasswordStrength() {
  if (passf1.value) {
    passwordStrength.value = validatePassword(passf1.value)
  }
}

// Validate email format
function checkEmail() {
  if (email.value) {
    emailValid.value = isValidEmail(email.value)
  } else {
    emailValid.value = true // Optional field
  }
}

// Validate username format
function checkUsername() {
  if (username.value) {
    usernameValid.value = validateUsername(username.value)
  }
}

// Validate phone format
function checkPhone() {
  if (phonenumber.value && phonenumber.value.trim() !== '') {
    phoneValid.value = isValidPhone(phonenumber.value)
  } else {
    phoneValid.value = true // Optional field
  }
}

// Computed password strength color
const strengthColor = computed(() => getStrengthColor(passwordStrength.value.strength))

// Envoi du code de vérification par email
async function sendEmailVerification() {
  if (!email.value || !emailValid.value) return

  emailVerificationSending.value = true
  emailVerificationMessage.value = ''

  try {
    const res = await fetch(import.meta.env.VITE_APP_API_VERIFICATION_SEND, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        contact: email.value,
        type: 'email'
      })
    })

    const data = await res.json().catch(() => ({}))
    if (data.success) {
      emailVerificationMessage.value = 'Code envoyé à votre adresse email.'
    } else {
      emailVerificationMessage.value = data.message || 'Erreur lors de l’envoi du code.'
    }
  } catch (e) {
    console.error('sendEmailVerification error:', e)
    emailVerificationMessage.value = 'Erreur réseau lors de l’envoi du code.'
  } finally {
    emailVerificationSending.value = false
  }
}

// Vérification du code reçu par email
async function verifyEmailCode() {
  if (!email.value || !emailCode.value) return

  emailVerificationChecking.value = true
  emailVerificationMessage.value = ''

  try {
    const res = await fetch(import.meta.env.VITE_APP_API_VERIFICATION_VERIFY, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        contact: email.value,
        code: emailCode.value,
        type: 'email'
      })
    })

    const data = await res.json().catch(() => ({}))
    if (data.success) {
      emailVerified.value = true
      emailVerificationMessage.value = 'Email vérifié.'
    } else {
      emailVerified.value = false
      emailVerificationMessage.value = data.message || 'Code invalide ou expiré.'
    }
  } catch (e) {
    console.error('verifyEmailCode error:', e)
    emailVerificationMessage.value = 'Erreur réseau lors de la vérification.'
  } finally {
    emailVerificationChecking.value = false
  }
}

function verif(){
  // Valider tous les champs
  checkUsername()
  checkEmail()
  checkPhone()
  checkPasswordStrength()
  
  // Vérifier username valide
  if (!usernameValid.value.valid) {
    return
  }
  
  // Vérifier email valide si rempli
  if (email.value && email.value.trim() !== '' && !emailValid.value) {
    return
  }
  
  // Vérifier téléphone valide si rempli
  if (phonenumber.value && phonenumber.value.trim() !== '' && !phoneValid.value) {
    return
  }
  
  // Vérifier force mot de passe
  if (!passwordStrength.value.valid) {
    return
  }

  // Si un email est saisi, on s'attend à ce qu'il ait été vérifié
  // (le backend vérifiera aussi et renverra email_not_verified si besoin)
  
  // Vérifier correspondance mots de passe
  if(passf1.value == passf2.value && passf2.value !="" && passf1.value != "")
    connect()
}

function connect(){  
  let d = new Date()
  let dc = `${d.getFullYear()}-${(d.getMonth()+1).toString().padStart(2, '0')}-${d.getDate().toString().padStart(2, '0')} ${d.getHours().toString().padStart(2, '0')}:${d.getMinutes().toString().padStart(2, '0')}:${d.getSeconds().toString().padStart(2, '0')}`

  fetch(import.meta.env.VITE_APP_API_CREATE_ACCOUNT, { 
    method: "POST", 
    mode: "cors", 
    credentials: "same-origin", 
    headers: { "Content-Type": "application/json"}, 
    body : JSON.stringify({
      "username" : username.value, 
      "password": passf1.value, 
      "email": email.value, 
      "phonenumber": phonenumber.value, 
      "datejoined": dc, 
      "lastlogin": dc
    }) 
  }).then((v)=>{return v.json()}).then(
  (v)=>{
    // Si on a un Token, c'est que la création a réussi
    if(v["Token"]){
      localStorage.setItem("log", 1)
      
      // Définir les cookies avec path=/ pour qu'ils soient accessibles partout
      document.cookie = `name=${v["Username"]}; expires=${v["Expiry"]}; path=/; Secure; SameSite=Lax`
      document.cookie = `sessionToken=${v["Token"]}; expires=${v["Expiry"]}; path=/; Secure; SameSite=Lax`
      
      // Mettre à jour le store après avoir défini les cookies
      const store = gls()
      store.log = 1
      store.updateFromCookies()
      
      // Attendre le prochain tick pour garantir que le store est à jour
      nextTick(() => {
        console.log("REDIRECTION !")
        router.push("/mail")
      })
    }
    // Erreur de validation (champs déjà pris)
    else if(v["username"] || v["phone"] || v["email"]) {
      // Reset les erreurs
      usernameR.value = 0
      emailR.value = 0
      phonenumberR.value = 0
      
      // Afficher les erreurs spécifiques
      if(v["username"] == "no"){
        usernameR.value = 1      
      }
      if(v["email"] == "no"){
        emailR.value = 1 
      }
      if(v["phone"] == "no"){
        phonenumberR.value = 1 
      }
    }
  })
}

function show(){
  if (passw.value=="password")
    passw.value = "text"
  else if(passw.value=="text")
    passw.value="password"
}

function show2(){
  if (passw2.value=="password")
    passw2.value = "text"
  else if(passw2.value=="text")
    passw2.value="password"
}
</script>

<style scoped>
.create-view {
  min-height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 32px;
  padding: 32px 16px;
  animation: fadeIn 0.5s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Hero Card */
.hero-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  max-width: 700px;
  padding: 40px 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  gap: 20px;
}

.dark .hero-card {
  background: #1C1C1E;
}

.icon-wrapper {
  width: 80px;
  height: 80px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.icon {
  width: 40px;
  height: 40px;
  color: white;
}

.hero-card h1 {
  font-family: roboto, sans-serif;
  font-size: 2.2rem;
  font-weight: 700;
  text-align: center;
  letter-spacing: 1px;
  color: #222;
  margin: 0;
}

.dark .hero-card h1 {
  color: white;
}

.hero-card h2 {
  font-family: roboto, sans-serif;
  font-size: 1.1rem;
  font-style: italic;
  text-align: center;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 0.5px;
  margin: 0;
}

/* Form Card */
.form-card {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: 700px;
  padding: 40px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  gap: 32px;
}

.dark .form-card {
  background: #1C1C1E;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-section h3 {
  font-family: roboto, sans-serif;
  font-size: 1.3rem;
  font-weight: 600;
  color: #222;
  margin: 0;
  padding-bottom: 12px;
  border-bottom: 2px solid rgba(0, 48, 143, 0.2);
}

.dark .form-section h3 {
  color: white;
  border-bottom-color: rgba(255, 255, 255, 0.2);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-family: roboto, sans-serif;
  font-size: 0.95rem;
  font-weight: 500;
  color: #333;
}

.dark .form-group label {
  color: #ddd;
}

.input-field {
  width: 100%;
  padding: 14px 18px;
  border-radius: 12px;
  border: 2px solid rgba(0, 48, 143, 0.2);
  font-size: 1rem;
  font-family: roboto, sans-serif;
  background: white;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.input-field:focus {
  outline: none;
  border-color: blue;
  box-shadow: 0 4px 16px rgba(0, 48, 143, 0.15);
}

.dark .input-field {
  background: rgba(30, 30, 40, 0.95);
  color: white;
  border-color: rgba(255, 255, 255, 0.2);
}

.input-with-suffix {
  display: flex;
  align-items: center;
  gap: 8px;
}

.suffix {
  font-family: roboto, sans-serif;
  font-size: 1rem;
  color: #666;
  white-space: nowrap;
}

.dark .suffix {
  color: #aaa;
}

.input-with-icon {
  position: relative;
  display: flex;
  align-items: center;
}

.input-with-icon .input-field {
  padding-right: 50px;
}

.toggle-password {
  position: absolute;
  right: 12px;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 4px 8px;
}

.error-message {
  font-family: roboto, sans-serif;
  font-size: 0.9rem;
  color: #ff3c3c;
  margin: 0;
}

/* Password Strength Indicator */
.password-strength {
  margin-top: 8px;
}

.strength-bar-container {
  width: 100%;
  height: 6px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 4px;
}

.dark .strength-bar-container {
  background: rgba(255, 255, 255, 0.1);
}

.strength-bar {
  height: 100%;
  transition: all 0.3s ease;
  border-radius: 3px;
}

.strength-label {
  font-size: 0.85rem;
  font-weight: 600;
  margin: 0;
  text-align: right;
}

.validation-errors {
  list-style: none;
  padding: 8px 12px;
  margin: 8px 0 0 0;
  background: rgba(255, 60, 60, 0.05);
  border-left: 3px solid #ff3c3c;
  border-radius: 4px;
}

.validation-errors li {
  font-size: 0.85rem;
  color: #ff3c3c;
  margin: 4px 0;
}

.validation-errors li::before {
  content: '• ';
  font-weight: bold;
}

.btn-submit {
  font-family: roboto, sans-serif;
  font-size: 1.2rem;
  padding: 16px 32px;
  border-radius: 16px;
  border: none;
  cursor: pointer;
  background: -webkit-linear-gradient(30deg, blue, red);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  font-weight: 600;
  margin-top: 8px;
}

.btn-submit:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.btn-icon {
  width: 20px;
  height: 20px;
}

.back-link {
  font-family: roboto, sans-serif;
  font-size: 1rem;
  color: #666;
  text-decoration: none;
  text-align: center;
  transition: color 0.3s ease;
}

.back-link:hover {
  color: blue;
}

.dark .back-link {
  color: #aaa;
}

.dark .back-link:hover {
  color: red;
}

/* Responsive */
@media (max-width: 767px) {
  .hero-card, .form-card {
    padding: 24px 20px;
    border-radius: 20px;
  }

  .hero-card h1 {
    font-size: 1.6rem;
  }

  .hero-card h2 {
    font-size: 0.95rem;
  }

  .form-section {
    gap: 16px;
  }

  .btn-submit {
    font-size: 1rem;
    padding: 14px 24px;
  }

  .input-with-suffix {
    flex-direction: column;
    align-items: flex-start;
  }

  .suffix {
    font-size: 0.9rem;
  }
}
</style>
