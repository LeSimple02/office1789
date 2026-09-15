<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(() => {
  setTimeout(() => {
    router.push('/account')
  }, 3000)
})
</script>

<template>
  <div class="success-wrap">
    <div class="success-card">
      <div class="check-circle">✓</div>
      <h2>Abonnement activé</h2>
      <p>Redirection vers votre compte...</p>
      <div class="progress-bar"><div class="progress-fill"></div></div>
    </div>
  </div>
</template>

<style scoped>
.success-wrap {
  display: flex; align-items: center; justify-content: center;
  min-height: 70vh; padding: 2rem 1.25rem;
  animation: fadeInUp .5s var(--ease) both;
}
.success-card {
  text-align: center; padding: 3rem 2.5rem;
  background: var(--bg-primary); border: 1px solid var(--border);
  border-radius: var(--r-2xl); box-shadow: var(--sh-lg);
  animation: popIn .4s var(--ease-spring) both;
  max-width: 400px; width: 100%;
}
.check-circle {
  width: 64px; height: 64px; margin: 0 auto 1.5rem;
  border-radius: 50%; background: var(--ink); color: var(--paper);
  display: flex; align-items: center; justify-content: center;
  font-size: 2rem; font-weight: 700;
  animation: popIn .3s var(--ease-spring) .2s both;
}
.success-card h2 { font-size: 1.5rem; font-weight: 800; letter-spacing: -.02em; margin-bottom: .5rem; }
.success-card p { font-size: .9375rem; color: var(--text-secondary); margin-bottom: 1.5rem; }
.progress-bar { height: 3px; background: var(--bg-tertiary); border-radius: var(--r-full); overflow: hidden; }
.progress-fill { height: 100%; background: var(--ink); border-radius: var(--r-full); animation: growBar 3s linear both; }
</style>
