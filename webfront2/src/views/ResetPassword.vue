<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { validatePassword, getStrengthColor } from '@/utils/validation'

const route = useRoute()
const router = useRouter()

const token = ref('')
const password = ref('')
const confirmPassword = ref('')
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const passwordStrength = ref({ valid: true, errors: [], strength: 0, strengthLabel: '' })
const loading = ref(false)
const success = ref(false)
const error = ref(false)
const errorMessage = ref('')

const strengthColor = computed(() => getStrengthColor(passwordStrength.value.strength))

onMounted(() => {
  token.value = route.query.token || ''
  if (!token.value) {
    error.value = true
    errorMessage.value = 'Token de réinitialisation manquant. Veuillez demander un nouveau lien.'
  }
})

function togglePassword() { showPassword.value = !showPassword.value }
function toggleConfirmPassword() { showConfirmPassword.value = !showConfirmPassword.value }
function checkPasswordStrength() {
  if (password.value) {
    passwordStrength.value = validatePassword(password.value)
  }
}

async function resetPassword() {
  checkPasswordStrength()
  if (!passwordStrength.value.valid) return
  if (password.value !== confirmPassword.value) return
  loading.value = true
  error.value = false
  errorMessage.value = ''
  try {
    const response = await fetch(import.meta.env.VITE_APP_API_PASSWORD_RESET_CONFIRM, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ token: token.value, new_password: password.value })
    })
    const data = await response.json()
    if (response.ok && data.success) {
      success.value = true
      setTimeout(() => { router.push('/login') }, 3000)
    } else {
      error.value = true
      errorMessage.value = data.error || 'Le lien de réinitialisation est invalide ou a expiré.'
    }
  } catch (err) {
    error.value = true
    errorMessage.value = 'Une erreur est survenue. Veuillez réessayer.'
    console.error('Reset password error:', err)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page-wrap">
    <!-- Form state -->
    <div v-if="!success && !error" class="card">
      <div class="head">
        <h1>{{ $t('resetPasswordTitle') }}</h1>
        <p>{{ $t('createNewSecurePassword') }}</p>
      </div>

      <div class="form">
        <div class="field">
          <label>{{ $t('password') }} *</label>
          <div class="pwd-wrap">
            <input v-model="password" :type="showPassword ? 'text' : 'password'" :placeholder="$t('enterNewPassword')" @input="checkPasswordStrength" minlength="8" required />
            <button type="button" @click="togglePassword" class="pwd-toggle">{{ showPassword ? '🙈' : '👁' }}</button>
          </div>
          <div v-if="password" class="strength">
            <div class="strength-bar-bg">
              <div class="strength-bar" :style="{ width: (passwordStrength.strength / 5 * 100) + '%', background: strengthColor }"></div>
            </div>
            <p class="strength-label" :style="{ color: strengthColor }">{{ passwordStrength.strengthLabel }}</p>
          </div>
          <ul v-if="password && passwordStrength.errors.length > 0" class="checklist">
            <li v-for="e in passwordStrength.errors" :key="e">{{ e }}</li>
          </ul>
        </div>

        <div class="field">
          <label>{{ $t('confirmPassword') }} *</label>
          <div class="pwd-wrap">
            <input v-model="confirmPassword" :type="showConfirmPassword ? 'text' : 'password'" :placeholder="$t('confirmNewPassword')" required />
            <button type="button" @click="toggleConfirmPassword" class="pwd-toggle">{{ showConfirmPassword ? '🙈' : '👁' }}</button>
          </div>
          <p v-if="password !== confirmPassword && confirmPassword" class="err">{{ $t('passwordsDoNotMatch') }}</p>
        </div>

        <button @click="resetPassword" class="btn-go" :disabled="loading">
          <span v-if="!loading">{{ $t('resetPassword') }}</span>
          <span v-else class="loader"></span>
        </button>
      </div>
    </div>

    <!-- Success -->
    <div v-if="success" class="card result">
      <h2>{{ $t('passwordResetSuccess') }}</h2>
      <p>{{ $t('passwordResetSuccessMessage') }}</p>
      <RouterLink to="/login" class="link-btn">{{ $t('backToLogin') }}</RouterLink>
    </div>

    <!-- Error -->
    <div v-if="error" class="card result">
      <h2>Erreur</h2>
      <p>{{ errorMessage }}</p>
      <RouterLink to="/forgot" class="link-btn">{{ $t('forgotp') }}</RouterLink>
    </div>
  </div>
</template>

<style scoped>
.page-wrap { max-width: 440px; margin: 0 auto; padding: 3rem 1.25rem 5rem; animation: fadeInUp .5s var(--ease) both; }
.card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); box-shadow: var(--sh-sm); overflow: hidden; animation: fadeInUp .4s var(--ease) .1s both; }
.head { text-align: center; padding: 2rem 1.75rem 1rem; animation: fadeInDown .4s var(--ease) both; }
.head h1 { font-size: 1.5rem; font-weight: 800; letter-spacing: -.025em; margin-bottom: .375rem; }
.head p { font-size: .875rem; color: var(--text-secondary); }
.form { padding: 0 1.75rem 1.75rem; display: flex; flex-direction: column; gap: 1.125rem; }
.field { display: flex; flex-direction: column; }
.field label { font-size: .8125rem; font-weight: 600; color: var(--text-secondary); margin-bottom: .375rem; }
.pwd-wrap { position: relative; }
.pwd-wrap input { width: 100%; padding: .6875rem 42px .6875rem .875rem; font-size: .9375rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); transition: border-color var(--tb), box-shadow var(--tb); }
.pwd-wrap input:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.pwd-wrap input::placeholder { color: var(--text-muted); }
.pwd-toggle { position: absolute; right: 6px; top: 50%; transform: translateY(-50%); background: none; border: none; padding: 4px 8px; font-size: 18px; cursor: pointer; opacity: .6; border-radius: var(--r-sm); transition: opacity var(--tb); }
.pwd-toggle:hover { opacity: 1; }
.strength { margin-top: .625rem; }
.strength-bar-bg { width: 100%; height: 4px; background: var(--bg-tertiary); border-radius: var(--r-full); overflow: hidden; }
.strength-bar { height: 100%; border-radius: var(--r-full); transition: width var(--tb); animation: growBar .4s var(--ease) both; }
.strength-label { font-size: .75rem; font-weight: 500; margin-top: .375rem; }
.checklist { list-style: none; padding: 0; margin: .5rem 0 0; }
.checklist li { font-size: .8125rem; color: var(--text-muted); margin-top: .25rem; }
.err { font-size: .8125rem; color: var(--text-primary); margin-top: .5rem; font-weight: 500; }
.btn-go { width: 100%; padding: .75rem; font-size: .9375rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); display: flex; align-items: center; justify-content: center; min-height: 44px; }
.btn-go:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-md); transform: translateY(-2px); }
.btn-go:disabled { opacity: .5; cursor: not-allowed; }
.loader { width: 18px; height: 18px; border: 2.5px solid var(--border); border-top-color: var(--paper); border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
.result { text-align: center; padding: 2.5rem 1.75rem; animation: popIn .3s var(--ease-spring) both; }
.result h2 { font-size: 1.25rem; font-weight: 700; margin-bottom: .5rem; }
.result p { font-size: .9375rem; color: var(--text-secondary); margin-bottom: 1.25rem; }
.link-btn { display: inline-block; padding: .625rem 1.5rem; font-size: .9375rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); transition: all var(--tb); }
.link-btn:hover { background: var(--ink-soft); border-color: var(--ink-soft); transform: translateY(-2px); box-shadow: var(--sh-sm); }
</style>
