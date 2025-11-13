<template>
  <div class="home-view">
    <!-- Hero Card avec Icon -->
    <div class="hero-banner">
      <div class="hero-icon-wrapper">
        <svg class="hero-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
        </svg>
      </div>
      <h1>{{ $t('presentation1I') }}</h1>
      <h2>Leurs privilèges, notre Révolution</h2>
      <div class="hero-actions">
        <RouterLink to="/createaccount" class="btn-create">
          {{ $t('create') }}
        </RouterLink>
        <RouterLink to="/about" class="btn-try">
          {{ $t('try') }}
        </RouterLink>
      </div>
    </div>

    <!-- Progress Bar -->
    <div class="progress-card">
      <ul class="progressbar">
        <li class="completed">
          <span class="marker">1</span>
          <p>{{ $t('publication') }}</p>
        </li>
        <li>
          <span class="marker">2</span>
          <p>{{ $t('pconv') }}</p>
        </li>
        <li>
          <span class="marker">3</span>
          <p>{{ $t('pdrive') }}</p>
        </li>
        <li>
          <span class="marker">4</span>
          <p>{{ $t('pemail') }}</p>
        </li>
        <li>
          <span class="marker">5</span>
          <p>{{ $t('publis') }}</p>
        </li>
      </ul>
    </div>

    <!-- Main Content -->
    <div class="main-section">
      <!-- Carrousel -->
      <div class="carousel-card">
        <div class="carousel-wrapper">
          <button @click="move(0)" class="carousel-btn-left">‹</button>
          <div class="carousel-viewport">
            <div class="carousel-track" :style="posc">
              <img src="@/assets/emailT.png" alt="Email" />
              <img src="@/assets/cloudT.png" alt="Cloud" />
              <img src="@/assets/dataT.png" alt="Data" />
              <img src="@/assets/writeT.png" alt="Write" />
            </div>
          </div>
          <button @click="move(1)" class="carousel-btn-right">›</button>
        </div>
        
        <div class="carousel-indicators">
          <span 
            v-for="(item, index) in presentation2" 
            :key="index"
            @click="moveTo(index)"
            :class="['dot', { active: movnb === index }]"
          ></span>
        </div>
        
        <div class="carousel-description" v-html="$t(presentation2[movnb])"></div>
      </div>

      <!-- News Section -->
      <div class="news-card">
        <h3>{{ $t('news') }}</h3>
        <ul v-html="$t('newsli')"></ul>
      </div>
    </div>

    <!-- Services Section -->
    <div class="services-section">
      <h2 class="section-title">Nos Services</h2>
      <p class="section-subtitle">Une suite complète d'outils pour votre productivité</p>
      
      <div class="services-grid">
        <div class="service-card">
          <div class="service-icon">📧</div>
          <h3>Messagerie</h3>
          <p>Email professionnel avec 10 Go de stockage, calendrier intégré et contacts centralisés</p>
        </div>
        
        <div class="service-card">
          <div class="service-icon">💬</div>
          <h3>Chat</h3>
          <p>Communication instantanée avec vos collaborateurs, groupes de discussion et appels intégrés</p>
        </div>
        
        <div class="service-card">
          <div class="service-icon">☁️</div>
          <h3>Stockage Cloud</h3>
          <p>50 Go de stockage sécurisé pour tous vos documents avec partage facile</p>
        </div>
        
        <div class="service-card">
          <div class="service-icon">📄</div>
          <h3>Suite Bureautique</h3>
          <p>Éditeur de texte, tableur et présentations avec OnlyOffice intégré</p>
        </div>
        
        <div class="service-card">
          <div class="service-icon">📅</div>
          <h3>Calendrier</h3>
          <p>Organisez vos rendez-vous et événements avec un calendrier complet</p>
        </div>
        
        <div class="service-card">
          <div class="service-icon">🔒</div>
          <h3>Sécurité</h3>
          <p>Vos données sont hébergées en Europe avec chiffrement de bout en bout</p>
        </div>
      </div>
    </div>

    <!-- French Power -->
    <div class="french-card">
      <span class="french-flag"></span>
      <img src="@/assets/europe.png" alt="Europe" />
      <p class="french-title">French/Europe Power</p>
      <p class="french-text">{{ $t('apo') }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue"

let mov = 0
let movnb = ref(0)
let posc = ref(`transform: translateX(0%)`)
let presentation2 = ["servicee", "servicec", "serviced", "servicew"]

function move(dir) {
  mov = dir ? mov -= 100 : mov += 100
  movnb.value = dir ? movnb.value += 1 : movnb.value -= 1
  
  if (mov < -300) {
    mov = 0
    movnb.value = 0
  }
  if (mov > 0) {
    mov = -300
    movnb.value = 3
  }
  
  posc.value = `transform: translateX(${mov}%); transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);`
}

function moveTo(index) {
  mov = -100 * index
  movnb.value = index
  posc.value = `transform: translateX(${mov}%); transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);`
}

// Auto-play
setInterval(() => move(1), 5000)
</script>

<style scoped>
.home-view {
  width: 100%;
  animation: fadeIn 0.5s ease-in;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

/* Hero Banner avec style moderne */
.hero-banner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  width: 100%;
  padding: 60px 32px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 32px;
  gap: 24px;
  position: relative;
  overflow: hidden;
}

.dark .hero-banner {
  background: #1C1C1E;
}

.hero-icon-wrapper {
  width: 100px;
  height: 100px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse 2s ease-in-out infinite;
  box-shadow: 0 8px 32px rgba(0, 0, 255, 0.3);
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.08); }
}

.hero-icon {
  width: 50px;
  height: 50px;
  color: white;
}

.hero-banner h1 {
  font-family: 'Roboto', sans-serif;
  font-size: 3rem;
  font-weight: 700;
  text-align: center;
  letter-spacing: 2px;
  color: #222;
  margin: 0;
}

.dark .hero-banner h1 {
  color: white;
}

.hero-banner h2 {
  font-family: 'Roboto', sans-serif;
  font-size: 1.5rem;
  font-style: italic;
  text-align: center;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 1px;
  margin: 0;
}

.hero-actions {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-top: 16px;
}

.btn-create,
.btn-try {
  font-family: 'Roboto', sans-serif;
  font-size: 1.1rem;
  padding: 0.875rem 2.5rem;
  border-radius: 32px;
  border: none;
  cursor: pointer;
  text-decoration: none;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-create {
  background: -webkit-linear-gradient(30deg, blue, red);
  color: white;
}

.btn-try {
  background: white;
  border: 2px solid blue;
  color: blue;
}

.dark .btn-try {
  background: rgba(30,30,40,0.95);
  color: white;
  border-color: red;
}

.btn-create:hover,
.btn-try:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.20);
}

/* Services Section */
.services-section {
  margin-bottom: 48px;
}

.section-title {
  font-family: roboto, sans-serif;
  font-size: 2.5rem;
  font-weight: 700;
  text-align: center;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0 0 8px 0;
  letter-spacing: 1px;
}

.section-subtitle {
  font-family: roboto, sans-serif;
  font-size: 1.1rem;
  text-align: center;
  color: #666;
  margin: 0 0 32px 0;
  font-style: italic;
}

.dark .section-subtitle {
  color: #aaa;
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
}

.service-card {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  padding: 32px 24px;
  text-align: center;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.dark .service-card {
  background: #1C1C1E;
}

.service-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.service-icon {
  font-size: 4rem;
  margin-bottom: 16px;
  animation: float 3s ease-in-out infinite;
}

.service-card:nth-child(2) .service-icon {
  animation-delay: 0.3s;
}

.service-card:nth-child(3) .service-icon {
  animation-delay: 0.6s;
}

.service-card:nth-child(4) .service-icon {
  animation-delay: 0.9s;
}

.service-card:nth-child(5) .service-icon {
  animation-delay: 1.2s;
}

.service-card:nth-child(6) .service-icon {
  animation-delay: 1.5s;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.service-card h3 {
  font-family: roboto, sans-serif;
  font-size: 1.25rem;
  font-weight: 600;
  color: #222;
  margin: 0 0 12px 0;
}

.dark .service-card h3 {
  color: white;
}

.service-card p {
  font-family: roboto, sans-serif;
  font-size: 0.95rem;
  color: #666;
  line-height: 1.6;
  margin: 0;
}

.dark .service-card p {
  color: #aaa;
}

/* Progress Bar */
.progress-card {
  margin-bottom: 32px;
}

.progressbar {
  display: flex;
  list-style: none;
  background: -webkit-linear-gradient(30deg, blue, red);
  color: white;
  border-radius: 10px;
  padding: 20px;
  margin: 0;
  gap: 0;
}

.progressbar li {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
  position: relative;
}

.progressbar .marker {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  border: 2px solid white;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 700;
  z-index: 1;
  margin-bottom: 8px;
  font-size: 0.9rem;
}

.progressbar .completed .marker {
  background: white;
  border-color: white;
  color: #06c;
}

.progressbar li:not(:last-child)::after {
  content: '';
  position: absolute;
  top: 16px;
  left: calc(50% + 16px);
  width: calc(100% - 32px);
  height: 3px;
  background: rgba(255, 255, 255, 0.4);
}

.progressbar .completed:not(:last-child)::after {
  background: white;
}

.progressbar p {
  font-size: 0.9rem;
  text-align: center;
  margin: 0;
}

/* Main Section */
.main-section {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 32px;
  margin-bottom: 32px;
  padding: 32px;
  background: url("@/assets/napo2.jpg") center/cover no-repeat;
  border-radius: 24px;
  position: relative;
  overflow: hidden;
}

.main-section::before {
  content: "";
  position: absolute;
  inset: 0;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(1px);
  z-index: 0;
  border-radius: 24px;
}

.dark .main-section::before {
  background: rgba(0, 0, 0, 0.6);
}

.carousel-card,
.news-card {
  position: relative;
  z-index: 1;
}

/* Carousel */
.carousel-card {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.carousel-wrapper {
  display: flex;
  align-items: center;
  gap: 16px;
  height: 300px;
}

.carousel-btn-left,
.carousel-btn-right {
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 50%;
  color: white;
  border: none;
  font-size: 48px;
  width: 56px;
  height: 56px;
  cursor: pointer;
  transition: all 0.3s ease;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

.carousel-btn-left:hover,
.carousel-btn-right:hover {
  transform: scale(1.1);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
}

.carousel-viewport {
  flex: 1;
  overflow: hidden;
  border-radius: 10px;
}

.carousel-track {
  display: flex;
  transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.carousel-track img {
  min-width: 100%;
  max-width: 100%;
  height: 300px;
  object-fit: contain;
}

.carousel-indicators {
  display: flex;
  justify-content: center;
  gap: 12px;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: all 0.3s ease;
}

.dot.active {
  background: white;
  transform: scale(1.3);
}

.carousel-description {
  font-family: roboto, sans-serif;
  background: linear-gradient(30deg, blue, red);
  padding: 20px;
  border-radius: 10px;
  color: white;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.37);
  line-height: 1.6;
}

/* News Card */
.news-card {
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 10px;
  padding: 20px;
  color: white;
  font-family: arial, sans-serif;
}

.news-card h3 {
  font-family: roboto, sans-serif;
  border-bottom: 1px solid rgba(255, 255, 255, 0.5);
  padding-bottom: 12px;
  margin-top: 0;
}

.news-card ul :deep(li) {
  margin-top: 20px;
  line-height: 1.6;
}

/* French Card amélioré */
.french-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 24px;
  padding: 60px 40px;
  background: -webkit-linear-gradient(135deg, rgba(0, 0, 139, 0.05), rgba(255, 0, 0, 0.05));
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  position: relative;
  overflow: hidden;
  border: 2px solid rgba(0, 0, 139, 0.1);
}

.dark .french-card {
  background: -webkit-linear-gradient(135deg, rgba(0, 0, 139, 0.15), rgba(255, 0, 0, 0.15));
  border-color: rgba(255, 255, 255, 0.1);
}

.french-card::before {
  content: '';
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: repeating-linear-gradient(
    45deg,
    transparent,
    transparent 10px,
    rgba(0, 0, 139, 0.03) 10px,
    rgba(0, 0, 139, 0.03) 20px
  );
  animation: slide 20s linear infinite;
  pointer-events: none;
}

@keyframes slide {
  0% { transform: translate(0, 0); }
  100% { transform: translate(50px, 50px); }
}

.french-flag,
.french-card img {
  width: 150px;
  height: 120px;
  border-radius: 16px;
  border: 3px solid rgba(255, 255, 255, 0.8);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  object-fit: cover;
  transition: all 0.3s ease;
  z-index: 1;
}

.french-card img:hover,
.french-flag:hover {
  transform: scale(1.05) rotate(2deg);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.25);
}

.french-flag {
  background: linear-gradient(90deg, blue 33%, white 33%, white 66%, red 66%);
}

.french-card img {
  background: none;
}

.french-title {
  font-family: 'Roboto', sans-serif;
  font-size: 2rem;
  font-weight: 700;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
  text-align: center;
  letter-spacing: 1px;
  z-index: 1;
}

.french-text {
  font-family: 'Roboto', sans-serif;
  font-size: 1.1rem;
  text-align: center;
  color: #333;
  max-width: 800px;
  line-height: 1.8;
  margin: 0;
  z-index: 1;
}

.dark .french-text {
  color: #ddd;
}

/* Responsive */
@media (max-width: 767px) {
  .hero-card {
    padding: 20px 12px;
    border-radius: 12px;
  }

  .hero-banner {
    min-height: 320px;
    padding: 40px 20px;
  }

  .hero-icon-wrapper {
    width: 80px;
    height: 80px;
  }

  .hero-icon {
    width: 40px;
    height: 40px;
  }

  .hero-banner h1 {
    font-size: 1.8rem;
  }

  .hero-banner h2 {
    font-size: 1rem;
  }

  .hero-actions {
    flex-direction: column;
    width: 100%;
  }

  .btn-create,
  .btn-try {
    width: 100%;
    justify-content: center;
    padding: 12px 20px;
    font-size: 1rem;
  }

  .services-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .section-title {
    font-size: 1.8rem;
  }

  .section-subtitle {
    font-size: 0.9rem;
  }

  .progressbar {
    flex-direction: column;
    align-items: flex-start;
    text-align: left;
    gap: 20px;
  }

  .progressbar li {
    flex-direction: row;
    align-items: center;
    gap: 12px;
  }

  .progressbar li::after {
    display: none;
  }

  .progressbar p {
    text-align: left;
  }

  .main-section {
    grid-template-columns: 1fr;
    gap: 40px;
  }

  .carousel-wrapper {
    height: 180px;
  }

  .carousel-track img {
    height: 180px;
  }

  .carousel-btn-left,
  .carousel-btn-right {
    width: 40px;
    height: 40px;
    font-size: 32px;
  }

  .french-card {
    flex-direction: column;
    gap: 12px;
    padding: 16px;
  }

  .french-flag,
  .french-card img {
    width: 60px;
    height: 50px;
  }

  .french-card p {
    font-size: 0.8rem;
  }
}
</style>
