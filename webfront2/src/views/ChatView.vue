<script setup>
//import {gls} from "@/stores/global.js"
import ChatConv from "@/components/ChatConv.vue"
import {ref, watch} from "vue"
import {useRoute} from 'vue-router'


let conv = ref([])
let convopt = ref([])

/*
function plus(){
	fetch(process.env.VUE_APP_API_CREATE_CONV, {method:"POST", mode:"cors", body : JSON.stringify({"username" : gls().username, "token" : gls().sessionT}) }).then(a=>a.text()).then(a=>{
		if(a != ""){
			conv.value= a.split(",")
			console.log(conv)
			for (let mess of conv.value){
				convopt.value.push({"user": mess, "link": "/chat/"+mess})
				
			}
			console.log(convopt)		
		}

	})
}
*/

let change = ref(window.location.href.split('/')[4] == 'edit' ? 1 : 0)
let calendar = ref(window.location.href.split('/')[4] == 'calendar' ? 1 : 0)


const route = useRoute()
function c(){
	change.value = window.location.href.split('/')[4] == 'edit' ? 1 : 0
	calendar.value = window.location.href.split('/')[3] == 'calendar' ? 1 : 0
	console.log(calendar.value)
	
}


watch(()=>route.path, c)


</script>
<template>
	
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
	<div id="aff" v-if="!change" >
		<img id="catE" src="@/assets/catE.png" />
		{{$t('affel')}}
	</div>
	<ChatConv v-if="change"/>
	
</template>


<style scoped>


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
		img{
			height: 50px;
	width: 50px;
	border-radius: 25px;
		}
		
	}
	
}

#conv{
	border-top : 2px solid black;
	display: flex;
	flex-direction: column;
	width: 300px;
	position: absolute;
	left: 70px;
	height: 100%;
	
	border-right : 2px solid black;
	padding: 10px;
	
	

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
