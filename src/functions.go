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
//	if logType == LOG_FATAL && UseCrashableLogger {
//		log.Printf("[%v-FATAL] %v", name, ARG)
//		return
//	}
	fmt.Printf("[%v-INFO] %v\n", name, ARG)
//	log.Printf("[%v-INFO] %v", name, ARG)

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

	// Data
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

	// Logging
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

	// Path Manipulation
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

	// File Manipulation
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

	// Embeded Directory And File Manipulation
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
//		_, FileName := path.Split(originFile)
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

	// URL_MANIPULATION
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

	// MISC
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
		// Realtime std::out logging
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				m := scanner.Text()
				fmt.Printf("%v\n", m)
				app.Dispatch(func() {
					//					w.Eval(fmt.Sprintf(`
					//z3 = document.createElement("pre"); z3.innerText = "%v"; z3.setAttribute("class", "log");
					//document.getElementsByClassName("console")[0].append(z3)
					//if (Math.round(document.getElementsByClassName("console")[0].scrollTop) >= (document.getElementsByClassName("console")[0].scrollHeight - 859)) {
					//    document.getElementsByClassName("console")[0].scroll(0, document.getElementsByClassName("console")[0].scrollHeight)
					//}
					//`, m))
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

//func bindFunctions(w webview.WebView) {
//	bind := func(name string, fn interface{}) {
//		if err := w.Bind(name, fn); err != nil {
//			fmt.Printf("Failed to bind %s: %v\n", name, err)
//		}
//	}
//
//	bind("makeDir", func(dir string) {
//		fmt.Println(os.Mkdir(dir, fs.ModeDir))
//	})
//
//	bind("readUrlText", func(url string) string {
//		resp, _ := http.Get(url)
//		bytes, _ := ioutil.ReadAll(resp.Body)
//		return string(bytes)
//	})
//
//	bind("readInsideZip", func(file string, entry string) string {
//		zr, _ := zip.OpenReader(file)
//		fi, _ := zr.Open(entry)
//		if (fi == nil) {
//			return ""
//		}
//		bytes, _ := ioutil.ReadAll(fi)
//		zr.Close()
//		return string(bytes)
//	})
//
//	bind("extractZip", func(file string, folder string) {
//		zr, _ := zip.OpenReader(file)
//		fmt.Println("------------")
//		fmt.Println(file)
//		fmt.Println("------------")
//		for _, f := range zr.File {
//			fbuf, _ := zr.Open(f.Name)
//			if fbuf != nil {
//				fbytes, _ := ioutil.ReadAll(fbuf)
//				fmt.Printf("%v Extracted to %v\n", f.Name, path.Join(folder, f.Name))
//				os.MkdirAll(path.Dir(path.Join(folder, f.Name)), fs.ModeDir)
//				fz, _ := os.Create(path.Join(folder, f.Name))
//				fmt.Println(path.Dir(path.Join(folder, f.Name)))
//				fz.Write(fbytes)
//				fz.Close()
//			}
//		}
//		zr.Close()
//	})
//
//	bind("open_run", func(f string) {
//		open.Run(f)
//	})
//
//	bind("extractZip_COMPL", func(file string, folder string) {
//		os.MkdirAll(folder, os.ModeDir)
//		zr, _ := zip.OpenReader(file)
//		for _, f := range zr.File {
//			fbuf, _ := zr.Open(f.Name)
//			fbytes, _ := ioutil.ReadAll(fbuf)
//			fmt.Printf("%v Extracted to %v\n", f.Name, path.Join(folder, f.Name))
//			fz, _ := os.Create(path.Join(folder, f.Name))
//			fz.Write(fbytes)
//			fz.Close()
//		}
//		zr.Close()
//		os.Remove(file)
//	})
//
//	bind("downloadUrl", func(file, url string) string {
//		out, _ := os.Create(strings.Trim(file, "\n\r"))
//		defer out.Close()
//
//		// Get the data
//		resp, _ := http.Get(strings.Trim(strings.ReplaceAll(url, " ", "%20"), "\n\r"))
//		defer resp.Body.Close()
//		// Write the body to file
//		_, _ = io.Copy(out, resp.Body)
//		return file
//	})
//
//	bind("WalkDir", func(dir string) []string {
//		dirs := make([]string, 1)
//		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
//			if dir != path && dir != "" {
//				dirs = append(dirs, path)
//			}
//			return err
//		})
//		if err != nil {
//			fmt.Println(err)
//		}
//		return dirs
//	})
//
//	bind("print", func(str any) { fmt.Println(str) })
//
//	bind("FileExist", func(filePath string) bool {
//		_, err := os.Stat(filePath)
//		return err == nil
//	})
//	bind("Embed_FileExist", func(filePath string) bool {
//		fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)))
//		_, fileName := path.Split(filePath)
//		_, err := fs.Stat(fileSys, fileName)
//		return err == nil
//	})
//
//	bind("ReadFile", func(filePath string) string {
//		content, _ := os.ReadFile(filePath)
//		return string(content)
//	})
//	bind("Embed_ReadFile", func(filePath string) string {
//		fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)))
//		_, fileName := path.Split(filePath)
//		content, _ := fs.ReadFile(fileSys, fileName)
//		return string(content)
//	})
//
//	bind("ReadFileAsDataUrl", func(filePath string) string {
//		content, _ := os.ReadFile(filePath)
//		return dataurl.EncodeBytes(content)
//	})
//	bind("Embed_ReadFileAsDataUrl", func(filePath string) string {
//		fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)))
//		_, fileName := path.Split(filePath)
//		content, _ := fs.ReadFile(fileSys, fileName)
//		return dataurl.EncodeBytes(content)
//	})
//
//	bind("CopyFile", func(filePath, newFilePath string) {
//		os.MkdirAll(path.Dir(newFilePath), fs.FileMode(os.O_CREATE))
//		z, _ := os.ReadFile(filePath)
//		os.WriteFile(newFilePath, z, fs.FileMode(os.O_CREATE))
//	})
//	bind("Embed_CopyFile", func(filePath, newFilePath string) {
//		os.MkdirAll(path.Dir(newFilePath), fs.FileMode(os.O_CREATE))
//		fileSys, _ := fs.Sub(data, path.Join("data", path.Dir(filePath)))
//		_, fileName := path.Split(filePath)
//		content, _ := fs.ReadFile(fileSys, fileName)
//		os.WriteFile(newFilePath, content, fs.FileMode(os.O_CREATE))
//	})
//
//	bind("RemoveFile", func(filePath string) { os.Remove(filePath) })
//	bind("RenameFile", func(filePath, newFilePath string) { os.Rename(filePath, newFilePath) })
//	bind("RemoveDir", func(path string) { os.RemoveAll(path) })
//
//	bind("WriteFile", func(filePath, content string) {
//		os.MkdirAll(path.Dir(filePath), fs.FileMode(os.O_CREATE))
//		q, _ := os.Create(filePath)
//		q.WriteString(content)
//		q.Close()
//	})
//	bind("LocalAppdata", func() string { return LocalAppData })
//	bind("RoamingAppdata", func() string { return RoamingAppData })
//	bind("edition", func() string { return edition })
//	bind("port", func() int { return port })
//	bind("execute", func(cwd string, prg string, args ...string) {
//		cmd := exec.Command(prg, args...)
//		cmd.Dir = cwd
//		stdout, err := cmd.StdoutPipe()
//		if err != nil {
//			fmt.Println(err)
//		}
//		err = cmd.Start()
//		fmt.Println("The command is running")
//		if err != nil {
//			fmt.Println(err)
//		}
//		// Realtime std::out logging
//		go func() {
//			scanner := bufio.NewScanner(stdout)
//			for scanner.Scan() {
//				m := scanner.Text()
//				fmt.Printf("%v\n", m)
//				w.Dispatch(func() {
////					w.Eval(fmt.Sprintf(`
////z3 = document.createElement("pre"); z3.innerText = "%v"; z3.setAttribute("class", "log");
////document.getElementsByClassName("console")[0].append(z3)
////if (Math.round(document.getElementsByClassName("console")[0].scrollTop) >= (document.getElementsByClassName("console")[0].scrollHeight - 859)) {
////    document.getElementsByClassName("console")[0].scroll(0, document.getElementsByClassName("console")[0].scrollHeight)
////}
////`, m))
//				})
//			}
//			cmd.Wait()
//		}()
//
//	})
//
//}
