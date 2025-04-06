import platform
import os
import distro




def lconfig():
    nameos = distro.os_release_info()['name']
    if nameos == "Fedora Linux":
    	os.system("""sudo dnf install nodejs
    	npm install vue
    	sudo dnf install cargo
    	""")
    
    



