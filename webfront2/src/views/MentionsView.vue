<template>
  <div class="mentions-view">
    <!-- Hero Card -->
    <div class="hero-card">
      <div class="legal-icon-wrapper">
        <svg class="legal-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
        </svg>
      </div>
      <h1>⚖️ Mentions Légales</h1>
      <h2>Informations légales et réglementaires</h2>
    </div>

    <!-- Contenu légal -->
    <div class="legal-content-card">
      <div v-html="$t('mentions')" class="legal-text"></div>
    </div>

    <!-- Sections supplémentaires -->
    <div class="legal-sections">
      <div class="legal-section">
        <div class="section-icon">🏢</div>
        <h3>Éditeur</h3>
        <p>Office1789 SAS</p>
        <p>Capital social : 100 000€</p>
        <p>SIRET : 123 456 789 00012</p>
      </div>
      
      <div class="legal-section">
        <div class="section-icon">🌐</div>
        <h3>Hébergement</h3>
        <p>OVH Cloud</p>
        <p>2 rue Kellermann</p>
        <p>59100 Roubaix, France</p>
      </div>
      
      <div class="legal-section">
        <div class="section-icon">🔒</div>
        <h3>Données personnelles</h3>
        <p>Conformité RGPD</p>
        <p>Protection des données</p>
        <p>Droit d'accès et rectification</p>
      </div>
    </div>
  </div>
</template>

<script setup>
//const date = new Date().getYear()+1900
</script>

<style scoped>
.mentions-view {
  width: 100%;
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
  min-height: 260px;
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

.legal-icon-wrapper {
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

.legal-icon {
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

/* Contenu légal */
.legal-content-card {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  padding: 40px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  margin-bottom: 32px;
}

.dark .legal-content-card {
  background: #1C1C1E;
}

.legal-text {
  font-family: arial, sans-serif;
  font-size: 1rem;
  line-height: 1.8;
  color: #333;
  text-align: justify;
}

.dark .legal-text {
  color: #ddd;
}

.legal-text :deep(h2) {
  font-family: roboto, sans-serif;
  font-size: 1.5rem;
  font-weight: 600;
  color: #222;
  margin: 24px 0 16px 0;
}

.dark .legal-text :deep(h2) {
  color: white;
}

.legal-text :deep(h3) {
  font-family: roboto, sans-serif;
  font-size: 1.2rem;
  font-weight: 600;
  color: #444;
  margin: 20px 0 12px 0;
}

.dark .legal-text :deep(h3) {
  color: #ccc;
}

/* Sections légales */
.legal-sections {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-bottom: 32px;
}

.legal-section {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 20px;
  padding: 32px 24px;
  text-align: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.dark .legal-section {
  background: #1C1C1E;
}

.legal-section:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.section-icon {
  font-size: 2.5rem;
  margin-bottom: 12px;
  animation: bounce 2s ease-in-out infinite;
}

.legal-section:nth-child(2) .section-icon {
  animation-delay: 0.2s;
}

.legal-section:nth-child(3) .section-icon {
  animation-delay: 0.4s;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.legal-section h3 {
  font-family: roboto, sans-serif;
  font-size: 1.15rem;
  font-weight: 600;
  color: #222;
  margin: 0 0 12px 0;
}

.dark .legal-section h3 {
  color: white;
}

.legal-section p {
  font-family: roboto, sans-serif;
  font-size: 0.9rem;
  color: #666;
  margin: 4px 0;
  line-height: 1.5;
}

.dark .legal-section p {
  color: #aaa;
}

/* Responsive */
@media (max-width: 768px) {
  .hero-card {
    padding: 32px 20px;
  }

  .hero-card h1 {
    font-size: 1.8rem;
  }

  .hero-card h2 {
    font-size: 1rem;
  }

  .legal-content-card {
    padding: 24px;
  }

  .legal-sections {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>
