<template>
  <div id="drive-v2">
    <!-- Hero Card -->
    <div class="drive-hero-card">
      <div class="hero-icon-wrapper">
        <svg class="hero-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z"/>
        </svg>
      </div>
      <h1>☁️ Drive Office1789</h1>
      <p class="hero-subtitle">Gérez vos fichiers en toute sécurité</p>
    </div>

    <!-- Drive Content Wrapper -->
    <div class="drive-content-wrapper">
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
        <button class="dv-btn" @click="showNewFileModal = true">📄 Nouveau fichier</button>
      </div>

      <nav class="dv-nav" aria-label="navigation drive">
        <button :class="['nav-item', { active: curr === 'drive' }]" @click="changec('drive')">📂 Mes fichiers</button>
        <button :class="['nav-item', { active: curr === 'shared' }]" @click="changec('shared')">🤝 Partagés</button>
        <button :class="['nav-item', { active: curr === 'trash' }]" @click="changec('trash')">🗑️ Corbeille</button>
      </nav>

      <div class="dv-footer">
        <small>{{ files.length }} éléments</small>
      </div>

      <!-- Storage Usage Bar - Only show in Drive and Shared, not in Trash -->
      <div v-if="curr !== 'trash'" class="storage-info">
        <div class="storage-header">
          <span class="storage-icon">💾</span>
          <span class="storage-text">Stockage</span>
        </div>
        <div class="storage-bar">
          <div class="storage-fill" :style="{ width: storagePercentage + '%', background: storageColor }"></div>
        </div>
        <div class="storage-details">
          <span>{{ humanSize(storageUsed) }} / {{ humanSize(storageLimit) }}</span>
          <span class="storage-percent">{{ storagePercentage }}%</span>
        </div>
        <div v-if="storagePercentage >= 90" class="storage-warning">
          ⚠️ Espace presque plein !
        </div>
      </div>
    </aside>

    <!-- MAIN -->
    <main class="dv-main">
      <div v-if="movingInProgress" class="moving-overlay" aria-live="polite" style="position:fixed;left:0;top:0;right:0;bottom:0;display:flex;align-items:center;justify-content:center;z-index:9999;pointer-events:none">
        <div style="background:rgba(0,0,0,0.6);color:#fff;padding:12px 18px;border-radius:10px;pointer-events:auto;">
          Déplacement en cours… ⏳
        </div>
      </div>
      <header class="dv-header">
       <div class="dv-breadcrumbs"
          @dragover.prevent
          @drop.prevent="onBreadcrumbDropRoot"
          role="navigation"
          aria-label="breadcrumbs">
        <button
          class="crumb"
          :class="{ 'crumb-target': dragTarget === 'breadcrumb-root', 'crumb-active': currentPath === '/' }"
          @click="goRoot"
          @dragover.prevent="dragTarget = 'breadcrumb-root'"
          @dragleave="dragTarget = null"
          @drop.prevent="onBreadcrumbDropRoot"
          aria-current="page"
          title="Root"
        >
          Root
        </button>

        <template v-for="(part, i) in currentParts" :key="i">
          <span class="slash">/</span>
          <button
            class="crumb"
            :class="{ 'crumb-target': dragTarget === `breadcrumb-${i}`, 'crumb-active': i === currentParts.length - 1 }"
            @click="goToBreadcrumb(i)"
            @dragover.prevent="dragTarget = `breadcrumb-${i}`"
            @dragleave="dragTarget = null"
            @drop.prevent="onBreadcrumbDrop(i)"
            :aria-label="`Aller à ${part}`"
          >
            {{ part }}
          </button>
        </template>
      </div>

        <div class="dv-tools">
          <input class="dv-search" v-model="searchQuery" placeholder="Rechercher..." />
          <button class="dv-btn small" :class="{ 'dv-btn small active': sMul}" @click="selectMultiple">Sélection Multiple</button>
          <label v-if="sMul" class="dv-btn small">Tout sélectionner :
            <input type="checkbox" id="select-mul" v-model="sMulAll" @click="selectAll" style="margin-right:6px" />
            </label>
          <button v-if="curr !== 'trash' && sMul" class="dv-btn danger" @click="moveSelectedToTrash">Mettre à la corbeille</button>
            <button v-else-if="curr === 'trash'&& sMul" class="dv-btn danger" @click="deleteSelectedPermanently">
              Supprimer définitivement
            </button>
          <button class="dv-btn green" v-if="sMul" @click="moveFile(selectedFile)">Déplacer fichiers</button>
          <button class="dv-btn green" v-if="sMul" @click="downloadSelectedAsZip(selectedFile)">Télécharger</button>
          <button class="dv-btn small" v-if="!sMul" @click="selectFirst">Sélectionner 1er</button>
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
          <!-- Remonter bouton (drop target + click) -->
          <div v-if="currentPath !== '/'" class="dv-up"
               @dragover.prevent="dragTarget = 'up'"
               @dragleave="dragTarget = null"
               @drop.prevent="onDropUp"
               :class="{ 'up-target': dragTarget === 'up' }"
          >
            <button class="dv-btn small" @click="goUp">⬆️ Remonter</button>
          </div>

          <div v-if="filteredFiles.length === 0" class="dv-empty">
            <div class="dv-empty-illustration" aria-hidden>
              <svg width="80" height="80" viewBox="0 0 24 24" fill="none">
                <rect width="24" height="24" rx="4" fill="#f0f4ff"/>
                <path d="M6 10c1.333-2 4-2 6 0 2 2.667 4.667 2.667 6 0" stroke="#a0b4ff" stroke-width="1.2" stroke-linecap="round" stroke-linejoin="round"/>
              </svg>
            </div>
            <p class="dv-empty-text">Aucun fichier ici</p>
            <p class="dv-empty-sub">Glisse-dépose un fichier/dossier ici ou clique sur <strong>Upload</strong>.</p>
          </div>

          <ul v-else class="dv-items" role="list">
            <li
              v-for="(it, index) in filteredFiles"
              :key="it.id"
              class="dv-item"
              :class="{
                active: selectedFile && selectedFile.id === it.id,
                'folder-hover': dragTarget === it.id
              }"
              @click="it.type === 'folder' ? openFolder(it) : openFile(it)"
              draggable="true"
              @dragstart="onDragStart($event, it)"
              @dragend="onDragEnd"
              @dragover.prevent="onItemDragOver(it)"
              @dragleave="onItemDragLeave(it)"
              @drop.prevent="it.type === 'folder' && onDropIntoFolder(it)"
            >
            <input v-if="sMul" @click.stop type="checkbox" v-model="selectedC[index]"/>
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

              <!-- si on est pas dans la corbeille, proposer mise à la corbeille -->
              <button v-if="it.folder !== 'trash'" class="dv-mini" @click.stop="moveToTrash(it.id)" title="Mettre à la corbeille">🗑️</button>
              

              <!-- si on est dans la corbeille, proposer suppression définitive (confirm) -->
              <button v-else class="dv-mini danger" @click.stop="confirmPermanentDelete(it)" title="Supprimer définitivement">🗑️❌</button>
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
                <button class="dv-btn green" @click="moveFile(selectedFile)">Déplacer fichier</button>
                <button v-if="curr !== 'trash'" class="dv-btn danger" @click="moveToTrash(selectedFile.id)">Mettre à la corbeille</button>
                
                <template v-else>
                  <button class="dv-btn" @click="restoreFile(selectedFile)">Restaurer</button>
                  <button class="dv-btn danger" @click="confirmPermanentDelete(selectedFile)">Supprimer définitivement</button>
                </template>
              </div>
            </div>

            <div class="detail-body">
              <div class="preview">
                <template v-if="selectedFile.url && (selectedFile.mime || '').startsWith('image/')">
                  <img :src="selectedFile.url" alt="preview" class="preview-img" :key="selectedFile.id" />
                </template>

                <template v-else-if="selectedFile.url && (selectedFile.mime || '').startsWith('video/')">
                  <video controls class="preview-video" :key="selectedFile.id">
                    <source :src="selectedFile.url" :type="selectedFile.mime" />
                    Ton navigateur ne supporte pas la lecture vidéo.
                  </video>
                </template>

                <template v-else-if="selectedFile.url && (selectedFile.mime || '').startsWith('audio/')">
                  <audio controls class="preview-audio" :key="selectedFile.id">
                    <source :src="selectedFile.url" :type="selectedFile.mime" />
                    Ton navigateur ne supporte pas la lecture audio.
                  </audio>
                </template>

                <!-- PDF inline preview -->
                <template v-else-if="selectedFile.url && (selectedFile.mime || '').toLowerCase() === 'application/pdf'">
                  <iframe :src="selectedFile.url" class="preview-pdf" :key="selectedFile.id" style="width:100%;height:400px;border:0"></iframe>
                </template>

                <!-- DOCX / other office: try opening in new tab (browser may download or display depending on support) -->
                <template v-else-if="selectedFile.url && ['docx','xlsx','pptx'].includes(selectedFile.name.split('.').pop().toLowerCase())">
                  <!-- OnlyOffice: open in modal -->
                  <div class="onlyoffice-placeholder" style="width:100%;height:100%;display:flex;align-items:center;justify-content:center;">
                    <div style="text-align:center">
                      <p style="margin:0 0 8px 0">Ce document doit être ouvert dans l'éditeur OnlyOffice.</p>
                      <div style="display:flex;gap:8px;justify-content:center">
                        <button class="dv-btn" @click="(showOnlyofficeModal = true, openInOnlyOffice(selectedFile))">Ouvrir dans OnlyOffice</button>
                        <button class="dv-btn" @click="downloadFile(selectedFile)">Télécharger</button>
                      </div>
                    </div>
                  </div>
                </template>

                <template v-else-if="selectedFile.text">
                  <pre class="preview-text" :key="selectedFile.id">{{ selectedFile.text }}</pre>
                </template>

                <template v-else>
                  <div class="no-preview">Aperçu non disponible pour ce fichier</div>
                </template>


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
    </div> <!-- End drive-content-wrapper -->

    <!-- hidden file input -->
    <input ref="uploadInput" type="file" multiple style="display:none" @change="handleUpload" />

    <!-- UPLOAD MODAL (avec file progress) -->
<div v-if="showUploadModal" class="modal-wrap" @keydown.esc="showUploadModal = false">
  <div class="modal" role="dialog" aria-modal="true" aria-label="Uploader des fichiers">
    <h3>Uploader des fichiers</h3>
    <p>Glisse-dépose ou clique pour sélectionner des fichiers.</p>

    <div class="upload-zone" @click="triggerFilePicker" @dragover.prevent @drop.prevent="handleUpload">
      <p>Déposez vos fichiers ici — ou cliquez</p>
      <small v-if="uploadQueue.length === 0" style="display:block; margin-top:8px; color:#666">
        Les fichiers seront envoyés et affichés ci-dessous avec une barre de progression.
      </small>
    </div>

    <input ref="uploadInput" type="file" multiple style="display:none" @change="handleUpload" />

    <div v-if="uploadQueue.length" style="margin-top:12px; max-height:40vh; overflow:auto; padding-right:4px">
      <div style="display:flex; justify-content:space-between; align-items:center; margin-bottom:8px">
        <strong>Uploads en cours</strong>
        <div style="display:flex; gap:8px; align-items:center">
          <small>{{ uploadsSummary }}</small>
          <button class="dv-btn small" @click="cancelAllUploads">Annuler tout</button>
        </div>
      </div>

      <!-- overall progress -->
      <div class="progress-wrap" style="margin-bottom:8px">
        <div class="progress-track">
          <div class="progress-fill" :style="{ width: overallProgress + '%' }"></div>
        </div>
        <div style="font-size:0.85rem; margin-top:6px; display:flex; justify-content:space-between">
          <span>Global</span>
          <span>{{ overallProgress }}%</span>
        </div>
      </div>

      <!-- per-file list -->
      <ul style="list-style:none; padding:0; margin:0; display:flex; flex-direction:column; gap:8px">
        <li v-for="(u, i) in uploadQueue" :key="u.id" style="display:flex; gap:12px; align-items:center;">
          <div style="flex:1">
            <div style="display:flex; justify-content:space-between; align-items:center">
              <div style="font-weight:600; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; max-width:60%">
                {{ u.name }}
              </div>
              <div style="font-size:0.85rem; color:#666">{{ u.progress }}% • {{ humanSize(u.size) }}</div>
            </div>

            <div class="progress-wrap" style="margin-top:6px">
              <div class="progress-track">
                <div class="progress-fill" :style="{ width: u.progress + '%' }"></div>
              </div>
            </div>

            <div style="display:flex; gap:8px; margin-top:6px; align-items:center">
              <small v-if="u.status === 'uploading'">Envoi…</small>
              <small v-else-if="u.status === 'done'" style="color:green">Terminé</small>
              <small v-else-if="u.status === 'error'" style="color:#c00">Erreur</small>
              <small v-else-if="u.status === 'cancelled'" style="color:#999">Annulé</small>
            </div>
          </div>

          <div style="display:flex; flex-direction:column; gap:6px; align-items:flex-end">
            <button v-if="u.status === 'uploading'" class="dv-mini" @click="cancelUpload(u.id)" title="Annuler">✖️</button>
            <button v-else class="dv-mini" @click="retryUpload(u.id)" title="Réessayer">↻</button>
          </div>
        </li>
      </ul>
    </div>

    <div class="modal-actions">
      <button class="dv-btn" @click="showUploadModal = false">Fermer</button>
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

    <!-- NEW FILE MODAL -->
    <div v-if="showNewFileModal" class="modal-wrap" @keydown.esc="showNewFileModal = false">
      <div class="modal" role="dialog" aria-modal="true" aria-label="Nouveau fichier">
        <h3>Nouveau fichier</h3>
        <div class="form-group">
          <label>Type de fichier</label>
          <select v-model="newFileType" style="width:100%; padding:8px; border-radius:6px; border:1px solid #e6e6ee; margin-bottom:12px;">
            <option value="docx">📄 Document Word (.docx)</option>
            <option value="xlsx">📊 Feuille Excel (.xlsx)</option>
            <option value="pptx">📽️ Présentation PowerPoint (.pptx)</option>
          </select>
        </div>
        <input v-model="newFileName" placeholder="Nom du fichier (sans extension)" class="compose-input" @keyup.enter="createFile" />
        <div class="modal-actions" style="justify-content:space-between;">
          <div>
            <button class="dv-btn ghost" @click="showNewFileModal = false">Annuler</button>
          </div>
          <div style="display:flex; gap:8px;">
            <button class="dv-btn primary" @click="createFile">Créer et ouvrir</button>
          </div>
        </div>
      </div>
      <div class="backdrop" @click="showNewFileModal = false"></div>
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
    <div v-if="showShareModal" class="modal-wrap" @keydown.esc="closeShareModal">
      <div class="modal" role="dialog" aria-modal="true">
        <h3>Partager "{{ selectedFile?.name }}"</h3>
        
        <div class="share-form">
          <div class="form-group">
            <label>Partager avec (nom d'utilisateur)</label>
            <input v-model="shareWithUsername" placeholder="Nom d'utilisateur" />
          </div>
          
          <div class="form-group">
            <label>Permission</label>
            <select v-model="sharePermission">
              <option value="editor">Éditeur</option>
              <option value="viewer">Lecteur seul</option>
            </select>
          </div>
        </div>

        <div class="shared-with-list" v-if="currentShares.length">
          <h4>Partagé avec</h4>
          <ul>
            <li v-for="share in currentShares" :key="share.id">
              {{ share.shared_with_name }}
              ({{ share.permission }})
              <button class="dv-mini danger" @click="unshareWith(share)">✖️</button>
            </li>
          </ul>
        </div>

        <div class="modal-actions">
          <button class="dv-btn" @click="closeShareModal">Fermer</button>
          <button class="dv-btn primary" @click="shareWithUser">Partager</button>
        </div>
      </div>
      <div class="backdrop" @click="closeShareModal"></div>
    </div>

    <!-- ONLYOFFICE MODAL -->
    <div v-if="showOnlyofficeModal" class="modal-wrap" @keydown.esc="showOnlyofficeModal = false">
      <div class="modal onlyoffice-modal" :class="{ 'fullscreen': isOnlyofficeFullscreen }" role="dialog" aria-modal="true">
        <div class="onlyoffice-header">
          <h3 style="margin:0; font-size:1rem;">{{ selectedFile?.name }}</h3>
          <div style="display:flex; gap:8px;">
            <button class="dv-btn small" @click="toggleOnlyofficeFullscreen" :title="isOnlyofficeFullscreen ? 'Quitter plein écran' : 'Plein écran'">
              {{ isOnlyofficeFullscreen ? '⛶ Quitter' : '⛶ Plein écran' }}
            </button>
            <button class="dv-btn" @click="showOnlyofficeModal = false">Fermer</button>
          </div>
        </div>
        <div id="onlyofficeModalContainer" class="onlyoffice-container" tabindex="0" @click="focusOnlyOfficeEditor"></div>
      </div>
      <div class="backdrop" @click="showOnlyofficeModal = false"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { gls } from '@/stores/global.js'


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
const isMobile = ref(window.innerWidth <= 750)
const searchQuery = ref('')
const selectedFile = ref(null)

// ---------- Storage info ----------
const storageUsed = ref(0)
const storageLimit = ref(0)
const storagePercentage = computed(() => {
  if (storageLimit.value === 0) return 0
  return Math.min(100, Math.round((storageUsed.value / storageLimit.value) * 100))
})
const storageColor = computed(() => {
  const pct = storagePercentage.value
  if (pct >= 90) return '#ef4444' // rouge
  if (pct >= 75) return '#f59e0b' // orange
  return '#10b981' // vert
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

let _onlyofficeScriptLoading = null
function loadOnlyOfficeScript() {
  if (window.DocEditor || (window.DocsAPI && window.DocsAPI.DocEditor)) return Promise.resolve()
  if (_onlyofficeScriptLoading) return _onlyofficeScriptLoading

  _onlyofficeScriptLoading = new Promise((resolve, reject) => {
    try {
      // determine base: either env VITE_API_ONLYOFFICE (can be full api.js URL or base) or fallback to localhost:8082
      let cfg = (import.meta.env.VITE_API_ONLYOFFICE || '').trim()
      let docserverBase = ''
      if (cfg) {
        // strip possible trailing /web-apps/.../api.js if user provided full path
        docserverBase = cfg.replace(/\/web-apps\/apps\/api\/documents\/api\.js\/?$/i, '').replace(/\/+$/, '')
      } else {
        docserverBase = `${window.location.protocol}//${window.location.hostname}:8082`
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


<style scoped>
/* Drive Hero Card */
.drive-hero-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 24px;
  background: rgba(245, 245, 247, 0.85);
  border-radius: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  margin: 16px 16px 0 16px;
  gap: 12px;
  position: relative;
  overflow: hidden;
}

.dark .drive-hero-card {
  background: #1C1C1E;
}

.drive-hero-card .hero-icon-wrapper {
  width: 70px;
  height: 70px;
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: pulse 2s ease-in-out infinite;
  box-shadow: 0 6px 24px rgba(0, 0, 255, 0.3);
}

@keyframes pulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.06); }
}

.drive-hero-card .hero-icon {
  width: 36px;
  height: 36px;
  color: white;
}

.drive-hero-card h1 {
  font-family: 'Roboto', sans-serif;
  font-size: 1.8rem;
  font-weight: 700;
  text-align: center;
  letter-spacing: 1px;
  color: #222;
  margin: 0;
}

.dark .drive-hero-card h1 {
  color: white;
}

.drive-hero-card .hero-subtitle {
  font-family: 'Roboto', sans-serif;
  font-size: 1rem;
  font-style: italic;
  text-align: center;
  background: -webkit-linear-gradient(30deg, blue, red);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 0.5px;
  margin: 0;
}

#drive-v2 {
  display: flex;
  flex-direction: column;
  height: 100vh;
  font-family: "Roboto", system-ui, -apple-system, "Segoe UI", Roboto, "Helvetica Neue", Arial;
  background: #f6f7fb;
  color: #111;
}
/* DARK MODE */
.dark #drive-v2 { background: #121212; color: #e5e5e5; }

/* Drive Content Wrapper */
.drive-content-wrapper {
  display: flex;
  flex: 1;
  overflow: hidden;
}

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
.dv-btn.ghost { background:grey; }
.dv-btn.small { padding:6px 8px; font-size:0.9rem; }
.dv-btn.small.active { padding:6px 8px; font-size:0.9rem; background: lightblue; color:#fff; }
.dv-nav { display:flex; flex-direction:column; gap:8px; margin-top:6px; }
.nav-item { background:transparent; border:none; padding:10px; text-align:left; border-radius:8px; cursor:pointer; color:#333; }
.dark .nav-item { color:#e5e5e5; }
.nav-item.active {  background: -webkit-linear-gradient(30deg, blue, red);; color:#fff; }

/* Storage Info */
.storage-info {
  margin-top: auto;
  margin-bottom: 8px;
  padding: 12px;
  background: #f9fafb;
  border-radius: 12px;
  border: 1px solid #e5e7eb;
}
.dark .storage-info {
  background: #1a1a1a;
  border-color: #2b2b2b;
}
.storage-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  font-weight: 600;
  font-size: 0.9rem;
}
.storage-icon {
  font-size: 1.2rem;
}
.storage-bar {
  height: 8px;
  background: #e5e7eb;
  border-radius: 999px;
  overflow: hidden;
  margin-bottom: 8px;
}
.dark .storage-bar {
  background: #2b2b2b;
}
.storage-fill {
  height: 100%;
  transition: width 0.3s ease, background 0.3s ease;
  border-radius: 999px;
}
.storage-details {
  display: flex;
  justify-content: space-between;
  font-size: 0.85rem;
  color: #666;
}
.dark .storage-details {
  color: #bbb;
}
.storage-percent {
  font-weight: 600;
}
.storage-warning {
  margin-top: 8px;
  padding: 6px 10px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: 8px;
  font-size: 0.8rem;
  color: #dc2626;
  text-align: center;
}
.dark .storage-warning {
  background: #3f1f1f;
  border-color: #7f1d1d;
  color: #fca5a5;
}

/* MAIN */
.dv-main { flex:1; padding:16px; display:flex; flex-direction:column; gap:12px; }
.dv-header { display:flex; justify-content:space-between; align-items:center; }
.dv-breadcrumbs { display:flex; align-items:center; gap:8px; }
.crumb { background:transparent; border:none; color:#0066cc; cursor:pointer; padding:6px 8px; border-radius:6px; }
.crumb-target { background:#e8f4ff; }
.slash { color:#888 }
.current { font-weight:600; }
.dv-tools { display:flex; gap:8px; align-items:center; }
.dv-search { padding:8px 10px; border-radius:8px; border:1px solid #e0e0e6; }
/* BODY SPLIT */
.dv-body { display:grid; grid-template-columns: 420px 1fr; gap:20px; height: calc(100% - 60px); }
/* LIST */
.dv-list { background:#fff; border-radius:8px; padding:12px; overflow:auto; border:1px solid #eaeaf2; min-height:200px; transition: border-color .08s, background .08s; }
.dark .dv-list { background:#151515; border-color:#2b2b2b; }
.dv-list.drag-over { background:#f3fff6; border-color:#bff0c5; }
.dv-empty { display:flex; flex-direction:column; align-items:center; justify-content:center; height:100%; padding:24px; text-align:center; color:#666; }
.dv-empty-illustration { margin-bottom:10px; }
.dv-empty-text { font-size:1.0rem; margin:6px 0; }
.dv-empty-sub { font-size:0.9rem; color:#999; }
.dv-items { list-style:none; margin:0; padding:0; display:flex; flex-direction:column; gap:8px; }
.dv-item { display:flex; justify-content:space-between; align-items:center; padding:10px; border-radius:8px; cursor:pointer; transition: background .08s; }
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

/* Highlight folder target */
.folder-hover { background: #d9f7e8 !important; }

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
.modal { width:420px; background:#fff; padding:50px; border-radius:12px; box-shadow:0 10px 40px rgba(0,0,0,0.12); z-index:2001; }
.dark .modal { background:#1e1e1e; color:#e5e5e5; }
.backdrop { position:fixed; inset:0; background:rgba(0,0,0,0.25); z-index:2000; }
.upload-zone { border:2px dashed #e6e6ee; border-radius:8px; padding:18px; text-align:center; cursor:pointer; margin:10px 0; }
.modal-actions { display:flex; gap:8px; justify-content:flex-end; margin-top:12px; }
.modal input, .modal textarea, .modal .compose-input { width:100%; padding:10px; border-radius:8px; border:1px solid #e6e6ee; }
.dark .modal input, .dark .modal .compose-input { background:#121212; border:1px solid #333; color:#eee; }

/* OnlyOffice modal styles */
.modal.onlyoffice-modal {
  width: 90vw;
  max-width: 1200px;
  height: 90vh;
  padding: 12px;
  display: flex;
  flex-direction: column;
  transition: all 0.2s ease;
}

.modal.onlyoffice-modal.fullscreen {
  width: 100vw;
  height: 100vh;
  max-width: 100vw;
  border-radius: 0;
  padding: 8px;
  transition: all 0.2s ease;
}

.onlyoffice-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e6e6ee;
}

.dark .onlyoffice-header {
  border-bottom-color: #2b2b2b;
}

.onlyoffice-container {
  width: 100%;
  height: calc(100% - 48px);
  min-height: 60vh;
  outline: none;
  /* Permettre le focus sur le conteneur */
}
.onlyoffice-container:focus {
  outline: none;
}

.modal.onlyoffice-modal.fullscreen .onlyoffice-container {
  height: calc(100vh - 56px);
}

/* responsive */
@media (max-width: 900px) {
  .dv-body { grid-template-columns: 1fr; }
  .dv-sidebar { display:none; }
  .info-panel { width:100%; border-left:none; padding-left:0; }
}

/* up target highlight */
.up-target { background: #f2fff0; border-radius:8px; padding:6px; }
.crumb-target { outline: 2px dashed rgba(0,110,255,0.15); border-radius:6px; }

/* small extras */
.dv-up { margin-bottom: 8px; display:flex; justify-content:flex-start; }


/* when OnlyOffice editor is open we enlarge the preview area */
.preview.office-open {
  min-height: calc(100vh - 140px);
  height: calc(100vh - 140px);
  padding: 0;
  margin: 0;
  overflow: hidden;
}

/* editor container should fill the preview */
#onlyofficeContainer {
  width: 100%;
  height: 100%;
  min-height: 60vh;

}

.dark .preview{ background:#1b1b1b; border-color:#2b2b2b; }
</style>