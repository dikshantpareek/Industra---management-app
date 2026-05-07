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

func CreateTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	projectID := chi.URLParam(r, "id")

	var body model.CreateTaskRequest
	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(body); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	var taskID string
	err := database.Tx(func(tx *sqlx.Tx) error {
		var err error
		taskID, err = dbhelper.CreateTask(tx, body.Title, body.Description, projectID, body.AssigneeID, body.DueDate)
		return err
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to create task")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]string{"task_id": taskID})
}

func GetTasksByProject(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	projectID := chi.URLParam(r, "id")

	tasks, err := dbhelper.GetTasksByProject(projectID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to fetch tasks")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "success",
		"data":    tasks,
	})
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	taskID := chi.URLParam(r, "id")

	var body model.UpdateTaskRequest
	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(body); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	err := database.Tx(func(tx *sqlx.Tx) error {
		return dbhelper.UpdateTask(tx, taskID, body.Title, body.Description, body.DueDate)
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update task")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	taskID := chi.URLParam(r, "id")

	var body model.UpdateStatusRequest
	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(body); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	if err := dbhelper.UpdateTaskStatus(taskID, body.Status); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to update status")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func AssignTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	taskID := chi.URLParam(r, "id")

	var body model.AssignTaskRequest
	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, "failed to parse body")
		return
	}
	if errs := utils.ValidateStruct(body); errs != nil {
		utils.RespondValidationError(w, errs)
		return
	}

	err := database.Tx(func(tx *sqlx.Tx) error {
		return dbhelper.AssignTask(tx, taskID, body.AssigneeID)
	})
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to assign task")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "success"})
}

func ArchiveTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.UserContext(r)
	if user == nil {
		utils.RespondError(w, http.StatusUnauthorized, nil, "unauthorized")
		return
	}

	taskID := chi.URLParam(r, "id")

	if err := dbhelper.ArchiveTask(taskID); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to delete task")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "success"})
}