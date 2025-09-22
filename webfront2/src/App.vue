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
    <button class="burger" @click="menuOpen = !menuOpen" aria-label="Menu">
      ☰
    </button>
    <nav :class="{ open: menuOpen }">
      <RouterLink v-for="key in menu" :key="key.name" v-bind:to="key.link" @click="menuOpen = false">{{ $t(key.name) }}</RouterLink>
      <div>
        <select v-model="$i18n.locale" @click="save()">
          <option v-for="locale in $i18n.availableLocales" :key="`locale-${locale}`" :value="locale">{{ locale }}</option>
        </select>
        <button @click="toggleDark()" id="toggle">
          <i inline-block align-middle i="dark:carbon-moon carbon-sun" />
          <span class="ml-2">{{ isDark ? "🌕" : "☀️" }}</span>
        </button>
      </div>
    </nav>
    <button id="bdeco" v-if="gls().log == 1" @click="popup = 1" title="Déconnexion">🔓</button>
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
  height: 50px;
  align-items: center;
  justify-content: center;
  position: relative;
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
  padding-left: 10px;
  border-left: 1px solid black;
  display: flex;
  align-items: center;
  height: 35px;
  position: relative;
  gap: 20px;
}
nav a.router-link-active {
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 10px;
  color: #fff;
  background: -webkit-linear-gradient(30deg, blue, red);
  box-shadow:
    0 0 16px 4px rgba(41,121,255,0.25),
    0 0 32px 8px rgba(237,73,86,0.18),
    0 2px 8px rgba(41,121,255,0.18);
  border: 2px solid #2979ff;
  outline: 2px solid #ed4956;
  filter: brightness(1.12) drop-shadow(0 2px 8px #2979ff);
  text-shadow: 0 2px 12px rgba(41,121,255,0.18);
  transform: translateY(-2px) scale(1.04);
}
nav a.router-link-active:hover {
  color: #fff;
  background: -webkit-linear-gradient(30deg, blue, red);
  box-shadow: 0 0 12px 2px rgba(41,121,255,0.18), 0 0 24px 4px rgba(237,73,86,0.12);
  filter: brightness(1.08);
  text-shadow: 0 2px 8px rgba(41,121,255,0.12);
}
nav a {
  padding: 0 1rem;
  font-size: 20px;
  border-left: 1px solid var(--color-border);
  font-family: roboto;
  color: black;
  text-decoration: none;
}
nav a:first-of-type {
  border: 0;
}
footer{
  position: fixed;
  bottom: 0%;
  right: 0%;
  background: black;
  padding: 10px;
  border-top-left-radius: 10px;
}
footer a{
  font-family: roboto;
  padding: 10px;
  text-decoration: none;
  color: white;
}
footer a:first-of-type {
  border: 0;
}
.dark footer{
  position: fixed;
  bottom: 0%;
  right: 0%;
  background: white;
  padding: 10px;
  border-top-left-radius: 10px;
}
.dark footer a{
  color: black;
}
#bdeco{
  background: rgba(255, 0,0,0.5);
  color: white;
  border: none;
  height: 40px;
  font-size: 20px;
  margin-left: 20px;
  border-radius: 20px;
}
#sb{
  width: 100%;
  height: 100%;
  position: absolute;
  background: rgba(0,0,0,0.6);
  left: 0%;
  top: 0%;
}
#pop{
  position: absolute;
  background: white;
  text-align: center;
  left: 50%;
  top: 50%;
  width: 300px;
  padding: 30px;
  border-radius: 10px;
  margin-left: -165px;
  margin-top: -50px;
  z-index: 1;
  font-family: arial;
}
#pop > #Cancel {
  background: rgba(255, 0,0,0.5);
  color: white;
  border: none;
}
#pop > #Yes {
  background: rgba(0, 200,0, 0.5);
  color: white;
  border: none;
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
  color: black;
}
.dark button{
  color: white;
}
select {
  background: -webkit-linear-gradient(30deg, blue, red);
  border-radius: 8px;
  color: black;
  padding: 4px 12px;
  border: none;
  font-size: 16px;
  margin-right: 8px;
  transition: background 0.3s, color 0.3s;
}
#toggle {
  border: none;
  background: none;
  font-size: 28px;
  cursor: pointer;
  margin-left: 8px;
  transition: color 0.3s;
}
header {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 24px;
  padding: 12px 32px;
  background: rgba(245,245,247,0.95);
  box-shadow: 0 2px 8px rgba(0,0,0,0.07);
  border-radius: 0 0 24px 24px;
}
.logo {
  width: 40px;
  height: 40px;
  margin-right: 16px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.10);
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
  padding: 4px;
  margin-left: auto;
  order: 2;
}
.dark .burger{
  color: white;
}

/* Responsive */
@media (max-width: 700px) {
  header {
    flex-direction: row;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    padding: 8px 16px;
  }
  .burger {
    display: block;
    order: 1;
  }
  nav {
    position: fixed;
    top: 60px;
    left: 0;
    right: 0;
    background: rgba(245,245,247,0.98);
    width: 100%;
    height: 100%;
    flex-direction: column;
    
    padding: 16px;
    gap: 8px;
    border-left: none;
    border-top: 1px solid var(--color-border);
    transform: translateY(-120%);
    transition: transform 0.3s ease;
    z-index: 10;
  }
  nav.open {
    transform: translateY(0);
  }
  nav a {
    padding: 8px 16px;
    border: none;
    width: 80%;
    
    border-bottom: 1px solid var(--color-border);
  }
  nav > div {
    flex-direction: row;
    justify-content: space-between;
    padding-top: 8px;
    border-top: 1px solid var(--color-border);
  }
  .logo {
    margin-right: 0;
    order: 0;
  }
  #bdeco {
    order: 3;
    margin-left: 0;
  }
  .dark nav {
    background: rgba(30,30,40,0.98);
  }
  footer{
    display: none;
  }
}
</style>
