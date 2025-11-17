<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(() => {
  // Rediriger vers le compte après 3 secondes
  setTimeout(() => {
    router.push('/account')
  }, 3000)
})
</script>

<template>
  <div class="success-container">
    <div class="success-card">
      <div class="success-icon">
        <svg width="80" height="80" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10" stroke="#28a745" fill="#d4edda"></circle>
          <polyline points="9 12 11 14 15 10" stroke="#28a745" stroke-width="3"></polyline>
        </svg>
      </div>
      
      <h1 class="success-title">Paiement réussi !</h1>
      <p class="success-message">
        Votre abonnement a été activé avec succès. Vous allez être redirigé vers votre compte.
      </p>
      
      <div class="success-details">
        <p>✅ Paiement validé</p>
        <p>✅ Abonnement activé</p>
        <p>✅ Accès complet à vos fonctionnalités</p>
      </div>

      <button class="btn-primary" @click="router.push('/account')">
        Retour au compte
      </button>
    </div>
  </div>
</template>

<style scoped>
.success-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.success-card {
  background: white;
  border-radius: 20px;
  padding: 60px 40px;
  max-width: 500px;
  width: 100%;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slideIn 0.5s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.success-icon {
  margin-bottom: 30px;
  animation: scaleIn 0.6s ease-out 0.2s both;
}

@keyframes scaleIn {
  from {
    transform: scale(0);
  }
  to {
    transform: scale(1);
  }
}

.success-title {
  font-size: 32px;
  font-weight: 700;
  color: #1a1d2e;
  margin-bottom: 16px;
}

.success-message {
  font-size: 16px;
  color: #6c757d;
  margin-bottom: 40px;
  line-height: 1.6;
}

.success-details {
  background: #f8f9fa;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 32px;
  text-align: left;
}

.success-details p {
  font-size: 15px;
  color: #28a745;
  margin: 12px 0;
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.btn-primary {
  background: linear-gradient(135deg, #007bff 0%, #0056b3 100%);
  color: white;
  border: none;
  padding: 16px 40px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 123, 255, 0.3);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 123, 255, 0.4);
}

@media (max-width: 768px) {
  .success-card {
    padding: 40px 24px;
  }
  
  .success-title {
    font-size: 26px;
  }
}
</style>
