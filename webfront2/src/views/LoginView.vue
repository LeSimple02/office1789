<script setup>
import {gls} from "@/stores/global"
import {ref} from "vue"
import router from "@/router/index"




let passw = ref("password")

let userl = ref('')
let passl = ref('')
let wrong = ref(0)

function connect(){

	
	
	fetch("http://localhost:8080/api/connect", {method: "POST", mode : "cors", headers: { "Content-Type": "application/json"}, body: JSON.stringify({username: userl.value, password : passl.value})}).then(a=>a.json()).then(a=>
	{
		
		
		if(a["Username"]!="no"){
			localStorage.setItem("log", 1)
			gls().log = 1
			gls().sessionT = a["Token"]
			gls().username = a["Username"]
			
			document.cookie = `name=${a["Username"]}; expires=${a["Expiry"]}; Secure`
			document.cookie = `sessionToken = ${a["Token"]}; expires=${a["Expiry"]}; Secure`
			router.push("/mail")
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
	<form id="ch" v-on:submit.prevent="connect">
			<div id="mass">
				<img v-if="!isDark" src="@/assets/logo.png" width="100" height="80" />
				<img v-if="isDark" src="@/assets/logol.png" width="100" height="80" />
				<p>Office1789</p>
			</div>
			<div id="ensi">
				<input v-model="userl" type="text" placeholder="     ID" id="in"/>
				<div id="in2d"><input v-model="passl" :type="passw" placeholder="    Password" id="in2"/><button type="button" v-on:click="show" id="show">👁</button></div>
				<p id="wrongP" v-if="wrong">{{$t('wrongL')}}</p>
			</div>
			<button id="lo" type="submit">{{$t('connection')}}</button>
			<div id="linka"><RouterLink to="/forgot">{{$t('forgot')}}</RouterLink>
			<RouterLink to="/createaccount">{{$t('create')}}</RouterLink></div>
	</form>
</template>

<style scoped>

#wrongP{
	position: relative;
	font-size: 10px; 
	color: red;
	text-align: center;
}

#mass {
	display: flex;
	text-align: center;
	justify-content: center;
	align-items: center;
	font-family: roboto;
}

#ch {
	position: absolute;
	display: flex;
	flex-direction : column;
	align-items: center;
	top: 50%;
	left: 50%;
	gap: 20px;
	margin-left: -250px;
	margin-top: -200px;
	height: 400px;
	width: 450px;
	padding: 50px;
	border-radius: 20px;
	
}

#lo {
	width: 30%;
	background: black;
	height: 40px;
	border: none;
	border-radius: 10px;
	color: white;
}

.dark #lo {
	width: 30%;
	background: white;
	color: black;
}

#in {
	height: 30px;
	width: 100%;	
}

#in2d{
	display: flex;
	align-items: center;
	width: 100%;
	height: 30px;
}

#in2{
	height: 100%;
	width: 100%;
}


#ensi{
	display: flex;
	flex-direction: column;
	gap: 20px;
	width: 80%;

}

#show {
	background:  none;
	color: #00308F;
	border: none;
	border-radius: 20px;
	margin-left: -10%;	
}


#show:hover {
	color:  white;
	background: #00308F;
}




.dark #linka > a {
	color: skyblue;

}

#linka > a:first-of-type {
	border: 0;
}

#linka > a {
  border-left: 1px solid green;
  text-decoration : none;
  font-family: arial;
  padding: 5px;
}
	
</style>
