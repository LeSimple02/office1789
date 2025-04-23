<script setup>
import {ref} from "vue"

let mov=0
let movnb = 0
let posc = ref(`transform: translateX(${mov}%)`)
let presentation2 = ["servicee", "servicec", "serviced", "servicew"]

function move(dir){
	mov = dir ? mov-=100 : mov+=100
	movnb = dir ? movnb+=1 : movnb-=1
	mov = mov<-300 ? 0 : mov
	mov = mov>0 ? -300 : mov
	movnb = movnb>3 ? 0 : movnb
	movnb = movnb<0 ? 3 : movnb
	posc.value = `transform: translateX(${mov}%)`
}

setInterval(move, 4000, 1)

</script>

<template>		
	<div id="colp1">
					<div>
					<h1>{{ $t('presentation1I')}}</h1><br/><br/>
					<p id="presentation" v-html="$t('presentation1')"></p>
					</div>
					<RouterLink to="/about" id="try">{{$t('try')}}</RouterLink>
					<img id="back" src="@/assets/background.png" />
				</div>
		
		<div id="ens">
				
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


#back{
	
	position: absolute;
	z-index: -1;
	width: 100%;
	height: 100%;
	object-fit: cover;
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

h1{
	display: block;
	font-size: 50px;
	margin-left: 10px;
}

#colp1{
	display: flex;
	position: relative;
	border-radius: 10px;
	overflow: hidden;
	
	align-items: center;
	justify-content: center;
	gap: 50px;
	margin-bottom: 100px;
}

#presentation{
	height: auto;
	width: 500px;
	margin-left: 1%;
	background: rgb(245,245,255, 0.7);
	padding: 10px;
	border-radius: 10px;

}

.dark #presentation{
	background: rgb(0, 0, 0, 0.7);
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

	position: relative;
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
	min-width: 220px;
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

@keyframes anim {
  from {opacity: 0;}
  to {opacity: 1;}
}


</style>
