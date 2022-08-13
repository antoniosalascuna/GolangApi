package api_git

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type FileDetails struct {
	Name    string `json:"Name,omitempty"`
	Lines   int    `json:"Lines,omitempty"`
	Bytes   int    `json:"Bytes,omitempty"`
	Content string `json:"Content,omitempty"`
}

//Obtengo los atributos de un Blob Object
func BlobDetails(RepositoryName string, FilePath string) (*FileDetails, error) {

	RepoOpen, err := OpenRepository(RepositoryName)
	FilesAtt := FileDetails{}
	if err != nil {
		return nil, err
	}
	//Esta parte se puede cambiar por un commiten especifico
	TreeHead := TreeCommitHead(RepoOpen)

	entry, err := TreeHead.File(FilePath)
	if err != nil {
		return nil, err
	}

	nameFile, err := TreeHead.FindEntry(FilePath)

	if err != nil {
		return nil, err
	}
	content, err := entry.Lines()
	if err != nil {
		return nil, err
	}
	FilesAtt.Lines = len(content)
	FilesAtt.Bytes = int(entry.Size)
	FilesAtt.Name = nameFile.Name

	return &FilesAtt, nil
}

func GetMediContent(RepositoryName string, ProjectName string, FilePath string) (string, error) {
	RepoOpen, err := OpenProjectRepository(RepositoryName, ProjectName)
	if err != nil {
		return "", err
	}

	TreeHead := TreeCommitHead(RepoOpen)

	entry, err := TreeHead.File(FilePath)
	if err != nil {
		return "", err
	}

	content, err := entry.Contents()

	if err != nil {
		return "", err
	}

	return content, err

}

func BlobReadmeDetails(RepositoryName string, ProjectName string, FilePath string) (*FileDetails, error) {

	RepoOpen, err := OpenProjectRepository(RepositoryName, ProjectName)
	FilesAtt := FileDetails{}
	if err != nil {
		return nil, err
	}
	//Esta parte se puede cambiar por un commiten especifico
	TreeHead := TreeCommitHead(RepoOpen)

	entry, err := TreeHead.File(FilePath)
	if err != nil {
		return nil, err
	}

	nameFile, err := TreeHead.FindEntry(FilePath)

	if err != nil {
		return nil, err
	}
	Lines, err := entry.Lines()
	if err != nil {
		return nil, err
	}

	content, err := entry.Contents()

	if err != nil {
		return nil, err
	}

	FilesAtt.Lines = len(Lines)
	FilesAtt.Bytes = int(entry.Size)
	FilesAtt.Name = nameFile.Name
	FilesAtt.Content = content

	return &FilesAtt, nil
}

func CreateInitialRepoFile(file FileCreate) (bool, error) {

	path := `/go_service/repos/` + file.ProjectName + `/` + file.RepositoryName

	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, err
	}

	createReadmefile := filepath.Join(path, file.NameFile)

	err = ioutil.WriteFile(createReadmefile, []byte(file.ContentFile), 0644)

	if err != nil {
		return false, err
	}

	worktree, err := repo.Worktree()

	if err != nil {
		return false, err
	}

	hash, err := worktree.Add(file.NameFile)

	if err != nil {
		return false, err
	}

	CommitHash, err := worktree.Commit(file.CommitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  file.CommiterName,
			Email: file.CommiterEmail,
			When:  time.Now(),
		},
	})

	if err != nil {
		return false, err
	}

	fmt.Println(hash)
	fmt.Println(CommitHash)

	return true, err
}
