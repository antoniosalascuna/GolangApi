package api_git

import (
	"encoding/json"
	"fmt"
	"go_service/tools"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type CloneRepoRequest struct {
	Url         string `json:"Url,omitempty"`
	ProjectName string `json:"ProjectName,omitempty"`
	RepoName    string `json:"RepoName,omitempty"`
}

func CloneRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	var CloneRepo CloneRepoRequest
	var response tools.Response
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &CloneRepo); err != nil {
			//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
		if CloneRepo.Url == "" {

			w.WriteHeader(http.StatusBadRequest)
			response.Message = "Provide a valid URL"
			response.Result = "Error"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		} else {

			status, err := CloneRepository(CloneRepo.RepoName, CloneRepo.ProjectName, CloneRepo.Url)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				response.Message = err.Error()
				response.Result = "Error"
				encodeData, _ := json.Marshal(response)
				fmt.Fprintf(w, string(encodeData))
				return
			}
			if status != nil {

				response.Message = "Repository has been Cloned Successfully"
				response.Result = "Success"
				encodeData, _ := json.Marshal(response)
				fmt.Fprintf(w, string(encodeData))
				return
			}
		}

	}
}

func CloneRepositoryWithUserNamePassWordHandler(w http.ResponseWriter, r *http.Request) {
	var CloneRepo CloneRepoRequest
	vars := mux.Vars(r)
	Username := vars["username"]
	Password := vars["password"]

	var response tools.Response
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &CloneRepo); err != nil {
			//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		status, err := CloneRepositoryWithUserNamePassword(CloneRepo.RepoName, CloneRepo.ProjectName, CloneRepo.Url, Username, Password)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			response.Message = err.Error()
			response.Result = "Error"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}
		if status {

			response.Message = "Repository has been Cloned Successfully"
			response.Result = "Success"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}

	}
}
