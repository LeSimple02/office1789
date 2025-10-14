import {gls} from "@/stores/global"
import router from "@/router";

function verifback(){



if(gls().log==1){
	
	fetch(import.meta.env.VITE_APP_API_INFO_USER, {method:"POST", mode:"cors", body : JSON.stringify({"Username" : gls().username, "Token" : gls().sessionT}) }).then(response=>{
	const contentType = response.headers.get('content-type');
	if (contentType && contentType.includes('application/json')) {
			
			return response.json();
        }
		
        else{
		gls().log = 0
		localStorage.setItem("log", 0)
		
        
        }}).then(a=>{
			
		if (a['Username']=='no'){
			
		gls().log = 0
		
		router.push("/")
		
	}

	}) .catch(err => {
		localStorage.removeItem("log")
		gls().sessionT = ""
			gls().username = ""
			document.cookie = "name=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
			document.cookie = "sessionToken=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
			router.push("/")
			
		}
		)

}
}

export default verifback
