<template>
  <div class="custom-domain-panel">
    <h2>🌐 Domaine personnalisé</h2>
    
    <div class="info-banner">
      <p>📧 Utilisez votre propre nom de domaine pour vos adresses email (ex: contact@votreentreprise.com)</p>
      <p class="plan-requirement">Disponible pour les plans <strong>Professional</strong> et <strong>Enterprise</strong></p>
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

const router = useRouter();
const loading = ref(false);
const verifying = ref(false);
const message = ref('');
const messageType = ref('success');
const showRemoveConfirm = ref(false);

const username = ref(sessionStorage.getItem('username') || '');
const currentDomain = ref(null);
const domainVerified = ref(false);
const verificationToken = ref('');
const verificationRecord = ref('');
const newDomain = ref('');

// Fetch current domain info
async function fetchDomainInfo() {
  const user = sessionStorage.getItem('username');
  const token = sessionStorage.getItem('token');

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

  const user = sessionStorage.getItem('username');
  const token = sessionStorage.getItem('token');

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

  const user = sessionStorage.getItem('username');
  const token = sessionStorage.getItem('token');

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

  const user = sessionStorage.getItem('username');
  const token = sessionStorage.getItem('token');

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
.custom-domain-panel {
  max-width: 900px;
  margin: 20px auto;
  padding: 20px;
}

h2 {
  color: #2c3e50;
  margin-bottom: 20px;
}

.info-banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px;
  border-radius: 12px;
  margin-bottom: 30px;
}

.info-banner p {
  margin: 8px 0;
}

.plan-requirement {
  font-size: 14px;
  opacity: 0.9;
}

.domain-status {
  margin-bottom: 30px;
}

.status-card {
  background: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.status-card.verified {
  border-left: 5px solid #28a745;
}

.status-card.pending {
  border-left: 5px solid #ffc107;
}

.domain-info {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
}

.domain-info h3 {
  margin: 0;
  color: #2c3e50;
  font-size: 24px;
}

.badge-verified, .badge-pending {
  padding: 6px 14px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: bold;
}

.badge-verified {
  background: #d4edda;
  color: #155724;
}

.badge-pending {
  background: #fff3cd;
  color: #856404;
}

.verification-instructions {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin: 20px 0;
}

.verification-instructions h4 {
  margin-top: 0;
  color: #2c3e50;
}

.dns-record {
  background: white;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  padding: 15px;
  margin: 15px 0;
}

.record-field {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 10px 0;
}

.record-field label {
  font-weight: 600;
  min-width: 100px;
  color: #555;
}

.record-field code {
  background: #f0f0f0;
  padding: 6px 12px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
}

.token-value {
  flex: 1;
  word-break: break-all;
}

.btn-copy {
  background: #667eea;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  white-space: nowrap;
}

.btn-copy:hover {
  background: #5568d3;
}

.help-text {
  margin: 15px 0;
  color: #666;
  font-size: 14px;
}

.help-text p {
  margin: 8px 0;
}

.btn-verify {
  background: #28a745;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 8px;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s;
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
  background: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  margin-bottom: 30px;
}

.add-domain-form h3 {
  margin-top: 0;
  color: #2c3e50;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #555;
  font-weight: 600;
}

.form-group input {
  width: 100%;
  padding: 12px;
  border: 2px solid #e0e0e0;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.form-group small {
  display: block;
  margin-top: 5px;
  color: #888;
  font-size: 13px;
}

.example-emails {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 8px;
  margin: 15px 0;
}

.example-emails ul {
  list-style: none;
  padding: 0;
  margin: 10px 0 0 0;
}

.example-emails li {
  padding: 5px 0;
  color: #667eea;
  font-weight: 500;
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
  background: rgba(0,0,0,0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 30px;
  border-radius: 12px;
  max-width: 500px;
  width: 90%;
}

.modal-content h3 {
  margin-top: 0;
  color: #dc3545;
}

.warning-text {
  color: #721c24;
  background: #f8d7da;
  padding: 10px;
  border-radius: 6px;
  margin: 15px 0;
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

.message.error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.help-section {
  background: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.help-section h3 {
  margin-top: 0;
  color: #2c3e50;
}

.dns-table {
  width: 100%;
  border-collapse: collapse;
  margin: 15px 0;
}

.dns-table th, .dns-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #e0e0e0;
}

.dns-table th {
  background: #f8f9fa;
  font-weight: 600;
  color: #555;
}

.dns-table code {
  background: #f0f0f0;
  padding: 4px 8px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.help-note {
  color: #666;
  font-size: 14px;
  margin-top: 15px;
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
