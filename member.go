package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	utils "github.com/wibecoderr/storex"
	"github.com/wibecoderr/storex/database"
	"github.com/wibecoderr/storex/database/dbhelper"
	"github.com/wibecoderr/storex/middleware"
	"github.com/wibecoderr/storex/model"
)

func AddProjectMember(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	projectID := chi.URLParam(r, "id")

	var body model.AddMemberRequest
	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(body); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	err := database.Tx(func(tx *sqlx.Tx) error {
		return dbhelper.AddProjectMember(tx, projectID, body.UserID)
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to add member")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func RemoveProjectMember(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	projectID := chi.URLParam(r, "id")
	userID := chi.URLParam(r, "uid")

	err := database.Tx(func(tx *sqlx.Tx) error {
		return dbhelper.RemoveProjectMember(tx, projectID, userID)
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to remove member")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "success"})
}