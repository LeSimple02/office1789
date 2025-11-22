<template>
  <div class="contact-view">
    <!-- Hero Card -->
    <div class="hero-card">
      <div class="contact-icon-wrapper">
        <svg class="contact-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
        </svg>
      </div>
      <h1>✉️ {{ $t('contactUs') }}</h1>
      <h2>{{ $t('hereToHelp') }}</h2>
    </div>

    <!-- Formulaire de contact -->
    <div class="contact-form-card">
      <form @submit.prevent="submitForm" class="contact-form">
        <div class="form-row">
          <div class="form-group">
            <label for="name">{{ $t('fullName') }}</label>
            <input 
              id="name" 
              v-model="formData.name" 
              type="text" 
              :placeholder="$t('yourName')" 
              required 
            />
          </div>
          
          <div class="form-group">
            <label for="email">{{ $t('email') }}</label>
            <input 
              id="email" 
              v-model="formData.email" 
              type="email" 
              :placeholder="$t('yourEmail')" 
              required 
            />
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label for="subject">{{ $t('subject') }}</label>
            <select id="subject" v-model="formData.subject" required>
              <option value="">{{ $t('selectSubject') }}</option>
              <option value="general">{{ $t('generalQuestion') }}</option>
              <option value="support">{{ $t('technicalSupport') }}</option>
              <option value="sales">{{ $t('sales') }}</option>
              <option value="feedback">{{ $t('feedback') }}</option>
            </select>
          </div>
        </div>

        <div class="form-group full-width">
          <label for="message">{{ $t('message') }}</label>
          <textarea 
            id="message" 
            v-model="formData.message" 
            :placeholder="$t('yourMessage')" 
            rows="6"
            required
          ></textarea>
        </div>

        <button type="submit" class="btn-submit">
          <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
          </svg>
          {{ $t('sendMessage') }}
        </button>

        <div v-if="submitStatus === 'success'" class="alert alert-success">
          ✅ {{ $t('messageSentSuccess') }}
        </div>
        <div v-if="submitStatus === 'error'" class="alert alert-error">
          ❌ {{ $t('error') }}
        </div>
      </form>
    </div>

    <!-- Info rapides -->
    <div class="info-cards">
      <div class="info-card">
        <div class="info-icon">📧</div>
        <h3>Email</h3>
        <p>contact@office1789.com</p>
      </div>
      
      <div class="info-card">
        <div class="info-icon">📞</div>
        <h3>Téléphone</h3>
        <p></p>
      </div>
      
      <div class="info-card">
        <div class="info-icon">🕐</div>
        <h3>Horaires</h3>
        <p>Lun-Ven : 9h-18h</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const formData = ref({
  name: '',
  email: '',
  phone: '',
  subject: '',
  message: ''
})

const submitStatus = ref(null)

function submitForm() {
  submitStatus.value = 'success'
  
  setTimeout(() => {
    submitStatus.value = null
    formData.value = {
      name: '',
      email: '',
      phone: '',
      subject: '',
      message: ''
    }
  }, 3000)
}
</script>

<style scoped>
.contact-view {
  width: 100%;
  animation: fadeIn 0.5s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.hero-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 280px;
  width: 100%;
  padding: 48px 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 32px;
  gap: 20px;
}

.dark .hero-card {
  background: #1C1C1E;
}

.contact-icon-wrapper {
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

.contact-icon {
  width: 40px;
  height: 40px;
  color: white;
}

.hero-card h1 {
  font-family: roboto, sans-serif;
  font-size: 2.5rem;
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
  font-size: 1.25rem;
  font-style: italic;
  text-align: center;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 0.5px;
  margin: 0;
}

.contact-form-card {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  padding: 40px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  margin-bottom: 32px;
}

.dark .contact-form-card {
  background: #1C1C1E;
}

.contact-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group.full-width {
  grid-column: 1 / -1;
}

.form-group label {
  font-family: roboto, sans-serif;
  font-size: 0.95rem;
  font-weight: 600;
  color: #333;
}

.dark .form-group label {
  color: #ddd;
}

.form-group input,
.form-group select,
.form-group textarea {
  padding: 12px 16px;
  border-radius: 12px;
  border: 1px solid #d0d0d0;
  font-size: 1rem;
  font-family: roboto, sans-serif;
  background: rgba(255,255,255,0.95);
  transition: all 0.2s;
}

.dark .form-group input,
.dark .form-group select,
.dark .form-group textarea {
  background: rgba(30,30,40,0.95);
  border-color: #444;
  color: #eee;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: blue;
  box-shadow: 0 0 0 3px rgba(0, 48, 143, 0.1);
}

.btn-submit {
  font-family: roboto, sans-serif;
  font-size: 1.1rem;
  padding: 14px 32px;
  border-radius: 12px;
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
  align-self: center;
  min-width: 240px;
}

.btn-submit:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.25);
}

.btn-icon {
  width: 20px;
  height: 20px;
}

.alert {
  padding: 16px;
  border-radius: 12px;
  text-align: center;
  font-weight: 500;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.alert-success {
  background: rgba(34, 197, 94, 0.15);
  color: #16a34a;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.alert-error {
  background: rgba(239, 68, 68, 0.15);
  color: #dc2626;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.info-cards {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-bottom: 32px;
}

.info-card {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 20px;
  padding: 32px 24px;
  text-align: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.dark .info-card {
  background: #1C1C1E;
}

.info-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.info-icon {
  font-size: 2.5rem;
  margin-bottom: 12px;
  animation: bounce 2s ease-in-out infinite;
}

.info-card:nth-child(2) .info-icon {
  animation-delay: 0.2s;
}

.info-card:nth-child(3) .info-icon {
  animation-delay: 0.4s;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.info-card h3 {
  font-family: roboto, sans-serif;
  font-size: 1.15rem;
  font-weight: 600;
  color: #222;
  margin: 0 0 8px 0;
}

.dark .info-card h3 {
  color: white;
}

.info-card p {
  font-family: roboto, sans-serif;
  font-size: 0.95rem;
  color: #666;
  margin: 0;
}

.dark .info-card p {
  color: #aaa;
}

@media (max-width: 768px) {
  .hero-card {
    padding: 32px 20px;
    min-height: 220px;
  }

  .hero-card h1 {
    font-size: 1.8rem;
  }

  .hero-card h2 {
    font-size: 1rem;
  }

  .contact-form-card {
    padding: 24px;
  }

  .form-row {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .info-cards {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .btn-submit {
    width: 100%;
  }
}
</style>
