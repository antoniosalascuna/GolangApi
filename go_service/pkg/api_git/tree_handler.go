package api_git

import (
	"encoding/json"
	"fmt"
	"go_service/tools"
	"io"
	"io/ioutil"
	"net/http"
	//"os"
)

type Directory struct {
	Name        string `json:"name"`
	ProjectName string `json:"ProjectName"`
	File        string `json:"file"`
}
type JsonData struct {
	Data string `json:"Data"`
}

func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	var directory Directory
	var response tools.Response
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &directory); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		files, err := ListPathFileRepository(directory.Name)
		if err != nil {
			response.Message = err.Error()
			response.Result = "Error"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}
		if files != nil {
			encodeData, _ := json.Marshal(files)
			fmt.Fprintf(w, string(encodeData))
			return
		}

	}
}

func ListTreeFileHandler(w http.ResponseWriter, r *http.Request) {
	var directory Directory
	var response tools.Response
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &directory); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		files, err := ContentTreeData(directory.ProjectName, directory.Name, directory.File)

		if err != nil {
			response.Message = err.Error()
			response.Result = "Error"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		encodeData, _ := json.Marshal(files)
		w.Write(encodeData)

		return

	}
}

/*
Method Http Content Blob File

func ListDataContentFile(w http.ResponseWriter, r *http.Request) {

	var directory Directory
	var response tools.Response

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &directory); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
		fileContent, err := examples.ListContenBlobFile(directory.Name, directory.File)

		//	branches, err := api_git.getBranches(directory.Name, 1, 1)
		if err != nil {
			response.Message = err.Error()
			response.Result = "Error"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}
		if fileContent != nil {
			encodeData, _ := json.Marshal(fileContent)
			fmt.Fprintf(w, string(encodeData))
			return
		}
	}

}*/

/* func JsonContenFileHandler(w http.ResponseWriter, r *http.Request) {
	var directory Directory
	var response tools.Response
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &directory); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		files, err := JsonContentData(directory.Name, directory.File)
		if err != nil {
			response.Message = err.Error()
			response.Result = "Error"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		}
		if files != nil {

			//encodeData, _ := json.Marshal(files)
			w.Write([]byte(files.DataContent))
			//fmt.Println(string(encodeData))
			//	fmt.Fprintf(w, string(encodeData))
			return
		}

	}
}*/
