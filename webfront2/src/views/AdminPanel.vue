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
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/admin/stats`, {
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
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/admin/users`, {
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
        backgroundColor: ['#808080', '#7c7c7c', '#404040'],
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
        backgroundColor: ['#727272', '#0a0a0a', '#7c7c7c', '#a7a7a7'],
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
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/admin/users/role`, {
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
    const response = await fetch(`${import.meta.env.VITE_APP_API}/api/admin/users/verify-contact`, {
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

<template>
  <div class="admin">
    <div class="admin-head">
      <h1>Administration</h1>
      <p>Gérez les utilisateurs et statistiques.</p>
    </div>

    <!-- Stats -->
    <div class="stats-grid">
      <div class="stat-card" v-for="(s, i) in [
        { icon: '👥', val: stats.total_users, label: 'Utilisateurs' },
        { icon: '✉', val: stats.total_verified_emails, label: 'Emails vérifiés' },
        { icon: '📱', val: stats.total_verified_phones, label: 'Téléphones vérifiés' },
        { icon: '📁', val: stats.total_files, label: 'Fichiers' },
        { icon: '📅', val: stats.total_calendar_events, label: 'Événements' },
        { icon: '⚠', val: stats.users_without_contacts, label: 'Sans contact' }
      ]" :key="i">
        <div class="stat-icon">{{ s.icon }}</div>
        <div class="stat-body">
          <span class="stat-val">{{ s.val }}</span>
          <span class="stat-lbl">{{ s.label }}</span>
        </div>
      </div>
    </div>

    <!-- Charts -->
    <div class="charts">
      <div class="chart-card">
        <h3>Vérifications</h3>
        <canvas ref="verificationChart"></canvas>
      </div>
      <div class="chart-card">
        <h3>Offres</h3>
        <canvas ref="offersChart"></canvas>
      </div>
    </div>

    <!-- Users table -->
    <div class="users-section">
      <div class="users-head">
        <h2>Utilisateurs</h2>
        <input v-model="searchQuery" type="text" placeholder="Rechercher..." class="search" />
      </div>
      <div class="table-wrap">
        <table class="table">
          <thead>
            <tr>
              <th>Utilisateur</th>
              <th>Email</th>
              <th>Téléphone</th>
              <th>Offre</th>
              <th>Rôle</th>
              <th>Inscrit le</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in filteredUsers" :key="user.user_id">
              <td class="cell-user">{{ user.username }}</td>
              <td class="cell-muted">
                <span v-if="user.email_verified" class="tag ok">✓</span>
                {{ user.email || '—' }}
              </td>
              <td class="cell-muted">
                <span v-if="user.phone_verified" class="tag ok">✓</span>
                {{ user.phonenumber || '—' }}
              </td>
              <td><span class="tag">{{ getOfferName(user.nboffer) }}</span></td>
              <td><span class="tag" :class="{ admin: user.role === 'admin' }">{{ user.role }}</span></td>
              <td class="cell-muted">{{ formatDate(user.date_joined) }}</td>
              <td class="cell-actions">
                <button @click="toggleRole(user)" class="act-btn" :class="{ danger: user.role === 'admin' }">
                  {{ user.role === 'admin' ? 'Rétrograder' : 'Promouvoir' }}
                </button>
                <button v-if="!user.email_verified" @click="verifyContact(user.user_id, 'email')" class="act-btn">Vérifier email</button>
                <button v-if="!user.phone_verified" @click="verifyContact(user.user_id, 'phone')" class="act-btn">Vérifier tel</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.admin { max-width: var(--content-max); margin: 0 auto; padding: 2.5rem 1.5rem 5rem; animation: fadeInUp .5s var(--ease) both; }
.admin-head { margin-bottom: 2rem; animation: fadeInDown .4s var(--ease) both; }
.admin-head h1 { font-size: 1.75rem; font-weight: 800; letter-spacing: -.025em; margin-bottom: .375rem; }
.admin-head p { font-size: .9375rem; color: var(--text-secondary); }

/* Stats */
.stats-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(180px, 1fr)); gap: 1rem; margin-bottom: 2rem; }
.stat-card { display: flex; align-items: center; gap: 1rem; padding: 1.25rem; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); transition: all var(--tb); animation: fadeInUp .4s var(--ease) both; }
.stat-card:hover { border-color: var(--border-strong); box-shadow: var(--sh-sm); transform: translateY(-2px); }
.stat-card:nth-child(2) { animation-delay: 40ms; }
.stat-card:nth-child(3) { animation-delay: 80ms; }
.stat-card:nth-child(4) { animation-delay: 120ms; }
.stat-card:nth-child(5) { animation-delay: 160ms; }
.stat-card:nth-child(6) { animation-delay: 200ms; }
.stat-icon { width: 44px; height: 44px; display: flex; align-items: center; justify-content: center; font-size: 1.25rem; border-radius: var(--r); background: var(--bg-tertiary); flex-shrink: 0; }
.stat-body { display: flex; flex-direction: column; }
.stat-val { font-size: 1.5rem; font-weight: 800; letter-spacing: -.02em; }
.stat-lbl { font-size: .75rem; color: var(--text-muted); }

/* Charts */
.charts { display: grid; grid-template-columns: 1fr 1fr; gap: 1.25rem; margin-bottom: 2rem; }
.chart-card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1.5rem; animation: fadeInUp .4s var(--ease) both; }
.chart-card:nth-child(2) { animation-delay: 80ms; }
.chart-card h3 { font-size: 1.0625rem; font-weight: 700; margin-bottom: 1rem; }
.chart-card canvas { max-width: 100%; }

/* Users */
.users-section { margin-bottom: 2rem; animation: fadeInUp .4s var(--ease) both; }
.users-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1rem; gap: 1rem; flex-wrap: wrap; }
.users-head h2 { font-size: 1.25rem; font-weight: 700; }
.search { width: 280px; max-width: 100%; padding: .5625rem .875rem; font-size: .875rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); transition: border-color var(--tb); }
.search:focus { outline: none; border-color: var(--ink); }
.search::placeholder { color: var(--text-muted); }
.table-wrap { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); overflow: auto; }
.table { width: 100%; border-collapse: collapse; }
.table th { padding: .75rem 1rem; font-size: .75rem; font-weight: 600; text-align: left; color: var(--text-secondary); background: var(--bg-secondary); border-bottom: 1px solid var(--border); text-transform: uppercase; letter-spacing: .03em; }
.table td { padding: .875rem 1rem; font-size: .875rem; color: var(--text-primary); border-bottom: 1px solid var(--border); }
.table tr { transition: background var(--tb); }
.table tr:hover td { background: var(--bg-secondary); }
.cell-user { font-weight: 600; }
.cell-muted { color: var(--text-secondary); }
.tag { display: inline-flex; padding: .125rem .5625rem; font-size: .75rem; font-weight: 600; border-radius: var(--r-full); background: var(--bg-tertiary); color: var(--text-secondary); }
.tag.ok { background: var(--ink); color: var(--paper); }
.tag.admin { background: var(--ink); color: var(--paper); }
.cell-actions { white-space: nowrap; }
.act-btn { padding: .375rem .75rem; font-size: .75rem; font-weight: 600; border: 1px solid var(--border-strong); background: var(--bg-primary); color: var(--text-primary); border-radius: var(--r-sm); cursor: pointer; transition: all var(--tb); margin-right: .375rem; }
.act-btn:hover { border-color: var(--ink); background: var(--bg-tertiary); transform: translateY(-1px); }
.act-btn.danger { color: var(--text-primary); }
.act-btn.danger:hover { border-color: var(--ink); }

@media (max-width: 768px) {
  .charts { grid-template-columns: 1fr; }
  .search { width: 100%; }
}
</style>
