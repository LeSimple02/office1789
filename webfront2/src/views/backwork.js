import {gls} from "@/stores/global"

function verifback(){
if(gls().log==1){
	fetch("http://localhost:8080/api/getinfop", {method:"POST", mode:"cors", body : JSON.stringify({"username" : gls().username, "token" : gls().sessionT}) }).then(a=>a.json()).then(a=>{
		if (a['Username']=="no"){
		gls().log = 0
	}

	})

}
}

export default verifback
