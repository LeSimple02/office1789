<script setup>
import {ref} from "vue"
import {gls} from "@/stores/global.js"

let dj = ref(0)
let lj = ref(0)
let domain = ref(0)
let nboffer = ref(0)
let phone  = ref(0)
let email = ref(0)



fetch(process.env.VUE_APP_API_INFO_USER, {method:"POST", mode:"cors", body : JSON.stringify({"username" : gls().username, "token" : gls().sessionT}) }).then(a=>a.json()).then(a=>{dj.value=a['DateJoined']; domain.value=a['Domain'];nboffer.value=a['Nboffer']; email.value = a['Email']; phone = a['PhoneNumber']; lj.value = a["LastLogin"]})


</script>
<template>
<div id="enso">
		<h1 id="title">{{$t('infop')}} :</h1>
		<div id="lis">
			<li id="edit"><RouterLink id="edit" to="/account/edit">{{$t('edit')}}</RouterLink></li>
			<li id="pic"><img src="@/assets/napo.png" />{{$t('picturep')}}</li>
			<li>{{$t('username')}} : {{gls().username}}</li>
			<li>{{$t('password')}} : ●●●●●</li>
			<li>{{$t('doble')}} : ❌</li>
			<li>{{$t('domainy')}} : {{domain}}</li>
			<li>{{$t('offery')}} : {{nboffer}}</li>
			<li>{{$t('emaily')}} : {{email}}</li>
			<li>{{$t('phoney')}} : {{phone}}</li>
			<li>{{$t('lastj')}} : {{new Date(lj).toDateString()}}</li>
			<li>{{$t('datej')}} : {{new Date(dj).toDateString()}}</li>	
		</div>
	</div>
</template>

<style scoped>

#edit{
	margin-top: -40px;
	color: grey;
	text-decoration: none;
}
#edit:hover{
	text-decoration: underline;
}

#lis {
	list-style: none;
	font-family: roboto;
	font-size: 20px;
	position: absolute;
	height: auto;
	padding: 10px;
	left: 50%;
	top: 20%;
	width: 25%;
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

</style>
