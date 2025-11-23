<script setup>
import {gls} from "@/stores/global"
import {ref, nextTick} from "vue"
import router from "@/router/index"
import { useDark } from '@vueuse/core'

let passw = ref("password")
let userl = ref('')
let passl = ref('')
let wrong = ref(0)
let require2FA = ref(false)
let totpCode = ref('')

const isDark = useDark()

function connect(){
	const payload = {
		username: userl.value, 
		password: passl.value
	}
	
	// Add TOTP code if in 2FA mode
	if (require2FA.value && totpCode.value) {
		payload.totp_code = totpCode.value
	}

	fetch(import.meta.env.VITE_APP_API_LOGIN, {
		method: "POST", 
		mode: "cors", 
		headers: { "Content-Type": "application/json"}, 
		body: JSON.stringify(payload)
	})
	.then(a => a.json())
	.then(a => {
		// Check if 2FA is required
		if (a.require_2fa) {
			require2FA.value = true
			wrong.value = 0
			return
		}
		
		// Check for invalid 2FA code
		if (a.error) {
			wrong.value = 2 // Special error for invalid 2FA code
			totpCode.value = ''
			return
		}
		
		// Successful login
		if (a["Username"] != "no") {
			localStorage.setItem("log", 1)
			
			// Définir les cookies avec path=/ pour qu'ils soient accessibles partout
			document.cookie = `name=${a["Username"]}; expires=${a["Expiry"]}; path=/; Secure; SameSite=Lax`
			document.cookie = `sessionToken=${a["Token"]}; expires=${a["Expiry"]}; path=/; Secure; SameSite=Lax`
			
			// Mettre à jour le store après avoir défini les cookies
			const store = gls()
			store.log = 1
			store.updateFromCookies()
			
			// Attendre le prochain tick pour garantir que le store est à jour
			nextTick(() => {
				router.push("/mail")
			})
		} else {
			wrong.value = 1
			require2FA.value = false
		}
	})
}

if (gls().log == 1){
	router.push("/mail")
}

function show(){
	if (passw.value=="password")
		passw.value = "text"
	else if(passw.value=="text")
		passw.value="password"
}
</script>

<template>
  <div class="login-container">
    <!-- Hero Card -->
    <div class="hero-card">
      <div class="login-icon-wrapper">
        <svg class="login-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"/>
        </svg>
      </div>
      <h1>🔐 {{ $t('connectionTitle') }}</h1>
      <h2>{{ $t('accessCollaborativeSpace') }}</h2>
    </div>

    <!-- Login Form Card -->
    <form class="login-card" @submit.prevent="connect">
      <div class="logo-section">
        <img v-if="!isDark" src="@/assets/logo.png" width="100" height="80" alt="Office1789 Logo"/>
        <img v-if="isDark" src="@/assets/logol.png" width="100" height="80" alt="Office1789 Logo"/>
        <p class="brand-name">Office1789</p>
      </div>

      <div class="form-inputs">
        <div class="input-group">
          <label for="username">{{ $t('identifier') }}</label>
          <input 
            v-model="userl" 
            type="text" 
            id="username"
            :placeholder="$t('yourIdentifier')" 
            autocomplete="username"
            required
          />
        </div>

        <div class="input-group">
          <label for="password">{{ $t('password') }}</label>
          <div class="password-wrapper">
            <input 
              v-model="passl" 
              :type="passw" 
              id="password"
              :placeholder="$t('yourPassword')" 
              autocomplete="current-password"
              required
            />
            <button type="button" @click="show" class="show-password" :aria-label="$t('showPassword')">
              <svg v-if="passw === 'password'" width="20" height="20" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"/>
              </svg>
              <svg v-else width="20" height="20" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 7c2.76 0 5 2.24 5 5 0 .65-.13 1.26-.36 1.83l2.92 2.92c1.51-1.26 2.7-2.89 3.43-4.75-1.73-4.39-6-7.5-11-7.5-1.4 0-2.74.25-3.98.7l2.16 2.16C10.74 7.13 11.35 7 12 7zM2 4.27l2.28 2.28.46.46A11.804 11.804 0 0 0 1 12c1.73 4.39 6 7.5 11 7.5 1.55 0 3.03-.3 4.38-.84l.42.42L19.73 22 21 20.73 3.27 3 2 4.27zM7.53 9.8l1.55 1.55c-.05.21-.08.43-.08.65 0 1.66 1.34 3 3 3 .22 0 .44-.03.65-.08l1.55 1.55c-.67.33-1.41.53-2.2.53-2.76 0-5-2.24-5-5 0-.79.2-1.53.53-2.2zm4.31-.78l3.15 3.15.02-.16c0-1.66-1.34-3-3-3l-.17.01z"/>
              </svg>
            </button>
          </div>
        </div>

        <!-- 2FA Code Input (shown only when required) -->
        <transition name="slide-fade">
          <div v-if="require2FA" class="input-group totp-group">
            <label for="totp">{{ $t('twoFactorCode') }}</label>
            <input 
              v-model="totpCode" 
              type="text" 
              id="totp"
              maxlength="8"
              placeholder="000000" 
              autocomplete="one-time-code"
              class="totp-input"
              required
              @input="totpCode = totpCode.replace(/[^0-9A-Z]/gi, '')"
            />
            <p class="totp-hint">{{ $t('enterSixDigitCode') }}</p>
          </div>
        </transition>

        <transition name="slide-fade">
          <div v-if="wrong === 1" class="error-message">
            ❌ {{$t('wrongL')}}
          </div>
          <div v-else-if="wrong === 2" class="error-message">
            ❌ {{ $t('invalidTwoFactorCode') }}
          </div>
        </transition>
      </div>

      <button type="submit" class="submit-btn">
        <svg width="20" height="20" fill="currentColor" viewBox="0 0 24 24">
          <path d="M10 17l5-5-5-5v10z"/>
        </svg>
        {{$t('connection')}}
      </button>

      <div class="links-section">
        <RouterLink to="/forgot" class="link">{{$t('forgot')}}</RouterLink>
        <RouterLink to="/createaccount" class="link">{{$t('create')}}</RouterLink>
      </div>
    </form>

    <!-- Info Card -->
    <div class="info-card">
      <p>💡 <strong>{{ $t('firstConnection') }}</strong> {{ $t('createAccountAndGet50GB') }}</p>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  width: 100%;
  max-width: 500px;
  margin: 0 auto;
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
  justify-content: center;
  min-height: 240px;
  width: 100%;
  padding: 40px 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 32px;
  gap: 16px;
}

.dark .hero-card {
  background: #1C1C1E;
}

.login-icon-wrapper {
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

.login-icon {
  width: 40px;
  height: 40px;
  color: white;
}

.hero-card h1 {
  font-family: 'Roboto', sans-serif;
  font-size: 2rem;
  font-weight: 700;
  text-align: center;
  letter-spacing: 0.5px;
  color: #222;
  margin: 0;
}

.dark .hero-card h1 {
  color: white;
}

.hero-card h2 {
  font-family: 'Roboto', sans-serif;
  font-size: 1.1rem;
  font-style: italic;
  text-align: center;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 0.3px;
  margin: 0;
}

/* Login Card */
.login-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
  width: 100%;
  padding: 40px 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 24px;
  animation: slideIn 0.6s ease-out;
}

@keyframes slideIn {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

.dark .login-card {
  background: #1C1C1E;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.logo-section img {
  animation: bounce 2s ease-in-out infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.brand-name {
  font-family: 'Roboto', sans-serif;
  font-size: 1.8rem;
  font-weight: 700;
  letter-spacing: 1.5px;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.form-inputs {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.input-group label {
  font-family: 'Roboto', sans-serif;
  font-size: 0.95rem;
  font-weight: 600;
  color: #333;
  margin-left: 4px;
}

.dark .input-group label {
  color: #ddd;
}

.input-group input {
  width: 100%;
  height: 48px;
  padding: 0 16px;
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-radius: 16px;
  font-family: 'Roboto', sans-serif;
  font-size: 1rem;
  background: white;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.dark .input-group input {
  background: #2C2C2E;
  border-color: rgba(255, 255, 255, 0.1);
  color: white;
}

.input-group input:focus {
  outline: none;
  border-color: blue;
  box-shadow: 0 0 0 3px rgba(0, 0, 255, 0.1);
  transform: translateY(-2px);
}

.password-wrapper {
  position: relative;
  width: 100%;
}

.password-wrapper input {
  padding-right: 48px;
}

.show-password {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.dark .show-password {
  color: #aaa;
}

.show-password:hover {
  background: rgba(0, 0, 0, 0.05);
  color: blue;
}

.dark .show-password:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.error-message {
  padding: 12px 16px;
  background: rgba(255, 60, 60, 0.1);
  border: 2px solid rgba(255, 60, 60, 0.3);
  border-radius: 12px;
  color: #ff3c3c;
  font-family: 'Roboto', sans-serif;
  font-size: 0.95rem;
  font-weight: 500;
  text-align: center;
}

.slide-fade-enter-active {
  transition: all 0.3s ease;
}

.slide-fade-leave-active {
  transition: all 0.3s ease;
}

.slide-fade-enter-from {
  transform: translateY(-10px);
  opacity: 0;
}

.slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}

.submit-btn {
  width: 100%;
  height: 56px;
  padding: 0 32px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border: none;
  border-radius: 20px;
  color: white;
  font-family: 'Roboto', sans-serif;
  font-size: 1.15rem;
  font-weight: 700;
  letter-spacing: 0.5px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.submit-btn:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.submit-btn:active {
  transform: translateY(-1px);
}

.links-section {
  display: flex;
  justify-content: center;
  gap: 24px;
  width: 100%;
  padding-top: 8px;
  border-top: 2px solid rgba(0, 0, 0, 0.08);
}

.dark .links-section {
  border-top-color: rgba(255, 255, 255, 0.08);
}

.link {
  font-family: 'Roboto', sans-serif;
  font-size: 0.95rem;
  font-weight: 500;
  color: blue;
  text-decoration: none;
  transition: all 0.3s ease;
  position: relative;
}

.dark .link {
  color: #6ec6ff;
}

.link::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: -webkit-linear-gradient(30deg, blue, red);
  transition: width 0.3s ease;
}

.link:hover {
  color: red;
  transform: translateY(-2px);
}

.link:hover::after {
  width: 100%;
}

/* Info Card */
.info-card {
  padding: 20px 24px;
  background: rgba(220, 240, 255, 0.85);
  border-radius: 20px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  border-left: 4px solid blue;
}

.dark .info-card {
  background: rgba(0, 50, 100, 0.2);
  border-left-color: #6ec6ff;
}

.info-card p {
  font-family: 'Roboto', sans-serif;
  font-size: 0.95rem;
  line-height: 1.6;
  color: #333;
  margin: 0;
}

.dark .info-card p {
  color: #ddd;
}

/* 2FA Input Styles */
.totp-group {
  background: rgba(59, 130, 246, 0.05);
  border: 1px solid #3b82f6;
  border-radius: 12px;
  padding: 16px;
  margin-top: 8px;
}

.dark .totp-group {
  background: rgba(59, 130, 246, 0.1);
  border-color: #60a5fa;
}

.totp-input {
  text-align: center;
  font-size: 24px;
  letter-spacing: 8px;
  font-family: 'Courier New', monospace;
  font-weight: bold;
  width: 100%;
  box-sizing: border-box;
  max-width: 100%;
}

.totp-hint {
  font-size: 12px;
  color: #6c757d;
  margin: 8px 0 0;
  text-align: center;
  line-height: 1.4;
}

.dark .totp-hint {
  color: #9ca3af;
}

/* Slide Fade Animation */
.slide-fade-enter-active {
  transition: all 0.3s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.2s ease-in;
}

.slide-fade-enter-from {
  transform: translateY(-10px);
  opacity: 0;
}

.slide-fade-leave-to {
  transform: translateY(-10px);
  opacity: 0;
}

/* Responsive */
@media (max-width: 767px) {
  .login-container {
    padding: 20px 12px;
  }

  .hero-card {
    padding: 32px 20px;
    min-height: 200px;
  }

  .hero-card h1 {
    font-size: 1.5rem;
  }

  .hero-card h2 {
    font-size: 0.95rem;
  }

  .login-card {
    padding: 32px 20px;
  }

  .logo-section img {
    width: 80px;
    height: 64px;
  }

  .brand-name {
    font-size: 1.4rem;
  }

  .submit-btn {
    height: 52px;
    font-size: 1rem;
  }

  .links-section {
    flex-direction: column;
    gap: 12px;
    align-items: center;
  }

  .info-card {
    padding: 16px 20px;
  }

  .info-card p {
    font-size: 0.9rem;
  }
}
</style>
