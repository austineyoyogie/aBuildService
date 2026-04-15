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

type CategoryResource struct {
	CategoryService service.CategoryService
}

func NewCategoryResource(categoryService service.CategoryService) *CategoryResource {
	return &CategoryResource{
		CategoryService: categoryService,
	}
}

func (app *CategoryResource) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var category products.Category
	if err := utils.ParseJSON(r, &category); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// validate category
	if err := utils.Validate.Struct(&category); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error", errors))
		return
	}
	// validate regex
	_, err := utils.ValidateSpaceRegex(category.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field name payload %s", err))
		return
	}
	// validate regex
	// check if name exist and case character match
	_, err = utils.ValidateDescriptionRegex(category.Description)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field description payload payload %s", err))
		return
	}
	// check if name exist and case character match
	categoryName := utils.Escape(category.Name)
	categories, _ := app.CategoryService.Exists(categoryName)
	if categoryName == categories.Name {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("category with that name %s already exists", category.Name))
		return
	}
	// save supply
	save, err := app.CategoryService.Save(&category)
	if err != nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, fmt.Errorf("invalid payload %s", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": save,
	})
}
func (app *CategoryResource) FindCategoryIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// get param url id
	params := r.PathValue("id")
	paramId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that category id not avaliable"))
		return
	}
	// find category request id
	categoryId, err := app.CategoryService.Find(int64(paramId))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": categoryId,
	})
}
func (app *CategoryResource) FindAllCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	categories, err := app.CategoryService.FindAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": categories,
	})
}
func (app *CategoryResource) UpdateCategoryIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("id")
	paramId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that category id not avaliable"))
		return
	}
	//var payload products Category
	category := &products.Category{}

	if err := utils.ParseJSON(r, &category); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// validate regex
	_, err = utils.ValidateSpaceRegex(category.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field name payload %s", err))
		return
	}
	// validate regex
	_, err = utils.ValidateDescriptionRegex(category.Description)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field description payload %s", err))
		return
	}
	// set Category.Id = request = paramId
	category.ID = int64(paramId)
	err = app.CategoryService.Update(category)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": "Category update successfully.",
	})
}
func (app *CategoryResource) DeleteCategoryIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("id")
	paramId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that category id not avaliable"))
		return
	}
	err = app.CategoryService.Delete(int64(paramId))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that category id not avaliable"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": "Delete category successfully.",
	})
}
