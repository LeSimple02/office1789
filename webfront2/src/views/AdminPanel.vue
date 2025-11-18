<template>
  <div class="admin-panel">
    <div class="admin-header">
      <h1>🔐 Panel Administrateur</h1>
      <p>Bienvenue {{ gls().username }}</p>
    </div>

    <!-- Stats Cards -->
    <div class="stats-grid">
      <div class="stat-card blue">
        <div class="stat-icon">👥</div>
        <div class="stat-content">
          <h3>{{ stats.total_users }}</h3>
          <p>Utilisateurs</p>
        </div>
      </div>
      
      <div class="stat-card green">
        <div class="stat-icon">✉️</div>
        <div class="stat-content">
          <h3>{{ stats.total_verified_emails }}</h3>
          <p>Emails vérifiés</p>
        </div>
      </div>
      
      <div class="stat-card purple">
        <div class="stat-icon">📱</div>
        <div class="stat-content">
          <h3>{{ stats.total_verified_phones }}</h3>
          <p>Téléphones vérifiés</p>
        </div>
      </div>
      
      <div class="stat-card red">
        <div class="stat-icon">⚠️</div>
        <div class="stat-content">
          <h3>{{ stats.users_without_contacts }}</h3>
          <p>Sans contact vérifié</p>
        </div>
      </div>
      
      <div class="stat-card orange">
        <div class="stat-icon">📁</div>
        <div class="stat-content">
          <h3>{{ stats.total_files }}</h3>
          <p>Fichiers stockés</p>
        </div>
      </div>
      
      <div class="stat-card teal">
        <div class="stat-icon">📅</div>
        <div class="stat-content">
          <h3>{{ stats.total_calendar_events }}</h3>
          <p>Événements calendrier</p>
        </div>
      </div>
    </div>

    <!-- Charts -->
    <div class="charts-grid">
      <div class="chart-card">
        <h2>📊 Distribution des vérifications</h2>
        <canvas ref="verificationChart"></canvas>
      </div>
      
      <div class="chart-card">
        <h2>📈 Abonnements</h2>
        <canvas ref="offersChart"></canvas>
      </div>
    </div>

    <!-- Users Table -->
    <div class="users-section">
      <div class="section-header">
        <h2>👥 Gestion des utilisateurs</h2>
        <input 
          v-model="searchQuery" 
          type="text" 
          placeholder="🔍 Rechercher un utilisateur..." 
          class="search-input"
        />
      </div>
      
      <div class="table-container">
        <table class="users-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Username</th>
              <th>Email</th>
              <th>Recovery Email</th>
              <th>Téléphone</th>
              <th>Rôle</th>
              <th>Offre</th>
              <th>Inscrit le</th>
              <th>Actions</th>
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
                  {{ user.role === 'admin' ? '⬇️ Rétrograder' : '⬆️ Promouvoir' }}
                </button>
                <button v-if="!user.recovery_email_verified && user.recovery_email" 
                        @click="verifyContact(user.user_id, 'email')" 
                        class="btn-action btn-verify">
                  ✉️ Vérifier Email
                </button>
                <button v-if="!user.phonenumber_verified && user.phonenumber" 
                        @click="verifyContact(user.user_id, 'phone')" 
                        class="btn-action btn-verify">
                  📱 Vérifier Tel
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 2rem;
}

.admin-header {
  text-align: center;
  color: white;
  margin-bottom: 2rem;
}

.admin-header h1 {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
  font-weight: 700;
}

.admin-header p {
  font-size: 1.1rem;
  opacity: 0.9;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: white;
  border-radius: 16px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  font-size: 3rem;
  line-height: 1;
}

.stat-content h3 {
  font-size: 2rem;
  font-weight: 700;
  margin: 0;
}

.stat-content p {
  font-size: 0.9rem;
  color: #666;
  margin: 0;
}

.stat-card.blue .stat-icon { filter: hue-rotate(200deg); }
.stat-card.green .stat-icon { filter: hue-rotate(100deg); }
.stat-card.purple .stat-icon { filter: hue-rotate(270deg); }
.stat-card.red .stat-icon { filter: hue-rotate(0deg); }
.stat-card.orange .stat-icon { filter: hue-rotate(30deg); }
.stat-card.teal .stat-icon { filter: hue-rotate(160deg); }

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.chart-card {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.chart-card h2 {
  margin: 0 0 1.5rem 0;
  font-size: 1.3rem;
  color: #333;
}

.chart-card canvas {
  max-height: 300px;
}

.users-section {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
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
  font-size: 1.5rem;
  color: #333;
}

.search-input {
  padding: 0.75rem 1.25rem;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  font-size: 1rem;
  width: 300px;
  transition: all 0.3s ease;
}

.search-input:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.users-table th {
  padding: 1rem;
  text-align: left;
  font-weight: 600;
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
}

.users-table tbody tr.admin-row:hover {
  background-color: #fde68a;
}

.users-table td {
  padding: 1rem;
}

.username-badge {
  background: #667eea;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 8px;
  font-weight: 600;
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
  border-radius: 8px;
  font-size: 0.85rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #667eea;
  color: white;
}

.btn-action:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.btn-action.btn-danger {
  background: #ef4444;
}

.btn-action.btn-danger:hover {
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.btn-action.btn-verify {
  background: #10b981;
}

.btn-action.btn-verify:hover {
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
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
