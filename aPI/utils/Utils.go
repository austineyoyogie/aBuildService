package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"html"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var Validate = validator.New()

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

/*====== NOT IN USE ======*/

func WriteAsJson(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
func WriteJSONAny(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// NOT WORKING
func ValidatePONE(password string) ([]byte, error) {
	return validatePhoneNumber(password)
}

func ValidatePH(phoneNumbers string) ([]byte, error) {
	for _, phoneNumber := range phoneNumbers {
		err, _ := validatePhoneNumber(string(phoneNumber))
		if err != nil {
			return []byte("All phone numbers must be in E.164 format"), nil
		}
	}
	return nil, nil
}

func ValidatePhoneNumbers(phoneNumbers string) error {
	for _, phoneNumber := range phoneNumbers {
		err, _ := validatePhoneNumber(string(phoneNumber))
		if err != nil {
			return errors.New("All phone numbers must be in E.164 format")
		}
	}
	return nil
}
func validatePhoneNumber(phoneNumber string) ([]byte, error) {
	e164Pattern := `^\+[1-9]\d{1,14}$`
	match, err := regexp.Match(e164Pattern, []byte(phoneNumber))
	if err != nil {
		log.Fatal(err.Error())
	}
	if !match {
		return []byte("Phone number must be in E.164 format"), nil

	}
	return nil, nil
}



// not IsMatch working. try again later

func IsMatch(value string) (string, error) {
	//US := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	//InTel := regexp.MustCompile(`^\+?[1-9]\d{0,3}[-.\s]?\d{1,14}$`)
	US := regexp.MustCompile(`^\(?(\d{3})\)?[-.\s]?(\d{3})[-.\s]?(\d{4})$`)
	if !US.MatchString(value) {
		return "", errors.New("numeric format must be valid")
	}
	return value, nil
}

func Escape(value string) string {
	return html.EscapeString(strings.TrimSpace(IsTitle(value)))
}
func IsTitle(value string) string {
	return cases.Title(language.English, cases.Compact).String(value)
}
func IsToLower(value string) string {
	return cases.Lower(language.English, cases.Compact).String(value)
}
func RegexValidate(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) < 2 || len(value) > 45 {
		return "", errors.New("field must be between 2 and 45 characters")
	}
	re := regexp.MustCompile(`^[a-zA-Z]+$`) // match without space
	if !re.MatchString(value) {
		return "", errors.New("character must be valid")
	}
	return value, nil
}
func RegexValidatePhone(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) < 2 || len(value) > 45 {
		return "", errors.New("field must be between 2 and 45 characters")
	}
	//re := regexp.MustCompile(`^[0-9-]+$`) // match without space
	//re := regexp.MustCompile(`^\d{3}\d{3}\d{4}$`) // match without space
	re := regexp.MustCompile(`^\+[1-9][1-9][0-9]{1,14}$`) // match without space
	if !re.MatchString(value) {
		return "", errors.New("character must be valid")
	}
	return value, nil
}

func ValidateSpaceRegex(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) < 2 || len(value) > 45 {
		return "", errors.New("field must be between 1 and 45 characters")
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9 ]+$`) // match with space and number
	if !re.MatchString(value) {
		return "", errors.New("character must be valid")
	}
	return value, nil
}

func ValidateDescriptionRegex(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) < 2 || len(value) > 255 {
		return "", errors.New("field must be between 3 and 45 characters")
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9- ]+$`) // match with space and number
	if !re.MatchString(value) {
		return "", errors.New("character must be valid")
	}
	return value, nil
}

func main111() {
	// Example usage
	email := "test@example.com"
	phone := "123-456-7890"
	name := "John Doe"

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	phoneRegex := regexp.MustCompile(`^\d{3}-\d{3}-\d{4}$`)
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)

	fmt.Printf("Email '%s' is valid: %t\n", email, emailRegex.MatchString(email))
	fmt.Printf("Phone '%s' is valid: %t\n", phone, phoneRegex.MatchString(phone))
	fmt.Printf("Name '%s' is valid: %t\n", name, nameRegex.MatchString(name))
}

func UserLoginErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("This %s field is required###.", IsToLower(err.Field()))
	default:
		return err.Error()
	}
}

// UserRegisterErrorMsg check profile model errors
func UserRegisterErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("This %s field is required.", IsToLower(err.Field()))
	case "min":
		return fmt.Sprintf("This %s field is less than required.", IsToLower(err.Field()))
	case "max":
		return fmt.Sprintf("This %s field is greater than required.", IsToLower(err.Field()))
	default:
		return err.Error()
	}
}

// UserProfileErrorMsg ProfileErrorMsg check profile model errors
func UserProfileErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("This %s field is required.", IsToLower(err.Field()))
	case "min":
		return fmt.Sprintf("This %s field is less than required.", IsToLower(err.Field()))
	case "max":
		return fmt.Sprintf("This %s field is greater than required.", IsToLower(err.Field()))
	case "gte":
		return fmt.Sprintf("This %s field is less than required.", IsToLower(err.Field()))
	case "lte":
		return fmt.Sprintf("This %s field is greater than required.", IsToLower(err.Field()))
	default:
		return err.Error()
	}
}
