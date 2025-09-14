<script setup>
import { ref, onMounted } from "vue"
import { gls } from "@/stores/global"
import router from "@/router/index"

if (gls().log != 1) {
  router.push("login")
}

let curr = ref('')
changec()
let showCompose = ref(false)
let isMobile = ref(window.innerWidth <= 750)

// confirmation de lecture (id → lu ?)
let readReceipts = ref({})

function changec() {
  const parts = window.location.href.split('/')
  curr.value = parts[4] ?? parts[3]
}

function openCompose() {
  showCompose.value = true
}

function closeCompose() {
  showCompose.value = false
}

function sendMail() {
  // logique d’envoi → ici on simule
  const mailId = Date.now()
  readReceipts.value[mailId] = false
  alert("📨 Mail envoyé ! (attente de confirmation de lecture)")
  closeCompose()

  // Simuler lecture au bout de 3 sec
  setTimeout(() => {
    readReceipts.value[mailId] = true
    alert("✅ Votre mail a été lu")
  }, 3000)
}

onMounted(() => {
  window.addEventListener("resize", () => {
    isMobile.value = window.innerWidth <= 750
  })
})
</script>

<template>
  <div id="mail-area">
    <aside id="sidebar">
      <button class="compose-btn" @click="openCompose">
        ➕ {{ $t('newMail') }}
      </button>
      <ul id="mail-menu" @click="changec">
        <li><RouterLink to="/mail" exact-active-class="active">{{$t('mail')}}</RouterLink></li>
        <li><RouterLink to="/mail/draft" exact-active-class="active">{{$t('draft')}}</RouterLink></li>
        <li><RouterLink to="/mail/send" exact-active-class="active">{{$t('send')}}</RouterLink></li>
        <li><RouterLink to="/mail/trash" exact-active-class="active">{{$t('trash')}}</RouterLink></li>
        <li><RouterLink to="/mail/sendem" exact-active-class="active">{{$t('sendem')}}</RouterLink></li>
      </ul>
    </aside>

    <main id="content">
      <header id="mail-header">
        <h2>{{$t(curr)}}</h2>
        <input class="search-input" :placeholder="$t('userids')" />
      </header>
      <section id="empty-state">
        <img src="@/assets/cat.png" alt="cat" class="cat-img" />
        <p>{{$t('nothing')}}</p>
      </section>
      <section id="aff-section">
        <img src="@/assets/catE.png" alt="catE" class="catE-img" />
        <p>{{$t('affel')}}</p>
      </section>
    </main>

    <!-- Popup de rédaction -->
    <div v-if="showCompose" :class="['compose-popup', isMobile ? 'mobile' : '']">
      <div class="compose-modal">
        <h3>{{ $t('composeMail') }}</h3>
        <input type="text" placeholder="Destination" class="compose-input" />
        <input type="text" placeholder="CC" class="compose-input" />
        <input type="text" placeholder="CCI" class="compose-input" />
        <input type="text" :placeholder="$t('subject')" class="compose-input" />
        <textarea :placeholder="$t('yourM')" class="compose-textarea"></textarea>
        <div class="compose-actions">
          <label class="attach-label">
            <span class="attach-btn">📎 {{$t("insert")}}</span>
            <input type="file" class="attach-input" style="display:none;" />
          </label>
          <button @click="closeCompose" class="close-btn">{{$t('cancel')}}</button>
          <button class="send-btn" @click="sendMail">{{$t('send')}}</button>
        </div>
      </div>
      <div class="compose-backdrop" @click="closeCompose"></div>
    </div>
  </div>
</template>

<style scoped>
/* Pièce jointe */
.attach-label {
  display: flex;
  align-items: center;
}
.attach-btn {
  background: #eee;
  color: #333;
  border-radius: 6px;
  padding: 8px 14px;
  cursor: pointer;
  font-size: 1rem;
  margin-right: 8px;
  border: 1px solid #ccc;
  transition: background 0.2s;
}
.attach-btn:hover {
  background: #e0e0e0;
}

#mail-area {
  display: flex;
  height: 100vh;
  background: #f6f7fb;
  font-family: 'Roboto', sans-serif;
}

.dark #mail-area {
  background: #1e1e2f;
  color: #ddd;
}

#sidebar {
  width: 220px;
  background: #fff;
  border-right: 1px solid #e0e0e0;
  padding: 24px 16px;
  display: flex;
  flex-direction: column;
}

.dark #sidebar {
  background: #2c2c3e;
  border-color: #444;
}

.compose-btn {
  background: -webkit-linear-gradient(30deg, blue, red);
  color: #fff;
  border: none;
  padding: 12px 0;
  border-radius: 8px;
  font-size: 1rem;
  margin-bottom: 24px;
  cursor: pointer;
  transition: background 0.2s;
}
.compose-btn:hover {
  background: -webkit-linear-gradient(30deg, blue, darkred);
}

#mail-menu {
  list-style: none;
  padding: 0;
  margin: 0;
}
#mail-menu li {
  margin-bottom: 16px;
}
#mail-menu a {
  text-decoration: none;
  color: #333;
  padding: 8px 12px;
  border-radius: 6px;
  display: block;
  transition: background 0.2s;
}
.dark #mail-menu a {
  color: #ddd;
}
#mail-menu a.active,
#mail-menu a.router-link-exact-active {
  background: -webkit-linear-gradient(30deg, blue, red);
  color: #fff;
}

#content {
  flex: 1;
  padding: 32px 40px;
  position: relative;
}

#mail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32px;
}
#mail-header h2 {
  font-size: 1.5rem;
  font-weight: 500;
}
.search-input {
  width: 220px;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #ccc;
}

#empty-state {
  text-align: center;
  margin-top: 60px;
}
.cat-img {
  width: 120px;
  height: 120px;
  margin-bottom: 16px;
}
#empty-state p {
  color: #888;
  font-size: 1.1rem;
}

#aff-section {
  text-align: center;
  margin-top: 40px;
}
.catE-img {
  width: 140px;
  height: 140px;
  margin-bottom: 12px;
}
#aff-section p {
  color: #555;
  font-size: 1.1rem;
}

.compose-popup {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}
.compose-modal {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.15);
  padding: 32px 28px;
  position: relative;
  z-index: 1001;
  width: 400px;
  display: flex;
  flex-direction: column;
}
.dark .compose-modal {
  background: #2c2c3e;
  color: #ddd;
}

/* Plein écran mobile */
.compose-popup.mobile .compose-modal {
  width: 100%;
  height: 100%;
  border-radius: 0;
  padding: 24px;
  justify-content: flex-start;
}

.compose-input {
  margin-bottom: 16px;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #ccc;
}
.compose-textarea {
  min-height: 100px;
  padding: 8px 12px;
  border-radius: 6px;
  border: 1px solid #ccc;
  margin-bottom: 16px;
  resize: vertical;
}
.compose-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
.close-btn, .send-btn {
  padding: 8px 18px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  font-size: 1rem;
}
.close-btn {
  background: #eee;
  color: #333;
}
.send-btn {
  background: -webkit-linear-gradient(30deg, blue, red);
  color: #fff;
}
.compose-backdrop {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,0.25);
  z-index: 1000;
}
</style>
