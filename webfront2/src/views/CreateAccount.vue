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

			fetch("http://127.0.0.1:8080/api/subscribe", { method: "POST", mode: "cors", credentials: "same-origin", headers: { "Content-Type": "application/json"}, body : JSON.stringify({"username" : username.value, "password": passf1.value, "email": email.value, "phonenumber": phonenumber.value, "datejoined": dc, "lastlogin": dc}) }).then((v)=>{return v.json()}).then(
			(v)=>{
				if(v["Username"] != "no" && v["Phone"] !="no" && v["Email"] !="no"){
					localStorage.setItem("log", 1)
					gls().log = 1
					gls().username = v["Username"]
					gls().sessionT = v["Token"]
					
					document.cookie = `name=${v["Username"]}; expires=${v["Expiry"]}; Secure`
			document.cookie = `sessionToken = ${v["Token"]}; expires=${v["Expiry"]}; Secure`
					
					router.push("/mail")
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
	<div id="forma">
		<h1 id="tac">{{$t("createac")}} :</h1>
		<div id="champ">
			<ul v-html="$t('createacl')">
			</ul>
			<ul>
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

</template>
<style>

.pr{
	position: fixed;
	font-size: 10px;
	color: red;
}

#forma{
	position: absolute;
	left: 50%;
	width: 650px;
	margin-left: -310px;
	top: 50%;
	
	
	margin-top: -250px;
	text-align: center;
}
#champ{
	display: flex;
	text-align: left;
	
}

#tac{
	font-family: roboto;
}

#forma ul {
	list-style: none;
	margin-top: 20px;
	font-family: arial;
}

#forma li{
	margin-top: 20px;
	height: 30px;
}

#bform{
	margin-top: 20px;
	background: black;
	color: white;
	border: none;
	
}

.dark #bform{
	margin-top: 20px;
	background: white;
	color: black;
	border: none;
	
}


.show2{
	
	color: black;
	border: none;
	background: none;
	height:20px;
	margin-left: -25px;
}





</style>
