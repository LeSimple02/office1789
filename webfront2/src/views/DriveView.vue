<template>
  <div id="drive-v2">
    <!-- SIDEBAR -->
    <aside class="dv-sidebar">
      <div class="dv-user">
        <div class="dv-avatar">{{ userName[0] }}</div>
        <div class="dv-user-meta">
          <div class="dv-name">{{ userName }}</div>
          <div class="dv-email">{{ userEmail }}</div>
        </div>
      </div>
      <div class="dv-actions">
        <button class="dv-btn primary" @click="openUploadModal">📤 Upload</button>
        <button class="dv-btn" @click="showNewFolderModal = true">📁 Nouveau dossier</button>
      </div>
      <nav class="dv-nav" aria-label="navigation drive">
        <button :class="['nav-item', { active: curr === 'drive' }]" @click="changec('drive')">📂 Mes fichiers</button>
        <button :class="['nav-item', { active: curr === 'shared' }]" @click="changec('shared')">🤝 Partagés</button>
        <button :class="['nav-item', { active: curr === 'trash' }]" @click="changec('trash')">🗑️ Corbeille</button>
      </nav>
      <div class="dv-footer">
        <small>{{ files.length }} éléments</small>
      </div>
    </aside>
    <!-- MAIN -->
    <main class="dv-main">
      <header class="dv-header">
        <div class="dv-breadcrumbs">
          <button class="crumb" @click="currentPath = '/'; changec('drive')">Root</button>
          <template v-for="(part, i) in currentPath.split('/').filter(p => p)" :key="i">
            <span class="slash">/</span>
            <button class="crumb" @click="currentPath = currentPath.split('/').slice(0, i+1).join('/') + '/'">{{ part }}</button>
          </template>
        </div>
        <div class="dv-tools">
          <input class="dv-search" v-model="searchQuery" placeholder="Rechercher..." />
          <button class="dv-btn small" @click="selectFirst">Sélectionner 1er</button>
          <button class="dv-btn small" @click="clearSelection">Effacer</button>
        </div>
      </header>
      <section class="dv-body">
        <!-- LEFT: list -->
        <aside
          class="dv-list"
          @dragover.prevent="dragOver = true"
          @dragleave.prevent="dragOver = false"
          @drop.prevent="onDrop"
          :class="{ 'drag-over': dragOver }"
        >
          <div v-if="filteredFiles.length === 0" class="dv-empty">
            <div class="dv-empty-illustration" aria-hidden>
              <svg width="80" height="80" viewBox="0 0 24 24" fill="none">
                <rect width="24" height="24" rx="4" fill="#f0f4ff"/>
                <path d="M6 10c1.333-2 4-2 6 0 2 2.667 4.667 2.667 6 0" stroke="#a0b4ff" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <p class="dv-empty-text">Aucun fichier ici</p>
            <p class="dv-empty-sub">Glisse-dépose un fichier ici ou clique sur <strong>Upload</strong>.</p>
          </div>
          <ul v-else class="dv-items" role="list">
            <li
              v-for="it in filteredFiles"
              :key="it.id"
              class="dv-item"
              :class="{ active: selectedFile && selectedFile.id === it.id }"
              @click="it.type === 'folder' ? openFolder(it) : openFile(it)"
            >
              <div class="left">
                <div class="icon" v-if="it.type === 'folder'">📁</div>
                <div class="icon" v-else-if="(it.mime || '').startsWith('image/')">🖼️</div>
                <div class="icon" v-else-if="(it.mime || '').startsWith('video/')">🎬</div>
                <div class="icon" v-else-if="(it.mime || '').startsWith('audio/')">🔊</div>
                <div class="icon" v-else>📄</div>
                <div class="meta">
                  <div class="name">{{ it.name }}</div>
                  <div class="sub">{{ it.type === 'folder' ? 'Dossier' : (it.sizeLabel || humanSize(it.size)) }} • {{ it.date }}</div>
                </div>
              </div>
              <div class="actions">
                <button class="dv-mini" @click.stop="startRename(it)" title="Renommer">✏️</button>
                <button class="dv-mini" @click.stop="downloadFile(it)" title="Télécharger">⬇️</button>
                <button v-if="it.folder !== 'trash'" class="dv-mini" @click.stop="moveToTrash([it.id])" title="Mettre à la corbeille">🗑️</button>
              </div>
            </li>
          </ul>
        </aside>
        <!-- RIGHT: detail / preview -->
        <section class="dv-detail">
          <div v-if="!selectedFile" class="dv-detail-empty">
            <div class="dv-empty-illustration-center" aria-hidden>
              <svg width="100" height="100" viewBox="0 0 24 24" fill="none">
                <rect width="24" height="24" rx="4" fill="#fff6f0"/>
                <path d="M8 14s1.5-4 4-4 4 4 4 4" stroke="#ffd4b3" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <p class="dv-empty-text-centered">Sélectionne un fichier pour voir l'aperçu</p>
          </div>
          <div v-else class="dv-detail-card">
            <div class="detail-top">
              <div class="detail-info">
                <h3 class="detail-name">{{ selectedFile.name }}</h3>
                <div class="detail-meta">{{ selectedFile.mime || (selectedFile.type === 'folder' ? 'Dossier' : '') }} • {{ selectedFile.date }}</div>
              </div>
              <div class="detail-actions">
                <button class="dv-btn" @click="openAction(selectedFile)">Ouvrir</button>
                <button class="dv-btn green" @click="downloadFile(selectedFile)">Télécharger</button>
                <button class="dv-btn" @click="startRename(selectedFile)">Renommer</button>
                <button class="dv-btn purple" @click="shareFile(selectedFile)">Partager</button>
                <button v-if="curr !== 'trash'" class="dv-btn danger" @click="moveToTrash([selectedFile.id])">Mettre à la corbeille</button>
                <template v-else>
                  <button class="dv-btn" @click="restoreFile(selectedFile)">Restaurer</button>
                  <button class="dv-btn danger" @click="confirmPermanentDelete(selectedFile)">Supprimer définitivement</button>
                </template>
              </div>
            </div>
            <div class="detail-body">
              <div class="preview">
                <!-- image -->
                <template v-if="selectedFile.url && (selectedFile.mime || '').startsWith('image/')">
                  <img :src="selectedFile.url" alt="preview" class="preview-img" :key="selectedFile.id" />
                </template>
                <!-- video -->
                <template v-else-if="selectedFile.url && (selectedFile.mime || '').startsWith('video/')">
                  <video controls class="preview-video" :key="selectedFile.id">
                    <source :src="selectedFile.url" :type="selectedFile.mime" />
                    Ton navigateur ne supporte pas la lecture vidéo.
                  </video>
                </template>
                <!-- audio -->
                <template v-else-if="selectedFile.url && (selectedFile.mime || '').startsWith('audio/')">
                  <audio controls class="preview-audio" :key="selectedFile.id">
                    <source :src="selectedFile.url" :type="selectedFile.mime" />
                    Ton navigateur ne supporte pas la lecture audio.
                  </audio>
                </template>
                <!-- texte -->
                <template v-else-if="selectedFile.text">
                  <pre class="preview-text" :key="selectedFile.id">{{ selectedFile.text }}</pre>
                </template>
                <!-- no preview -->
                <template v-else>
                  <div class="no-preview">Aperçu non disponible pour ce fichier</div>
                </template>
              </div>
              <div class="info-panel">
                <p><strong>Nom :</strong> {{ selectedFile.name }}</p>
                <p v-if="selectedFile.size"><strong>Taille :</strong> {{ selectedFile.sizeLabel || humanSize(selectedFile.size) }}</p>
                <p><strong>Type :</strong> {{ selectedFile.mime || selectedFile.type }}</p>
                <p><strong>Propriétaire :</strong> {{ selectedFile.owner }}</p>
                <p><strong>Emplacement :</strong> {{ selectedFile.parentPath || '/' }}</p>
              </div>
            </div>
          </div>
        </section>
      </section>
    </main>
    <!-- hidden file input -->
    <input ref="uploadInput" type="file" multiple style="display:none" @change="handleUpload" />
    <!-- UPLOAD MODAL -->
    <div v-if="showUploadModal" class="modal-wrap" @keydown.esc="showUploadModal = false">
      <div class="modal" role="dialog" aria-modal="true" aria-label="Uploader des fichiers">
        <h3>Uploader des fichiers</h3>
        <p>Glisse-dépose ou clique pour sélectionner des fichiers.</p>
        <div class="upload-zone" @click="triggerFilePicker" @dragover.prevent @drop.prevent="onDrop">
          <p>Déposez vos fichiers ici — ou cliquez</p>
        </div>
        <div class="modal-actions">
          <button class="dv-btn" @click="showUploadModal = false">Annuler</button>
        </div>
      </div>
      <div class="backdrop" @click="showUploadModal = false"></div>
    </div>
    <!-- NEW FOLDER MODAL -->
    <div v-if="showNewFolderModal" class="modal-wrap" @keydown.esc="showNewFolderModal = false">
      <div class="modal" role="dialog" aria-modal="true" aria-label="Nouveau dossier">
        <h3>Nouveau dossier</h3>
        <input v-model="newFolderName" placeholder="Nom du dossier" class="compose-input" @keyup.enter="createFolder" />
        <div class="modal-actions" style="justify-content:space-between;">
          <div>
            <button class="dv-btn ghost" @click="showNewFolderModal = false">Annuler</button>
          </div>
          <div style="display:flex; gap:8px;">
            <button class="dv-btn" @click="createFolder">Créer</button>
          </div>
        </div>
      </div>
      <div class="backdrop" @click="showNewFolderModal = false"></div>
    </div>
    <!-- RENAME MODAL -->
    <div v-if="showRenameModal" class="modal-wrap" @keydown.esc="showRenameModal = false">
      <div class="modal" role="dialog" aria-modal="true" aria-label="Renommer">
        <h3>Renommer</h3>
        <input v-model="renameValue" @keyup.enter="applyRename" />
        <div class="modal-actions">
          <button class="dv-btn" @click="showRenameModal = false">Annuler</button>
          <button class="dv-btn primary" @click="applyRename">Valider</button>
        </div>
      </div>
      <div class="backdrop" @click="showRenameModal = false"></div>
    </div>
    <!-- SHARE MODAL -->
    <div v-if="showShareModal" class="modal-wrap" @keydown.esc="showShareModal = false">
      <div class="modal" role="dialog" aria-modal="true" aria-label="Partager">
        <h3>Lien de partage</h3>
        <input readonly :value="shareLink" @focus="$event.target.select()" />
        <div class="modal-actions">
          <button class="dv-btn" @click="copyShareLink">Copier</button>
          <button class="dv-btn" @click="showShareModal = false">Fermer</button>
        </div>
      </div>
      <div class="backdrop" @click="showShareModal = false"></div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { gls } from '@/stores/global'
import router from '@/router/index'
// simple auth redirect
if (gls().log != 1) {
  router.push('login')
}
/* état */
const curr = ref('drive') // drive | shared | trash
const currentPath = ref('/')
const userName = ref(gls().username || 'Alice')
const userEmail = ref((gls().username || 'alice') + '@example.com')
const isMobile = ref(window.innerWidth <= 750)
const searchQuery = ref('')
const selectedFile = ref(null)
const uploadInput = ref(null)
const dragOver = ref(false)
/* modals */
const showUploadModal = ref(false)
const showNewFolderModal = ref(false)
const showRenameModal = ref(false)
const showShareModal = ref(false)
const renameValue = ref('')
const newFolderName = ref('')
const shareLink = ref('')
/* helper taille */
function humanSize(bytes) {
  if (bytes == null) return ''
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}
/* sample image for demo */
const sampleImgDataUrl = `data:image/svg+xml;utf8,${encodeURIComponent(
  `<svg xmlns='http://www.w3.org/2000/svg' width='600' height='400'><rect width='100%' height='100%' fill='#f6fbff'/><text x='50%' y='50%' font-size='32' text-anchor='middle' fill='#6aa8ff' dy='.35em'>Image de démo</text></svg>`
)}`
/* données initiales */
const files = ref([
  { id: 1, name: 'Documents', type: 'folder', parentPath: '/', date: new Date().toLocaleString(), folder: 'drive' },
  { id: 2, name: 'Photos', type: 'folder', parentPath: '/', date: new Date().toLocaleString(), folder: 'drive' },
  { id: 3, name: 'projet.pdf', type: 'file', mime: 'application/pdf', size: 234000, sizeLabel: humanSize(234000), parentPath: '/Documents/', date: new Date().toLocaleString(), folder: 'drive' },
  { id: 4, name: 'vacances.mp4', type: 'file', mime: 'video/mp4', size: 1200000, sizeLabel: humanSize(1200000), parentPath: '/Photos/', date: new Date().toLocaleString(), folder: 'drive', url: sampleImgDataUrl },
  { id: 99, name: 'ancien.txt', type: 'file', mime: 'text/plain', size: 220, sizeLabel: humanSize(220), parentPath: '/', date: new Date().toLocaleString(), folder: 'trash', text: 'Contenu supprimé (exemple).' }
])
const currLabel = computed(() => (curr.value === 'drive' ? 'Mes fichiers' : curr.value === 'shared' ? 'Partagés' : 'Corbeille'))
/* filtered files selon dossier courant + recherche */
const filteredFiles = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  let list
  if (curr.value === 'trash') list = files.value.filter(f => f.folder === 'trash')
  else if (curr.value === 'shared') list = files.value.filter(f => f.folder === 'shared' || f.folder === 'drive')
  else list = files.value.filter(f => f.folder === 'drive' && f.parentPath === currentPath.value)
  if (!q) return list
  return list.filter(f => f.name.toLowerCase().includes(q))
})
/* navigation */
function changec(folder) {
  curr.value = folder
  currentPath.value = '/'
  selectedFile.value = null
}
function openFolder(folder) {
  if (folder.type !== 'folder') return
  currentPath.value = folder.parentPath + folder.name + '/'
  selectedFile.value = null
}
/* upload */
function openUploadModal() { showUploadModal.value = true }
function triggerFilePicker() { if (uploadInput.value) uploadInput.value.click() }
function handleUpload(e) {
  const fileList = e.target?.files || e.dataTransfer?.files
  if (!fileList || fileList.length === 0) return
  Array.from(fileList).forEach(f => {
    const id = Date.now() + Math.floor(Math.random() * 1000)
    const item = {
      id,
      name: f.name,
      type: 'file',
      mime: f.type || 'application/octet-stream',
      size: f.size,
      sizeLabel: humanSize(f.size),
      parentPath: currentPath.value,
      date: new Date().toLocaleString(),
      folder: 'drive',
      owner: userEmail.value
    }
    // create preview URLs / read text
    if (f.type && f.type.startsWith('image/')) {
      item.url = URL.createObjectURL(f)
    } else if (f.type && f.type.startsWith('video/')) {
      item.url = URL.createObjectURL(f)
    } else if (f.type && f.type.startsWith('audio/')) {
      item.url = URL.createObjectURL(f)
    } else if (f.type && f.type.startsWith('text/')) {
      const reader = new FileReader()
      reader.onload = ev => { item.text = ev.target.result }
      reader.readAsText(f.slice(0, 200000))
    }
    files.value.push(item)
  })
  showUploadModal.value = false
  dragOver.value = false
  if (uploadInput.value) uploadInput.value.value = ''
}
/* drag & drop */
function onDrop(e) {
  dragOver.value = false
  handleUpload(e)
}
/* open / preview / download */
function openFile(it) { selectedFile.value = it }
function openAction(it) {
  if (!it) return
  if (it.url) { window.open(it.url, '_blank'); return }
  if (it.text) {
    const w = window.open('', '_blank')
    w.document.write(`<pre style="white-space:pre-wrap;font-family:monospace">${escapeHtml(it.text)}</pre>`)
    w.document.title = it.name
    return
  }
  const blob = new Blob([`Simulation d'ouverture : ${it.name}`], { type: it.mime || 'text/plain' })
  const url = URL.createObjectURL(blob)
  window.open(url, '_blank')
  setTimeout(() => URL.revokeObjectURL(url), 2000)
}
function downloadFile(it) {
  if (!it) return
  if (it.url) {
    const a = document.createElement('a'); a.href = it.url; a.download = it.name; document.body.appendChild(a); a.click(); a.remove(); return
  }
  const blob = new Blob([`Téléchargement simulé : ${it.name}`], { type: it.mime || 'application/octet-stream' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a'); a.href = url; a.download = it.name; document.body.appendChild(a); a.click(); a.remove()
  URL.revokeObjectURL(url)
}
/* corbeille / restaurer / suppression définitive */
function moveToTrash(ids = null) {
  const idsArr = Array.isArray(ids) ? ids : (selectedFile.value ? [selectedFile.value.id] : [])
  idsArr.forEach(id => {
    const idx = files.value.findIndex(f => f.id === id)
    if (idx !== -1) files.value[idx].folder = 'trash'
  })
  // if the currently selected file was moved to trash, we keep it selected (but curr may differ), clear selection if needed
  if (selectedFile.value && idsArr.includes(selectedFile.value.id)) selectedFile.value = files.value.find(f => f.id === selectedFile.value.id) || null
}
function restoreFile(it) {
  if (!it) return
  const idx = files.value.findIndex(f => f.id === it.id)
  if (idx !== -1) files.value[idx].folder = 'drive'
  selectedFile.value = null
}
function deletePermanently(ids = null) {
  const idsArr = Array.isArray(ids) ? ids : (selectedFile.value ? [selectedFile.value.id] : [])
  files.value = files.value.filter(f => !idsArr.includes(f.id))
  selectedFile.value = null
}
function confirmPermanentDelete(file) {
  if (!file) return
  const ok = window.confirm(`Supprimer définitivement "${file.name}" ? Cette action est irréversible.`)
  if (!ok) return
  deletePermanently([file.id])
}
/* rename */
function startRename(it) {
  renameValue.value = it.name
  selectedFile.value = it
  showRenameModal.value = true
}
function applyRename() {
  if (!selectedFile.value) { showRenameModal.value = false; return }
  const idx = files.value.findIndex(f => f.id === selectedFile.value.id)
  if (idx !== -1) files.value[idx].name = renameValue.value || files.value[idx].name
  showRenameModal.value = false
  renameValue.value = ''
}
/* share */
function shareFile(it) {
  shareLink.value = `${window.location.origin}/share/${it.id}`
  showShareModal.value = true
}
function copyShareLink() { if (!shareLink.value) return; navigator.clipboard?.writeText(shareLink.value) }
/* selection simple */
function selectFirst() { if (filteredFiles.value.length) openFile(filteredFiles.value[0]) }
function clearSelection() { selectedFile.value = null }
/* create folder (modal) */
function createFolder() {
  const name = (newFolderName.value || 'Nouveau dossier').trim()
  if (!name) return
  const id = Date.now() + Math.floor(Math.random() * 1000)
  files.value.push({ id, name, type: 'folder', parentPath: currentPath.value, date: new Date().toLocaleString(), folder: 'drive' })
  newFolderName.value = ''
  showNewFolderModal.value = false
}
/* small util */
function escapeHtml(str = '') { return String(str).replace(/[&<>"']/g, s => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[s])) }
onMounted(() => {
  window.addEventListener('resize', () => { isMobile.value = window.innerWidth <= 750 })
})
</script>
<style scoped>
#drive-v2 {
  display: flex;
  height: 100vh;
  font-family: "Roboto", system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial;
  background: #f6f7fb;
  color: #111;
}
/* DARK MODE */
.dark #drive-v2 { background: #121212; color: #e5e5e5; }
/* SIDEBAR */
.dv-sidebar {
  width: 240px;
  padding: 16px;
  background: #fff;
  border-right: 1px solid #e6e6ee;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.dark .dv-sidebar { background: #161616; border-color: #2b2b2b; }
.dv-user { display:flex; align-items:center; gap:12px; }
.dv-avatar { width:44px; height:44px; border-radius:8px; display:flex; align-items:center; justify-content:center; font-weight:700; background: -webkit-linear-gradient(30deg, blue, red); color:#fff; }
.dv-name { font-weight:600; }
.dv-email { font-size:0.85rem; color:#666; }
.dark .dv-email { color:#bbb; }
.dv-actions { display:flex; gap:8px; flex-direction:column; }
.dv-btn { padding:10px 12px; border-radius:8px; border:none; cursor:pointer; background:transparent; text-align:left; }
.dv-btn.primary {  background: -webkit-linear-gradient(30deg, blue, red); color:#fff; }
.dv-btn.ghost { border:1px solid #eee; background:#fff; }
.dv-btn.small { padding:6px 8px; font-size:0.9rem; }
.dv-nav { display:flex; flex-direction:column; gap:8px; margin-top:6px; }
.nav-item { background:transparent; border:none; padding:10px; text-align:left; border-radius:8px; cursor:pointer; color:#333; }
.nav-item.active {  background: -webkit-linear-gradient(30deg, blue, red);; color:#fff; }
/* MAIN */
.dv-main { flex:1; padding:16px; display:flex; flex-direction:column; gap:12px; }
.dv-header { display:flex; justify-content:space-between; align-items:center; }
.dv-breadcrumbs { display:flex; align-items:center; gap:8px; }
.crumb { background:transparent; border:none; color:#0066cc; cursor:pointer; }
.current { font-weight:600; }
.dv-tools { display:flex; gap:8px; align-items:center; }
.dv-search { padding:8px 10px; border-radius:8px; border:1px solid #e0e0e6; }
/* BODY SPLIT */
.dv-body { display:grid; grid-template-columns: 420px 1fr; gap:20px; height: calc(100% - 60px); }
/* LIST */
.dv-list { background:#fff; border-radius:8px; padding:12px; overflow:auto; border:1px solid #eaeaf2; min-height:200px; transition: border-color .12s, background .12s; }
.dark .dv-list { background:#151515; border-color:#2b2b2b; }
.dv-list.drag-over { background:#f3fff6; border-color:#bff0c5; }
.dv-empty { display:flex; flex-direction:column; align-items:center; justify-content:center; height:100%; padding:24px; text-align:center; color:#666; }
.dv-empty-illustration { margin-bottom:10px; }
.dv-empty-text { font-size:1.0rem; margin:6px 0; }
.dv-empty-sub { font-size:0.9rem; color:#999; }
.dv-items { list-style:none; margin:0; padding:0; display:flex; flex-direction:column; gap:8px; }
.dv-item { display:flex; justify-content:space-between; align-items:center; padding:10px; border-radius:8px; cursor:pointer; transition: background .12s; }
.dv-item:hover { background:#fbfbfe; }
.dv-item.active { background:#e9f1ff; }
.dark .dv-item:hover { background:#1b1b1b; }
.dark .dv-item.active { background:#2b2b35; }
.dv-item .left { display:flex; gap:12px; align-items:center; }
.icon { width:36px; font-size:1.1rem; text-align:center; }
.meta .name { font-weight:600; }
.meta .sub { font-size:0.85rem; color:#777; }
.dark .meta .sub { color:#bbb; }
.dv-item .actions { display:flex; gap:6px; }
.dv-mini { background:transparent; border:none; padding:6px; cursor:pointer; border-radius:6px; }
.dv-mini:hover { background:#f0f0f7; }
/* DETAIL */
.dv-detail { background:#fff; border-radius:8px; padding:18px; border:1px solid #eaeaf2; display:flex; flex-direction:column; gap:12px; min-height:200px; }
.dark .dv-detail { background:#151515; border-color:#2b2b2b; color:#e5e5e5; }
.dv-detail-empty { display:flex; flex-direction:column; align-items:center; justify-content:center; text-align:center; color:#666; padding:24px; }
.dv-empty-illustration-center { margin-bottom:10px; }
.dv-empty-text-centered { font-size:1.0rem; }
.detail-top { display:flex; justify-content:space-between; align-items:flex-start; gap:12px; }
.detail-info { display:flex; flex-direction:column; gap:6px; }
.detail-name { margin:0; font-size:1.15rem; }
.detail-meta { color:#666; font-size:0.95rem; }
.dark .detail-meta { color:#bbb; }
.detail-actions { display:flex; gap:8px; flex-wrap:wrap; align-items:center; }
.dv-btn.green { background:#e6ffe6; color:#008000; }
.dv-btn.purple { background:#f3e8ff; color:#6a00b6; }
.dv-btn.danger { background:#ffecec; color:#a00; }
/* detail body */
.detail-body { display:flex; gap:12px; margin-top:12px; align-items:flex-start; }
.preview { flex:1; display:flex; align-items:center; justify-content:center; min-height:140px; border-radius:8px; background:#fbfbfe; border:1px dashed #eee; padding:12px; }
.preview-img { max-width:100%; max-height:280px; border-radius:8px; }
.preview-video { max-width:100%; max-height:280px; border-radius:8px; background:#000; }
.preview-audio { width:100%; }
.preview-text { width:100%; max-height:280px; overflow:auto; background:#fff; padding:12px; border-radius:8px; font-family:monospace; }
.no-preview { color:#777; }
.info-panel { width:220px; background:#fff; padding-left:12px; color:#444; font-size:0.95rem; }
.dark .info-panel { color:white; background:none; }
/* MODALS (centered) */
.modal-wrap { position:fixed; inset:0; display:flex; align-items:center; justify-content:center; z-index:2000; }
.modal { width:420px; background:#fff; padding:20px; border-radius:12px; box-shadow:0 10px 40px rgba(0,0,0,0.12); z-index:2001; }
.dark .modal { background:#1e1e1e; color:#e5e5e5; }
.backdrop { position:fixed; inset:0; background:rgba(0,0,0,0.25); z-index:2000; }
.upload-zone { border:2px dashed #e6e6ee; border-radius:8px; padding:18px; text-align:center; cursor:pointer; margin:10px 0; }
.modal-actions { display:flex; gap:8px; justify-content:flex-end; margin-top:12px; }
.modal input, .modal textarea, .modal .compose-input { width:100%; padding:10px; border-radius:8px; border:1px solid #e6e6ee; }
.dark .modal input, .dark .modal .compose-input { background:#121212; border:1px solid #333; color:#eee; }
/* responsive */
@media (max-width: 900px) {
  .dv-body { grid-template-columns: 1fr; }
  .dv-sidebar { display:none; }
  .info-panel { width:100%; border-left:none; padding-left:0; }
}
</style>
