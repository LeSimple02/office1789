<script setup>
import {ref} from "vue"

let id = ref("")
let sendm = ref(false)
let loading = ref(false)

async function send(){
  if(id.value == "") {
    sendm.value = 'nothing'
    return
  }
  loading.value = true
  sendm.value = false
  try {
    const response = await fetch(import.meta.env.VITE_APP_API_PASSWORD_RESET_REQUEST, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ identifier: id.value })
    })
    const data = await response.json()
    if (response.ok || data.success) {
      sendm.value = true
    } else {
      sendm.value = true
    }
  } catch (error) {
    console.error('Error:', error)
    sendm.value = true
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="page-wrap">
    <div class="card">
      <div class="card-head">
        <h1>{{ $t('forgotp') }}</h1>
        <p>{{ $t('recoverAccountAccess') }}</p>
      </div>
      <div class="form">
        <input v-model="id" type="text" :placeholder="$t('identifierOrEmail')" required />
        <button @click="send" class="btn-go" :disabled="loading">
          <span v-if="!loading">{{ $t('sendEmail') }}</span>
          <span v-else class="loader"></span>
        </button>
        <p v-if="sendm=='nothing'" class="msg err">{{ $t('pleaseEnterIdentifier') }}</p>
        <p v-if="sendm==true" class="msg ok">{{ $t('emailHasBeenSent') }}</p>
        <RouterLink to="/login" class="back">← {{ $t('backToLogin') }}</RouterLink>
      </div>
    </div>
    <div class="tip">{{ $t('forgotPasswordInfo') }}</div>
  </div>
</template>

<style scoped>
.page-wrap { max-width: 440px; margin: 0 auto; padding: 3rem 1.25rem 5rem; animation: fadeInUp .5s var(--ease) both; }
.card { background: var(--bg-primary); border: 1px solid var(--border); border-radius: var(--r-lg); box-shadow: var(--sh-sm); overflow: hidden; animation: fadeInUp .4s var(--ease) .1s both; }
.card-head { text-align: center; padding: 2rem 1.75rem 1rem; animation: fadeInDown .4s var(--ease) both; }
.card-head h1 { font-size: 1.5rem; font-weight: 800; letter-spacing: -.025em; margin-bottom: .375rem; }
.card-head p { font-size: .875rem; color: var(--text-secondary); }
.form { padding: 0 1.75rem 1.75rem; display: flex; flex-direction: column; gap: 1rem; }
.form input { width: 100%; padding: .6875rem .875rem; font-size: .9375rem; color: var(--text-primary); background: var(--bg-primary); border: 1px solid var(--border-strong); border-radius: var(--r); transition: border-color var(--tb), box-shadow var(--tb); }
.form input:focus { outline: none; border-color: var(--ink); box-shadow: 0 0 0 3px var(--bg-tertiary); }
.form input::placeholder { color: var(--text-muted); }
.btn-go { width: 100%; padding: .75rem; font-size: .9375rem; font-weight: 600; color: var(--paper); background: var(--ink); border: 1px solid var(--ink); border-radius: var(--r); cursor: pointer; transition: all var(--tb); display: flex; align-items: center; justify-content: center; min-height: 44px; }
.btn-go:hover:not(:disabled) { background: var(--ink-soft); border-color: var(--ink-soft); box-shadow: var(--sh-md); transform: translateY(-2px); }
.btn-go:disabled { opacity: .5; cursor: not-allowed; }
.msg { padding: .625rem .875rem; font-size: .8125rem; border-radius: var(--r); border: 1px solid var(--border); animation: fadeIn .2s var(--ease) both; }
.msg.err { background: var(--bg-tertiary); color: var(--text-primary); }
.msg.ok { background: var(--bg-tertiary); color: var(--text-primary); }
.back { font-size: .8125rem; font-weight: 500; color: var(--text-secondary); text-align: center; transition: color var(--tb); }
.back:hover { color: var(--text-primary); text-decoration: underline; text-underline-offset: 3px; }
.tip { margin-top: 1.25rem; padding: .875rem 1.125rem; font-size: .8125rem; color: var(--text-secondary); background: var(--bg-secondary); border: 1px solid var(--border); border-radius: var(--r); line-height: 1.5; animation: fadeInUp .4s var(--ease) .3s both; }
.loader { width: 18px; height: 18px; border: 2.5px solid var(--border); border-top-color: var(--paper); border-radius: 50%; animation: spin .7s linear infinite; display: inline-block; }
</style>
