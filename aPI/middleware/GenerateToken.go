package middleware

import (
	"aIBuildService/aPI/models"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey}
}

func (maker *JWTMaker) CreateToken(user *models.User, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserTokenClaims(user, duration)
	if err != nil {
		return "", nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("error signing token: %w", err)
	}
	return tokenStr, claims, nil
}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

func (maker *JWTMaker) RenewAccessToken(UUID string, email string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := RenewAccessTokenClaims(UUID, email, duration)
	if err != nil {
		return "", nil, err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("error signing token: %w", err)
	}
	return tokenStr, claims, nil
}



var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
	ErrEmailInUse         = errors.New("email already in use")
)

var cfs = config.LoadConfig()

type JWTClaim struct {
	UserId    string `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	jwt.StandardClaims
}

func GenerateAccessToken(user *models.User) (tokenStr string, err error) {
	claims := &JWTClaim{
		UserId: strconv.Itoa(int(user.ID)),
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			Id:        "https://www.lts.co.uk",
			Issuer:    "https://www.lts.co.uk",
			Subject:   user.Email,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 1).Unix(), // 1 minutes // time.Hour * 24 = 1 day
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString(cfs.JWT.SecretTokenKey)
	return
}

func GenerateRefreshToken(user *models.User) (tokenStr string, err error) {
	claims := &JWTClaim{
		UserId: strconv.Itoa(int(user.ID)),
		Email:  user.Email,
		StandardClaims: jwt.StandardClaims{
			Id:        "https://www.lts.co.uk",
			Issuer:    "https://www.lts.co.uk",
			Subject:   user.Email,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString(cfs.JWT.SecretTokenKey)
	return
}

func GenerateNewRefreshToken(refresh *models.RefreshToken) (tokenStr string, err error) {
	claims := &JWTClaim{
		UserId: strconv.Itoa(int(refresh.UserId)),
		Email:  refresh.Email,
		StandardClaims: jwt.StandardClaims{

			Issuer:    "https://www.lts.co.uk",
			Subject:   refresh.Email,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Minute * 2).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString(cfs.JWT.SecretTokenKey)
	return
}

func VerifyToken(c *gin.Context) error {
	tokenStr := ExtractToken(c)
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return cfs.JWT.SecretTokenKey, nil
		})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return nil
	}
	if claims.ExpiresAt > time.Now().Add(time.Minute*1).Unix() {
		err = errors.New("token has expired")
		return nil
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	useBearer := c.Request.Header.Get("Authorization")
	if useBearer == "" {
		c.JSON(401, gin.H{"error": "required an access token!!"})
		c.Abort()
		return ""
	}
	if len(strings.Split(useBearer, " ")) == 2 {
		return strings.Split(useBearer, " ")[1]
	}
	return ""
}
*/

// https://neon.com/guides/golang-jwt  => insert jwt into database
// https://dev.to/neelp03/securing-your-go-api-with-jwt-authentication-4amj
// https://www.codingexplorations.com/blog/creating-a-restful-api-with-jwt-authentication-in-go
// https://www.youtube.com/watch?v=Xqk-5lynrCQ

/*
// Key type for context values
type contextKey string

const (
	// UserIDKey is the key for user ID in the request context
	UserIDKey contextKey = "userID"
)

// AuthMiddleware checks JWT tokens and adds user info to the request context
func AuthMiddleware(authService *service.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check Bearer token format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validate the token
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Extract user ID from claims
			userIDStr, ok := claims["sub"].(string)
			if !ok {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userId, err := uuid.Parse(userIDStr)
			if err != nil {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			// Add user ID to request context
			ctx := context.WithValue(r.Context(), UserIDKey, userId)

			// Call the next handler with the enhanced context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ValidateToken verifies a JWT token and returns the claims
func ValidateToken(tokenString string) (jwt.MapClaims, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// Extract and validate claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

