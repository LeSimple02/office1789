<script setup>
import {ref} from "vue"
import {gls} from "@/stores/global"
import router from "@/router/index"

let passw = ref("password")
let passw2 = ref("password")

let username = ref('')
let passf1 = ref('')
let passf2 = ref('')
let email = ref('')
let phonenumber = ref('')

let usernameR = ref('')
let emailR = ref('')
let phonenumberR = ref('')

function verif(){
	if(passf1.value == passf2.value && passf2.value !="" && passf1.value != "")
		connect()
}

function connect(){	
	let d = new Date()
	let dc = `${d.getFullYear()}-${(d.getMonth()+1).toString().padStart(2, '0')}-${d.getDate()} ${d.getHours()}:${d.getMinutes()}`

	fetch("http://127.0.0.1:8080/api/subscribe", { 
		method: "POST", 
		mode: "cors", 
		credentials: "same-origin", 
		headers: { "Content-Type": "application/json"}, 
		body : JSON.stringify({
			"username" : username.value, 
			"password": passf1.value, 
			"email": email.value, 
			"phonenumber": phonenumber.value, 
			"datejoined": dc, 
			"lastlogin": dc
		}) 
	}).then((v)=>{return v.json()}).then(
	(v)=>{
		if(v["Username"] != "no" && v["Phone"] !="no" && v["Email"] !="no"){
			localStorage.setItem("log", 1)
			gls().log = 1
			gls().username = v["Username"]
			gls().sessionT = v["Token"]
			
			document.cookie = `name=${v["Username"]}; expires=${v["Expiry"]}; Secure`
			document.cookie = `sessionToken=${v["Token"]}; expires=${v["Expiry"]}; Secure`
			
			router.push("mail")
		}
		else {
			if(v["Username"]){
				usernameR.value = 1			
			}
			if(v["Email"]){
				emailR.value = 1 
			}
			if(v["Phone"]){
				phonenumberR.value = 1 
			}
		}
	})
}

function show(){
	if (passw.value=="password")
		passw.value = "text"
	else if(passw.value=="text")
		passw.value="password"
}
function show2(){
	if (passw2.value=="password")
		passw2.value = "text"
	else if(passw2.value=="text")
		passw2.value="password"
}
</script>

<template>
<div id="create-bg">
	<div id="forma">
		<h1 id="tac">{{$t("createac")}} :</h1>
		<div id="champ">
			<ul id="col1" v-html="$t('createacl')"></ul>
			<ul id="col2">
				<p class="pr" v-if="usernameR">{{$t('dejaUP')}}</p>
				<li><input v-model="username" type="text" required/>@office1789.com</li>
				
				<li><input v-model="passf1" :type="passw" required/><input type="button" value="👁" @click="show()" class="show2" /></li>
				<li><input v-model="passf2" :type="passw2" required/><input type="button" value="👁" @click="show2()" class="show2"/></li>
				<p class="pr" v-if="emailR">{{$t('dejaEP')}}</p>
				<li><input v-model="email" type="text"/></li>
				<p class="pr" v-if="phonenumberR">{{$t('dejaPP')}}</p>
				<li><input v-model="phonenumber" type="text"/></li>
			</ul>
		</div>
		<p v-if="passf1!=passf2 && passf2!=''" style="color: red">{{$t('passwordd')}}</p>
		<button id="bform" @click="verif()">Submit</button>
	</div>
</div>
</template>

<style scoped>
#create-bg {
  min-height: 100vh;
  width: 100vw;
  background: linear-gradient(120deg, #e0e7ef 0%, #f5f5f7 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}
#create-bg::before {
  content: "";
  position: absolute;
  inset: 0;
  background: rgba(240,240,245,0.7);
  backdrop-filter: blur(16px);
  z-index: 0;
}

/* messages d'erreur */
.pr {
  font-size: 11px;
  color: red;
  margin: 0;
  grid-column: 2; /* force l'affichage côté inputs */
}

/* carte */
#forma {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 28px;
  width: 100%;
  max-width: 700px; /* un peu plus large pour la grille */
  padding: 48px 32px 32px 32px;
  border-radius: 24px;
  background: rgba(255,255,255,0.85);
  box-shadow: 0 8px 32px rgba(0,0,0,0.18);
  z-index: 1;
  margin: 0 auto;
}
#forma button {
  background: -webkit-linear-gradient(30deg, #00308F, #ff3c3c);
  color: white;
  border-radius: 20px;
  border: none;
  padding: 10px 18px;
  cursor: pointer;
}

#col1 ::v-deep{
    
    
    list-style: none;
  margin: 0;
  padding: 0;
  font-family: arial;
  display: flex;
  flex-direction: column;
  gap: 20px;
  text-align: left;
  li {
      
      
      align-items: center;
      display: flex;
      height: 40px;

  }
    
}

/* dark mode */
.dark #forma {
  background: rgba(30,30,40,0.7);
  backdrop-filter: blur(16px);
  color: #eee;
}

/* --- disposition grille --- */
#champ {
  display: grid;
  grid-template-columns: 1fr 1fr; /* gauche = règles, droite = inputs */
  gap: 20px;
  align-items: start;
  width: 100%;
}

#tac {
  font-family: roboto;
  text-align: center;
  margin: 0;
}

/* listes */
#col2 {
  list-style: none;
  margin: 0;
  padding: 0;
  
  font-family: arial;
  display: flex;
  flex-direction: column;
  gap: 20px;
  text-align: left;
  li{
    display: flex;
    gap: 10px;
  }
}

/* chaque ligne d'inputs */
#forma li {
  display: flex;
  justify-content: flex-start;
  align-items: center;
  width: 100%;
}

/* inputs */
#forma input[type="text"],
#forma input[type="password"],
#forma input:not([type="button"]) {
  width: 100%;
  padding: 6px 10px;
  border-radius: 8px;
  border: 1px solid #d0d0d0;
  box-sizing: border-box;
  background: rgba(245,245,247,0.95);
  height: 40px;
}

/* icône œil */
.show2 {
  color: black;
  border: none;
  background: none;
  height: 20px;
  margin-left: -10px;
  cursor: pointer;
}

/* responsive : repasse en colonne */
@media (max-width: 600px) {
  #champ {
    grid-template-columns: 1fr;
  }
}
</style>
