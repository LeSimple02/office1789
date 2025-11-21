<template>
  <div class="organization-container">
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
          <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
          <polyline points="9 22 9 12 15 12 15 22"></polyline>
        </svg>
      </div>
      
      <h1 class="hero-title">🏢 Gestion de l'organisation</h1>
      <p class="hero-subtitle">Créez et gérez les comptes de votre équipe</p>
    </div>

    <!-- Organization Status Card -->
    <div v-if="organizationInfo" class="status-card">
        <div class="org-header">
          <div class="org-icon-badge">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
              <path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"></path>
            </svg>
          </div>
          <div class="org-details">
            <h2 class="org-name">{{ organizationInfo.organization_name }}</h2>
            <span v-if="accountType === 'organization_owner'" class="role-badge owner">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
              Propriétaire
            </span>
            <span v-else class="role-badge member">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
              Membre
            </span>
          </div>
        </div>
        
        <div class="members-stats">
          <div class="stat-item">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
              <circle cx="9" cy="7" r="4"></circle>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
            </svg>
            <div class="stat-content">
              <span class="stat-value">{{ organizationInfo.current_count }}/{{ organizationInfo.max_members }}</span>
              <span class="stat-label">Membres actifs</span>
            </div>
          </div>
          
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: Math.min((organizationInfo.current_count / organizationInfo.max_members * 100), 100) + '%' }"></div>
          </div>
        </div>
    </div>

    <!-- Create Sub-Account Form -->
    <div v-if="accountType === 'organization_owner' && organizationInfo && organizationInfo.current_count < organizationInfo.max_members" class="form-card">
      <div class="form-header">
        <div class="form-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="8.5" cy="7" r="4"></circle>
            <line x1="20" y1="8" x2="20" y2="14"></line>
            <line x1="23" y1="11" x2="17" y2="11"></line>
          </svg>
        </div>
        <div>
          <h3 class="form-title">Créer un sous-compte</h3>
          <p class="form-subtitle">Ajoutez un nouveau membre à votre organisation</p>
        </div>
      </div>

      <div class="info-banner">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="12" y1="16" x2="12" y2="12"></line>
          <line x1="12" y1="8" x2="12.01" y2="8"></line>
        </svg>
        <span>Les sous-comptes héritent de votre offre et partagent le stockage de l'organisation</span>
      </div>

      <form @submit.prevent="createSubAccount" class="create-form">
        <div class="form-row">
          <div class="form-field">
            <label>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                <circle cx="12" cy="7" r="4"></circle>
              </svg>
              Nom d'utilisateur *
            </label>
            <input 
              type="text" 
              v-model="newSubAccount.username" 
              placeholder="john.doe"
              required 
              pattern="[a-z0-9._-]+"
            />
            <span class="field-hint">@office1789.com sera ajouté automatiquement</span>
          </div>

          <div class="form-field">
            <label>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
              </svg>
              Mot de passe *
            </label>
            <div class="password-field">
              <input 
                :type="showPassword ? 'text' : 'password'" 
                v-model="newSubAccount.password" 
                placeholder="••••••••"
                required 
                minlength="8"
                @input="checkPasswordStrength"
              />
              <button type="button" @click="togglePasswordVisibility" class="toggle-btn">
                <svg v-if="!showPassword" width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"/>
                </svg>
                <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 7c2.76 0 5 2.24 5 5 0 .65-.13 1.26-.36 1.83l2.92 2.92c1.51-1.26 2.7-2.89 3.43-4.75-1.73-4.39-6-7.5-11-7.5-1.4 0-2.74.25-3.98.7l2.16 2.16C10.74 7.13 11.35 7 12 7zM2 4.27l2.28 2.28.46.46A11.804 11.804 0 0 0 1 12c1.73 4.39 6 7.5 11 7.5 1.55 0 3.03-.3 4.38-.84l.42.42L19.73 22 21 20.73 3.27 3 2 4.27zM7.53 9.8l1.55 1.55c-.05.21-.08.43-.08.65 0 1.66 1.34 3 3 3 .22 0 .44-.03.65-.08l1.55 1.55c-.67.33-1.41.53-2.2.53-2.76 0-5-2.24-5-5 0-.79.2-1.53.53-2.2zm4.31-.78l3.15 3.15.02-.16c0-1.66-1.34-3-3-3l-.17.01z"/>
                </svg>
              </button>
            </div>
            <div v-if="newSubAccount.password" class="strength-indicator">
              <div class="strength-bar-bg">
                <div class="strength-bar-fill" :style="{ width: passwordStrength.width, backgroundColor: passwordStrength.color }"></div>
              </div>
              <span class="strength-text" :style="{ color: passwordStrength.color }">{{ passwordStrength.text }}</span>
            </div>
          </div>
        </div>

        <div class="form-row">
          <div class="form-field">
            <label>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"></path>
                <polyline points="22,6 12,13 2,6"></polyline>
              </svg>
              Email de récupération
            </label>
            <input 
              type="email" 
              v-model="newSubAccount.recoveryEmail" 
              placeholder="john@example.com"
            />
          </div>

          <div class="form-field">
            <label>
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path>
              </svg>
              Téléphone
            </label>
            <input 
              type="tel" 
              v-model="newSubAccount.phone" 
              placeholder="+33 6 XX XX XX XX"
            />
          </div>
        </div>

        <button type="submit" class="submit-btn" :disabled="loading">
          <svg v-if="!loading" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="20 6 9 17 4 12"></polyline>
          </svg>
          <div v-else class="btn-spinner"></div>
          {{ loading ? 'Création en cours...' : 'Créer le sous-compte' }}
        </button>
      </form>
    </div>

    <!-- Limit Reached -->
    <div v-if="accountType === 'organization_owner' && organizationInfo && organizationInfo.current_count >= organizationInfo.max_members" class="limit-card">
      <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="15" y1="9" x2="9" y2="15"></line>
        <line x1="9" y1="9" x2="15" y2="15"></line>
      </svg>
      <h3>Limite atteinte</h3>
      <p>Vous avez atteint la limite de {{ organizationInfo.max_members }} comptes.</p>
      <router-link v-if="userInfo.Nboffer < 3" to="/account" class="upgrade-link">
        Passer à Enterprise pour 20 comptes
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="5" y1="12" x2="19" y2="12"></line>
          <polyline points="12 5 19 12 12 19"></polyline>
        </svg>
      </router-link>
    </div>

    <!-- Members Table -->
    <div v-if="members.length > 0" class="members-card">
      <div class="members-header">
        <div class="members-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
            <circle cx="9" cy="7" r="4"></circle>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
            <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
          </svg>
        </div>
        <div>
          <h3 class="members-title">Membres de l'organisation</h3>
          <p class="members-subtitle">Gérez les accès de votre équipe</p>
        </div>
      </div>

      <div class="table-wrapper">
        <table class="members-table">
          <thead>
            <tr>
              <th>Utilisateur</th>
              <th>Email</th>
              <th>Rôle</th>
              <th>Offre</th>
              <th>Stockage</th>
              <th>Ajouté le</th>
              <th v-if="accountType === 'organization_owner'">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="member in members" :key="member.user_id" class="member-row">
              <td>
                <div class="member-cell">
                  <div class="member-avatar">
                    <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
                    </svg>
                  </div>
                  <span class="member-name">{{ member.username }}</span>
                </div>
              </td>
              <td class="email-cell">{{ member.email }}</td>
              <td>
                <span v-if="member.account_type === 'organization_owner'" class="role-tag owner">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
                  </svg>
                  Propriétaire
                </span>
                <span v-else class="role-tag member">Membre</span>
              </td>
              <td>
                <span v-if="member.account_type === 'organization_member'" class="plan-badge plan-badge-subaccount">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="margin-right: 4px;">
                    <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
                    <circle cx="9" cy="7" r="4"></circle>
                    <path d="M23 21v-2a4 4 0 0 0-3-3.87"></path>
                    <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
                  </svg>
                  Sous-compte
                </span>
                <span v-else-if="member.account_type === 'organization_owner'" class="plan-badge" :class="'plan-badge-' + getPlanClass(userInfo.Nboffer || 3)">
                  {{ getPlanName(userInfo.Nboffer || 3) }}
                </span>
                <span v-else class="plan-badge" :class="'plan-badge-' + getPlanClass(member.nboffer || 0)">
                  {{ getPlanName(member.nboffer || 0) }}
                </span>
              </td>
              <td class="storage-cell">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="display: inline; margin-right: 6px;">
                  <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"></path>
                </svg>
                <span v-if="member.account_type === 'organization_member'">
                  Partagé
                </span>
                <span v-else-if="member.account_type === 'organization_owner'">
                  {{ getStorageSize(userInfo.Nboffer || 3) }}
                </span>
                <span v-else>
                  {{ getStorageSize(member.nboffer || 0) }}
                </span>
              </td>
              <td class="date-cell">{{ formatDate(member.date_joined) }}</td>
              <td v-if="accountType === 'organization_owner'">
                <button 
                  v-if="member.account_type !== 'organization_owner'" 
                  @click="deleteSubAccount(member.user_id)" 
                  class="delete-btn"
                  :disabled="loading"
                >
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3 6 5 6 21 6"></polyline>
                    <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path>
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Alert Messages -->
    <transition name="alert-slide">
      <div v-if="message" :class="['alert-message', messageType === 'success' ? 'success' : 'error']">
        <svg v-if="messageType === 'success'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="20 6 9 17 4 12"></polyline>
        </svg>
        <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="12" r="10"></circle>
          <line x1="15" y1="9" x2="9" y2="15"></line>
          <line x1="9" y1="9" x2="15" y2="15"></line>
        </svg>
        <span>{{ message }}</span>
      </div>
    </transition>
  </div>
</template>

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
  let color = '#6c757d';
  
  if (password.length >= 8) strength++;
  if (password.length >= 12) strength++;
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++;
  if (/\d/.test(password)) strength++;
  if (/[^a-zA-Z0-9]/.test(password)) strength++;
  
  if (strength <= 2) {
    text = 'Faible';
    color = '#dc3545';
  } else if (strength === 3) {
    text = 'Moyen';
    color = '#ffc107';
  } else if (strength === 4) {
    text = 'Fort';
    color = '#28a745';
  } else {
    text = 'Très fort';
    color = '#155724';
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

  const username = gls().username;
  const token = gls().sessionT;

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

  const username = gls().username;
  const token = gls().sessionT;

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

<style scoped>
/* Animation d'entrée */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes slideIn {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Container */
.organization-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  animation: fadeIn 0.5s ease-in;
}

.organization-container > * {
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

/* Status Card */
.status-card {
  position: relative;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  padding: 40px;
  margin-bottom: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  overflow: hidden;
  animation: slideIn 0.6s ease-out;
}

.dark .status-card {
  background: #1C1C1E;
}

.org-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 32px;
}

.org-icon-badge {
  width: 72px;
  height: 72px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 24px rgba(0, 0, 255, 0.3);
}

.org-icon-badge svg {
  color: white;
}

.org-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.org-name {
  font-family: 'Roboto', sans-serif;
  font-size: 24px;
  font-weight: 700;
  color: #222;
  margin: 0 0 8px 0;
}

.dark .org-name {
  color: white;
}

.role-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 12px;
  font-family: 'Roboto', sans-serif;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.3px;
  white-space: nowrap;
  flex-shrink: 0;
}

.role-badge.owner {
  background: linear-gradient(135deg, #FFD700 0%, #FFA500 100%);
  color: #8B4000;
  box-shadow: 0 4px 12px rgba(255, 215, 0, 0.3);
}

.role-badge.member {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

.members-stats {
  background: rgba(245, 245, 247, 0.6);
  border-radius: 20px;
  padding: 24px;
}

.dark .members-stats {
  background: rgba(255, 255, 255, 0.05);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 12px;
}

.dark .stat-item {
  background: rgba(255, 255, 255, 0.05);
}

.stat-item svg {
  color: blue;
  flex-shrink: 0;
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 6px;
  flex: 1;
}

.stat-value {
  font-family: 'Roboto', sans-serif;
  font-size: 24px;
  font-weight: 700;
  color: #222;
}

.dark .stat-value {
  color: white;
}

.stat-label {
  font-family: 'Roboto', sans-serif;
  font-size: 14px;
  color: #6c757d;
  font-weight: 500;
}

.dark .stat-label {
  color: #9ca3af;
}

.progress-bar {
  height: 8px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  overflow: hidden;
}

.dark .progress-bar {
  background: rgba(255, 255, 255, 0.1);
}

.progress-fill {
  height: 100%;
  background: -webkit-linear-gradient(90deg, blue, red);
  border-radius: 4px;
  transition: width 0.6s ease;
}

/* Form Card */
.form-card {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  padding: 40px;
  margin-bottom: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  animation: slideIn 0.7s ease-out;
}

.dark .form-card {
  background: #1C1C1E;
}

.form-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 24px;
}

.form-icon {
  width: 56px;
  height: 56px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(0, 0, 255, 0.25);
}

.form-icon svg {
  color: white;
}

.form-title {
  font-family: 'Roboto', sans-serif;
  font-size: 20px;
  font-weight: 700;
  color: #222;
  margin: 0 0 4px 0;
}

.dark .form-title {
  color: white;
}

.form-subtitle {
  font-family: 'Roboto', sans-serif;
  font-size: 14px;
  color: #6c757d;
  margin: 0;
}

.dark .form-subtitle {
  color: #9ca3af;
}

.info-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background: rgba(255, 193, 7, 0.1);
  border-left: 4px solid #FFC107;
  border-radius: 12px;
  margin-bottom: 32px;
  font-family: 'Roboto', sans-serif;
  font-size: 14px;
  color: #856404;
}

.dark .info-banner {
  background: rgba(255, 193, 7, 0.15);
  color: #FFD54F;
}

.info-banner svg {
  flex-shrink: 0;
  color: #FFC107;
}

/* Form */
.create-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 24px;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-field label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-family: 'Roboto', sans-serif;
  font-size: 14px;
  font-weight: 600;
  color: #495057;
}

.dark .form-field label {
  color: #e5e7eb;
}

.form-field label svg {
  color: blue;
  flex-shrink: 0;
}

.form-field input {
  width: 100%;
  height: 48px;
  padding: 0 16px;
  border: 2px solid #e9ecef;
  border-radius: 12px;
  font-family: 'Roboto', sans-serif;
  font-size: 15px;
  background: white;
  transition: all 0.3s ease;
  box-sizing: border-box;
  color: #495057;
}

.dark .form-field input {
  background: #2C2C2E;
  border-color: rgba(255, 255, 255, 0.15);
  color: white;
}

.form-field input:focus {
  outline: none;
  border-color: blue;
  background: white;
  box-shadow: 0 0 0 4px rgba(0, 0, 255, 0.1);
}

.dark .form-field input:focus {
  background: #2C2C2E;
  border-color: #6ec6ff;
  box-shadow: 0 0 0 4px rgba(110, 198, 255, 0.2);
}

.field-hint {
  font-family: 'Roboto', sans-serif;
  font-size: 13px;
  color: #6c757d;
  margin-top: 4px;
}

.dark .field-hint {
  color: #9ca3af;
}

.password-field {
  position: relative;
}

.password-field input {
  padding-right: 52px;
}

.toggle-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #6c757d;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.dark .toggle-btn {
  color: #9ca3af;
}

.toggle-btn:hover {
  background: rgba(0, 0, 0, 0.05);
  color: #667eea;
}

.dark .toggle-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.strength-indicator {
  margin-top: 12px;
}

.strength-bar-bg {
  height: 6px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 8px;
}

.dark .strength-bar-bg {
  background: rgba(255, 255, 255, 0.1);
}

.strength-bar-fill {
  height: 100%;
  border-radius: 3px;
  transition: all 0.3s ease;
}

.strength-text {
  font-family: 'Roboto', sans-serif;
  font-size: 13px;
  font-weight: 600;
}

/* Submit Button */
.submit-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  width: 100%;
  height: 56px;
  padding: 0 32px;
  background: -webkit-linear-gradient(30deg, blue, red);
  color: white;
  border: none;
  border-radius: 20px;
  font-family: 'Roboto', sans-serif;
  font-size: 1.15rem;
  font-weight: 700;
  letter-spacing: 0.5px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.25);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(-1px);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-spinner {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

/* Limit Card */
.limit-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  padding: 48px 32px;
  background: rgba(255, 235, 59, 0.1);
  border: 2px dashed #FFC107;
  border-radius: 32px;
  text-align: center;
  margin-bottom: 32px;
}

.limit-card svg {
  color: #FFC107;
}

.limit-card h3 {
  font-size: 20px;
  font-weight: 700;
  color: #1a1d2e;
  margin: 0;
}

.dark .limit-card h3 {
  color: white;
}

.limit-card p {
  font-size: 15px;
  color: #6c757d;
  margin: 0;
}

.dark .limit-card p {
  color: #9ca3af;
}

.upgrade-link {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  text-decoration: none;
  border-radius: 12px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.upgrade-link:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(102, 126, 234, 0.4);
}

/* Members Card */
.members-card {
  background: rgba(245, 245, 247, 0.85);
  border-radius: 32px;
  padding: 40px;
  margin-bottom: 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.10);
  animation: slideIn 0.8s ease-out;
}

.dark .members-card {
  background: #1C1C1E;
}

.members-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
}

.members-icon {
  width: 56px;
  height: 56px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(0, 0, 255, 0.25);
}

.members-icon svg {
  color: white;
}

.members-title {
  font-family: 'Roboto', sans-serif;
  font-size: 20px;
  font-weight: 700;
  color: #222;
  margin: 0 0 4px 0;
}

.dark .members-title {
  color: white;
}

.members-subtitle {
  font-family: 'Roboto', sans-serif;
  font-size: 14px;
  color: #6c757d;
  margin: 0;
}

.dark .members-subtitle {
  color: #9ca3af;
}

/* Table */
.table-wrapper {
  overflow-x: auto;
  border-radius: 16px;
}

.members-table {
  width: 100%;
  border-collapse: collapse;
  background: #f8f9fa;
  border-radius: 12px;
  overflow: hidden;
}

.dark .members-table {
  background: #2C2C2E;
}

.members-table thead {
  background: #e9ecef;
}

.dark .members-table thead {
  background: rgba(255, 255, 255, 0.05);
}

.members-table th {
  padding: 16px 20px;
  text-align: left;
  font-size: 13px;
  font-weight: 700;
  color: #667eea;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.members-table td {
  padding: 20px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
  font-size: 14px;
  color: #495057;
}

.dark .members-table td {
  border-top-color: rgba(255, 255, 255, 0.05);
  color: #e5e7eb;
}

.member-row {
  transition: all 0.3s ease;
}

.member-row:hover {
  background: rgba(102, 126, 234, 0.05);
}

.member-cell {
  display: flex;
  align-items: center;
  gap: 12px;
}

.member-avatar {
  width: 36px;
  height: 36px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.member-avatar svg {
  color: white;
}

.member-name {
  font-family: 'Roboto', sans-serif;
  font-weight: 600;
  color: #222;
}

.dark .member-name {
  color: white;
}

.email-cell {
  font-family: 'Roboto', sans-serif;
  color: #6c757d;
}

.dark .email-cell {
  color: #9ca3af;
}

.role-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 8px;
  font-family: 'Roboto', sans-serif;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.role-tag.owner {
  background: rgba(255, 215, 0, 0.15);
  color: #FFA500;
}

.role-tag.member {
  background: rgba(102, 126, 234, 0.1);
  color: #667eea;
}

/* Plan Badges */
.plan-badge {
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 8px;
  font-family: 'Roboto', sans-serif;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.plan-badge-free {
  background: #e9ecef;
  color: #495057;
}

.dark .plan-badge-free {
  background: #3a3d4a;
  color: #9ca3af;
}

.plan-badge-standard {
  background: #cfe2ff;
  color: #084298;
}

.dark .plan-badge-standard {
  background: rgba(13, 110, 253, 0.25);
  color: #6ea8fe;
}

.plan-badge-professional {
  background: #d1e7dd;
  color: #0f5132;
}

.dark .plan-badge-professional {
  background: rgba(25, 135, 84, 0.25);
  color: #75b798;
}

.plan-badge-enterprise {
  background: linear-gradient(135deg, rgba(255, 215, 0, 0.2) 0%, rgba(255, 165, 0, 0.2) 100%);
  color: #cc7a00;
  font-weight: 700;
}

.dark .plan-badge-enterprise {
  background: linear-gradient(135deg, rgba(255, 215, 0, 0.25) 0%, rgba(255, 165, 0, 0.25) 100%);
  color: #FFD700;
}

/* Sub-account Badge */
.plan-badge-subaccount {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.2) 0%, rgba(118, 75, 162, 0.2) 100%);
  color: #667eea;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
}

.dark .plan-badge-subaccount {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.25) 0%, rgba(118, 75, 162, 0.25) 100%);
  color: #a78bfa;
}

/* Storage Cell */
.storage-cell {
  font-family: 'Roboto', sans-serif;
  color: #495057;
  font-size: 13px;
  font-weight: 500;
}

.dark .storage-cell {
  color: #e5e7eb;
}

.date-cell {
  font-family: 'Roboto', sans-serif;
  color: #6c757d;
  font-size: 13px;
}

.dark .date-cell {
  color: #9ca3af;
}

.delete-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: rgba(220, 53, 69, 0.1);
  color: #dc3545;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.delete-btn:hover:not(:disabled) {
  background: #dc3545;
  color: white;
  transform: scale(1.1);
}

.delete-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Alert Message */
.alert-message {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 24px;
  border-radius: 16px;
  font-family: 'Roboto', sans-serif;
  font-size: 15px;
  font-weight: 500;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  animation: slideIn 0.4s ease-out;
}

.alert-message.success {
  background: linear-gradient(135deg, #d4edda 0%, #c3e6cb 100%);
  color: #155724;
}

.alert-message.error {
  background: linear-gradient(135deg, #f8d7da 0%, #f5c6cb 100%);
  color: #721c24;
}

.alert-message svg {
  flex-shrink: 0;
}

.alert-slide-enter-active,
.alert-slide-leave-active {
  transition: all 0.3s ease;
}

.alert-slide-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.alert-slide-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

/* Responsive */
@media (max-width: 768px) {
  .organization-container {
    padding: 16px;
  }

  .hero-section,
  .status-card,
  .form-card,
  .members-card {
    padding: 24px 20px;
    border-radius: 24px;
  }

  .back-link {
    position: static;
    margin-bottom: 16px;
  }

  .form-row {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .org-header {
    flex-direction: column;
    text-align: center;
  }

  .members-table {
    font-size: 13px;
  }

  .members-table th,
  .members-table td {
    padding: 12px;
  }
}
</style>
