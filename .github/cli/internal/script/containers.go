package script

import (
	"log"
	"os"
	"text/template"

	"github.com/ucl-arc-tre/global/internal/meta"
)

type ContainerBuildParams struct {
	Image    string
	FileName string
	Context  string
}

func printContainerBuildScript(dir string, t *template.Template) {
	err := t.Execute(os.Stdout, ContainerBuildParams{
		Image:    meta.Image(dir),
		FileName: meta.BundleFilename(dir),
		Context:  dir,
	})
	if err != nil {
		log.Fatal(err)
	}
}
