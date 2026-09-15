<script setup>
</script>

<template>
  <div class="page-wrap">
    <div class="head">
      <h1>{{ $t('contactUs') }}</h1>
      <p>{{ $t('hereToHelp') }}</p>
    </div>
    <div class="card">
      <p class="label">{{ $t('contactEmailTitle') }}</p>
      <a class="email" href="mailto:contact@office1789.com">contact@office1789.com</a>
      <p class="note">{{ $t('contactEmailNote') }}</p>
    </div>
  </div>
</template>

<style scoped>
.page-wrap { max-width: 520px; margin: 0 auto; padding: 3rem 1.25rem 5rem; animation: fadeInUp .5s var(--ease) both; }
.head { text-align: center; margin-bottom: 1.5rem; animation: fadeInDown .4s var(--ease) both; }
.head h1 { font-size: 1.75rem; font-weight: 800; letter-spacing: -.025em; margin-bottom: .375rem; }
.head p { font-size: .9375rem; color: var(--text-secondary); }
.card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 2rem 1.75rem; text-align: center; box-shadow: var(--sh-sm); animation: fadeInUp .4s var(--ease) .1s both; }
.label { font-size: .8125rem; font-weight: 600; color: var(--text-secondary); margin-bottom: .75rem; letter-spacing: .01em; }
.email { display: inline-block; font-size: 1.125rem; font-weight: 700; color: var(--text-primary); margin-bottom: 1rem; transition: opacity var(--tb); }
.email:hover { opacity: .7; }
.note { font-size: .875rem; color: var(--text-muted); line-height: 1.6; }
</style>
