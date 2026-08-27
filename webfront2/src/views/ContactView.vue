<template>
  <div class="contact-view">
    <div class="hero-card">
      <div class="contact-icon-wrapper">
        <svg class="contact-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
          />
        </svg>
      </div>
      <h1>{{ $t('contactUs') }}</h1>
      <h2>{{ $t('hereToHelp') }}</h2>
    </div>

    <div class="contact-card">
      <p class="contact-label">{{ $t('contactEmailTitle') }}</p>
      <a class="contact-email" href="mailto:contact@office1789.com">contact@office1789.com</a>
      <p class="contact-note">{{ $t('contactEmailNote') }}</p>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.contact-view {
  width: 100%;
  animation: fadeIn 0.5s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.hero-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 220px;
  width: 100%;
  padding: 40px 24px;
  background: linear-gradient(135deg, #f5f5ff, #ffecef);
  border-radius: 24px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
  gap: 14px;
}

.dark .hero-card {
  background: linear-gradient(135deg, #111827, #1c1c1e);
}

.contact-icon-wrapper {
  width: 72px;
  height: 72px;
  background: linear-gradient(30deg, #1f70ff, #ff4d4f);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%,
  100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.contact-icon {
  width: 36px;
  height: 36px;
  color: white;
}

.hero-card h1 {
  font-family: roboto, sans-serif;
  font-size: 2.1rem;
  font-weight: 700;
  text-align: center;
  color: #222;
  margin: 0;
}

.dark .hero-card h1 {
  color: #fff;
}

.hero-card h2 {
  font-family: roboto, sans-serif;
  font-size: 1.05rem;
  text-align: center;
  background: linear-gradient(30deg, #1f70ff, #ff4d4f);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  margin: 0;
}

.contact-card {
  background: #ffffff;
  border-radius: 20px;
  padding: 32px 24px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
  display: grid;
  gap: 10px;
  text-align: center;
}

.dark .contact-card {
  background: #111827;
}

.contact-label {
  font-size: 0.95rem;
  letter-spacing: 0.3px;
  color: #666;
  margin: 0;
}

.dark .contact-label {
  color: #aaa;
}

.contact-email {
  font-size: 1.4rem;
  font-weight: 700;
  color: #1f70ff;
  text-decoration: none;
  word-break: break-all;
}

.contact-email:hover {
  text-decoration: underline;
}

.dark .contact-email {
  color: #6ea8ff;
}

.contact-note {
  font-size: 0.95rem;
  color: #444;
  margin: 0;
}

.dark .contact-note {
  color: #ccc;
}

@media (max-width: 768px) {
  .hero-card {
    padding: 28px 18px;
  }

  .hero-card h1 {
    font-size: 1.8rem;
  }

  .contact-email {
    font-size: 1.1rem;
  }
}
</style>
