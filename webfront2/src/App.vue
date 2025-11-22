<script setup>
import {gls} from "@/stores/global"
import { RouterLink, RouterView} from 'vue-router'
import {ref, computed} from "vue"
import router from "@/router/index"
import { useDark, useToggle } from '@vueuse/core'
import verifback from "./views/backwork.js"


const isDark = useDark()
const toggleDark = useToggle(isDark)
const menuOpen = ref(false)

// Empêcher le scroll quand le menu est ouvert
const toggleMenu = () => {
  menuOpen.value = !menuOpen.value
  if (menuOpen.value) {
    document.body.classList.add('menu-open')
  } else {
    document.body.classList.remove('menu-open')
  }
}

const closeMenu = () => {
  menuOpen.value = false
  document.body.classList.remove('menu-open')
}

let menu1 = [{name : "Home", link : "/"},{name : "Login", link : "/login"},{name : "About", link : "/about"}]
let menu2 = [{name : "Mail", link : "/mail"}, {name : "Drive", link : "/drive"}, {name : "Chat", link : "/chat"}, {name : "Calendar", link : "/calendar"}, {name : "Account", link : "/account"}]
verifback()

let menu = computed(()=>{
  if(gls().log == 1){
    return menu2
  }
  else{
    return menu1
  }
})
let popup = ref(0)
function decof(){
  localStorage.setItem("log", 0)
  gls().log = localStorage.getItem("log")
  popup.value = 0
  router.push('/login')
}
</script>

<template>
  <header>
    <img class="logo" v-if="!isDark" src="@/assets/logo.png" />
    <img class="logo" v-if="isDark" src="@/assets/logol.png" />
    <button class="burger" @click="toggleMenu" aria-label="Menu">
      ☰
    </button>
    <nav :class="{ open: menuOpen }">
      <RouterLink
        v-for="key in menu"
        :key="key.name"
        :to="key.link"
        @click="closeMenu"
      >
        {{ $t(key.name) }}
      </RouterLink>

      <div class="nav-settings">
        <select v-model="$i18n.locale" @click="save()">
          <option
            v-for="locale in $i18n.availableLocales"
            :key="`locale-${locale}`"
            :value="locale"
          >
            {{ locale }}
          </option>
        </select>
        <button @click="toggleDark()" id="toggle">
          <i inline-block align-middle i="dark:carbon-moon carbon-sun" />
          <span class="ml-2">{{ isDark ? "🌕" : "☀️" }}</span>
        </button>
      </div>

      <!-- Bouton déconnexion toujours en dernier, desktop + mobile -->
      <button
        id="bdeco"
        v-if="gls().log == 1"
        @click="popup = 1; closeMenu()"
        title="Déconnexion"
      >
        Déconnexion
      </button>
    </nav>
  </header>
  <div v-if="popup==1" id="sb"></div>
  <div v-if="popup==1" id="pop">
    <p>Are you sure you want to log out?</p>
    <button @click="popup=0" id="Cancel">Cancel</button>
    <button @click="decof" id="Yes">Yes</button>
  </div>
  <RouterView></RouterView>
  <footer>
    <RouterLink to="/legalesmentions">{{ $t('legal') }}</RouterLink>
    <RouterLink to="/contact">Contact</RouterLink>
  </footer>
</template>

<script>
export default {
  methods:{
    save(){
      localStorage.setItem("lang", this.$i18n.locale)
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  width: 100%;
  overflow-x: hidden;
  margin: 0;
  padding: 0;
}

#app {
  width: 100%;
  overflow-x: hidden;
  min-height: 100vh;
}

@font-face {
  font-family: roboto;
  src: url("@/assets/roboto/RobotoCondensed-Regular.ttf");
}
.logo{
  width: 35px;
  height: 35px;
}
header{
  display: flex;
  height: 100px;
  align-items: center;
  justify-content: center;
  position: sticky;
  top: 0;
  background: white;
  z-index: 10000;
}
h1 {
  font-family: roboto;
}
p{
  font-family: arial;
}
header > nav > div {
  display: flex;
  align-items: center;
}
nav {
  padding-left: 16px;
  border-left: 1px solid black;
  display: flex;
  align-items: center;
  height: 35px;
  position: relative;
  gap: 28px;
}
nav a.router-link-active {
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 16px;
  color: #fff;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  border: none;
  outline: none;
  position: relative;
  filter: brightness(1.05);
  transform: translateY(-2px) scale(1.04);
  transition: all 0.3s ease;
}
nav a.router-link-active::before {
  content: '';
  position: absolute;
  inset: -3px;
  border-radius: 18px;
  padding: 3px;
  background: linear-gradient(45deg, #2979ff, #ed4956, #2979ff);
  -webkit-mask: linear-gradient(#fff 0 0) content-box, linear-gradient(#fff 0 0);
  -webkit-mask-composite: xor;
  mask-composite: exclude;
  animation: borderGlow 3s linear infinite;
}
@keyframes borderGlow {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.6; }
}
nav a.router-link-active:hover {
  color: #fff;
  background: -webkit-linear-gradient(30deg, blue, red);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.25);
  filter: brightness(1.08);
  transform: translateY(-3px) scale(1.06);
}
nav a {
  padding: 10px;
  font-size: 20px;
  border-left: 1px solid var(--color-border);
  font-family: 'Roboto', sans-serif;
  color: black;
  text-decoration: none;
  transition: all 0.3s ease;
}
nav a:hover {
  transform: translateY(-1px);
}
nav a:first-of-type {
  border: 0;
}
footer{
  position: fixed;
  bottom: 0;
  right: 0;
  background: rgba(0, 0, 0, 0.95);
  padding: 16px 24px;
  border-top-left-radius: 24px;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(10px);
  z-index: 90;
}
footer a{
  font-family: 'Roboto', sans-serif;
  padding: 0px 8px 0px 8px;
  text-decoration: none;
  color: white;
  font-weight: 500;
  font-size: 14px;
  transition: all 0.3s ease;
  border-radius: 12px;
  display: inline-block;
}
footer a:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateY(-2px);
}
footer a:first-of-type {
  border: 0;
}
.dark footer{
  background: rgba(255, 255, 255, 0.95);
}
.dark footer a{
  color: black;
}
.dark footer a:hover {
  background: rgba(0, 0, 0, 0.05);
}
#bdeco{
  background: -webkit-linear-gradient(30deg, #ff4444, #ff0000);
  color: white;
  border: none;
  height: 44px;
  padding: 0 24px;
  font-size: 18px;
  font-family: 'Roboto', sans-serif;
  font-weight: 600;
  margin-left: 20px;
  border-radius: 22px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(255, 68, 68, 0.3);
  display: flex;
  align-items: center;
  gap: 8px;
}
#bdeco:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 68, 68, 0.4);
  background: -webkit-linear-gradient(30deg, #ff0000, #cc0000);
}
#bdeco:active {
  transform: translateY(0);
}
#sb{
  width: 100%;
  height: 100%;
  position: fixed;
  background: rgba(0, 0, 0, 0.7);
  left: 0;
  top: 0;
  z-index: 999;
  backdrop-filter: blur(8px);
  animation: fadeIn 0.3s ease;
}
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
#pop{
  position: fixed;
  background: white;
  text-align: center;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 90%;
  max-width: 420px;
  padding: 40px 32px;
  border-radius: 32px;
  z-index: 1000;
  font-family: 'Roboto', sans-serif;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: slideIn 0.3s ease;
}
@keyframes slideIn {
  from { transform: translate(-50%, -60%); opacity: 0; }
  to { transform: translate(-50%, -50%); opacity: 1; }
}
#pop p {
  font-size: 1.2rem;
  font-weight: 500;
  color: #333;
  margin: 0 0 32px 0;
  line-height: 1.6;
}
#pop > #Cancel {
  background: rgba(150, 150, 150, 0.2);
  color: #333;
  border: 2px solid #ddd;
  padding: 12px 32px;
  border-radius: 16px;
  font-family: 'Roboto', sans-serif;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-right: 12px;
}
#pop > #Cancel:hover {
  background: rgba(150, 150, 150, 0.3);
  border-color: #999;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
#pop > #Yes {
  background: -webkit-linear-gradient(30deg, #ff4444, #ff0000);
  color: white;
  border: none;
  padding: 12px 32px;
  border-radius: 16px;
  font-family: 'Roboto', sans-serif;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(255, 68, 68, 0.3);
}
#pop > #Yes:hover {
  background: -webkit-linear-gradient(30deg, #ff0000, #cc0000);
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(255, 68, 68, 0.4);
}
#toggle{
  border: none;
  background: none;
  font-size: 30px;
}
.dark > body{
  background: #3B1C32;
  color: #fff;
}
.dark nav a.router-link-exact-active{
  background:  white;
  color: black;
}
.dark nav a{
  color:  white;
}
.dark nav {
  border-left:  1px solid white;
}
.dark nav a:hover{
  color: #fff;
  background: -webkit-linear-gradient(30deg, blue, red);
  box-shadow: 0 0 12px 2px rgba(41,121,255,0.28), 0 0 24px 4px rgba(237,73,86,0.18);
  filter: brightness(1.12);
  text-shadow: 0 2px 8px rgba(41,121,255,0.18);
}
.dark #pop{
  background: #1C1C1E;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.6);
}
.dark #pop p {
  color: white;
}
.dark #pop > #Cancel {
  background: rgba(100, 100, 100, 0.3);
  border-color: #444;
  color: white;
}
.dark #pop > #Cancel:hover {
  background: rgba(100, 100, 100, 0.5);
  border-color: #666;
}
.dark button{
  color: white;
}
select {
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 12px;
  color: white;
  padding: 8px 16px;
  border: none;
  font-size: 15px;
  font-family: 'Roboto', sans-serif;
  font-weight: 600;
  margin-right: 12px;
  transition: all 0.3s ease;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}
select option {
  background: white;
  color: #333;
  padding: 8px;
}
.dark select option {
  background: #2C2C2E;
  color: white;
}
select:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}
#toggle {
  border: none;
  
  
  display: flex;
  padding: 0;
  font-size: 28px;
  cursor: pointer;
  
  transition: all 0.3s ease;
  position: relative;
 
  
}
#toggle:hover {
  transform: scale(1.15) rotate(15deg);
}
.dark #toggle:hover {
  transform: scale(1.15) rotate(-15deg);
}
header {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 32px;
  padding: 20px 48px;
  background: rgba(255, 255, 255, 0.98);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border-radius: 0 0 32px 32px;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
  position: sticky;
  top: 0;
  z-index: 10000;
}
header:hover {
  box-shadow: 0 6px 28px rgba(0, 0, 0, 0.12);
}
.logo {
  width: 48px;
  height: 48px;
  margin-right: 8px;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  transition: all 0.3s ease;
  cursor: pointer;
}
.logo:hover {
  transform: scale(1.08) rotate(5deg);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.18);
}
.dark nav {
  background: rgba(30,30,40,0.85);
}
.dark nav a {
  color: #eee;
}
.dark nav a.router-link-active,
.dark nav a.router-link-exact-active,
.dark nav a:hover {
  background: -webkit-linear-gradient(30deg, blue, red);
  color: #fff;
}
.dark select{
  background: -webkit-linear-gradient(30deg, blue, red);
}
.dark header {
  background: rgba(30,30,40,0.95);
}
.dark .logo {
  box-shadow: 0 2px 8px rgba(0,0,0,0.30);
}

/* Burger menu */
.burger {
  display: none;
  background: none;
  border: none;
  font-size: 28px;
  cursor: pointer;
  color: black;
  padding: 8px;
  z-index: 10001;
  transition: all 0.3s ease;
}
.burger:hover {
  transform: scale(1.1);
}
.dark .burger{
  color: white;
}

/* Responsive */
@media (max-width: 768px) {
  header {
    flex-wrap: nowrap;
    padding: 12px 20px;
    gap: 12px;
    position: relative;
  }
  
  .logo {
    width: 40px;
    height: 40px;
    margin-right: auto;
  }
  
  .burger {
    display: block;
    margin-left: auto;
  }

  nav {
    position: fixed;
    top: 66px;
    left: 0;
    right: 0;
    background: white;
    width: 100%;
    height: calc(100vh - 66px);
    display: flex;
    flex-direction: column;
    align-items: stretch;
    justify-content: flex-start;
    padding: 24px 20px 40px 20px;
    gap: 12px;
    border-left: none;
    border-top: none;
    transform: translateX(-100%);
    transition: transform 0.3s ease;
    z-index: 9999;
    overflow-y: auto;
    overflow-x: hidden;
    box-shadow: 4px 0 20px rgba(0, 0, 0, 0.1);
  }
  
  nav.open {
    transform: translateX(0);
  }
  
  body.menu-open {
    overflow: hidden !important;
    position: fixed;
    width: 100%;
    height: 100%;
  }
  
  body.menu-open #app {
    overflow: hidden;
  }
  
  nav a {
    padding: 14px 20px;
    border: none;
    border-bottom: 1px solid #e5e7eb;
    width: 100%;
    text-align: left;
    font-size: 18px;
    border-radius: 8px;
    margin-bottom: 4px;
    
  }
  
  nav a:hover {
    background: rgba(0, 0, 0, 0.05);
  }
  
  .nav-settings {
    flex-direction: column;
    gap: 12px;
    padding-top: 16px;
    margin-top: 16px;
    border-top: 2px solid #e5e7eb;
    align-items: stretch;
    display: flex;
  }

  nav #bdeco {
    
    width: 20%;
    
    position: relative;
    left: 50%;
    margin-left: -10%;
    
    background: linear-gradient(135deg, #ef4444, #dc2626);
    height: auto;
    padding: 14px 20px;
    font-size: 1rem;
    border-radius: 12px;
    font-weight: 600;
    color: white;
    border: none;
    cursor: pointer;
    box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
    text-align: center;
  }

  .nav-settings select,
  .nav-settings button {
    width: 100%;
    justify-content: center;
  }
  
  .dark nav {
    background: rgba(30, 30, 40, 0.98);
    border-top-color: #444;
  }
  
  .dark nav a {
    border-bottom-color: #444;
  }
  
  .dark .nav-settings {
    border-top-color: #444;
  }
  
  footer {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    border-radius: 0;
    padding: 12px 16px;
    text-align: center;
    display: flex;
    justify-content: center;
    gap: 12px;
  }
}

@media (max-width: 480px) {
  header {
    padding: 10px 16px;
  }
  
  .logo {
    width: 36px;
    height: 36px;
  }
  
  nav {
    top: 60px;
    height: calc(100vh - 60px);
  }
  
  nav a {
    font-size: 16px;
    padding: 12px 16px;
  }
}
</style>
