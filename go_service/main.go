package main

import (
	"fmt"
	"go_service/pkg/api_git"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Init router
	r := mux.NewRouter()
	fmt.Println("PRUEBA sirve !!!!!!")
	//r.HandleFunc("/", api_git.HomeHandler).Methods("GET")
	r.HandleFunc("/GetDiffCommitHistory", api_git.GetDiffCommitsHandler).Methods("GET")

	//r.HandleFunc("/CopyRepoFromTag", api_git.CopyRepoFromTagHandler).Methods("POST")
	r.HandleFunc("/CreateMerge", api_git.MergeHandler).Methods("POST")
	r.HandleFunc("/Diff", api_git.DiffHandler).Methods("POST")
	//	r.HandleFunc("/PullRequest", api_git.PRHandler).Methods("POST")

	//PullRequest
	r.HandleFunc("/CreatePr", api_git.InsertOnePRHandler).Methods("POST")
	r.HandleFunc("/GetPr", api_git.GetAllPr).Methods("GET")
	r.HandleFunc("/GetOnePr/{id}", api_git.GetPrByIDHandler).Methods("GET")
	r.HandleFunc("/UpdatePr/{id}", api_git.UpdatePr).Methods("PUT")
	r.HandleFunc("/UpdateStatusPr/{id}", api_git.UpdateStatPR).Methods("POST")
	r.HandleFunc("/DeletePr/{id}", api_git.DeleteOnePr).Methods("DELETE")
	r.HandleFunc("/LocalPullRequest", api_git.PRHandlerV2).Methods("POST")
	r.HandleFunc("/GetPrRepoNameProjectName/{RepoName}/{ProjectName}", api_git.GetPrByRepoNameAndProjectNameHandler).Methods("GET")

	//Tree Routes
	r.HandleFunc("/ListTreePathFiles", api_git.ListTreeFileHandler).Methods("POST", "OPTIONS")

	//Files
	r.HandleFunc(`/FileMediaContent`, api_git.MediaFileHandler).Methods("GET", "OPTIONS")
	r.HandleFunc(`/FileSvgContent`, api_git.MediaSvgHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/FileDetails", api_git.FileDetailsHandler).Methods("POST")
	r.HandleFunc("/ListPathFiles", api_git.ListFilesHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/CreateFile", api_git.CreateInitialFileHandler).Methods("POST", "OPTIONS")

	//Commits
	r.HandleFunc(`/GetAllCommits`, api_git.GetAllCommitsHandler).Methods("GET", "OPTIONS")
	r.HandleFunc(`/GetHistoryCommits/{ProjectName}/{RepoName}`, api_git.GetHistoryCommitsHandler).Methods("GET", "OPTIONS")

	//Cambiar por el branch {sha1} y enviarlo por un query params
	r.HandleFunc(`/{RepoName}/{ProjectName}/Commit/{Sha1}`, api_git.GetCommitTreeHandler).Methods("GET", "OPTIONS")

	//Branch
	r.HandleFunc(`/GetAllBranches/{ProjectName}/{RepoName}`, api_git.BranchesHandler).Methods("GET", "OPTIONS")
	r.HandleFunc(`/GetBranch/{RepoName}`, api_git.GetBranchHandler).Methods("GET", "OPTIONS")
	r.HandleFunc(`/CreateBranch`, api_git.CreateBranchHandler).Methods("POST", "OPTIONS")

	//Change
	r.HandleFunc(`/Change/{ProjectName}/{RepoName}`, api_git.ChangesTree).Methods("GET", "OPTIONS")

	//createRepository
	r.HandleFunc("/NewRepository", api_git.CreateRepositoryHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/Readme", api_git.GetFileReadmeHandler).Methods("GET", "OPTIONS")

	//CloneRepository
	r.HandleFunc("/CloneRepository", api_git.CloneRepositoryHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/CloneRepositoryWithUsernamePassword/{username}/{password}", api_git.CloneRepositoryWithUserNamePassWordHandler).Methods("POST", "OPTIONS")

	handler := cors.Default().Handler(r)
	// Start server
	log.Fatal(http.ListenAndServe(":3001", handler))

}
