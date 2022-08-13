package api_git

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

/*
*Funcion que obtiene un commit por medio de su Hash
 */
func getCommit(hash string, repo git.Repository) (*object.Commit, error) {
	hashObject := plumbing.NewHash(hash)
	commit, err := repo.CommitObject(hashObject)
	if err != nil {
		return nil, err
	}

	return commit, nil
}

func diffToHead(path string, hash string) (*object.Patch, error) {
	repo, err := git.PlainOpen(path)
	var patch *object.Patch = nil

	if repo != nil {
		commit, err := getCommit(hash, *repo)
		if err == nil && commit != nil {

			ref, _ := repo.Head()
			commitHead, _ := repo.CommitObject(ref.Hash())
			//Compare with the head of the repo
			patch, err = commit.Patch(commitHead)
			if err == nil && patch != nil {
				return patch, nil
			}
		}
	}

	if err != nil {
		return nil, err
	}
	return patch, err

}

/*
*Funcion que compara los dos repositorios (URL_main_repo/URL_Diff_repo)
comprueba si existen cambios en los repositorios y retorna cuales
son los archivos que se modificaron
*/
func diffTreeRepos(diffRequest DiffRequest) (object.Changes, error) {
	if diffRequest.UrlMain != "" && diffRequest.UrlDiff != "" {

		mainRepo, err := git.PlainOpen(diffRequest.UrlMain)
		diffRepo, err := git.PlainOpen(diffRequest.UrlDiff)

		var refMainHead *plumbing.Reference
		var refDiffHead *plumbing.Reference
		var commitMain *object.Commit
		var commitDiff *object.Commit
		var objectChanges object.Changes

		if diffRequest.HashMain != "" {
			commitMain, err = getCommit(diffRequest.HashMain, *mainRepo)
		} else {
			refMainHead, err = mainRepo.Head()
			commitMain, err = mainRepo.CommitObject(refMainHead.Hash())
		}

		if diffRequest.HashDiff != "" {
			commitDiff, err = getCommit(diffRequest.HashDiff, *diffRepo)
		} else {
			refDiffHead, err = diffRepo.Head()
			commitDiff, err = diffRepo.CommitObject(refDiffHead.Hash())
		}

		if err == nil {
			treeMain, err := commitMain.Tree()
			treeDiff, err := commitDiff.Tree()
			if err == nil {
				objectChanges, err = object.DiffTree(treeMain, treeDiff)
			}
		}

		if err != nil {
			fmt.Print("Error" + err.Error())
			return nil, err
		}
		return objectChanges, nil
	} else {
		return nil, errors.New("No contiene la direccion de algunos de los repositorios")
	}

}

func diffTreeBranches(diffRequest DiffRequestBranches) (object.Changes, error) {
	/* if diffRequest.BaseBranchName != "" && diffRequest.RefBranchName != "" { */

	repo, err := OpenProjectRepository(diffRequest.UrlRepoMain, diffRequest.ProjecName)

	if err != nil {
		fmt.Print("Error" + err.Error())
		return nil, err
	}

	var commitMain *object.Commit
	var commitDiff *object.Commit
	var refMainHead *plumbing.Reference
	var refDiffHead *plumbing.Reference
	var objectChanges object.Changes

	if diffRequest.BaseBranchName == "" {
		refMainHead, err = repo.Head()
	} else {
		refMainHead, err = repo.Reference(plumbing.ReferenceName(BranchPrefix+diffRequest.BaseBranchName), true)
	}

	if err != nil {
		fmt.Print("Error" + err.Error())
		return nil, err
	}

	refDiffHead, err = repo.Reference(plumbing.ReferenceName(BranchPrefix+diffRequest.RefBranchName), true)

	if err != nil {
		fmt.Print("Error" + err.Error())
		return nil, err
	}

	commitMain, err = repo.CommitObject(refMainHead.Hash())

	if err != nil {
		fmt.Print("Error" + err.Error())
		return nil, err
	}

	commitDiff, err = repo.CommitObject(refDiffHead.Hash())

	if err == nil {

		treeMain, err := commitMain.Tree()

		if err != nil {
			fmt.Print("Error" + err.Error())
			return nil, err
		}

		treeDiff, err := commitDiff.Tree()

		if err != nil {
			fmt.Print("Error" + err.Error())
			return nil, err
		}
		if err == nil {

			objectChanges, err = object.DiffTree(treeMain, treeDiff)

			if err != nil {
				fmt.Print("Error" + err.Error())
				return nil, err
			}
		}
	}

	if err != nil {
		fmt.Print("Error" + err.Error())
		return nil, err
	}
	return objectChanges, nil
	/* } else {
		return nil, errors.New("No contiene la direccion de algunos de los repositorios")
	} */

}
