import subprocess
import webbrowser
import docker




if __name__ == "__main__":
    print("Start of Postgresql docker")
    client = docker.from_env()
    container = client.containers.get("postgres_db")
    container.start()
    print("Start of Backend go")
    subprocess.Popen("air", cwd="backend")
    print("Start of webfront")
    process = subprocess.Popen("npm run dev", cwd="webfront2", shell=True)
    process.wait()
    print("Open Browser")
    webbrowser.open("http://localhost:5173")
    """
    J'attends la suite pour construire une config très pro
    mode =str()
    lang = ["en", "fr", "es"]
    ilang = int()
    print("Language :")
    for i in lang:
        print(i)
    mode = input("")
    ilang = lang.index(mode)
    print("Do you have a DB")
    """
    