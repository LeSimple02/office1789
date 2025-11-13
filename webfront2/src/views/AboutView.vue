<template>
  <div class="about-view">
    <!-- Hero Card -->
    <div class="hero-card">
      <div class="about-icon-wrapper">
        <svg class="about-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
      </div>
      <h1>ℹ️ À propos d'Office1789</h1>
      <h2>Découvrez notre suite collaborative souveraine</h2>
    </div>

    <!-- Mission Section -->
    <div class="mission-section">
      <h2 class="section-title">🎯 Notre Mission</h2>
      <div class="mission-grid">
        <div class="mission-card">
          <div class="mission-icon">🚀</div>
          <h3>Innovation</h3>
          <p>Développer des outils collaboratifs modernes et performants qui rivalisent avec les géants américains</p>
        </div>
        
        <div class="mission-card">
          <div class="mission-icon">🇫🇷</div>
          <h3>Souveraineté</h3>
          <p>Garantir l'hébergement de vos données en France et en Europe pour une totale maîtrise</p>
        </div>
        
        <div class="mission-card">
          <div class="mission-icon">🔒</div>
          <h3>Confidentialité</h3>
          <p>Protéger vos données avec un chiffrement de bout en bout et sans revente à des tiers</p>
        </div>
        
        <div class="mission-card">
          <div class="mission-icon">♻️</div>
          <h3>Éco-responsable</h3>
          <p>Utiliser des serveurs verts et optimiser notre infrastructure pour réduire notre empreinte carbone</p>
        </div>
      </div>
    </div>

    <!-- Values Section -->
    <div class="values-section">
      <h2 class="section-title">💎 Nos Valeurs</h2>
      <div class="values-content">
        <div class="value-item">
          <span class="value-number">01</span>
          <div class="value-info">
            <h3>Transparence</h3>
            <p>Nous croyons en une communication ouverte avec nos utilisateurs. Pas de conditions cachées, pas de collecte de données à votre insu.</p>
          </div>
        </div>
        
        <div class="value-item">
          <span class="value-number">02</span>
          <div class="value-info">
            <h3>Accessibilité</h3>
            <p>Des tarifs justes et une offre gratuite complète pour rendre le numérique accessible à tous, particuliers comme entreprises.</p>
          </div>
        </div>
        
        <div class="value-item">
          <span class="value-number">03</span>
          <div class="value-info">
            <h3>Performance</h3>
            <p>Des outils rapides, fiables et disponibles 24/7 avec une infrastructure cloud de dernière génération.</p>
          </div>
        </div>
        
        <div class="value-item">
          <span class="value-number">04</span>
          <div class="value-info">
            <h3>Innovation Continue</h3>
            <p>Nous améliorons constamment nos services en écoutant vos retours et en intégrant les dernières technologies.</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Why Choose Us Section -->
    <div class="why-section">
      <h2 class="section-title">⭐ Pourquoi Office1789 ?</h2>
      <div class="why-grid">
        <div class="why-card">
          <h3>📊 Statistiques</h3>
          <div class="stats">
            <div class="stat-item">
              <span class="stat-value">99.9%</span>
              <span class="stat-label">Disponibilité</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">50 Go</span>
              <span class="stat-label">Stockage</span>
            </div>
            <div class="stat-item">
              <span class="stat-value">100%</span>
              <span class="stat-label">Français</span>
            </div>
          </div>
        </div>
        
        <div class="why-card">
          <h3>🎁 Avantages</h3>
          <ul class="advantages-list">
            <li>✅ Sans publicité</li>
            <li>✅ Open Source</li>
            <li>✅ Support réactif</li>
            <li>✅ Migration facile</li>
            <li>✅ Applications mobiles</li>
            <li>✅ Intégrations tierces</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Contenu About -->
    <div class="about-content-card">
      <div class="about" v-html="$t('aboutd')"></div>
    </div>

    <!-- Offres -->
    <div class="offers-section">
      <h2 class="section-title">📦 Nos Offres</h2>
      <div id="offer" v-html="$t('offer')"></div>
    </div>

    <!-- CTA -->
    <div class="cta-card">
      <h3>Prêt à rejoindre Office1789 ?</h3>
      <p>Créez votre compte gratuitement et profitez de tous nos services</p>
      <RouterLink to="/createaccount" class="btn-cta">{{$t('createac')}}</RouterLink>
    </div>
  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.about-view {
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

.about-icon-wrapper {
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

.about-icon {
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

/* Mission Section */
.mission-section {
  margin-bottom: 48px;
}

.mission-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.mission-card {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  padding: 32px;
  text-align: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.dark .mission-card {
  background: #1C1C1E;
}

.mission-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.mission-icon {
  font-size: 4rem;
  margin-bottom: 16px;
  animation: bounce 2s ease-in-out infinite;
}

.mission-card:nth-child(2) .mission-icon {
  animation-delay: 0.2s;
}

.mission-card:nth-child(3) .mission-icon {
  animation-delay: 0.4s;
}

.mission-card:nth-child(4) .mission-icon {
  animation-delay: 0.6s;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.mission-card h3 {
  font-family: roboto, sans-serif;
  font-size: 1.4rem;
  font-weight: 600;
  color: #222;
  margin: 0 0 12px 0;
}

.dark .mission-card h3 {
  color: white;
}

.mission-card p {
  font-family: roboto, sans-serif;
  font-size: 1rem;
  color: #666;
  line-height: 1.6;
  margin: 0;
}

.dark .mission-card p {
  color: #aaa;
}

/* Values Section */
.values-section {
  margin-bottom: 48px;
}

.values-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.value-item {
  display: flex;
  gap: 24px;
  align-items: flex-start;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  padding: 32px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.dark .value-item {
  background: #1C1C1E;
}

.value-item:hover {
  transform: translateX(10px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.value-number {
  font-family: roboto, sans-serif;
  font-size: 3rem;
  font-weight: 700;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  flex-shrink: 0;
  line-height: 1;
}

.value-info h3 {
  font-family: roboto, sans-serif;
  font-size: 1.4rem;
  font-weight: 600;
  color: #222;
  margin: 0 0 8px 0;
}

.dark .value-info h3 {
  color: white;
}

.value-info p {
  font-family: roboto, sans-serif;
  font-size: 1rem;
  color: #666;
  line-height: 1.6;
  margin: 0;
}

.dark .value-info p {
  color: #aaa;
}

/* Why Section */
.why-section {
  margin-bottom: 48px;
}

.why-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
}

.why-card {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  padding: 32px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.dark .why-card {
  background: #1C1C1E;
}

.why-card h3 {
  font-family: roboto, sans-serif;
  font-size: 1.5rem;
  font-weight: 600;
  color: #222;
  margin: 0 0 24px 0;
  text-align: center;
}

.dark .why-card h3 {
  color: white;
}

.stats {
  display: flex;
  justify-content: space-around;
  gap: 16px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.stat-value {
  font-family: roboto, sans-serif;
  font-size: 2.5rem;
  font-weight: 700;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.stat-label {
  font-family: roboto, sans-serif;
  font-size: 0.95rem;
  color: #666;
  font-weight: 500;
}

.dark .stat-label {
  color: #aaa;
}

.advantages-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.advantages-list li {
  font-family: roboto, sans-serif;
  font-size: 1rem;
  color: #333;
}

.dark .advantages-list li {
  color: #ddd;
}

/* Contenu About */
.about-content-card {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  padding: 40px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  margin-bottom: 32px;
}

.dark .about-content-card {
  background: #1C1C1E;
}

.about {
  font-family: arial, sans-serif;
  font-size: 1.05rem;
  line-height: 1.8;
  color: #333;
}

.dark .about {
  color: #ddd;
}

.about :deep(h1) {
  font-family: roboto, sans-serif;
  font-size: 2rem;
  font-weight: 600;
  color: #222;
  margin-bottom: 20px;
  width: 100%;
  max-width: none;
}

.dark .about :deep(h1) {
  color: white;
}

/* Section des offres */
.offers-section {
  margin-bottom: 32px;
}

.section-title {
  font-family: roboto, sans-serif;
  font-size: 2rem;
  font-weight: 700;
  text-align: center;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin-bottom: 32px;
}

#offer {
  display: flex;
  flex-wrap: wrap;
  align-items: stretch;
  justify-content: center;
  gap: 32px;
}

#offer :deep(.offerm) {
  background: rgba(245, 245, 247, 0.85);
  border: 2px solid transparent;
  border-radius: 20px;
  padding: 30px;
  width: 270px;
  height: auto;
  font-family: arial, sans-serif;
  transition: all 0.3s ease;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.dark #offer :deep(.offerm) {
  background: #1C1C1E;
  border-color: transparent;
}

#offer :deep(.offerm):hover {
  transform: translateY(-8px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
  border-color: blue;
}

#offer :deep(.offerm > h3) {
  text-align: center;
  padding: 12px;
  margin: -30px -30px 20px -30px;
  height: auto;
  color: white;
  background: -webkit-linear-gradient(30deg, blue, red);
  transition: all 0.3s ease;
  font-size: 1.3rem;
  font-weight: 700;
}

#offer :deep(.offerm):hover > h3 {
  background: -webkit-linear-gradient(30deg, red, blue);
}

#offer :deep(.lioffer) {
  list-style: none;
  font-size: 1rem;
  text-align: left;
  padding: 0;
  margin: 0;
}

#offer :deep(.lioffer > li) {
  margin-top: 16px;
  padding-left: 24px;
  position: relative;
  color: #333;
}

.dark #offer :deep(.lioffer > li) {
  color: #ddd;
}

#offer :deep(.lioffer > li)::before {
  content: '✓';
  position: absolute;
  left: 0;
  color: green;
  font-weight: 700;
}

/* CTA Card */
.cta-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 32px;
  background: rgba(220, 220, 240, 0.95);
  border-radius: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  gap: 16px;
  text-align: center;
}

.dark .cta-card {
  background: rgba(0, 50, 0, 0.2);
}

.cta-card h3 {
  font-family: roboto, sans-serif;
  font-size: 1.8rem;
  font-weight: 700;
  color: #222;
  margin: 0;
}

.dark .cta-card h3 {
  color: white;
}

.cta-card p {
  font-family: roboto, sans-serif;
  font-size: 1.1rem;
  color: #666;
  margin: 0;
}

.dark .cta-card p {
  color: #aaa;
}

.btn-cta {
  font-family: roboto, sans-serif;
  font-size: 1.2rem;
  padding: 1rem 3rem;
  border-radius: 32px;
  border: none;
  cursor: pointer;
  background: -webkit-linear-gradient(30deg, blue, red);
  color: white;
  text-decoration: none;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  font-weight: 600;
  margin-top: 16px;
}

.btn-cta:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
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

  .about-content-card {
    padding: 24px;
  }

  .mission-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .values-content {
    gap: 16px;
  }

  .value-item {
    flex-direction: column;
    padding: 24px;
    align-items: center;
    text-align: center;
  }

  .value-item:hover {
    transform: translateY(-5px);
  }

  .why-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .stats {
    flex-direction: column;
  }

  .advantages-list {
    grid-template-columns: 1fr;
  }

  #offer {
    gap: 20px;
  }

  #offer :deep(.offerm) {
    width: 100%;
    max-width: 320px;
  }

  .cta-card {
    padding: 32px 20px;
  }

  .cta-card h3 {
    font-size: 1.4rem;
  }

  .btn-cta {
    width: 100%;
    font-size: 1rem;
    padding: 0.875rem 2rem;
  }
}
</style>
