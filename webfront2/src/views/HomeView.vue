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
					<div style="display: flex; align-items: center; gap: 20px;">
					<RouterLink id="try2" to="/createaccount">{{$t('create')}}</RouterLink>
					<RouterLink to="/about" id="try">{{$t('try')}}</RouterLink>
					</div>
					<h2>Leurs privilèges, notre Révolution</h2>
					
					
	</div>
		
		<div id="ens">
				<div id="prevision">
					<ul id="progressbar">
						<li><span class="activeb"></span><span class="active">1</span><p>{{$t('publication')}}</p></li>
						<li><span class="nactiveb"></span><span class="nactive">2</span><p>{{$t('pconv')}}</p></li>
						<li><span class="nactiveb"></span><span class="nactive">3</span><p>{{$t('pdrive')}}</p></li>
						<li><span class="nactiveb"></span><span class="nactive">4</span><p>{{$t('pemail')}}</p></li>
						<li><span class="nactive">5</span><p>{{$t('publis')}}</p></li>
					</ul>
					
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


#progressbar{
	display: flex;
	list-style: none;
	width: auto;
	font-family: roboto;
background: -webkit-linear-gradient(30deg, blue, red);;
	color: white;
	border-radius: 10px;
	
	
	li {
		display: flex;
		flex-direction: column;
		width: 300px;
		align-items: center;
		
		
		
		.nactive {
			width: 20px;
			height: 20px;
			border-radius: 20px;
			background: white;
			border: 1px solid skyblue;
			position: absolute;
			margin-top: -10px;
			text-align: center;
			color: red;
			
		}
		.active {
			background: skyblue;
			width: 20px;
			height: 20px;
			border-radius: 20px;
			border: 1px solid skyblue;
			position: absolute;
			margin-top: -10px;
			text-align: center;
			color: red;
		}
		
		.nactiveb {
			width: 100%;
			height: 1px;
			background: white;
			margin-left: 300px;
			
			
		}
		.activeb {
			width: 100%;
			height: 1px;
			background: red;
			margin-left: 300px;
			
		}
		
		
		
		
	}
	

	
}

#back{
	
	position: absolute;
	z-index: -1;
	width: 100%;
	height: 100%;
	object-fit: cover;
}

#ens {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 100%;
    margin-top: 32px;
    padding: 32px 0 24px 0;
    position: relative;
    background: url("@/assets/napo2.jpg") center/cover no-repeat;
    border-radius: 24px;
    gap: 32px;
    overflow: hidden;
}

#ens::before {
    content: "";
    position: absolute;
    inset: 0;
    background: rgba(255,255,255,0.6);
    backdrop-filter: blur(1px);
    z-index: 0;
    border-radius: 24px;
}

#prevision,
#ensA {
    position: relative;
    z-index: 1;
}




#ens2{
	display: flex;
	flex-wrap: wrap;
	min-height: 35em;
	flex-direction: column;
	gap: 20px;
	text-align: justify;
	
	width: 50%;
	
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


#colp1 {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 220px;
    width: 100%;
    padding: 48px 0 32px 0;
    background: rgba(245,245,247,0.85);
    border-radius: 32px;
    box-shadow: 0 8px 32px 0 rgba(0,0,0,0.10);
    margin-bottom: 32px;
    position: relative;
    overflow: hidden;
    gap: 24px;
}

#colp1 h1 {
    font-family: roboto;
    font-size: 3rem;
    font-weight: 700;
    text-align: center;
    margin-bottom: 12px;
    letter-spacing: 2px;
    color: #222;
}
.dark #colp1 h1{
	color: white;
}
#colp1 h2 {
    font-family: roboto;
    font-size: 1.5rem;
    font-style: italic;
    text-align: center;
    margin-top: 8px;
    background: -webkit-linear-gradient(30deg, blue, red);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    letter-spacing: 1px;
}

#colp1 > div {
    display: flex;
    align-items: center;
    gap: 20px;
    margin-bottom: 8px;
}

#try, #try2 {
	font-family: roboto;
    font-size: 1.1rem;
    padding: 0.75rem 2rem;
    border-radius: 32px;
    border: none;
    cursor: pointer;
	text-decoration: none;
    transition: background 0.3s, color 0.3s, box-shadow 0.3s;
    box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

#try {
    background: transparent;
    border: 2px solid red;
    color: red;
}

#try2 {
    background: #06c;
    color: white;
}

#try:hover,
#try2:hover {
    background: -webkit-linear-gradient(30deg, blue, red);
    color: white;
    box-shadow: 0 4px 16px rgba(0,0,0,0.12);
}


#presentation2 {
    font-family: roboto;
    background: linear-gradient(30deg, blue, red);
    padding: 20px;
    border-radius: 10px;
    color: white;
    box-shadow: 0 8px 32px 0 rgba(0,0,0,0.37); /* Ombre floue */
}


#try:hover,
#try2:hover,
#carrousel button:hover {
    background: -webkit-linear-gradient(30deg, blue, red);
    color: white;
    transition: background 0.3s, color 0.3s;
}

#carrousel{
	display: flex;
	align-items: center;
	justify-content: center;
	height: 300px;
	width: 100%;
	
	
}

#actu{
	
	width: auto;
	min-width: 220px;
	height: auto;
	border-radius: 10px;
	padding: 10px;
	color : white;
	font-family: arial;
	background: -webkit-linear-gradient(30deg, blue, red);;
	
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



#carrousel img {
	max-width: 100%;

}

#carrousel button {
	background : -webkit-linear-gradient(30deg, blue, red);
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

#psf {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 32px;
    padding: 48px 0;
    margin-top: 32px;
    background: rgba(220, 220, 240, 0.95);
    border-radius: 24px;
    box-shadow: 0 8px 32px 0 rgba(0,0,0,0.18);
    position: relative;
    overflow: hidden;
    flex-wrap: wrap;
}

#psf span, #psf img {
    width: 120px;
    height: 100px;
    border-radius: 50px;
    border: 1px solid #bbb;
    box-shadow: 0 2px 8px rgba(0,0,0,0.10);
    object-fit: cover;
    background: linear-gradient(90deg, blue 33%, white 33%, white 66%, red 66%);
}

#psf img {
    background: none;
}

#psf p {
    font-family: roboto;
    font-size: 1.1rem;
    color: #222;
    margin: 0 8px;
    text-align: center;
    font-weight: 500;
    letter-spacing: 1px;
}

@keyframes anim {
  from {opacity: 0;}
  to {opacity: 1;}
}

.dark #psf{
	background: rgba(0, 50, 0, 0.2);
	
}

</style>
