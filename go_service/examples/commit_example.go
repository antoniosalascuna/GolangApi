package examples

// change package to main and put it inside the main.go file if u need to see
// the running example

import (
	"fmt"

	"log"

	//	"strings"

	//"github.com/go-git/go-billy/v5/osfs"

	"github.com/go-git/go-git/v5"

	"github.com/go-git/go-git/v5/plumbing/object"
)

func main() {

	repo, _ := git.PlainOpen("../monkeytest.git")

	ref, _ := repo.Head()

	commit, _ := repo.CommitObject(ref.Hash())

	var author = commit.Author

	var committer = commit.Committer

	var hash = commit.Hash

	var message = commit.Message

	var time = commit.Author.When

	fmt.Println("Author of the commit:", author)
	fmt.Println("Commiter of this commit: ", committer)
	fmt.Println("Message of this commit: ", message)
	fmt.Println("hash of this commit: ", hash)
	fmt.Println("Time stamp of the commit: ", time)

}

func ListFile(Path string, Namefile string) (string, error) {
	repo, err := git.PlainOpen(Path)
	var file string
	if repo != nil {
		treeIter, errW := repo.TreeObjects()

		if treeIter != nil {
			treeIter.ForEach(func(t *object.Tree) error {

				// ... get the files iterator and print the file

				t.Files().ForEach(func(f *object.File) error {

					if f.Name == Namefile {

						fmt.Printf(f.Name, "entra aqui")

						file = f.Name
					} else {
						log.Printf("There is not exist a file ")
					}
					return nil
				})
				return nil
			})
		}
		if errW != nil {
			log.Printf("There is not a tree %s", errW)
		}
	}
	return file, err
}

/*Funcion que retorna el contenido de cada Blob File
con la ruta dada por parametro*/
func ListContenBlobFile(repoPath string, fileP string) ([]string, error) {

	repo, err := git.PlainOpen(repoPath)

	var filepath []string

	var datafile []string

	ref, err := repo.Head()

	commit, err := repo.CommitObject(ref.Hash())

	tree, err := commit.Tree()

	/*For que llena el array para que
	Se pueda obtener el main tree de repositorio
	*/
	tree.Files().ForEach(func(f *object.File) error {

		filepath = append(filepath, f.Name)

		fmt.Printf("File Hashe and Path: %s    %s\n", f.Hash, f.Name)

		return nil
	})

	//Retorna el file de la ruta que se le pasa
	//Ejemplo
	/*
		repoPath = /var/www/git/repoexample.git
		fileP = css/assets

		retorna -> ERROR porque no es un archivo esta path es de una carpeta

		fileP = css/styles.css

		retorna ->  "body{\n    background-color: aqua;\n}\n\nheader{\n    position: relative;\n}"
	*/

	treefile, err := tree.File(fileP)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(filepath); i++ {

		if filepath[i] == fileP {

			content, err := treefile.Contents()

			datafile = append(datafile, string(content))

			if err != nil {

			}
			fmt.Printf("FilePath and parameter: %s    %s\n", filepath[i], fileP)
			fmt.Println(content)
			break
		}

	}

	fmt.Println(treefile)

	return datafile, err
}
