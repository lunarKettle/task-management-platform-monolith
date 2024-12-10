package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/transport/dto"
	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/usecases"
	"github.com/lunarKettle/task-management-platform-monolith/pkg/utils"
)

type ProjectHandlers struct {
	usecases *usecases.ProjectUseCases
}

func NewProjectHandlers(usecases *usecases.ProjectUseCases) *ProjectHandlers {
	return &ProjectHandlers{
		usecases: usecases,
	}
}

type handler = func(w http.ResponseWriter, r *http.Request) error

func (h *ProjectHandlers) RegisterRoutes(mux *http.ServeMux, errorHandler func(handler) http.Handler) {
	mux.Handle("GET /projects", errorHandler(h.getAllProjects))
	mux.Handle("GET /projects/{id}", errorHandler(h.getProject))
	mux.Handle("POST /projects", errorHandler(h.createProject))
	mux.Handle("PUT /projects", errorHandler(h.updateProject))
	mux.Handle("DELETE /projects/{id}", errorHandler(h.deleteProject))

	mux.Handle("GET /teams", errorHandler(h.getAllTeams))
	mux.Handle("GET /teams/{id}", errorHandler(h.getTeam))
	mux.Handle("POST /teams", errorHandler(h.createTeam))
	mux.Handle("PUT /teams", errorHandler(h.updateTeam))
	mux.Handle("DELETE /teams/{id}", errorHandler(h.deleteTeam))

	mux.Handle("GET /members", errorHandler(h.getAllMembers))

	mux.Handle("GET /tasks/{id}", errorHandler(h.getTask))
	mux.Handle("GET /tasks", errorHandler(h.getTasks))
	mux.Handle("POST /tasks", errorHandler(h.createTask))
	mux.Handle("PUT /tasks", errorHandler(h.updateTask))
	mux.Handle("DELETE /tasks/{id}", errorHandler(h.deleteTask))
}

func (h *ProjectHandlers) getAllProjects(w http.ResponseWriter, r *http.Request) error {
	projects, err := h.usecases.GetAllProjects(r.Context())

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		err = fmt.Errorf("failed to encode project to JSON: %w", err)
		return err
	}
	return err
}

func (h *ProjectHandlers) getProject(w http.ResponseWriter, r *http.Request) error {
	var id uint32

	id, err := utils.ExtractIDFromPath(r.URL.Path)
	if err != nil {
		return fmt.Errorf("failed to get id from url path: %w", err)
	}

	query := usecases.NewGetProjectByIDQuery(id)
	project, err := h.usecases.GetProjectByID(r.Context(), query)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(project); err != nil {
		err = fmt.Errorf("failed to encode project to JSON: %w", err)
		return err
	}
	return err
}

func (h *ProjectHandlers) createProject(w http.ResponseWriter, r *http.Request) error {
	var requestData dto.CreateProjectRequestDTO

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		err = fmt.Errorf("error while decoding request body: %w", err)
		return err
	}
	defer r.Body.Close()

	cmd := usecases.NewCreateProjectCommand(
		requestData.Name,
		requestData.Description,
		requestData.PlannedEndDate,
		requestData.Status,
		requestData.Priority,
		requestData.TeamId,
		requestData.Budget,
	)

	id, err := h.usecases.CreateProject(r.Context(), cmd)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"id": id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		err = fmt.Errorf("failed to encode project to JSON: %w", err)
		log.Print(err)
		return err
	}
	return err
}

func (h *ProjectHandlers) updateProject(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPut {
		return fmt.Errorf("invalid request method: %s", r.Method)
	}

	var requestData dto.UpdateProjectRequestDTO

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		err = fmt.Errorf("error while decoding request body: %w", err)
		return err
	}
	defer r.Body.Close()

	cmd := usecases.NewUpdateProjectCommand(
		requestData.ID,
		requestData.Name,
		requestData.Description,
		requestData.StartDate,
		requestData.PlannedEndDate,
		requestData.ActualEndDate,
		requestData.Status,
		requestData.Priority,
		requestData.TeamId,
		requestData.Budget,
	)

	err = h.usecases.UpdateProject(r.Context(), cmd)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *ProjectHandlers) deleteProject(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, err := strconv.ParseUint(parts[len(parts)-1], 10, 32)

	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return err
	}

	cmd := usecases.NewDeleteProjectCommand(uint32(id))
	err = h.usecases.DeleteProject(r.Context(), cmd)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *ProjectHandlers) getAllTeams(w http.ResponseWriter, r *http.Request) error {
	teams, err := h.usecases.GetAllTeams(r.Context())

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(teams); err != nil {
		err = fmt.Errorf("failed to encode team to JSON: %w", err)
		return err
	}
	return err
}

func (h *ProjectHandlers) getTeam(w http.ResponseWriter, r *http.Request) error {
	var id uint32

	id, err := utils.ExtractIDFromPath(r.URL.Path)
	if err != nil {
		return fmt.Errorf("failed to get id from url path: %w", err)
	}

	query := usecases.NewGetTeamByIDQuery(id)
	team, err := h.usecases.GetTeamByID(query)

	if err != nil {
		return err
	}

	membersDTO := make([]dto.MemberDTO, len(team.Members))

	for i, v := range team.Members {
		membersDTO[i] = memberModelToDTO(v)
	}

	responseData := dto.GetTeamResponseDTO{
		ID:        team.ID,
		Name:      team.Name,
		Members:   membersDTO,
		ManagerID: team.ManagerID,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(responseData); err != nil {
		err = fmt.Errorf("failed to encode team to JSON: %w", err)
		return err
	}
	return err
}

func (h *ProjectHandlers) createTeam(w http.ResponseWriter, r *http.Request) error {
	var requestData dto.CreateTeamRequestDTO

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		err = fmt.Errorf("error while decoding request body: %w", err)
		return err
	}
	defer r.Body.Close()

	members := make([]usecases.Member, len(requestData.Members))
	for i, v := range requestData.Members {
		members[i] = *usecases.NewMember(v.ID, v.Role)
	}

	cmd := usecases.NewCreateTeamCommand(
		requestData.Name,
		members,
		requestData.ManagerID,
	)

	id, err := h.usecases.CreateTeam(r.Context(), cmd)

	// добавить добавление участников команды через usecase

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"id": id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		err = fmt.Errorf("failed to encode team to JSON: %w", err)
		log.Print(err)
		return err
	}
	return err
}

func (h *ProjectHandlers) updateTeam(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPut {
		return fmt.Errorf("invalid request method: %s", r.Method)
	}

	var requestData dto.UpdateTeamRequestDTO

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		err = fmt.Errorf("error while decoding request body: %w", err)
		return err
	}
	defer r.Body.Close()

	members := make([]usecases.Member, len(requestData.Members))
	for i, v := range requestData.Members {
		members[i] = *usecases.NewMember(v.ID, v.Role)
	}

	cmd := usecases.NewUpdateTeamCommand(
		requestData.ID,
		requestData.Name,
		members,
		requestData.ManagerID,
	)

	// добавить изменение участников

	err = h.usecases.UpdateTeam(r.Context(), cmd)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *ProjectHandlers) deleteTeam(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, err := strconv.ParseUint(parts[len(parts)-1], 10, 32)

	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return err
	}

	cmd := usecases.NewDeleteTeamCommand(uint32(id))
	err = h.usecases.DeleteTeam(r.Context(), cmd)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *ProjectHandlers) getAllMembers(w http.ResponseWriter, r *http.Request) error {
	members, err := h.usecases.GetAllMembers(r.Context())

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(members); err != nil {
		err = fmt.Errorf("failed to encode team to JSON: %w", err)
		return err
	}
	return err
}

func (h *ProjectHandlers) getTask(w http.ResponseWriter, r *http.Request) error {
	id, err := utils.ExtractIDFromPath(r.URL.Path)
	if err != nil {
		return fmt.Errorf("failed to extract id: %w", err)
	}

	query := usecases.NewGetTaskByIDQuery(id)
	task, err := h.usecases.GetTaskByID(r.Context(), query)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		return fmt.Errorf("failed to encode task to JSON: %w", err)
	}

	return nil
}

func (h *ProjectHandlers) getTasks(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()

	employeeID, _ := strconv.Atoi(query.Get("employee_id"))
	projectID, _ := strconv.Atoi(query.Get("project_id"))
	isCompleted := query.Get("is_completed")

	filter := usecases.TaskFilter{
		EmployeeID:  uint32(employeeID),
		ProjectID:   uint32(projectID),
		IsCompleted: parseBool(isCompleted),
	}

	tasks, err := h.usecases.GetTasks(filter)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		return fmt.Errorf("failed to encode task to JSON: %w", err)
	}

	return nil
}

func (h *ProjectHandlers) createTask(w http.ResponseWriter, r *http.Request) error {
	var requestData dto.CreateTaskRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		return fmt.Errorf("error decoding request body: %w", err)
	}
	defer r.Body.Close()

	cmd := usecases.NewCreateTaskCommand(
		requestData.Description,
		requestData.EmployeeID,
		requestData.ProjectID,
		requestData.IsCompleted,
	)

	id, err := h.usecases.CreateTask(r.Context(), cmd)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{"id": id}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}

	return nil
}

func (h *ProjectHandlers) updateTask(w http.ResponseWriter, r *http.Request) error {
	var requestData dto.UpdateTaskRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		return fmt.Errorf("error decoding request body: %w", err)
	}
	defer r.Body.Close()

	cmd := usecases.NewUpdateTaskCommand(
		requestData.ID,
		requestData.Description,
		requestData.EmployeeID,
		requestData.ProjectID,
		requestData.IsCompleted,
	)

	if err := h.usecases.UpdateTask(r.Context(), cmd); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *ProjectHandlers) deleteTask(w http.ResponseWriter, r *http.Request) error {
	id, err := utils.ExtractIDFromPath(r.URL.Path)
	if err != nil {
		return fmt.Errorf("failed to extract id: %w", err)
	}

	cmd := usecases.NewDeleteTaskCommand(id)
	if err := h.usecases.DeleteTask(r.Context(), cmd); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
