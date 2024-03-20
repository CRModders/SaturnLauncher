import os

def remove_file(file_path):
    if os.path.exists(file_path):
        os.remove(file_path)

def copy_file(source, destination):
    with open(source, "rb") as src_file, open(destination, "wb") as dest_file:
        dest_file.write(src_file.read())

# os.system("go clean -modcache")
# os.system("go mod tidy")
os.system("go get")

# Compiling Dev - Go
print("Compiling Dev - Go")
os.chdir("src")

remove_file("win/favicon.ico")
remove_file("winres/icon.png")

copy_file("data/img/dev_logo.ico", "win/favicon.ico")
copy_file("data/img/dev_logo.png", "winres/icon.png")

os.system("go install github.com/tc-hib/go-winres@latest")
os.system("go-winres make")
os.system("go build -ldflags=\"-X main.edition=dev -X main.serverDir={0}\" -o ../bin/Saturn_launcher_Dev.exe".format(os.path.join(os.getcwd(), "win")))

# Compiling Stable - Go
print("Compiling Stable - Go")

remove_file("win/favicon.ico")
remove_file("winres/icon.png")

copy_file("data/img/logo.ico", "win/favicon.ico")
copy_file("data/img/logo.png", "winres/icon.png")

os.system("go install github.com/tc-hib/go-winres@latest")
os.system("go-winres make")
os.system("go build -ldflags=\"-H windowsgui\" -o ../bin/Saturn_Launcher.exe")

remove_file("rsrc_windows_386.syso")
remove_file("rsrc_windows_amd64.syso")

os.chdir("../bin")
os.system("start Saturn_launcher.exe")
