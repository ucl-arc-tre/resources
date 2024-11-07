package script

import (
	"log"
	"os"
	"path"
	"text/template"
)

const (
	buildTemplateFilename = "build.tmpl.sh"
)

func Print(dir string) error {
	assertIsDir(dir)
	templatePath := path.Join(dirParent(dir), buildTemplateFilename)
	t := template.Must(template.ParseFiles(templatePath))
	switch base := path.Base(dirParent(dir)); base {
	case "containers":
		printContainerBuildScript(dir, t)
	default:
		log.Fatalf("unknown package type [%v]", base)
	}
	return nil
}

func assertIsDir(dir string) {
	info, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err)
	} else if !info.IsDir() {
		log.Fatal("Was not directory")
	}
}

func dirParent(dir string) string {
	return path.Dir(dir)
}
