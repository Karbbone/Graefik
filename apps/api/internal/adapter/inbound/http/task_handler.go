package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/Karbbone/Graefik/apps/api/internal/core/domain"
	"github.com/Karbbone/Graefik/apps/api/internal/core/port"
	"github.com/labstack/echo/v5"
)

// TaskHandler traduit les requêtes HTTP vers le port inbound TaskService.
type TaskHandler struct {
	service port.TaskService
}

// NewTaskHandler crée le handler à partir du service (port inbound).
func NewTaskHandler(service port.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// --- DTO (contrat HTTP, découplé du domaine) ---

type createTaskRequest struct {
	Title string `json:"title"`
}

type taskResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"createdAt"`
}

func toTaskResponse(t *domain.Task) taskResponse {
	return taskResponse{
		ID:        t.ID,
		Title:     t.Title,
		Done:      t.Done,
		CreatedAt: t.CreatedAt.UTC().Format(time.RFC3339),
	}
}

// List — GET /api/tasks
func (h *TaskHandler) List(c *echo.Context) error {
	tasks, err := h.service.List(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	res := make([]taskResponse, 0, len(tasks))
	for _, t := range tasks {
		res = append(res, toTaskResponse(t))
	}
	return c.JSON(http.StatusOK, res)
}

// Create — POST /api/tasks
func (h *TaskHandler) Create(c *echo.Context) error {
	var req createTaskRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "corps de requête invalide")
	}

	task, err := h.service.Create(c.Request().Context(), req.Title)
	if err != nil {
		if errors.Is(err, domain.ErrTaskTitleRequired) {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, toTaskResponse(task))
}
