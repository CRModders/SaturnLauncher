package main

import (
	"embed"
	"fmt"
	"math/rand"
	"os"
	"time"

	webview "github.com/webview/webview_go"
)

var (
	edition string
	serverDir string

	LocalAppData   = os.Getenv("localappdata")
	RoamingAppData = os.Getenv("appdata")
	//go:embed win
	win embed.FS
	//go:embed data
	data embed.FS
	
	dir = LocalAppData
	settingsDir = dir + "\\Saturn_Launcher\\settings.json"
	
	port = randInt(255, 9999)
)

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func main() {
	app := webview.New(edition == "dev")
//	defer app.Destroy()

	app.SetSize(900, 570, webview.HintFixed)
	app.Navigate(fmt.Sprintf("http://localhost:%v", port))

	SetApplicationIcon(app, "CLIENT_LOGO")
	BindBuiltInFunctions(app)

	switch edition {
		case "dev":
			app.SetTitle("Saturn Launcher 1.0.0 | Dev")
			go DevServer(
				serverDir,
				fmt.Sprintf("localhost:%v", port),
				)
		default:
			app.SetTitle("Saturn Launcher 1.0.0")
			go EmbededServer(
				win,
				fmt.Sprintf("localhost:%v", port),
				)
	}
	app.Run()
}