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

	if(document.cookie.split(';')[0])
	{
		username = document.cookie.split(';')[0].split('=')[1] 
	}

	let sessionT = ref(0)
	if(document.cookie.split(';')[1])
	{
		sessionT = document.cookie.split(';')[1].split('=')[1];
	}

	return {log, lang, username, sessionT}
        
    }
)
