import {ref} from 'vue'
import { defineStore } from 'pinia'


export let gls = defineStore('main',()=>
    {   
	let log = ref(0)
	if (localStorage.getItem("log"))
	{
		log.value = localStorage.getItem("log")
	}

	let lang = ref(navigator.language.split('-')[0])
	if (localStorage.getItem("lang"))
	{
		lang.value = localStorage.getItem("lang")
	}

	let username = ref(0)
	let sessionT = ref(0)
	
	
	let decodedCookie = decodeURIComponent(document.cookie);
	let co = decodedCookie.split("; ");
	
	for (let i=0; i!=co.length; i++){
		let aco = co[i].split("=")
		
		if(aco[0]=="name"){
			username = aco[1]
			
		}
		else if(aco[0]=="sessionToken"){
			sessionT = aco[1]
		}
	}
	return {log, lang, username, sessionT}
        
    }
)
