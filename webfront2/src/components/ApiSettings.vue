<script>
import { API_CONFIG, testConnection } from '@/config/api.js'

export default {
  name: 'ApiSettings',
  data() {
    return {
      currentURL: API_CONFIG.getBaseURL(),
      customURL: '',
      presets: API_CONFIG.presets,
      testing: false,
      connectionStatus: null
    }
  },
  computed: {
    connectionClass() {
      return this.connectionStatus ? 'success' : 'error'
    }
  },
  methods: {
    async testCurrentConnection() {
      this.testing = true
      this.connectionStatus = await testConnection(this.currentURL)
      this.testing = false
    },
    selectPreset(name, url) {
      if (name === 'wifi') {
        const ip = prompt('Entrez votre IP locale (ex: 192.168.1.100):')
        if (ip) {
          const customUrl = `http://${ip}:8080`
          API_CONFIG.setBaseURL(customUrl)
        }
      } else {
        API_CONFIG.setBaseURL(url)
      }
    },
    setCustomURL() {
      if (this.customURL) {
        API_CONFIG.setBaseURL(this.customURL)
      }
    }
  }
}
</script>

<template>
  <div class="settings">
    <div class="setting-row">
      <label class="label">URL actuelle</label>
      <code class="url-display">{{ currentURL }}</code>
      <button @click="testCurrentConnection" class="btn-test" :disabled="testing">
        <span v-if="testing" class="spinner"></span>
        {{ testing ? 'Test...' : 'Tester' }}
      </button>
    </div>
    <div v-if="connectionStatus !== null" class="status" :class="connectionClass">
      {{ connectionStatus ? 'Connecté' : 'Échec de connexion' }}
    </div>
    <div class="presets">
      <label class="label">Préréglages</label>
      <div class="preset-list">
        <button v-for="(url, name) in presets" :key="name" @click="selectPreset(name, url)" class="preset-btn">
          {{ name }}
        </button>
      </div>
    </div>
    <div class="custom-row">
      <label class="label">URL personnalisée</label>
      <div class="custom-input">
        <input v-model="customURL" type="text" placeholder="https://api.example.com" />
        <button @click="setCustomURL" class="btn-set">Définir</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings { display: flex; flex-direction: column; gap: 1.25rem; }
.setting-row { display: flex; align-items: center; gap: .75rem; flex-wrap: wrap; }
.label { font-size: .8125rem; font-weight: 600; color: var(--text-secondary); }
.url-display { padding: .375rem .75rem; font-size: .8125rem; font-family: var(--font-mono); background: var(--bg-tertiary); border: 1px solid var(--border); border-radius: var(--r-sm); flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; }
.btn-test { padding: .375rem .875rem; font-size: .8125rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r-sm); cursor: pointer; transition: all var(--tb); display: flex; align-items: center; gap: .375rem; }
.btn-test:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); }
.btn-test:disabled { opacity: .5; cursor: not-allowed; }
.status { padding: .5rem .875rem; font-size: .8125rem; font-weight: 500; border-radius: var(--r); border: 1px solid var(--border); }
.status.success { background: var(--bg-tertiary); color: var(--text-primary); }
.status.error { background: var(--bg-tertiary); color: var(--text-primary); }
.presets { display: flex; flex-direction: column; gap: .5rem; }
.preset-list { display: flex; gap: .5rem; flex-wrap: wrap; }
.preset-btn { padding: .375rem .875rem; font-size: .8125rem; font-weight: 500; background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r-sm); cursor: pointer; transition: all var(--tb); }
.preset-btn:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.custom-row { display: flex; flex-direction: column; gap: .5rem; }
.custom-input { display: flex; gap: .5rem; }
.custom-input input { flex: 1; padding: .5rem .75rem; font-size: .875rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); }
.custom-input input:focus { outline: none; border-color: var(--ink); }
.btn-set { padding: .5rem 1rem; font-size: .8125rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-set:hover { background: var(--ink-soft); border-color: var(--ink-soft); }
.spinner { width: 14px; height: 14px; border: 2px solid var(--border); border-top-color: var(--paper); border-radius: 50%; animation: spin .6s linear infinite; display: inline-block; }
</style>
