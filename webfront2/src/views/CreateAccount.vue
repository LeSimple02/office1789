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

let emailCode = ref('')
let emailVerified = ref(false)
let emailVerificationSending = ref(false)
let emailVerificationChecking = ref(false)
let emailVerificationMessage = ref('')

let passwordStrength = ref({ valid: true, errors: [], strength: 0, strengthLabel: '' })
let emailValid = ref(true)
let usernameValid = ref({ valid: true, error: '' })
let phoneValid = ref(true)

function checkPasswordStrength() {
  if (passf1.value) {
    passwordStrength.value = validatePassword(passf1.value)
  }
}

function checkEmail() {
  if (email.value) {
    emailValid.value = isValidEmail(email.value)
  } else {
    emailValid.value = true
  }
}

function checkUsername() {
  if (username.value) {
    usernameValid.value = validateUsername(username.value)
  }
}

function checkPhone() {
  if (phonenumber.value && phonenumber.value.trim() !== '') {
    phoneValid.value = isValidPhone(phonenumber.value)
  } else {
    phoneValid.value = true
  }
}

const strengthColor = computed(() => getStrengthColor(passwordStrength.value.strength))

async function sendEmailVerification() {
  if (!email.value || !emailValid.value) return
  emailVerificationSending.value = true
  emailVerificationMessage.value = ''
  try {
    const res = await fetch(import.meta.env.VITE_APP_API_VERIFICATION_SEND, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ contact: email.value, type: 'email' })
    })
    const data = await res.json().catch(() => ({}))
    if (data.success) {
      emailVerificationMessage.value = 'Code envoyé à votre adresse email.'
    } else {
      emailVerificationMessage.value = data.message || "Erreur lors de l'envoi du code."
    }
  } catch (e) {
    console.error('sendEmailVerification error:', e)
    emailVerificationMessage.value = "Erreur réseau lors de l'envoi du code."
  } finally {
    emailVerificationSending.value = false
  }
}

async function verifyEmailCode() {
  if (!email.value || !emailCode.value) return
  emailVerificationChecking.value = true
  emailVerificationMessage.value = ''
  try {
    const res = await fetch(import.meta.env.VITE_APP_API_VERIFICATION_VERIFY, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ contact: email.value, code: emailCode.value, type: 'email' })
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
  checkUsername()
  checkEmail()
  checkPhone()
  checkPasswordStrength()
  if (!usernameValid.value.valid) return
  if (email.value && email.value.trim() !== '' && !emailValid.value) return
  if (phonenumber.value && phonenumber.value.trim() !== '' && !phoneValid.value) return
  if (!passwordStrength.value.valid) return
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
    if(v["Token"]){
      localStorage.setItem("log", 1)
      document.cookie = `name=${v["Username"]}; expires=${v["Expiry"]}; path=/; Secure; SameSite=Lax`
      document.cookie = `sessionToken=${v["Token"]}; expires=${v["Expiry"]}; path=/; Secure; SameSite=Lax`
      const store = gls()
      store.log = 1
      store.updateFromCookies()
      nextTick(() => {
        router.push("/mail")
      })
    }
    else if(v["username"] || v["phone"] || v["email"]) {
      usernameR.value = 0
      emailR.value = 0
      phonenumberR.value = 0
      if(v["username"] == "no") usernameR.value = 1
      if(v["email"] == "no") emailR.value = 1
      if(v["phone"] == "no") phonenumberR.value = 1
    }
  })
}

function show(){
  if (passw.value=="password") passw.value = "text"
  else if(passw.value=="text") passw.value="password"
}
function show2(){
  if (passw2.value=="password") passw2.value = "text"
  else if(passw2.value=="text") passw2.value="password"
}
</script>

<template>
  <div class="page-wrap">
    <div class="head">
      <h1>{{ $t('createac') }}</h1>
      <p>{{ $t('joinDigitalRevolution') }}</p>
    </div>

    <div class="card">
      <!-- Account section -->
      <div class="section">
        <h3 class="section-title">{{ $t('accountInformation') }}</h3>

        <div class="field">
          <label>{{ $t('usernameField') }}</label>
          <div class="input-suffix">
            <input v-model="username" type="text" minlength="3" maxlength="20" pattern="[a-zA-Z0-9_]+" @blur="checkUsername" required />
            <span class="suffix">@office1789.com</span>
          </div>
          <p v-if="!usernameValid.valid" class="err">{{ usernameValid.error }}</p>
          <p v-if="usernameR" class="err">{{ $t('dejaUP') }}</p>
        </div>

        <div class="field">
          <label>{{ $t('passwordField') }}</label>
          <div class="pwd-wrap">
            <input v-model="passf1" :type="passw" minlength="8" @input="checkPasswordStrength" required />
            <button type="button" @click="show" class="pwd-toggle">👁</button>
          </div>
          <div v-if="passf1" class="strength">
            <div class="strength-bar-bg">
              <div class="strength-bar" :style="{ width: (passwordStrength.strength / 5 * 100) + '%', background: strengthColor }"></div>
            </div>
            <p class="strength-label" :style="{ color: strengthColor }">{{ passwordStrength.strengthLabel }}</p>
          </div>
          <ul v-if="passf1 && passwordStrength.errors.length > 0" class="checklist">
            <li v-for="e in passwordStrength.errors" :key="e">{{ e }}</li>
          </ul>
        </div>

        <div class="field">
          <label>{{ $t('confirmPasswordField') }}</label>
          <div class="pwd-wrap">
            <input v-model="passf2" :type="passw2" required />
            <button type="button" @click="show2" class="pwd-toggle">👁</button>
          </div>
          <p v-if="passf1!=passf2 && passf2!=''" class="err">{{ $t('passwordd') }}</p>
        </div>
      </div>

      <!-- Personal section -->
      <div class="section">
        <h3 class="section-title">{{ $t('personalInformation') }}</h3>

        <div class="field">
          <label>{{ $t('emailField') }}</label>
          <input v-model="email" type="email" @blur="checkEmail" />
          <div v-if="email && emailValid" class="verify-row">
            <button type="button" class="btn-sec" @click="sendEmailVerification" :disabled="emailVerificationSending">
              <span v-if="emailVerificationSending" class="spinner"></span>
              {{ emailVerificationSending ? '...' : 'Envoyer le code' }}
            </button>
          </div>
          <div v-if="email" class="verify-row">
            <input v-model="emailCode" type="text" placeholder="Code reçu par email" />
            <button type="button" class="btn-sec" @click="verifyEmailCode" :disabled="emailVerificationChecking || !emailCode">
              {{ emailVerificationChecking ? '...' : 'Vérifier' }}
            </button>
          </div>
          <p v-if="emailVerificationMessage" class="verify-msg">{{ emailVerificationMessage }}</p>
          <p v-if="email && !emailValid" class="err">{{ $t('invalidEmailFormat') }}</p>
          <p v-if="emailR" class="err">{{ $t('dejaEP') }}</p>
        </div>
      </div>

      <button @click="verif" class="btn-submit">{{ $t('createac') }}</button>
    </div>

    <div class="login-link">
      <RouterLink to="/login" class="link">{{ $t('connection') }}</RouterLink>
    </div>
  </div>
</template>

<style scoped>
.page-wrap { max-width: 480px; margin: 0 auto; padding: 3rem 1.25rem 5rem; animation: fadeInUp .5s var(--ease) both; }
.head { text-align: center; margin-bottom: 1.75rem; animation: fadeInDown .4s var(--ease) both; }
.head h1 { font-size: 1.625rem; font-weight: 800; letter-spacing: -.025em; margin-bottom: .375rem; }
.head p { font-size: .875rem; color: var(--text-secondary); }
.card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 2rem 1.75rem; box-shadow: var(--sh-sm); animation: fadeInUp .4s var(--ease) .1s both; }
.section { margin-bottom: 1.5rem; }
.section-title { font-size: 1.0625rem; font-weight: 700; margin-bottom: 1.125rem; padding-bottom: .625rem; border-bottom: 1px solid var(--border); }
.field { margin-bottom: 1.125rem; display: flex; flex-direction: column; }
.field label { font-size: .8125rem; font-weight: 600; color: var(--text-secondary); margin-bottom: .375rem; }
.field input { width: 100%; padding: .625rem .875rem; font-size: .9375rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); transition: border-color var(--tb), box-shadow var(--tb); }
.field input:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.field input::placeholder { color: var(--text-muted); }
.input-suffix { display: flex; }
.input-suffix input { flex: 1; border-radius: var(--r) 0 0 var(--r); }
.suffix { display: flex; align-items: center; padding: 0 .875rem; font-size: .75rem; font-weight: 500; color: var(--text-muted); background: var(--bg-tertiary); border: 1px solid var(--border-strong); border-left: none; border-radius: 0 var(--r) var(--r) 0; white-space: nowrap; }
.pwd-wrap { position: relative; }
.pwd-wrap input { padding-right: 42px; }
.pwd-toggle { position: absolute; right: 6px; top: 50%; transform: translateY(-50%); background: none; border: none; padding: 4px 8px; font-size: 18px; cursor: pointer; opacity: .6; border-radius: var(--r-sm); transition: opacity var(--tb); }
.pwd-toggle:hover { opacity: 1; }
.strength { margin-top: .625rem; }
.strength-bar-bg { width: 100%; height: 4px; background: var(--bg-tertiary); border-radius: var(--r-full); overflow: hidden; }
.strength-bar { height: 100%; border-radius: var(--r-full); transition: width var(--tb); animation: growBar .4s var(--ease) both; }
.strength-label { font-size: .75rem; font-weight: 500; margin-top: .375rem; }
.checklist { list-style: none; padding: 0; margin: .5rem 0 0; }
.checklist li { font-size: .8125rem; color: var(--text-muted); margin-top: .25rem; animation: fadeIn .2s var(--ease) both; }
.err { font-size: .8125rem; color: var(--text-primary); margin-top: .5rem; font-weight: 500; }
.verify-row { display: flex; gap: .5rem; margin-top: .625rem; }
.verify-row input { flex: 1; }
.btn-sec { padding: .625rem 1rem; font-size: .8125rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; transition: all var(--tb); white-space: nowrap; display: flex; align-items: center; gap: .375rem; }
.btn-sec:hover:not(:disabled) { border-color: var(--ink); background: var(--bg-tertiary); }
.btn-sec:disabled { opacity: .5; cursor: not-allowed; }
.spinner { width: 14px; height: 14px; border: 2px solid var(--border); border-top-color: var(--ink); border-radius: 50%; animation: spin .6s linear infinite; display: inline-block; }
.verify-msg { font-size: .8125rem; color: var(--text-secondary); margin-top: .5rem; animation: fadeIn .2s var(--ease) both; }
.btn-submit { width: 100%; margin-top: .5rem; padding: .75rem; font-size: .9375rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-submit:hover { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-md); transform: translateY(-2px); }
.login-link { text-align: center; margin-top: 1.25rem; }
.link { font-size: .8125rem; font-weight: 500; color: var(--text-secondary); transition: color var(--tb); }
.link:hover { color: var(--text-primary); text-decoration: underline; text-underline-offset: 3px; }
</style>
