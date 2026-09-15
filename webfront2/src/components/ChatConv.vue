<script setup>
import {gls} from "@/stores/global.js"
import {ref} from "vue"
import CallAc from "@/components/CallAc.vue"

let conv = ref([])
let convopt = ref([])
let url = ref(window.location.href.split("/")[4])

fetch(process.env.VUE_APP_API_INFO_CONV, {method:"POST", mode:"cors", body : JSON.stringify({"username" : gls().username, "token" : gls().sessionT}) }).then(a=>a.text()).then(a=>{
  if(a != ""){
    conv.value = a.split(",")
    for (let mess of conv.value){
      convopt.value.push({"user": mess, "link": "/chat/"+mess})
    }
  }
})

let callA = ref(0)

function call(){
  callA.value = 1
  let sound = new Audio(require("@/assets/call.mp3"))
  sound.play()
}
</script>

<template>
  <div class="conv">
    <CallAc v-if="callA == 1" />
    <div class="conv-head">
      <h3>Conversations</h3>
      <button @click="call" class="call-btn">Appeler</button>
    </div>
    <div class="conv-list">
      <RouterLink v-for="c in convopt" :key="c.user" :to="c.link" class="conv-item">
        <div class="conv-avatar">{{ c.user.charAt(0).toUpperCase() }}</div>
        <span class="conv-name">{{ c.user }}</span>
      </RouterLink>
      <div v-if="convopt.length === 0" class="empty">Aucune conversation</div>
    </div>
  </div>
</template>

<style scoped>
.conv { max-width: 480px; margin: 0 auto; padding: 2rem 1.25rem 4rem; animation: fadeInUp .5s var(--ease) both; }
.conv-head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 1.25rem; animation: fadeInDown .4s var(--ease) both; }
.conv-head h3 { font-size: 1.25rem; font-weight: 700; }
.call-btn { padding: .5rem 1.125rem; font-size: .8125rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); }
.call-btn:hover { background: var(--ink-soft); border-color: var(--ink-soft); transform: translateY(-1px); }
.conv-list { display: flex; flex-direction: column; gap: .5rem; }
.conv-item { display: flex; align-items: center; gap: .875rem; padding: .875rem 1rem; background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r); transition: all var(--tb); animation: fadeInLeft .3s var(--ease) both; }
.conv-item:hover { border-color: var(--border-strong); background: var(--bg-secondary); transform: translateX(4px); }
.conv-avatar { width: 36px; height: 36px; border-radius: 50%; background: var(--bg-tertiary); border: 1px solid var(--border); display: flex; align-items: center; justify-content: center; font-size: .875rem; font-weight: 700; flex-shrink: 0; }
.conv-name { font-size: .9375rem; font-weight: 500; color: var(--text-primary); }
.empty { text-align: center; padding: 2rem 1rem; font-size: .875rem; color: var(--text-muted); }
</style>
