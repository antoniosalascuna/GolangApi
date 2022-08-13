package api_git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type BranchInfo struct {
	Count int      `json:"Count"`
	Names []string `json:"BranchesName"`
}

const baseRepoStoreDir = "/go_service/repos/"

/*Get all branches of Repository*/
func GetBranches(RepoPath string, ProjectName string) ([]string, string, error) {

	path := baseRepoStoreDir + ProjectName + "/" + RepoPath

	RepoOpen, err := OpenRepository(path)

	if err != nil {
		return nil, "", err
	}

	//BrInfo := BranchInfo{}

	var branchesNames []string
	var branchesWithoutMaster []string

	branches, err := RepoOpen.Branches()
	if err != nil {
		return nil, "", err
	}

	defaultbranch, err := RepoOpen.Head()
	if err != nil {
		return nil, "", err
	}

	headName := strings.TrimPrefix(defaultbranch.Name().String(), BranchPrefix)

	if err != nil {
		return nil, "", err
	}

	_ = branches.ForEach(func(branch *plumbing.Reference) error {

		branchesNames = append(branchesNames, strings.TrimPrefix(branch.Name().String(), BranchPrefix))

		return nil
	})

	for i := range branchesNames {
		if branchesNames[i] == headName {
			continue
		} else {
			branchesWithoutMaster = append(branchesWithoutMaster, branchesNames[i])
		}
	}

	return branchesWithoutMaster, headName, nil
}

/*Get a Branch By Name*/
func GetBranch(RepoPath string, branchName string, ProjectName string) {

	path := baseRepoStoreDir + RepoPath

	repo, err := OpenRepository(path)

	if err != nil {
		//return err
		fmt.Println(err)
	}

	//Ref al Head
	repohead, err := repo.Head()

	//Get All references from repo
	references, err := repo.References()

	//Get all branches to reposity
	braches, err := repo.Branches()

	//Get a branch
	branch, err := repo.Branch(branchName)

	//Muestra el historial de commits de branch

	if err != nil {
		//return err
		fmt.Println(err)
	}

	fmt.Println(branch)
	fmt.Println(braches)
	fmt.Println(references)
	fmt.Println(repohead)
	//return nil
}

func createBranch(NewbranchName string, ProjecName string, repoPath string, branchName string) error {

	if repoPath != "" {

		//	path := baseRepoStoreDir + repoPath

		//repo, err := OpenRepository(path)

		repo, err := OpenProjectRepository(repoPath, ProjecName)

		if err != nil {
			return err
		}

		//Apunto la referencia al HEAD para crear una rama apartir del HEAD
		if branchName != "" {

		}
		headRef, err := repo.Head()
		if err != nil {
			return err
		}

		//newBranchReference := plumbing.NewBranchReferenceName(branchName)
		newBranch := plumbing.ReferenceName("refs/heads/" + branchName)

		ref := plumbing.NewHashReference(newBranch, headRef.Hash())

		//Asigno la nueva ref del branch con el head
		err = repo.Storer.SetReference(ref)

		if err != nil {
			return err
		}

		fmt.Println(ref)
	}

	return nil
}

func InitialCommitChange(RepoName string, ProjectName string) error {

	repo, err := OpenProjectRepository(RepoName, ProjectName)

	if err != nil {
		return err
	}

	var changes []*object.Change

	headRepo, err := repo.Head()

	commiHead, err := repo.CommitObject(headRepo.Hash())

	tree, err := commiHead.Tree()

	tree.Files().ForEach(func(f *object.File) error {

		e, _ := tree.FindEntry(f.Name)

		ch := &object.Change{
			From: object.ChangeEntry{},
			To: object.ChangeEntry{
				Name:      f.Name,
				Tree:      tree,
				TreeEntry: *e,
			},
		}

		changes = append(changes, ch)
		return nil
	})

	fmt.Println(changes)

	for _, value := range changes {

		fmt.Println(value.Patch())
	}

	//change, err := ch.Action()

	//from, to, err := ch.Files()
	/* references, err := repo.References()

	references.ForEach(func(r *plumbing.Reference) error {

		fmt.Println(r.Strings())

		return nil
	}) */

	//fmt.Println(references)

	//repoHead, err := repo.Head()

	/* 	var commits []*object.Commit

	   	if err != nil {
	   		return err
	   	}

	   	commitIter, err := repo.CommitObjects()

	   	commitIter.ForEach(func(c *object.Commit) error {

	   		commits = append(commits, c)

	   		return nil
	   	})

	   	var prev *object.Commit

	   	for _, commit := range commits {

	   		if prev == nil {
	   			prev = commit
	   			continue
	   		}

	   		tree, err := commit.Tree()
	   		if err != nil {
	   			return err
	   		}

	   		prevTree, err := prev.Tree()
	   		if err != nil {
	   			return err
	   		}

	   		changes, err := object.DiffTree(tree, prevTree)

	   		for _, change := range changes {
	   			_, _, err := change.Files()

	   			if err != nil {
	   				return err
	   			}
	   		}

	   		prev = commit
	   	} */

	return err

}

func getBranches(repo *git.Repository, skip, limit int) ([]string, int, error) {

	var branchNames []string

	branches, err := repo.Branches()

	if err != nil {
		return nil, 0, err
	}

	i := 0

	count := 0

	_ = branches.ForEach(func(r *plumbing.Reference) error {
		count++

		if i < skip {
			i++
			return nil
		} else if limit != 0 && count > skip+limit {
			return nil
		}

		branchNames = append(branchNames, strings.TrimPrefix(r.Name().String(), BranchPrefix))
		return nil
	})

	return branchNames, count, nil
}
