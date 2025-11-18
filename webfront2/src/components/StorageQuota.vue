<template>
  <div class="storage-quota">
    <div v-if="loading" class="loading">Chargement...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="quota-display">
      <div class="quota-header">
        <h3>💾 Stockage</h3>
        <span class="offer-badge" :class="'offer-' + storageInfo.nboffer">
          {{ getOfferName(storageInfo.nboffer) }}
        </span>
      </div>
      
      <div class="quota-bar">
        <div class="quota-progress" :style="{ width: progressWidth + '%' }" :class="progressClass"></div>
      </div>
      
      <div class="quota-text">
        <span v-if="storageInfo.unlimited">
          {{ formatBytes(storageInfo.current_usage) }} utilisés
          <span class="unlimited-badge">✨ Illimité</span>
        </span>
        <span v-else>
          {{ formatBytes(storageInfo.current_usage) }} / {{ formatBytes(storageInfo.storage_limit) }}
          ({{ Math.round(storageInfo.percentage_used) }}%)
        </span>
      </div>
      
      <div v-if="!storageInfo.unlimited && storageInfo.remaining < 0" class="quota-warning">
        ⚠️ Quota dépassé ! Supprimez des fichiers pour libérer de l'espace.
      </div>
      <div v-else-if="!storageInfo.unlimited && storageInfo.percentage_used > 90" class="quota-warning">
        ⚠️ Espace presque plein. Pensez à améliorer votre offre.
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { gls } from '@/stores/global'

const storageInfo = ref({})
const loading = ref(true)
const error = ref(null)

const progressWidth = computed(() => {
  if (storageInfo.value.unlimited) return 50
  const percent = storageInfo.value.percentage_used || 0
  return Math.min(percent, 100)
})

const progressClass = computed(() => {
  const percent = storageInfo.value.percentage_used || 0
  if (percent >= 95) return 'critical'
  if (percent >= 80) return 'warning'
  return 'normal'
})

onMounted(async () => {
  await loadStorageInfo()
})

async function loadStorageInfo() {
  loading.value = true
  error.value = null
  
  try {
    const response = await fetch('http://localhost:8080/api/drive/storage', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: gls().username,
        token: gls().sessionT
      })
    })
    
    if (response.ok) {
      storageInfo.value = await response.json()
    } else {
      error.value = 'Erreur de chargement'
    }
  } catch (err) {
    error.value = 'Erreur réseau'
    console.error('Error loading storage info:', err)
  } finally {
    loading.value = false
  }
}

function formatBytes(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

function getOfferName(nboffer) {
  const names = {
    0: 'Free',
    1: 'Standard',
    2: 'Professional',
    3: 'Enterprise'
  }
  return names[nboffer] || 'Free'
}

// Expose function to parent to refresh after upload/delete
defineExpose({
  refresh: loadStorageInfo
})
</script>

<style scoped>
.storage-quota {
  background: white;
  border-radius: 12px;
  padding: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 1rem;
}

.loading, .error {
  text-align: center;
  padding: 1rem;
  color: #666;
}

.error {
  color: #ef4444;
}

.quota-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.quota-header h3 {
  margin: 0;
  font-size: 1.1rem;
  color: #333;
}

.offer-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 8px;
  font-weight: 600;
  font-size: 0.75rem;
}

.offer-badge.offer-0 {
  background: #e5e7eb;
  color: #374151;
}

.offer-badge.offer-1 {
  background: #dbeafe;
  color: #1e40af;
}

.offer-badge.offer-2 {
  background: #ede9fe;
  color: #5b21b6;
}

.offer-badge.offer-3 {
  background: #fef3c7;
  color: #92400e;
}

.quota-bar {
  width: 100%;
  height: 20px;
  background: #e5e7eb;
  border-radius: 10px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.quota-progress {
  height: 100%;
  transition: width 0.3s ease, background-color 0.3s ease;
  border-radius: 10px;
}

.quota-progress.normal {
  background: linear-gradient(90deg, #10b981, #059669);
}

.quota-progress.warning {
  background: linear-gradient(90deg, #f59e0b, #d97706);
}

.quota-progress.critical {
  background: linear-gradient(90deg, #ef4444, #dc2626);
}

.quota-text {
  font-size: 0.9rem;
  color: #666;
  text-align: center;
}

.unlimited-badge {
  display: inline-block;
  margin-left: 0.5rem;
  padding: 0.15rem 0.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 600;
}

.quota-warning {
  margin-top: 0.75rem;
  padding: 0.75rem;
  background: #fef3c7;
  border-left: 4px solid #f59e0b;
  border-radius: 8px;
  font-size: 0.85rem;
  color: #92400e;
}
</style>
