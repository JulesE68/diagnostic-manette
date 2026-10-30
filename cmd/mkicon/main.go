// Génère rsrc_windows_amd64.syso (icône + manifeste) à partir de ui/icon.png.
package main

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/tc-hib/winres"
	"github.com/tc-hib/winres/version"
)

func main() {
	f, err := os.Open("ui/icon.png")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	ico, err := winres.NewIconFromResizedImage(img, nil)
	if err != nil {
		log.Fatal(err)
	}
	rs := winres.ResourceSet{}
	if err := rs.SetIcon(winres.ID(1), ico); err != nil {
		log.Fatal(err)
	}
	rs.SetManifest(winres.AppManifest{
		DPIAwareness:        winres.DPIPerMonitorV2,
		UseCommonControlsV6: true,
		ExecutionLevel:      winres.AsInvoker,
	})
	v := version.Info{}
	v.SetFileVersion("1.0.0.0")
	v.SetProductVersion("1.0.0.0")
	v.Set(0, version.ProductName, "Diagnostic Manette")
	v.Set(0, version.FileDescription, "Diagnostic Manette Xbox")
	v.Set(0, version.OriginalFilename, "DiagnosticManette.exe")
	rs.SetVersionInfo(v)
	out, err := os.Create("rsrc_windows_amd64.syso")
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	if err := rs.WriteObject(out, winres.ArchAMD64); err != nil {
		log.Fatal(err)
	}
}
