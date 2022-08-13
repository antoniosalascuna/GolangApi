package api_git

import (
	"bytes"
	"fmt"
	"go_service/tools"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

//Variable to point to Base Repository
const baseRepoDir = "/go_service/repos/"
const BranchPrefix = "refs/heads/"

//CreateMerge Method to create merge between to branches
func CreateMerge(mergeRequest MergeRequest) (string, error) {

	repo, err := OpenProjectRepository(mergeRequest.RepoName, mergeRequest.ProjecName)

	var result string = ""
	if repo != nil {
		reference, _ := getBranch(*repo, mergeRequest.TargetBranch)
		if reference != nil {
			isReference, err := checkoutBranch(*reference, *repo)

			if isReference {

				cmd := exec.Command("git", "merge", mergeRequest.Branch)
				cmd.Dir = baseRepoDir + mergeRequest.ProjecName + "/" + mergeRequest.RepoName

				var stderr bytes.Buffer

				cmd.Stderr = &stderr
				resultCmd, err := cmd.Output()

				if err != nil {
					fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
					return "", err
				}

				result = tools.BytesToString(resultCmd)
			} else if err != nil {
				return result, err
			}
		}

		return result, nil
	} else {
		return "", err
	}

}

//Checkout a branch
func checkoutBranch(reference plumbing.Reference, repo git.Repository) (bool, error) {
	w, err := repo.Worktree()
	status, err := w.Status()
	fmt.Println(status)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: reference.Name(),
		Force:  true,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

//Get specific branch reference
func getBranch(repo git.Repository, branchName string) (*plumbing.Reference, error) {

	reference, err := repo.Reference(plumbing.ReferenceName(BranchPrefix+branchName), true)
	if err != nil {
		return nil, err
	} else {
		return reference, nil
	}
}
