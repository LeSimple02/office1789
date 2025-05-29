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
					<h1>{{ $t('presentation1I')}}</h1>
					<h2>Leurs privilèges, notre Révolution</h2>
					<div style="display: flex; align-items: center; gap: 20px;">
					
					<RouterLink id="try2" to="/createaccount">{{$t('create')}}</RouterLink>
					<RouterLink to="/about" id="try">{{$t('try')}}</RouterLink>
					</div>
					
					
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
		<div id="psf">
		<span></span>
		<img src="@/assets/europe.png">
		<p>French/Europe Power</p>
		<p>{{$t('apo')}}</p>
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
	align-items: center;
	width: 50%;
	margin-top: 20px;
	background: #f5f5f7;
	width: 100%;
	
	
	
	
	
}

#psf{

	display: flex;
	padding: 100px;
	justify-content: center;
	align-items: center;
	margin-top: 20px;
	background : rgb(220, 220, 240);
	gap: 20px;
	p{
	font-family: roboto;
	}
	

	span{
		background: linear-gradient(90deg, blue 33%, white 33%, white 66%, red 66%);
		width: 120px;
		height: 100px;
		border: 1px solid grey;
		border-radius: 50px;
	}
	img{
		width: 120px;
		height: 100px;
		
		border: 1px solid grey;
		border-radius: 50px;
	}
}


#ens2{
	display: flex;
	flex-wrap: wrap;
	min-height: 35em;
	flex-direction: column;
	gap: 20px;
	text-align: justify;
	
	width: 55%;
	
}

#ensA{
	display: flex;
	flex-wrap: wrap;
	gap: 80px;
	width: 50%;
	
	justify-content: center;
}

h1{
	display: block;
	text-align: center;
	font-size: 50px;
	
	
}

h2{
	display: block;
	text-align: center;
	
	
        font-family: roboto;
	background : -webkit-linear-gradient(30deg, blue, red);
	font-style: italic;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
	
}

#colp1{
	display: grid;
	

	
	background: #f5f5f7;
	width: 100%;
	
	padding-bottom: 50px;
	
	justify-content: center;
	align-items: center;
	
	
	
}

#presentation2{
	font-family: roboto;
}

#try{
	border: 2px solid red;
	color: red;
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

#try2{
	border: 2px solid #06c;
	color: #06c;
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



#carrousel{
	display: flex;
	align-items: center;
	justify-content: center;
	height: 300px;
	width: 100%;
	
	
}

#actu{
	
	width: 30%;
	min-width: 220px;
	height: auto;
	padding: 10px;
	font-family: arial;
	background: rgb(240, 200, 200);
	
}


#actu > h3 {
	font-family: roboto;
	border-bottom: 1px solid black;
}

.dark #actu > h3 {
	border-bottom: 1px solid white;
}

.dark #colp1{
	background: #1C1C1E;
}

.dark #ens{
	background: #1C1C1E;
}


#carrousel img {
	width: 100%;

}

#carrousel button {
	background: black;
	border-radius: 20px;
	color: white;
	font-style: bold;
	
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
