package resource

import (
	"aIBuildService/aPI/service"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type TwilioResource struct {
	db            *gorm.DB
	TwilioService service.TwilioService
}

type QRResponse struct {
	QRCode    string `json:"qrCode"`
	FactorSid string `json:"factorSid"`
}

func NewTwilioResource(twilioService service.TwilioService) *TwilioResource {
	return &TwilioResource{
		TwilioService: twilioService,
	}
}

func (tw *TwilioResource) SendSMSOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := tw.TwilioService.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err := tw.TwilioService.SendSMSOTP(user.PhoneNumber); err != nil {
		log.Printf("Failed to send SMS OTP: %v", err)
		http.Error(w, "Failed to send SMS OTP", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (tw *TwilioResource) VerifySMSOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := tw.TwilioService.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	verified, err := tw.TwilioService.VerifySMSOTP(user.PhoneNumber, req.Code)
	if err != nil {
		log.Printf("Failed to verify SMS OTP: %v", err)
		http.Error(w, "Failed to verify SMS OTP", http.StatusInternalServerError)
		return
	}
	if !verified {
		http.Error(w, "Invalid OTP", http.StatusUnauthorized)
		return
	}
	user.SMSEnabled = true
	if err := tw.TwilioService.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (tw *TwilioResource) CreateTOTPFactorHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := tw.TwilioService.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Check Id Variable again
	sid, uri, err := tw.TwilioService.CreateTOTPFactor(strconv.FormatInt(user.ID, 10), user.FirstName)
	if err != nil {
		log.Printf("Failed to create TOTP factor: %v", err)
		http.Error(w, "Failed to create TOTP factor", http.StatusInternalServerError)
		return
	}
	user.TOTPFactorSid = sid
	if err := tw.TwilioService.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	response := QRResponse{
		QRCode:    uri,
		FactorSid: sid,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func (tw *TwilioResource) VerifyFactorHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := tw.TwilioService.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Check Id Variable again
	verified, err := tw.TwilioService.VerifyFactor(user.TOTPFactorSid, req.Code, strconv.FormatInt(user.ID, 10))
	if err != nil {
		log.Printf("Failed to verify factor: %v", err)
		http.Error(w, "Failed to verify factor", http.StatusInternalServerError)
		return
	}
	if !verified {
		http.Error(w, "Invalid code", http.StatusUnauthorized)
		return
	}
	user.TOTPEnabled = true
	user.IsAuthenticated = true
	if err := tw.TwilioService.UpdateUser(user); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func (tw *TwilioResource) CreateTOTPChallengeHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	user, err := tw.TwilioService.GetUserByEmail(req.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	challengeSid, err := tw.TwilioService.CreateTOTPChallenge(user.TOTPFactorSid, req.Code, strconv.FormatInt(user.ID, 10))
	if err != nil {
		log.Printf("Failed to create TOTP challenge: %v", err)
		http.Error(w, "Failed to create TOTP challenge", http.StatusInternalServerError)
		return
	}
	user.IsAuthenticated = true
	if err := tw.TwilioService.UpdateUser(user); err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"challengeSid": challengeSid})
}
