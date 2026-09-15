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
		if (a.require_2fa) {
			require2FA.value = true
			wrong.value = 0
			return
		}
		if (a.error) {
			wrong.value = 2
			totpCode.value = ''
			return
		}
		if (a["Username"] != "no") {
			localStorage.setItem("log", 1)
			document.cookie = `name=${a["Username"]}; expires=${a["Expiry"]}; path=/; Secure; SameSite=Lax`
			document.cookie = `sessionToken=${a["Token"]}; expires=${a["Expiry"]}; path=/; Secure; SameSite=Lax`
			const store = gls()
			store.log = 1
			store.updateFromCookies()
			nextTick(() => {
				router.push("/mail")
			})
		} else {
			wrong.value = 1
			require2FA.value = false
		}
	})
}

function show(){
	if (passw.value=="password")
		passw.value = "text"
	else if(passw.value=="text")
		passw.value="password"
}
</script>

<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-head">
        <img v-if="!isDark" src="@/assets/logo.png" class="login-logo" alt="Office1789" />
        <img v-if="isDark" src="@/assets/logol.png" class="login-logo" alt="Office1789" />
        <h1 class="login-title">{{ $t('connectionTitle') }}</h1>
        <p class="login-sub">{{ $t('accessCollaborativeSpace') }}</p>
      </div>

      <form class="login-form" @submit.prevent="connect">
        <div class="field">
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

        <div class="field">
          <label for="password">{{ $t('password') }}</label>
          <div class="pwd-wrap">
            <input
              v-model="passl"
              :type="passw"
              id="password"
              :placeholder="$t('yourPassword')"
              autocomplete="current-password"
              required
            />
            <button type="button" @click="show" class="pwd-toggle" :aria-label="$t('showPassword')">
              <svg v-if="passw === 'password'" width="18" height="18" fill="currentColor" viewBox="0 0 24 24"><path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"/></svg>
              <svg v-else width="18" height="18" fill="currentColor" viewBox="0 0 24 24"><path d="M12 7c2.76 0 5 2.24 5 5 0 .65-.13 1.26-.36 1.83l2.92 2.92c1.51-1.26 2.7-2.89 3.43-4.75-1.73-4.39-6-7.5-11-7.5-1.4 0-2.74.25-3.98.7l2.16 2.16C10.74 7.13 11.35 7 12 7zM2 4.27l2.28 2.28.46.46A11.804 11.804 0 0 0 1 12c1.73 4.39 6 7.5 11 7.5 1.55 0 3.03-.3 4.38-.84l.42.42L19.73 22 21 20.73 3.27 3 2 4.27z"/></svg>
            </button>
          </div>
        </div>

        <transition name="slide">
          <div v-if="require2FA" class="field totp">
            <label for="totp">{{ $t('twoFactorCode') }}</label>
            <input
              v-model="totpCode"
              type="text"
              id="totp"
              maxlength="8"
              placeholder="000000"
              autocomplete="one-time-code"
              required
              @input="totpCode = totpCode.replace(/[^0-9A-Z]/gi, '')"
            />
            <p class="totp-hint">{{ $t('enterSixDigitCode') }}</p>
          </div>
        </transition>

        <transition name="slide">
          <div v-if="wrong === 1" class="alert">⚠ {{$t('wrongL')}}</div>
          <div v-else-if="wrong === 2" class="alert">⚠ {{ $t('invalidTwoFactorCode') }}</div>
        </transition>

        <button type="submit" class="login-btn">{{ $t('connection') }}</button>
      </form>

      <div class="login-links">
        <RouterLink to="/forgot" class="link">{{ $t('forgot') }}</RouterLink>
        <RouterLink to="/createaccount" class="link">{{ $t('create') }}</RouterLink>
      </div>
    </div>

    <div class="login-tip">
      <strong>{{ $t('firstConnection') }}</strong> {{ $t('createAccountAndGet50GB') }}
    </div>
  </div>
</template>

<style scoped>
.login-page {
  max-width: 420px;
  margin: 0 auto;
  padding: 3rem 1.25rem 5rem;
  animation: fadeInUp .5s var(--ease) both;
}

/* ===== CARD ===== */
.login-card {
  background: var(--bg-primary);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  padding: 2.25rem 1.75rem;
  box-shadow: var(--sh-sm);
  animation: fadeInUp .4s var(--ease) .1s both;
}

/* ===== HEAD ===== */
.login-head {
  text-align: center;
  margin-bottom: 1.75rem;
  animation: fadeInDown .4s var(--ease) both;
}
.login-logo {
  width: 44px;
  height: 36px;
  border-radius: var(--r-sm);
  margin: 0 auto .875rem;
  transition: transform var(--tb);
}
.login-logo:hover { transform: scale(1.06); }
.login-title {
  font-size: 1.5rem;
  font-weight: 800;
  letter-spacing: -.025em;
  margin-bottom: .375rem;
}
.login-sub {
  font-size: .875rem;
  color: var(--text-secondary);
}

/* ===== FORM ===== */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.125rem;
}
.field {
  display: flex;
  flex-direction: column;
  animation: fadeIn .3s var(--ease) both;
}
.field label {
  font-size: .8125rem;
  font-weight: 600;
  color: var(--text-secondary);
  margin-bottom: .375rem;
  letter-spacing: .01em;
}
.field input {
  width: 100%;
  padding: .6875rem .875rem;
  font-size: .9375rem;
  color: var(--text-primary);
  background: var(--bg-primary);
  border: 1px solid var(--border-strong);
  border-radius: var(--r);
  transition: border-color var(--tb), box-shadow var(--tb);
}
.field input:focus {
  outline: none;
  border-color: var(--ink);
  box-shadow: 0 0 0 3px var(--bg-tertiary);
}
.field input::placeholder { color: var(--text-muted); }

/* ===== PASSWORD ===== */
.pwd-wrap { position: relative; }
.pwd-wrap input { padding-right: 42px; }
.pwd-toggle {
  position: absolute;
  right: 6px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  padding: 6px;
  color: var(--text-muted);
  border-radius: var(--r-sm);
  cursor: pointer;
  transition: color var(--tb);
}
.pwd-toggle:hover { color: var(--text-primary); }

/* ===== TOTP ===== */
.totp input {
  font-family: var(--font-mono);
  letter-spacing: .3em;
  font-size: 1.125rem;
  text-align: center;
}
.totp-hint {
  font-size: .75rem;
  color: var(--text-muted);
  margin-top: .375rem;
}

/* ===== ALERT ===== */
.alert {
  padding: .625rem .875rem;
  font-size: .8125rem;
  font-weight: 500;
  color: var(--text-primary);
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: var(--r);
  animation: fadeIn .2s var(--ease) both;
}

/* ===== BUTTON ===== */
.login-btn {
  width: 100%;
  padding: .75rem;
  margin-top: .5rem;
  font-size: .9375rem;
  font-weight: 600;
  color: var(--paper);
  background: var(--ink);
  border: 1px solid var(--ink);
  border-radius: var(--r);
  cursor: pointer;
  transition: all var(--tb);
  position: relative;
  overflow: hidden;
}
.login-btn:hover {
  background: var(--ink-soft);
  border-color: var(--ink-soft);
  box-shadow: var(--sh-md);
  transform: translateY(-2px);
}
.login-btn:active { transform: translateY(0); }

/* ===== LINKS ===== */
.login-links {
  display: flex;
  justify-content: space-between;
  margin-top: 1.25rem;
  gap: .75rem;
}
.link {
  font-size: .8125rem;
  font-weight: 500;
  color: var(--text-secondary);
  transition: color var(--tb);
}
.link:hover { color: var(--text-primary); text-decoration: underline; text-underline-offset: 3px; }

/* ===== TIP ===== */
.login-tip {
  margin-top: 1.25rem;
  padding: .875rem 1.125rem;
  font-size: .8125rem;
  color: var(--text-secondary);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--r);
  line-height: 1.5;
  animation: fadeInUp .4s var(--ease) .3s both;
}
.login-tip strong { color: var(--text-primary); font-weight: 600; }

/* ===== TRANSITIONS ===== */
.slide-enter-active { transition: all .3s var(--ease); }
.slide-leave-active { transition: all .2s var(--ease); }
.slide-enter-from { opacity: 0; transform: translateY(-8px); }
.slide-leave-to { opacity: 0; transform: translateY(-8px); }
</style>
