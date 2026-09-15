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

onMounted(async () => { await loadStorageInfo() })

async function loadStorageInfo() {
  loading.value = true
  error.value = null
  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/drive/storage-info`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: gls().username, token: gls().sessionT })
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
  const names = { 0: 'Free', 1: 'Standard', 2: 'Professional', 3: 'Enterprise' }
  return names[nboffer] || 'Free'
}

defineExpose({ refresh: loadStorageInfo })
</script>

<template>
  <div class="quota">
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <span>Chargement...</span>
    </div>
    <div v-else-if="error" class="err">{{ error }}</div>
    <div v-else class="quota-body">
      <div class="quota-head">
        <span class="quota-label">Stockage</span>
        <span class="quota-offer">{{ getOfferName(storageInfo.nboffer) }}</span>
      </div>
      <div class="bar-bg">
        <div class="bar-fill" :style="{ width: progressWidth + '%' }"></div>
      </div>
      <div class="quota-stats">
        <span>{{ formatBytes(storageInfo.used_bytes || 0) }}</span>
        <span v-if="storageInfo.unlimited">Illimité</span>
        <span v-else>{{ formatBytes(storageInfo.total_bytes || 0) }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.quota { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1.25rem; animation: fadeInUp .4s var(--ease) both; }
.loading { display: flex; align-items: center; gap: .625rem; font-size: .8125rem; color: var(--text-muted); }
.spinner { width: 16px; height: 16px; border: 2px solid var(--border); border-top-color: var(--ink); border-radius: 50%; animation: spin .7s linear infinite; }
.err { font-size: .8125rem; color: var(--text-primary); }
.quota-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: .75rem; }
.quota-label { font-size: .8125rem; font-weight: 600; color: var(--text-secondary); }
.quota-offer { font-size: .75rem; font-weight: 600; padding: .2rem .6rem; border-radius: var(--r-full); background: var(--bg-tertiary); color: var(--text-secondary); }
.bar-bg { height: 6px; background: var(--bg-tertiary); border-radius: var(--r-full); overflow: hidden; }
.bar-fill { height: 100%; background: var(--ink); border-radius: var(--r-full); transition: width var(--ts); animation: growBar .6s var(--ease) both; }
.quota-stats { display: flex; justify-content: space-between; margin-top: .5rem; font-size: .75rem; color: var(--text-muted); }
</style>
