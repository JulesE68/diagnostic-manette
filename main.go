//go:build windows

// Diagnostic Manette : fenêtre Windows autonome qui affiche l'interface de test
// (ui/app.html) dans WebView2, le moteur Edge intégré à Windows 10 et 11.
package main

import (
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jchv/go-webview2"
	"golang.org/x/sys/windows"
)

//go:embed ui/app.html
var appHTML string

func messageBox(title, text string, flags uint32) {
	t, _ := windows.UTF16PtrFromString(title)
	m, _ := windows.UTF16PtrFromString(text)
	windows.MessageBox(0, m, t, flags)
}

// exeDir renvoie le dossier du .exe, pour enregistrer les rapports à côté (usage portable).
func exeDir() string {
	p, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(p)
}

func main() {
	// Données du moteur web dans le dossier temporaire : rien n'est écrit à côté du .exe
	// ni dans AppData en dehors de ce dossier jetable.
	dataPath := filepath.Join(os.TempDir(), "DiagnosticManette-WebView2")

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     false,
		AutoFocus: true,
		DataPath:  dataPath,
		WindowOptions: webview2.WindowOptions{
			Title:  "Diagnostic Manette Xbox",
			Width:  1180,
			Height: 800,
			IconId: 1,
			Center: true,
		},
	})
	if w == nil {
		messageBox("Diagnostic Manette",
			"Impossible de démarrer l'interface.\n\nCe programme a besoin de Microsoft Edge WebView2 Runtime, "+
				"déjà installé sur Windows 11 et sur la plupart des Windows 10.\n\n"+
				"Si besoin, installe-le depuis :\nhttps://developer.microsoft.com/microsoft-edge/webview2/",
			windows.MB_ICONERROR)
		return
	}
	defer w.Destroy()

	w.SetSize(1000, 700, webview2.HintMin)

	// Enregistre le rapport texte à côté du .exe et renvoie le chemin complet.
	_ = w.Bind("saveReport", func(text string) (string, error) {
		name := "Rapport manette " + time.Now().Format("2006-01-02 15h04") + ".txt"
		path := filepath.Join(exeDir(), name)
		// BOM UTF-8 pour que le Bloc-notes affiche bien les accents.
		data := append([]byte{0xEF, 0xBB, 0xBF}, []byte(strings.ReplaceAll(text, "\r\n", "\n"))...)
		data = []byte(strings.ReplaceAll(string(data), "\n", "\r\n"))
		if err := os.WriteFile(path, data, 0o644); err != nil {
			// Dossier protégé en écriture (ex. Program Files) : repli sur les Documents.
			home, _ := os.UserHomeDir()
			alt := filepath.Join(home, "Documents", name)
			if err2 := os.WriteFile(alt, data, 0o644); err2 != nil {
				return "", fmt.Errorf("%v", err)
			}
			return alt, nil
		}
		return path, nil
	})

	// Ouvre le rapport avec l'application par défaut (Bloc-notes en général).
	_ = w.Bind("openPath", func(path string) {
		go func() {
			verb, _ := windows.UTF16PtrFromString("open")
			file, _ := windows.UTF16PtrFromString(path)
			_ = windows.ShellExecute(0, verb, file, nil, nil, windows.SW_SHOWNORMAL)
		}()
	})

	// Ouvre l'Explorateur avec le rapport sélectionné.
	_ = w.Bind("showInFolder", func(path string) {
		go func() { _ = exec.Command("explorer.exe", "/select,", path).Start() }()
	})

	_ = w.Bind("quitApp", func() {
		w.Dispatch(func() { w.Terminate() })
	})

	// Page servie en local (127.0.0.1 uniquement) : contexte sécurisé, donc accès
	// garanti à l'API Gamepad, aux vibrations et au presse-papiers.
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		w.SetHtml(appHTML)
	} else {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				http.NotFound(rw, r)
				return
			}
			rw.Header().Set("Content-Type", "text/html; charset=utf-8")
			rw.Header().Set("Cache-Control", "no-store")
			_, _ = rw.Write([]byte(appHTML))
		})
		go func() { _ = http.Serve(ln, mux) }()
		w.Navigate("http://" + ln.Addr().String() + "/")
	}
	w.Run()
}

