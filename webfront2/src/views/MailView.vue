<template>
  <div class="mail-view">
    <!-- Hero Card -->
    <div class="hero-card">
      <div class="mail-icon-wrapper">
        <svg class="mail-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
        </svg>
      </div>
      <h1>📧 Messagerie Office1789</h1>
      <h2>Vos emails, simplement et en toute sécurité</h2>
      
      <button @click="openMailInNewWindow" class="btn-open-mail">
        <svg class="btn-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"/>
        </svg>
        Ouvrir la Messagerie
      </button>
    </div>

    <!-- Features Section -->
    <div class="features-section">
      <div class="feature-card">
        <div class="feature-icon">📬</div>
        <h3>Sécurisé</h3>
        <p>Vos emails sont protégés avec le chiffrement de bout en bout</p>
      </div>
      
      <div class="feature-card">
        <div class="feature-icon">⚡</div>
        <h3>Rapide</h3>
        <p>Interface ultra-réactive pour une gestion efficace de vos messages</p>
      </div>
      
      <div class="feature-card">
        <div class="feature-icon">🔒</div>
        <h3>Privé</h3>
        <p>Vos données restent privées, aucun tracking publicitaire</p>
      </div>
    </div>

    <!-- Info Box -->
    <div class="info-card">
      <svg class="info-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
      </svg>
      <p>La messagerie s'ouvrira dans une nouvelle fenêtre pour une meilleure expérience utilisateur</p>
    </div>
  </div>
</template>

<script setup>
import { gls } from "@/stores/global"
import router from "@/router/index"

if (gls().log != 1) {
  router.push("login")
}

const mailUrl = "http://localhost:8081"

function openMailInNewWindow() {
  window.open(mailUrl, '_blank')
}
</script>

<style scoped>
.mail-view {
  width: 100%;
  animation: fadeIn 0.5s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* Hero Card */
.hero-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
  width: 100%;
  padding: 48px 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 32px;
  gap: 24px;
}

.dark .hero-card {
  background: #1C1C1E;
}

.mail-icon-wrapper {
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

.mail-icon {
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

.btn-open-mail {
  font-family: roboto, sans-serif;
  font-size: 1.2rem;
  padding: 1rem 2.5rem;
  border-radius: 32px;
  border: none;
  cursor: pointer;
  background: -webkit-linear-gradient(30deg, blue, red);
  color: white;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  font-weight: 600;
}

.btn-open-mail:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.btn-icon {
  width: 24px;
  height: 24px;
}

/* Features Section */
.features-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-bottom: 32px;
}

.feature-card {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  padding: 32px 24px;
  text-align: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.dark .feature-card {
  background: #1C1C1E;
}

.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.feature-icon {
  font-size: 3rem;
  margin-bottom: 16px;
  animation: bounce 2s ease-in-out infinite;
}

.feature-card:nth-child(2) .feature-icon {
  animation-delay: 0.2s;
}

.feature-card:nth-child(3) .feature-icon {
  animation-delay: 0.4s;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.feature-card h3 {
  font-family: roboto, sans-serif;
  font-size: 1.25rem;
  font-weight: 600;
  color: #222;
  margin: 0 0 12px 0;
}

.dark .feature-card h3 {
  color: white;
}

.feature-card p {
  font-family: roboto, sans-serif;
  font-size: 0.95rem;
  color: #666;
  line-height: 1.6;
  margin: 0;
}

.dark .feature-card p {
  color: #aaa;
}

/* Info Card */
.info-card {
  display: flex;
  align-items: center;
  gap: 16px;
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
  font-size: 1rem;
  color: #222;
  margin: 0;
  line-height: 1.6;
}

.dark .info-card p {
  color: white;
}

/* Responsive */
@media (max-width: 767px) {
  .hero-card {
    padding: 32px 20px;
    border-radius: 20px;
  }

  .hero-card h1 {
    font-size: 1.8rem;
  }

  .hero-card h2 {
    font-size: 1rem;
  }

  .btn-open-mail {
    width: 100%;
    justify-content: center;
    font-size: 1rem;
    padding: 0.875rem 1.5rem;
  }

  .features-section {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .info-card {
    flex-direction: column;
    text-align: center;
    padding: 20px;
  }
}
</style>
