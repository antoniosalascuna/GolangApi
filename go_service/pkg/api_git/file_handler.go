package api_git

import (
	"encoding/json"
	"fmt"
	"go_service/tools"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type MedContent struct {
	MediaContent string `string:"Name"`
}

type FileCreate struct {
	RepositoryName string `json:"RepositoryName,omitempty"`
	ProjectName    string `json:"ProjectName,omitempty"`
	ContentFile    string `json:"ContentFile,omitempty"`
	NameFile       string `json:"NameFile,omitempty"`
	CommiterName   string `json:"CommiterName,omitempty"`
	CommiterEmail  string `json:"CommiterEmail,omitempty"`
	CommitMessage  string `json:"CommitMessage,omitempty"`
}

func FileDetailsHandler(w http.ResponseWriter, r *http.Request) {
	var directory Directory
	var response tools.Response
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if body != nil {
		if err := json.Unmarshal(body, &directory); err != nil {
			//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		files, err := BlobDetails(directory.Name, directory.File)
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

func MediaFileHandler(w http.ResponseWriter, r *http.Request) {
	/*
		vars := mux.Vars(r)

		RepoName := vars["Name"]
		FilePath := vars["File"]
		fmt.Println(FilePath)
		fmt.Println(RepoName) */
	FilePath := r.FormValue("File")
	RepoName := r.FormValue("Name")
	ProjectName := r.FormValue("ProjectName")
	fmt.Println(FilePath)
	fmt.Println(RepoName)
	//var directory Directory
	//var response tools.Response
	//body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	/*if body != nil {
	if err := json.Unmarshal(body, &directory); err != nil {
		//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}*/
	MediaFile, err := GetMediContent(RepoName, ProjectName, FilePath)

	if err != nil {
		return
	}
	fmt.Println(MediaFile)

	w.Header().Set("Content-Type", "image/png")

	w.Write([]byte(MediaFile))

}

func MediaSvgHandler(w http.ResponseWriter, r *http.Request) {
	/*
		vars := mux.Vars(r)

		RepoName := vars["Name"]
		FilePath := vars["File"]
		fmt.Println(FilePath)
		fmt.Println(RepoName) */
	FilePath := r.FormValue("File")
	RepoName := r.FormValue("Name")
	ProjectName := r.FormValue("ProjectName")
	fmt.Println(FilePath)
	fmt.Println(RepoName)
	//var directory Directory
	//var response tools.Response
	//body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	/*if body != nil {
	if err := json.Unmarshal(body, &directory); err != nil {
		//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}*/
	MediaFile, err := GetMediContent(RepoName, ProjectName, FilePath)

	if err != nil {
		return
	}
	fmt.Println(MediaFile)

	w.Header().Set("Content-Type", "image/svg+xml")

	w.Write([]byte(MediaFile))

}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func GetFileReadmeHandler(w http.ResponseWriter, r *http.Request) {

	var response tools.Response

	name := r.FormValue("name")
	ProjectName := r.FormValue("ProjectName")
	File := r.FormValue("File")

	files, err := BlobReadmeDetails(name, ProjectName, File)
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

func CreateInitialFileHandler(w http.ResponseWriter, r *http.Request) {

	var response tools.Response
	var NewCreateFile FileCreate
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		panic(err)
	}

	if body != nil {

		if err := json.Unmarshal(body, &NewCreateFile); err != nil {
			//	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}

		if NewCreateFile.NameFile == "" {
			response.Message = "File Name Empty"
			response.Result = "Error"
			encodeData, _ := json.Marshal(response)
			fmt.Fprintf(w, string(encodeData))
			return
		} else {
			status, err := CreateInitialRepoFile(NewCreateFile)
			if err != nil {
				response.Message = err.Error()
				response.Result = "Error"
				encodeData, _ := json.Marshal(response)
				fmt.Fprintf(w, string(encodeData))
				return
			}
			if status != false {
				response.Message = "Success File Create"
				response.Result = "Success"
				encodeData, _ := json.Marshal(response)
				fmt.Fprintf(w, string(encodeData))
				return
			}
		}

	}

}
