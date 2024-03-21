package main

/*
#include <windows.h>

void set_resource_icon(const void *ptr, char* name) {
	HINSTANCE hInstance = GetModuleHandle(NULL);
	HICON iconBig = (HICON)LoadImage(hInstance, name, IMAGE_ICON, GetSystemMetrics(SM_CXICON), GetSystemMetrics(SM_CXICON), LR_DEFAULTCOLOR);
	HICON iconSml = (HICON)LoadImage(hInstance, name, IMAGE_ICON, GetSystemMetrics(SM_CXSMICON), GetSystemMetrics(SM_CYSMICON), LR_DEFAULTCOLOR);
	if (iconSml) SendMessage((HWND)ptr, WM_SETICON, ICON_SMALL, (LPARAM)iconSml);
	if (iconBig) SendMessage((HWND)ptr, WM_SETICON, ICON_BIG, (LPARAM)iconBig);
}

void set_external_icon(const void *ptr, char* iconPath) {
	HICON iconBig = LoadImage(NULL, iconPath, IMAGE_ICON, GetSystemMetrics(SM_CXICON), GetSystemMetrics(SM_CXICON), LR_LOADFROMFILE);
	HICON iconSml = LoadImage(NULL, iconPath, IMAGE_ICON, GetSystemMetrics(SM_CXSMICON), GetSystemMetrics(SM_CXSMICON), LR_LOADFROMFILE);
	if (iconSml) SendMessage((HWND)ptr, WM_SETICON, ICON_SMALL, (LPARAM)iconSml);
	if (iconBig) SendMessage((HWND)ptr, WM_SETICON, ICON_BIG, (LPARAM)iconBig);
}
*/
import "C"
import "unsafe"
import webview "github.com/webview/webview_go"

func SetApplicationIcon(app webview.WebView, name string)  {
	hwnd := app.Window()
	cstr := C.CString(name)
	defer C.free(unsafe.Pointer(cstr))
	C.set_resource_icon(hwnd, cstr)
}

func SetApplicationIconFromExternal(app webview.WebView, path string) {
	hwnd := app.Window()
	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))
	C.set_external_icon(hwnd, cstr)
}