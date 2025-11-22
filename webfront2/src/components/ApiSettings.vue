<template>
  <div class="api-settings">
    <h3>⚙️ Configuration API</h3>
    
    <div class="current-url">
      <label>URL actuelle :</label>
      <code>{{ currentURL }}</code>
      <button @click="testCurrentConnection" :disabled="testing">
        {{ testing ? '⏳' : '🔍' }} Test
      </button>
      <span v-if="connectionStatus !== null" :class="connectionClass">
        {{ connectionStatus ? '✅ OK' : '❌ Erreur' }}
      </span>
    </div>

    <div class="presets">
      <h4>Préréglages :</h4>
      <button 
        v-for="(url, name) in presets" 
        :key="name"
        @click="selectPreset(name, url)"
        :class="{ active: currentURL === url }"
      >
        {{ name }}
      </button>
    </div>

    <div class="custom-url">
      <h4>URL personnalisée :</h4>
      <input 
        v-model="customURL" 
        type="text" 
        placeholder="http://192.168.1.xxx:8080"
      />
      <button @click="setCustomURL" :disabled="!customURL">
        Appliquer
      </button>
    </div>

    <div class="ip-helper">
      <details>
        <summary>💡 Comment trouver mon IP locale ?</summary>
        <p><strong>Windows:</strong> Ouvrez PowerShell et tapez <code>ipconfig</code></p>
        <p><strong>Mac/Linux:</strong> Ouvrez Terminal et tapez <code>ifconfig</code></p>
        <p>Cherchez "Adresse IPv4" ou "inet" (ex: 192.168.1.100)</p>
      </details>
    </div>
  </div>
</template>

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

<style scoped>
.api-settings {
  padding: 20px;
  max-width: 600px;
  margin: 0 auto;
}

.current-url {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  padding: 15px;
  background: #f5f5f5;
  border-radius: 8px;
}

.current-url code {
  flex: 1;
  padding: 5px 10px;
  background: white;
  border-radius: 4px;
  font-size: 14px;
}

.current-url button {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  background: #007bff;
  color: white;
  cursor: pointer;
}

.current-url button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.success {
  color: #28a745;
  font-weight: bold;
}

.error {
  color: #dc3545;
  font-weight: bold;
}

.presets {
  margin: 20px 0;
}

.presets button {
  margin: 5px;
  padding: 10px 20px;
  border: 2px solid #ddd;
  border-radius: 6px;
  background: white;
  cursor: pointer;
  transition: all 0.2s;
}

.presets button:hover {
  border-color: #007bff;
  background: #f0f8ff;
}

.presets button.active {
  border-color: #007bff;
  background: #007bff;
  color: white;
}

.custom-url {
  margin: 20px 0;
}

.custom-url input {
  width: 70%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.custom-url button {
  margin-left: 10px;
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  background: #28a745;
  color: white;
  cursor: pointer;
}

.ip-helper {
  margin-top: 20px;
  padding: 15px;
  background: #fff3cd;
  border-radius: 8px;
}

.ip-helper summary {
  cursor: pointer;
  font-weight: bold;
  color: #856404;
}

.ip-helper p {
  margin: 10px 0 5px;
  color: #856404;
}

.ip-helper code {
  background: white;
  padding: 2px 6px;
  border-radius: 3px;
}
</style>
