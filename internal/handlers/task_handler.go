package handlers

import (
	"github.com/milwad-dev/do-it/internal/utils"
	"net/http"
)

func (db *DBHandler) GetLatestTasks(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]string)

	query := "SELECT * FROM tasks ORDER BY created_at DESC"
	rows, err := db.Query(query)
	if err != nil {
		data["message"] = err.Error()
		utils.JsonResponse(w, data)

		return
	}

	for rows.Next() {
		
	}

	utils.JsonResponse(w, data)
}
