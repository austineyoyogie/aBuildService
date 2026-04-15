package resource

import (
	"aIBuildService/aPI/models"
	"aIBuildService/aPI/service"
	"aIBuildService/aPI/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

type RoleResource struct {
	RoleService service.RoleService
}

func NewRoleResource(roleService service.RoleService) *RoleResource {
	return &RoleResource{
		RoleService: roleService,
	}
}

func (app *RoleResource) CreateRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload models.Role
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// validate payload
	if err := utils.Validate.Struct(&payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error", errors))
		return
	}
	// validate regex
	_, err := utils.RegexValidate(payload.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Field payload %s", err))
		return
	}
	// validate regex
	_, err = utils.RegexValidate(payload.Permissions)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Field payload payload %s", err))
		return
	}
	// check if role name exist
	role, _ := app.RoleService.Exists(payload.Name)
	if payload.Name == role.Name {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Role with that name %s already exists", payload.Name))
		return
	}
	// save user payload
	save, err := app.RoleService.Save(&payload)
	if err != nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, fmt.Errorf("Unvalid payload &v", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{"data": save})
}
func (app *RoleResource) FindRoleIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("id")
	roleId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that role id not avaliable"))
		return
	}
	rolePermission, err := app.RoleService.Find(int64(roleId))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": rolePermission,
	})
}
func (app *RoleResource) FindAllRoleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	role, err := app.RoleService.FindAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": role,
	})
}
func (app *RoleResource) UpdateRoleIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("id")
	roleId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that role id not avaliable"))
		return
	}
	//var payload models.Role
	rolePayload := &models.Role{}
	if err := utils.ParseJSON(r, &rolePayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// validate regex
	_, err = utils.RegexValidate(rolePayload.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("Field payload %s", err))
		return
	}
	// set Role.Id = request = roleId
	rolePayload.ID = int64(roleId)
	err = app.RoleService.Update(rolePayload)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": "Role update successfully.",
	})
}
func (app *RoleResource) DeleteRoleIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("id")
	roleId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that role id not avaliable"))
		return
	}
	err = app.RoleService.Delete(int64(roleId))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that role id not avaliable"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": "Delete role successfully.",
	})
}
