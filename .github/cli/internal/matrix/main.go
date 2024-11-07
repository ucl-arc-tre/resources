package matrix

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	defaultBranchName = "main"
	packageRegex      = `^((?:containers|snaps)\/[\w_-]+)\/.*$`
)

type MatrixElement struct {
	Directory string `json:"directory"`
}

type GitHubMatrix struct {
	Include []MatrixElement `json:"include"`
}

func Print(repoRoot string) error {
	elements := map[string]MatrixElement{}
	for _, path := range modifiedFilePaths(repoRoot) {
		if dir := l2DirectoryFromFile(path); dir != nil {
			elements[*dir] = MatrixElement{
				Directory: *dir,
			}
		}
	}
	matrix := GitHubMatrix{
		Include: values(elements),
	}
	marshaledElements, err := json.Marshal(matrix)
	if err == nil {
		fmt.Println(string(marshaledElements))
	} else {
		return err
	}
	return nil
}

func modifiedFilePaths(repoRoot string) []string {
	repo, err := git.PlainOpen(repoRoot)
	if err != nil {
		log.Fatal(err)
	}
	var other *object.Commit
	if defaultBranchIsCheckedout(repo) {
		other = previousCommit(repo)
	} else {
		other = defaultBranchHeadCommit(repo)
	}
	patch, err := headCommit(repo).Patch(other)
	if err != nil {
		log.Fatal(err)
	}
	filePaths := []string{}
	for _, filePatch := range patch.FilePatches() {
		from, to := filePatch.Files()
		if from != nil {
			filePaths = append(filePaths, from.Path())
		}
		if to != nil {
			filePaths = append(filePaths, to.Path())
		}
	}
	return filePaths
}

func defaultBranchIsCheckedout(repo *git.Repository) bool {
	head := mustHead(repo)
	return head.Name().IsBranch() && head.Name().Short() == defaultBranchName
}

func defaultBranchHeadCommit(repo *git.Repository) *object.Commit {
	ref, err := repo.Reference("refs/remotes/origin/"+defaultBranchName, true)
	if err != nil {
		log.Fatal(err)
	}
	commitHistory, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		log.Fatal(err)
	}
	lastCommit, err := commitHistory.Next()
	if err != nil {
		log.Fatal(err)
	}
	return lastCommit
}

func headCommit(repo *git.Repository) *object.Commit {
	head := mustHead(repo)
	commitHistory, err := repo.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		log.Fatal(err)
	}
	lastCommit, err := commitHistory.Next()
	if err != nil {
		log.Fatal(err)
	}
	return lastCommit
}

func previousCommit(repo *git.Repository) *object.Commit {
	head := mustHead(repo)
	commitHistory, err := repo.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		log.Fatal(err)
	}
	_, _ = commitHistory.Next()
	commit, err := commitHistory.Next()
	if err != nil {
		log.Fatal(err)
	}
	return commit
}

func mustHead(repo *git.Repository) plumbing.Reference {
	head, err := repo.Head()
	if err != nil {
		log.Fatal(err)
	}
	return *head
}

// Find the relative level 2 directory from a path e,g,
// containers/hello-world/version.txt -> containers/hello-world
// if it does not exist then return nil
func l2DirectoryFromFile(path string) *string {
	re := regexp.MustCompile(packageRegex)
	matches := re.FindStringSubmatch(path)
	if len(matches) < 2 {
		return nil
	} else {
		return &(matches[1])
	}
}

func values[Map ~map[K]V, K comparable, V any](m Map) []V {
	vs := []V{}
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}
