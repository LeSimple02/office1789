<script setup>
import { ref } from "vue"
import { gls } from "@/stores/global.js"

let dj = ref(0)
let lj = ref(0)
let domain = ref(0)
let nboffer = ref(0)
let phone  = ref(0)
let email = ref(0)

fetch(import.meta.env.VITE_APP_API_INFO_USER, {
  method: "POST",
  mode: "cors",
  body: JSON.stringify({ 
    "username": gls().username, 
    "token": gls().sessionT 
  })
})
.then(res => res.json())
.then(data => {
  dj.value = data['DateJoined']
  domain.value = data['Domain']
  nboffer.value = data['Nboffer']
  email.value = data['Email']
  phone.value = data['PhoneNumber']
  lj.value = data["LastLogin"]
})
</script>

<template>
  <div class="profile-container">
    <header class="profile-header">
      <h1 class="title">{{ $t('infop') }}</h1>
      <RouterLink to="/account/edit" class="edit-link">{{ $t('edit') }}</RouterLink>
    </header>

    <section class="profile-card">
      <div class="profile-pic">
        <img src="@/assets/napo.png" alt="profile" />
        <span>{{ $t('picturep') }}</span>
      </div>

      <ul class="info-list">
        <li><strong>{{ $t('username') }} :</strong> {{ gls().username }}</li>
        <li><strong>{{ $t('password') }} :</strong> ●●●●●</li>
        <li><strong>{{ $t('doble') }} :</strong> ❌</li>
        <li><strong>{{ $t('domainy') }} :</strong> {{ domain }}</li>
        <li><strong>{{ $t('offery') }} :</strong> {{ nboffer }}</li>
        <li><strong>{{ $t('emaily') }} :</strong> {{ email }}</li>
        <li><strong>{{ $t('phoney') }} :</strong> {{ phone }}</li>
        <li><strong>{{ $t('lastj') }} :</strong> {{ new Date(lj).toDateString() }}</li>
        <li><strong>{{ $t('datej') }} :</strong> {{ new Date(dj).toDateString() }}</li>
      </ul>
    </section>
  </div>
</template>

<style scoped>


* { box-sizing: border-box; font-family: 'Roboto', sans-serif; }

.profile-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 16px;
  background: #f4f6fb;
  min-height: 100vh;
}
.dark .profile-container { background: #0f1220; }

.profile-header {
  width: 100%;
  max-width: 900px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
.title {
  font-size: 28px;
  font-weight: 700;
  color: #222;
}
.dark .title { color: #eee; }

.edit-link {
  color: grey;
  text-decoration: none;
  font-size: 16px;
}
.edit-link:hover { text-decoration: underline; }

.profile-card {
  width: 100%;
  max-width: 900px;
  background: #fff;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 8px 30px rgba(15,20,40,0.08);
  display: flex;
  flex-direction: column;
  gap: 24px;
}
.dark .profile-card {
  background: #1f2230;
  color: #eee;
  box-shadow: 0 6px 24px rgba(0,0,0,0.4);
}

.profile-pic {
  display: flex;
  flex-direction: column;
  align-items: center;
  font-weight: 600;
}
.profile-pic img {
  width: 140px;
  height: 140px;
  border-radius: 50%;
  margin-bottom: 12px;
}

.info-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: grid;
  gap: 12px;
  font-size: 1rem;
}
.info-list li strong {
  font-weight: 600;
  margin-right: 6px;
}

/* Responsive */
@media (max-width: 768px) {
  .title { font-size: 22px; }
  .profile-card { padding: 16px; }
  .profile-pic img { width: 100px; height: 100px; }
}
@media (max-width: 480px) {
  .title { font-size: 20px; text-align: center; }
  .profile-header { flex-direction: column; gap: 8px; }
  .info-list { font-size: 0.9rem; }
}
</style>
