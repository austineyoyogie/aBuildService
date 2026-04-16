package routes

import (
	"net/http"
)

func RouteHandler() *http.ServeMux {
	mux := http.NewServeMux()
	//TokenMaker  *middleware.JWTMaker
	//jwtMaker = *middleware.JWTMaker



	//mux.Handle("/v3/", http.StripPrefix("/v3", mux))
	//mux.HandleFunc("GET /welcome", welcomeResource.WelcomeGetHandler)
	/*-------------------------------------------------------------*/
	//mux.HandleFunc("POST /login", userResource.LoginUserHandler)
	//mux.HandleFunc("POST /register", userResource.RegisterUserHandler)
	//mux.HandleFunc("GET /verify/account", userResource.GetUserTokenEmailVerificationHandler)

	// get to see if it works with for front end
	// mux.HandleFunc("POST /refresh/{refresh_token}", userResource.RenewAccessTokenSessionHandler)
	//mux.HandleFunc("POST /refresh", userResource.RenewAccessTokenSessionHandler)
	//
	//mux.HandleFunc("POST /revoke/{id}", userResource.RevokedAccessTokenSessionHandler)
	//mux.HandleFunc("POST /session/{id}", userResource.DeleteRefreshTokenSessionHandler)
	//
	//mux.HandleFunc("POST /reset/password", userResource.ResetUserPasswordTokenHandler)
	//mux.HandleFunc("GET /verify/password", userResource.VerifyUserPasswordTokenHandler)
	//mux.HandleFunc("POST /update/password", userResource.UpdateUserPasswordTokenHandler)

	mux.HandleFunc("POST /verify/send-sms", twilioResource.SendSMSOTPHandler)
	mux.HandleFunc("POST /verify/verify-sms", twilioResource.VerifySMSOTPHandler)
	mux.HandleFunc("POST /verify/create-totp", twilioResource.CreateTOTPFactorHandler)
	mux.HandleFunc("POST /verify/verify-factor", twilioResource.VerifyFactorHandler)
	mux.HandleFunc("POST /verify/create-totp-challenge", twilioResource.CreateTOTPChallengeHandler)
	/*-------------------------------------------------------------*/
	// error verifying token: ....invalid token: error parsing token: token signature is invalid: signature is invalid
	//mux.Handle("GET /profile", middleware.AuthMiddlewareFunc(jwtMaker)(http.HandlerFunc(userResource.GetAllUsersHandler)))
	//mux.Handle("GET /show/{uuid}", middleware.AuthMiddlewareFunc(jwtMaker)(http.HandlerFunc(userResource.GetUserByIdHandler)))
	//mux.HandleFunc("GET /profile", userResource.GetAllUsersHandler)
	//mux.HandleFunc("GET /show/{uuid}", userResource.GetUserByIdHandler)

	//mux.HandleFunc("GET /profile", userResource.FindAllUsersHandler)
	//mux.HandleFunc("PUT /update/profile/{uuid}", userResource.UpdateUserByIdHandler)
	//mux.HandleFunc("DELETE /delete/profile/{uuid}", userResource.DeleteUserByIdHandler)
	/*-------------------------------------------------------------*/
	//mux.HandleFunc("POST /role", roleResource.CreateRoleHandler)
	mux.HandleFunc("GET /role/{id}", roleResource.FindRoleIdHandler)
	mux.HandleFunc("GET /role", roleResource.FindAllRoleHandler)
	mux.HandleFunc("PUT /role/{id}", roleResource.UpdateRoleIdHandler)
	mux.HandleFunc("DELETE /role/{id}", roleResource.DeleteRoleIdHandler)
	/*-------------------------------------------------------------*/
	mux.HandleFunc("POST /supply", supplyResource.CreateSupplyHandler)
	mux.HandleFunc("GET /supply/{id}", supplyResource.FindSupplyIdHandler)
	mux.HandleFunc("GET /supplier", supplyResource.FindAllSupplyHandler)
	mux.HandleFunc("PUT /supply/{id}", supplyResource.UpdateSupplyIdHandler)
	mux.HandleFunc("DELETE /supply/{id}", supplyResource.DeleteSupplyIdHandler)
	/*-------------------------------------------------------------*/
	mux.HandleFunc("POST /category", categoryResource.CreateCategoryHandler)
	mux.HandleFunc("GET /category/{id}", categoryResource.FindCategoryIdHandler)
	mux.HandleFunc("GET /categories", categoryResource.FindAllCategoryHandler)
	mux.HandleFunc("PUT /category/{id}", categoryResource.UpdateCategoryIdHandler)
	mux.HandleFunc("DELETE /category/{id}", categoryResource.DeleteCategoryIdHandler)
	/*-------------------------------------------------------------*/
	mux.HandleFunc("POST /product", productResource.CreateProductHandler)
	mux.HandleFunc("GET /product/{id}", productResource.FindProductIdHandler)
	mux.HandleFunc("GET /products", productResource.FindAllProductHandler)
	mux.HandleFunc("PUT /product/{id}", productResource.UpdateProductIdHandler)
	mux.HandleFunc("DELETE /product/{id}", productResource.DeleteProductIdHandler)
	return mux
}
