package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/wibecoderr/storex/handler"
	"github.com/wibecoderr/storex/middleware"
)

func SetUpRoutes(r chi.Router) {
	handler.InitFirebase()

	// public routes
	r.Post("/register", handler.RegisterUser)
	r.Post("/login", handler.LoginUser)
	r.Post("/auth/google", handler.FirebaseLogin)
	r.Post("/auth/github", handler.FirebaseLogin)
	r.Post("/auth/linkedin", handler.FirebaseLogin)

	// authenticated routes (admin + member)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Post("/logout", handler.LogoutUser)

		// dashboard
		r.Get("/dashboard", handler.GetDashboard)

		// projects
		r.Get("/projects", handler.GetProjects)
		r.Get("/projects/{id}", handler.GetProjectByID)

		// tasks
		r.Get("/projects/{id}/tasks", handler.GetTasksByProject)
		r.Patch("/tasks/{id}/status", handler.UpdateTaskStatus)
	})

	// admin only routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.RoleMiddleware("admin"))

		r.Post("/register/employee", handler.CreateEmployee)

		// projects
		r.Post("/projects", handler.CreateProject)
		r.Put("/projects/{id}", handler.UpdateProject)
		r.Delete("/projects/{id}", handler.ArchiveProject)

		// members
		r.Post("/projects/{id}/members", handler.AddProjectMember)
		r.Delete("/projects/{id}/members/{uid}", handler.RemoveProjectMember)

		// tasks
		r.Post("/projects/{id}/tasks", handler.CreateTask)
		r.Put("/tasks/{id}", handler.UpdateTask)
		r.Patch("/tasks/{id}/assign", handler.AssignTask)
		r.Delete("/tasks/{id}", handler.ArchiveTask)
	})
}