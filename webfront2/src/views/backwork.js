import {gls} from "@/stores/global"
import router from "@/router";

function verifback(){
	const store = gls()
	
	if(store.log == 1){
		// Mettre à jour les valeurs depuis les cookies
		store.updateFromCookies()
		
		// Vérifier que nous avons bien un username et un token
		if(!store.username || !store.sessionT) {
			store.log = 0
			localStorage.setItem("log", 0)
			router.push("/login")
			return
		}
		
		fetch(import.meta.env.VITE_APP_API + '/api/session/check', {
			method: "POST",
			mode: "cors",
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include',
			body: JSON.stringify({
				"username": store.username,
				"token": store.sessionT
			})
		})
		.then(response => {
			if (!response.ok) {
				throw new Error('Session check failed')
			}
			return response.json()
		})
		.then(data => {
			if (!data.connected) {
				// Session invalide
				store.log = 0
				localStorage.setItem("log", 0)
				store.sessionT = ""
				store.username = ""
				document.cookie = "name=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
				document.cookie = "sessionToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
				router.push("/login")
			}
		})
		.catch(err => {
			console.error('Session verification error:', err)
			store.log = 0
			localStorage.setItem("log", 0)
			store.sessionT = ""
			store.username = ""
			document.cookie = "name=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
			document.cookie = "sessionToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
			router.push("/login")
		})
	}
}

export default verifback
