<script setup>
import {ref} from "vue"
import {gls} from "@/stores/global.js"

let dj = ref(0)
let lj = ref(0)
let domain = ref(0)
let nboffer = ref(0)
let phone  = ref(0)
let email = ref(0)

let passwordt = ref('password')
let passwordt2 = ref('password')

let newusername = ref('')
let newphone = ref('')
let newemail = ref('')
let newoffer = ref(0)

let passf1 = ref('')
let passf2 = ref('')

let usernameR = ref('')
let emailR = ref('')
let phonenumberR = ref('')

fetch(import.meta.env.VITE_APP_API_INFO_USER, {method:"POST", mode:"cors", body : JSON.stringify({"username" : gls().username, "token" : gls().sessionT}) }).then(a=>a.json()).then(a=>{dj.value=a['DateJoined']; domain.value=a['Domain'];nboffer.value=a['Nboffer']; email.value = a['Email']; phone.value = a['PhoneNumber']; lj.value = a["LastLogin"]})

function send(){
	if(passf1.value==passf2.value){
		fetch(process.env.VUE_APP_API_INFO_CHANGE, {method:"POST", mode:"cors", body : JSON.stringify({"lastusername" : gls().username, "username" : newusername.value, "phonenumber": newphone.value, "email": newemail.value, "nboffer": newoffer.value, "password": passf2.value, "token": gls().sessionT }) })
		.then(a=>a.json()).then((a)=>{
			if(newusername.value != "" && a["Username"] != "no" && a["Email"] != "no" && a["Phone"] != "no"){
				
				console.log(a["Username"])
				document.cookie = `name=${gls().username}; expires=Fri, 31 Dec 1900 23:59:59 GMT; Secure`
				document.cookie = `sessionToken=${gls().sessionT}; expires=Fri, 31 Dec 1900 23:59:59 GMT; Secure`
				gls().sessionT = a["Token"]
				gls().username = a["Username"]
				document.cookie = `name=${a["Username"]}; expires=${a["Expiry"]}; Secure`
				document.cookie = `sessionToken = ${a["Token"]}; expires=${a["Expiry"]}; Secure`
				window.location.href = "/account"
			}
			else {
				
				if(a["Username"]){
					usernameR.value = 1			
				}
				if(a["Email"]){
					emailR.value = 1 
				}
				if(a["Phone"]){
					phonenumberR.value = 1 
				}
					
			}})
			
			

	}
}


function show(){
			if (passwordt.value=="password")
				passwordt.value = "text"
			else if(passwordt.value=="text")
				passwordt.value="password"
}

function show2(){
			if (passwordt2.value=="password")
				passwordt2.value = "text"
			else if(passwordt2.value=="text")
				passwordt2.value="password"
}
</script>

<template>
<h1 id="title">{{$t('infop')}} :</h1>
<div id="enso">
		<li id="edit"><RouterLink id="edit" to="/account">⬅️</RouterLink></li>
		<li id="pic"><img src="@/assets/napo.png" />{{$t('picturep')}}</li>
		<div id="enso2">
			<div id="lis">
				<li>{{$t('username')}} :</li>
				<li style="padding-top: 20px; padding-bottom: 20px;">{{$t('password')}} :</li>
				<li>{{$t('doble')}} :</li>
				<li>{{$t('domainy')}} :</li>
				<li>{{$t('offery')}} :</li>
				<li>{{$t('emaily')}} :</li>
				<li>{{$t('phoney')}} :</li>
				<li>{{$t('lastj')}} : {{new Date(lj).toDateString()}}</li>
				<li>{{$t('datej')}} : {{new Date(dj).toDateString()}}</li>		  		
			</div>
			<div id="lisrep">
				<p class="pr" v-if="usernameR">{{$t('dejaUP')}}</p>
				<li><input v-model="newusername" v-bind:placeholder="gls().username"/></li>
				<li id="passc" style="display: grid;"><p v-if="passf1!=passf2 && passf2!=''" style="color: red; position: fixed; margin-top: -10px; font-size:10px;">{{$t('passwordd')}}</p><div><input v-bind:type="passwordt" v-model="passf1" v-bind:placeholder="$t('passwordN')"/><button type="button" v-on:click="show" id="show">👁</button></div><div><input v-model="passf2" style="margin-top: 10px;" v-bind:placeholder="$t('repassword')" v-bind:type="passwordt2"/><button type="button" v-on:click="show2" id="show">👁</button></div></li>
				<li style="display: flex; gap :20px;">❌<button class="buttons">config</button></li>
				<li><select><option>@office1789</option></select></li>
				<li><select><option>1</option><option>2</option><option>3</option></select>&nbsp;<RouterLink style="text-decoration: none; color: red; font-size: 10px;" to="/about">⚠️{{ $t("About") }}</RouterLink></li>
				<p class="pr" v-if="emailR">{{$t('dejaEP')}}</p>
				<li><input v-model="newemail" v-bind:placeholder="email"/></li>
				<p class="pr" v-if="phonenumberR">{{$t('dejaPP')}}</p>
				<li><input  v-model="newphone"  v-bind:placeholder="phone"/></li>			
			</div>
		</div>
		<div id="choice">
			<button @click="send">✅</button>
			<RouterLink  style="text-decoration:none;" to="/account">❌</RouterLink>
		</div>	
	</div>
</template>

<style scoped>

#lisrep{
	
		.buttons{
			border: none;
			background: lightgrey;
			border-radius: 5px;
		}
	
} 

.pr{
	position: fixed;
	font-size: 10px;
	color: red;
}

#passc{
	display: flex;
}

#passc > div > button{
	background: none;
	border: none;
}
#choice{
	border: none;
	background: none;
	font-size: 30px;
	position: absolute;
	left: 80%;
	top: 90%;
}

#choice > button{
	border: none;
	background: none;
	font-size: 30px;
}

input{
	border-radius: 2px;
	border : none;
}

#edit{
	margin-top: -40px;
	color: grey;
	text-decoration: none;
	
}

#enso2{
	display: flex;
}

#enso{

	list-style: none;
	font-family: roboto;
	font-size: 20px;
	position: absolute;
	height: auto;
	padding: 10px;
	left: 50%;
	top: 20%;
	width: auto;
	min-width: 200px;
	text-align: left;
	margin-left: -10%;
	border-radius: 30px;

}

#pic{
	display: grid;
	align-items: center;
	text-align: center;
	justify-content: center;
	font-weight: bold;
}
#pic >img{
	width: 200px;
	border-radius: 100px;
	
}

li{
	margin-top: 20px;
}

#title{
	font-family: roboto;
	font-size: 40px;
	position: absolute;
	left: 3%;
	
}

.dark #lisrep{
		
		.buttons{
			border: none;
			background: grey;
			border-radius: 5px;
		}
	
} 

</style>
