package api_git

import (
	"fmt"
	"go_service/tools"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Home")

	tools.ConnectionDB()

}
