package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"encoding/json"
	"strings"

	"github.com/vincent-petithory/dataurl"
	webview "github.com/webview/webview_go"
)

var UseCrashableLogger = false

type LogType int
const (
	LOG_FATAL LogType = iota
	LOG_INFO
)

func LogDynamic(args ...interface{}) {
	Log(LOG_INFO, "SL", args...)
}

func Log(logType LogType, name string, args ...interface{}) {
	ARG := ""
	for _, i := range args {
		ARG += fmt.Sprintf("%v", i)
	}

	fmt.Printf("[%v-INFO] %v\n", name, ARG)

}

type Settings struct {
	BackgroundType string `json:"backgroundType"`
	Theme          string `json:"theme"`
}

func readSettings(filename string) (Settings, error) {
    settings := Settings{} // Default values in case the file is new

    content, err := ioutil.ReadFile(filename)
    if err != nil {
        return settings, nil // Return defaults without an error if the file doesn't exist
    }

    err = json.Unmarshal(content, &settings)
    if err != nil {
        return settings, err
    }

    return settings, nil
}

func writeSettings(filename string, settings Settings) error {
    jsonBytes, err := json.MarshalIndent(settings, "", "    ") 
    if err != nil {
        return err
    }

    return os.WriteFile(filename, jsonBytes, 0644)
}

func BindBuiltInFunctions(app webview.WebView) {
	BindFunc := func(FuncName string, Function interface{}) {
		LogDynamic(
			fmt.Sprintf(
				"Bound Function: %v", FuncName,
				),
			)
		err := app.Bind(FuncName, Function)
		if err != nil {
			LogDynamic(
				fmt.Sprintf("Failed To Bind Function %v; Err: %v",
					FuncName,
					err),
				)
		}
	}

	BindFunc("GetSettings", func() Settings {
        settings, err := readSettings(settingsDir)
        if err != nil {
            LogDynamic(err)
        }
        return settings
    })

    BindFunc("SaveSettings", func(newSettings Settings) {
        err := writeSettings(settingsDir, newSettings)
        if err != nil {
            LogDynamic(err)
        }
    })

	BindFunc("GetAppdataPath", func(isLocal bool) string {
		if isLocal {
			return LocalAppData
		}
		return RoamingAppData
	})

	BindFunc("IsDevEnviornment", func() bool {
		return edition == "dev"
	})

	BindFunc("GetServerPort", func() int {
		return port
	})

	BindFunc("UseCrashableLogger", func(isCrashable bool) {
		UseCrashableLogger = isCrashable
	})

	BindFunc("IsUsingCrashableLogger", func() bool {
		return UseCrashableLogger
	})

	BindFunc("Log", func(name string, args ...string) {
		ARG := ""
		for _, i := range args {
			ARG += fmt.Sprintf("%v", i)
		}
		fmt.Printf("[%v-INFO] %v\n", name, ARG)
	})

	BindFunc("CreateDirectory", func(DirName string) {
		err := os.MkdirAll(DirName, fs.ModeDir)
		if err != nil { LogDynamic(err) }
	})

	BindFunc("RemoveDirectory", func(DirName string) {
		err := os.RemoveAll(DirName)
		if err != nil { LogDynamic(err) }
	})

	BindFunc("Rename", func(OldPath, NewPath string) {
		err := os.Rename(OldPath, NewPath)
		if err != nil { LogDynamic(err) }
	})

	BindFunc("PathExists", func(path string) bool {
		_, err := os.Stat(path)
		return err == nil
	})

	BindFunc("WriteToFile", func(Name, Contents string) {
		err := os.WriteFile(Name, []byte(Contents), fs.FileMode(os.O_CREATE))
		if err != nil { LogDynamic(err) }
	})

	BindFunc("RemoveFile", func(FilePath string) {
		err := os.Remove(FilePath)
		if err != nil { LogDynamic(err) }
	})

	BindFunc("ReadFile", func(FilePath string) string {
		content, err := os.ReadFile(FilePath)
		if err != nil { LogDynamic(err) }
		return string(content)
	})

	BindFunc("CopyFileTo", func(originFile, newFile string) {
		content, err := os.ReadFile(originFile)
		if err != nil { LogDynamic(err) }
		err = os.MkdirAll(path.Dir(newFile), fs.ModeDir)
		if err != nil { LogDynamic(err) }
		err = os.WriteFile(newFile, content, fs.ModePerm)
		if err != nil { LogDynamic(err) }
	})

	BindFunc("FileToDataURL", func(FilePath string) string {
		content, err := os.ReadFile(FilePath)
		if err != nil { LogDynamic(err) }
		return dataurl.EncodeBytes(content)
	})

	BindFunc("WalkDir", func(dir string) []string {
		dirs := make([]string, 1)
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if dir != path && dir != "" {
				dirs = append(dirs, path)
			}
			return err
		})
		if err != nil { LogDynamic(err) }
		return dirs
	})

	BindFunc("ExtractZip", func(file string, folder string) {
		zr, err := zip.OpenReader(file)
		if err != nil { LogDynamic(err) }
		for _, f := range zr.File {
			fbuf, _ := zr.Open(f.Name)
			if fbuf != nil {
				fbytes, _ := ioutil.ReadAll(fbuf)
				fmt.Printf("%v Extracted to %v\n", f.Name, path.Join(folder, f.Name))
				os.MkdirAll(path.Dir(path.Join(folder, f.Name)), fs.ModeDir)
				fz, _ := os.Create(path.Join(folder, f.Name))
				fmt.Println(path.Dir(path.Join(folder, f.Name)))
				fz.Write(fbytes)
				fz.Close()
			}
		}
		zr.Close()
	})

	BindFunc("FAUTO_ExtractZip", func(file string, folder string) {
		os.MkdirAll(folder, os.ModeDir)
		zr, err := zip.OpenReader(file)
		if err != nil { LogDynamic(err) }
		for _, f := range zr.File {
			fbuf, _ := zr.Open(f.Name)
			if fbuf != nil {
				fbytes, _ := ioutil.ReadAll(fbuf)
				fmt.Printf("%v Extracted to %v\n", f.Name, path.Join(folder, f.Name))
				os.MkdirAll(path.Dir(path.Join(folder, f.Name)), fs.ModeDir)
				fz, _ := os.Create(path.Join(folder, f.Name))
				fmt.Println(path.Dir(path.Join(folder, f.Name)))
				fz.Write(fbytes)
				fz.Close()
			}
		}
		zr.Close()
		os.Remove(file)
	})

	BindFunc("ReadEntryInsideZip", func(file string, entry string) string {
		zr, err := zip.OpenReader(file)
		if err != nil { LogDynamic(err) }
		fi, err := zr.Open(entry)
		if err != nil { LogDynamic(err) }
		bytes, _ := ioutil.ReadAll(fi)
		zr.Close()
		return string(bytes)
	})

	BindFunc("EMBED_PathExists", func(FilePath string) bool {
		EmbededFileSystem, err := fs.Sub(
			data,
			path.Join("data", FilePath),
			)
		if err != nil { LogDynamic(err) }
		_, FileName := path.Split(FilePath)
		_, err = fs.Stat(EmbededFileSystem, FileName)
		return err == nil
	})

	BindFunc("EMBED_ReadFile", func(FilePath string) string {
		EmbededFileSystem, err := fs.Sub(
			data,
			path.Join("data", FilePath),
			)
		if err != nil { LogDynamic(err) }
		_, FileName := path.Split(FilePath)
		content, err := fs.ReadFile(EmbededFileSystem, FileName)
		if err != nil { LogDynamic(err) }
		return string(content)
	})

	BindFunc("EMBED_CopyFileTo", func(originFile, newFile string) {
		EmbededFileSystem, err := fs.Sub(
			data,
			path.Join("data", originFile),
			)
		if err != nil { LogDynamic(err) }
		content, err := fs.ReadFile(EmbededFileSystem, originFile)
		if err != nil { LogDynamic(err) }
		err = os.MkdirAll(path.Dir(newFile), fs.ModeDir)
		if err != nil { LogDynamic(err) }
		err = os.WriteFile(newFile, content, fs.ModePerm)
		if err != nil { LogDynamic(err) }
	})

	BindFunc("EMBED_FileToDataURL", func(originFile, newFile string) string {
		EmbededFileSystem, err := fs.Sub(
			data,
			path.Join("data", originFile),
			)
		if err != nil { LogDynamic(err) }
		_, FileName := path.Split(originFile)
		content, err := fs.ReadFile(EmbededFileSystem, FileName)
		if err != nil { LogDynamic(err) }
		return dataurl.EncodeBytes(content)
	})

	BindFunc("ReadFileURLToText", func(URL string) string {
		resp, err := http.Get(URL)
		if err != nil { LogDynamic(err) }
		content, _ := ioutil.ReadAll(resp.Body)
		if err != nil { LogDynamic(err) }
		return string(content)
	})

	BindFunc("DownloadFileFromURL", func(URL string, filepath string) {
		out, err := os.Create(strings.Trim(filepath, "\n\r"))
		if err != nil { LogDynamic(err) }
		defer out.Close()
		resp, err := http.Get(strings.Trim(strings.ReplaceAll(URL, " ", "%20"), "\n\r"))
		if err != nil { LogDynamic(err) }
		defer resp.Body.Close()
		io.Copy(out, resp.Body)
	})

	BindFunc("ExecuteCommand", func(cwd string, prg string, args ...string) {
		cmd := exec.Command(prg, args...)
		cmd.Dir = cwd
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
		}
		err = cmd.Start()
		fmt.Println("The command is running")
		if err != nil {
			fmt.Println(err)
		}
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Printf("%v\n", m)
				app.Dispatch(func() {
				})
			}
			cmd.Wait()
		}()

	})

	BindFunc("Open", func(f string) {
		err := open.Run(f)
		if err != nil { LogDynamic(err) }
	})
}