<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { gls } from '@/stores/global.js'
import { SERVICES_CONFIG } from '../config/services'


// ---------- helpers pour résolution d'URL API ----------
function resolveAPI(pathOrUrl) {
  const base = (import.meta.env.VITE_APP_API || '').replace(/\/+$/, '')
  if (!pathOrUrl) return base || ''
  if (pathOrUrl.startsWith('http://') || pathOrUrl.startsWith('https://')) return pathOrUrl
  const path = pathOrUrl.startsWith('/') ? pathOrUrl : '/' + pathOrUrl
  if (!base) return path
  return base + path
}

// Lecture des variables d'env (préférer valeurs spécifiques si fournies)
const API_GETFILES = resolveAPI(import.meta.env.VITE_API_DRIVE || '/api/drive/getfiles')
const API_UPLOAD = resolveAPI(import.meta.env.VITE_API_DRIVE_UPLOAD || '/api/drive/upload')
const API_DOWNLOAD = resolveAPI(import.meta.env.VITE_API_DRIVE_DOWNLOAD || '/api/drive/download')
const API_RENAME = resolveAPI(import.meta.env.VITE_API_DRIVE_RENAME || '/api/drive/rename')
const API_DELETE = resolveAPI(import.meta.env.VITE_API_DRIVE_DELETE || '/api/drive/delete')
const API_TRASH = resolveAPI(import.meta.env.VITE_API_DRIVE_TRASH || '/api/drive/trash')
const API_GTRASH = resolveAPI(import.meta.env.VITE_API_DRIVE_GTRASH || '/api/drive/gettrash')
const API_RESTORE = resolveAPI(import.meta.env.VITE_API_DRIVE_RESTORE || '/api/drive/restore')
const API_CREATE_FOLDER = resolveAPI(import.meta.env.VITE_API_DRIVE_CREATE_FOLDER || '/api/drive/createFolder')
const API_CREATE_FILE = resolveAPI(import.meta.env.VITE_API_DRIVE_CREATE_FILE || '/api/drive/createFile')
const API_MOVE_FILE = resolveAPI(import.meta.env.VITE_API_DRIVE_MOVE_FILE || '/api/drive/moveFile')
const API_MOVE_FOLDER = resolveAPI(import.meta.env.VITE_API_DRIVE_MOVE_FOLDER || '/api/drive/moveFolder')
const API_SHARE = resolveAPI(import.meta.env.VITE_API_DRIVE_SHARE || '/api/drive/share')
const API_UNSHARE = resolveAPI(import.meta.env.VITE_API_DRIVE_UNSHARE || '/api/drive/unshare')
const API_GET_SHARES = resolveAPI(import.meta.env.VITE_API_DRIVE_GET_SHARES || '/api/drive/shares')
const API_STORAGE_INFO = resolveAPI(import.meta.env.VITE_API_DRIVE_STORAGE_INFO || '/api/drive/storage-info')

// ---------- état / UI ----------
const userName = gls().username
const userEmail = ref('')
const isMobile = computed(() => window.innerWidth <= 900)

// Fonction de sélection adaptée mobile
function selectItem(it) {
  if (it.type === 'folder') {
    openFolder(it)
  } else {
    openFile(it)
  }
}
const searchQuery = ref('')
const selectedFile = ref(null)

// ---------- Storage info ----------
const storageUsed = ref(0)
const storageLimit = ref(0)
const isUnlimitedStorage = ref(false)
const storagePercentage = computed(() => {
  if (isUnlimitedStorage.value) return 0
  if (storageLimit.value === 0 || storageLimit.value === null) return 0
  const pct = (storageUsed.value / storageLimit.value) * 100
  return Math.min(100, Math.round(pct))
})
const storageColor = computed(() => {
  const pct = storagePercentage.value
  if (pct >= 90) return '#404040' // rouge
  if (pct >= 75) return '#a7a7a7' // orange
  return '#808080' // vert
})

const curr = ref('drive') // drive | shared | trash
const currentPath = ref('/') // toujours finie par '/'
const currentParts = computed(() => {
  const p = normalizePath(currentPath.value || '/')
  if (p === '/') return []
  return p.replace(/^\//, '').replace(/\/$/, '').split('/')
})

const files = ref([])

// moving lock
const movingInProgress = ref(false)

function changec(folder) {
  curr.value = folder
  currentPath.value = folder === 'trash' ? '/.trash/' : '/'
  selectedFile.value = null
}

function openFolder(folder) {
  if (!folder || folder.type !== 'folder') return
  currentPath.value = (folder.parentPath || '/') + folder.name + '/'
  selectedFile.value = null
}

// drag/drop basic state
const dragOver = ref(false)
const draggedItem = ref(null)
const dragTarget = ref(null)

// upload UI state & queue
const uploadInput = ref(null)
const showUploadModal = ref(false)
const showMobileFab = ref(false)
const showNewFolderModal = ref(false)
const showNewFileModal = ref(false)
const showRenameModal = ref(false)
const showShareModal = ref(false)
const renameValue = ref('')
const newFolderName = ref('')
const newFileName = ref('')
const newFileType = ref('docx')
const shareLink = ref('')
const showOnlyofficeModal = ref(false) // montre le modal OnlyOffice
// share modal state (was missing -> runtime error when clicking "Partager")
const shareWithUsername = ref('')
const sharePermission = ref('editor')
const currentShares = ref([])
// upload queue: item = { id, file, name, size, progress, status, xhr }
const uploadQueue = ref([])

// overall progress computed (weighted by file size)
const overallProgress = computed(() => {
  if (uploadQueue.value.length === 0) return 0
  let totalSize = 0
  let weighted = 0
  uploadQueue.value.forEach(u => {
    const s = u.size || 0
    totalSize += s
    weighted += (u.progress || 0) * s
  })
  if (totalSize === 0) {
    const avg = Math.round(uploadQueue.value.reduce((a,b)=>a+(b.progress||0),0) / uploadQueue.value.length)
    return avg
  }
  return Math.round(weighted / totalSize)
})

const uploadsSummary = computed(() => {
  const total = uploadQueue.value.length
  const done = uploadQueue.value.filter(u => u.status === 'done').length
  const err = uploadQueue.value.filter(u => u.status === 'error').length
  const uploading = uploadQueue.value.filter(u => u.status === 'uploading').length
  return `${done}/${total} • ${uploading} en cours • ${err} erreurs`
})

// ---------- util ----------
function humanSize(bytes) {
  if (bytes == null) return ''
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`
}

function escapeHtml(str) { return String(str).replace(/[&<>"']/g, s => ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[s])) }

function mapServerFile(row) {
  const id = row.file_id ?? row.fileId ?? row.id
  const fileName = row.file_name ?? row.fileName ?? ''
  let p = row.file_path ?? row.filePath ?? '/'
  if (!p) p = '/'
  if (p !== '/' && p.slice(-1) !== '/') p = p + '/'
  const parentPath = p
  const size = (row.file_size !== null && row.file_size !== undefined) ? Number(row.file_size) : null
  const mime = row.file_type ?? row.fileType ?? null
  const date = row.date_uploaded ?? row.dateUploaded ?? row.date ?? ''
  const lowerType = (mime || '').toString().toLowerCase()
  const isFolder = lowerType.includes('folder') || lowerType.includes('directory') || (row.is_folder === true)
  const type = isFolder ? 'folder' : 'file'
  let url = null
  if (!isFolder) {
    const token = gls().sessionT || ''
    const uname = encodeURIComponent(userName || '')
    url = `${API_DOWNLOAD}?file_id=${encodeURIComponent(id)}&token=${encodeURIComponent(token)}&username=${uname}`
  }

  let folderFlag = 'drive'
  if (String(parentPath).includes('/.trash/') || parentPath === '/.trash/') folderFlag = 'trash'
  else if (row.shared === true || row.is_shared === true) folderFlag = 'shared'

  return {
    id,
    name: fileName,
    type,
    mime,
    size,
    sizeLabel: size ? humanSize(size) : null,
    parentPath,
    date,
    folder: folderFlag,
    owner: userName || row.owner || '',
    url,
    _server: row
  }
}

var sMul = ref(false)
var sMulAll = ref(false)
var selectedC = ref(new Array(files.value.length).fill(false))

function selectAll() {
  sMulAll.value = !sMulAll.value
  for (var i = 0; i < selectedC.value.length; i++) {
    selectedC.value[i] = sMulAll.value
  }
}

function selectMultiple() {
  selectedC.value = new Array(files.value.length).fill(false)
  sMul.value = !sMul.value
}

// ---------- appels API ----------
async function fetchFiles() {
  try {
    const payload = { username: userName, token: gls().sessionT }
    let api
    if (curr.value === 'shared') {
      api = API_GET_SHARES
    } else if (curr.value === 'trash') {
      api = API_GTRASH
    } else {
      api = API_GETFILES
    }
    const res = await fetch(api, {
      method: 'POST',
      mode: 'cors',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    if (!res.ok) {
      console.warn('getfiles failed', res.status, res.statusText)
      files.value = []
      return
    }
    const j = await res.json()
    const rows = Array.isArray(j.files) ? j.files : (Array.isArray(j) ? j : (j.data || []))
    files.value = rows.map(mapServerFile)
    // reset selectedC length to current files length
    selectedC.value = new Array(files.value.length).fill(false)
    if (j.user_email) userEmail.value = j.user_email
  } catch (err) {
    console.error('fetchFiles error', err)
  }
}

// Fetch storage info
async function fetchStorageInfo() {
  try {
    const res = await fetch(API_STORAGE_INFO, {
      method: 'POST',
      mode: 'cors',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username: userName, token: gls().sessionT }),
    })
    if (!res.ok) {
      console.warn('storage-info failed', res.status)
      return
    }
    const data = await res.json()
    storageUsed.value = data.current_usage || 0
    storageLimit.value = data.storage_limit || 0
    isUnlimitedStorage.value = data.unlimited || false
    // Handle unlimited storage (Enterprise)
    if (data.unlimited) {
      storageLimit.value = storageUsed.value * 2 || 1024 * 1024 * 1024 // Show bar at 50% if unlimited
    }
  } catch (err) {
    console.error('fetchStorageInfo error', err)
  }
}

watch(curr, async () => {
  await fetchFiles()
  await fetchStorageInfo()
})

// normalise un chemin pour qu'il soit '/' ou commence et termine par '/'
function normalizePath(p) {
  if (!p) return '/'
  let s = String(p)
  s = s.replace(/\/+/g, '/') // collapse multiple slashes
  if (s === '/') return '/'
  if (!s.startsWith('/')) s = '/' + s
  if (!s.endsWith('/')) s = s + '/'
  return s
}

// navigue vers la racine
function goRoot() {
  currentPath.value = '/'
  selectedFile.value = null
  dragTarget.value = null
}

// navigue vers le crumb d'index i (0-based)
function goToBreadcrumb(i) {
  const parts = currentParts.value.slice(0, i + 1)
  currentPath.value = '/' + (parts.length ? parts.join('/') + '/' : '')
  currentPath.value = normalizePath(currentPath.value)
  selectedFile.value = null
  dragTarget.value = null
  fetchFiles().catch(()=>{})
}

// remonte d'un niveau
function goUp() {
  if (!currentPath.value || currentPath.value === '/') return
  const parts = currentParts.value.slice()
  parts.pop()
  currentPath.value = '/' + (parts.length ? parts.join('/') + '/' : '')
  currentPath.value = normalizePath(currentPath.value)
  selectedFile.value = null
  dragTarget.value = null
  fetchFiles().catch(()=>{})
}

async function moveSelectedToTrash() {
  const selectedIds = files.value
    .filter((f, i) => selectedC.value[i])
    .map(f => f.id)

  if (selectedIds.length === 0) {
    window.alert("Aucun élément sélectionné")
    return
  }

  const confirmMsg = `Mettre ${selectedIds.length} élément(s) à la corbeille ?`
  if (!window.confirm(confirmMsg)) return

  for (const id of selectedIds) {
    await moveToTrash(id)
  }

  sMulAll.value = false
  selectedC.value = new Array(files.value.length).fill(false)
  await fetchFiles()
}

// createFolder, upload logic, rename/delete, restore, moveToTrash remain unchanged
// (I keep your existing implementations for those — include them below as in your original code)
// For brevity I will re-use your implementations: createFolder, upload functions, renameFileRequest, deleteFileRequest, restoreFile, moveToTrash, etc.
// --- paste the unmodified functions from your original script for createFolder/upload/rename/delete/restore/trash ---
// (IN PRACTICE: keep the implementations you already have; they were included earlier in your file)

async function createFolder() {
  const name = (newFolderName.value || '').trim()
  if (!name) {
    window.alert('Le nom du dossier est vide')
    return
  }

  const payload = {
    username: userName,
    token: gls().sessionT,
    parent_path: currentPath.value || '/',
    folder_name: name
  }

  try {
    const res = await fetch(API_CREATE_FOLDER, {
      method: 'POST',
      mode: 'cors',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    if (res.status === 401) {
      window.alert('Session invalide')
      return
    }
    if (res.status === 409) {
      const j = await res.json().catch(()=>({}))
      window.alert(j.error || 'Un fichier/dossier du même nom existe déjà')
      return
    }
    if (!res.ok) {
      const txt = await res.text().catch(()=> '')
      console.error('createFolder failed', res.status, txt)
      window.alert('Erreur lors de la création du dossier')
      return
    }

    await fetchFiles()
    showNewFolderModal.value = false
    newFolderName.value = ''
  } catch (err) {
    console.error('createFolder error', err)
    window.alert('Impossible de créer le dossier')
  }
}

async function createFile() {
  const name = (newFileName.value || '').trim()
  if (!name) {
    window.alert('Le nom du fichier est vide')
    return
  }

  const payload = {
    username: userName,
    token: gls().sessionT,
    parent_path: currentPath.value || '/',
    file_name: name,
    file_type: newFileType.value
  }

  try {
    const res = await fetch(API_CREATE_FILE, {
      method: 'POST',
      mode: 'cors',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })

    if (res.status === 401) {
      window.alert('Session invalide')
      return
    }
    if (res.status === 409) {
      const j = await res.json().catch(()=>({}))
      window.alert(j.error || 'Un fichier du même nom existe déjà')
      return
    }
    if (!res.ok) {
      const txt = await res.text().catch(()=> '')
      console.error('createFile failed', res.status, txt)
      window.alert('Erreur lors de la création du fichier')
      return
    }

    const result = await res.json()
    await fetchFiles()
    
    // Ouvrir automatiquement le fichier dans OnlyOffice
    const newFile = files.value.find(f => f.id === result.file_id)
    if (newFile) {
      selectedFile.value = newFile
      showOnlyofficeModal.value = true
      await openInOnlyOffice(newFile)
    }
    
    showNewFileModal.value = false
    newFileName.value = ''
  } catch (err) {
    console.error('createFile error', err)
    window.alert('Impossible de créer le fichier')
  }
}

// Upload + progress functions (unchanged) - copy from original
function openUploadModal() { showUploadModal.value = true }
function triggerFilePicker() { if (uploadInput.value) uploadInput.value.click() }
function handleUpload(e) {
  let fileList = null
  if (e && e.target && e.target.files) fileList = e.target.files
  else if (e && e.dataTransfer && e.dataTransfer.files) fileList = e.dataTransfer.files
  if (!fileList || fileList.length === 0) return
  showUploadModal.value = true
  uploadFilesFromInput(fileList)
  if (uploadInput.value) uploadInput.value.value = ''
}
function uploadFilesFromInput(fileList) {
  const arr = Array.from(fileList)
  arr.forEach((f, idx) => {
    const id = `${Date.now()}_${idx}_${Math.random().toString(36).slice(2,8)}`
    const item = { id, file: f, name: f.name, size: f.size, progress: 0, status: 'queued', xhr: null }
    uploadQueue.value.push(item)
    setTimeout(() => startUpload(item), 10 * idx)
  })
}
function startUpload(item) {
  if (!item || !item.file) return
  item.status = 'uploading'
  item.progress = 0

  const fd = new FormData()
  fd.append('username', userName)
  if (gls().sessionT) fd.append('token', gls().sessionT)
  fd.append('parent_path', currentPath.value || '/')
  fd.append('files', item.file, item.file.name)

  const xhr = new XMLHttpRequest()
  item.xhr = xhr
  xhr.open('POST', API_UPLOAD, true)

  xhr.upload.onprogress = function (ev) {
    if (!ev.lengthComputable) return
    const percent = Math.round((ev.loaded / ev.total) * 100)
    item.progress = percent
  }

  xhr.onload = async function () {
    item.xhr = null
    if (xhr.status >= 200 && xhr.status < 300) {
      try {
        item.progress = 100
        item.status = 'done'
        await fetchFiles()
        await fetchStorageInfo()
      } catch (e) {
        item.status = 'done'
        await fetchFiles()
        await fetchStorageInfo()
      }
    } else {
      item.status = 'error'
      console.error('upload failed', xhr.status, xhr.responseText)
    }
  }
  xhr.onerror = function () { item.status = 'error'; item.xhr = null }
  xhr.onabort = function () { item.status = 'cancelled'; item.xhr = null }
  xhr.send(fd)
}
function cancelUpload(id) {
  const idx = uploadQueue.value.findIndex(u => u.id === id)
  if (idx === -1) return
  const u = uploadQueue.value[idx]
  if (u.xhr) {
    try { u.xhr.abort() } catch (e) {}
  } else { u.status = 'cancelled' }
}
function retryUpload(id) {
  const idx = uploadQueue.value.findIndex(u => u.id === id)
  if (idx === -1) return
  const u = uploadQueue.value[idx]
  if (u.status === 'uploading') return
  u.progress = 0
  u.status = 'queued'
  setTimeout(() => startUpload(u), 50)
}
function cancelAllUploads() {
  uploadQueue.value.forEach(u => {
    if (u.xhr) try { u.xhr.abort() } catch(e) {}
    else u.status = 'cancelled'
  })
}
watch(uploadQueue, (q) => {
  if (!q || q.length === 0) return
  const allFinished = q.every(u => ['done','error','cancelled'].includes(u.status))
  if (allFinished) {
    setTimeout(() => {
      showUploadModal.value = false
      uploadQueue.value = uploadQueue.value.filter(u => u.status !== 'done')
    }, 700)
  }
}, { deep: true })

// ---------- rename / delete via API ----------
async function renameFileRequest(fileId, newName) {
  try {
    const payload = { username: userName, token: gls().sessionT, file_id: fileId, new_name: newName }
    const res = await fetch(API_RENAME, {
      method: 'POST',
      mode: 'cors',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(payload)
    })
    if (!res.ok) throw new Error('rename failed')
    await fetchFiles()
  } catch (err) { console.error(err) }
}

async function deleteFileRequest(fileId) {
  try {
    const payload = { username: userName, token: gls().sessionT, file_id: fileId }
    const res = await fetch(API_DELETE, {
      method: 'POST',
      mode: 'cors',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(payload)
    })
    if (!res.ok) throw new Error('delete failed')
    await fetchFiles()
  } catch (err) { console.error(err) }
}

async function deleteSelectedPermanently() {
  const selectedIds = files.value
    .filter((f, i) => selectedC.value[i])
    .map(f => f.id)

  if (selectedIds.length === 0) {
    window.alert("Aucun élément sélectionné")
    return
  }

  const confirmMsg = `Supprimer définitivement ${selectedIds.length} élément(s) ? Cette action est irréversible.`
  if (!window.confirm(confirmMsg)) return

  for (const id of selectedIds) {
    await deleteFileRequest(id)
  }

  sMulAll.value = false
  selectedC.value = new Array(files.value.length).fill(false)
  await fetchFiles()
}

// ---------- UI actions (move/rename/download/etc.) ----------
function openFile(it) { selectedFile.value = it }

function openAction(it) {
  if (!it) return
  
  // Ajouter cette condition pour les fichiers Office
  if (['docx','xlsx','pptx'].includes(it.name.split('.').pop().toLowerCase())) {
    showOnlyofficeModal.value = true
    openInOnlyOffice(it)
    return
  }

  // Reste du code existant...
  if (it.url) {
    window.open(it.url, '_blank')
    return
  }
  if (it.text) {
    const w = window.open('', '_blank')
    w.document.write(`<pre style="white-space:pre-wrap;font-family:monospace">${escapeHtml(it.text)}</pre>`)
    w.document.title = it.name
    return
  }
  window.alert('Aperçu non disponible pour ce fichier. Vérifiez que le backend fournit une URL de téléchargement ou le contenu.')
}

function downloadFile(it) {
  if (!it) return
  if (it.url) {
    const sep = it.url.includes('?') ? '&' : '?'
    const downloadUrl = it.url + sep + 'download=1'
    const a = document.createElement('a')
    a.href = downloadUrl
    a.download = it.name || ''
    document.body.appendChild(a)
    a.click()
    a.remove()
    return
  }

  if (it.id) {
    const token = gls().sessionT || ''
    const uname = encodeURIComponent(userName || '')
    const url = `${API_DOWNLOAD}?file_id=${encodeURIComponent(it.id)}&token=${encodeURIComponent(token)}&username=${uname}&download=1`
    window.open(url, '_blank')
    return
  }

  window.alert('Téléchargement non disponible pour ce fichier.')
}

function startRename(it) {
  renameValue.value = it.name
  selectedFile.value = it
  showRenameModal.value = true
}
async function applyRename() {
  if (!selectedFile.value) { showRenameModal.value = false; return }
  const idx = files.value.findIndex(f => f.id === selectedFile.value.id)
  if (idx !== -1) {
    const newName = (renameValue.value || files.value[idx].name).trim()
    if (!newName) return
    await renameFileRequest(selectedFile.value.id, newName)
  }
  showRenameModal.value = false
  renameValue.value = ''
}

async function restoreFile(file) {
  if (!file) return
  const payload = { username: userName, token: gls().sessionT, file_id: file.id }
  try {
    const res = await fetch(API_RESTORE, {
      method: 'POST',
      mode: 'cors',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (!res.ok) throw new Error('restore failed')
    await fetchFiles()
    selectedFile.value = null
  } catch (err) {
    console.error('Erreur restauration', err)
    window.alert('Impossible de restaurer le fichier.')
  }
}

async function moveToTrash(fileId) {
  if (!fileId) return;

  const serverId = (() => {
    const f = files.value.find(f => f.id === fileId);
    return f?._server?.file_id ?? fileId;
  })();

  try {
    const payload = { username: userName, token: gls().sessionT, file_id: serverId }
    const res = await fetch(API_TRASH, {
      method: 'POST',
      mode: 'cors',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload)
    })
    if (!res.ok) {
      const txt = await res.text().catch(() => '')
      console.error('Échec mise à la corbeille', res.status, txt)
      window.alert('Échec mise à la corbeille : ' + (txt || res.statusText))
      return
    }
    files.value.forEach(f => {
      if (f.id === fileId) {
        f.folder = 'trash'
        f.parentPath = '/.trash/'
      }
    })
    await fetchFiles()
    selectedFile.value = null
  } catch (err) {
    console.error('Erreur moveToTrash', err)
    window.alert('Erreur lors de la mise à la corbeille')
  }
}

function confirmPermanentDelete(file) {
  if (!file) return
  const ok = window.confirm(`Supprimer définitivement "${file.name}" ? Cette action est irréversible.`)
  if (!ok) return
  deleteFileRequest(file.id)
}

// selection helpers
function selectFirst() { if (filteredFiles.value.length) openFile(filteredFiles.value[0]) }
function clearSelection() { selectedFile.value = null }

function ensureSelectedStillVisible() {
  if (!selectedFile.value) return
  const exists = files.value.find(f => f.id === selectedFile.value.id)
  if (!exists) { selectedFile.value = null; return }
  const visible = filteredFiles.value.find(f => f.id === selectedFile.value.id)
  if (!visible) selectedFile.value = null
}

// ---------- MOVE helpers & network calls ----------

// get server id if available
function serverIdFor(item) {
  return item?._server?.file_id ?? item?.id
}

// low-level API calls
async function doMoveFileById(serverFileId, destinationPath) {
  if (!serverFileId) throw new Error('missing server file id')
  const dest = normalizePath(destinationPath || '/')
  const payload = { username: userName, token: gls().sessionT, file_id: serverFileId, destination_path: dest }
  const res = await fetch(API_MOVE_FILE, {
    method: 'POST',
    mode: 'cors',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  })
  if (!res.ok) {
    const txt = await res.text().catch(()=> '')
    throw new Error(txt || `status ${res.status}`)
  }
  return await res.json().catch(()=> ({}))
}

async function doMoveFolderByPath(folderPath, destinationPath) {
  if (!folderPath) throw new Error('missing folder path')
  const dest = normalizePath(destinationPath || '/')
  const payload = { username: userName, token: gls().sessionT, folder_path: normalizePath(folderPath), destination_path: dest }
  const res = await fetch(API_MOVE_FOLDER, {
    method: 'POST',
    mode: 'cors',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload)
  })
  if (!res.ok) {
    const txt = await res.text().catch(()=> '')
    throw new Error(txt || `status ${res.status}`)
  }
  return await res.json().catch(()=> ({}))
}

// high-level UI function called by button (handles single or multiple)
async function moveFile(item) {
  // if multiple selection enabled -> move all selected files (files only)
  if (sMul.value) {
    const checkedPairs = files.value.map((f, i) => ({ f, checked: !!selectedC.value[i] }))
      .filter(x => x.checked)
    if (checkedPairs.length === 0) {
      window.alert('Aucun élément sélectionné')
      return
    }

    const defaultDest = currentPath.value || '/'
    const userInput = window.prompt(`Chemin de destination pour ${checkedPairs.length} élément(s) (ex: / ou /dossier/):`, defaultDest)
    if (userInput === null) return
    const dest = normalizePath(userInput)

    // only files; if some entries are folders, skip them or prompt? we skip folders here
    const filePairs = checkedPairs.filter(x => x.f.type !== 'folder')
    if (filePairs.length === 0) {
      window.alert('Aucun fichier (seulement des dossiers) sélectionné pour déplacement.')
      return
    }

    movingInProgress.value = true
    try {
      const failures = []
      for (const p of filePairs) {
        const sid = serverIdFor(p.f)
        if (!sid) {
          // update locally if no server id
          updateFileParentPath(p.f.id, dest)
          continue
        }
        try {
          await doMoveFileById(sid, dest)
        } catch (err) {
          console.error('move file error', err)
          failures.push({ id: p.f.id, err: err.message || String(err) })
        }
      }
      if (failures.length) window.alert(`${failures.length} fichier(s) n'ont pas pu être déplacés.`)
      await fetchFiles()
      // reset selection UI
      sMulAll.value = false
      selectedC.value = new Array(files.value.length).fill(false)
    } finally {
      movingInProgress.value = false
    }
    return
  }

  // single item mode: if item is folder -> route to folder move
  const defaultDest = currentPath.value || '/'
  const userInput = window.prompt('Chemin de destination (ex: / ou /dossier/):', defaultDest)
  if (userInput === null) return
  const dest = normalizePath(userInput)

  if (item && item.type === 'folder') {
    const folderPath = normalizePath((item.parentPath || '/') + item.name + '/')
    movingInProgress.value = true
    try {
      await doMoveFolderByPath(folderPath, dest)
      await fetchFiles()
    } catch (err) {
      console.error('move folder failed', err)
      window.alert('Impossible de déplacer le dossier: ' + (err.message || err))
      await fetchFiles()
    } finally {
      movingInProgress.value = false
    }
    return
  }

  // single file
  const fid = serverIdFor(item)
  if (!fid) {
    window.alert('Impossible de retrouver l\'identifiant serveur du fichier.')
    return
  }
  movingInProgress.value = true
  try {
    await doMoveFileById(fid, dest)
    await fetchFiles()
  } catch (err) {
    console.error('move file failed', err)
    window.alert('Impossible de déplacer le fichier: ' + (err.message || err))
    await fetchFiles()
  } finally {
    movingInProgress.value = false
  }
}

// ---------- Drag/drop handlers (async, guarded by movingInProgress) ----------
function onDragStart(event, it) {
  draggedItem.value = { id: it.id, type: it.type }
  try {
    if (event.dataTransfer) {
      event.dataTransfer.setData('application/x-drive-item', JSON.stringify({ id: it.id, type: it.type }))
      event.dataTransfer.effectAllowed = 'move'
    }
  } catch (err) {}
}
function onDragEnd() { draggedItem.value = null; dragTarget.value = null }
function onItemDragOver(it) { if (it && it.type === 'folder') dragTarget.value = it.id }
function onItemDragLeave(it) { if (it && dragTarget.value === it.id) dragTarget.value = null }

async function onDropIntoFolder(folder) {
  if (movingInProgress.value) return
  if (!draggedItem.value || !folder || folder.type !== 'folder') return
  movingInProgress.value = true
  try {
    const moving = files.value.find(f => f.id === draggedItem.value.id)
    if (!moving) return
    if (moving.type === 'folder') {
      const oldPath = normalizePath((moving.parentPath || '/') + moving.name + '/')
      const newParent = normalizePath((folder.parentPath || '/') + folder.name + '/')
      if (newParent.indexOf(oldPath) === 0 || oldPath === newParent) return
      await moveFolderToPath(moving, newParent)
    } else {
      await moveFileToFolder(moving, folder)
    }
    await fetchFiles()
  } catch (err) {
    console.error('dropIntoFolder error', err)
    window.alert('Erreur lors du déplacement : ' + (err.message || err))
    await fetchFiles()
  } finally {
    draggedItem.value = null
    dragTarget.value = null
    movingInProgress.value = false
  }
}

async function onDropUp() {
  if (movingInProgress.value) return
  if (!draggedItem.value) return
  movingInProgress.value = true
  try {
    const moving = files.value.find(f => f.id === draggedItem.value.id)
    if (!moving) return
    const parts = currentParts.value.slice()
    parts.pop()
    const parentPath = normalizePath('/' + (parts.length ? parts.join('/') + '/' : ''))
    if (moving.type === 'folder') {
      const oldPath = normalizePath((moving.parentPath || '/') + moving.name + '/')
      if (parentPath.indexOf(oldPath) === 0) return
      await moveFolderToPath(moving, parentPath)
    } else {
      await moveFileToPath(moving, parentPath)
    }
    await fetchFiles()
  } catch (err) {
    console.error('dropUp error', err)
    window.alert('Erreur lors du déplacement : ' + (err.message || err))
    await fetchFiles()
  } finally {
    draggedItem.value = null
    dragTarget.value = null
    movingInProgress.value = false
  }
}

async function onDrop(e) {
  dragOver.value = false
  if (movingInProgress.value) return
  if (draggedItem.value) {
    movingInProgress.value = true
    try {
      const moving = files.value.find(f => f.id === draggedItem.value.id)
      if (!moving) { return }
      if (moving.type === 'folder') await moveFolderToPath(moving, currentPath.value)
      else await moveFileToPath(moving, currentPath.value)
      await fetchFiles()
    } catch (err) {
      console.error('onDrop error', err)
      window.alert('Erreur lors du déplacement : ' + (err.message || err))
      await fetchFiles()
    } finally {
      draggedItem.value = null
      dragTarget.value = null
      movingInProgress.value = false
    }
    return
  }
  handleUpload(e)
}

async function onBreadcrumbDrop(i) {
  if (movingInProgress.value) return
  if (!draggedItem.value) return
  movingInProgress.value = true
  try {
    const parts = currentParts.value.slice(0, i + 1)
    const targetPath = normalizePath('/' + (parts.length ? parts.join('/') + '/' : ''))
    const moving = files.value.find(f => f.id === draggedItem.value.id)
    if (!moving) return
    if (moving.type === 'folder') {
      const oldPath = normalizePath((moving.parentPath || '/') + moving.name + '/')
      if (targetPath.indexOf(oldPath) === 0) return
      await moveFolderToPath(moving, targetPath)
    } else {
      await moveFileToPath(moving, targetPath)
    }
    await fetchFiles()
  } catch (err) {
    console.error('breadcrumbDrop error', err)
    window.alert('Erreur lors du déplacement : ' + (err.message || err))
    await fetchFiles()
  } finally {
    draggedItem.value = null
    dragTarget.value = null
    movingInProgress.value = false
  }
}

async function onBreadcrumbDropRoot() {
  if (movingInProgress.value) return
  if (!draggedItem.value) return
  movingInProgress.value = true
  try {
    const targetPath = '/'
    const moving = files.value.find(f => f.id === draggedItem.value.id)
    if (!moving) return
    if (moving.type === 'folder') {
      const oldPath = normalizePath((moving.parentPath || '/') + moving.name + '/')
      if (targetPath.indexOf(oldPath) === 0) return
      await moveFolderToPath(moving, targetPath)
    } else {
      await moveFileToPath(moving, targetPath)
    }
    await fetchFiles()
  } catch (err) {
    console.error('breadcrumbDropRoot error', err)
    window.alert('Erreur lors du déplacement : ' + (err.message || err))
    await fetchFiles()
  } finally {
    draggedItem.value = null
    dragTarget.value = null
    movingInProgress.value = false
  }
}

// ---------- move helpers (async, optimistic UI + server call + resync) ----------
function updateFileParentPath(id, newParentPath) {
  const idx = files.value.findIndex(f => f.id === id)
  if (idx !== -1) files.value[idx].parentPath = newParentPath
}

async function moveFileToFolder(file, folder) {
  if (!file || !folder || folder.type !== 'folder') return
  const newParent = normalizePath((folder.parentPath || '/') + folder.name + '/')
  await moveFileToPath(file, newParent)
}

async function moveFileToPath(file, newParentPath) {
  if (!file) return
  let dest = newParentPath || '/'
  dest = normalizePath(dest)

  const sid = serverIdFor(file)
  if (!sid) {
    updateFileParentPath(file.id, dest)
    ensureSelectedStillVisible()
    return
  }

  // optimistic update locally
  const oldParent = file.parentPath
  updateFileParentPath(file.id, dest)
  ensureSelectedStillVisible()

  try {
    await doMoveFileById(sid, dest)
    await fetchFiles()
  } catch (err) {
    console.error('moveFileToPath failed', err)
    window.alert('Erreur lors du déplacement : ' + (err.message || err))
    await fetchFiles()
  }
}

async function moveFolderToFolder(folderToMove, targetFolder) {
  if (!folderToMove || !targetFolder) return
  if (folderToMove.id === targetFolder.id) return
  const newParent = normalizePath((targetFolder.parentPath || '/') + targetFolder.name + '/')
  await moveFolderToPath(folderToMove, newParent)
}

async function moveFolderToPath(folderToMove, newParentPath) {
  if (!folderToMove) return
  let targetPath = newParentPath || '/'
  targetPath = normalizePath(targetPath)

  const oldPath = normalizePath((folderToMove.parentPath || '/') + folderToMove.name + '/')
  if (targetPath.indexOf(oldPath) === 0) {
    window.alert("Impossible de déplacer un dossier dans lui-même ou dans un de ses sous-dossiers.")
    return
  }

  const sid = serverIdFor(folderToMove)
  // local-only fallback if no server id
  if (!sid) {
    const idx = files.value.findIndex(f => f.id === folderToMove.id)
    if (idx !== -1) files.value[idx].parentPath = targetPath
    const newPath = normalizePath(targetPath + folderToMove.name + '/')
    files.value.forEach(f => {
      if (f.parentPath && f.parentPath.indexOf(oldPath) === 0) {
        f.parentPath = f.parentPath.replace(oldPath, newPath)
      }
    })
    ensureSelectedStillVisible()
    return
  }

  // optimistic local update
  const idx = files.value.findIndex(f => f.id === folderToMove.id)
  if (idx !== -1) files.value[idx].parentPath = targetPath
  const newPath = normalizePath(targetPath + folderToMove.name + '/')
  files.value.forEach(f => {
    if (f.parentPath && f.parentPath.indexOf(oldPath) === 0) {
      f.parentPath = f.parentPath.replace(oldPath, newPath)
    }
  })
  ensureSelectedStillVisible()

  try {
    await doMoveFolderByPath(oldPath, targetPath)
    await fetchFiles()
  } catch (err) {
    console.error('moveFolderToPath failed', err)
    window.alert('Erreur lors du déplacement du dossier : ' + (err.message || err))
    await fetchFiles()
  }
}

// filtered files
const filteredFiles = computed(() => {
  const q = (searchQuery.value || '').trim().toLowerCase()
  let list = []
  if (curr.value === 'trash') {
    list = files.value.filter(f => f.folder === 'trash' && f.parentPath === currentPath.value)
  } else if (curr.value === 'shared') {
    list = files.value.filter(f => f.folder === 'shared' && f.parentPath === currentPath.value)
  } else {
    list = files.value.filter(f => f.folder === 'drive' && f.parentPath === currentPath.value)
  }
  if (!q) return list
  return list.filter(f => (f.name || '').toLowerCase().includes(q))
})

// Nombre de fichiers hors corbeille
const fileCount = computed(() => {
  return files.value.filter(f => f.folder !== 'trash').length
})

let _onlyofficeScriptLoading = null
function loadOnlyOfficeScript() {
  if (window.DocEditor || (window.DocsAPI && window.DocsAPI.DocEditor)) return Promise.resolve()
  if (_onlyofficeScriptLoading) return _onlyofficeScriptLoading

  _onlyofficeScriptLoading = new Promise((resolve, reject) => {
    try {
      // Determine OnlyOffice base URL from environment or fallback
      let cfg = (import.meta.env.VITE_ONLYOFFICE_URL || '').trim()
      let docserverBase = ''
      if (cfg) {
        // strip possible trailing /web-apps/.../api.js if user provided full path
        docserverBase = cfg.replace(/\/web-apps\/apps\/api\/documents\/api\.js\/?$/i, '').replace(/\/+$/, '')
      } else {
        // Utiliser la config centralisée
        docserverBase = SERVICES_CONFIG.getDocsURL()
      }
      const src = `${docserverBase}/web-apps/apps/api/documents/api.js`

      // avoid double-insert
      if (document.querySelector(`script[src="${src}"]`)) {
        // attach to existing script load if not ready
        const existing = document.querySelector(`script[src="${src}"]`)
        existing.addEventListener('load', () => setTimeout(() => resolve(), 50))
        existing.addEventListener('error', (e) => reject(e))
        return
      }

      const s = document.createElement('script')
      s.src = src
      s.async = true
      s.onload = () => setTimeout(() => resolve(), 50)
      s.onerror = (e) => {
        console.error('Failed to load OnlyOffice api.js from', src, e)
        reject(new Error('Failed to load OnlyOffice api.js'))
      }
      document.head.appendChild(s)
    } catch (err) {
      reject(err)
    }
  })
  return _onlyofficeScriptLoading
}

async function openInOnlyOffice(file) {
  try {
    const fileId = file.file_id || file.id
    if (!fileId) return
    const username = gls().username || localStorage.getItem('username') || ''
    const token = gls().sessionT || localStorage.getItem('token') || ''
    if (!username || !token) { alert('Session invalide'); return }

    // fetch config WITH credentials to allow backend session validation
    const apiBase = resolveAPI('')
    const cfgUrl = `${apiBase.replace(/\/+$/, '')}/api/onlyoffice/config?file_id=${encodeURIComponent(fileId)}&token=${encodeURIComponent(token)}&username=${encodeURIComponent(username)}`
    const cfgResp = await fetch(cfgUrl, { credentials: 'include', headers: { 'Accept': 'application/json' } })
    if (!cfgResp.ok) {
      const txt = await cfgResp.text().catch(()=> '')
      console.error('OnlyOffice config fetch failed', cfgResp.status, txt)
      alert('Impossible de récupérer la configuration OnlyOffice (voir console).')
      return
    }
    const cfg = await cfgResp.json()
    console.log('OnlyOffice config:', cfg)

    // load api.js
    await loadOnlyOfficeScript()

    const container = document.getElementById('onlyofficeModalContainer') || document.getElementById('onlyofficeContainer')
    if (!container) { console.error('onlyoffice container missing'); return }
    container.innerHTML = ''

    const DocEditorCtor = window.DocEditor || (window.DocsAPI && window.DocsAPI.DocEditor)
    if (!DocEditorCtor) {
      console.error('DocEditor constructor not found after loading api.js')
      alert('Impossible d\'initialiser OnlyOffice (api.js manquant).')
      return
    }

    // ensure editor config uses full-width/height
    cfg.width = '100%'
    cfg.height = '100%'

    // instantiate and keep reference
    try {
      const containerId = (container && container.id) ? container.id : 'onlyofficeContainer'
      const editor = new DocEditorCtor(containerId, cfg)
      window._onlyofficeEditor = editor
      console.log('OnlyOffice editor created in', containerId)
      
      // Auto-focus avec plusieurs tentatives pour garantir le focus
      const attemptFocus = (attempts = 0) => {
        if (attempts > 5) return // Max 5 tentatives
        
        setTimeout(() => {
          const iframe = container.querySelector('iframe')
          if (iframe) {
            try {
              // Focus sur l'iframe
              iframe.focus()
              
              // Focus sur le contenu de l'iframe si accessible
              if (iframe.contentWindow) {
                iframe.contentWindow.focus()
              }
              
              // Focus sur le document à l'intérieur si accessible
              if (iframe.contentDocument) {
                iframe.contentDocument.body?.focus()
              }
              
              // Simuler un clic pour activer l'éditeur
              iframe.click()
              
              console.log(`OnlyOffice focus attempt ${attempts + 1} successful`)
            } catch (e) {
              console.log(`OnlyOffice focus attempt ${attempts + 1} failed:`, e)
              // Réessayer après un délai plus long
              attemptFocus(attempts + 1)
            }
          } else if (attempts < 5) {
            // Si l'iframe n'est pas encore chargée, réessayer
            attemptFocus(attempts + 1)
          }
        }, 300 + (attempts * 200)) // Délais progressifs: 300ms, 500ms, 700ms, etc.
      }
      
      // Lancer les tentatives de focus
      attemptFocus()
    } catch (err) {
      console.error('Failed to instantiate OnlyOffice editor', err)
      alert('Erreur lors de l\'initialisation de l\'éditeur OnlyOffice (voir console).')
    }
  } catch (err) {
    console.error('openInOnlyOffice failed:', err)
    alert('Erreur lors de l\'ouverture dans OnlyOffice. Voir console.')
  }
}

function shareFile(file){
  // open share modal and load current shares for this file (if owner)
  selectedFile.value = file
  shareWithUsername.value = ''
  sharePermission.value = 'editor'
  currentShares.value = []
  showShareModal.value = true
  loadCurrentShares(file)
}

async function loadCurrentShares(file) {
  if (!file) return
  const payload = { username: userName, token: gls().sessionT, file_id: file.file_id ?? file.id }
  try {
    const res = await fetch(API_GET_SHARES, {
      method: 'POST',
      mode: 'cors',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(payload)
    })
    if (!res.ok) {
      console.error('loadCurrentShares failed', res.status, await res.text().catch(()=>''))
      currentShares.value = []
      return
    }
    const j = await res.json()
    currentShares.value = Array.isArray(j.shared_files) ? j.shared_files : []
  } catch (err) {
    console.error('loadCurrentShares error', err)
    currentShares.value = []
  }
}

async function shareWithUser() {
  if (!selectedFile.value) return
  if (!shareWithUsername.value || !shareWithUsername.value.trim()) {
    window.alert('Nom d\'utilisateur cible requis')
    return
  }
  const payload = {
    username: userName,
    token: gls().sessionT,
    file_id: selectedFile.value.file_id ?? selectedFile.value.id,
    share_with_username: shareWithUsername.value.trim(),
    permission: sharePermission.value || 'editor'
  }
  try {
    const res = await fetch(API_SHARE, {
      method: 'POST',
      mode: 'cors',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(payload)
    })
    if (!res.ok) {
      const txt = await res.text().catch(()=> '')
      console.error('share failed', res.status, txt)
      window.alert('Échec du partage : ' + (txt || res.statusText))
      return
    }
    await loadCurrentShares(selectedFile.value)
    shareWithUsername.value = ''
  } catch (err) {
    console.error('shareWithUser error', err)
    window.alert('Erreur lors du partage')
  }
}

async function unshareWith(share) {
  if (!selectedFile.value || !share) return
  if (!window.confirm(`Arrêter le partage avec ${share.shared_with || share.shared_with_name || share.shared_with_username || share.shared_with_name}?`)) return
  const payload = {
    username: userName,
    token: gls().sessionT,
    file_id: selectedFile.value.file_id ?? selectedFile.value.id,
    unshare_with_username: share.shared_with || share.shared_with_name || share.shared_with_username
  }
  try {
    const res = await fetch(API_UNSHARE, {
      method: 'POST',
      mode: 'cors',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(payload)
    })
    if (!res.ok) {
      const txt = await res.text().catch(()=> '')
      console.error('unshare failed', res.status, txt)
      window.alert('Échec annulation partage : ' + (txt || res.statusText))
      return
    }
    await loadCurrentShares(selectedFile.value)
  } catch (err) {
    console.error('unshareWith error', err)
    window.alert('Erreur lors de l\'annulation du partage')
  }
}

// small helper to close modal
function closeShareModal() {
  showShareModal.value = false
  selectedFile.value = null
  shareWithUsername.value = ''
  sharePermission.value = 'editor'
  currentShares.value = []
}

// OnlyOffice fullscreen state
const isOnlyofficeFullscreen = ref(false)

function toggleOnlyofficeFullscreen() {
  isOnlyofficeFullscreen.value = !isOnlyofficeFullscreen.value
  
  // Auto-focus le conteneur après transition
  setTimeout(() => {
    focusOnlyOfficeEditor()
  }, 100)
}

// Fonction pour forcer le focus sur l'éditeur OnlyOffice
function focusOnlyOfficeEditor() {
  const container = document.getElementById('onlyofficeModalContainer')
  if (!container) return
  
  // Focus sur le conteneur
  container.focus()
  
  // Essayer de focus sur l'iframe OnlyOffice
  const iframe = container.querySelector('iframe')
  if (iframe) {
    try {
      iframe.focus()
      
      // Focus sur le contenu de l'iframe
      if (iframe.contentWindow) {
        iframe.contentWindow.focus()
      }
      
      // Focus sur le document à l'intérieur
      if (iframe.contentDocument && iframe.contentDocument.body) {
        iframe.contentDocument.body.focus()
      }
      
      console.log('OnlyOffice editor focused')
    } catch (e) {
      console.log('Cannot focus OnlyOffice iframe:', e)
    }
  }
}

// destroy onlyoffice editor when modal closed
watch(showOnlyofficeModal, (v) => {
  if (!v) {
    isOnlyofficeFullscreen.value = false
    try {
      if (window._onlyofficeEditor && typeof window._onlyofficeEditor.destroy === 'function') {
        window._onlyofficeEditor.destroy()
      }
    } catch (e) { console.error('destroy onlyoffice editor failed', e) }
    try { const el = document.getElementById('onlyofficeModalContainer'); if (el) el.innerHTML = '' } catch(e) {}
    window._onlyofficeEditor = null
  }
})

onMounted(async () => {
  try {
    await fetchFiles()
    await fetchStorageInfo()
  } catch (e) {
    console.error('fetchFiles onMounted failed', e)
  }
})

// optionnel : relancer fetchFiles quand la session est définie/renouvelée
watch(() => gls().sessionT, (newToken) => {
  if (newToken) {
    fetchFiles().catch(err => console.error('fetchFiles on session change failed', err))
    fetchStorageInfo().catch(err => console.error('fetchStorageInfo on session change failed', err))
  }
})

// don't auto-open OnlyOffice on selection; user opens editor via the "Ouvrir" button which opens the modal
</script>

<template>
  <div class="drive-page">
    <input ref="uploadInput" type="file" multiple style="display:none" @change="handleUpload" />

    <!-- Hero -->
    <div class="hero">
      <div class="hero-badge">Drive</div>
      <h1>Vos fichiers</h1>
      <p class="hero-sub">Gérez vos documents en toute sécurité.</p>
    </div>

    <!-- Storage bar -->
    <div class="storage">
      <div class="storage-head">
        <span class="storage-label">Stockage</span>
        <span class="storage-pct">{{ isUnlimitedStorage ? 'Illimité' : storagePercentage + '%' }}</span>
      </div>
      <div class="storage-bar-bg">
        <div class="storage-bar-fill" :style="{ width: (isUnlimitedStorage ? '50%' : storagePercentage + '%'), background: storageColor }"></div>
      </div>
      <span class="storage-detail">{{ humanSize(storageUsed) }}{{ isUnlimitedStorage ? '' : ' / ' + humanSize(storageLimit) }}</span>
    </div>

    <!-- Main layout -->
    <div class="layout">
      <!-- Sidebar -->
      <aside class="sidebar">
        <nav class="side-nav">
          <button :class="{ active: curr === 'drive' }" @click="changec('drive')" class="nav-btn">Drive</button>
          <button :class="{ active: curr === 'shared' }" @click="changec('shared')" class="nav-btn">Partagés</button>
          <button :class="{ active: curr === 'trash' }" @click="changec('trash')" class="nav-btn">Corbeille</button>
        </nav>
        <div class="side-user">
          <div class="side-avatar">{{ (userName || '?').charAt(0).toUpperCase() }}</div>
          <div class="side-user-info">
            <span class="side-name">{{ userName }}</span>
            <span v-if="userEmail" class="side-email">{{ userEmail }}</span>
          </div>
        </div>
      </aside>

      <!-- Content -->
      <div class="content">
        <!-- Toolbar -->
        <div class="toolbar">
          <div class="breadcrumbs">
            <button @click="goRoot" class="crumb" :class="{ 'crumb-active': currentPath === '/' }">Root</button>
            <template v-for="(part, i) in currentParts" :key="i">
              <span class="sep">/</span>
              <button @click="goToBreadcrumb(i)" class="crumb" :class="{ 'crumb-active': i === currentParts.length - 1 }">{{ part }}</button>
            </template>
          </div>
          <div class="toolbar-actions">
            <input v-model="searchQuery" type="text" placeholder="Rechercher..." class="search" />
            <button @click="showNewFolderModal = true" class="tool-btn">Dossier</button>
            <button @click="showNewFileModal = true" class="tool-btn">Fichier</button>
            <button @click="triggerFilePicker" class="tool-btn primary">Upload</button>
          </div>
        </div>

        <!-- File list -->
        <div class="files" @dragover.prevent="dragOver = true" @dragleave.prevent="dragOver = false" @drop.prevent="onDrop">
          <div v-if="filteredFiles.length === 0" class="empty-state">
            <p>Aucun fichier</p>
          </div>
          <div v-for="(f, i) in filteredFiles" :key="f.id" class="file-item"
            :class="{ selected: selectedFile?.id === f.id, dragtarget: dragTarget === f.id }"
            @click="selectItem(f)"
            draggable="true" @dragstart="onDragStart($event, f)" @dragend="onDragEnd" @dragover="onItemDragOver(f)" @dragleave="onItemDragLeave(f)" @drop.prevent="onDropIntoFolder(f)">
            <span class="file-icon">{{ f.type === 'folder' ? '📁' : '📄' }}</span>
            <span class="file-name">{{ f.name }}</span>
            <span v-if="f.sizeLabel" class="file-size">{{ f.sizeLabel }}</span>
            <span class="file-date">{{ f.date }}</span>
            <div class="file-actions" @click.stop>
              <button v-if="f.type !== 'folder'" @click="openAction(f)" class="fa-btn">Ouvrir</button>
              <button v-if="f.type !== 'folder'" @click="downloadFile(f)" class="fa-btn">DL</button>
              <button @click="startRename(f)" class="fa-btn">Renommer</button>
              <button @click="shareFile(f)" class="fa-btn">Partager</button>
              <button v-if="curr === 'trash'" @click="restoreFile(f)" class="fa-btn">Restaurer</button>
              <button v-else @click="moveToTrash(f.id)" class="fa-btn">Supprimer</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- FAB -->
    <button @click="showMobileFab = !showMobileFab" class="fab">+</button>
    <div v-if="showMobileFab" class="fab-menu">
      <button @click="triggerFilePicker; showMobileFab = false" class="fab-item">Upload</button>
      <button @click="showNewFolderModal = true; showMobileFab = false" class="fab-item">Dossier</button>
      <button @click="showNewFileModal = true; showMobileFab = false" class="fab-item">Fichier</button>
    </div>

    <!-- Moving overlay -->
    <div v-if="movingInProgress" class="moving-overlay">
      <div class="spinner-lg"></div>
      <span>Déplacement en cours...</span>
    </div>

    <!-- New Folder Modal -->
    <div v-if="showNewFolderModal" class="overlay" @click.self="showNewFolderModal = false">
      <div class="dialog">
        <div class="dialog-head"><h3>Nouveau dossier</h3><button @click="showNewFolderModal = false" class="close">×</button></div>
        <div class="dialog-body"><input v-model="newFolderName" type="text" placeholder="Nom du dossier" @keyup.enter="createFolder" /></div>
        <div class="dialog-foot"><button @click="showNewFolderModal = false" class="btn-cancel">Annuler</button><button @click="createFolder" class="btn-confirm">Créer</button></div>
      </div>
    </div>

    <!-- New File Modal -->
    <div v-if="showNewFileModal" class="overlay" @click.self="showNewFileModal = false">
      <div class="dialog">
        <div class="dialog-head"><h3>Nouveau fichier</h3><button @click="showNewFileModal = false" class="close">×</button></div>
        <div class="dialog-body">
          <input v-model="newFileName" type="text" placeholder="Nom du fichier" class="mb" />
          <select v-model="newFileType" class="mb">
            <option value="docx">Document Word</option>
            <option value="xlsx">Tableur Excel</option>
            <option value="pptx">Présentation PowerPoint</option>
          </select>
        </div>
        <div class="dialog-foot"><button @click="showNewFileModal = false" class="btn-cancel">Annuler</button><button @click="createFile" class="btn-confirm">Créer</button></div>
      </div>
    </div>

    <!-- Rename Modal -->
    <div v-if="showRenameModal" class="overlay" @click.self="showRenameModal = false">
      <div class="dialog">
        <div class="dialog-head"><h3>Renommer</h3><button @click="showRenameModal = false" class="close">×</button></div>
        <div class="dialog-body"><input v-model="renameValue" type="text" @keyup.enter="applyRename" /></div>
        <div class="dialog-foot"><button @click="showRenameModal = false" class="btn-cancel">Annuler</button><button @click="applyRename" class="btn-confirm">Renommer</button></div>
      </div>
    </div>

    <!-- Upload Modal -->
    <div v-if="showUploadModal" class="overlay" @click.self="showUploadModal = false">
      <div class="dialog">
        <div class="dialog-head"><h3>Uploads ({{ uploadsSummary }})</h3><button @click="showUploadModal = false" class="close">×</button></div>
        <div class="dialog-body">
          <div class="progress-bar-bg"><div class="progress-bar-fill" :style="{ width: overallProgress + '%' }"></div></div>
          <p class="progress-pct">{{ overallProgress }}%</p>
        </div>
        <div class="dialog-foot"><button @click="cancelAllUploads" class="btn-cancel">Annuler tout</button></div>
      </div>
    </div>

    <!-- Share Modal -->
    <div v-if="showShareModal" class="overlay" @click.self="closeShareModal">
      <div class="dialog">
        <div class="dialog-head"><h3>Partager "{{ selectedFile?.name }}"</h3><button @click="closeShareModal" class="close">×</button></div>
        <div class="dialog-body">
          <div class="share-row">
            <input v-model="shareWithUsername" type="text" placeholder="Nom d'utilisateur" />
            <select v-model="sharePermission"><option value="editor">Éditeur</option><option value="viewer">Lecteur</option></select>
            <button @click="shareWithUser" class="btn-share">Partager</button>
          </div>
          <div v-if="currentShares.length > 0" class="shares-list">
            <div v-for="s in currentShares" :key="s.shared_with" class="share-item">
              <span>{{ s.shared_with }}</span>
              <button @click="unshareWith(s)" class="btn-unshare">Retirer</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- OnlyOffice Modal -->
    <div v-if="showOnlyofficeModal" class="overlay onlyoffice-overlay" @click.self="showOnlyofficeModal = false">
      <div class="dialog onlyoffice-dialog" :class="{ fullscreen: isOnlyofficeFullscreen }">
        <div class="dialog-head">
          <h3>{{ selectedFile?.name }}</h3>
          <div>
            <button @click="toggleOnlyofficeFullscreen" class="close">{{ isOnlyofficeFullscreen ? '🗗' : '🗖' }}</button>
            <button @click="showOnlyofficeModal = false" class="close">×</button>
          </div>
        </div>
        <div class="onlyoffice-body"><div id="onlyofficeModalContainer" style="width:100%;height:100%"></div></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.drive-page { max-width: var(--content-max); margin: 0 auto; padding: 2rem 1.5rem 5rem; animation: fadeInUp .5s var(--ease) both; }

/* Hero */
.hero { text-align: center; padding: 2.5rem 2rem; background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--r-2xl); margin-bottom: 1.5rem; animation: fadeInUp .5s var(--ease) both; }
.hero-badge { display: inline-block; padding: .3rem .7rem; font-size: .75rem; font-weight: 700; border-radius: var(--r-full); background: var(--bg-tertiary); color: var(--text-secondary); margin-bottom: .75rem; animation: fadeInDown .4s var(--ease) both; }
.hero h1 { font-size: clamp(1.75rem, 4vw, 2.25rem); font-weight: 800; letter-spacing: -.03em; margin-bottom: .375rem; animation: fadeInUp .5s var(--ease) .1s both; }
.hero-sub { font-size: .9375rem; color: var(--text-secondary); animation: fadeInUp .5s var(--ease) .2s both; }

/* Storage */
.storage { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1rem 1.25rem; margin-bottom: 1.5rem; animation: fadeInUp .4s var(--ease) .1s both; }
.storage-head { display: flex; justify-content: space-between; margin-bottom: .5rem; }
.storage-label { font-size: .8125rem; font-weight: 600; color: var(--text-secondary); }
.storage-pct { font-size: .8125rem; font-weight: 600; color: var(--text-primary); }
.storage-bar-bg { height: 6px; background: var(--bg-tertiary); border-radius: var(--r-full); overflow: hidden; }
.storage-bar-fill { height: 100%; border-radius: var(--r-full); transition: width var(--ts); animation: growBar .6s var(--ease) both; }
.storage-detail { font-size: .75rem; color: var(--text-muted); margin-top: .375rem; display: block; }

/* Layout */
.layout { display: grid; grid-template-columns: 200px 1fr; gap: 1.25rem; }
.sidebar { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); padding: 1rem; display: flex; flex-direction: column; gap: .5rem; animation: fadeInLeft .4s var(--ease) both; }
.side-nav { display: flex; flex-direction: column; gap: .25rem; }
.nav-btn { padding: .625rem .875rem; font-size: .8125rem; font-weight: 500; color: var(--text-secondary); background: transparent; border: none; border-radius: var(--r); cursor: pointer; transition: all var(--tb); text-align: left; }
.nav-btn:hover { color: var(--text-primary); background: var(--bg-tertiary); }
.nav-btn.active { color: var(--paper); background: var(--ink); }
.side-user { margin-top: auto; padding-top: 1rem; border-top: 1px solid var(--border); display: flex; align-items: center; gap: .625rem; }
.side-avatar { width: 32px; height: 32px; border-radius: 50%; background: var(--bg-tertiary); border: 1px solid var(--border); display: flex; align-items: center; justify-content: center; font-size: .8125rem; font-weight: 700; flex-shrink: 0; }
.side-user-info { display: flex; flex-direction: column; overflow: hidden; }
.side-name { font-size: .8125rem; font-weight: 600; }
.side-email { font-size: .6875rem; color: var(--text-muted); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* Content */
.content { display: flex; flex-direction: column; gap: .875rem; animation: fadeInRight .4s var(--ease) both; }
.toolbar { display: flex; align-items: center; justify-content: space-between; gap: .75rem; flex-wrap: wrap; }
.breadcrumbs { display: flex; align-items: center; gap: .25rem; flex-wrap: wrap; }
.crumb { padding: .25rem .5625rem; font-size: .75rem; font-weight: 500; color: var(--text-secondary); background: transparent; border: none; border-radius: var(--r-sm); cursor: pointer; transition: all var(--tb); }
.crumb:hover { color: var(--text-primary); background: var(--bg-tertiary); }
.crumb-active { color: var(--text-primary); font-weight: 600; }
.sep { color: var(--text-muted); font-size: .75rem; }
.toolbar-actions { display: flex; gap: .375rem; align-items: center; flex-wrap: wrap; }
.search { width: 160px; padding: .375rem .625rem; font-size: .8125rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r-sm); transition: border-color var(--tb); }
.search:focus { outline: none; border-color: var(--ink); }
.search::placeholder { color: var(--text-muted); }
.tool-btn { padding: .375rem .75rem; font-size: .8125rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r-sm); cursor: pointer; transition: all var(--tb); }
.tool-btn:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.tool-btn.primary { color: var(--paper); background: var(--ink); border-color: var(--ink); }
.tool-btn.primary:hover { background: var(--ink-soft); border-color: var(--ink-soft); }

/* Files */
.files { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); overflow: hidden; min-height: 300px; }
.empty-state { display: flex; align-items: center; justify-content: center; min-height: 300px; color: var(--text-muted); font-size: .875rem; }
.file-item { display: flex; align-items: center; gap: .75rem; padding: .75rem 1rem; border-bottom: 1px solid var(--border); cursor: pointer; transition: background var(--tb); animation: fadeIn .2s var(--ease) both; }
.file-item:last-child { border-bottom: none; }
.file-item:hover { background: var(--bg-secondary); }
.file-item.selected { background: var(--bg-tertiary); }
.file-item.dragtarget { background: var(--bg-tertiary); border: 2px dashed var(--ink); }
.file-icon { font-size: 1.25rem; flex-shrink: 0; }
.file-name { flex: 1; font-size: .875rem; font-weight: 500; color: var(--text-primary); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.file-size { font-size: .75rem; color: var(--text-muted); white-space: nowrap; }
.file-date { font-size: .75rem; color: var(--text-muted); white-space: nowrap; }
.file-actions { display: flex; gap: .25rem; opacity: 0; transition: opacity var(--tb); }
.file-item:hover .file-actions { opacity: 1; }
.fa-btn { padding: .25rem .5625rem; font-size: .6875rem; font-weight: 600; color: var(--text-secondary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r-sm); cursor: pointer; transition: all var(--tb); white-space: nowrap; }
.fa-btn:hover { border-color: var(--ink); color: var(--text-primary); background: var(--bg-tertiary); }

/* FAB */
.fab { position: fixed; bottom: 1.5rem; right: 1.5rem; width: 52px; height: 52px; border-radius: 50%; background: var(--ink); color: var(--paper); border: none; font-size: 24px; cursor: pointer; box-shadow: var(--sh-lg); z-index: 100; transition: transform var(--tb); display: flex; align-items: center; justify-content: center; }
.fab:hover { transform: scale(1.08) rotate(90deg); box-shadow: var(--sh-xl); }
.fab-menu { position: fixed; bottom: 5.5rem; right: 1.5rem; display: flex; flex-direction: column; gap: .5rem; z-index: 100; }
.fab-item { padding: .625rem 1.125rem; font-size: .8125rem; font-weight: 600; background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; box-shadow: var(--sh-md); transition: all var(--tb); white-space: nowrap; animation: fadeInUp .2s var(--ease) both; }
.fab-item:hover { border-color: var(--ink); box-shadow: var(--sh-lg); transform: translateX(-4px); }

/* Moving overlay */
.moving-overlay { position: fixed; inset: 0; background: var(--bg-overlay); backdrop-filter: blur(4px); z-index: 200; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 1rem; font-size: 1rem; font-weight: 600; color: var(--text-primary); animation: fadeIn .2s var(--ease) both; }
.spinner-lg { width: 40px; height: 40px; border: 3px solid var(--border); border-top-color: var(--ink); border-radius: 50%; animation: spin .8s linear infinite; }

/* Modals */
.overlay { position: fixed; inset: 0; z-index: 999; background: var(--bg-overlay); backdrop-filter: blur(6px); display: flex; align-items: center; justify-content: center; animation: fadeIn .2s var(--ease) both; }
.dialog { width: calc(100% - 32px); max-width: 480px; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); box-shadow: var(--sh-xl); animation: popIn .25s var(--ease-spring) both; max-height: 90vh; display: flex; flex-direction: column; }
.dialog-head { display: flex; align-items: center; justify-content: space-between; padding: 1rem 1.25rem; border-bottom: 1px solid var(--border); }
.dialog-head h3 { font-size: 1rem; font-weight: 700; }
.dialog-head > div { display: flex; gap: .375rem; }
.close { width: 30px; height: 30px; display: flex; align-items: center; justify-content: center; font-size: 18px; line-height: 1; background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r-sm); cursor: pointer; color: var(--text-secondary); transition: all var(--tb); }
.close:hover { border-color: var(--ink); background: var(--bg-tertiary); transform: rotate(90deg); color: var(--text-primary); }
.dialog-body { padding: 1.25rem; overflow-y: auto; }
.dialog-body input, .dialog-body select { width: 100%; padding: .5625rem .75rem; font-size: .875rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); transition: border-color var(--tb), box-shadow var(--tb); }
.dialog-body input:focus, .dialog-body select:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.mb { margin-bottom: .625rem; }
.dialog-foot { display: flex; gap: .5rem; padding: 1rem 1.25rem; border-top: 1px solid var(--border); background: var(--bg-secondary); }
.btn-cancel { flex: 1; padding: .5625rem; font-size: .8125rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-cancel:hover { border-color: var(--ink); background: var(--bg-tertiary); }
.btn-confirm { flex: 1; padding: .5625rem; font-size: .8125rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.btn-confirm:hover { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-sm); }

/* Upload progress */
.progress-bar-bg { height: 6px; background: var(--bg-tertiary); border-radius: var(--r-full); overflow: hidden; }
.progress-bar-fill { height: 100%; background: var(--ink); border-radius: var(--r-full); transition: width var(--tb); }
.progress-pct { text-align: center; margin-top: .625rem; font-size: .8125rem; font-weight: 600; color: var(--text-primary); }

/* Share */
.share-row { display: flex; gap: .375rem; margin-bottom: 1rem; }
.share-row input { flex: 1; }
.share-row select { width: auto; min-width: 100px; }
.btn-share { padding: .5625rem .875rem; font-size: .8125rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); white-space: nowrap; }
.btn-share:hover { background: var(--ink-soft); border-color: var(--ink-soft); }
.shares-list { display: flex; flex-direction: column; gap: .375rem; }
.share-item { display: flex; align-items: center; justify-content: space-between; padding: .5rem .75rem; background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--r); animation: fadeIn .2s var(--ease) both; }
.share-item span { font-size: .8125rem; font-weight: 500; }
.btn-unshare { padding: .25rem .5625rem; font-size: .6875rem; font-weight: 600; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r-sm); cursor: pointer; transition: all var(--tb); }
.btn-unshare:hover { border-color: var(--ink); }

/* OnlyOffice */
.onlyoffice-overlay { z-index: 1001; }
.onlyoffice-dialog { max-width: 900px; height: 85vh; }
.onlyoffice-dialog.fullscreen { max-width: 100vw; width: 100vw; height: 100vh; border-radius: 0; border: none; }
.onlyoffice-body { flex: 1; overflow: hidden; padding: 0; }

@media (max-width: 768px) {
  .layout { grid-template-columns: 1fr; }
  .sidebar { flex-direction: row; }
  .side-nav { flex-direction: row; overflow-x: auto; }
  .nav-btn { white-space: nowrap; }
  .side-user { display: none; }
  .toolbar { flex-direction: column; align-items: stretch; }
  .search { width: 100%; }
  .file-actions { opacity: 1; }
}
</style>
