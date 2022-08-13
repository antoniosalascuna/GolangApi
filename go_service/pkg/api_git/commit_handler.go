package api_git

import (
	"encoding/json"
	"fmt"
	"go_service/tools"

	"net/http"

	"github.com/gorilla/mux"
	//"github.com/go-git/go-billy/v5/osfs"
)

/*Obtengo todos las path de los commits que hay en el branch que este el repo*/
func GetAllCommitsHandler(w http.ResponseWriter, r *http.Request) {

	RepoName := r.FormValue("Name")

	GetAllSHACommits(RepoName)
}

func GetHistoryCommitsHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	RepoName := vars["RepoName"]
	ProjectName := vars["ProjectName"]
	var response tools.Response

	RepoHistory, err := GetCommitHistory(RepoName, ProjectName)

	if err != nil {
		response.Message = err.Error()
		response.Result = "Error"
		encodeData, _ := json.Marshal(response)
		fmt.Fprintf(w, string(encodeData))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encodeData, _ := json.Marshal(RepoHistory)
	w.Write(encodeData)
}

func GetCommitTreeHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	RepoName := vars["RepoName"]
	Sha1 := vars["Sha1"]
	ProjectName := vars["ProjectName"]
	var response tools.Response

	ChangeName, err := getTreeCommit(RepoName, Sha1, ProjectName)

	if err != nil {
		response.Message = err.Error()
		response.Result = "Error"
		encodeData, _ := json.Marshal(response)
		fmt.Fprintf(w, string(encodeData))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encodeData, _ := json.Marshal(ChangeName)
	w.Write(encodeData)
}

func GetDiffCommitsHandler(w http.ResponseWriter, r *http.Request) {

	var response tools.Response

	RepoName := r.FormValue("RepoName")
	BaseBranchNamePR := r.FormValue("BaseBranchNamePR")
	RefBranchNamePR := r.FormValue("RefBranchNamePR")
	ProjecName := r.FormValue("ProjecName")

	CommitsInfo, err := GetCommitDiff(RepoName, BaseBranchNamePR, RefBranchNamePR, ProjecName)

	if err != nil {
		response.Message = err.Error()
		response.Result = "Error"

	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encodeData, _ := json.Marshal(CommitsInfo)
	w.Write(encodeData)
	return

}
