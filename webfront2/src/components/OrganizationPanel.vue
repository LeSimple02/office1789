<template>
  <div class="organization-panel">
    <h2>🏢 Gestion de l'organisation</h2>

    <!-- Status information -->
    <div v-if="organizationInfo" class="org-status">
      <div class="status-card">
        <h3>{{ organizationInfo.organization_name }}</h3>
        <p class="account-type">
          <span v-if="accountType === 'organization_owner'" class="badge-owner">👑 Propriétaire</span>
          <span v-else class="badge-member">👤 Membre</span>
        </p>
        <p class="member-count">
          <strong>{{ organizationInfo.current_count }}/{{ organizationInfo.max_members }}</strong> comptes utilisés
        </p>
      </div>
    </div>

    <!-- Create sub-account form (only for owner) -->
    <div v-if="accountType === 'organization_owner' && organizationInfo && organizationInfo.current_count < organizationInfo.max_members" class="create-subaccount">
      <h3>Créer un nouveau sous-compte</h3>
      <form @submit.prevent="createSubAccount">
        <div class="form-group">
          <label for="sub_username">Nom d'utilisateur *</label>
          <input 
            type="text" 
            id="sub_username" 
            v-model="newSubAccount.username" 
            required 
            placeholder="john.doe"
            pattern="[a-z0-9._-]+"
            title="Lettres minuscules, chiffres, .-_ uniquement"
          />
          <small>@office1789.com sera automatiquement ajouté</small>
        </div>

        <div class="form-group">
          <label for="sub_password">Mot de passe *</label>
          <input 
            type="password" 
            id="sub_password" 
            v-model="newSubAccount.password" 
            required 
            minlength="8"
            placeholder="Minimum 8 caractères"
          />
        </div>

        <div class="form-group">
          <label for="sub_email">Email de récupération</label>
          <input 
            type="email" 
            id="sub_email" 
            v-model="newSubAccount.recoveryEmail" 
            placeholder="john@example.com (optionnel)"
          />
        </div>

        <div class="form-group">
          <label for="sub_phone">Téléphone</label>
          <input 
            type="tel" 
            id="sub_phone" 
            v-model="newSubAccount.phone" 
            placeholder="+33 6 XX XX XX XX (optionnel)"
          />
        </div>

        <div class="form-group" v-if="!organizationInfo.organization_id">
          <label for="org_name">Nom de l'organisation</label>
          <input 
            type="text" 
            id="org_name" 
            v-model="newSubAccount.organizationName" 
            placeholder="Ex: Mon Entreprise SAS"
          />
          <small>Seulement pour le premier sous-compte</small>
        </div>

        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? '⏳ Création...' : '✅ Créer le sous-compte' }}
        </button>
      </form>
    </div>

    <!-- Limit reached message -->
    <div v-if="accountType === 'organization_owner' && organizationInfo && organizationInfo.current_count >= organizationInfo.max_members" class="limit-reached">
      <p>⚠️ Vous avez atteint la limite de {{ organizationInfo.max_members }} comptes.</p>
      <p v-if="userInfo.Nboffer < 3">
        <router-link to="/account">Passez à l'offre Enterprise (20 comptes)</router-link>
      </p>
    </div>

    <!-- Members list -->
    <div v-if="members.length > 0" class="members-list">
      <h3>Membres de l'organisation</h3>
      <table>
        <thead>
          <tr>
            <th>Utilisateur</th>
            <th>Email</th>
            <th>Type</th>
            <th>Date de création</th>
            <th v-if="accountType === 'organization_owner'">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="member in members" :key="member.user_id">
            <td>{{ member.username }}</td>
            <td>{{ member.email }}</td>
            <td>
              <span v-if="member.account_type === 'organization_owner'" class="badge-owner">👑 Propriétaire</span>
              <span v-else class="badge-member">👤 Membre</span>
            </td>
            <td>{{ formatDate(member.date_joined) }}</td>
            <td v-if="accountType === 'organization_owner'">
              <button 
                v-if="member.account_type !== 'organization_owner'" 
                @click="deleteSubAccount(member.user_id)" 
                class="btn-delete"
                :disabled="loading"
              >
                🗑️ Supprimer
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Error/Success messages -->
    <div v-if="message" :class="['message', messageType]">
      {{ message }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';

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

// Fetch organization info and members
async function fetchOrganizationData() {
  const username = sessionStorage.getItem('username');
  const token = sessionStorage.getItem('token');

  if (!username || !token) {
    router.push('/login');
    return;
  }

  try {
    // Get user info first
    const userResponse = await fetch('http://localhost:8080/api/getinfop', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, token })
    });

    if (!userResponse.ok) throw new Error('Failed to get user info');
    
    const userData = await userResponse.json();
    userInfo.value = userData;

    // Only fetch organization data if user has Pro or Enterprise
    if (userData.Nboffer >= 2) {
      const orgResponse = await fetch('http://localhost:8080/api/organization/members', {
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

  const username = sessionStorage.getItem('username');
  const token = sessionStorage.getItem('token');

  loading.value = true;
  message.value = '';

  try {
    const response = await fetch('http://localhost:8080/api/organization/create-subaccount', {
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

  const username = sessionStorage.getItem('username');
  const token = sessionStorage.getItem('token');

  loading.value = true;
  message.value = '';

  try {
    const response = await fetch('http://localhost:8080/api/organization/delete-member', {
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

onMounted(() => {
  fetchOrganizationData();
});
</script>

<style scoped>
.organization-panel {
  max-width: 1000px;
  margin: 20px auto;
  padding: 20px;
}

h2 {
  color: #2c3e50;
  margin-bottom: 20px;
}

.org-status {
  margin-bottom: 30px;
}

.status-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 4px 15px rgba(0,0,0,0.1);
}

.status-card h3 {
  margin: 0 0 10px 0;
  font-size: 24px;
}

.account-type {
  margin: 10px 0;
}

.badge-owner, .badge-member {
  padding: 5px 12px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: bold;
}

.badge-owner {
  background: #ffd700;
  color: #333;
}

.badge-member {
  background: rgba(255,255,255,0.3);
  color: white;
}

.member-count {
  margin: 10px 0 0 0;
  font-size: 18px;
}

.create-subaccount {
  background: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  margin-bottom: 30px;
}

.create-subaccount h3 {
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
  font-weight: 500;
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

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  padding: 14px 28px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.limit-reached {
  background: #fff3cd;
  border: 1px solid #ffc107;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 30px;
  text-align: center;
}

.limit-reached a {
  color: #667eea;
  font-weight: bold;
  text-decoration: none;
}

.limit-reached a:hover {
  text-decoration: underline;
}

.members-list {
  background: white;
  padding: 25px;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.members-list h3 {
  margin-top: 0;
  color: #2c3e50;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 15px;
}

thead {
  background: #f8f9fa;
}

th {
  padding: 12px;
  text-align: left;
  color: #555;
  font-weight: 600;
  border-bottom: 2px solid #e0e0e0;
}

td {
  padding: 12px;
  border-bottom: 1px solid #e0e0e0;
}

tbody tr:hover {
  background: #f8f9fa;
}

.btn-delete {
  background: #dc3545;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.2s;
}

.btn-delete:hover:not(:disabled) {
  background: #c82333;
}

.btn-delete:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.message {
  padding: 15px 20px;
  border-radius: 8px;
  margin-top: 20px;
  font-weight: 500;
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

/* Responsive */
@media (max-width: 768px) {
  .organization-panel {
    padding: 10px;
  }

  table {
    font-size: 14px;
  }

  th, td {
    padding: 8px;
  }

  .btn-primary {
    width: 100%;
  }
}
</style>
