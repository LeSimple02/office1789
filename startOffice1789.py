import subprocess
import webbrowser


if __name__ == "__main__":
    print("Start of Postgresql docker")
    process = subprocess.Popen("docker start focused_bhabha")
    process.wait()
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
    