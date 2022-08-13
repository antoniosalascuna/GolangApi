package api_git

import (
	"encoding/json"
	"fmt"
	"go_service/tools"
	"io"
	"io/ioutil"
	"net/http"
)

type MergeRequest struct {
	RepoName     string `json:"RepoName"`
	Branch       string `json:"branch"`
	TargetBranch string `json:"targetBranch"`
	ProjecName   string `json:"ProjecName"`
}

func MergeHandler(w http.ResponseWriter, r *http.Request) {
	var mergeRequest MergeRequest
	var response tools.Response
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		response.Message = err.Error()
		response.Result = "Error"
	}
	if body != nil {
		if err := json.Unmarshal(body, &mergeRequest); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		result, err := CreateMerge(mergeRequest)
		if err != nil {
			response.Message = err.Error()
			response.Result = "Error"
		} else if result != "" {
			response.Message = result
			response.Result = "Success"
		}

	}

	encodeData, _ := json.Marshal(response)
	fmt.Fprintf(w, string(encodeData))
	return
}
