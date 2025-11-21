<template>
  <div class="custom-domain-container">
    <!-- Hero Section -->
    <div class="hero-section">
      <router-link to="/account" class="back-link">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="19" y1="12" x2="5" y2="12"></line>
          <polyline points="12 19 5 12 12 5"></polyline>
        </svg>
        Retour
      </router-link>
      
      <div class="hero-icon-wrapper">
        <div class="icon-glow"></div>
        <svg class="hero-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="2" y1="12" x2="22" y2="12"></line>
          <path d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"></path>
        </svg>
      </div>
      
      <h1 class="hero-title">🌐 Domaine personnalisé</h1>
      <p class="hero-subtitle">Utilisez votre propre domaine pour vos emails professionnels</p>
      <div class="plan-badge-hero">Professional & Enterprise uniquement</div>
    </div>

    <!-- Current domain status -->
    <div v-if="currentDomain" class="domain-status">
      <div :class="['status-card', domainVerified ? 'verified' : 'pending']">
        <div class="domain-info">
          <h3>{{ currentDomain }}</h3>
          <span v-if="domainVerified" class="badge-verified">✅ Vérifié</span>
          <span v-else class="badge-pending">⏳ En attente de vérification</span>
        </div>

        <!-- Verification instructions -->
        <div v-if="!domainVerified && verificationToken" class="verification-instructions">
          <h4>📋 Instructions de vérification</h4>
          <p>Ajoutez cet enregistrement TXT à votre DNS :</p>
          
          <div class="dns-record">
            <div class="record-details">
              <div class="record-field">
                <label>Type :</label>
                <code>TXT</code>
              </div>
              <div class="record-field">
                <label>Nom/Host :</label>
                <code>@</code> ou <code>{{ currentDomain }}</code>
              </div>
              <div class="record-field">
                <label>Valeur :</label>
                <code class="token-value">{{ verificationRecord }}</code>
                <button @click="copyToClipboard(verificationRecord)" class="btn-copy">📋 Copier</button>
              </div>
            </div>
          </div>

          <div class="help-text">
            <p><strong>⏱️ Propagation DNS :</strong> La propagation peut prendre de quelques minutes à 48 heures.</p>
            <p><strong>💡 Aide :</strong> Consultez la documentation de votre registrar (OVH, Gandi, Cloudflare, etc.)</p>
          </div>

          <button @click="verifyDomain" :disabled="verifying" class="btn-verify">
            {{ verifying ? '🔍 Vérification...' : '🔍 Vérifier le domaine' }}
          </button>
        </div>

        <!-- Actions -->
        <div class="domain-actions">
          <button v-if="domainVerified" @click="showRemoveConfirm = true" class="btn-remove">
            🗑️ Retirer le domaine
          </button>
          <button v-else @click="showRemoveConfirm = true" class="btn-cancel">
            ❌ Annuler
          </button>
        </div>
      </div>
    </div>

    <!-- Add domain form -->
    <div v-else class="add-domain-form">
      <h3>Ajouter un domaine personnalisé</h3>
      
      <form @submit.prevent="addDomain">
        <div class="form-group">
          <label for="domain">Nom de domaine</label>
          <input 
            type="text" 
            id="domain" 
            v-model="newDomain" 
            required 
            placeholder="exemple.com"
            pattern="^[a-zA-Z0-9][a-zA-Z0-9\-\.]{0,251}[a-zA-Z0-9]$"
            title="Domaine valide (ex: votreentreprise.com)"
          />
          <small>⚠️ N'incluez pas le sous-domaine (pas de www, mail, etc.)</small>
        </div>

        <div class="example-emails">
          <p><strong>Exemples d'adresses créées :</strong></p>
          <ul>
            <li v-if="newDomain">contact@{{ newDomain || 'exemple.com' }}</li>
            <li v-if="newDomain">support@{{ newDomain || 'exemple.com' }}</li>
            <li v-if="newDomain">{{ username }}@{{ newDomain || 'exemple.com' }}</li>
          </ul>
        </div>

        <button type="submit" :disabled="loading || !newDomain" class="btn-primary">
          {{ loading ? '⏳ Ajout...' : '✅ Ajouter le domaine' }}
        </button>
      </form>
    </div>

    <!-- Remove confirmation modal -->
    <div v-if="showRemoveConfirm" class="modal-overlay" @click.self="showRemoveConfirm = false">
      <div class="modal-content">
        <h3>⚠️ Confirmer la suppression</h3>
        <p>Êtes-vous sûr de vouloir retirer le domaine <strong>{{ currentDomain }}</strong> ?</p>
        <p class="warning-text">Toutes les adresses email utilisant ce domaine cesseront de fonctionner.</p>
        
        <div class="modal-actions">
          <button @click="showRemoveConfirm = false" class="btn-cancel">Annuler</button>
          <button @click="removeDomain" :disabled="loading" class="btn-confirm-remove">
            {{ loading ? '⏳ Suppression...' : '🗑️ Confirmer la suppression' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Messages -->
    <div v-if="message" :class="['message', messageType]">
      {{ message }}
    </div>

    <!-- DNS Configuration help -->
    <div class="help-section">
      <h3>❓ Configuration DNS requise</h3>
      <p>Pour utiliser votre domaine personnalisé, vous devez configurer les enregistrements DNS suivants :</p>
      
      <table class="dns-table">
        <thead>
          <tr>
            <th>Type</th>
            <th>Nom/Host</th>
            <th>Valeur</th>
            <th>Priorité</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td><code>MX</code></td>
            <td><code>@</code></td>
            <td><code>mail.office1789.com</code></td>
            <td><code>10</code></td>
          </tr>
          <tr>
            <td><code>TXT</code></td>
            <td><code>@</code></td>
            <td><code>v=spf1 include:office1789.com ~all</code></td>
            <td><code>-</code></td>
          </tr>
          <tr>
            <td><code>TXT</code></td>
            <td><code>_dmarc</code></td>
            <td><code>v=DMARC1; p=quarantine; rua=mailto:postmaster@office1789.com</code></td>
            <td><code>-</code></td>
          </tr>
        </tbody>
      </table>

      <p class="help-note">📚 Ces enregistrements garantissent la délivrabilité de vos emails et la sécurité de votre domaine.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { gls } from '@/stores/global.js';

const router = useRouter();
const loading = ref(false);
const verifying = ref(false);
const message = ref('');
const messageType = ref('success');
const showRemoveConfirm = ref(false);

const username = ref(gls().username || '');
const currentDomain = ref(null);
const domainVerified = ref(false);
const verificationToken = ref('');
const verificationRecord = ref('');
const newDomain = ref('');

// Fetch current domain info
async function fetchDomainInfo() {
  const user = gls().username;
  const token = gls().sessionT;

  if (!user || !token) {
    router.push('/login');
    return;
  }

  try {
    const response = await fetch('http://localhost:8080/api/domain/info', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: user, token })
    });

    if (response.ok) {
      const data = await response.json();
      if (data.success && data.domain) {
        currentDomain.value = data.domain;
        domainVerified.value = data.verified;
        verificationToken.value = data.verification_token || '';
        verificationRecord.value = data.verification_record || '';
      }
    }
  } catch (error) {
    console.error('Error fetching domain info:', error);
  }
}

// Add custom domain
async function addDomain() {
  if (loading.value || !newDomain.value) return;

  const user = gls().username;
  const token = gls().sessionT;

  loading.value = true;
  message.value = '';

  try {
    const response = await fetch('http://localhost:8080/api/domain/add', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: user,
        token,
        domain: newDomain.value
      })
    });

    const data = await response.json();

    if (data.success) {
      message.value = '✅ ' + data.message;
      messageType.value = 'success';
      currentDomain.value = data.domain;
      domainVerified.value = data.verified;
      verificationToken.value = data.verification_token || '';
      verificationRecord.value = data.verification_record || '';
      newDomain.value = '';
    } else {
      message.value = '❌ ' + data.message;
      messageType.value = 'error';
    }
  } catch (error) {
    message.value = `❌ Erreur réseau: ${error.message}`;
    messageType.value = 'error';
  } finally {
    loading.value = false;
  }
}

// Verify domain
async function verifyDomain() {
  if (verifying.value) return;

  const user = gls().username;
  const token = gls().sessionT;

  verifying.value = true;
  message.value = '';

  try {
    const response = await fetch('http://localhost:8080/api/domain/verify', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: user, token })
    });

    const data = await response.json();

    if (data.success) {
      message.value = '✅ ' + data.message;
      messageType.value = 'success';
      domainVerified.value = true;
    } else {
      message.value = '❌ ' + data.message;
      messageType.value = 'error';
      
      if (data.found_txt_records) {
        message.value += `\n\nEnregistrements TXT trouvés: ${data.found_txt_records.join(', ')}`;
      }
    }
  } catch (error) {
    message.value = `❌ Erreur réseau: ${error.message}`;
    messageType.value = 'error';
  } finally {
    verifying.value = false;
  }
}

// Remove domain
async function removeDomain() {
  if (loading.value) return;

  const user = gls().username;
  const token = gls().sessionT;

  loading.value = true;
  message.value = '';

  try {
    const response = await fetch('http://localhost:8080/api/domain/remove', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: user, token })
    });

    const data = await response.json();

    if (data.success) {
      message.value = '✅ ' + data.message;
      messageType.value = 'success';
      currentDomain.value = null;
      domainVerified.value = false;
      verificationToken.value = '';
      verificationRecord.value = '';
      showRemoveConfirm.value = false;
    } else {
      message.value = '❌ ' + data.message;
      messageType.value = 'error';
    }
  } catch (error) {
    message.value = `❌ Erreur réseau: ${error.message}`;
    messageType.value = 'error';
  } finally {
    loading.value = false;
  }
}

// Copy to clipboard
function copyToClipboard(text) {
  navigator.clipboard.writeText(text).then(() => {
    message.value = '✅ Copié dans le presse-papier !';
    messageType.value = 'success';
    setTimeout(() => { message.value = ''; }, 3000);
  }).catch(() => {
    message.value = '❌ Erreur lors de la copie';
    messageType.value = 'error';
  });
}

onMounted(() => {
  fetchDomainInfo();
});
</script>

<style scoped>
* { box-sizing: border-box; font-family: 'Roboto', sans-serif; }

/* Animations */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideIn {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

/* Container */
.custom-domain-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  animation: fadeIn 0.5s ease-in;
}

.custom-domain-container > * {
  width: 100%;
  max-width: 900px;
}

/* Hero Section */
.hero-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 48px 40px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  margin-bottom: 32px;
  gap: 20px;
  position: relative;
  overflow: hidden;
  animation: slideIn 0.6s ease-out;
}

.dark .hero-section {
  background: #1C1C1E;
}

.back-link {
  position: absolute;
  top: 20px;
  left: 24px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-family: 'Roboto', sans-serif;
  color: blue;
  text-decoration: none;
  font-weight: 600;
  font-size: 15px;
  transition: all 0.3s ease;
}

.dark .back-link {
  color: #6ec6ff;
}

.back-link:hover {
  gap: 12px;
  color: rgba(0, 0, 255, 0.8);
}

.back-link svg {
  transition: transform 0.3s ease;
}

.back-link:hover svg {
  transform: translateX(-4px);
}

.hero-icon-wrapper {
  width: 80px;
  height: 80px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse 2s ease-in-out infinite;
  position: relative;
  overflow: hidden;
}

.icon-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 60%;
  height: 60%;
  transform: translate(-50%, -50%);
  background: rgba(255, 255, 255, 0.4);
  border-radius: 50%;
  filter: blur(12px);
  z-index: 0;
}

.hero-icon-wrapper:hover {
  transform: scale(1.05);
  transition: transform 0.3s ease;
}

.hero-icon {
  width: 40px;
  height: 40px;
  color: white;
  position: relative;
  z-index: 1;
}

.hero-title {
  font-family: 'Roboto', sans-serif;
  font-size: 2rem;
  font-weight: 700;
  text-align: center;
  letter-spacing: 0.5px;
  color: #222;
  margin: 8px 0 0 0;
  line-height: 1.2;
}

.dark .hero-title {
  color: white;
}

.hero-subtitle {
  font-family: 'Roboto', sans-serif;
  font-size: 1.1rem;
  font-style: italic;
  text-align: center;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 0.3px;
  margin: 0;
}

.dark .hero-subtitle {
  background: -webkit-linear-gradient(30deg, #6ec6ff, #ff6b9d);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.plan-badge-hero {
  display: inline-block;
  padding: 8px 20px;
  background: linear-gradient(135deg, rgba(255, 215, 0, 0.2) 0%, rgba(255, 165, 0, 0.2) 100%);
  color: #cc7a00;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.dark .plan-badge-hero {
  background: linear-gradient(135deg, rgba(255, 215, 0, 0.25) 0%, rgba(255, 165, 0, 0.25) 100%);
  color: #FFD700;
}

/* Cards */
.domain-status {
  margin-bottom: 32px;
  animation: slideIn 0.7s ease-out;
}

.status-card {
  background: rgba(245, 245, 247, 0.85);
  padding: 40px;
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  position: relative;
  overflow: hidden;
}

.dark .status-card {
  background: #1C1C1E;
}

.status-card.verified {
  border-left: 6px solid #28a745;
}

.status-card.pending {
  border-left: 6px solid #ffc107;
}

.domain-info {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
}

.domain-info h3 {
  margin: 0;
  font-family: 'Roboto', sans-serif;
  color: #222;
  font-size: 26px;
  font-weight: 700;
}

.dark .domain-info h3 {
  color: white;
}

.badge-verified, .badge-pending {
  padding: 8px 16px;
  border-radius: 12px;
  font-size: 13px;
  font-family: 'Roboto', sans-serif;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.badge-verified {
  background: linear-gradient(135deg, #d4edda 0%, #c3e6cb 100%);
  color: #155724;
}

.badge-pending {
  background: linear-gradient(135deg, #fff3cd 0%, #ffeaa7 100%);
  color: #856404;
}

/* Verification Instructions */
.verification-instructions {
  background: rgba(102, 126, 234, 0.08);
  padding: 32px;
  border-radius: 20px;
  margin: 24px 0;
  border: 2px solid rgba(102, 126, 234, 0.2);
}

.dark .verification-instructions {
  background: rgba(102, 126, 234, 0.15);
  border-color: rgba(102, 126, 234, 0.3);
}

.verification-instructions h4 {
  margin-top: 0;
  font-family: 'Roboto', sans-serif;
  color: #222;
  font-size: 20px;
  font-weight: 700;
  margin-bottom: 16px;
}

.dark .verification-instructions h4 {
  color: white;
}

.dns-record {
  background: white;
  border: 2px solid rgba(102, 126, 234, 0.3);
  border-radius: 16px;
  padding: 24px;
  margin: 20px 0;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
}

.dark .dns-record {
  background: #2C2C2E;
  border-color: rgba(102, 126, 234, 0.4);
}

.record-field {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 16px 0;
}

.record-field label {
  font-family: 'Roboto', sans-serif;
  font-weight: 700;
  min-width: 120px;
  color: #555;
  font-size: 14px;
}

.dark .record-field label {
  color: #e5e7eb;
}

.record-field code {
  background: rgba(102, 126, 234, 0.1);
  padding: 8px 14px;
  border-radius: 8px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  color: #667eea;
  font-weight: 600;
}

.dark .record-field code {
  background: rgba(102, 126, 234, 0.2);
  color: #a78bfa;
}

.token-value {
  flex: 1;
  word-break: break-all;
}

/* Buttons */
.btn-copy {
  background: -webkit-linear-gradient(30deg, blue, red);
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 10px;
  cursor: pointer;
  font-family: 'Roboto', sans-serif;
  font-size: 13px;
  font-weight: 600;
  white-space: nowrap;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-copy:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(102, 126, 234, 0.4);
}

.btn-verify {
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
  border: none;
  padding: 14px 32px;
  border-radius: 12px;
  font-family: 'Roboto', sans-serif;
  font-weight: 700;
  font-size: 15px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 6px 20px rgba(40, 167, 69, 0.3);
  margin-top: 20px;
}

.btn-verify:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(40, 167, 69, 0.4);
}

.btn-verify:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.help-text {
  margin: 20px 0;
  font-family: 'Roboto', sans-serif;
  color: #666;
  font-size: 14px;
  line-height: 1.6;
}

.dark .help-text {
  color: #9ca3af;
}

.help-text p {
  margin: 12px 0;
}

.btn-verify:hover:not(:disabled) {
  background: #218838;
}

.btn-verify:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.domain-actions {
  margin-top: 20px;
  display: flex;
  gap: 10px;
}

.btn-remove, .btn-cancel {
  background: #dc3545;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
}

.btn-cancel {
  background: #6c757d;
}

.btn-remove:hover, .btn-cancel:hover {
  opacity: 0.9;
}

.add-domain-form {
  background: rgba(245, 245, 247, 0.85);
  padding: 40px;
  border-radius: 32px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.10);
  margin-bottom: 32px;
  animation: slideIn 0.7s ease-out;
}

.dark .add-domain-form {
  background: #1C1C1E;
  box-shadow: 0 8px 32px rgba(0,0,0,0.3);
}

.add-domain-form h3 {
  margin-top: 0;
  font-family: 'Roboto', sans-serif;
  font-size: 24px;
  font-weight: 700;
  color: #222;
  margin-bottom: 24px;
}

.dark .add-domain-form h3 {
  color: white;
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  margin-bottom: 10px;
  color: #333;
  font-family: 'Roboto', sans-serif;
  font-weight: 700;
  font-size: 14px;
}

.dark .form-group label {
  color: #e5e7eb;
}

.form-group input {
  width: 100%;
  padding: 14px 16px;
  border: 2px solid rgba(102, 126, 234, 0.2);
  border-radius: 12px;
  font-family: 'Roboto', sans-serif;
  font-size: 15px;
  transition: all 0.3s ease;
  background: white;
}

.dark .form-group input {
  background: #2C2C2E;
  border-color: rgba(102, 126, 234, 0.3);
  color: white;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.dark .form-group input:focus {
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.2);
}

.form-group small {
  display: block;
  margin-top: 8px;
  color: #888;
  font-family: 'Roboto', sans-serif;
  font-size: 12px;
}

.dark .form-group small {
  color: #9ca3af;
}

.example-emails {
  background: rgba(102, 126, 234, 0.08);
  padding: 20px;
  border-radius: 16px;
  margin: 20px 0;
  border: 2px solid rgba(102, 126, 234, 0.15);
}

.dark .example-emails {
  background: rgba(102, 126, 234, 0.15);
  border-color: rgba(102, 126, 234, 0.25);
}

.example-emails p {
  font-family: 'Roboto', sans-serif;
  font-weight: 700;
  margin: 0 0 12px 0;
  color: #333;
}

.dark .example-emails p {
  color: #e5e7eb;
}

.example-emails ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.example-emails li {
  padding: 8px 0;
  color: #667eea;
  font-family: 'Roboto', sans-serif;
  font-weight: 600;
  font-size: 14px;
}

.dark .example-emails li {
  color: #a78bfa;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 14px 28px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: transform 0.2s;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.7);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease;
}

.modal-content {
  background: rgba(245, 245, 247, 0.98);
  padding: 40px;
  border-radius: 24px;
  max-width: 500px;
  width: 90%;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slideIn 0.3s ease-out;
}

.dark .modal-content {
  background: #1C1C1E;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.6);
}

.modal-content h3 {
  margin-top: 0;
  font-family: 'Roboto', sans-serif;
  font-size: 22px;
  font-weight: 700;
  color: #dc3545;
  margin-bottom: 16px;
}

.dark .modal-content h3 {
  color: #ff6b6b;
}

.modal-content p {
  font-family: 'Roboto', sans-serif;
  color: #555;
  line-height: 1.6;
  margin: 12px 0;
}

.dark .modal-content p {
  color: #d1d5db;
}

.warning-text {
  color: #721c24;
  background: rgba(220, 53, 69, 0.15);
  padding: 14px;
  border-radius: 10px;
  margin: 20px 0;
  border-left: 4px solid #dc3545;
  font-family: 'Roboto', sans-serif;
  font-weight: 600;
}

.dark .warning-text {
  color: #ffaaaa;
  background: rgba(220, 53, 69, 0.2);
  border-left-color: #ff6b6b;
}

.modal-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
}

.btn-confirm-remove {
  background: #dc3545;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 600;
}

.btn-confirm-remove:hover:not(:disabled) {
  background: #c82333;
}

.btn-confirm-remove:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.message {
  padding: 15px 20px;
  border-radius: 8px;
  margin: 20px 0;
  font-weight: 500;
  white-space: pre-line;
}

.message.success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.dark .message.success {
  background: rgba(40, 167, 69, 0.2);
  color: #75b798;
  border-color: rgba(40, 167, 69, 0.4);
}

.message.error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.dark .message.error {
  background: rgba(220, 53, 69, 0.2);
  color: #ffaaaa;
  border-color: rgba(220, 53, 69, 0.4);
}

.help-section {
  background: rgba(245, 245, 247, 0.85);
  padding: 32px;
  border-radius: 24px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.08);
  margin-top: 32px;
  animation: slideIn 0.8s ease-out;
}

.dark .help-section {
  background: #1C1C1E;
  box-shadow: 0 4px 20px rgba(0,0,0,0.3);
}

.help-section h3 {
  margin-top: 0;
  font-family: 'Roboto', sans-serif;
  font-size: 22px;
  font-weight: 700;
  color: #222;
  margin-bottom: 16px;
}

.dark .help-section h3 {
  color: white;
}

.help-section p {
  font-family: 'Roboto', sans-serif;
  color: #555;
  line-height: 1.6;
  margin-bottom: 20px;
}

.dark .help-section p {
  color: #d1d5db;
}

.dns-table {
  width: 100%;
  border-collapse: collapse;
  margin: 20px 0;
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.dark .dns-table {
  background: #2C2C2E;
}

.dns-table th, .dns-table td {
  padding: 14px 16px;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
  font-family: 'Roboto', sans-serif;
}

.dark .dns-table th, .dark .dns-table td {
  border-bottom-color: #3a3d4a;
}

.dns-table th {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
  font-weight: 700;
  color: #333;
  font-size: 14px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.dark .dns-table th {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%);
  color: #e5e7eb;
}

.dns-table td {
  color: #555;
}

.dark .dns-table td {
  color: #d1d5db;
}

.dns-table code {
  background: rgba(102, 126, 234, 0.1);
  padding: 6px 10px;
  border-radius: 6px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #667eea;
  font-weight: 600;
}

.dark .dns-table code {
  background: rgba(102, 126, 234, 0.2);
  color: #a78bfa;
}

.help-note {
  font-family: 'Roboto', sans-serif;
  color: #666;
  font-size: 14px;
  margin-top: 20px;
  line-height: 1.6;
  padding: 16px;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 10px;
  border-left: 4px solid #667eea;
}

.dark .help-note {
  color: #9ca3af;
  background: rgba(102, 126, 234, 0.1);
  border-left-color: #a78bfa;
}

@media (max-width: 768px) {
  .custom-domain-panel {
    padding: 10px;
  }

  .record-field {
    flex-direction: column;
    align-items: flex-start;
  }

  .dns-table {
    font-size: 13px;
  }

  .dns-table th, .dns-table td {
    padding: 8px;
  }
}
</style>
