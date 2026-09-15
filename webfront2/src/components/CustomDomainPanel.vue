<script setup>
import { ref, onMounted } from 'vue'
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
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/domain/info`, {
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
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/domain/add`, {
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
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/domain/verify`, {
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
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/domain/remove`, {
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

<template>
  <div class="domain-panel">
    <div class="head">
      <h2>Domaine personnalisé</h2>
      <p>Configurez votre propre domaine pour vos emails.</p>
    </div>

    <!-- Current domain -->
    <div v-if="currentDomain" class="card current">
      <div class="current-head">
        <div>
          <span class="domain-name">{{ currentDomain }}</span>
          <span class="status" :class="{ verified: domainVerified, unverified: !domainVerified }">
            {{ domainVerified ? 'Vérifié' : 'Non vérifié' }}
          </span>
        </div>
        <button @click="showRemoveConfirm = true" class="btn-remove">Supprimer</button>
      </div>

      <!-- Verification info -->
      <div v-if="!domainVerified" class="verify-section">
        <h3>Vérification</h3>
        <p>Ajoutez cet enregistrement DNS TXT à votre domaine :</p>
        <div class="dns-record" @click="copyToClipboard(verificationRecord)">
          <code>{{ verificationRecord }}</code>
          <span class="copy-hint">Cliquer pour copier</span>
        </div>
        <button @click="verifyDomain" class="btn-verify" :disabled="verifying">
          <span v-if="verifying" class="spinner"></span>
          {{ verifying ? 'Vérification...' : 'Vérifier le domaine' }}
        </button>
      </div>
    </div>

    <!-- Add domain -->
    <div v-else class="card add">
      <h3>Ajouter un domaine</h3>
      <div class="add-row">
        <input v-model="newDomain" type="text" placeholder="exemple.com" @keyup.enter="addDomain" />
        <button @click="addDomain" class="btn-add" :disabled="loading">
          <span v-if="loading" class="spinner"></span>
          {{ loading ? '...' : 'Ajouter' }}
        </button>
      </div>
    </div>

    <!-- Message -->
    <p v-if="message" class="msg" :class="messageType">{{ message }}</p>

    <!-- Remove confirm modal -->
    <div v-if="showRemoveConfirm" class="overlay" @click.self="showRemoveConfirm = false">
      <div class="dialog">
        <h3>Supprimer le domaine ?</h3>
        <p>Cette action est irréversible.</p>
        <div class="dialog-actions">
          <button @click="showRemoveConfirm = false" class="btn-cancel">Annuler</button>
          <button @click="removeDomain" class="btn-confirm" :disabled="loading">
            <span v-if="loading" class="spinner"></span>
            {{ loading ? '...' : 'Supprimer' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.domain-panel { max-width: 560px; margin: 0 auto; padding: 2rem 1.25rem 4rem; animation: fadeInUp .5s var(--ease) both; }
.head { margin-bottom: 1.5rem; animation: fadeInDown .4s var(--ease) both; }
.head h2 { font-size: 1.375rem; font-weight: 800; letter-spacing: -.02em; margin-bottom: .375rem; }
.head p { font-size: .875rem; color: var(--text-secondary); }
.card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1.75rem; margin-bottom: 1rem; animation: fadeInUp .4s var(--ease) both; }
.current-head { display: flex; align-items: center; justify-content: space-between; gap: 1rem; flex-wrap: wrap; }
.domain-name { font-size: 1.125rem; font-weight: 700; }
.status { display: inline-block; padding: .2rem .65rem; font-size: .75rem; font-weight: 600; border-radius: var(--r-full); margin-left: .5rem; }
.status.verified { background: var(--bg-tertiary); color: var(--text-primary); }
.status.unverified { background: var(--bg-tertiary); color: var(--text-muted); border: 1px solid var(--border); }
.btn-remove { padding: .5rem 1.125rem; font-size: .8125rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-remove:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.verify-section { margin-top: 1.5rem; padding-top: 1.5rem; border-top: 1px solid var(--border); animation: fadeIn .3s var(--ease) both; }
.verify-section h3 { font-size: 1rem; font-weight: 700; margin-bottom: .5rem; }
.verify-section p { font-size: .8125rem; color: var(--text-secondary); margin-bottom: .75rem; }
.dns-record { display: flex; align-items: center; gap: .75rem; padding: .75rem 1rem; background: var(--bg-tertiary); border: 1px solid var(--border); border-radius: var(--r); cursor: pointer; transition: border-color var(--tb); animation: fadeIn .3s var(--ease) both; }
.dns-record:hover { border-color: var(--border-strong); }
.dns-record code { font-family: var(--font-mono); font-size: .8125rem; color: var(--text-primary); flex: 1; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.copy-hint { font-size: .75rem; color: var(--text-muted); white-space: nowrap; }
.btn-verify { margin-top: 1rem; padding: .625rem 1.25rem; font-size: .875rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); display: flex; align-items: center; gap: .375rem; }
.btn-verify:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-sm); }
.btn-verify:disabled { opacity: .5; cursor: not-allowed; }
.add h3 { font-size: 1.0625rem; font-weight: 700; margin-bottom: 1rem; }
.add-row { display: flex; gap: .625rem; }
.add-row input { flex: 1; padding: .625rem .875rem; font-size: .9375rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); transition: border-color var(--tb); }
.add-row input:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.add-row input::placeholder { color: var(--text-muted); }
.btn-add { padding: .625rem 1.25rem; font-size: .875rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); display: flex; align-items: center; gap: .375rem; white-space: nowrap; }
.btn-add:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); }
.btn-add:disabled { opacity: .5; cursor: not-allowed; }
.msg { padding: .75rem 1rem; font-size: .8125rem; border-radius: var(--r); border: 1px solid var(--border); white-space: pre-wrap; animation: fadeIn .2s var(--ease) both; }
.msg.success { background: var(--bg-tertiary); color: var(--text-primary); }
.msg.error { background: var(--bg-tertiary); color: var(--text-primary); }
.overlay { position: fixed; inset: 0; z-index: 999; background: var(--bg-overlay); backdrop-filter: blur(6px); display: flex; align-items: center; justify-content: center; animation: fadeIn .2s var(--ease) both; }
.dialog { width: calc(100% - 32px); max-width: 360px; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); box-shadow: var(--sh-xl); padding: 1.75rem; text-align: center; animation: popIn .25s var(--ease-spring) both; }
.dialog h3 { font-size: 1.125rem; font-weight: 700; margin-bottom: .375rem; }
.dialog p { font-size: .875rem; color: var(--text-secondary); margin-bottom: 1.25rem; }
.dialog-actions { display: flex; gap: .625rem; }
.btn-cancel { flex: 1; padding: .625rem; font-size: .875rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-cancel:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.btn-confirm { flex: 1; padding: .625rem; font-size: .875rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); display: flex; align-items: center; justify-content: center; gap: .375rem; }
.btn-confirm:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); }
.btn-confirm:disabled { opacity: .5; }
.spinner { width: 14px; height: 14px; border: 2px solid var(--border); border-top-color: var(--ink); border-radius: 50%; animation: spin .6s linear infinite; display: inline-block; }
.btn-verify .spinner, .btn-confirm .spinner { border-top-color: var(--paper); }
</style>
