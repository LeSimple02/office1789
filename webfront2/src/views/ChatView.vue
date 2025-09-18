<template>
  <div class="chat-container">
    <div class="main-layout">
      <!-- Affichage normal (desktop) ou liste des conversations (mobile) -->
      <aside class="sidebar" v-if="!isMobileView || !selectedDiscussion">
        <button class="btn btn-primary" @click="showCreateModal = true">
          <span type="icon">＋</span>{{$t("newG")}}
        </button>
        <h2 class="subtitle">Messages</h2>
        <div class="sidebar-list">
          <div v-for="user in users" :key="user.id" :class="['sidebar-card', selectedDiscussion && selectedDiscussion.id === user.id && !selectedDiscussion.isGroup ? 'active' : '']" @click="selectDiscussion(user, false)">
            <img class="sidebar-avatar" :src="user.avatar" alt="avatar" />
            <span class="sidebar-name">{{ user.name }}</span>
          </div>
        </div>
        <h2 class="subtitle">Groupes</h2>
        <div class="sidebar-list">
          <div v-for="group in groups" :key="group.id" :class="['sidebar-card', selectedDiscussion && selectedDiscussion.id === group.id && selectedDiscussion.isGroup ? 'active' : '']" @click="selectDiscussion(group, true)">
            <img class="sidebar-avatar" :src="group.avatar" alt="avatar" />
            <span class="sidebar-name">{{ group.name }}</span>
            <button class="btn btn-danger" @click.stop="openDeleteModal(group)"><span class="icon">🗑️</span></button>
          </div>
        </div>
      </aside>

      <!-- Affichage des messages (toujours visible en desktop, plein écran en mobile) -->
      <section class="messages-panel" :class="{ 'mobile-fullscreen': isMobileView }">
        <div v-if="selectedDiscussion" class="messages-modal">
          <div class="conv-header">
            <button v-if="showBackButton" class="btn" @click="backToList">
              <span class="icon">←</span>
            </button>
            <h2>{{ messagesModalTitle }}</h2>
            <div class="conv-actions">
              <button class="btn" @click="call('audio')">
                <span class="icon">📞</span>
              </button>
              <button class="btn" @click="call('video')">
                <span class="icon">🎥</span>
              </button>
              <button v-if="selectedDiscussion && selectedDiscussion.isGroup" class="btn" @click="showGroupUsers = true">
                <span class="icon">👥</span>
              </button>
              <div class="dropdown">
                <button class="btn" @click="toggleMenu">
                  <span class="icon">⚙️</span>
                </button>
                <div v-if="showMenu" class="dropdown-menu">
                  <button @click="blockUserOrGroup">Bloquer</button>
                  <button @click="reportUserOrGroup">Signaler</button>
                </div>
              </div>
            </div>
            <!-- Modal liste des utilisateurs du groupe -->
            <div v-if="showGroupUsers" class="modal-overlay">
              <div class="modal">
                <h3>Utilisateurs du groupe</h3>
                <ul class="group-users-list">
                  <li v-for="user in groupUsers" :key="user.id" class="group-user-item">
                    <span>{{ user.name }}</span>
                    <button v-if="isAdmin" class="btn btn-danger" @click="removeUserFromGroup(user)">Supprimer</button>
                  </li>
                </ul>
                <div class="modal-actions">
                  <button class="btn" @click="showGroupUsers = false">Fermer</button>
                </div>
              </div>
            </div>
          </div>
          <div :class="{messages_list: 1}">
            <div v-for="msg in currentMessages" :key="msg.id" :class="['message', msg.fromMe ? 'me' : 'other']">
              <span class="message-content">{{ msg.text }}</span>
              <span class="message-date">{{ msg.date }}</span>
            </div>
          </div>
          <div class="modal-actions">
          <label class="attach-label">
                        <span class="attach-btn">📎</span>
                        <input type="file" class="attach-input" style="display:none;" />
                    </label>
            <input v-model="newMessageText" class="input" placeholder="Écrire un message..." @keyup.enter="sendMessage" />
            <button class="btn btn-primary" @click="sendMessage">Envoyer</button>
          </div>
        </div>
      </section>
    </div>
    <!-- Modal Création -->
    <div v-if="showCreateModal" class="modal-overlay">
      <div class="modal">
        <h2>{{$t("newG")}}</h2>
        <input v-model="newGroupName" class="input" :placeholder="$t('username')" /><button class="btn btn-primary">
        +
      </button>
        <div class="modal-actions">
          <button class="btn btn-primary" @click="createGroup">Créer</button>
          <button class="btn" @click="showCreateModal = false">Annuler</button>
        </div>
      </div>
    </div>
    <!-- Modal Suppression -->
    <div v-if="showDeleteModal" class="modal-overlay">
      <div class="modal">
        <h2>Supprimer le groupe</h2>
        <p>Voulez-vous vraiment supprimer <b>{{ groupToDelete?.name }}</b> ?</p>
        <div class="modal-actions">
          <button class="btn btn-danger" @click="deleteGroup">Supprimer</button>
          <button class="btn" @click="showDeleteModal = false">Annuler</button>
        </div>
      </div>
    </div>
    <!-- Modal Call -->
    <div v-if="showCallModal" class="modal-overlay">
      <div class="modal call-modal">
        <img
          v-if="selectedDiscussion"
          class="call-avatar"
          :src="selectedDiscussion.avatar"
          alt="avatar"
        />
        <h2 v-if="selectedDiscussion">{{ selectedDiscussion.name }}</h2>
        <p>Appel {{ callType === 'audio' ? 'audio' : 'vidéo' }} en cours...</p>
        <div class="modal-actions">
          <button class="btn btn-danger" @click="declineCall">
            ❌ Raccrocher
          </button>
          <button class="btn btn-primary" @click="acceptCall">
            ✅ Décrocher
          </button>
        </div>
      </div>
    </div>
    <!-- Modal In-Call -->
    <div v-if="showInCallModal" class="modal-overlay">
      <div class="modal in-call-modal">
        <h2>En communication avec {{ selectedDiscussion?.name }}</h2>

        <!-- Vidéo si activée -->
        <video
          v-if="callType === 'video'"
          ref="videoRef"
          autoplay
          playsinline
          muted
          class="call-video"
        ></video>
        <p v-else>🎤 Appel audio en cours...</p>
        <div class="modal-actions">
          <button class="btn btn-danger" @click="endCall">
            🔴 Raccrocher
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';

const showInCallModal = ref(false);
const localStream = ref(null);
const videoRef = ref(null);
const showCallModal = ref(false);
const callType = ref(null);
const showGroupUsers = ref(false);
const isAdmin = ref(true);
const groupUsers = ref([
  { id: 1, name: 'Tk78' },
  { id: 2, name: 'Bob' },
  { id: 3, name: 'Charlie' }
]);
const isMobileView = ref(false);
const showBackButton = ref(false);

function removeUserFromGroup(user) {
  groupUsers.value = groupUsers.value.filter(u => u.id !== user.id);
  alert(user.name + ' supprimé du groupe !');
}

const groups = ref([
  { id: 1, name: 'Général', avatar: 'https://randomuser.me/api/portraits/lego/1.jpg' },
  { id: 2, name: 'Projet X', avatar: 'https://randomuser.me/api/portraits/lego/2.jpg' },
]);
const users = ref([
  { id: 1, name: 'Tk78', avatar: 'https://i.ytimg.com/vi/970_h2pxfo8/maxresdefault.jpg' },
  { id: 2, name: 'Bob', avatar: 'https://randomuser.me/api/portraits/men/2.jpg' },
]);
const showCreateModal = ref(false);
const showDeleteModal = ref(false);
const newGroupName = ref('');
const newGroupAvatar = ref('');
const groupToDelete = ref(null);
const selectedDiscussion = ref(null);
const currentMessages = ref([]);
const messagesModalTitle = ref('');
const newMessageText = ref('');
const messagesData = ref({
  user_1: [
    { id: 1, text: 'Salut Tk78 !', fromMe: true, date: '10:01' },
    { id: 2, text: 'Coucou, ça va ?', fromMe: false, date: '10:02' },
  ],
  user_2: [
    { id: 1, text: 'Hey Bob !', fromMe: true, date: '09:55' },
    { id: 2, text: 'Yo !', fromMe: false, date: '09:56' },
  ],
  group_1: [
    { id: 1, text: 'Bienvenue dans le groupe Général !', fromMe: false, date: 'Hier' },
    { id: 2, text: 'Hello tout le monde', fromMe: true, date: 'Aujourd’hui' },
  ],
  group_2: [
    { id: 1, text: 'Projet X commence !', fromMe: false, date: 'Hier' },
  ],
});

function selectDiscussion(item, isGroup) {
  selectedDiscussion.value = { ...item, isGroup };
  let key = isGroup ? `group_${item.id}` : `user_${item.id}`;
  currentMessages.value = messagesData.value[key] || [];
  messagesModalTitle.value = isGroup ? `Groupe : ${item.name}` : `${item.name}`;
  newMessageText.value = '';

  if (window.innerWidth <= 600) {
    isMobileView.value = true;
    showBackButton.value = true;
  }
}

function backToList() {
  isMobileView.value = false;
  showBackButton.value = false;
  selectedDiscussion.value = null;
}

function createGroup() {
  if (newGroupName.value.trim()) {
    groups.value.push({
      id: Date.now(),
      name: newGroupName.value.trim(),
      avatar: newGroupAvatar.value.trim() || 'https://randomuser.me/api/portraits/lego/3.jpg'
    });
    newGroupName.value = '';
    newGroupAvatar.value = '';
    showCreateModal.value = false;
  }
}

function openDeleteModal(group) {
  groupToDelete.value = group;
  showDeleteModal.value = true;
}

function deleteGroup() {
  groups.value = groups.value.filter(g => g.id !== groupToDelete.value.id);
  groupToDelete.value = null;
  showDeleteModal.value = false;
}


var messages_list = ref('false')
function sendMessage() {
  if (!newMessageText.value.trim() || !selectedDiscussion.value) return;
  let key = selectedDiscussion.value.isGroup ? `group_${selectedDiscussion.value.id}` : `user_${selectedDiscussion.value.id}`;
  const msg = {
    id: Date.now(),
    text: newMessageText.value,
    fromMe: true,
    date: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  };
  //messages_list.value.scrollTop = messages_list.value.scrollHeight
  if (!messagesData.value[key]) messagesData.value[key] = [];
  messagesData.value[key].push(msg);
  newMessageText.value = '';
}

const showMenu = ref(false);
function toggleMenu() {
  showMenu.value = !showMenu.value;
}

const ringtone = new Audio("/src/assets/call.mp3");

function call(type) {
  callType.value = type;
  showCallModal.value = true;
  ringtone.loop = true;
  ringtone.play();
}

function declineCall() {
  showCallModal.value = false;
  callType.value = null;
  ringtone.loop = false;
  ringtone.pause();
}

async function acceptCall() {
  showCallModal.value = false;
  showInCallModal.value = true;
  ringtone.loop = false;
  ringtone.pause();
  try {
    const stream = await navigator.mediaDevices.getUserMedia({
      video: callType.value === 'video',
      audio: true
    });
    localStream.value = stream;
    if (videoRef.value) {
      videoRef.value.srcObject = stream;
    }
  } catch (err) {
    alert("Impossible d’accéder à la caméra/micro : " + err.message);
  }
}

function blockUserOrGroup() {
  alert('Utilisateur/Groupe bloqué !');
  showMenu.value = false;
}

function reportUserOrGroup() {
  alert('Utilisateur/Groupe signalé !');
  showMenu.value = false;
}

function endCall() {
  if (localStream.value) {
    localStream.value.getTracks().forEach(track => track.stop());
    localStream.value = null;
  }
  showInCallModal.value = false;
  callType.value = null;
}

function handleResize() {
  if (window.innerWidth <= 600 && selectedDiscussion.value) {
    isMobileView.value = true;
    showBackButton.value = true;
  } else {
    isMobileView.value = false;
    showBackButton.value = false;
  }
}

onMounted(() => {
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
:root {
  --bg: #fafafa;
  --card: #fff;
  --border: #dbdbdb;
  --primary: #3897f0;
  --danger: #ed4956;
  --text: #262626;
  --muted: #8e8e8e;
  --shadow: 0 2px 16px rgba(0,0,0,0.08);
}
.chat-container {
  min-height: 100vh;
  width: 100%;
  background: var(--bg);
  font-family: 'Segoe UI', 'Roboto', Arial, sans-serif;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 0 18px 0;
  max-width: 500px;
  margin: 0 auto;
}
.title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text);
  letter-spacing: 1px;
}
.btn {
  border: none;
  border-radius: 999px;
  padding: 10px 22px;
  font-size: 1rem;
  background: var(--card);
  color: var(--text);
  cursor: pointer;
  box-shadow: 0 1px 8px rgba(0,0,0,0.07);
  transition: background 0.2s, color 0.2s, box-shadow 0.2s;
  margin-left: 8px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.btn-primary {
  background: -webkit-linear-gradient(30deg, blue, red);
  color: #fff;
  font-weight: 600;
  box-shadow: 0 2px 12px rgba(0,0,0,0.12);
}
.btn-primary:hover {
  filter: brightness(1.08);
}
.dark .btn-danger {
  background: var(--danger);
  border: 1px solid red;
  color: #fff;
  font-weight: 600;
}
.btn-danger{
  color: black;
}
.btn-danger:hover {
  background: #b71c1c;
  color: white;
}
.subtitle {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--muted);
  margin: 18px 0 8px 8px;
}
.main-layout {
  display: flex;
  margin: 32px auto 0 auto;
  background: var(--card);
  border-radius: 24px;
  box-shadow: var(--shadow);
  min-height: 600px;
}
.sidebar {
  width: 320px;
  border-right: 1px solid var(--border);
  border-radius: 24px 0 0 24px;
  padding: 24px 0 0 0;
  display: flex;
  flex-direction: column;
  gap: 18px;
}
.sidebar-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 0 12px;
}
.sidebar-card {
  border: 1px solid var(--border);
  border-radius: 18px;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 14px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
  cursor: pointer;
  position: relative;
  transition: box-shadow 0.2s, background 0.2s;
}
.sidebar-card.active {
  background: -webkit-linear-gradient(30deg, blue, red);
  color: #fff;
}
.sidebar-card.active .sidebar-name {
  color: #fff;
}
.sidebar-card:hover {
  box-shadow: 0 8px 32px rgba(0,0,0,0.12);
}
.sidebar-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #d6249f;
  background: #eee;
  box-shadow: 0 2px 8px rgba(214,36,159,0.08);
}
.sidebar-name {
  font-size: 1rem;
  font-weight: 500;
  color: var(--text);
}
.sidebar-card .btn-danger {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 10px;
  font-size: 0.9em;
}
.messages-panel {
  flex: 1;
  padding: 32px 32px 0 32px;
  display: flex;
  flex-direction: column;
  align-items: stretch;
}
.messages-modal {
  margin: 0 auto;
  padding: 0;
  background: none;
  border-radius: 18px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.06);
  min-height: 500px;
  width: 100%;
  display: flex;
  flex-direction: column;
  max-height: 500px;
  overflow: hidden;
}
.messages_list {
  overflow: scroll;
  margin-bottom: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px 18px 0 18px;
  min-height: 40%;
}
.message {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  box-shadow: none;
  border-radius: 16px;
  padding: 7px 13px;
  box-shadow: 0 1px 6px rgba(41,121,255,0.04);
  font-size: 0.97em;
  max-width: 400px;
  width: fit-content;
  border: none;
}
.message.me {
  align-self: flex-end;
  background: #2979ff;
  color: #fff;
  box-shadow: none;
}
.message-content {
  margin-bottom: 1px;
}
.message-date {
  font-size: 0.78em;
  color: var(--muted);
  align-self: flex-end;
}
.messages-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--muted);
  font-size: 1.2em;
  gap: 12px;
}
.input {
  width: 90%;
  margin: 14px auto;
  padding: 12px 16px;
  border-radius: 5px;
  border: 1px solid var(--border);
  font-size: 1rem;
  background: #f7f7f7;
  color: var(--text);
  outline: none;
  transition: border 0.2s;
  display: block;
}
.input:focus {
  border-color: #d6249f;
}
.modal-overlay {
  position: fixed;
  font-family: roboto;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(38,38,38,0.18);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal {
  background: white;
  width: 25%;
  border-radius: 24px;
  box-shadow: 0 8px 32px rgba(214,36,159,0.18);
  padding: 40px 36px 32px 36px;
  text-align: center;
  align-items: center;
  justify-content: center;
}
.dark .modal{
   background: black;
}
.modal-actions {
  display: flex;
  justify-content: center;
  align-items: center;
  bottom : 0%;
  position: relative;
  gap: 12px;
  margin: 18px 0 8px 0;
}
.dark .modal-actions input{
  color: black;
}
.icon {
  font-size: 1.3em;
}
.conv-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 10px;
}
.conv-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}
.dropdown {
  position: relative;
}
.dropdown-menu {
  position: absolute;
  top: 36px;
  right: 0;
  background: var(--card);
  border: 1px solid var(--border);
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  padding: 8px 0;
  z-index: 10;
  min-width: 120px;
  display: flex;
  flex-direction: column;
}
.dropdown-menu button {
  background: none;
  border: none;
  padding: 8px 16px;
  text-align: left;
  color: var(--text);
  cursor: pointer;
  font-size: 1em;
}
.dropdown-menu button:hover {
  background: #f3f6fa;
}
@media (max-width: 600px) {
  .conv-header{
    background: black;
  }
  .modal-overlay{
    position: absolute;
    background: black;
    width: 100%;
    height: 100%;
    
  }
  .call-modal {
 
  width: 100%;
  height: 100%;
  position: absolute;
  background: rgb(100,100,100);
  button{
    margin-top: -50px;
  }
}
  .modal-actions{
    position: absolute;
    align-items: left;
    top:100%;
    
    
    margin-top: -50px;
    height: 30px;
    width: 100%;
    .input{
      border-radius: 12px;
      height:30px;
      width: 60%;
    }
    button{
      border-radius: 0;
      height:30px;
      width: %;
    }
  }
  .header, .discussion-list {
    max-width: 98vw;
    padding-left: 2vw;
    padding-right: 2vw;
  }
  .discussion-card {
    padding: 12px 8px;
  }
  .discussion-avatar {
    width: 36px;
    height: 36px;
  }
  .main-layout {
    flex-direction: column;
  }
  .sidebar {
    width: 100%;
    border-right: none;
    border-radius: 24px 24px 0 0;
  }
  .messages-panel {
    width: 100%;
    padding: 0;
    border-radius: 0;
  }
  .messages-panel.mobile-fullscreen {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 100;
    background: var(--card);
  }
  .conv-header {
    padding: 12px 16px;
    border-bottom: 1px solid var(--border);
  }
}
.call-modal {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.call-avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  margin-bottom: 16px;
  border: 3px solid #3897f0;
}
.in-call-modal {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.call-video {
  width: 100%;
  max-width: 500px;
  border-radius: 16px;
  margin: 16px 0;
  background: black;
}
.group-users-list {
  list-style: none;
  padding: 0;
  margin: 18px 0;
}
.group-user-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid #e3e6ee;
}
.group-user-item span {
  flex: 1;
}
.group-user-item button {
  margin-left: 12px;
}
</style>
