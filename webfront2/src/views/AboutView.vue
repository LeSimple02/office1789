<script setup lang="ts">

</script>
<template>
	<div id="enso">
	<div class="about" v-html="$t('aboutd')"></div>
	<div id="offer" v-html="$t('offer')"></div>
	<RouterLink to="/createaccount">{{$t('createac')}}</RouterLink>
  </div>
</template>

<style scoped>

.about > h1{
	font-family: roboto;
	width: 300px;

}

.about{
	font-family: arial;
}

#enso{
	align-items: center;
	display: flex;
	flex-direction: column;
	height: auto;
}

#enso > a{
	
	background: black;
	color: white;
	font-family: arial;
	text-decoration: none;
	border-radius: 5px;
	font-size: 20px;
	width: auto;
	max-width: 300px;
	text-align: center;
	margin-top: 50px;
	padding: 20px;
}

.dark #enso > a{
	
	background: white;
	color: black;
	
}




#offer{
	display: flex;
	flex-wrap: wrap;
	align-items: center;
	justify-content: center;
	margin-top: 60px;
	gap: 100px;
	
}

#offer ::v-deep{

.offerm{
	border: 2px solid grey;
	width: 270px;
	padding: 30px;
	height: auto;
	font-family: arial;
	background: transparent;
	transition: 0.5s ease-out;

}

.offerm:hover{
	border: 2px solid black;
}

.offerm > h3{
	text-align: center;
	padding-top: 10px;
	width: 330px;
	margin-top: -30px;
	height: 40px;
	margin-left: -30px;
	background: rgba(255, 0,0,0.5);
}
.offerm:hover > h3{
	transition: 0.5s ease-out;
	background: rgba(255, 0,0,0.8);
}

.dark .offerm{
	border: 2px solid grey;
	

}

.dark .offerm:hover{
	border: 2px solid white;
}


.lioffer{
	margin-top: 30px;
	list-style: none;
	font-size: 20px;
	text-align: left;
}

.lioffer > li{
	margin-top: 30px;
}

}

</style>
