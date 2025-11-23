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

	let username = ref("")
	let sessionT = ref("")
	
	// Fonction pour lire les cookies et mettre à jour les refs
	const updateFromCookies = () => {
		let decodedCookie = decodeURIComponent(document.cookie);
		let co = decodedCookie.split("; ");
		
		for (let i=0; i!=co.length; i++){
			let aco = co[i].split("=")
			
			if(aco[0]=="name"){
				username.value = aco[1]
			}
			else if(aco[0]=="sessionToken"){
				sessionT.value = aco[1]
			}
		}
	}
	
	// Initialiser depuis les cookies
	updateFromCookies()
	
	return {log, lang, username, sessionT, updateFromCookies}
        
    }
)
