<script setup>
import {ref} from "vue"

let id = ref("")
let sendm= ref(false)

function send(){
	if(id.value != "")
		sendm.value = true
	else{
		sendm.value = 'nothing'
	}

}


</script>
<template>
	<div id="formem">
			<label for="fpassw">{{$t('forgotp')}} :</label> 
			<input v-model="id" name="fpassw" type="text" required/>
			<button id="bform" @click="send">Submit</button>
			<p v-if="sendm=='nothing'" style="color: red">Enter your id</p>
			<p v-if="sendm==true"  style="color: red">Un email vient de vous être envoyé</p>
	</div>

</template>
<style>

#formem {
	position: absolute;
	left: 50%;
	top: 50%;
	width: 200px;
	margin-left: -100px;
	text-align: center;
	font-family : arial;
}

#formem > button {
	margin-top: 10px;
	background: black;
	border : none;
	color: white;
}

.dark div > button {
	background: black;
	border : none;
	color: white;
}


</style>
