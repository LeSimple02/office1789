<script setup>




</script>
<template>
	<div id="callarea">
	<img src="@/assets/napo.png" />
	{{$t('callt')}}...
	</div>		
</template>


<style scoped>
#callarea{
	background: skyblue;
	display: flex;
	flex-direction: column;
	font-family: roboto;
	position: absolute;
	left: 0%;
	top: 0%;
	width: 100%;
	height: 100%;
	align-items: center;
	justify-content: center;
	text-align: center;
	img{
		width: 80px;
		height: 80px;
		border-radius:40px; 
	}
	
}

</style>
