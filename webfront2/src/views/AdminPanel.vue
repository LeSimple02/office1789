<template>
  <div class="admin-panel">
    <div class="admin-header">
      <h1>{{ $t('adminPanel') }}</h1>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">👥</div>
        <div class="stat-content">
          <h3>{{ stats.total_users }}</h3>
          <p>{{ $t('users') }}</p>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">✉️</div>
        <div class="stat-content">
          <h3>{{ stats.total_verified_emails }}</h3>
          <p>{{ $t('verifiedEmails') }}</p>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">📱</div>
        <div class="stat-content">
          <h3>{{ stats.total_verified_phones }}</h3>
          <p>{{ $t('verifiedPhones') }}</p>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">⚠️</div>
        <div class="stat-content">
          <h3>{{ stats.users_without_contacts }}</h3>
          <p>{{ $t('withoutVerifiedContact') }}</p>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">📁</div>
        <div class="stat-content">
          <h3>{{ stats.total_files }}</h3>
          <p>{{ $t('storedFiles') }}</p>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon">📅</div>
        <div class="stat-content">
          <h3>{{ stats.total_calendar_events }}</h3>
          <p>{{ $t('calendarEvents') }}</p>
        </div>
      </div>
    </div>

    <!-- Charts -->
    <div class="charts-grid">
      <div class="chart-card">
        <h2>📊 {{ $t('verificationDistribution') }}</h2>
        <canvas ref="verificationChart"></canvas>
      </div>
      
      <div class="chart-card">
        <h2>📈 {{ $t('subscriptions') }}</h2>
        <canvas ref="offersChart"></canvas>
      </div>
    </div>

    <!-- Users Table -->
    <div class="users-section">
      <div class="section-header">
        <h2>👥 {{ $t('userManagement') }}</h2>
        <input 
          v-model="searchQuery" 
          type="text" 
          :placeholder="'🔍 ' + $t('searchUser')" 
          class="search-input"
        />
      </div>
      
      <div class="table-container">
        <table class="users-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>{{ $t('username') }}</th>
              <th>{{ $t('email') }}</th>
              <th>{{ $t('recoveryEmail') }}</th>
              <th>{{ $t('phone') }}</th>
              <th>{{ $t('role') }}</th>
              <th>{{ $t('offer') }}</th>
              <th>{{ $t('registeredOn') }}</th>
              <th>{{ $t('actions') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in filteredUsers" :key="user.user_id" :class="{ 'admin-row': user.role === 'admin' }">
              <td>{{ user.user_id }}</td>
              <td class="username-cell">
                <span class="username-badge">{{ user.username }}</span>
              </td>
              <td>{{ user.email }}</td>
              <td>
                <span v-if="user.recovery_email" class="email-cell">
                  {{ user.recovery_email }}
                  <span v-if="user.recovery_email_verified" class="verified-badge">✓</span>
                  <span v-else class="unverified-badge">✗</span>
                </span>
                <span v-else class="empty-cell">-</span>
              </td>
              <td>
                <span v-if="user.phonenumber" class="phone-cell">
                  {{ user.phonenumber }}
                  <span v-if="user.phonenumber_verified" class="verified-badge">✓</span>
                  <span v-else class="unverified-badge">✗</span>
                </span>
                <span v-else class="empty-cell">-</span>
              </td>
              <td>
                <span class="role-badge" :class="user.role">{{ user.role }}</span>
              </td>
              <td>
                <span class="offer-badge" :class="'offer-' + user.nboffer">
                  {{ getOfferName(user.nboffer) }}
                </span>
              </td>
              <td>{{ formatDate(user.date_joined) }}</td>
              <td class="actions-cell">
                <button @click="toggleRole(user)" class="btn-action" :class="{ 'btn-danger': user.role === 'admin' }">
                  {{ user.role === 'admin' ? '⬇️ ' + $t('demote') : '⬆️ ' + $t('promote') }}
                </button>
                <button v-if="!user.recovery_email_verified && user.recovery_email" 
                        @click="verifyContact(user.user_id, 'email')" 
                        class="btn-action btn-verify">
                  ✉️ {{ $t('verify') }} {{ $t('email') }}
                </button>
                <button v-if="!user.phonenumber_verified && user.phonenumber" 
                        @click="verifyContact(user.user_id, 'phone')" 
                        class="btn-action btn-verify">
                  📱 {{ $t('verify') }} {{ $t('phone') }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { gls } from '@/stores/global'
import { Chart, registerables } from 'chart.js'

Chart.register(...registerables)

const stats = ref({
  total_users: 0,
  total_verified_emails: 0,
  total_verified_phones: 0,
  users_without_contacts: 0,
  total_files: 0,
  total_calendar_events: 0
})

const users = ref([])
const searchQuery = ref('')
const verificationChart = ref(null)
const offersChart = ref(null)

const filteredUsers = computed(() => {
  if (!searchQuery.value) return users.value
  
  const query = searchQuery.value.toLowerCase()
  return users.value.filter(user => 
    user.username.toLowerCase().includes(query) ||
    user.email.toLowerCase().includes(query) ||
    (user.recovery_email && user.recovery_email.toLowerCase().includes(query)) ||
    (user.phonenumber && user.phonenumber.includes(query))
  )
})

onMounted(async () => {
  await loadStats()
  await loadUsers()
  createCharts()
})

async function loadStats() {
  try {
    const response = await fetch('http://localhost:8080/api/admin/stats', {
      headers: {
        'Authorization': gls().sessionT
      }
    })
    
    if (response.ok) {
      stats.value = await response.json()
    }
  } catch (error) {
    console.error('Error loading stats:', error)
  }
}

async function loadUsers() {
  try {
    const response = await fetch('http://localhost:8080/api/admin/users', {
      headers: {
        'Authorization': gls().sessionT
      }
    })
    
    if (response.ok) {
      users.value = await response.json()
    }
  } catch (error) {
    console.error('Error loading users:', error)
  }
}

function createCharts() {
  // Chart 1: Vérifications
  const ctx1 = verificationChart.value.getContext('2d')
  new Chart(ctx1, {
    type: 'doughnut',
    data: {
      labels: ['Emails vérifiés', 'Téléphones vérifiés', 'Sans contact'],
      datasets: [{
        data: [
          stats.value.total_verified_emails,
          stats.value.total_verified_phones,
          stats.value.users_without_contacts
        ],
        backgroundColor: ['#10b981', '#8b5cf6', '#ef4444'],
        borderWidth: 2,
        borderColor: '#fff'
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: true,
      plugins: {
        legend: {
          position: 'bottom',
          labels: { color: '#333', font: { size: 12 } }
        }
      }
    }
  })

  // Chart 2: Offres
  const offerCounts = [0, 0, 0, 0]
  users.value.forEach(user => {
    if (user.nboffer >= 0 && user.nboffer <= 3) {
      offerCounts[user.nboffer]++
    }
  })

  const ctx2 = offersChart.value.getContext('2d')
  new Chart(ctx2, {
    type: 'bar',
    data: {
      labels: ['Free', 'Standard', 'Professional', 'Enterprise'],
      datasets: [{
        label: 'Nombre d\'utilisateurs',
        data: offerCounts,
        backgroundColor: ['#6b7280', '#3b82f6', '#8b5cf6', '#f59e0b'],
        borderRadius: 8
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: true,
      scales: {
        y: {
          beginAtZero: true,
          ticks: { stepSize: 1 }
        }
      },
      plugins: {
        legend: { display: false }
      }
    }
  })
}

async function toggleRole(user) {
  const newRole = user.role === 'admin' ? 'user' : 'admin'
  
  if (!confirm(`Êtes-vous sûr de vouloir changer le rôle de ${user.username} en ${newRole} ?`)) {
    return
  }
  
  try {
    const response = await fetch('http://localhost:8080/api/admin/users/role', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': gls().sessionT
      },
      body: JSON.stringify({
        user_id: user.user_id,
        role: newRole
      })
    })
    
    if (response.ok) {
      user.role = newRole
      alert(`✅ Rôle mis à jour avec succès !`)
    } else {
      alert('❌ Erreur lors de la mise à jour du rôle')
    }
  } catch (error) {
    console.error('Error updating role:', error)
    alert('❌ Erreur réseau')
  }
}

async function verifyContact(userId, contactType) {
  try {
    const response = await fetch('http://localhost:8080/api/admin/users/verify-contact', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': gls().sessionT
      },
      body: JSON.stringify({
        user_id: userId,
        contact_type: contactType
      })
    })
    
    if (response.ok) {
      alert(`✅ Contact vérifié avec succès !`)
      await loadUsers()
    } else {
      alert('❌ Erreur lors de la vérification')
    }
  } catch (error) {
    console.error('Error verifying contact:', error)
    alert('❌ Erreur réseau')
  }
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

function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleDateString('fr-FR', { 
    day: '2-digit', 
    month: '2-digit', 
    year: 'numeric' 
  })
}
</script>

<style scoped>
.admin-panel {
  min-height: 100vh;
  background: #f5f5f5;
  padding: 2rem;
}

.admin-header {
  margin-bottom: 2rem;
}

.admin-header h1 {
  font-size: 2rem;
  font-weight: 400;
  color: #333;
  margin: 0;
  font-family: Roboto, sans-serif;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.stat-card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e5e7eb;
  transition: box-shadow 0.2s ease;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  font-size: 2.5rem;
  line-height: 1;
  opacity: 0.8;
}

.stat-content h3 {
  font-size: 1.8rem;
  font-weight: 500;
  margin: 0;
  color: #333;
}

.stat-content p {
  font-size: 0.9rem;
  color: #666;
  margin: 0;
  font-weight: 400;
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 2rem;
  margin-bottom: 3rem;
}

.chart-card {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e5e7eb;
}

.chart-card h2 {
  margin: 0 0 1.5rem 0;
  font-size: 1.2rem;
  color: #333;
  font-weight: 500;
}

.chart-card canvas {
  max-height: 300px;
}

.users-section {
  background: white;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  border: 1px solid #e5e7eb;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
  gap: 1rem;
}

.section-header h2 {
  margin: 0;
  font-size: 1.3rem;
  color: #333;
  font-weight: 500;
}

.search-input {
  padding: 0.6rem 1rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 0.95rem;
  width: 300px;
  transition: border-color 0.2s ease;
  font-family: Roboto, sans-serif;
}

.search-input:focus {
  outline: none;
  border-color: #999;
}

.table-container {
  overflow-x: auto;
}

.users-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.9rem;
}

.users-table thead {
  background: #f9fafb;
  color: #333;
  border-bottom: 2px solid #e5e7eb;
}

.users-table th {
  padding: 1rem;
  text-align: left;
  font-weight: 500;
}

.users-table tbody tr {
  border-bottom: 1px solid #e5e7eb;
  transition: background-color 0.2s ease;
}

.users-table tbody tr:hover {
  background-color: #f9fafb;
}

.users-table tbody tr.admin-row {
  background-color: #fef3c7;
  border-left: 3px solid #f59e0b;
}

.users-table tbody tr.admin-row:hover {
  background-color: #fde68a;
}

.users-table td {
  padding: 1rem;
}

.username-badge {
  background: #e5e7eb;
  color: #374151;
  padding: 0.25rem 0.75rem;
  border-radius: 6px;
  font-weight: 500;
}

.verified-badge {
  color: #10b981;
  font-weight: bold;
  margin-left: 0.5rem;
}

.unverified-badge {
  color: #ef4444;
  font-weight: bold;
  margin-left: 0.5rem;
}

.empty-cell {
  color: #9ca3af;
  font-style: italic;
}

.role-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 8px;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 0.75rem;
}

.role-badge.admin {
  background: #fbbf24;
  color: #78350f;
}

.role-badge.user {
  background: #e5e7eb;
  color: #374151;
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

.actions-cell {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.btn-action {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.2s ease;
  background: #4b5563;
  color: white;
}

.btn-action:hover {
  background: #374151;
}

.btn-action.btn-danger {
  background: #ef4444;
}

.btn-action.btn-danger:hover {
  background: #dc2626;
}

.btn-action.btn-verify {
  background: #10b981;
}

.btn-action.btn-verify:hover {
  background: #059669;
}

@media (max-width: 768px) {
  .admin-panel {
    padding: 1rem;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .charts-grid {
    grid-template-columns: 1fr;
  }

  .search-input {
    width: 100%;
  }

  .users-table {
    font-size: 0.8rem;
  }

  .users-table th,
  .users-table td {
    padding: 0.5rem;
  }
}
</style>
