<script setup>
import {ref} from "vue"
import {gls} from "@/stores/global"
import router from "@/router/index"



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




</script>
<template>
	<div id="forma">
		<h1 id="tac">{{$t("contact")}} :</h1>
		<div id="champ">
			<ul v-html="$t('contactacl')">
			</ul>
			<ul>
				<li><input v-model="email" type="text" required/></li>
				<li><input v-model="username" type="text" required/>@office1789.com</li>
				
				<li><input /></li>
				<li><input v-model="phonenumber" type="text"/></li>
				<li><textarea style="height: 200px;"></textarea></li>
				
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
	position: relative;
	width: 100%;
	display: flex;
	flex-direction: column;
	align-items: center;
	top: 50%;
	height: auto;
	margin-top: 150px;
	text-align: center;
}
#champ{
	display: flex;
	text-align: left;
	justify-content: center;
	
	
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
	height: auto;
}

#bform{
	
	background: black;
	color: white;
	border: none;
	width: 100px;
	//margin: auto;
	
	
}

.dark #bform{
	margin-top: 20px;
	background: white;
	color: black;
	border: none;
	
}





</style>
