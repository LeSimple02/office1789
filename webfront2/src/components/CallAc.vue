<script setup>
</script>

<template>
  <div class="call-area">
    <img src="@/assets/napo.png" alt="Call" class="call-img" />
    <p class="call-text">{{ $t('callt') }}...</p>
  </div>
</template>

<style scoped>
.call-area {
  position: fixed; inset: 0; z-index: 9999;
  display: flex; flex-direction: column;
  align-items: center; justify-content: center;
  gap: 1.5rem; text-align: center;
  background: var(--bg-primary);
  animation: fadeIn .3s var(--ease) both;
}
.call-img {
  width: 80px; height: 80px; border-radius: 50%;
  border: 2px solid var(--border-strong);
  animation: pulseScale 1.5s var(--ease) infinite;
}
.call-text { font-size: 1.125rem; font-weight: 600; color: var(--text-primary); }
</style>
