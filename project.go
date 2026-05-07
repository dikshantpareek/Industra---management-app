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

func CreateProject(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	var body model.CreateProjectRequest
	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(body); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	var projectID string
	err := database.Tx(func(tx *sqlx.Tx) error {
		var err error
		projectID, err = dbhelper.CreateProject(tx, body.Name, body.Description, user.UserId)
		if err != nil {
			return err
		}
		return dbhelper.AddProjectMember(tx, projectID, user.UserId)
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create project")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{"project_id": projectID})
}

func GetProjects(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	projects, err := dbhelper.GetProjects(user.UserId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch projects")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    projects,
	})
}

func GetProjectByID(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	projectID := chi.URLParam(r, "id")

	project, err := dbhelper.GetProjectByID(projectID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err, "project not found")
		return
	}

	members, err := dbhelper.GetProjectMembers(projectID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch members")
		return
	}

	tasks, err := dbhelper.GetTasksByProject(projectID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch tasks")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "success",
		"data": map[string]interface{}{
			"project": project,
			"members": members,
			"tasks":   tasks,
		},
	})
}

func UpdateProject(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	projectID := chi.URLParam(r, "id")

	var body model.UpdateProjectRequest
	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(body); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	err := database.Tx(func(tx *sqlx.Tx) error {
		return dbhelper.UpdateProject(tx, projectID, body.Name, body.Description)
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update project")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func ArchiveProject(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	projectID := chi.URLParam(r, "id")

	if err := dbhelper.ArchiveProject(projectID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete project")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "success"})
}