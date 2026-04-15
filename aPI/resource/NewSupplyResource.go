package resource

import (
	"aIBuildService/aPI/models/products"
	"aIBuildService/aPI/service"
	"aIBuildService/aPI/utils"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

type SupplyResource struct {
	SupplyService service.SupplyService
}

func NewSupplyResource(supplyService service.SupplyService) *SupplyResource {
	return &SupplyResource{
		SupplyService: supplyService,
	}
}

func (app *SupplyResource) CreateSupplyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var supply products.Supply
	if err := utils.ParseJSON(r, &supply); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// validate supply
	if err := utils.Validate.Struct(&supply); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error", errors))
		return
	}
	// validate regex
	_, err := utils.ValidateSpaceRegex(supply.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field payload %s", err))
		return
	}
	// validate regex
	_, err = utils.ValidateSpaceRegex(supply.Country)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field payload payload %s", err))
		return
	}
	// check if name exist and case character match
	supplyName := utils.Escape(supply.Name)
	supplier, _ := app.SupplyService.Exists(supplyName)
	if supplyName == supplier.Name {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("category with that name %s already exists", supply.Name))
		return
	}
	// save supply
	save, err := app.SupplyService.Save(&supply)
	if err != nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, fmt.Errorf("invalid payload %s", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": save,
	})
}
func (app *SupplyResource) FindSupplyIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// get param url id
	params := r.PathValue("id")
	paramId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that supply id not avaliable"))
		return
	}
	// find supply request id
	supplyId, err := app.SupplyService.Find(int64(paramId))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": supplyId,
	})
}
func (app *SupplyResource) FindAllSupplyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	supplier, err := app.SupplyService.FindAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": supplier,
	})
}
func (app *SupplyResource) UpdateSupplyIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("id")
	paramId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that supply id not avaliable"))
		return
	}
	//var payload products.Supply
	supply := &products.Supply{}

	if err := utils.ParseJSON(r, &supply); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// validate regex
	_, err = utils.RegexValidate(supply.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field payload %s", err))
		return
	}
	// set Supply.Id = request = paramId
	supply.ID = int64(paramId)
	err = app.SupplyService.Update(supply)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": "Supply update successfully.",
	})
}
func (app *SupplyResource) DeleteSupplyIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("id")
	paramId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that supply id not avaliable"))
		return
	}
	err = app.SupplyService.Delete(int64(paramId))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that supply id not avaliable"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": "Delete supply successfully.",
	})
}
