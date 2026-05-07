package handler

import (
	"net/http"

	utils "github.com/wibecoderr/storex"
	"github.com/wibecoderr/storex/database/dbhelper"
	"github.com/wibecoderr/storex/middleware"
)

func GetDashboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	dashboard, err := dbhelper.GetDashboard(user.UserId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch dashboard")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    dashboard,
	})
}