import {gls} from "@/stores/global"

function verifback(){
if(gls().log==1){
	fetch(process.env.VUE_APP_API_INFO_USER, {method:"POST", mode:"cors", body : JSON.stringify({"username" : gls().username, "token" : gls().sessionT}) }).then(response=>{
	const contentType = response.headers.
            get('content-type');
        
	if (contentType && contentType.
                includes('application/json')) {
            return response.json();
        }
        else{
		gls().log = 0
		
        
        }}).then(a=>a.json()).then(a=>{
		if (a['Username']=="no"){
		gls().log = 0
	}

	})

}
}

export default verifback
