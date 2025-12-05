package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// KeycloakJWTMiddleware validates JWT tokens from Keycloak
type KeycloakJWTMiddleware struct {
	baseURL        string
	realm          string
	clientID       string
	publicKeys     map[string]*rsa.PublicKey
	publicKeysLock sync.RWMutex
	httpClient     *http.Client
}

// NewJWTMiddleware creates a new Keycloak JWT middleware
func NewJWTMiddleware(baseURL, realm, clientID string) *KeycloakJWTMiddleware {
	return &KeycloakJWTMiddleware{
		baseURL:    baseURL,
		realm:      realm,
		clientID:   clientID,
		publicKeys: make(map[string]*rsa.PublicKey),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// JWKSResponse represents Keycloak JWKS response
type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// Claims represents JWT claims
type Claims struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	AuthorizedParty   string `json:"azp"`
	RealmAccess       struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
	jwt.RegisteredClaims
}

// getPublicKey retrieves the public key from Keycloak JWKS endpoint
func (m *KeycloakJWTMiddleware) getPublicKey(kid string) (*rsa.PublicKey, error) {
	m.publicKeysLock.RLock()
	if key, ok := m.publicKeys[kid]; ok {
		m.publicKeysLock.RUnlock()
		return key, nil
	}
	m.publicKeysLock.RUnlock()

	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", m.baseURL, m.realm)
	resp, err := m.httpClient.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status %d", resp.StatusCode)
	}

	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	// Parse and cache all keys
	m.publicKeysLock.Lock()
	defer m.publicKeysLock.Unlock()

	for _, jwk := range jwks.Keys {
		if jwk.Kty != "RSA" {
			continue
		}

		nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
		if err != nil {
			continue
		}

		eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
		if err != nil {
			continue
		}

		var eInt int
		for _, b := range eBytes {
			eInt = eInt<<8 | int(b)
		}

		publicKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: eInt,
		}

		m.publicKeys[jwk.Kid] = publicKey
	}

	if key, ok := m.publicKeys[kid]; ok {
		return key, nil
	}

	return nil, fmt.Errorf("key with kid %s not found", kid)
}

// Validate validates JWT token and returns claims
func (m *KeycloakJWTMiddleware) Validate(tokenString string) (*Claims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &Claims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("kid not found in token header")
	}

	publicKey, err := m.getPublicKey(kid)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	claims := &Claims{}
	token, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	expectedIssuer := fmt.Sprintf("%s/realms/%s", m.baseURL, m.realm)
	if claims.Issuer != expectedIssuer {
		return nil, fmt.Errorf("invalid issuer: expected %s, got %s", expectedIssuer, claims.Issuer)
	}

	audValid := false
	if len(claims.Audience) > 0 {
		for _, aud := range claims.Audience {
			if aud == m.clientID || aud == m.realm {
				audValid = true
				break
			}
		}
	}
	if !audValid && claims.AuthorizedParty != "" {
		if claims.AuthorizedParty == m.clientID {
			audValid = true
		}
	}

	if !audValid && len(claims.Audience) == 0 {
		audValid = true
	}
	if !audValid {
		return nil, fmt.Errorf("invalid audience")
	}

	return claims, nil
}

// RequireAuth creates a Fiber middleware that requires authentication
func (m *KeycloakJWTMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := parts[1]

		claims, err := m.Validate(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token",
				"details": err.Error(),
			})
		}

		c.Locals("user", claims)
		c.Locals("userID", claims.Sub)
		c.Locals("email", claims.Email)
		c.Locals("username", claims.PreferredUsername)

		return c.Next()
	}
}

// RequireRole creates a Fiber middleware that requires a specific role
func (m *KeycloakJWTMiddleware) RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := m.RequireAuth()(c); err != nil {
			return err
		}

		claims, ok := c.Locals("user").(*Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User claims not found",
			})
		}

		for _, r := range claims.RealmAccess.Roles {
			if r == role {
				return c.Next()
			}
		}

		for _, resource := range claims.ResourceAccess {
			for _, r := range resource.Roles {
				if r == role {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": fmt.Sprintf("Role '%s' is required", role),
		})
	}
}
