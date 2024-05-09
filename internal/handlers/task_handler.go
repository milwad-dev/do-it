package handlers

import (
	"github.com/milwad-dev/do-it/internal/utils"
	"net/http"
)

func (db *DBHandler) GetLatestTasks(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]string)

	utils.JsonResponse(w, data)
}
