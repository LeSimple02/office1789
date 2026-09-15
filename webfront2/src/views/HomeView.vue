<template>
  <div class="home">
    <!-- HERO -->
    <section class="hero">
      <div class="hero-badge">Office1789</div>
      <h1 class="hero-title">{{ $t('presentation1I') }}</h1>
      <p class="hero-sub">{{ $t('theirPrivilegesOurRevolution') }}</p>
      <div class="hero-cta">
        <RouterLink to="/createaccount" class="cta-primary">{{ $t('create') }}</RouterLink>
        <RouterLink to="/about" class="cta-secondary">{{ $t('try') }}</RouterLink>
      </div>
    </section>

    <!-- ROADMAP -->
    <section class="roadmap">
      <ul class="steps">
        <li class="step done">
          <span class="step-num">1</span>
          <span class="step-label">{{ $t('publication') }}</span>
        </li>
        <li class="step">
          <span class="step-num">2</span>
          <span class="step-label">{{ $t('pconv') }}</span>
        </li>
        <li class="step">
          <span class="step-num">3</span>
          <span class="step-label">{{ $t('pdrive') }}</span>
        </li>
        <li class="step">
          <span class="step-num">4</span>
          <span class="step-label">{{ $t('pemail') }}</span>
        </li>
        <li class="step">
          <span class="step-num">5</span>
          <span class="step-label">{{ $t('publis') }}</span>
        </li>
      </ul>
    </section>

    <!-- CAROUSEL + NEWS -->
    <section class="split">
      <div class="carousel">
        <div class="carousel-track-wrap">
          <button class="carousel-arrow left" @click="move(0)" aria-label="Prev">‹</button>
          <div class="carousel-window">
            <div class="carousel-strip" :style="posc">
              <img src="@/assets/emailT.png" alt="Email" />
              <img src="@/assets/cloudT.png" alt="Cloud" />
              <img src="@/assets/dataT.png" alt="Data" />
              <img src="@/assets/writeT.png" alt="Write" />
            </div>
          </div>
          <button class="carousel-arrow right" @click="move(1)" aria-label="Next">›</button>
        </div>
        <div class="carousel-dots">
          <button
            v-for="(item, index) in presentation2"
            :key="index"
            @click="moveTo(index)"
            :class="['dot', { on: movnb === index }]"
            :aria-label="'Slide ' + (index + 1)"
          ></button>
        </div>
        <p class="carousel-text" v-html="$t(presentation2[movnb])"></p>
      </div>

      <aside class="news">
        <h3 class="news-title">{{ $t('news') }}</h3>
        <ul class="news-list" v-html="$t('newsli')"></ul>
      </aside>
    </section>

    <!-- SERVICES -->
    <section class="services">
      <h2 class="block-title">{{ $t('ourServices') }}</h2>
      <p class="block-sub">{{ $t('completeProductivityTools') }}</p>
      <div class="services-row">
        <article class="svc" v-for="(svc, i) in [
          { icon: '📧', title: $t('messaging'), desc: $t('professionalEmailWith10GB') },
          { icon: '💬', title: $t('chat'), desc: $t('instantCommunicationWithColleagues') },
          { icon: '☁️', title: $t('cloudStorage'), desc: $t('secureStorage50GB') },
          { icon: '📄', title: $t('officeSuite'), desc: $t('textEditorSpreadsheetPresentation') },
          { icon: '📅', title: $t('calendarService'), desc: $t('organizeAppointmentsEvents') },
          { icon: '🔒', title: 'Sécurité', desc: 'Vos données sont hébergées en Europe avec chiffrement de bout en bout' }
        ]" :key="i">
          <div class="svc-icon">{{ svc.icon }}</div>
          <h3 class="svc-title">{{ svc.title }}</h3>
          <p class="svc-desc">{{ svc.desc }}</p>
        </article>
      </div>
    </section>

    <!-- FR -->
    <section class="fr-card">
      <span class="fr-flag"></span>
      <img src="@/assets/europe.png" alt="Europe" class="fr-img" />
      <h3 class="fr-title">French/Europe Power</h3>
      <p class="fr-text">{{ $t('apo') }}</p>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from "vue"

let mov = 0
let movnb = ref(0)
let posc = ref(`transform: translateX(0%)`)
let presentation2 = ["servicee", "servicec", "serviced", "servicew"]

function move(dir) {
  mov = dir ? mov -= 100 : mov += 100
  movnb.value = dir ? movnb.value += 1 : movnb.value -= 1
  if (mov < -300) { mov = 0; movnb.value = 0 }
  if (mov > 0) { mov = -300; movnb.value = 3 }
  posc.value = `transform: translateX(${mov}%); transition: transform .6s var(--ease);`
}

function moveTo(index) {
  mov = -100 * index
  movnb.value = index
  posc.value = `transform: translateX(${mov}%); transition: transform .6s var(--ease);`
}

let intervalId
onMounted(() => { intervalId = setInterval(() => move(1), 5000) })
onUnmounted(() => { clearInterval(intervalId) })
</script>

<style scoped>
.home {
  max-width: var(--content-max);
  margin: 0 auto;
  padding: 2.5rem 1.5rem 5rem;
}

/* ===== HERO ===== */
.hero {
  text-align: center;
  padding: 4rem 2rem 3rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--r-2xl);
  margin-bottom: 2rem;
  animation: fadeInUp .6s var(--ease) both;
}
.hero-badge {
  display: inline-block;
  padding: .35rem .8rem;
  font-size: .75rem;
  font-weight: 700;
  letter-spacing: .05em;
  border-radius: var(--r-full);
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  margin-bottom: 1.25rem;
  animation: fadeInDown .4s var(--ease) both;
}
.hero-title {
  font-size: clamp(2rem, 5vw, 3.25rem);
  font-weight: 800;
  letter-spacing: -.035em;
  line-height: 1.08;
  margin-bottom: .75rem;
  animation: fadeInUp .5s var(--ease) .1s both;
}
.hero-sub {
  font-size: clamp(1rem, 2vw, 1.1875rem);
  color: var(--text-secondary);
  max-width: 520px;
  margin: 0 auto 2rem;
  line-height: 1.6;
  animation: fadeInUp .5s var(--ease) .2s both;
}
.hero-cta {
  display: flex;
  gap: .75rem;
  justify-content: center;
  flex-wrap: wrap;
  animation: fadeInUp .5s var(--ease) .3s both;
}
.cta-primary, .cta-secondary {
  padding: .75rem 1.75rem;
  font-size: .9375rem;
  font-weight: 600;
  border-radius: var(--r);
  transition: transform var(--tb), box-shadow var(--tb), background var(--tb), border-color var(--tb);
  display: inline-flex;
  align-items: center;
}
.cta-primary {
  background: var(--ink);
  color: var(--paper);
  border: 1px solid var(--ink);
}
.cta-primary:hover {
  background: var(--ink-soft);
  border-color: var(--ink-soft);
  transform: translateY(-2px);
  box-shadow: var(--sh-md);
}
.cta-secondary {
  background: var(--bg-primary);
  color: var(--text-primary);
  border: 1px solid var(--border-strong);
}
.cta-secondary:hover {
  border-color: var(--ink);
  transform: translateY(-2px);
  box-shadow: var(--sh-sm);
}

/* ===== ROADMAP ===== */
.roadmap {
  margin-bottom: 2rem;
  animation: fadeInUp .5s var(--ease) .15s both;
}
.steps {
  display: flex;
  gap: 0;
  padding: 1.25rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
}
.step {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  gap: .5rem;
}
.step:not(:last-child)::after {
  content: '';
  position: absolute;
  top: 14px;
  left: calc(50% + 16px);
  width: calc(100% - 32px);
  height: 2px;
  background: var(--border);
}
.step.done:not(:last-child)::after { background: var(--ink); }
.step-num {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: .8125rem;
  font-weight: 700;
  background: var(--bg-tertiary);
  border: 2px solid var(--border-strong);
  color: var(--text-muted);
  z-index: 1;
  transition: all var(--tb);
}
.step.done .step-num {
  background: var(--ink);
  border-color: var(--ink);
  color: var(--paper);
  animation: popIn .3s var(--ease-spring) both;
}
.step-label {
  font-size: .75rem;
  font-weight: 500;
  color: var(--text-secondary);
  text-align: center;
}

/* ===== SPLIT (carousel + news) ===== */
.split {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 1.5rem;
  margin-bottom: 2rem;
}
.carousel {
  background: var(--bg-primary);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1rem;
  animation: fadeInLeft .5s var(--ease) both;
}
.carousel-track-wrap {
  display: flex;
  align-items: center;
  gap: 1rem;
}
.carousel-arrow {
  width: 40px;
  height: 40px;
  flex-shrink: 0;
  border-radius: 50%;
  border: 1px solid var(--border-strong);
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 22px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--tb);
}
.carousel-arrow:hover {
  border-color: var(--ink);
  background: var(--bg-tertiary);
  transform: scale(1.08);
}
.carousel-window {
  flex: 1;
  overflow: hidden;
  border-radius: var(--r);
}
.carousel-strip {
  display: flex;
}
.carousel-strip img {
  min-width: 100%;
  height: 220px;
  object-fit: cover;
  border-radius: var(--r);
}
.carousel-dots {
  display: flex;
  justify-content: center;
  gap: .5rem;
}
.carousel-dots .dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  border: none;
  background: var(--border-strong);
  cursor: pointer;
  padding: 0;
  transition: all var(--tb);
}
.carousel-dots .dot.on {
  background: var(--ink);
  transform: scale(1.4);
}
.carousel-text {
  padding: .875rem 1.125rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--r);
  font-size: .9375rem;
  color: var(--text-secondary);
  line-height: 1.6;
}
.news {
  background: var(--bg-primary);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  padding: 1.5rem;
  animation: fadeInRight .5s var(--ease) both;
}
.news-title {
  font-size: 1rem;
  font-weight: 700;
  padding-bottom: .75rem;
  margin-bottom: 1rem;
  border-bottom: 1px solid var(--border);
}
.news-list { list-style: none; padding: 0; margin: 0; }
.news-list :deep(li) {
  font-size: .875rem;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-top: .75rem;
  padding-left: 14px;
  position: relative;
  animation: fadeIn .3s var(--ease) both;
}
.news-list :deep(li)::before {
  content: '';
  position: absolute;
  left: 0;
  top: 9px;
  width: 5px;
  height: 5px;
  border-radius: 50%;
  background: var(--text-muted);
}

/* ===== SERVICES ===== */
.services {
  margin-bottom: 2rem;
  animation: fadeInUp .5s var(--ease) .2s both;
}
.block-title {
  font-size: 1.75rem;
  font-weight: 800;
  letter-spacing: -.025em;
  text-align: center;
  margin-bottom: .375rem;
}
.block-sub {
  font-size: 1rem;
  color: var(--text-secondary);
  text-align: center;
  margin-bottom: 2rem;
}
.services-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.25rem;
}
.svc {
  background: var(--bg-primary);
  border: 1px solid var(--border);
  border-radius: var(--r-lg);
  padding: 1.75rem 1.5rem;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  transition: transform var(--tb), border-color var(--tb), box-shadow var(--tb);
  animation: fadeInUp .4s var(--ease) both;
}
.svc:nth-child(2) { animation-delay: 60ms; }
.svc:nth-child(3) { animation-delay: 120ms; }
.svc:nth-child(4) { animation-delay: 180ms; }
.svc:nth-child(5) { animation-delay: 240ms; }
.svc:nth-child(6) { animation-delay: 300ms; }
.svc:hover {
  border-color: var(--border-strong);
  box-shadow: var(--sh-md);
  transform: translateY(-5px);
}
.svc-icon {
  font-size: 2rem;
  margin-bottom: .875rem;
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--r-lg);
  background: var(--bg-tertiary);
  animation: float 3s var(--ease) infinite;
}
.svc:nth-child(2) .svc-icon { animation-delay: .3s; }
.svc:nth-child(3) .svc-icon { animation-delay: .6s; }
.svc:nth-child(4) .svc-icon { animation-delay: .9s; }
.svc:nth-child(5) .svc-icon { animation-delay: 1.2s; }
.svc:nth-child(6) .svc-icon { animation-delay: 1.5s; }
.svc-title {
  font-size: 1.0625rem;
  font-weight: 700;
  margin-bottom: .5rem;
}
.svc-desc {
  font-size: .875rem;
  color: var(--text-secondary);
  line-height: 1.6;
}

/* ===== FR CARD ===== */
.fr-card {
  text-align: center;
  padding: 3rem 2rem;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: var(--r-2xl);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.25rem;
  animation: fadeInUp .5s var(--ease) .25s both;
}
.fr-flag {
  width: 100px;
  height: 64px;
  border-radius: var(--r);
  background: linear-gradient(90deg, #0055A4 33%, #fff 33%, #fff 66%, #EF4135 66%);
  border: 1px solid var(--border);
  transition: transform var(--tb);
}
.fr-flag:hover { transform: scale(1.06); }
.fr-img {
  width: 100px;
  height: 64px;
  object-fit: cover;
  border-radius: var(--r);
  border: 1px solid var(--border);
}
.fr-title {
  font-size: 1.5rem;
  font-weight: 800;
  letter-spacing: -.02em;
}
.fr-text {
  font-size: 1rem;
  color: var(--text-secondary);
  max-width: 560px;
  line-height: 1.7;
}

/* ===== RESPONSIVE ===== */
@media (max-width: 768px) {
  .home { padding: 1.5rem 1rem 3rem; }
  .hero { padding: 2.5rem 1.25rem 2rem; }
  .hero-cta { flex-direction: column; width: 100%; }
  .cta-primary, .cta-secondary { width: 100%; justify-content: center; }
  .split { grid-template-columns: 1fr; }
  .services-row { grid-template-columns: 1fr; }
  .steps { flex-direction: column; align-items: flex-start; gap: 1rem; padding: 1rem; }
  .step { flex-direction: row; align-items: center; }
  .step::after { display: none !important; }
  .carousel-strip img { height: 160px; }
}
</style>
