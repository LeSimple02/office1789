<script setup>
</script>

<template>
  <div class="page-wrap">
    <div class="head">
      <h1>Mentions Légales</h1>
      <p>Informations légales et réglementaires</p>
    </div>
    <div class="legal-body" v-html="$t('mentions')"></div>
    <div class="grid">
      <div class="legal-card" v-for="s in [
        { icon: '🏢', title: 'Éditeur', lines: ['Office1789 SAS', 'Capital social : 100 000€', 'SIRET : 123 456 789 00012'] },
        { icon: '🌐', title: 'Hébergement', lines: ['OVH Cloud', '2 rue Kellermann', '59100 Roubaix, France'] },
        { icon: '🔒', title: 'Données personnelles', lines: ['Conformité RGPD', 'Protection des données', 'Droit d\'accès et rectification'] }
      ]" :key="s.title">
        <div class="legal-icon">{{ s.icon }}</div>
        <h3>{{ s.title }}</h3>
        <p v-for="l in s.lines" :key="l">{{ l }}</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.page-wrap { max-width: 760px; margin: 0 auto; padding: 3rem 1.25rem 5rem; animation: fadeInUp .5s var(--ease) both; }
.head { text-align: center; margin-bottom: 2rem; animation: fadeInDown .4s var(--ease) both; }
.head h1 { font-size: 1.75rem; font-weight: 800; letter-spacing: -.025em; margin-bottom: .375rem; }
.head p { font-size: .9375rem; color: var(--text-secondary); }
.legal-body { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 2rem 1.75rem; box-shadow: var(--sh-sm); font-size: .9375rem; color: var(--text-secondary); line-height: 1.7; margin-bottom: 2rem; animation: fadeInUp .4s var(--ease) .1s both; }
.legal-body :deep(h2) { font-size: 1.125rem; font-weight: 700; color: var(--text-primary); margin: 1.5rem 0 .75rem; }
.legal-body :deep(h3) { font-size: 1rem; font-weight: 600; color: var(--text-primary); margin: 1.125rem 0 .5rem; }
.legal-body :deep(p) { margin-bottom: .75rem; }
.legal-body :deep(a) { color: var(--text-primary); text-decoration: underline; text-underline-offset: 2px; }
.legal-body :deep(ul) { padding-left: 1.25rem; margin-bottom: .75rem; }
.grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1.25rem; }
.legal-card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1.75rem 1.5rem; animation: fadeInUp .4s var(--ease) both; transition: border-color var(--tb), box-shadow var(--tb); }
.legal-card:nth-child(2) { animation-delay: 80ms; }
.legal-card:nth-child(3) { animation-delay: 160ms; }
.legal-card:hover { border-color: var(--border-strong); box-shadow: var(--sh-sm); }
.legal-icon { font-size: 1.5rem; margin-bottom: .75rem; width: 48px; height: 48px; display: flex; align-items: center; justify-content: center; border-radius: var(--r); background: var(--bg-tertiary); }
.legal-card h3 { font-size: 1.0625rem; font-weight: 700; margin-bottom: .5rem; }
.legal-card p { font-size: .875rem; color: var(--text-secondary); line-height: 1.6; margin-top: .25rem; }
@media (max-width: 768px) { .grid { grid-template-columns: 1fr; } }
</style>
