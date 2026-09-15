<script setup>
import MailAccessSSO from "@/components/MailAccessSSO.vue"
</script>

<template>
  <div class="page-wrap">
    <div class="hero">
      <div class="hero-badge">Messagerie</div>
      <h1>Vos emails, simplement</h1>
      <p>Sécurisé, rapide et sans tracking publicitaire.</p>
      <MailAccessSSO />
    </div>
    <div class="grid">
      <div class="feat" v-for="f in [
        { icon: '📬', title: 'Sécurisé', desc: 'Chiffrement de bout en bout' },
        { icon: '⚡', title: 'Rapide', desc: 'Interface ultra-réactive' },
        { icon: '🔒', title: 'Privé', desc: 'Aucun tracking publicitaire' }
      ]" :key="f.title">
        <div class="feat-icon">{{ f.icon }}</div>
        <h3>{{ f.title }}</h3>
        <p>{{ f.desc }}</p>
      </div>
    </div>
    <div class="note-box">La messagerie s'ouvrira dans une nouvelle fenêtre.</div>
  </div>
</template>

<style scoped>
.page-wrap { max-width: var(--content-max); margin: 0 auto; padding: 2.5rem 1.5rem 5rem; animation: fadeInUp .5s var(--ease) both; }
.hero { text-align: center; padding: 3rem 2rem; background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--r-2xl); margin-bottom: 2rem; animation: fadeInUp .5s var(--ease) both; }
.hero-badge { display: inline-block; padding: .3rem .7rem; font-size: .75rem; font-weight: 700; border-radius: var(--r-full); background: var(--bg-tertiary); color: var(--text-secondary); margin-bottom: 1rem; animation: fadeInDown .4s var(--ease) both; }
.hero h1 { font-size: clamp(1.75rem, 4vw, 2.5rem); font-weight: 800; letter-spacing: -.03em; margin-bottom: .5rem; animation: fadeInUp .5s var(--ease) .1s both; }
.hero p { font-size: 1rem; color: var(--text-secondary); margin-bottom: 1.5rem; animation: fadeInUp .5s var(--ease) .2s both; }
.grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1.25rem; margin-bottom: 2rem; }
.feat { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1.75rem 1.5rem; text-align: center; transition: transform var(--tb), border-color var(--tb), box-shadow var(--tb); animation: fadeInUp .4s var(--ease) both; }
.feat:nth-child(2) { animation-delay: 80ms; }
.feat:nth-child(3) { animation-delay: 160ms; }
.feat:hover { border-color: var(--border-strong); box-shadow: var(--sh-md); transform: translateY(-4px); }
.feat-icon { width: 48px; height: 48px; margin: 0 auto .875rem; display: flex; align-items: center; justify-content: center; font-size: 1.5rem; border-radius: var(--r); background: var(--bg-tertiary); animation: float 3s var(--ease) infinite; }
.feat:nth-child(2) .feat-icon { animation-delay: .3s; }
.feat:nth-child(3) .feat-icon { animation-delay: .6s; }
.feat h3 { font-size: 1.0625rem; font-weight: 700; margin-bottom: .375rem; }
.feat p { font-size: .875rem; color: var(--text-secondary); line-height: 1.6; }
.note-box { padding: 1rem 1.25rem; background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--r-lg); font-size: .875rem; color: var(--text-secondary); line-height: 1.6; animation: fadeInUp .4s var(--ease) .3s both; }
@media (max-width: 768px) { .grid { grid-template-columns: 1fr; } }
</style>
