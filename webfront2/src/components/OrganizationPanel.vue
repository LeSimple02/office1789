<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { gls } from '@/stores/global.js';

const router = useRouter();
const loading = ref(false);
const message = ref('');
const messageType = ref('success');

const userInfo = ref({});
const accountType = ref('personal');
const organizationInfo = ref(null);
const members = ref([]);

const newSubAccount = ref({
  username: '',
  password: '',
  recoveryEmail: '',
  phone: '',
  organizationName: ''
});

const showPassword = ref(false);
const passwordStrength = ref({
  class: '',
  text: '',
  width: '0%'
});

function togglePasswordVisibility() {
  showPassword.value = !showPassword.value;
}

function checkPasswordStrength() {
  const password = newSubAccount.value.password;
  let strength = 0;
  let text = '';
  let color = '#737373';
  
  if (password.length >= 8) strength++;
  if (password.length >= 12) strength++;
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++;
  if (/\d/.test(password)) strength++;
  if (/[^a-zA-Z0-9]/.test(password)) strength++;
  
  if (strength <= 2) {
    text = 'Faible';
    color = '#696969';
  } else if (strength === 3) {
    text = 'Moyen';
    color = '#bebebe';
  } else if (strength === 4) {
    text = 'Fort';
    color = '#767676';
  } else {
    text = 'Très fort';
    color = '#3d3d3d';
  }
  
  passwordStrength.value = {
    text,
    color,
    width: (strength * 20) + '%'
  };
}

// Fetch organization info and members
async function fetchOrganizationData() {
  console.log('🔍 OrganizationPanel: fetchOrganizationData called')
  const username = gls().username;
  const token = gls().sessionT;
  
  console.log('📝 Session check:', { username, token: token ? 'exists' : 'missing' })

  if (!username || !token) {
    console.log('❌ No session - redirecting to login')
    router.push('/login');
    return;
  }

  try {
    // Get user info first
    const userResponse = await fetch(import.meta.env.VITE_APP_API_GETINFOP, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, token })
    });

    if (!userResponse.ok) throw new Error('Failed to get user info');
    
    const userData = await userResponse.json();
    userInfo.value = userData;

    // Only fetch organization data if user has Pro or Enterprise
    if (userData.Nboffer >= 2) {
      const orgResponse = await fetch(`${import.meta.env.VITE_APP_API}/api/organization/members`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, token })
      });

      if (orgResponse.ok) {
        const orgData = await orgResponse.json();
        if (orgData.success) {
          organizationInfo.value = orgData;
          members.value = orgData.members || [];
          
          // Find current user's account type
          const currentUser = members.value.find(m => m.username === username);
          if (currentUser) {
            accountType.value = currentUser.account_type;
          }
        }
      }
    }
  } catch (error) {
    console.error('Error fetching organization data:', error);
  }
}

// Create sub-account
async function createSubAccount() {
  if (loading.value) return;

  const username = gls().username;
  const token = gls().sessionT;

  loading.value = true;
  message.value = '';

  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/organization/create-subaccount`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username,
        token,
        sub_username: newSubAccount.value.username,
        sub_password: newSubAccount.value.password,
        sub_email: newSubAccount.value.recoveryEmail,
        sub_phone_number: newSubAccount.value.phone,
        organization_name: newSubAccount.value.organizationName
      })
    });

    const data = await response.json();

    if (data.success) {
      message.value = `✅ Sous-compte créé avec succès ! Email: ${data.sub_account_email}`;
      messageType.value = 'success';
      
      // Reset form
      newSubAccount.value = {
        username: '',
        password: '',
        recoveryEmail: '',
        phone: '',
        organizationName: ''
      };

      // Refresh data
      await fetchOrganizationData();
    } else {
      message.value = `❌ ${data.message}`;
      messageType.value = 'error';
    }
  } catch (error) {
    message.value = `❌ Erreur réseau: ${error.message}`;
    messageType.value = 'error';
  } finally {
    loading.value = false;
  }
}

// Delete sub-account
async function deleteSubAccount(subAccountId) {
  if (!confirm('Êtes-vous sûr de vouloir supprimer ce sous-compte ? Cette action est irréversible.')) {
    return;
  }

  const username = gls().username;
  const token = gls().sessionT;

  loading.value = true;
  message.value = '';

  try {
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/organization/delete-member`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username,
        token,
        sub_account_id: subAccountId
      })
    });

    const data = await response.json();

    if (data.success) {
      message.value = '✅ Sous-compte supprimé avec succès';
      messageType.value = 'success';
      await fetchOrganizationData();
    } else {
      message.value = `❌ ${data.message}`;
      messageType.value = 'error';
    }
  } catch (error) {
    message.value = `❌ Erreur réseau: ${error.message}`;
    messageType.value = 'error';
  } finally {
    loading.value = false;
  }
}

// Format date
function formatDate(dateString) {
  const date = new Date(dateString);
  return date.toLocaleDateString('fr-FR', { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  });
}

// Get plan name
function getPlanName(nboffer) {
  const plans = {
    0: 'Free',
    1: 'Standard',
    2: 'Professional',
    3: 'Enterprise'
  };
  return plans[nboffer] || 'Free';
}

// Get plan CSS class
function getPlanClass(nboffer) {
  const classes = {
    0: 'free',
    1: 'standard',
    2: 'professional',
    3: 'enterprise'
  };
  return classes[nboffer] || 'free';
}

// Get storage size
function getStorageSize(nboffer) {
  const storage = {
    0: '1GB',
    1: '50GB',
    2: '200GB',
    3: 'Illimité'
  };
  return storage[nboffer] || '1GB';
}

onMounted(() => {
  fetchOrganizationData();
});
</script>

<template>
  <div class="org">
    <div class="head">
      <h2>Organisation</h2>
      <p>Gérez votre équipe et sous-comptes.</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="loading">
      <div class="spinner"></div>
      <span>Chargement...</span>
    </div>

    <template v-else>
      <!-- No org (Free plan) -->
      <div v-if="!organizationInfo" class="card empty">
        <h3>Plan Free</h3>
        <p>Passez à un plan Professional ou Enterprise pour gérer une organisation.</p>
      </div>

      <!-- Org info -->
      <template v-else>
        <!-- Create sub-account -->
        <div class="card create">
          <h3>Créer un sous-compte</h3>
          <div class="form">
            <div class="row">
              <input v-model="newSubAccount.username" placeholder="Nom d'utilisateur" />
            </div>
            <div class="row">
              <div class="pwd-wrap">
                <input v-model="newSubAccount.password" :type="showPassword ? 'text' : 'password'" placeholder="Mot de passe" @input="checkPasswordStrength" />
                <button type="button" @click="togglePasswordVisibility" class="pwd-toggle">{{ showPassword ? '🙈' : '👁' }}</button>
              </div>
              <div v-if="newSubAccount.password" class="strength">
                <div class="bar-bg"><div class="bar-fill" :style="{ width: passwordStrength.width, background: passwordStrength.color }"></div></div>
                <span class="strength-text" :style="{ color: passwordStrength.color }">{{ passwordStrength.text }}</span>
              </div>
            </div>
            <div class="row">
              <input v-model="newSubAccount.recoveryEmail" type="email" placeholder="Email de récupération" />
            </div>
            <div class="row">
              <input v-model="newSubAccount.phone" type="tel" placeholder="Téléphone (optionnel)" />
            </div>
            <div class="row">
              <input v-model="newSubAccount.organizationName" placeholder="Nom d'organisation (optionnel)" />
            </div>
            <button @click="createSubAccount" class="btn-create" :disabled="loading">
              <span v-if="loading" class="spinner"></span>
              {{ loading ? '...' : 'Créer le sous-compte' }}
            </button>
          </div>
        </div>

        <!-- Members list -->
        <div class="card members">
          <h3>Membres ({{ members.length }})</h3>
          <div class="member-list">
            <div v-for="m in members" :key="m.id" class="member">
              <div class="member-avatar">{{ (m.username || '?').charAt(0).toUpperCase() }}</div>
              <div class="member-info">
                <span class="member-name">{{ m.username }}</span>
                <span class="member-type">{{ m.account_type }}</span>
              </div>
              <span class="member-plan">{{ getPlanName(m.nboffer) }}</span>
              <button v-if="m.username !== gls().username" @click="deleteSubAccount(m.id)" class="btn-del">Supprimer</button>
            </div>
            <div v-if="members.length === 0" class="empty-msg">Aucun membre</div>
          </div>
        </div>
      </template>

      <p v-if="message" class="msg" :class="messageType">{{ message }}</p>
    </template>
  </div>
</template>

<style scoped>
.org { max-width: 560px; margin: 0 auto; padding: 2rem 1.25rem 4rem; animation: fadeInUp .5s var(--ease) both; }
.head { margin-bottom: 1.5rem; animation: fadeInDown .4s var(--ease) both; }
.head h2 { font-size: 1.375rem; font-weight: 800; letter-spacing: -.02em; margin-bottom: .375rem; }
.head p { font-size: .875rem; color: var(--text-secondary); }
.loading { display: flex; align-items: center; justify-content: center; gap: .625rem; padding: 3rem 2rem; color: var(--text-muted); font-size: .875rem; }
.spinner { width: 20px; height: 20px; border: 2.5px solid var(--border); border-top-color: var(--ink); border-radius: 50%; animation: spin .7s linear infinite; }
.card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1.75rem; margin-bottom: 1rem; animation: fadeInUp .4s var(--ease) both; }
.card.empty { text-align: center; }
.card.empty h3 { font-size: 1.125rem; font-weight: 700; margin-bottom: .5rem; }
.card.empty p { font-size: .875rem; color: var(--text-secondary); }
.create h3, .members h3 { font-size: 1.0625rem; font-weight: 700; margin-bottom: 1.125rem; }
.form { display: flex; flex-direction: column; gap: .875rem; }
.row { display: flex; flex-direction: column; gap: .375rem; }
.row input { width: 100%; padding: .5625rem .75rem; font-size: .875rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); transition: border-color var(--tb), box-shadow var(--tb); }
.row input:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.row input::placeholder { color: var(--text-muted); }
.pwd-wrap { position: relative; }
.pwd-wrap input { padding-right: 38px; }
.pwd-toggle { position: absolute; right: 6px; top: 50%; transform: translateY(-50%); background: none; border: none; font-size: 16px; cursor: pointer; opacity: .6; padding: 4px; border-radius: var(--r-sm); transition: opacity var(--tb); }
.pwd-toggle:hover { opacity: 1; }
.strength { display: flex; align-items: center; gap: .5rem; }
.bar-bg { flex: 1; height: 4px; background: var(--bg-tertiary); border-radius: var(--r-full); overflow: hidden; }
.bar-fill { height: 100%; border-radius: var(--r-full); transition: width var(--tb); animation: growBar .3s var(--ease) both; }
.strength-text { font-size: .75rem; font-weight: 500; white-space: nowrap; }
.btn-create { padding: .625rem 1.25rem; font-size: .875rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); display: flex; align-items: center; justify-content: center; gap: .375rem; }
.btn-create:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-sm); }
.btn-create:disabled { opacity: .5; cursor: not-allowed; }
.member-list { display: flex; flex-direction: column; gap: .5rem; }
.member { display: flex; align-items: center; gap: .75rem; padding: .75rem .875rem; background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--r); transition: all var(--tb); animation: fadeInLeft .3s var(--ease) both; }
.member:hover { border-color: var(--border-strong); }
.member-avatar { width: 32px; height: 32px; border-radius: 50%; background: var(--bg-tertiary); border: 1px solid var(--border); display: flex; align-items: center; justify-content: center; font-size: .8125rem; font-weight: 700; flex-shrink: 0; }
.member-info { flex: 1; display: flex; flex-direction: column; }
.member-name { font-size: .875rem; font-weight: 600; }
.member-type { font-size: .75rem; color: var(--text-muted); }
.member-plan { font-size: .75rem; font-weight: 600; padding: .125rem .5rem; border-radius: var(--r-full); background: var(--bg-tertiary); color: var(--text-secondary); }
.btn-del { padding: .375rem .75rem; font-size: .75rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r-sm); cursor: pointer; transition: all var(--tb); }
.btn-del:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.empty-msg { text-align: center; padding: 1.5rem 1rem; font-size: .8125rem; color: var(--text-muted); }
.msg { padding: .75rem 1rem; font-size: .8125rem; border-radius: var(--r); border: 1px solid var(--border); animation: fadeIn .2s var(--ease) both; }
.msg.success { background: var(--bg-tertiary); color: var(--text-primary); }
.msg.error { background: var(--bg-tertiary); color: var(--text-primary); }
</style>
