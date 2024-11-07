package check

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ucl-arc-tre/global/internal/meta"
	"github.com/ucl-arc-tre/global/internal/repo"
)

const (
	versionFilename = "version.txt"
)

type ShouldBeBumped bool

type PackageDirectories map[string]ShouldBeBumped

// Check that if there is a change in a directory then
// if there is a version.txt file then it has been bumped
func CheckVersionBumps(repoRoot string) error {
	dirs := dirsWithVersionFiles(repoRoot)
	for _, filePath := range repo.ModifiedFilePaths(repoRoot) {
		for dir, _ := range dirs {
			if fileExistsInDirectory(filePath, dir) {
				dirs[dir] = ShouldBeBumped(true)
			}
		}
	}
	errorMessage := ""
	for dir, shouldBeBumped := range dirs {
		versionFilepath := path.Join(dir, versionFilename)
		if bool(shouldBeBumped) && !versionIsBumped(repo.FilePatch(repoRoot, versionFilepath)) {
			errorMessage += fmt.Sprintf("[%v/version.txt] has not been bumped", dir)
		}
	}
	if errorMessage != "" {
		log.Fatal(errorMessage)
	}
	return nil
}

func versionIsBumped(patch *repo.Patch) bool {
	if patch == nil || patch.New == "" {
		return false
	}
	if patch.Old == "" { // has been created
		return true
	}
	oldVersion := meta.MakeVersionFromFileContent(patch.Old)
	newVersion := meta.MakeVersionFromFileContent(patch.New)
	return newVersion.IsGreater(oldVersion)
}

func dirsWithVersionFiles(repoRoot string) PackageDirectories {
	dirs := PackageDirectories{}
	_ = filepath.Walk(repoRoot, func(filePath string, info os.FileInfo, err error) error {
		if err == nil && filepath.Base(filePath) == versionFilename {
			dir, _ := filepath.Rel(repoRoot, path.Dir(filePath))
			dirs[dir] = false
		}
		return nil
	})
	return dirs
}

func fileExistsInDirectory(filepath string, dir string) bool {
	return strings.HasPrefix(filepath, dir)
}

func exists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			log.Fatal(err)
		}
	}
	return true
}
