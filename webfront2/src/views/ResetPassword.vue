<template>
  <div class="reset-password-view">
    <!-- Hero Card -->
    <div class="hero-card">
      <div class="icon-wrapper">
        <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
        </svg>
      </div>
      <h1>🔐 Réinitialiser le mot de passe</h1>
      <h2>Créez un nouveau mot de passe sécurisé</h2>
    </div>

    <!-- Form Card -->
    <div v-if="!success && !error" class="form-card">
      <div class="form-section">
        <h3>Nouveau mot de passe</h3>
        
        <div class="form-group">
          <label>Mot de passe *</label>
          <div class="input-with-icon">
            <input 
              v-model="password" 
              :type="showPassword ? 'text' : 'password'" 
              class="input-field" 
              placeholder="Entrez votre nouveau mot de passe"
              @input="checkPasswordStrength"
              minlength="8"
              required 
            />
            <button @click="togglePassword" class="toggle-password" type="button">
              {{ showPassword ? '🙈' : '👁' }}
            </button>
          </div>
          
          <!-- Indicateur de force du mot de passe -->
          <div v-if="password" class="password-strength">
            <div class="strength-bar-container">
              <div 
                class="strength-bar" 
                :style="{ width: (passwordStrength.strength / 5 * 100) + '%', backgroundColor: strengthColor }"
              ></div>
            </div>
            <p class="strength-label" :style="{ color: strengthColor }">{{ passwordStrength.strengthLabel }}</p>
          </div>
          <ul v-if="password && passwordStrength.errors.length > 0" class="validation-errors">
            <li v-for="error in passwordStrength.errors" :key="error">{{ error }}</li>
          </ul>
        </div>

        <div class="form-group">
          <label>Confirmer le mot de passe *</label>
          <div class="input-with-icon">
            <input 
              v-model="confirmPassword" 
              :type="showConfirmPassword ? 'text' : 'password'" 
              class="input-field" 
              placeholder="Confirmez votre nouveau mot de passe"
              required 
            />
            <button @click="toggleConfirmPassword" class="toggle-password" type="button">
              {{ showConfirmPassword ? '🙈' : '👁' }}
            </button>
          </div>
          <p v-if="password !== confirmPassword && confirmPassword" class="error-message">
            ❌ Les mots de passe ne correspondent pas
          </p>
        </div>

        <button @click="resetPassword" class="submit-btn" :disabled="loading">
          <span v-if="!loading">Réinitialiser le mot de passe</span>
          <span v-else>⏳ Réinitialisation en cours...</span>
        </button>
      </div>
    </div>

    <!-- Success Message -->
    <div v-if="success" class="message-card success">
      <div class="icon-wrapper success">
        <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
      </div>
      <h2>✅ Mot de passe réinitialisé avec succès !</h2>
      <p>Vous pouvez maintenant vous connecter avec votre nouveau mot de passe.</p>
      <router-link to="/login" class="btn-link">Se connecter</router-link>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="message-card error">
      <div class="icon-wrapper error">
        <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
      </div>
      <h2>❌ Erreur</h2>
      <p>{{ errorMessage }}</p>
      <router-link to="/forgot" class="btn-link">Demander un nouveau lien</router-link>
    </div>
  </div>
</template>

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
  // Extraire le token de l'URL
  token.value = route.query.token || ''
  if (!token.value) {
    error.value = true
    errorMessage.value = 'Token de réinitialisation manquant. Veuillez demander un nouveau lien.'
  }
})

function togglePassword() {
  showPassword.value = !showPassword.value
}

function toggleConfirmPassword() {
  showConfirmPassword.value = !showConfirmPassword.value
}

function checkPasswordStrength() {
  if (password.value) {
    passwordStrength.value = validatePassword(password.value)
  }
}

async function resetPassword() {
  // Valider le mot de passe
  checkPasswordStrength()
  if (!passwordStrength.value.valid) {
    return
  }

  // Vérifier que les mots de passe correspondent
  if (password.value !== confirmPassword.value) {
    return
  }

  loading.value = true
  error.value = false
  errorMessage.value = ''

  try {
    const response = await fetch('http://localhost:8080/api/password/reset/confirm', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        token: token.value,
        new_password: password.value
      })
    })

    const data = await response.json()

    if (response.ok && data.success) {
      success.value = true
      // Rediriger vers la page de connexion après 3 secondes
      setTimeout(() => {
        router.push('/login')
      }, 3000)
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

<style scoped>
.reset-password-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2rem;
}

.hero-card, .form-card, .message-card {
  background: white;
  border-radius: 20px;
  padding: 3rem;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  max-width: 500px;
  width: 100%;
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.hero-card {
  text-align: center;
}

.icon-wrapper {
  width: 80px;
  height: 80px;
  margin: 0 auto 1.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
}

.icon-wrapper.success {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.icon-wrapper.error {
  background: linear-gradient(135deg, #eb3349 0%, #f45c43 100%);
}

.icon {
  width: 40px;
  height: 40px;
  color: white;
}

h1 {
  font-size: 2rem;
  color: #2d3748;
  margin-bottom: 0.5rem;
  font-weight: 700;
}

h2 {
  font-size: 1.1rem;
  color: #718096;
  font-weight: 400;
}

.form-section h3 {
  font-size: 1.3rem;
  color: #2d3748;
  margin-bottom: 1.5rem;
  font-weight: 600;
}

.form-group {
  margin-bottom: 1.5rem;
}

label {
  display: block;
  font-weight: 600;
  color: #4a5568;
  margin-bottom: 0.5rem;
  font-size: 0.95rem;
}

.input-with-icon {
  position: relative;
  display: flex;
  align-items: center;
}

.input-field {
  width: 100%;
  padding: 0.875rem 3rem 0.875rem 1rem;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  font-size: 1rem;
  transition: all 0.3s ease;
  background: #f7fafc;
}

.input-field:focus {
  outline: none;
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.toggle-password {
  position: absolute;
  right: 1rem;
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 0.5rem;
  transition: transform 0.2s;
}

.toggle-password:hover {
  transform: scale(1.1);
}

.password-strength {
  margin-top: 0.75rem;
}

.strength-bar-container {
  width: 100%;
  height: 8px;
  background-color: #e2e8f0;
  border-radius: 10px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.strength-bar {
  height: 100%;
  transition: width 0.3s ease, background-color 0.3s ease;
  border-radius: 10px;
}

.strength-label {
  font-size: 0.875rem;
  font-weight: 600;
  margin: 0;
}

.validation-errors {
  margin-top: 0.5rem;
  padding-left: 1.25rem;
  color: #e53e3e;
  font-size: 0.875rem;
}

.validation-errors li {
  margin-bottom: 0.25rem;
}

.error-message {
  color: #e53e3e;
  font-size: 0.875rem;
  margin-top: 0.5rem;
  font-weight: 500;
}

.submit-btn {
  width: 100%;
  padding: 1rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 1.5rem;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.message-card {
  text-align: center;
}

.message-card.success h2 {
  color: #11998e;
}

.message-card.error h2 {
  color: #e53e3e;
}

.message-card p {
  font-size: 1.1rem;
  color: #4a5568;
  margin: 1.5rem 0;
  line-height: 1.6;
}

.btn-link {
  display: inline-block;
  padding: 0.875rem 2rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  text-decoration: none;
  border-radius: 12px;
  font-weight: 600;
  transition: all 0.3s ease;
  margin-top: 1rem;
}

.btn-link:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
}

@media (max-width: 640px) {
  .reset-password-view {
    padding: 1rem;
  }

  .hero-card, .form-card, .message-card {
    padding: 2rem;
  }

  h1 {
    font-size: 1.5rem;
  }

  h2 {
    font-size: 1rem;
  }
}
</style>
