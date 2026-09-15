<script setup lang="ts">
</script>

<template>
  <div class="about">
    <!-- HERO -->
    <section class="hero">
      <div class="hero-badge">À propos</div>
      <h1>{{ $t('aboutOffice1789') }}</h1>
      <p>{{ $t('discoverCollaborativeSuite') }}</p>
    </section>

    <!-- MISSION -->
    <section class="block">
      <h2 class="block-title">{{ $t('ourMission') }}</h2>
      <div class="mission-grid">
        <div class="m-card" v-for="(m, i) in [
          { icon: '🚀', title: $t('innovation'), desc: $t('developModernTools') },
          { icon: '🇫🇷', title: $t('sovereignty'), desc: $t('guaranteeDataHosting') },
          { icon: '🔒', title: $t('confidentiality'), desc: $t('protectDataEndToEnd') },
          { icon: '♻️', title: $t('ecoResponsible'), desc: $t('useGreenServers') }
        ]" :key="i">
          <div class="m-icon">{{ m.icon }}</div>
          <h3>{{ m.title }}</h3>
          <p>{{ m.desc }}</p>
        </div>
      </div>
    </section>

    <!-- VALUES -->
    <section class="block">
      <h2 class="block-title">{{ $t('ourValues') }}</h2>
      <div class="values">
        <div class="value" v-for="(v, i) in [
          { n: '01', title: $t('transparency'), desc: $t('openCommunication') },
          { n: '02', title: $t('accessibility'), desc: $t('fairPricing') },
          { n: '03', title: $t('performance'), desc: $t('fastReliableTools') },
          { n: '04', title: $t('continuousInnovation'), desc: $t('constantlyImprove') }
        ]" :key="i">
          <span class="v-num">{{ v.n }}</span>
          <div class="v-body">
            <h3>{{ v.title }}</h3>
            <p>{{ v.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- WHY -->
    <section class="block">
      <h2 class="block-title">{{ $t('whyOffice1789') }}</h2>
      <div class="why-grid">
        <div class="why-card">
          <h3>{{ $t('statistics') }}</h3>
          <div class="stats">
            <div class="stat" v-for="(s, i) in [
              { val: '99.9%', label: $t('availability') },
              { val: '50 Go', label: $t('storage') },
              { val: '100%', label: $t('french') }
            ]" :key="i">
              <span class="stat-val">{{ s.val }}</span>
              <span class="stat-lbl">{{ s.label }}</span>
            </div>
          </div>
        </div>
        <div class="why-card">
          <h3>{{ $t('advantages') }}</h3>
          <ul class="adv-list">
            <li v-for="(a, i) in [$t('noAds'), $t('openSource'), $t('reactiveSupport'), $t('easyMigration'), $t('mobileApps'), $t('thirdPartyIntegrations')]" :key="i">{{ a }}</li>
          </ul>
        </div>
      </div>
    </section>

    <!-- CONTENT -->
    <div class="content-card" v-html="$t('aboutd')"></div>

    <!-- OFFERS -->
    <section class="block">
      <h2 class="block-title">{{ $t('ourOffers') }}</h2>
      <div v-html="$t('offer')" class="offers"></div>
    </section>

    <!-- CTA -->
    <section class="cta">
      <h3>{{ $t('readyToJoin') }}</h3>
      <p>{{ $t('createAccountEnjoyServices') }}</p>
      <RouterLink to="/createaccount" class="cta-btn">{{ $t('createac') }}</RouterLink>
    </section>
  </div>
</template>

<style scoped>
.about { max-width: var(--content-max); margin: 0 auto; padding: 2.5rem 1.5rem 5rem; animation: fadeInUp .5s var(--ease) both; }

/* HERO */
.hero { text-align: center; padding: 3rem 2rem; background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--r-2xl); margin-bottom: 2.5rem; animation: fadeInUp .5s var(--ease) both; }
.hero-badge { display: inline-block; padding: .3rem .7rem; font-size: .75rem; font-weight: 700; border-radius: var(--r-full); background: var(--bg-tertiary); color: var(--text-secondary); margin-bottom: 1rem; animation: fadeInDown .4s var(--ease) both; }
.hero h1 { font-size: clamp(1.75rem, 4vw, 2.5rem); font-weight: 800; letter-spacing: -.03em; margin-bottom: .5rem; animation: fadeInUp .5s var(--ease) .1s both; }
.hero p { font-size: 1rem; color: var(--text-secondary); animation: fadeInUp .5s var(--ease) .2s both; }

/* BLOCK */
.block { margin-bottom: 2.5rem; animation: fadeInUp .5s var(--ease) both; }
.block-title { font-size: 1.5rem; font-weight: 800; letter-spacing: -.025em; margin-bottom: 1.5rem; }

/* MISSION */
.mission-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 1.25rem; }
.m-card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1.5rem; text-align: center; transition: transform var(--tb), border-color var(--tb), box-shadow var(--tb); animation: fadeInUp .4s var(--ease) both; }
.m-card:nth-child(2) { animation-delay: 60ms; }
.m-card:nth-child(3) { animation-delay: 120ms; }
.m-card:nth-child(4) { animation-delay: 180ms; }
.m-card:hover { border-color: var(--border-strong); box-shadow: var(--sh-md); transform: translateY(-4px); }
.m-icon { width: 48px; height: 48px; margin: 0 auto .875rem; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; border-radius: var(--r); background: var(--bg-tertiary); }
.m-card h3 { font-size: 1rem; font-weight: 700; margin-bottom: .375rem; }
.m-card p { font-size: .8125rem; color: var(--text-secondary); line-height: 1.6; }

/* VALUES */
.values { display: flex; flex-direction: column; }
.value { display: flex; align-items: flex-start; gap: 1.25rem; padding: 1.25rem 0; border-bottom: 1px solid var(--border); animation: fadeInLeft .4s var(--ease) both; }
.value:last-child { border-bottom: none; }
.v-num { font-size: 1.75rem; font-weight: 800; letter-spacing: -.03em; color: var(--text-muted); flex-shrink: 0; min-width: 48px; line-height: 1; }
.v-body h3 { font-size: 1.0625rem; font-weight: 700; margin-bottom: .25rem; }
.v-body p { font-size: .9375rem; color: var(--text-secondary); line-height: 1.6; }

/* WHY */
.why-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; }
.why-card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1.75rem; animation: fadeInUp .4s var(--ease) both; }
.why-card:nth-child(2) { animation-delay: 80ms; }
.why-card h3 { font-size: 1.125rem; font-weight: 700; margin-bottom: 1.25rem; }
.stats { display: flex; gap: 1.5rem; }
.stat { text-align: center; flex: 1; }
.stat-val { display: block; font-size: 1.75rem; font-weight: 800; letter-spacing: -.02em; animation: countUp .6s var(--ease) both; }
.stat-lbl { display: block; font-size: .8125rem; color: var(--text-muted); margin-top: .25rem; }
.adv-list { list-style: none; padding: 0; margin: 0; }
.adv-list li { font-size: .9375rem; color: var(--text-secondary); margin-bottom: .625rem; display: flex; align-items: center; gap: .5rem; animation: fadeIn .3s var(--ease) both; }
.adv-list li::before { content: ''; width: 6px; height: 6px; border-radius: 50%; background: var(--ink); flex-shrink: 0; }

/* CONTENT */
.content-card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 2rem 1.75rem; margin-bottom: 2.5rem; line-height: 1.7; font-size: .9375rem; color: var(--text-secondary); animation: fadeInUp .5s var(--ease) both; }
.content-card :deep(h2) { font-size: 1.25rem; color: var(--text-primary); margin: 1.25rem 0 .75rem; }
.content-card :deep(h3) { font-size: 1.0625rem; color: var(--text-primary); margin: 1rem 0 .5rem; }

/* OFFERS */
.offers { line-height: 1.8; font-size: .9375rem; color: var(--text-secondary); }
.offers :deep(table) { width: 100%; border-collapse: collapse; margin-top: 1rem; }
.offers :deep(th), .offers :deep(td) { padding: .75rem 1rem; border: 1px solid var(--border); font-size: .875rem; }
.offers :deep(th) { background: var(--bg-secondary); font-weight: 700; text-align: left; }

/* CTA */
.cta { text-align: center; padding: 3rem 2rem; background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--r-2xl); animation: fadeInUp .5s var(--ease) both; }
.cta h3 { font-size: 1.5rem; font-weight: 800; letter-spacing: -.02em; margin-bottom: .5rem; }
.cta p { font-size: 1rem; color: var(--text-secondary); margin-bottom: 1.5rem; }
.cta-btn { display: inline-block; padding: .75rem 2rem; font-size: .9375rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); transition: all var(--tb); }
.cta-btn:hover { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-md); transform: translateY(-2px); }

/* RESPONSIVE */
@media (max-width: 768px) {
  .mission-grid { grid-template-columns: 1fr 1fr; }
  .why-grid { grid-template-columns: 1fr; }
  .stats { flex-direction: column; gap: 1rem; }
}
@media (max-width: 480px) { .mission-grid { grid-template-columns: 1fr; } }
</style>
