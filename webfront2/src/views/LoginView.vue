<script setup>
import {gls} from "@/stores/global"
import {ref} from "vue"
import { useDark } from '@vueuse/core'
import router from "@/router/index"


let passw = ref("password")

let userl = ref('')
let passl = ref('')
let wrong = ref(0)

const isDark = useDark()


function connect(){

	
	
	fetch(import.meta.env.VITE_API_LOGIN, {method: "POST", mode : "cors", headers: { "Content-Type": "application/json"}, body: JSON.stringify({username: userl.value, password : passl.value})}).then(a=>a.json()).then(a=>
	{
		
		
		if(a["Username"]!="no"){
			localStorage.setItem("log", 1)
			gls().log = 1
			gls().sessionT = a["Token"]
			gls().username = a["Username"]
			
			document.cookie = `name=${a["Username"]}; expires=${a["Expiry"]}; Secure`
			document.cookie = `sessionToken = ${a["Token"]}; expires=${a["Expiry"]}; Secure`
			router.push("mail")
		}
		else{
			wrong.value = 1
			
		}
	
	})
}


if (gls().log == 1){
	router.push("/mail")
}


function show(){
			if (passw.value=="password")
				passw.value = "text"
			else if(passw.value=="text")
				passw.value="password"
}

</script>

<template>
  <div id="login-bg">
    <form id="ch" v-on:submit.prevent="connect">
      <div id="mass">
        <img v-if="!isDark" src="@/assets/logo.png" width="100" height="80" />
        <img v-if="isDark" src="@/assets/logol.png" width="100" height="80" />
        <p class="brand">Office1789</p>
      </div>
      <div id="ensi">
        <input v-model="userl" type="text" placeholder="Identifiant" id="in" autocomplete="username"/>
        <div id="in2d">
          <input v-model="passl" :type="passw" placeholder="Mot de passe" id="in2" autocomplete="current-password"/>
          <button type="button" v-on:click="show" id="show" aria-label="Afficher le mot de passe">👁</button>
        </div>
        <transition name="fade">
          <p id="wrongP" v-if="wrong">{{$t('wrongL')}}</p>
        </transition>
      </div>
      <button id="lo" type="submit">{{$t('connection')}}</button>
      <div id="linka">
        <RouterLink to="/forgot">{{$t('forgot')}}</RouterLink>
        <RouterLink to="/createaccount">{{$t('create')}}</RouterLink>
      </div>
    </form>
  </div>
</template>

<style scoped>

/* Modern background with blur and gradient */
#login-bg {
  min-height: 100vh;
  width: 100vw;
  background: linear-gradient(120deg, #e0e7ef 0%, #f5f5f7 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}
#login-bg::before {
  content: "";
  position: absolute;
  inset: 0;
  background: rgba(240,240,245,0.7);
  backdrop-filter: blur(16px);
  z-index: 0;
}

.dark #login-bg {
  background: linear-gradient(120deg, #23243a 0%, #1a1b26 100%);
}
.dark #login-bg::before {
  background: rgba(30,30,40,0.7);
  backdrop-filter: blur(16px);
}


/* Centered and elevated login card */
#ch {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 28px;
  width: 100%;
  max-width: 400px;
  padding: 48px 32px 32px 32px;
  border-radius: 24px;
  background: rgba(255,255,255,0.85);
  box-shadow: 0 8px 32px rgba(0,0,0,0.18);
  z-index: 1;
  margin: 0 auto;
}

.dark #ch {
  background: rgba(30,30,40,0.92);
  color: #eee;
}

#mass {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}
.brand {
  font-family: roboto;
  font-size: 1.5rem;
  font-weight: 700;
  letter-spacing: 2px;
  color: #00308F;
  margin-top: 8px;
}
.dark .brand {
  color: #fff;
}

#ensi {
  display: flex;
  flex-direction: column;
  gap: 18px;
  width: 100%;
}

#in, #in2 {
  height: 44px;
  width: 90%;
  border-radius: 12px;
  border: 1px solid #d0d0d0;
  padding: 0 16px;
  font-size: 1.1rem;
  font-family: roboto;
  background: rgba(245,245,247,0.95);
  transition: border 0.2s, box-shadow 0.2s;
  box-shadow: 0 2px 8px rgba(0,0,0,0.06);
}
#in:focus, #in2:focus {
  border: 1.5px solid #00308F;
  box-shadow: 0 4px 16px rgba(0,48,143,0.10);
  outline: none;
}

.dark #in, .dark #in2 {
  background: rgba(30,30,40,0.95);
  color: #eee;
  border: 1px solid #444;
}

#in2d {
  display: flex;
  align-items: center;
  width: 100%;
  position: relative;
}

#show {
  background: none;
  color: #00308F;
  border: none;
  border-radius: 50%;
  margin-left: -40px;
  font-size: 1.3rem;
  cursor: pointer;
  padding: 6px;
  transition: background 0.2s, color 0.2s;
  z-index: 2;
}
#show:hover {
  color: white;
  background: #00308F;
}

#wrongP {
  position: relative;
  font-size: 0.95rem;
  color: #ff3c3c;
  text-align: center;
  margin-top: 4px;
  font-family: roboto;
  font-weight: 500;
  letter-spacing: 1px;
}

/* Fade transition for error */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}

#lo {
  width: 100%;
  height: 44px;
  background: -webkit-linear-gradient(30deg, #00308F, #ff3c3c);
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 1.2rem;
  font-family: roboto;
  font-weight: 700;
  letter-spacing: 1px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.10);
  cursor: pointer;
  transition: background 0.3s, box-shadow 0.2s;
}
#lo:hover {
  background: -webkit-linear-gradient(30deg, #ff3c3c, #00308F);
  box-shadow: 0 4px 16px rgba(0,0,0,0.18);
}

.dark #lo {
  background: -webkit-linear-gradient(30deg, #00308F, #ff3c3c);
  color: #fff;
}

#linka {
  display: flex;
  justify-content: center;
  gap: 18px;
  width: 100%;
  margin-top: 8px;
}
#linka > a {
  border-left: 1px solid #00308F;
  text-decoration: none;
  font-family: arial;
  padding: 5px 12px;
  color: #00308F;
  font-size: 1rem;
  transition: color 0.2s;
}
#linka > a:first-of-type {
  border: 0;
}
#linka > a:hover {
  color: #ff3c3c;
  text-decoration: underline;
}

.dark #linka > a {
  color: #6ec6ff;
  border-left: 1px solid #6ec6ff;
}
.dark #linka > a:hover {
  color: #ff3c3c;
}

/* Responsive */
@media (max-width: 600px) {
  #ch {
    max-width: 95vw;
    padding: 24px 8px 16px 8px;
    border-radius: 12px;
  }
  #mass img {
    width: 70px;
    height: 56px;
  }
  .brand {
    font-size: 1.1rem;
  }
}
</style>
