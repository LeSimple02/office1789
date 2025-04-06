<script setup>
import {gls} from "@/stores/global"
import { RouterLink, RouterView} from 'vue-router'
import {ref, computed} from "vue"
import router from "@/router/index"
import { useDark, useToggle } from '@vueuse/core'
import verifback from "./views/backwork.js"

const isDark = useDark()
const toggleDark = useToggle(isDark)

let menu1 = [{name : "Home", link : "/"},{name : "Login", link : "/login"},{name : "About", link : "/about"}]
let menu2 = [{name : "Mail", link : "/mail"}, {name : "Drive", link : "/drive"}, {name : "Chat", link : "/chat"}, {name : "Account", link : "/account"}]


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
       <nav>
        <RouterLink v-for="key in menu" :key="key.name" v-bind:to="key.link">{{ $t(key.name) }}</RouterLink>
      </nav>
      <button id="bdeco" v-if="gls().log == 1" @click="popup = 1">⏼</button>
      <div>
        <select v-model="$i18n.locale" @click="save()">
            <option v-for="locale in $i18n.availableLocales" :key="`locale-${locale}`" :value="locale">{{ locale }}</option>
          </select>
        <button @click="toggleDark()" id="toggle">
          <i inline-block align-middle i="dark:carbon-moon carbon-sun" />

          <span class="ml-2">{{ isDark ? "🌕" : "☀️" }}</span>
          </button>
      </div>
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
}

h1 {
  font-family: roboto;
}

p{
  font-family: arial;
}



header > div {
  position: absolute;
  display: flex;
  align-items: center;
  gap: 20px;
  right: 0%;
  
}
nav {
  padding-left: 10px;
  border-left: 1px solid black;
  display: flex;
  align-items: center;
  height: 35px;  
}

nav a.router-link-exact-active {
  background: black;
  border-radius: 10px;
  color: white;
}

nav a.router-link-exact-active:hover {
  color: white;
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

select{
  background: black;
  border-radius: 10px;
  color: white;
}

#bdeco{
  background: rgba(255, 0,0,0.5);
  color: white;
  border: none;
  height: 40px;
  font-size: 20px;
  margin-left: 20px;
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
  color:  rgb(200, 200, 200);
}

.dark nav a.router-link-exact-active:hover{
  background:  white;
  color: black;
}

.dark #pop{
  color: black;
}

.dark select{
  background: white;
  color: black;
}

.dark button{
  color: white;
}

</style>
