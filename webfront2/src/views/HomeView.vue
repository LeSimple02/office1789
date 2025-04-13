<script setup>
import {ref} from "vue"

let mov=0
let movnb = 0
let posc = ref(`transform: translateX(${mov}%)`)
let presentation2 = ["servicee", "servicec", "serviced", "servicew"]

function move(dir){
	mov = dir ? mov-=100 : mov+=100
	movnb = dir ? movnb+=1 : movnb-=1
	mov = mov<-300 ? -300 : mov
	mov = mov>0 ? 0 : mov
	movnb = movnb>3 ? 3 : movnb
	movnb = movnb<0 ? 0 : movnb
	
	posc.value = `transform: translateX(${mov}%)`
}

</script>

<template>		
		<div id="ens">
				<div id="colp1">
					<h1>{{ $t('presentation1I')}}</h1><br/><br/>
					<p id="presentation" v-html="$t('presentation1')"></p>
					<RouterLink to="/about" id="try">{{$t('try')}}</RouterLink>
				</div>
				<div id="ensA">
					<div id="ens2">
						<div id="carrousel">
							<button @click="move(0)">&lt;</button>
							<div id="carrouselc">
								<div :style="posc" id="carrouselm">
									<img src="@/assets/emailT.png" />
									<img src="@/assets/cloudT.png" />
									<img src="@/assets/dataT.png" />
									<img src="@/assets/writeT.png" />
								</div>
							</div>
							<button @click="move(1)">></button>
						</div>
						<p id="presentation2" v-html="$t(presentation2[movnb])"></p>
					</div>
					<div id="actu"><h3>{{$t('news')}}</h3>
					<ul v-html="$t('newsli')"></ul>
					</div>	
				</div>
		</div>
	
</template>
<style scoped>

ul ::v-deep{
	li {
		margin-top : 20px;
	}
	
}


#ens{
	display: flex;
	flex-direction: column;
	flex-wrap: wrap;
	gap: 80px;
	text-align: justify;
	width: 50%;
	margin: auto;
}

#ens2{
	display: flex;
	flex-wrap: wrap;
	
	flex-direction: column;
	gap: 20px;
	text-align: justify;
	height: auto;
	width: 55%;
	
}

#ensA{
	display: flex;
	flex-wrap: wrap;
	gap: 80px;
	width: 100%;
	
	justify-content: center;
}

#colp1{
	display: flex;
	flex-direction: column;


	
}

#presentation{
	height: 200px;

}

#try{
	background: black;
	color: white;
	font-family: arial;
	text-decoration: none;
	border-radius: 5px;
	font-size: 20px;
	width: 300px;
	max-width: 300px;
	text-align: center;
	padding: 20px;

}

.dark #try{
	background: white;
	color: black;

}

#carrousel{
	display: flex;
	align-items: center;
	justify-content: center;
	height: 300px;
	width: 100%;
	
	
}

#actu{
	border: 1px solid black;
	width: 30%;
	height: 400px;
	padding: 10px;
	font-family: arial
}

.dark #actu{
	border: 1px solid white;
}

#actu > h3 {
	font-family: roboto;
	border-bottom: 1px solid black;
}

.dark #actu > h3 {
	border-bottom: 1px solid white;
}



#carrousel img {
	width: 100%;

}

#carrousel button {
	background: black;
	border-radius: 20px;
	color: white;
	padding: 10px;
	border: none;
	font-size: 30px;
}

#carrouselc{
	display: flex;
	align-items: center;
	
	overflow: hidden;
}
#carrouselm{
	display: flex;
	transition: ease 1s;
}


</style>
