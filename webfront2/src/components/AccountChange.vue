<script setup>
import { ref } from "vue"
import { gls } from "@/stores/global.js"

let dj = ref(0)
let lj = ref(0)
let domain = ref('')
let nboffer = ref(0)
let phone  = ref('')
let email = ref('')

let passwordt = ref('password')
let passwordt2 = ref('password')

let newusername = ref('')
let newphone = ref('')
let newemail = ref('')
let newoffer = ref(0)

let passf1 = ref('')
let passf2 = ref('')

let usernameR = ref(false)
let emailR = ref(false)
let phonenumberR = ref(false)

// Récupère les infos utilisateur (comme avant)
fetch(import.meta.env.VITE_APP_API_INFO_USER, {
  method: "POST",
  mode: "cors",
  body: JSON.stringify({ "username": gls().username, "token": gls().sessionT })
})
  .then(r => r.json())
  .then(a => {
    dj.value = a['DateJoined']
    domain.value = a['Domain']
    nboffer.value = a['Nboffer']
    email.value = a['Email']
    phone.value = a['PhoneNumber']
    lj.value = a["LastLogin"]
  })
  .catch(() => {
    // ignore fetch errors for now
  })

function send() {
  // Reset previous errors
  usernameR.value = false
  emailR.value = false
  phonenumberR.value = false

  if (passf1.value !== passf2.value) {
    return // you may show an error — handled in template
  }

  fetch(import.meta.env.VITE_APP_API_INFO_CHANGE, {
    method: "POST",
    mode: "cors",
    body: JSON.stringify({
      "lastusername": gls().username,
      "username": newusername.value,
      "phonenumber": newphone.value,
      "email": newemail.value,
      "nboffer": newoffer.value,
      "password": passf2.value,
      "token": gls().sessionT
    })
  })
    .then(r => r.json())
    .then(a => {
      if (newusername.value !== "" && a["Username"] !== "no" && a["Email"] !== "no" && a["Phone"] !== "no") {
        // succes: update cookies / store
        document.cookie = `name=${gls().username}; expires=Fri, 31 Dec 1900 23:59:59 GMT; Secure`
        document.cookie = `sessionToken=${gls().sessionT}; expires=Fri, 31 Dec 1900 23:59:59 GMT; Secure`
        gls().sessionT = a["Token"]
        gls().username = a["Username"]
        document.cookie = `name=${a["Username"]}; expires=${a["Expiry"]}; Secure`
        document.cookie = `sessionToken=${a["Token"]}; expires=${a["Expiry"]}; Secure`
        window.location.href = "/account"
      } else {
        // Show field-specific errors if present
        if (a["Username"]) usernameR.value = true
        if (a["Email"]) emailR.value = true
        if (a["Phone"]) phonenumberR.value = true
      }
    })
    .catch(() => {
      // handle error silently for now
    })
}

function togglePassword1() {
  passwordt.value = (passwordt.value === "password") ? "text" : "password"
}
function togglePassword2() {
  passwordt2.value = (passwordt2.value === "password") ? "text" : "password"
}
</script>

<template>
  <div class="container">
    <h1 id="title">{{ $t('infop') }}</h1>

    <section id="enso" class="card">
      <!-- top: back link + avatar -->
      <div class="top-row">
        <RouterLink class="back" to="/account">⬅️</RouterLink>

        <div id="pic">
          <img src="@/assets/napo.png" alt="avatar" />
          <p class="pic-label">{{ $t('picturep') }}</p>
        </div>
      </div>

      <!-- grid with two columns: label | field -->
      <form class="form-grid" @submit.prevent="send" novalidate>
        <!-- Username -->
        <div class="form-label">{{ $t('username') }}</div>
        <div class="form-field">
          <input v-model="newusername" :placeholder="gls().username" />
          <p v-if="usernameR" class="error small">{{ $t('dejaUP') }}</p>
        </div>

        <!-- Password (two inputs stacked in field column) -->
        <div class="form-label">{{ $t('password') }}</div>
        <div class="form-field">
          <div class="password-row">
            <input :type="passwordt" v-model="passf1" :placeholder="$t('passwordN')" />
            <button type="button" class="icon-btn" @click="togglePassword1" aria-label="toggle password">👁</button>
          </div>
          <div class="password-row">
            <input :type="passwordt2" v-model="passf2" :placeholder="$t('repassword')" />
            <button type="button" class="icon-btn" @click="togglePassword2" aria-label="toggle confirm">👁</button>
          </div>
          <p v-if="passf1 !== passf2 && passf2 !== ''" class="error small">{{ $t('passwordd') }}</p>
        </div>

        <!-- Doble (example static) -->
        <div class="form-label">{{ $t('doble') }}</div>
        <div class="form-field">
          <button class="small-btn">config</button>
        </div>

        <!-- Domain (read-only) -->
        <div class="form-label">{{ $t('domainy') }}</div>
        <div class="form-field">
          <input readonly :value="domain || '-'" />
        </div>

        <!-- Offer selector -->
        <div class="form-label">{{ $t('offery') }}</div>
        <div class="form-field">
          <select v-model="newoffer">
            <option value="0">0</option>
            <option value="1">1</option>
            <option value="2">2</option>
          </select>
          <RouterLink class="warn" to="/about">⚠️ {{ $t('About') }}</RouterLink>
        </div>

        <!-- Email -->
        <div class="form-label">{{ $t('emaily') }}</div>
        <div class="form-field">
          <input v-model="newemail" :placeholder="email" />
          <p v-if="emailR" class="error small">{{ $t('dejaEP') }}</p>
        </div>

        <!-- Phone -->
        <div class="form-label">{{ $t('phoney') }}</div>
        <div class="form-field">
          <input v-model="newphone" :placeholder="phone" />
          <p v-if="phonenumberR" class="error small">{{ $t('dejaPP') }}</p>
        </div>

        <!-- Last login (read-only) -->
        <div class="form-label">{{ $t('lastj') }}</div>
        <div class="form-field"><span>{{ lj ? (new Date(lj).toDateString()) : '-' }}</span></div>

        <!-- Date joined -->
        <div class="form-label">{{ $t('datej') }}</div>
        <div class="form-field"><span>{{ dj ? (new Date(dj).toDateString()) : '-' }}</span></div>
      </form>

      <!-- actions -->
      <div class="actions">
        <button type="button" class="btn primary" @click="send">✔ {{ $t('save') || 'Save' }}</button>
        <RouterLink class="btn ghost" to="/account">✖ {{ $t('cancel') || 'Cancel' }}</RouterLink>
      </div>
    </section>
  </div>
</template>

<style scoped>


* { box-sizing: border-box; font-family: 'Roboto', sans-serif; }

.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  
  padding: 28px 16px;
  background: #f4f6fb;
  min-height: 100vh;
  
}
.dark .container { background: #0f1220; }

#title {
  width: 100%;
  max-width: 900px;
  margin: 0 auto 16px;
  font-size: 26px;
  font-weight: 700;
  color: #222;
  text-align: left;
}
.dark #title { color: #eee; }

/* Card */
.card {
  width: 100%;
  max-width: 900px;
  background: #fff;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 8px 30px rgba(15,20,40,0.06);
}
.dark .card { background: #1f2230; color: #eee; box-shadow: 0 6px 24px rgba(0,0,0,0.4); }

.top-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 18px;
}
.back {
  color: #666;
  text-decoration: none;
  font-size: 18px;
}
#pic {
  display: flex;
  align-items: center;
  gap: 12px;
}
#pic img {
  width: 84px;
  height: 84px;
  border-radius: 999px;
  object-fit: cover;
  box-shadow: 0 4px 12px rgba(0,0,0,0.08);
}
.pic-label { font-weight: 600; color: #333; }
.dark .pic-label { color: #ddd; }

/* Grid: labels and fields aligned exactly opposite */
.form-grid {
  display: grid;
  grid-template-columns: 1fr 2fr; /* label | field */
  gap: 12px 20px;
  align-items: center;
  margin-top: 6px;
}

/* Labels (right-aligned so they're directly facing fields) */
.form-label {
  justify-self: end;
  text-align: right;
  color: #444;
  font-weight: 500;
  padding-right: 6px;
}
.dark .form-label { color: #ddd; }

/* Field column */
.form-field { justify-self: start; width: 100%; }

/* Inputs */
input[type="text"],
input[type="password"],
input[readonly],
select,
textarea {
  width: 100%;
  padding: 10px 12px;
  border-radius: 10px;
  border: 1px solid #e2e6ef;
  background: #fff;
  font-size: 14px;
}
.dark input, .dark select, .dark textarea, .dark input[readonly] {
  background: #2b2d3b; color: #eee; border-color: #444;
}

/* Password rows */
.password-row {
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
}
.icon-btn {
  border: none;
  background: none;
  cursor: pointer;
  padding: 6px;
  font-size: 18px;
}

/* small buttons inside form */
.small-btn {
  padding: 8px 12px;
  border-radius: 8px;
  border: 1px solid #e6e6e6;
  background: #f4f6fb;
  cursor: pointer;
}

/* error text */
.error { color: #d9534f; margin-top: 6px; }
.small { font-size: 12px; margin: 6px 0 0 0; }

/* Actions */
.actions {
  margin-top: 18px;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}
.btn {
  padding: 10px 16px;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  font-weight: 600;
}
.btn.primary {
  background: linear-gradient(135deg,#4facfe 0%,#00f2fe 100%);
  color: white;
  box-shadow: 0 6px 18px rgba(64,150,255,0.12);
}
.btn.ghost {
  background: transparent;
  color: #c0392b;
  border: 1px solid #eee;
}

/* Responsive: on small screens show stacked layout where label is above field */
@media (max-width: 820px) {
  .form-grid {
    grid-template-columns: 1fr; /* single column */
  }
  .form-label {
    justify-self: start;
    text-align: left;
    padding-right: 0;
  }
  #pic {
    justify-self: center;
  }
  .actions {
    justify-content: center;
  }
  #title { font-size: 20px; }
}

/* Very small screens tweak */
@media (max-width: 420px) {
  #pic img { width: 72px; height:72px; }
  .btn { padding: 10px 12px; font-size: 14px; }
}
</style>
