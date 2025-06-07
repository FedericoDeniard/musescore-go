package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/FedericoDeniard/musescore-go/src/config" // Importar tu config
	customErrors "github.com/FedericoDeniard/musescore-go/src/utils/error"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CognitoJWTClaims representa las claims específicas de Cognito
type CognitoJWTClaims struct {
	Sub      string `json:"sub"`
	Email    string `json:"email,omitempty"`
	Username string `json:"cognito:username,omitempty"`
	TokenUse string `json:"token_use,omitempty"`
	Audience string `json:"aud,omitempty"`
	jwt.RegisteredClaims
}

// JWK representa una clave JSON Web Key
type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
	Alg string `json:"alg"`
}

// JWKSet representa un conjunto de claves JWKS
type JWKSet struct {
	Keys []JWK `json:"keys"`
}

// JWKSCache almacena las claves JWKS en caché
type JWKSCache struct {
	keys      *JWKSet
	lastFetch time.Time
	cacheTTL  time.Duration
}

var jwksCache *JWKSCache

func init() {
	jwksCache = &JWKSCache{
		cacheTTL: 24 * time.Hour, // Cache por 1 hora
	}
}

// base64URLDecode decodifica una cadena base64 URL-safe
func base64URLDecode(s string) ([]byte, error) {
	// Agregar padding si es necesario
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}

	// Reemplazar caracteres URL-safe
	s = strings.ReplaceAll(s, "-", "+")
	s = strings.ReplaceAll(s, "_", "/")

	return base64.StdEncoding.DecodeString(s)
}

// jwkToRSAPublicKey convierte una JWK a una clave pública RSA
func jwkToRSAPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decodificar N (modulus)
	nBytes, err := base64URLDecode(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("error decodificando N: %v", err)
	}

	// Decodificar E (exponente)
	eBytes, err := base64URLDecode(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("error decodificando E: %v", err)
	}

	// Convertir a big.Int
	n := new(big.Int).SetBytes(nBytes)

	// E es típicamente 65537 (0x010001)
	e := 0
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: n,
		E: e,
	}, nil
}

// getJWKS obtiene las claves JWKS de Cognito (con caché)
func getJWKS() (*JWKSet, error) {
	now := time.Now()

	// Si tenemos claves en caché y no han expirado, las devolvemos
	if jwksCache.keys != nil && now.Sub(jwksCache.lastFetch) < jwksCache.cacheTTL {
		return jwksCache.keys, nil
	}

	// Construir la URL del JWKS usando tu config
	region := config.KEYS.AWS_DEFAULT_REGION
	userPoolID := config.KEYS.AWS_USER_POOL_ID

	if region == "" || userPoolID == "" {
		return nil, fmt.Errorf("AWS_DEFAULT_REGION y AWS_USER_POOL_ID son requeridos")
	}

	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
		region, userPoolID)

	// Hacer la petición HTTP
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo JWKS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error HTTP %d obteniendo JWKS", resp.StatusCode)
	}

	// Decodificar la respuesta JSON
	var jwks JWKSet
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("error decodificando JWKS: %v", err)
	}

	// Actualizar caché
	jwksCache.keys = &jwks
	jwksCache.lastFetch = now

	return &jwks, nil
}

// getPublicKey obtiene la clave pública para verificar el JWT
func getPublicKey(token *jwt.Token) (interface{}, error) {
	// Verificar que el algoritmo sea RS256
	if token.Header["alg"] != "RS256" {
		return nil, fmt.Errorf("algoritmo no soportado: %v", token.Header["alg"])
	}

	// Obtener el kid (key ID) del header del token
	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("kid no encontrado en el header del token")
	}

	// Obtener las claves JWKS
	keySet, err := getJWKS()
	if err != nil {
		return nil, err
	}

	// Buscar la clave por kid
	var selectedKey *JWK
	for _, key := range keySet.Keys {
		if key.Kid == kid {
			selectedKey = &key
			break
		}
	}

	if selectedKey == nil {
		return nil, fmt.Errorf("clave con kid %s no encontrada", kid)
	}

	// Convertir JWK a clave RSA pública
	rsaKey, err := jwkToRSAPublicKey(*selectedKey)
	if err != nil {
		return nil, fmt.Errorf("error convirtiendo JWK a RSA: %v", err)
	}

	return rsaKey, nil
}

// ValidateJWT es el middleware de Gin para validar tokens JWT de Cognito
func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el header Authorization
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			customErrors.HandleError(c, &customErrors.HttpError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Tu sesión ha expirado. Por favor, inicia sesión nuevamente.",
			})
			return
		}

		// Extraer el token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			customErrors.HandleError(c, &customErrors.HttpError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Tu sesión ha expirado. Por favor, inicia sesión nuevamente.",
			})
			return
		}

		// Parsear y validar el token
		token, err := jwt.ParseWithClaims(tokenString, &CognitoJWTClaims{}, getPublicKey)
		if err != nil {
			customErrors.HandleError(c, &customErrors.HttpError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Tu sesión ha expirado. Por favor, inicia sesión nuevamente.",
			})
			return
		}

		if !token.Valid {
			customErrors.HandleError(c, &customErrors.HttpError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Tu sesión ha expirado. Por favor, inicia sesión nuevamente.",
			})
			return
		}

		// Obtener las claims
		claims, ok := token.Claims.(*CognitoJWTClaims)
		if !ok {
			customErrors.HandleError(c, &customErrors.HttpError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Tu sesión ha expirado. Por favor, inicia sesión nuevamente.",
			})
			return
		}

		// Validar que sea un token ID
		if claims.TokenUse != "id" {
			customErrors.HandleError(c, &customErrors.HttpError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Tu sesión ha expirado. Por favor, inicia sesión nuevamente.",
			})
			return
		}

		// Validar audience
		clientID := config.KEYS.AWS_USER_POOL_CLIENT_ID
		if clientID != "" && claims.Audience != clientID {
			customErrors.HandleError(c, &customErrors.HttpError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Tu sesión ha expirado. Por favor, inicia sesión nuevamente.",
			})
			return
		}

		// Validar issuer
		region := config.KEYS.AWS_DEFAULT_REGION
		userPoolID := config.KEYS.AWS_USER_POOL_ID
		expectedIssuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, userPoolID)
		if claims.Issuer != expectedIssuer {
			customErrors.HandleError(c, &customErrors.HttpError{
				StatusCode: http.StatusUnauthorized,
				Message:    "Tu sesión ha expirado. Por favor, inicia sesión nuevamente.",
			})
			return
		}

		// Agregar el usuario al contexto de Gin
		c.Set("user", claims)

		// Continuar con el siguiente handler
		c.Next()
	}
}

// GetUserFromContext obtiene el usuario del contexto de Gin
func GetUserFromContext(c *gin.Context) (*CognitoJWTClaims, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	claims, ok := user.(*CognitoJWTClaims)
	return claims, ok
}
