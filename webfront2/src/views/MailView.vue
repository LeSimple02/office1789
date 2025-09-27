<script setup>
import { ref, onMounted, computed, watch } from "vue"
import { gls } from "@/stores/global"
import router from "@/router/index"

if (gls().log != 1) {
  router.push("login")
}

let curr = ref('')
let mails = ref([]) // chaque mail aura maintenant une propriété `folder`: 'inbox'|'draft'|'send'|'trash'
let selectedMail = ref(null)
let showCompose = ref(false)
let isMobile = ref(window.innerWidth <= 750)

// confirmation de lecture
let readReceipts = ref({})

// pour l'édition de brouillon / formulaire de composition
let composeForm = ref({
  to: '',
  cc: '',
  bcc: '',
  subject: '',
  body: ''
})
let editingDraftId = ref(null) // id du brouillon en cours d'édition (ou null)
let blockedSenders = ref([]) // liste d'expéditeurs bloqués

// email de l'utilisateur (si disponible dans le store sinon fallback)
let userEmail = ref(gls().email || "alice@example.com")

function changec() {
  const parts = window.location.href.split('/')
  curr.value = parts[4] ?? parts[3]
}

function openCompose(prefill = null) {
  // prefill peut être un mail (pour modifier un brouillon) ou null
  if (prefill) {
    composeForm.value.to = prefill.to || ''
    composeForm.value.cc = prefill.cc || ''
    composeForm.value.bcc = prefill.bcc || ''
    composeForm.value.subject = prefill.subject || ''
    composeForm.value.body = prefill.body || ''
    editingDraftId.value = prefill.id
  } else {
    // reset
    composeForm.value.to = ''
    composeForm.value.cc = ''
    composeForm.value.bcc = ''
    composeForm.value.subject = ''
    composeForm.value.body = ''
    editingDraftId.value = null
  }
  showCompose.value = true
}
function closeCompose() {
  showCompose.value = false
  // ne pas vider editingDraftId automatiquement (on le fait après sauvegarde/envoi)
}

// envoi d'un mail (depuis la popup). Si on était en train d'éditer un brouillon, le mail est mis en 'send'
function sendMail() {
  const now = new Date()
  const mailId = editingDraftId.value ?? Date.now()

  const newMail = {
    id: mailId,
    from: userEmail.value, // utiliser l'email réel de l'utilisateur
    to: composeForm.value.to,
    cc: composeForm.value.cc,
    bcc: composeForm.value.bcc,
    subject: composeForm.value.subject || " (sans sujet) ",
    body: composeForm.value.body,
    date: now.toLocaleString(),
    folder: "send"
  }

  if (editingDraftId.value) {
    // remplacer le brouillon existant par la version envoyée
    const idx = mails.value.findIndex(m => m.id === editingDraftId.value)
    if (idx !== -1) {
      mails.value[idx] = { ...mails.value[idx], ...newMail }
    } else {
      mails.value.push(newMail)
    }
    editingDraftId.value = null
  } else {
    mails.value.push(newMail)
  }

  readReceipts.value[mailId] = false

  closeCompose()

  setTimeout(() => {
    readReceipts.value[mailId] = true
  }, 3000)
}

// sauvegarder un brouillon (nouveau ou mise à jour)
function saveDraft() {
  const now = new Date()
  const mailId = editingDraftId.value ?? Date.now()

  const draftMail = {
    id: mailId,
    from: userEmail.value,
    to: composeForm.value.to,
    cc: composeForm.value.cc,
    bcc: composeForm.value.bcc,
    subject: composeForm.value.subject || " (sans sujet) ",
    body: composeForm.value.body,
    date: now.toLocaleString(),
    folder: "draft"
  }

  if (editingDraftId.value) {
    const idx = mails.value.findIndex(m => m.id === editingDraftId.value)
    if (idx !== -1) {
      mails.value[idx] = { ...mails.value[idx], ...draftMail }
    } else {
      mails.value.push(draftMail)
    }
    editingDraftId.value = null
  } else {
    mails.value.push(draftMail)
  }

  // garder la popup ouverte pour continuer à éditer ou fermer selon préférence : on ferme
  closeCompose()
}

// ouvrir la popup en mode édition de brouillon (utilisé aussi depuis le bouton dans la liste)
function editDraft() {
  if (!selectedMail.value) return
  if (selectedMail.value.folder !== 'draft') return
  openCompose(selectedMail.value)
}

// suppression : on déplace le mail vers la corbeille (folder = 'trash')
function deleteMail(m) {
  const idx = mails.value.findIndex(x => x.id === m.id)
  if (idx !== -1) {
    mails.value[idx].folder = 'trash'
  }
  if (selectedMail.value && selectedMail.value.id === m.id) {
    selectedMail.value = null
  }
}

// blocage d'un expéditeur : on ajoute à blockedSenders et on déplace ses mails vers la corbeille
function blockSender(m) {
  if (!m || !m.from) return
  // protection : ne pas bloquer si c'est l'utilisateur lui-même
  if (m.from === userEmail.value) return
  if (!blockedSenders.value.includes(m.from)) {
    blockedSenders.value.push(m.from)
  }
  // déplacer tous les mails de cet expéditeur dans la corbeille
  mails.value.forEach(mail => {
    if (mail.from === m.from) mail.folder = 'trash'
  })
  if (selectedMail.value && selectedMail.value.from === m.from) {
    selectedMail.value = null
  }
}

function openMail(m) {
  selectedMail.value = m
}

// computed : mails visibles selon le dossier courant
const filteredMails = computed(() => {
  const folderKey = (curr.value === 'mail' || !curr.value) ? 'inbox' : curr.value
  // si pas de folder défini sur un mail, on considère 'inbox'
  return mails.value.filter(m => (m.folder || 'inbox') === folderKey)
})

// WATCHER : quand on change de dossier (filteredMails change), mettre à jour selectedMail
watch(
  filteredMails,
  (newList) => {
    if (!newList || newList.length === 0) {
      // pas de mail dans la nouvelle boîte => vider la vue détail
      selectedMail.value = null
      return
    }
    // si le mail actuellement affiché n'appartient plus à la nouvelle liste, afficher le premier
    const stillThere = selectedMail.value && newList.some(m => m.id === selectedMail.value.id)
    if (!stillThere) {
      selectedMail.value = newList[0]
    }
  },
  { immediate: true }
)

onMounted(() => {
  window.addEventListener("resize", () => {
    isMobile.value = window.innerWidth <= 750
  })
})
</script>

<template>
  <div id="mail-area">
    <!-- Sidebar -->
    <aside id="sidebar">
      <button class="compose-btn" @click="openCompose">
        ➕ {{ $t('newMail') }}
      </button>
      <ul id="mail-menu" @click="changec">
        <li><RouterLink to="/mail" exact-active-class="active">{{$t('mail')}}</RouterLink></li>
        <li><RouterLink to="/mail/draft" exact-active-class="active">{{$t('draft')}}</RouterLink></li>
        <li><RouterLink to="/mail/send" exact-active-class="active">{{$t('send')}}</RouterLink></li>
        <li><RouterLink to="/mail/trash" exact-active-class="active">{{$t('trash')}}</RouterLink></li>
      </ul>
    </aside>

    <!-- Zone principale -->
    <main id="content">
      <header id="mail-header">
        <h2>{{$t(curr)}}</h2>
        <input class="search-input" :placeholder="$t('userids')" />
      </header>

      <div id="mail-split">
        <!-- Colonne gauche : liste -->
        <section id="mail-list">
          <div v-if="filteredMails.length === 0" id="empty-state">
            <img src="@/assets/cat.png" alt="cat" class="cat-img" />
            <p>{{$t('nothing')}}</p>
          </div>
          <ul v-else>
            <li
              v-for="m in filteredMails"
              :key="m.id"
              class="mail-item"
              @click="openMail(m)"
              :class="{ active: selectedMail && selectedMail.id === m.id }"
            >
              <div class="mail-header">
                <strong>{{ m.from }}</strong>
                <span class="mail-date">{{ m.date }}</span>
              </div>

              <div class="mail-subject">{{ m.subject }}</div>

              <!-- EDIT DIRECT POUR LES BROUILLONS : visible dans la liste sans ouvrir -->
              <button
                v-if="m.folder === 'draft'"
                class="list-edit-btn"
                @click.stop="openCompose(m)"
                :title="$t ? $t('editDraft') : 'Modifier le brouillon'"
              >
                ✏️
              </button>

              <span v-if="readReceipts[m.id]" class="mail-status">✅</span>
              <span v-else class="mail-status">⏳</span>
            </li>
          </ul>
        </section>

        <!-- Colonne droite : détail -->
        <section id="mail-detail">
          <div v-if="!selectedMail" class="empty-detail">
            <img src="@/assets/catE.png" alt="catE" class="catE-img" />
            <p>{{$t('affel')}}</p>
          </div>
          <div v-else class="mail-content">
            <div class="mail-actions">
              <button class="action-btn delete-btn" @click="deleteMail(selectedMail)">{{$t ? $t('delete') : 'Supprimer'}}</button>
              <button
                v-if="selectedMail && selectedMail.from !== userEmail"
                class="action-btn block-btn"
                @click="blockSender(selectedMail)"
              >
                {{$t ? $t('block') : 'Bloquer'}}
              </button>
              <button
                v-if="selectedMail.folder === 'draft'"
                class="action-btn edit-btn"
                @click="editDraft"
              >
                {{ $t ? $t('editDraft') : 'Modifier le brouillon' }}
              </button>
            </div>

            <h3>{{ selectedMail.subject }}</h3>
            <p class="mail-meta">
              <strong>{{ selectedMail.from }}</strong> – {{ selectedMail.date }}
            </p>
            <div class="mail-body">
              {{ selectedMail.body }}
            </div>
          </div>
        </section>
      </div>
    </main>

    <!-- Popup rédaction -->
    <div v-if="showCompose" :class="['compose-popup', isMobile ? 'mobile' : '']">
      <div class="compose-modal">
        <h3>{{ $t('composeMail') }}</h3>
        <input type="text" v-model="composeForm.to" placeholder="Destination" class="compose-input" />
        <input type="text" v-model="composeForm.cc" placeholder="CC" class="compose-input" />
        <input type="text" v-model="composeForm.bcc" placeholder="CCI" class="compose-input" />
        <input type="text" v-model="composeForm.subject" :placeholder="$t('subject')" class="compose-input" />
        <textarea v-model="composeForm.body" :placeholder="$t('yourM')" class="compose-textarea"></textarea>
        <div class="compose-actions">
          <label class="attach-label">
            <span class="attach-btn">📎 {{$t("insert")}}</span>
            <input type="file" class="attach-input" style="display:none;" />
          </label>
          <button @click="closeCompose" class="close-btn">{{$t('cancel')}}</button>

          <!-- Si on édite un brouillon, on propose "Mettre à jour" ; toujours possibilité d'enregistrer un brouillon -->
          <button class="close-btn" @click="saveDraft">{{ $t ? $t('saveDraft') : 'Enregistrer brouillon' }}</button>
          <button class="send-btn" @click="sendMail">{{$t('send')}}</button>
        </div>
      </div>
      <div class="compose-backdrop" @click="closeCompose"></div>
    </div>
  </div>
</template>

<style scoped>
/* Sidebar + layout */
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

/* Header */
#content {
  flex: 1;
  padding: 32px 40px;
  position: relative;
}
#mail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
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

/* Split: liste + détail */
#mail-split {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 20px;
  height: calc(100% - 60px);
}

/* Colonne gauche : liste */
#mail-list {
  overflow-y: auto;
  border-right: 1px solid #ddd;
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
.mail-item {
  background: #fff;
  border-bottom: 1px solid #eee;
  padding: 12px 14px;
  cursor: pointer;
  transition: background 0.2s;
  position: relative;
  list-style: none;
}
.mail-item:hover {
  background: #f0f0f7;
}
.mail-item.active {
  background: #e6f0ff;
}
.dark .mail-item {
  background: #2c2c3e;
  border-color: #444;
}
.mail-header {
  display: flex;
  justify-content: space-between;
  font-size: 0.85rem;
  color: #555;
}
.mail-date {
  font-size: 0.75rem;
  color: #999;
}
.mail-subject {
  font-weight: bold;
  margin: 4px 0;
}
.mail-status {
  font-size: 0.8rem;
  color: #888;
}

/* bouton edit directement dans la liste (pour brouillons) */
.list-edit-btn {
  position: absolute;
  right: 10px;
  top: 12px;
  border: none;
  background: transparent;
  cursor: pointer;
  font-size: 0.95rem;
}

/* Colonne droite : détail */
#mail-detail {
  padding: 20px;
}
.empty-detail {
  text-align: center;
  margin-top: 40px;
}
.catE-img {
  width: 140px;
  height: 140px;
  margin-bottom: 12px;
}
.empty-detail p {
  color: #555;
  font-size: 1.1rem;
}
.mail-content h3 {
  font-size: 1.3rem;
  margin-bottom: 10px;
}
.mail-meta {
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 20px;
}
.mail-body {
  font-size: 1rem;
  line-height: 1.5;
}

/* Actions spécifiques aux mails (supprimer / bloquer / modifier) */
.mail-actions {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}
.action-btn {
  padding: 6px 10px;
  border-radius: 6px;
  border: none;
  cursor: pointer;
  font-size: 0.9rem;
}
.delete-btn {
  background: #ffecec;
  color: #a00;
}
.block-btn {
  background: #fff4e6;
  color: #b65a00;
}
.edit-btn {
  background: #eef7ff;
  color: #0065a3;
}

/* Compose modal (inchangé structurellement) */
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
