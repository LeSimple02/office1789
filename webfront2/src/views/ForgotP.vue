<template>
  <div class="forgot-view">
    <!-- Hero Card -->
    <div class="forgot-card">
      <div class="icon-wrapper">
        <svg class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"/>
        </svg>
      </div>
      <h1>🔑 {{$t('forgotp')}}</h1>
      <h2>{{ $t('recoverAccountAccess') }}</h2>

      <div class="form-content">
        <input 
          v-model="id" 
          type="text" 
          :placeholder="$t('identifierOrEmail')" 
          class="input-field"
          required 
        />
        
        <button @click="send" class="btn-submit" :disabled="loading">
          <svg v-if="!loading" class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
          </svg>
          <span v-if="!loading">{{ $t('sendEmail') }}</span>
          <span v-else>⏳ {{ $t('sendingInProgress') }}</span>
        </button>

        <p v-if="sendm=='nothing'" class="message error">❌ {{ $t('pleaseEnterIdentifier') }}</p>
        <p v-if="sendm==true" class="message success">✅ {{ $t('emailHasBeenSent') }}</p>

        <RouterLink to="/login" class="back-link">← {{ $t('backToLogin') }}</RouterLink>
      </div>
    </div>

    <!-- Info Card -->
    <div class="info-card">
      <svg class="info-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
      </svg>
      <p>{{ $t('forgotPasswordInfo') }}</p>
    </div>
  </div>
</template>

<script setup>
import {ref} from "vue"

let id = ref("")
let sendm = ref(false)
let loading = ref(false)

async function send(){
  if(id.value == "") {
    sendm.value = 'nothing'
    return
  }

  loading.value = true
  sendm.value = false

  try {
    const response = await fetch('http://localhost:8080/api/password/reset/request', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        identifier: id.value
      })
    })

    const data = await response.json()
    
    // Toujours afficher un message de succès pour des raisons de sécurité
    // (ne pas révéler si un compte existe ou non)
    if (response.ok || data.success) {
      sendm.value = true
    } else {
      // En cas d'erreur serveur, afficher quand même le message de succès
      sendm.value = true
    }
  } catch (error) {
    console.error('Error:', error)
    // Même en cas d'erreur, afficher le message de succès (sécurité)
    sendm.value = true
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.forgot-view {
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

/* Forgot Card */
.forgot-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  max-width: 480px;
  padding: 48px 40px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  gap: 24px;
}

.dark .forgot-card {
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

.forgot-card h1 {
  font-family: roboto, sans-serif;
  font-size: 2.2rem;
  font-weight: 700;
  text-align: center;
  letter-spacing: 1px;
  color: #222;
  margin: 0;
}

.dark .forgot-card h1 {
  color: white;
}

.forgot-card h2 {
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

.form-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
  width: 100%;
}

.input-field {
  width: 100%;
  padding: 16px 20px;
  border-radius: 16px;
  border: 2px solid rgba(0, 48, 143, 0.2);
  font-size: 1rem;
  font-family: roboto, sans-serif;
  background: white;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
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

.btn-submit {
  font-family: roboto, sans-serif;
  font-size: 1.1rem;
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
}

.btn-submit:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.btn-icon {
  width: 20px;
  height: 20px;
}

.message {
  font-family: roboto, sans-serif;
  font-size: 1rem;
  padding: 12px 16px;
  border-radius: 12px;
  text-align: center;
  font-weight: 500;
}

.message.error {
  background: rgba(255, 60, 60, 0.1);
  color: #ff3c3c;
  border: 2px solid rgba(255, 60, 60, 0.3);
}

.message.success {
  background: rgba(0, 200, 0, 0.1);
  color: #00c800;
  border: 2px solid rgba(0, 200, 0, 0.3);
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

/* Info Card */
.info-card {
  display: flex;
  align-items: center;
  gap: 16px;
  max-width: 600px;
  padding: 24px 28px;
  background: rgba(220, 220, 240, 0.95);
  border-radius: 24px;
  border-left: 4px solid blue;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.dark .info-card {
  background: rgba(0, 50, 0, 0.2);
  border-left-color: #06c;
}

.info-icon {
  width: 28px;
  height: 28px;
  color: #06c;
  flex-shrink: 0;
}

.info-card p {
  font-family: roboto, sans-serif;
  font-size: 0.95rem;
  color: #222;
  margin: 0;
  line-height: 1.6;
}

.dark .info-card p {
  color: white;
}

/* Responsive */
@media (max-width: 767px) {
  .forgot-card {
    padding: 32px 24px;
    border-radius: 20px;
  }

  .forgot-card h1 {
    font-size: 1.6rem;
  }

  .forgot-card h2 {
    font-size: 0.95rem;
  }

  .btn-submit {
    font-size: 1rem;
    padding: 14px 24px;
  }

  .info-card {
    flex-direction: column;
    text-align: center;
  }
}
</style>
