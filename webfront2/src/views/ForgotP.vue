<script setup>
import {ref} from "vue"

let id = ref("")
let sendm= ref(false)

function send(){
	if(id.value != "")
		sendm.value = true
	else{
		sendm.value = 'nothing'
	}
}
</script>

<template>
  <div id="forgot-bg">
    <div id="formem">
      <h1 class="title">{{$t('forgotp')}}</h1>
      <input v-model="id" name="fpassw" type="text" placeholder="Identifiant ou email" required />
      <button id="bform" @click="send">Submit</button>
      <p v-if="sendm=='nothing'" class="error">Enter your id</p>
      <p v-if="sendm==true" class="success">Un email vient de vous être envoyé</p>
    </div>
  </div>
</template>

<style scoped>
/* Fond comme login/create */
#forgot-bg {
  min-height: 100vh;
  width: 100vw;
  background: linear-gradient(120deg, #e0e7ef 0%, #f5f5f7 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}
#forgot-bg::before {
  content: "";
  position: absolute;
  inset: 0;
  background: rgba(240,240,245,0.7);
  backdrop-filter: blur(16px);
  z-index: 0;
}

.dark #forgot-bg {
  background: linear-gradient(120deg, #23243a 0%, #1a1b26 100%);
}
.dark #forgot-bg::before {
  background: rgba(30,30,40,0.7);
  backdrop-filter: blur(16px);
}

/* Carte centrée */
#formem {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 18px;
  width: 100%;
  max-width: 400px;
  padding: 48px 32px 32px 32px;
  border-radius: 24px;
  background: rgba(255,255,255,0.85);
  box-shadow: 0 8px 32px rgba(0,0,0,0.18);
  z-index: 1;
  margin: 0 auto;
  text-align: center;
  font-family: roboto, arial, sans-serif;
}

.dark #formem {
  background: rgba(30,30,40,0.92);
  color: #eee;
}

/* Titre */
.title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #00308F;
  margin-bottom: 8px;
}
.dark .title {
  color: #fff;
}

/* Champ input */
#formem input {
  height: 44px;
  width: 100%;
  border-radius: 12px;
  border: 1px solid #d0d0d0;
  padding: 0 16px;
  font-size: 1rem;
  background: rgba(245,245,247,0.95);
  transition: border 0.2s, box-shadow 0.2s;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}
#formem input:focus {
  border: 1.5px solid #00308F;
  box-shadow: 0 4px 16px rgba(0,48,143,0.10);
  outline: none;
}
.dark #formem input {
  background: rgba(30,30,40,0.95);
  color: #eee;
  border: 1px solid #444;
}

/* Bouton */
#bform {
  width: 100%;
  height: 44px;
  background: -webkit-linear-gradient(30deg, #00308F, #ff3c3c);
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 1.1rem;
  font-weight: 700;
  letter-spacing: 1px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.10);
  cursor: pointer;
  transition: background 0.3s, box-shadow 0.2s;
}
#bform:hover {
  background: -webkit-linear-gradient(30deg, #ff3c3c, #00308F);
  box-shadow: 0 4px 16px rgba(0,0,0,0.18);
}

/* Messages */
.error {
  color: #ff3c3c;
  font-weight: 500;
  font-size: 0.95rem;
}
.success {
  color: green;
  font-weight: 500;
  font-size: 0.95rem;
}

/* Responsive */
@media (max-width: 600px) {
  #formem {
    max-width: 95vw;
    padding: 24px 8px 16px 8px;
    border-radius: 12px;
  }
  .title {
    font-size: 1.2rem;
  }
}
</style>
