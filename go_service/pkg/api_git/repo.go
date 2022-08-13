package api_git

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

//Create repo from bash
func CreateRepoFromBash(name string) (bool, error) {
	out, err := exec.Command("/bin/sh", "../../../etc/git-create-repo.sh", name).Output()
	if err != nil {
		return false, err
	}
	if out != nil {
		return true, nil
	}
	return false, err
}

/* func CopyRepoFromTag(tagRequest TagRequest) (bool, error) {
	isCreate, err := CreateRepoFromBash(tagRequest.Name)
	if err != nil {
		return false, err
	}
	if isCreate {
		path := "/var/www/git/" + tagRequest.Name + ".git"
		repo, err := CloneRepository(path, tagRequest.Url)
		if err != nil {
			if err.Error() == "repository already exists" {
				repo, err = git.PlainOpen(path)
				if err != nil {
					return false, err
				}
			} else {
				return false, err
			}

		}
		if repo != nil {
			tagReference, err := TagExists(tagRequest.Tag, *repo)
			if err != nil {
				return false, err
			}
			if tagReference.Name().String() != "" {
				result, err := createBranchFromTag(&tagReference, repo)
				if err != nil {
					return false, err
				}
				if result {
					return true, nil
				}
			} else {
				return false, errors.New("Tag does not exist")
			}
		}
	} else {
		return false, err
	}
	return false, err
} */

//Function to open repository
func OpenRepository(RepoPath string) (*git.Repository, error) {

	if RepoPath != "" {
		ConcatRepoPath := /*baseRepoDir +*/ RepoPath + "/" + ".git"

		repository, err := git.PlainOpen(ConcatRepoPath)

		if err != nil {
			return nil, git.ErrRepositoryNotExists
		}
		return repository, err
	}
	return nil, git.ErrRepositoryNotExists
}

func OpenProjectRepository(RepoName string, ProjectName string) (*git.Repository, error) {

	if RepoName != "" {
		ConcatRepoPath := "/go_service/repos/" + ProjectName + "/" + RepoName

		repository, err := git.PlainOpen(ConcatRepoPath)

		if err != nil {
			return nil, git.ErrRepositoryNotExists
		}
		return repository, err
	}
	return nil, git.ErrRepositoryNotExists
}

func createRepository(NameRpository string, ProyectName string, Readmefile bool, ignoreFile bool, isPublic bool, UserName string, userEmail string) ([]string, error) {

	path := `/go_service/repos/` + ProyectName + `/` + NameRpository

	var branchesName []string
	r, err := git.PlainInit(path, false)

	if err != nil {
		return nil, err
	}

	//Get workTree

	worktree, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	//CREATE A FILE (.gitignore or README) INTO HEAD BRANCH

	if Readmefile && ignoreFile {

		createReadmefile := filepath.Join(path, "README.md")

		err = ioutil.WriteFile(createReadmefile, []byte("#"+NameRpository), 0644)

		if err != nil {
			return nil, err
		}

		_, err := worktree.Add(`README.md`)

		if err != nil {
			return nil, err
		}

		createGitIgnoreFile := filepath.Join(path, ".gitignore")

		err = ioutil.WriteFile(createGitIgnoreFile, []byte(" "), 0644)

		if err != nil {
			return nil, err
		}

		BlobFileAdd, err := worktree.Add(`.gitignore`)

		if err != nil {
			return nil, err
		}

		status, err := worktree.Status()
		if err != nil {
			return nil, err
		}

		fmt.Println(status)

		CommitHash, err := worktree.Commit("initial Commit", &git.CommitOptions{
			Author: &object.Signature{
				Name:  UserName,
				Email: userEmail,
				When:  time.Now(),
			},
		})

		if err != nil {
			return nil, err
		}

		branches, err := r.Branches()

		_ = branches.ForEach(func(branch *plumbing.Reference) error {
			branchesName = append(branchesName, strings.TrimPrefix(branch.Name().String(), BranchPrefix))
			return nil
		})

		if err != nil {
			return nil, err
		}

		fmt.Println(CommitHash)

		fmt.Println(BlobFileAdd)
		
		return branchesName, err

	}

	if Readmefile && !ignoreFile {
		createReadmefile := filepath.Join(path, "README.md")

		err = ioutil.WriteFile(createReadmefile, []byte(" # "+NameRpository), 0644)

		if err != nil {
			return nil, err
		}

		_, err := worktree.Add(`README.md`)

		if err != nil {
			return nil, err
		}

		status, err := worktree.Status()
		if err != nil {
			return nil, err
		}

		fmt.Println(status)

		CommitHash, err := worktree.Commit("initial Commit", &git.CommitOptions{
			Author: &object.Signature{
				Name:  UserName,
				Email: userEmail,
				When:  time.Now(),
			},
		})

		if err != nil {
			return nil, err
		}

		branches, err := r.Branches()

		_ = branches.ForEach(func(branch *plumbing.Reference) error {
			branchesName = append(branchesName, strings.TrimPrefix(branch.Name().String(), BranchPrefix))
			return nil
		})

		if err != nil {
			return nil, err
		}

		fmt.Println(CommitHash)
	}

	if ignoreFile && !Readmefile {
		createGitIgnoreFile := filepath.Join(path, ".gitignore")

		err = ioutil.WriteFile(createGitIgnoreFile, []byte(" "), 0644)

		if err != nil {
			return nil, err
		}

		_, err := worktree.Add(`.gitignore`)

		if err != nil {
			return nil, err
		}

		status, err := worktree.Status()
		if err != nil {
			return nil, err
		}

		fmt.Println(status)

		CommitHash, err := worktree.Commit("initial Commit", &git.CommitOptions{
			Author: &object.Signature{
				Name:  UserName,
				Email: userEmail,
				When:  time.Now(),
			},
		})

		if err != nil {
			return nil, err
		}

		branches, err := r.Branches()

		_ = branches.ForEach(func(branch *plumbing.Reference) error {
			branchesName = append(branchesName, strings.TrimPrefix(branch.Name().String(), BranchPrefix))
			return nil
		})

		if err != nil {
			return nil, err
		}

		fmt.Println(CommitHash)

	}
	return branchesName, err
}
