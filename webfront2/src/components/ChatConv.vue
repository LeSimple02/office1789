<script setup>
import {gls} from "@/stores/global.js"
import {ref} from "vue"
import CallAc from "@/components/CallAc.vue"


let conv = ref([])
let convopt = ref([])
let url = ref(window.location.href.split("/")[4])


fetch(process.env.VUE_APP_API_INFO_CONV, {method:"POST", mode:"cors", body : JSON.stringify({"username" : gls().username, "token" : gls().sessionT}) }).then(a=>a.text()).then(a=>{
	if(a != ""){
		conv.value= a.split(",")
		console.log(conv)
		for (let mess of conv.value){
			convopt.value.push({"user": mess, "link": "/chat/"+mess})
			
		}
		console.log(convopt)		
	}

})

let callA = ref(0)

function call(){
	
	callA.value = 1
	let sound = new Audio(require("@/assets/call.mp3"))
	sound.play()

}

/*
fetch(process.env.VUE_APP_API_INFO_MESSAGE, {method:"POST", mode:"cors", body : JSON.stringify({"username" : gls().username, "token" : gls().sessionT, "conv": url.value}) }).then(a=>a.text()).then(a=>{
	if(a != ""){
		conv.value= a.split(",")
		console.log(conv)
		for (let mess of conv.value){
			convopt.value.push({"user": mess, "link": "/chat/"+mess})
			
		}
		console.log(convopt)		
	}

})*/



</script>
<template>
	<div id="opt">
		<RouterLink to="/chat">💬</RouterLink>
		<RouterLink to="/chat/notif">📅</RouterLink>
		<RouterLink to="/chat/notif">⚙️</RouterLink>
		
	</div>
	<div id="conv">
		<input v-bind:placeholder="$t('userids')"/>
		<span style="margin-top: 10px;border-bottom: 1px solid grey;"></span>
		<div v-if="conv==[]"  id="not">
			<img id="cat" src="@/assets/cat.png" />
					{{$t('nothing')}}
		</div>
		<div v-if="conv!=[]" id="lconv">
			<div v-for="m in convopt" :key="m" ><img src="@/assets/napo.png" /><RouterLink v-bind:to="m.link">{{m.user}}</RouterLink></div>
		</div>
		
	</div>
	<div id="message">
		<div id="barm">
			<h2>{{url}}</h2>
			<button @click="call">📞</button>
			<button>⚙️</button>
		</div>
		<div id="messconv">
				<div class="messageR">f est dérivable !</div>
				<div class="messageS">Malheureux ! <br />Diantre</div>
		</div>
		<div id="sendbar">
			<textarea rows = "5" cols = "10" name = "description" placeholder="Sendmessage"></textarea>
			<button>📨</button>
			<button>📎</button>
		</div>
	</div>
	<CallAc v-if="callA"></CallAc>
	
</template>


<style scoped>







#message{
	position: fixed;
	width: 50%;
	height: 100%- 50px;
	top: 50px;
	font-family: roboto;
	left: 50%;
	margin-left: -25%;
	bottom: 20px;
	justify-content: center;
	#barm{
		display: flex;
		height: 50px;
		align-items: center;
		position: relative;
		left: 50%;
		
		width:70%;
		border-bottom: 2px solid black;
		margin-left: -25%;
		text-align: center;
		justify-content: center;
		
		button{
			border-radius:	15px;
			height: 30px;
			margin-left: 10px;
			background: none;
				
			position: relative;
			left: 40%;
		}
	}
	
	#sendbar{
		display: flex;
		align-items: center;
		justify-content: center;
		text-align: center;
		position: absolute;
		bottom: 0px;
		width: 70%;
		left: 50%;
		margin-left: -25%;
		textarea{
			width: 50%;
			height: 50px;
			border-radius: 10px;
			
		}
		button{
			height: 40px;
			width: 40px;
			font-size: 20px;
			margin-left: 10px;
			border-radius:	30px;
			background: none;
			border: 1px solid black;	
		}
	}
	#messconv{

	display: grid;
	position : relative;
	
	width: 70%;
	left: 50%;
	margin-left: -25%;
	.messageR{
		position: relative;
		background: lightgrey;
		border-radius: 10px;
		padding: 10px;
		width: auto;
		left: 0%;
		max-width: 200px;
		
		font-family: arial;
		margin-top: 20px;
	}
	.messageS{
		position: relative;
		background: skyblue;
		border-radius: 10px;
		padding: 10px;
		width: auto;
		left: 100%;
		margin-left: -220px;
		max-width: 200px;
		font-family: arial;
		margin-top: 20px;
	}
	
}
}

#lconv{
	
	position: relative;
	left: 0%;
	
	div{
		list-style: none;
		
		padding: 20px;
		display: flex;
		padding: 25px; 
		width: 100%;
		height: 20px;
		align-items: center;
		
		
		width: auto;
		border-bottom: 1px solid grey;
		
		a {
			text-decoration: none;
			font-family: roboto;
			font-weight: bold;
			color: grey;
			margin-left: 10px;
			
			
		}
			a.router-link-exact-active {
	background: lightblue;
	padding: 10px;
	border-radius: 5px;
	color: black;
}
		img{
			height: 50px;
	width: 50px;
	border-radius: 25px;
		}
		
	}
	
}


#opt{
	position: absolute;
	left: 0%;
	
	display: flex;
	flex-direction: column;
	width: 50px;
	height: 100%;
	background: lightgrey;
	align-items: center;
	text-align: center;
	padding: 10px;
	font-size: 30px;
	
	
	a{
		
		padding: 10px;
		text-decoration: none;
	};
}

#conv{
	
	display: flex;
	flex-direction: column;
	width: 300px;
	position: absolute;
	left: 70px;
	height: 100%;
	
	border-right : 2px solid black;
	padding: 10px;
	
	input{
		width: 100%;	
	}

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

#catE{
	width: 200px;
	height: 200px;
	position: relative;
	left: 50%;
	margin-left: -100px;
	
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

</style>
