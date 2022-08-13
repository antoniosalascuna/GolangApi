package api_git

import (
	"encoding/json"
	"fmt"
	"go_service/tools"
	"io"
	"io/ioutil"
	"net/http"
)

type RepoDetails struct {
	RepoName    string `json:"RepoName"`
	ProjectName string `json:"ProjectName"`
	Readme      bool   `json:"Readme"`
	Gitignore   bool   `json:"gitignore"`
	IsPublic    bool   `json:"isPublic"`
	UserName    string `json:"UserName"`
	UserEmail   string `json:"UserEmail"`
}

func CreateRepositoryHandler(w http.ResponseWriter, r *http.Request) {

	var repo RepoDetails

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	var response tools.Response

	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &repo); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		branchName, err := createRepository(repo.RepoName, repo.ProjectName, repo.Readme, repo.Gitignore, repo.IsPublic, repo.UserName, repo.UserEmail)

		if err != nil {
			response.Message = err.Error()
			response.Result = "Error"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}

		if len(branchName) == 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			response.Message = "Empty Repository"
			response.Result = "Success"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return

		}

		if len(branchName) != 0 {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			response.Message = "Repository Create"
			response.Result = "Success"
			encodeData, _ := json.Marshal(branchName)
			fmt.Fprintf(w, string(encodeData))
			return

		}

	}

}
