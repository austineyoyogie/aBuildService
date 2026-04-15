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

type ProductResource struct {
	ProductService service.ProductService
}

func NewProductResource(productService service.ProductService) *ProductResource {
	return &ProductResource{
		ProductService: productService,
	}
}

func (app *ProductResource) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var product products.Product
	var productCatgory products.ProductCategory
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// validate product
	if err := utils.Validate.Struct(&product); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error", errors))
		return
	}
	// validate regex
	_, err := utils.ValidateSpaceRegex(product.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field payload %s", err))
		return
	}
	// validate regex
	// check if name exist and case character match
	_, err = utils.ValidateDescriptionRegex(product.Description)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field payload payload %s", err))
		return
	}
	// check if name exist and case character match
	productName := utils.Escape(product.Name)
	ifExistProduct, _ := app.ProductService.Exists(productName)
	if productName == ifExistProduct.Name {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("product with that name %s already exists", product.Name))
		return
	}
	// save product
	save, err := app.ProductService.Save(&product)
	if err != nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, fmt.Errorf("invalid save payload %s", err))
		return
	}
	// add to user role
	productCatgory.ProductId = uint(save.ID)
	_, err = app.ProductService.AddToProductCategory(&productCatgory)

	if err != nil {
		utils.WriteError(w, http.StatusUnprocessableEntity, fmt.Errorf("invalid product category payload entry %s", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": save,
	})
}
func (app *ProductResource) FindProductIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// get param url id
	params := r.PathValue("id")
	paramId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that product id not avaliable"))
		return
	}
	// find product request id
	productId, err := app.ProductService.Find(int64(paramId))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": productId,
	})
}
func (app *ProductResource) FindAllProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	allProducts, err := app.ProductService.FindAll()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": allProducts,
	})
}
func (app *ProductResource) UpdateProductIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("id")
	paramId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that product id not avaliable"))
		return
	}
	//var payload products
	product := &products.Product{}
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	// validate regex
	_, err = utils.ValidateSpaceRegex(product.Name)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field name payload %s", err))
		return
	}
	// validate regex
	_, err = utils.ValidateDescriptionRegex(product.Description)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("field description payload %s", err))
		return
	}
	// set Product.Id = request = paramId
	product.ID = int64(paramId)
	err = app.ProductService.Update(product)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": "Product update successfully.",
	})
}
func (app *ProductResource) DeleteProductIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("id")
	paramId, err := strconv.ParseUint(params, 10, 64)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that product id not avaliable"))
		return
	}
	err = app.ProductService.Delete(int64(paramId))
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that product id not avaliable"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": "Delete product successfully.",
	})
}
