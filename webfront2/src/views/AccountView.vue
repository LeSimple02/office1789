<script setup>
import AccountProfile from "@/components/AccountProfile.vue"
import AccountChange from "@/components/AccountChange.vue"
import { computed } from "vue"
import { useRoute } from 'vue-router'

const route = useRoute()

// Détermine quel composant afficher en fonction de la route
const showProfile = computed(() => route.path === '/account')
const showEdit = computed(() => route.path === '/account/edit')

</script>

<template>
<AccountProfile v-if="showProfile"></AccountProfile>
<AccountChange v-else-if="showEdit"></AccountChange>	
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
