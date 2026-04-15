package resource

import (
	"aIBuildService/aPI/middleware"
	"aIBuildService/aPI/models"
	"aIBuildService/aPI/service"
	"aIBuildService/aPI/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

const MAXIMUM = 3
const AWAIT = 5

var (
	MethodNotAllowed = errors.New("Method not allowed ")
)

type UserResource struct {
	UserService service.UserService
	TokenMaker  *middleware.JWTMaker
}

func NewUserResource(userService service.UserService, secretKey string) *UserResource {
	return &UserResource{
		UserService: userService,
		TokenMaker:  middleware.NewJWTMaker(secretKey),
	}
}

func (app *UserResource) UserResourceHandlerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /dashboard", middleware.Logger(middleware.AuthMiddlewareFunc(app.TokenMaker)(http.HandlerFunc(app.WelcomeGetHandler))))
	mux.HandleFunc("POST /login", middleware.Logger(http.HandlerFunc(app.LoginUserHandler)))
	mux.HandleFunc("POST /register", middleware.Logger(http.HandlerFunc(app.RegisterUserHandler)))
	mux.HandleFunc("GET /verify/account", middleware.Logger(http.HandlerFunc(app.GetUserTokenEmailVerificationHandler)))
	mux.HandleFunc("POST /reset/password", middleware.Logger(http.HandlerFunc(app.ResetUserPasswordTokenHandler)))
	mux.HandleFunc("GET /verify/password", middleware.Logger(http.HandlerFunc(app.VerifyUserPasswordTokenHandler)))
	mux.HandleFunc("POST /update/password", middleware.Logger(http.HandlerFunc(app.UpdateUserPasswordTokenHandler)))
	mux.HandleFunc("GET /refresh/{token}", middleware.Logger(http.HandlerFunc(app.RenewAccessTokenHandler)))
	mux.HandleFunc("POST /revoke/{id}", middleware.Logger(http.HandlerFunc(app.RevokedAccessTokenHandler)))
	mux.HandleFunc("GET /profile", middleware.Logger(middleware.AuthMiddlewareFunc(app.TokenMaker)(http.HandlerFunc(app.GetAllUsersHandler))))
	mux.HandleFunc("GET /show/{uuid}", middleware.Logger(middleware.AuthMiddlewareFunc(app.TokenMaker)(http.HandlerFunc(app.GetUserByIdHandler))))
	mux.HandleFunc("PUT /update/profile/{uuid}", middleware.Logger(middleware.AuthMiddlewareFunc(app.TokenMaker)(http.HandlerFunc(app.UpdateUserByIdHandler))))
	mux.HandleFunc("DELETE /delete/profile/{uuid}", middleware.Logger(middleware.AuthMiddlewareFunc(app.TokenMaker)(http.HandlerFunc(app.DeleteUserByIdHandler))))
	mux.HandleFunc("POST /blacklist/{token}", middleware.Logger(http.HandlerFunc(app.BlackListRefreshTokenHandler)))
}
func (app *UserResource) WelcomeGetHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Welcome Resource"))
	if err != nil {
		return
	}
}
func (app *UserResource) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"error": "Error - Method not allowed",
		})
		return
	}
	var payload models.CreateUserRequest
	var userRoles models.UserRole

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "Error - Invalid body: An error occurred.",
		})
		return
	}
	if err := utils.Validate.Struct(&payload); err != nil {
		_ = err.(validator.ValidationErrors)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "Error - Invalid. All fields is required.",
		})
		return
	}
	_, err := utils.RegexValidate(payload.FirstName)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
			"error": "Error - That field an invalid characters.",
		})
		return
	}
	_, err = utils.RegexValidate(payload.LastName)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
			"error": "Error - That field an invalid characters.",
		})
		return
	}
	_, err = utils.RegexValidatePhone(payload.PhoneNumber)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
			"error": "Error - That field an invalid phone numbers.",
		})
		return
	}
	existsUser, _ := app.UserService.GetUserByEmail(payload.Email)
	if payload.Email == existsUser.Email {
		utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
			"error": "Error - User with that email already exists.",
		})
		return
	}
	// Still in development =>TRY TO CHANGE TO USER REQUEST OR REPS
	user := models.User{
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		Email:       payload.Email,
		Password:    payload.Password,
		PhoneNumber: payload.PhoneNumber,
	}
	userSaveId, err := app.UserService.CreateUser(&user)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "Error - Bad status unprocessable entity.",
		})
		return
	}
	userRoles.UserId = userSaveId.ID
	_, err = app.UserService.AddToUserRole(&userRoles)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "Error - Bad status unprocessable entity.",
		})
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Success - Account has been created successfully.",
		"data":    userSaveId,
	})
}
func (app *UserResource) LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"error": "Error - Method not allowed",
		})
		return
	}
	var user models.User
	var credentials models.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utils.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "Error - Invalid body: An error occurred.",
		})
		return
	}

	if err := utils.Validate.Struct(&credentials); err != nil {
		_ = err.(validator.ValidationErrors)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "Error - Invalid. All fields is required.",
		})
		return
	}
	existsUser, _ := app.UserService.GetUserByEmail(credentials.Email)
	if existsUser.Email == user.Email {
		utils.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "Invalid - Error an occurred. Please try again",
		})
		return
	}
	disabledUser, _ := app.UserService.IsDisabledUser(credentials.Email)
	if disabledUser.Disabled.Bool == true {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Error - Credentials error. Contact help and service before try again.",
		})
		return
	}
	verifiedEmail, _ := app.UserService.VerifyUserByEmail(credentials.Email)
	if verifiedEmail.Verified.Bool == false {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Error - Invalid: An email verification required.",
		})
		return
	}
	enabledUser, _ := app.UserService.IsEnabledUser(credentials.Email)
	if enabledUser.Enabled.Bool == false {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Error - Account is not enabled.",
		})
		return
	}
	password := credentials.Password
	hashPassword := existsUser.Password
	if existsUser.FailedAttempts >= MAXIMUM && time.Since(existsUser.LastFailedAttemptTime.Time) < AWAIT*time.Minute {
		utils.WriteJSON(w, http.StatusForbidden, map[string]any{
			"error": "Error - Account is locked. Try again in 5 minutes.",
		})
		return
	}
	err := utils.ComparePassword(hashPassword, password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			if err := app.UserService.IncrementFailedLoginAttempts(existsUser); err != nil {
				log.Printf("Error incrementing failed attempts: %v", err)
			}
			utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
				"error": "Error - An invalid - unauthorized credentials.",
			})
			return
		}
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Error - An invalid an account credentials.",
		})
		return
	}
	if utils.HashPasswordIsExpired(existsUser.PasswordLastChangedAt) {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Error - Password has expired: Please update your password",
		})
		return
	}
	accessTokenTTL := 20 * time.Second // Define LL time at 15*time.Minute
	accessToken, accessClaims, err := app.TokenMaker.CreateToken(existsUser, accessTokenTTL)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Error - Creating access token claims",
		})
		return
	}
	refreshTokenTTL := 24 * time.Hour // Define LL time at 24*time.Hour
	refreshToken, refreshClaims, err := app.TokenMaker.CreateToken(existsUser, refreshTokenTTL)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Error - Creating access token claims",
		})
		return
	}
	err = app.UserService.ClearRefreshToken(existsUser.Email)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Error - Clear refresh token session",
		})
		return
	}
	newToken, err := app.UserService.StoreRefreshToken(&models.RefreshToken{
		ID:           refreshClaims.RegisteredClaims.ID,
		UUID:         existsUser.UUID,
		UserEmail:    existsUser.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    models.CustomTime(refreshClaims.RegisteredClaims.ExpiresAt.Time),
	})
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Error - Creating session store token",
		})
		return
	}
	response := &models.LoginCredentialsResponse{
		UUID:                  newToken.UUID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User:                  existsUser,
	}
	if err := app.UserService.ResetFailedLoginAttempts(existsUser); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Error - resetting attempts failed",
		})
	}
	iPRemote := r.RemoteAddr
	userAgent := r.Header.Get("User-Agent")
	browserAgent := middleware.DetectBrowser(userAgent) // not sure if it works
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Success - Your have log in successfully",
		"data":    response,
		"browser": browserAgent,
		"address": iPRemote,
	})
}
func (app *UserResource) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]any{
			"error": "Error - method not allowed.",
		})
		return
	}
	users, err := app.UserService.GetAllUsers()
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Error - internal server error.",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, map[string]map[string]interface{}{
		"data": {"users": users},
	})
}
func (app *UserResource) GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("uuid")
	if params == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that ID not avaliable"))
		return
	}
	userId, _ := app.UserService.GetUserById(params)
	if params != userId.UUID {
		utils.WriteError(w, http.StatusNotAcceptable, fmt.Errorf("[Error]: Field user payload doesn't exists"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, map[string]map[string]interface{}{
		"data": {"user": userId},
	})
}
func (app *UserResource) UpdateUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var payload models.UpdateUserRequest
	params := r.PathValue("uuid")
	if params == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that ID not avaliable"))
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "Error - Invalid body: An error occurred.",
		})
		return
	}
	_, err := utils.RegexValidate(payload.FirstName)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("[Error]: field invalid characters"))
		return
	}
	_, err = utils.RegexValidate(payload.LastName)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("[Error]: field invalid characters"))
		return
	}
	_, err = utils.RegexValidatePhone(payload.PhoneNumber)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("[Error]: That field an invalid phone numbers"))
		return
	}
	userById, err := app.UserService.GetUserById(params)
	if params != userById.UUID {
		utils.WriteError(w, http.StatusNotAcceptable, fmt.Errorf("[Error]: Field user payload doesn't exists"))
		return
	}
	updateUser := models.User{
		UUID:        params,
		FirstName:   payload.FirstName,
		LastName:    payload.LastName,
		PhoneNumber: payload.PhoneNumber,
	}
	err = app.UserService.UpdateUserById(&updateUser)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("[Error:] Internel server error"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": "Update user successfully.",
	})
}
func (app *UserResource) DeleteUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	params := r.PathValue("uuid")
	if params == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("[Error]: ID payload"))
		return
	}
	userById, err := app.UserService.GetUserById(params)
	if params != userById.UUID {
		utils.WriteError(w, http.StatusNotAcceptable, fmt.Errorf("[Error]: Field user payload doesn't exists"))
		return
	}
	err = app.UserService.DeleteUserById(params)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("[Error]: That role id not avaliable"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"success": "Delete user successfully.",
	})
}
func (app *UserResource) GetUserTokenEmailVerificationHandler(w http.ResponseWriter, r *http.Request) {
	UUID := r.URL.Query().Get("uuid")
	UToken := r.URL.Query().Get("token")
	if UUID == "" || UToken == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}
	// Still in development
	// You have disabled accounts. Please enable your accounts to login.",
	// "message": "Your accounts have been locked. Please contact your administrator to resolve the issues.",
	user, _ := app.UserService.GetUserTokenEmailVerification(UUID, UToken)
	if user.VerifyToken == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("[Error}: Your accounts already activated. Go ahead and log in"))
		return
	}
	err := app.UserService.UpdateUserEmailVerification(user)
	if err != nil {
		http.Error(w, "error verifying token", http.StatusUnauthorized)
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": "Your account has be activated.",
	})
}
func (app *UserResource) RenewAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	params := r.PathValue("token")
	refreshClaims, err := app.TokenMaker.VerifyToken(params)
	if err != nil {
		http.Error(w, "Invalid token can't be verified....", http.StatusUnauthorized)
		return
	}
	// Still in development => Tomorrow assignment
	// Here need to check if there exist user email refresh: Remove it before insert another new token.
	refreshToken, err := app.UserService.GetRefreshToken(&models.RefreshToken{
		ID:        refreshClaims.RegisteredClaims.ID,
		UUID:      refreshClaims.UUID,
		UserEmail: refreshClaims.Email,
	})
	if err != nil {
		http.Error(w, "error getting session", http.StatusInternalServerError)
		return
	}
	if refreshToken.IsRevoked {
		http.Error(w, "session revoked", http.StatusUnauthorized)
		return
	}
	if refreshToken.UserEmail != refreshClaims.Email {
		http.Error(w, "invalid... session", http.StatusUnauthorized)
		return
	}
	accessTokenTTL := 20 * time.Second // Define LL time at 15*time.Minute
	accessToken, accessClaims, err := app.TokenMaker.RenewAccessToken(refreshClaims.UUID, refreshClaims.Email, accessTokenTTL)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}
	response := models.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}
	fmt.Println("RefreshToken: ", response)
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": response,
	})
}
func (app *UserResource) LogOutUserHandler(w http.ResponseWriter, r *http.Request) {
	params := r.PathValue("token")
	refreshClaims, err := app.TokenMaker.VerifyToken(params)
	if err != nil {
		http.Error(w, "Invalid token can't be verified....", http.StatusUnauthorized)
		return
	}
	// Still in development => Tomorrow assignment
	// Here need to check if there exist user email refresh: Remove it before insert another new token.
	refreshToken, err := app.UserService.GetRefreshToken(&models.RefreshToken{
		ID:        refreshClaims.RegisteredClaims.ID,
		UUID:      refreshClaims.UUID,
		UserEmail: refreshClaims.Email,
	})
	if err != nil {
		http.Error(w, "error getting session", http.StatusInternalServerError)
		return
	}
	if refreshToken.IsRevoked {
		http.Error(w, "session revoked", http.StatusUnauthorized)
		return
	}
	if refreshToken.UserEmail != refreshClaims.Email {
		http.Error(w, "invalid... session", http.StatusUnauthorized)
		return
	}
	accessTokenTTL := 20 * time.Second // Define LL time at 15*time.Minute
	accessToken, accessClaims, err := app.TokenMaker.RenewAccessToken(refreshClaims.UUID, refreshClaims.Email, accessTokenTTL)
	if err != nil {
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}
	response := models.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"data": response,
	})
}
func (app *UserResource) RevokedAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	valueId := r.PathValue("id")
	if valueId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that session ID not avaliable"))
		return
	}
	// POST: http://localhost:8080/api/v1/revoke/bc8ea776-7426-4752-b776-2281f6e6e6b2
	err := app.UserService.RevokeRefreshToken(valueId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error revoking session"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (app *UserResource) BlackListRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	params := r.PathValue("token")
	refreshClaims, err := app.TokenMaker.VerifyToken(params)
	if err != nil {
		http.Error(w, "Invalid token can't be verified....", http.StatusUnauthorized)
		return
	}
	// Still in development => Tomorrow assignment
	// Here need to check if there exist user email refresh: Remove it before insert another new token.

	refreshToken, err := app.UserService.GetRefreshToken(&models.RefreshToken{
		ID:        refreshClaims.RegisteredClaims.ID,
		UUID:      refreshClaims.UUID,
		UserEmail: refreshClaims.Email,
	})
	if err != nil {
		http.Error(w, "error getting session", http.StatusInternalServerError)
		return
	}
	if refreshToken.IsRevoked {
		http.Error(w, "session revoked", http.StatusUnauthorized)
		return
	}
	if refreshToken.UserEmail != refreshClaims.Email {
		http.Error(w, "invalid...AAAAA session", http.StatusUnauthorized)
		return
	}
	err = app.UserService.DeleteRefreshToken(refreshClaims.UUID, refreshClaims.Email)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error revoking session"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (app *UserResource) ResetUserPasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, MethodNotAllowed)
		return
	}
	var user models.User
	var credentials models.ForgotPasswordRequest
	// Still in development
	// THIS DOESN'T WORK ON LOGIN, DOESN'T WORK ON REGISTER
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		utils.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "Error - Invalid: An error occurred. Please try again.",
		})
		return
	}
	if err := utils.Validate.Struct(&credentials); err != nil {
		_ = err.(validator.ValidationErrors)
		utils.WriteJSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "Error - Invalid. A valid email required.",
		})
		return
	}

	existsUser, _ := app.UserService.VerifyUserByEmail(credentials.Email)
	if existsUser.Email == user.Email {
		utils.WriteJSON(w, http.StatusNotAcceptable, map[string]any{
			"error": "Error - User with that email doesn't exists.",
		})
		return
	}
	emailToken := models.User{
		Email: credentials.Email,
	}
	// sent reset password token
	//err := app.UserService.ResetUserPasswordToken(existsUser)  // TRY THIS TO SEE IF IT WORKS
	err := app.UserService.ResetUserPasswordToken(&emailToken)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("server error. try it again"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "We sent you an email for reset password",
	})
}
func (app *UserResource) VerifyUserPasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	UToken := r.URL.Query().Get("token")
	if UToken == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request"))
		return
	}
	user, _ := app.UserService.VerifyUserPasswordToken(UToken)
	if user.VerifyToken == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error invalid token request"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Password link have been verified",
		"data":    user,
	})
}
func (app *UserResource) UpdateUserPasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, MethodNotAllowed)
		return
	}
	var user models.User
	var credentials models.ChangePasswordWithTokenRequest

	if err := utils.ParseJSON(r, &credentials); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}
	if err := utils.Validate.Struct(&credentials); err != nil {
		validationError := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error", validationError))
		return
	}
	existsUser, _ := app.UserService.VerifyUserPasswordToken(credentials.Token)
	if existsUser.VerifyToken == user.VerifyToken {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("that user or email doesn't exists"))
		return
	}
	updatePassword := models.User{
		VerifyToken: credentials.Token,
		Password:    credentials.Password,
	}
	// Still in development => TO DO NOTE->

	err := app.UserService.UpdateUserPasswordToken(&updatePassword)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("server error. Try it again"))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "New password has been successfully.",
	})
}
