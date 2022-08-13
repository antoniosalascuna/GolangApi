package api_git

import (
	"fmt"
	"io"
	"path"

	"github.com/emirpasic/gods/trees/binaryheap"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	commitgraph_fmt "github.com/go-git/go-git/v5/plumbing/format/commitgraph"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/object/commitgraph"
)

type commitAndPaths struct {
	commit commitgraph.CommitNode
	// Paths that are still on the branch represented by commit
	paths []string
	// Set of hashes for the paths
	hashes map[string]plumbing.Hash
}

type CommitInfo struct {
	CommitID       plumbing.Hash `json:"CommitID,omitempty"`
	Committer      string        `json:"Committer,omitempty"`
	Message        string        `json:"Message,omitempty"`
	Autor          string        `json:"Autor,omitempty"`
	CommitTime     string        `json:"Date,omitempty"`
	CommitString   string        `json:"CommitString,omitempty"`
	BaseBranchName string        `json:"BaseBranchName,omitempty"`
	RefBranchName  string        `json:"RefBranchName,omitempty"`
}
type Commit struct {
	CountCommits int          `json:"CountCommits"`
	CommitInfo   []CommitInfo `json:"CommitHistory,omitempty"`
}

type ActionPaths struct {
	Action       string
	Path         string
	PatchContent string
}
type ChangeDescription struct {
	ActionsAndPaths []ActionPaths `json:"Paths,omitempty"`
	CommitContent   string        `json:"CommitContent,omitempty"`
	CommitMessage   string        `json:"CommitMessage,omitempty"`
}

func TreeCommitHead(repo *git.Repository) (repotree *object.Tree) {

	headrepo, _ := repo.Head()

	commit, _ := repo.CommitObject(headrepo.Hash())

	tree, _ := commit.Tree()

	return tree

}

//Get all SHA commits
func GetAllSHACommits(RepoPath string) {

	var CommitHash []string
	RepoOpen, err := OpenRepository(RepoPath)

	if err != nil {
		return
	}

	repoobjects, err := RepoOpen.CommitObjects()
	if err != nil {
		return
	}

	if err != nil {
		return
	}

	repoobjects.ForEach(func(c *object.Commit) error {

		CommitHash = append(CommitHash, c.Hash.String())

		return nil
	})

}

func GetCommitDiff(RepoName string, BaseBranchNamePR string, RefBranchNamePR string, ProjecName string) ([]CommitInfo, error) {
	var result []plumbing.Hash

	//var resultRefCommits []plumbing.Hash

	var commitReF []plumbing.Hash

	commitmap := []CommitInfo{}

	CommitHistoryRef, err := getBranchCommitHistory(RepoName, ProjecName, RefBranchNamePR)

	if err != nil {
		return commitmap, err
	}
	CommitHistoryBase, err := getBranchCommitHistory(RepoName, ProjecName, BaseBranchNamePR)

	if err != nil {
		return commitmap, err
	}

	for i := 0; i < len(CommitHistoryBase.CommitInfo); i++ {

		fmt.Println(CommitHistoryBase.CommitInfo[i].CommitID)

		for j := 0; j < len(CommitHistoryRef.CommitInfo); j++ {

			fmt.Println(CommitHistoryRef.CommitInfo[j].CommitID)

			/*Obtengo los diferentes commits iD entre ramas*/
			if CommitHistoryBase.CommitInfo[i].CommitID == CommitHistoryRef.CommitInfo[j].CommitID {

				result = append(result, CommitHistoryRef.CommitInfo[j].CommitID)
			}
		}
	}

	for _, refdata := range CommitHistoryRef.CommitInfo {

		result = append(result, refdata.CommitID)
	}

	for i := 0; i < len(result); i++ {
		if isUnique(result[i], result) {
			commitReF = append(commitReF, result[i])
		}
	}

	repo, err := OpenProjectRepository(RepoName, ProjecName)

	for _, CommitId := range commitReF {

		commitObject, err := repo.CommitObject(CommitId)

		if err != nil {
			return commitmap, err
		}

		commitmap = append(commitmap, CommitInfo{
			CommitID:       commitObject.ID(),
			Committer:      commitObject.Committer.Name,
			Message:        commitObject.Message,
			Autor:          commitObject.Author.Name,
			CommitTime:     commitObject.Committer.When.Local().String(),
			CommitString:   commitObject.ID().String(),
			RefBranchName:  RefBranchNamePR,
			BaseBranchName: BaseBranchNamePR,
		})

	}

	return commitmap, err

}

// Funcion que toma los valores unicos de un array
func isUnique(value plumbing.Hash, array []plumbing.Hash) bool {
	var contador int = 0
	for i := 0; i < len(array); i++ {
		if array[i] == value {
			contador++
		}
	}
	return contador == 1
}

//Get commit history by defaoult branch in this case Head()
func GetCommitHistory(RepoPath string, ProjectName string) (*Commit, error) {

	repo, err := OpenProjectRepository(RepoPath, ProjectName)
	if err != nil {
		return nil, err
	}
	// We instantiate a new repository targeting the given path (the .git folder)
	fs := osfs.New("/go_services/repos/" + ProjectName + "/" + RepoPath + ".git")
	if _, err := fs.Stat(git.GitDirName); err == nil {
		fs, err = fs.Chroot(git.GitDirName)
		fmt.Println(err)
	}

	CommitNodeIndes, file := getCommitNodeIndex(repo, fs)

	if file != nil {
		defer file.Close()
	}

	HeadReference, err := repo.Head()

	if err != nil {
		return nil, err
	}
	//Esta linea de codigo hace lo mismo que el comando de Git Log
	CommitBranchHistory, err := repo.Log(&git.LogOptions{From: HeadReference.Hash()})

	if err != nil {
		return nil, err
	}
	commitmap := []CommitInfo{}
	var count int
	//Recorro cada commit que existe en una branch para obtener la informacion de cada uno
	//lo almaceno en CommitMap como tipo array para devolver la respuesto tipo [{array de objectos}]
	CommitBranchHistory.ForEach(func(c *object.Commit) error {
		count++
		//commitID = append(commitID, c.Hash)
		commitNode, err := CommitNodeIndes.Get(c.Hash)
		if err != nil {
			return err
		}

		commitmap = append(commitmap, CommitInfo{
			CommitID:     c.ID(),
			Committer:    c.Committer.Name,
			Message:      c.Message,
			Autor:        c.Author.Name,
			CommitTime:   commitNode.CommitTime().Local().String(),
			CommitString: c.ID().String(),
		})
		/*
			Date Format

			2021-11-01 16:22:09 +0000 UTC ---->> commitNode.CommitTime().Local().String()
			2021-11-01 10:22:09 -0600 -0600 --- >> commitNode.CommitTime().String()
			time.Date(2021, time.November, 1, 10, 22, 9, 0, time.Location(\"\")) -->> commitNode.CommitTime().GoString()
		*/
		return nil
	})

	//Esta variable lo que hace es que dentro de un struct Commit
	// Solo para agregarle la cantidad de commits que posee un branch
	commitFull := Commit{count, commitmap}

	return &commitFull, nil

}

func getBranchCommitHistory(RepoName string, ProjectName string, branchName string) (*Commit, error) {

	repo, err := OpenProjectRepository(RepoName, ProjectName)
	var branchReference *plumbing.Reference
	if err != nil {
		return nil, err
	}

	branchReference, err = repo.Reference(plumbing.ReferenceName(BranchPrefix+branchName), true)
	if err != nil {
		return nil, err
	}

	CommitBranchHistory, err := repo.Log(&git.LogOptions{From: branchReference.Hash()})

	if err != nil {
		return nil, err
	}
	commitmap := []CommitInfo{}
	var count int

	CommitBranchHistory.ForEach(func(c *object.Commit) error {
		count++
		//commitID = append(commitID, c.Hash)
		//commitNode, err := CommitNodeIndes.Get(c.Hash)
		if err != nil {
			return err
		}

		commitmap = append(commitmap, CommitInfo{
			CommitID:     c.ID(),
			Committer:    c.Committer.Name,
			Message:      c.Message,
			Autor:        c.Author.Name,
			CommitTime:   c.Author.When.String(),
			CommitString: c.ID().String(),
		})
		/*
			Date Format

			2021-11-01 16:22:09 +0000 UTC ---->> commitNode.CommitTime().Local().String()
			2021-11-01 10:22:09 -0600 -0600 --- >> commitNode.CommitTime().String()
			time.Date(2021, time.November, 1, 10, 22, 9, 0, time.Location(\"\")) -->> commitNode.CommitTime().GoString()
		*/
		return nil
	})

	commitFull := Commit{count, commitmap}

	return &commitFull, nil

}

func getTreeCommit(RepoName string, Sha1Id string, ProjectName string) (*ChangeDescription, error) {

	repo, err := OpenProjectRepository(RepoName, ProjectName)
	if err != nil {
		return nil, err
	}

	//ChangeD := ChangeDescription{}

	ActPaths := []ActionPaths{}

	/* 	ref, err := repo.Head()

	   	commitHead, err := repo.CommitObject(ref.Hash())


	   	fmt.Println(commitHead)*/

	h, err := repo.ResolveRevision(plumbing.Revision(Sha1Id))

	if err != nil {
		return nil, err
	}

	commitDiff, err := repo.CommitObject(*h)

	if err != nil {
		return nil, err
	}

	commitIter, err := repo.Log(&git.LogOptions{From: commitDiff.Hash})

	if err != nil {
		return nil, err
	}
	defer commitIter.Close()

	//var arraychuncks []string

	//Se declara una variable que va a contener prevCommit y el prevTree
	var prevCommit *object.Commit
	var prevTree *object.Tree

	//Recorre todos los commit para obtener los cambios del commit que viene por parametro
	for {
		commit, err := commitIter.Next()

		if err == io.EOF {

			var changes []*object.Change

			commitTree, err := prevCommit.Tree()
			if err != nil {
				return nil, err
			}

			commitTree.Files().ForEach(func(f *object.File) error {

				e, _ := commitTree.FindEntry(f.Name)

				ch := &object.Change{
					From: object.ChangeEntry{},
					To: object.ChangeEntry{
						Name:      f.Name,
						Tree:      commitTree,
						TreeEntry: *e,
					},
				}

				changes = append(changes, ch)

				return nil
			})

			for _, value := range changes {

				patch, err := value.Patch()

				//from, to, err := c.Files()

				if err != nil {
					return nil, err
				}

				mekrlitre, err := value.Action()

				if err != nil {
					return nil, err
				}

				ActPaths = append(ActPaths, ActionPaths{
					Action:       string(mekrlitre.String()),
					Path:         patch.Stats().String(),
					PatchContent: patch.String(),
				})

			}
			ChangeFull := ChangeDescription{ActPaths, "", prevCommit.Message}

			return &ChangeFull, err
		}
		//Obtengo  el arbol de archivos del commit actual
		currentTree, err := commit.Tree()

		if err != nil {
			return nil, err
		}

		if prevCommit == nil {
			prevCommit = commit
			prevTree = currentTree
			continue
		}

		//fmt.Println(Sha1Id)
		//fmt.Println(prevCommit.Hash.String())

		if Sha1Id == prevCommit.Hash.String() {

			//	ContentDataTreeF, err := ContentTreeData(RepoName, )

			changes, err := currentTree.Diff(prevTree)

			if err != nil {
				return nil, err
			}

			for _, c := range changes {

				patch, err := c.Patch()

				//from, to, err := c.Files()

				if err != nil {
					return nil, err
				}

				mekrlitre, err := c.Action()

				if err != nil {
					return nil, err
				}

				ActPaths = append(ActPaths, ActionPaths{
					Action:       string(mekrlitre.String()),
					Path:         patch.Stats().String(),
					PatchContent: patch.String(),
				})

			}

			ChangeFull := ChangeDescription{ActPaths, "", prevCommit.Message}

			return &ChangeFull, err

		}

		prevCommit = commit
		prevTree = currentTree

	}

}

func getCommitNodeIndex(r *git.Repository, fs billy.Filesystem) (commitgraph.CommitNodeIndex, io.ReadCloser) {

	file, err := fs.Open(path.Join("objects", "info", "commit-graph"))

	if err == nil {
		//obtengo el index commit(nodo)
		index, err := commitgraph_fmt.OpenFileIndex(file)
		if err == nil {
			return commitgraph.NewGraphCommitNodeIndex(index, r.Storer), file
		}
		file.Close()
	}
	return commitgraph.NewObjectCommitNodeIndex(r.Storer), nil
}

func getcommitTree(c commitgraph.CommitNode, treePath string) (*object.Tree, error) {

	tree, err := c.Tree()

	if err != nil {
		return nil, err
	}

	if treePath != "" {
		tree, err = tree.Tree(treePath)

		if err != nil {
			return nil, err
		}
	}
	return tree, nil
}

func getFileHashes(c commitgraph.CommitNode, treePath string, paths []string) (map[string]plumbing.Hash, error) {

	tree, err := getcommitTree(c, treePath)

	if err == object.ErrDirectoryNotFound {
		return make(map[string]plumbing.Hash), nil
	}
	if err != nil {
		return nil, err
	}

	hashes := make(map[string]plumbing.Hash)

	for _, path := range paths {
		if path != "" {
			entry, err := tree.FindEntry(path)

			if err == nil {
				hashes[path] = plumbing.Hash(entry.Hash)
			}
		} else {
			hashes[path] = plumbing.Hash(tree.Hash)
		}
	}

	return hashes, nil
}

// GetLastCommitForPaths returns last commit information

func getLastCommitForPaths(c commitgraph.CommitNode, treePath string, paths []string) (map[string]*object.Commit, error) {

	// We do a tree traversal with nodes sorted by commit time
	heap := binaryheap.NewWith(func(a, b interface{}) int {
		if a.(*commitAndPaths).commit.CommitTime().Before(b.(*commitAndPaths).commit.CommitTime()) {
			return 1
		}
		return -1
	})

	resultNode := make(map[string]commitgraph.CommitNode)

	initialHashes, err := getFileHashes(c, treePath, paths)

	if err != nil {
		return nil, err
	}
	// Start search from the root commit and with full set of paths
	heap.Push(&commitAndPaths{c, paths, initialHashes})

	for {
		cIn, ok := heap.Pop()

		if !ok {
			break
		}
		current := cIn.(*commitAndPaths)

		// Load the parent commits for the one we are currently examining

		numParents := current.commit.NumParents()

		var parents []commitgraph.CommitNode

		for i := 0; i < numParents; i++ {

			parent, err := current.commit.ParentNode(i)

			if err != nil {
				break
			}
			parents = append(parents, parent)
		}

		// Examine the current commit and set of interesting paths

		pathUnchanged := make([]bool, len(current.paths))

		parentHashes := make([]map[string]plumbing.Hash, len(parents))

		for j, parent := range parents {

			parentHashes[j], err = getFileHashes(parent, treePath, current.paths)

			if err != nil {
				break
			}

			for i, path := range current.paths {
				if parentHashes[j][path] == current.hashes[path] {
					pathUnchanged[i] = true
				}
			}
		}

		var remainingPaths []string

		for i, path := range current.paths {

			if resultNode[path] == nil {

				if pathUnchanged[i] {
					remainingPaths = append(remainingPaths, path)
				} else {
					resultNode[path] = current.commit
				}
			}
		}

		if len(remainingPaths) > 0 {

			for j, parent := range parents {

				remainingPathsForParent := make([]string, 0, len(remainingPaths))

				newRemainingPaths := make([]string, 0, len(remainingPaths))

				for _, path := range remainingPaths {

					if parentHashes[j][path] == current.hashes[path] {

						remainingPathsForParent = append(remainingPathsForParent, path)

					} else {
						newRemainingPaths = append(newRemainingPaths, path)
					}
				}
				if remainingPathsForParent != nil {
					heap.Push(&commitAndPaths{parent, remainingPathsForParent, parentHashes[j]})
				}

				if len(newRemainingPaths) == 0 {
					break
				} else {
					remainingPaths = newRemainingPaths
				}
			}
		}
	}

	result := make(map[string]*object.Commit)

	for path, commitNode := range resultNode {

		var err error

		result[path], err = commitNode.Commit()

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
