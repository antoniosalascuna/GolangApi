package api_git

import (
	"encoding/json"
	"fmt"
	"go_service/tools"
	"io"
	"io/ioutil"

	//"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type BranchResBody struct {
	BranchName    string   `json:"BranchName,omitempty"`
	ProjectName   string   `json:"ProjectName,omitempty"`
	RepoName      string   `json:"RepoName,omitempty"`
	NewBranchName string   `json:"NewBranchName,omitempty"`
	HeadBranch    string   `json:"HeadBranch,omitempty"`
	Branches      []string `json:"Branches,omitempty"`
}

func BranchesHandler(w http.ResponseWriter, r *http.Request) {

	var response tools.Response
	BranchRes := BranchResBody{}
	vars := mux.Vars(r)
	RepoName := vars["RepoName"]
	ProjecName := vars["ProjectName"]

	branches, branchHead, err := GetBranches(RepoName, ProjecName)

	BranchRes.Branches = branches
	BranchRes.HeadBranch = branchHead

	if err != nil {
		response.Message = err.Error()
		response.Result = "Error"
		encodeData, _ := json.Marshal(response)
		fmt.Fprintf(w, string(encodeData))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encodeData, _ := json.Marshal(BranchRes)
	w.Write(encodeData)

}

func GetBranchHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	RepoName := vars["RepoName"]
	BranchName := r.FormValue("BranchName")
	ProjectName := r.FormValue("ProjectName")

	GetBranch(RepoName, BranchName, ProjectName)

	/* 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encodeData, _ := json.Marshal(BranchInfo)
	w.Write(encodeData) */

}

func ChangesTree(w http.ResponseWriter, r *http.Request) {

	//var response tools.Response
	vars := mux.Vars(r)
	RepoName := vars["RepoName"]
	ProjectName := vars["ProjectName"]

	err := InitialCommitChange(RepoName, ProjectName)
	if err != nil {
		panic(err)
	}
	/* 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encodeData, _ := json.Marshal(BranchInfo)
	w.Write(encodeData) */

}

func CreateBranchHandler(w http.ResponseWriter, r *http.Request) {

	var brachinfo BranchResBody
	var response tools.Response
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if body != nil {
		if err := json.Unmarshal(body, &brachinfo); err != nil {
			//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		err := createBranch(brachinfo.NewBranchName, brachinfo.ProjectName, brachinfo.RepoName, brachinfo.BranchName)
		if err != nil {
			response.Message = err.Error()
			response.Result = "Error"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}

	}

	/* 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	encodeData, _ := json.Marshal(BranchInfo)
	w.Write(encodeData) */

}
