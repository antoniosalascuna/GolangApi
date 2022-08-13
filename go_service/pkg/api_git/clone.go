package api_git

import (
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

//Clone a repository in a specific directory
func CloneRepository(repoName string, ProjectName string, url string) (*git.Repository, error) {

	if repoName != "" || url != "" {
		pathClone := `/go_service/repos/` + ProjectName + `/` + repoName

		repository, err := git.PlainClone(pathClone, false, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})

		if err != nil {
			return nil, err
		}
		return repository, nil
	}

	return nil, git.ErrRepositoryNotExists
}

func CloneRepositoryWithUserNamePassword(repoName string, ProjectName string, url string, username string, password string) (bool, error) {

	if repoName != "" || url != "" {
		pathClone := `/go_service/repos/` + ProjectName + `/` + repoName

		_, err := git.PlainClone(pathClone, false, &git.CloneOptions{

			Auth:     &http.BasicAuth{Username: username, Password: password},
			URL:      url,
			Progress: os.Stdout,
		})

		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, git.ErrRepositoryNotExists
}
