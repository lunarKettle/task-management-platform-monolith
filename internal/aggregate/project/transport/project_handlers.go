package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/transport/dto"
	"github.com/lunarKettle/task-management-platform-monolith/internal/aggregate/project/usecases"
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

func (h *ProjectHandlers) RegisterRoutes(mux *http.ServeMux, eh func(handler) http.Handler) {
	mux.Handle("GET /projects/{id}", eh(h.getProject))
	mux.Handle("POST /projects", eh(h.createProject))
	mux.Handle("PUT /projects", eh(h.updateProject))
	mux.Handle("DELETE /projects/{id}", eh(h.deleteProject))

	mux.Handle("GET /teams/{id}", eh(h.getTeam))
	mux.Handle("POST /teams", eh(h.createTeam))
	mux.Handle("PUT /teams", eh(h.updateTeam))
	mux.Handle("DELETE /teams/{id}", eh(h.deleteTeam))
}

func (h *ProjectHandlers) getProject(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, err := strconv.ParseUint(parts[len(parts)-1], 10, 32)

	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return err
	}

	query := usecases.NewGetProjectByIDQuery(uint32(id))
	project, err := h.usecases.GetProjectByID(query)

	if err != nil {
		return err
	}

	if err := json.NewEncoder(w).Encode(project); err != nil {
		err = fmt.Errorf("failed to encode project to JSON: %w", err)
		return err
	}
	return err
}

func (h *ProjectHandlers) createProject(w http.ResponseWriter, r *http.Request) error {
	var requestData struct {
		Name           string    `json:"name"`
		Description    string    `json:"description"`
		PlannedEndDate time.Time `json:"planned_end_date"`
		Status         string    `json:"status"`
		Priority       uint32    `json:"priority"`
		TeamId         uint32    `json:"team_id"`
		Budget         float64   `json:"budget"`
	}

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

	id, err := h.usecases.CreateProject(cmd)

	if err != nil {
		return err
	}

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

	var requestData struct {
		ID             uint32    `json:"id"`
		Name           string    `json:"name"`
		Description    string    `json:"description"`
		StartDate      time.Time `json:"start_date"`
		PlannedEndDate time.Time `json:"planned_end_date"`
		ActualEndDate  time.Time `json:"actual_end_date"`
		Status         string    `json:"status"`
		Priority       uint32    `json:"priority"`
		TeamId         uint32    `json:"team_id"`
		Budget         float64   `json:"budget"`
	}

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

	err = h.usecases.UpdateProject(cmd)

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
	err = h.usecases.DeleteProject(cmd)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}

func (h *ProjectHandlers) getTeam(w http.ResponseWriter, r *http.Request) error {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id, err := strconv.ParseUint(parts[len(parts)-1], 10, 32)

	if err != nil {
		http.Error(w, "invalid team id", http.StatusBadRequest)
		return err
	}

	query := usecases.NewGetTeamByIDQuery(uint32(id))
	team, err := h.usecases.GetTeamByID(query)

	if err != nil {
		return err
	}

	membersDTO := make([]dto.MemberDTO, len(team.Members))

	for i, v := range team.Members {
		membersDTO[i] = memberModelToDTO(v)
	}

	responseData := dto.GetTeamResponseDTO{
		ID:      team.ID,
		Name:    team.Name,
		Members: membersDTO,
	}

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
	)

	id, err := h.usecases.CreateTeam(cmd)

	// добавить добавление участников команды через usecase

	if err != nil {
		return err
	}

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
	)

	// добавить изменение участников

	err = h.usecases.UpdateTeam(cmd)

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
	err = h.usecases.DeleteTeam(cmd)

	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)

	return nil
}
