package repo

import (
	"log"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	defaultBranchName = "main"
)

type Patch struct {
	Old string
	New string
}

func FilePatch(repoRoot string, filePath string) *Patch {
	for _, filePatch := range patchToOtherCommit(repoRoot).FilePatches() {
		_, to := filePatch.Files()
		if to == nil || to.Path() != filePath {
			continue
		}
		patch := Patch{}
		for _, chunk := range filePatch.Chunks() {
			if chunk.Type() == diff.Add {
				patch.New += chunk.Content()
			} else if chunk.Type() == diff.Delete {
				patch.Old += chunk.Content()
			}
		}
		return &patch
	}
	return nil
}

// Get all files that have been modifed between the current commit
// and either origin/main if we're on a branch or the previous main
// commit if we're on main
func ModifiedFilePaths(repoRoot string) []string {
	filePaths := []string{}
	for _, filePatch := range patchToOtherCommit(repoRoot).FilePatches() {
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

func patchToOtherCommit(repoRoot string) *object.Patch {
	repo := must(git.PlainOpen(repoRoot))
	return must(diffCommit(repo).Patch(headCommit(repo)))
}

func diffCommit(repo *git.Repository) *object.Commit {
	if defaultBranchIsCheckedout(repo) {
		return previousCommit(repo)
	} else {
		return defaultBranchHeadCommit(repo)
	}
}

func defaultBranchIsCheckedout(repo *git.Repository) bool {
	head := must(repo.Head())
	return head.Name().IsBranch() && head.Name().Short() == defaultBranchName
}

func defaultBranchHeadCommit(repo *git.Repository) *object.Commit {
	ref := must(repo.Reference("refs/remotes/origin/"+defaultBranchName, true))
	commitHistory := must(repo.Log(&git.LogOptions{From: ref.Hash()}))
	lastCommit := must(commitHistory.Next())
	return lastCommit
}

func headCommit(repo *git.Repository) *object.Commit {
	head := must(repo.Head())
	commitHistory := must(repo.Log(&git.LogOptions{From: head.Hash()}))
	return must(commitHistory.Next())
}

func previousCommit(repo *git.Repository) *object.Commit {
	head := must(repo.Head())
	commitHistory := must(repo.Log(&git.LogOptions{From: head.Hash()}))
	_, _ = commitHistory.Next()
	return must(commitHistory.Next())
}

func must[T any](value T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return value
}
