<template>
  <div class="chat-view">
    <!-- Hero Card -->
    <div class="hero-card">
      <div class="chat-icon-wrapper">
        <svg class="chat-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"/>
        </svg>
      </div>
      <h1>💬 Chat Office1789</h1>
      <h2>Communiquez en temps réel avec votre équipe</h2>
      
      <MatrixAccessSSO />
    </div>

    <!-- Features Section -->
    <div class="features-section">
      <div class="feature-card">
        <div class="feature-icon">💬</div>
        <h3>Instantané</h3>
        <p>Messagerie en temps réel pour une communication fluide</p>
      </div>
      
      <div class="feature-card">
        <div class="feature-icon">👥</div>
        <h3>Groupes</h3>
        <p>Créez des groupes de discussion pour vos projets</p>
      </div>
      
      <div class="feature-card">
        <div class="feature-icon">📞</div>
        <h3>Appels</h3>
        <p>Appels audio et vidéo intégrés</p>
      </div>
    </div>

    <!-- Info Box -->
    <div class="info-card">
      <svg class="info-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
      </svg>
      <p>Le chat s'ouvrira dans une nouvelle fenêtre pour une meilleure expérience utilisateur</p>
    </div>
  </div>
</template>

<script setup>
import { gls } from "@/stores/global"
import router from "@/router/index"
import MatrixAccessSSO from "@/components/MatrixAccessSSO.vue"

if (gls().log != 1) {
  router.push("login")
}
</script>

<style scoped>
.chat-view {
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

.chat-icon-wrapper {
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

.chat-icon {
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

.btn-open-chat {
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

.btn-open-chat:hover {
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

  .btn-open-chat {
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
