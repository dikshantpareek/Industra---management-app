package dbhelper

import (
	"github.com/jmoiron/sqlx"
	"github.com/wibecoderr/storex/database"
	"github.com/wibecoderr/storex/model"
)

// ─── PROJECTS ───────────────────────────────────────────────

func CreateProject(tx *sqlx.Tx, name, description, ownerID string) (string, error) {
	sql := `INSERT INTO projects (name, description, owner_id)
            VALUES ($1, $2, $3)
            RETURNING id`
	var id string
	err := tx.Get(&id, sql, name, description, ownerID)
	return id, err
}

func GetProjects(userID string) ([]model.Project, error) {
	sql := `SELECT p.id, p.name, p.description, p.owner_id, p.created_at
            FROM projects p
            INNER JOIN project_members pm ON pm.project_id = p.id
            WHERE pm.user_id = $1 AND p.archived_at IS NULL`
	var projects []model.Project
	err := database.DB.Select(&projects, sql, userID)
	return projects, err
}

func GetProjectByID(projectID string) (model.Project, error) {
	sql := `SELECT id, name, description, owner_id, created_at
            FROM projects
            WHERE id = $1 AND archived_at IS NULL`
	var project model.Project
	err := database.DB.Get(&project, sql, projectID)
	return project, err
}

func UpdateProject(tx *sqlx.Tx, projectID, name, description string) error {
	sql := `UPDATE projects SET name = $1, description = $2
            WHERE id = $3 AND archived_at IS NULL`
	_, err := tx.Exec(sql, name, description, projectID)
	return err
}

func ArchiveProject(projectID string) error {
	sql := `UPDATE projects SET archived_at = now()
            WHERE id = $1 AND archived_at IS NULL`
	_, err := database.DB.Exec(sql, projectID)
	return err
}

// ─── PROJECT MEMBERS ────────────────────────────────────────

func AddProjectMember(tx *sqlx.Tx, projectID, userID string) error {
	sql := `INSERT INTO project_members (project_id, user_id)
            VALUES ($1, $2)
            ON CONFLICT DO NOTHING`
	_, err := tx.Exec(sql, projectID, userID)
	return err
}

func RemoveProjectMember(tx *sqlx.Tx, projectID, userID string) error {
	sql := `DELETE FROM project_members
            WHERE project_id = $1 AND user_id = $2`
	_, err := tx.Exec(sql, projectID, userID)
	return err
}

func GetProjectMembers(projectID string) ([]model.ProjectMember, error) {
	sql := `SELECT u.id, u.name, u.email, u.role
            FROM users u
            INNER JOIN project_members pm ON pm.user_id = u.id
            WHERE pm.project_id = $1 AND u.archived_at IS NULL`
	var members []model.ProjectMember
	err := database.DB.Select(&members, sql, projectID)
	return members, err
}

func IsProjectMember(projectID, userID string) (bool, error) {
	sql := `SELECT count(*) > 0 FROM project_members
            WHERE project_id = $1 AND user_id = $2`
	var exists bool
	err := database.DB.Get(&exists, sql, projectID, userID)
	return exists, err
}

// ─── TASKS ──────────────────────────────────────────────────

func CreateTask(tx *sqlx.Tx, title, description, projectID, assigneeID, dueDate string) (string, error) {
	sql := `INSERT INTO tasks (title, description, project_id, assignee_id, due_date)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id`
	var id string
	err := tx.Get(&id, sql, title, description, projectID, assigneeID, dueDate)
	return id, err
}

func GetTasksByProject(projectID string) ([]model.Task, error) {
	sql := `SELECT id, title, description, project_id, assignee_id, status, due_date, created_at
            FROM tasks
            WHERE project_id = $1 AND archived_at IS NULL
            ORDER BY created_at DESC`
	var tasks []model.Task
	err := database.DB.Select(&tasks, sql, projectID)
	return tasks, err
}

func GetTaskByID(taskID string) (model.Task, error) {
	sql := `SELECT id, title, description, project_id, assignee_id, status, due_date, created_at
            FROM tasks
            WHERE id = $1 AND archived_at IS NULL`
	var task model.Task
	err := database.DB.Get(&task, sql, taskID)
	return task, err
}

func UpdateTask(tx *sqlx.Tx, taskID, title, description, dueDate string) error {
	sql := `UPDATE tasks SET title = $1, description = $2, due_date = $3
            WHERE id = $4 AND archived_at IS NULL`
	_, err := tx.Exec(sql, title, description, dueDate, taskID)
	return err
}

func UpdateTaskStatus(taskID, status string) error {
	sql := `UPDATE tasks SET status = $1
            WHERE id = $2 AND archived_at IS NULL`
	_, err := database.DB.Exec(sql, status, taskID)
	return err
}

func AssignTask(tx *sqlx.Tx, taskID, assigneeID string) error {
	sql := `UPDATE tasks SET assignee_id = $1
            WHERE id = $2 AND archived_at IS NULL`
	_, err := tx.Exec(sql, assigneeID, taskID)
	return err
}

func ArchiveTask(taskID string) error {
	sql := `UPDATE tasks SET archived_at = now()
            WHERE id = $1 AND archived_at IS NULL`
	_, err := database.DB.Exec(sql, taskID)
	return err
}

// ─── DASHBOARD ──────────────────────────────────────────────

func GetDashboard(userID string) (model.Dashboard, error) {
	sql := `SELECT
              COUNT(*) FILTER (WHERE t.assignee_id = $1)                                        AS my_tasks,
              COUNT(*) FILTER (WHERE t.status = 'todo' AND t.assignee_id = $1)                  AS todo,
              COUNT(*) FILTER (WHERE t.status = 'in_progress' AND t.assignee_id = $1)           AS in_progress,
              COUNT(*) FILTER (WHERE t.status = 'done' AND t.assignee_id = $1)                  AS done,
              COUNT(*) FILTER (WHERE t.due_date < now() AND t.status != 'done' AND t.assignee_id = $1) AS overdue
            FROM tasks t
            INNER JOIN project_members pm ON pm.project_id = t.project_id
            WHERE pm.user_id = $1 AND t.archived_at IS NULL`
	var dashboard model.Dashboard
	err := database.DB.Get(&dashboard, sql, userID)
	return dashboard, err
}