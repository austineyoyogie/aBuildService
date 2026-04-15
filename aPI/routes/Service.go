package routes

import (
	"aIBuildService/aPI/config"
	"aIBuildService/aPI/config/database"
	"aIBuildService/aPI/resource"
	"aIBuildService/aPI/service/implementation"
	"os"
)

var (
	CD = database.CDriver()
	LC = config.LoadConfig()

	welcomeServiceImpl = implementation.NewWelcomeServiceImpl(CD)
	welcomeResource    = resource.NewWelcomeResource(welcomeServiceImpl)
	/*----------------------------------------------------*/
	/*----------------------------------------------------*/
	roleServiceImpl = implementation.NewRoleServiceImpl(CD)
	roleResource    = resource.NewRoleResource(roleServiceImpl)
	/*----------------------------------------------------*/
	twilioServiceImpl = implementation.NewTwilioServiceImpl(
		CD,
		os.Getenv("TWILIO_ACCOUNT_SID"),
		os.Getenv("TWILIO_AUTH_TOKEN"),
		os.Getenv("TWILIO_VERIFY_SID"))
	twilioResource = resource.NewTwilioResource(twilioServiceImpl)
	/*----------------------------------------------------*/
	supplyServiceImpl = implementation.NewSupplyServiceImpl(CD)
	supplyResource    = resource.NewSupplyResource(supplyServiceImpl)
	/*----------------------------------------------------*/
	categoryServiceImpl = implementation.NewCategoryServiceImpl(CD)
	categoryResource    = resource.NewCategoryResource(categoryServiceImpl)
	/*----------------------------------------------------*/
	productServiceImpl = implementation.NewProductServiceImpl(CD)
	productResource    = resource.NewProductResource(productServiceImpl)
)
