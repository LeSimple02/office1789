<script setup>
import {ref} from "vue"
import {gls} from "@/stores/global"
import router from "@/router/index"


if(gls().log != 1){
	router.push("/login")
}


let curr = ref('')

changec()
function changec(){
	if(window.location.href.split('/')[4]!=undefined)
		curr.value = window.location.href.split('/')[4]
	else{
		curr.value = window.location.href.split('/')[3]
	}
	
}



</script>

<template>	
		<div id="area">	
			<div id="chooseA">
			<ul id="chooseL" @click="changec">
				<li>📥 <RouterLink to="/mail">{{$t('mail')}}</RouterLink></li>
				<li>✏️ <RouterLink to="/mail/draft">{{$t('draft')}}</RouterLink></li>
				<li>✉️ <RouterLink to="/mail/send">{{$t('send')}}</RouterLink></li>
				<li>🗑️ <RouterLink to="/mail/trash">{{$t('trash')}}</RouterLink></li>
				<li>📨 <RouterLink to="/mail/sendem">{{$t('sendem')}}</RouterLink></li>
			</ul>
			
			</div>
			<div id="ls">
				<div id="bs">
					{{$t(curr)}}
				</div>
				<input v-bind:placeholder="$t('userids')"/>
				<div id="not">
				<img id="cat" src="@/assets/cat.png" />
					{{$t('nothing')}}
				</div>
			</div>
			<div id="aff">
				<img id="catE" src="@/assets/catE.png" />
				{{$t('affel')}}
			</div>
		
		</div>
	
</template>
<style scoped>

#area{
	overflow: hidden;
	position: absolute;
	width: 100%;
	left: 0%;
	
	height: 100%;
}

#aff {
	position: absolute;
	border-radius: 20px;
	width: 60%;
	left: 40%;
	display: grid;
	font-family : roboto;
	text-align: center;
}

#catE{
	width: 200px;
	height: 200px;
	position: relative;
	left: 50%;
	margin-left: -100px;
	
}

#chooseA{
	
	position: absolute;
	
	height: 100%;
	left: -10px;
	width: 200px;
	background: lightgrey;
	font-family: roboto;
	
	padding: 20px;
	z-index: 2;
	border-top: 2px solid black;
	
}
#ls{
	position: absolute;
	min-height: 100%;
	width:300px;
	border-right: 2px solid black;
	border-top: 2px solid black;
	left: 230px;
	border-top-right-radius: 10px;
	overflow: hidden;
	z-index: 1;
	
	input{
		position: relative;
		left: 5%;
		margin-top: 10px;
		width: 90%;
		height: 20px;
			
	}
	
	
}
#bs{
	
	position: relative;
	text-align: center;
	width: 100%;
	background: red;
	font-family: roboto;
	
	color: white;
	border-top-right-radius: 9px;
}

#not
{
	width: 200px;
	height: 150px;
	display: grid;
	position: relative;
	left: 150px;
	margin-left: -100px;
	margin-top: 50px;
	font-family: roboto;
	text-align: center;
}

#cat{
	width: 100px;
	height: 100px;
	position: relative;
	left: 50%;
	margin-left: -50px;
}

ul{
	list-style: none;
	margin-left: -20px;
}
li{
	margin-bottom: 10px; 
	
}

li > a {
	padding: 3px;
	border-radius: 3px;
	text-decoration: none;
	color: black;
}

li > a.router-link-exact-active {
	background: grey;
	color: white;
}

.dark li >a {
	color: white
}

.dark li > a.router-link-exact-active {
	background: lightgrey;
	color: black;
}

.dark #chooseA {
	background: grey;
	border-top: 2px solid white;
}

.dark #ls {
	border : 2px solid white;
}


</style>
