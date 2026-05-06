// model/project.go
package model

import "time"

type Project struct {
	ID          string    `db:"id"          json:"id"`
	Name        string    `db:"name"        json:"name"`
	Description string    `db:"description" json:"description"`
	OwnerID     string    `db:"owner_id"    json:"owner_id"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
}

type ProjectMember struct {
	ID    string `db:"id"    json:"id"`
	Name  string `db:"name"  json:"name"`
	Email string `db:"email" json:"email"`
	Role  string `db:"role"  json:"role"`
}

// model/task.go
package model

import "time"

type Task struct {
	ID          string     `db:"id"          json:"id"`
	Title       string     `db:"title"       json:"title"`
	Description string     `db:"description" json:"description"`
	ProjectID   string     `db:"project_id"  json:"project_id"`
	AssigneeID  string     `db:"assignee_id" json:"assignee_id"`
	Status      string     `db:"status"      json:"status"`
	DueDate     *time.Time `db:"due_date"    json:"due_date"`
	CreatedAt   time.Time  `db:"created_at"  json:"created_at"`
}

// model/dashboard.go
package model

type Dashboard struct {
	MyTasks    int `db:"my_tasks"    json:"my_tasks"`
	Todo       int `db:"todo"        json:"todo"`
	InProgress int `db:"in_progress" json:"in_progress"`
	Done       int `db:"done"        json:"done"`
	Overdue    int `db:"overdue"     json:"overdue"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"        validate:"required"`
	Description string `json:"description"`
}

type AddMemberRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

type CreateTaskRequest struct {
	Title       string `json:"title"       validate:"required"`
	Description string `json:"description"`
	AssigneeID  string `json:"assignee_id" validate:"required,uuid"`
	DueDate     string `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"       validate:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=todo in_progress done"`
}

type AssignTaskRequest struct {
	AssigneeID string `json:"assignee_id" validate:"required,uuid"`
}
type UserCxt struct {
    UserId    string
    SessionId string
}