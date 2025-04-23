<script setup>



const drives = document.createElement("script")
drives.setAttribute("src", "http://localhost/web-apps/apps/api/documents/api.js")


const config = {
  document: {
    fileType: "docx",
    key: "Khirz6zTPdfd7",
    title: "Example Document Title.docx",
    url: "https://example.com/url-to-example-document.docx",
  },
  documentType: "word",
  editorConfig: {
    callbackUrl: "https://example.com/url-to-callback.ashx",
  },
};

const docEditor = new drives.DocsAPI.DocEditor("placeholder", config);
console.log(docEditor)

/*
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
*/

</script>
<template>

	<div id="placeholder"></div>



</template>
<style scoped>

iframe{
	position: absolute;
	left:0%;
	overflow: hidden;
	border: none;
	width: 100%;
	height: 100%;
}


</style>
