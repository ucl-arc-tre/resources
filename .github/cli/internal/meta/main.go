package meta

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
}

func (v Version) IsGreater(other Version) bool {
	if v.Major > other.Major {
		return true
	} else if v.Major == other.Major && v.Minor > other.Minor {
		return true
	} else if v.Major == other.Major && v.Minor == other.Minor && v.Patch > other.Patch {
		return true
	} else {
		return false
	}
}

func BundleFilename(dir string) string {
	name := path.Base(dir)
	version := makeVersionFromDir(dir)
	return fmt.Sprintf("%s_%v_%v_%v.tar.gz", name, version.Major, version.Minor, version.Patch)
}

func Image(dir string) string {
	app := path.Base(dir)
	version := makeVersionFromDir(dir)
	return fmt.Sprintf("%s:%v.%v.%v", app, version.Major, version.Minor, version.Patch)
}

func MakeVersionFromFileContent(content string) Version {
	parts := strings.Split(content, ".")
	if len(parts) != 3 {
		log.Fatal("Version had the wrong structure. Must be semver")
	}
	version := Version{
		Major: mustStringToInt(parts[0]),
		Minor: mustStringToInt(parts[1]),
		Patch: mustStringToInt(parts[2]),
	}
	return version
}

func makeVersionFromDir(dir string) Version {
	filePath := path.Join(dir, "version.txt")
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return MakeVersionFromFileContent(string(fileBytes))

}

func mustStringToInt(s string) int {
	if i, err := strconv.Atoi(strings.TrimSuffix(s, "\n")); err == nil {
		return i
	} else {
		log.Fatal(err)
		return 0
	}
}
