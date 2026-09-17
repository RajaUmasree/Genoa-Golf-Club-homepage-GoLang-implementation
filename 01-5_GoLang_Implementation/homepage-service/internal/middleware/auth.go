package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/genoagolfclub/homepage-service/internal/models"
)

// jwtClaims is the minimal claim set this service cares about.
type jwtClaims struct {
	Exp int64 `json:"exp"`
}

type authError string

func (e authError) Error() string { return string(e) }

const (
	errInvalidToken = authError("invalid token")
	errExpiredToken = authError("expired token")
)

// verifyHS256JWT validates a compact JWT's structure, HS256 signature, and
// expiry using only the standard library — no external JWT dependency is
// required. This mirrors what a library such as golang-jwt/jwt would do for
// the subset of functionality this service needs, and can be swapped for
// that library directly (drop-in) if the project standardizes on it.
func verifyHS256JWT(token, secret string) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errInvalidToken
	}

	headerB64, payloadB64, sigB64 := parts[0], parts[1], parts[2]

	headerJSON, err := base64.RawURLEncoding.DecodeString(headerB64)
	if err != nil {
		return errInvalidToken
	}
	var header struct {
		Alg string `json:"alg"`
	}
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return errInvalidToken
	}
	if header.Alg != "HS256" {
		return errInvalidToken
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(headerB64 + "." + payloadB64))
	expectedSig := mac.Sum(nil)

	actualSig, err := base64.RawURLEncoding.DecodeString(sigB64)
	if err != nil {
		return errInvalidToken
	}
	if !hmac.Equal(expectedSig, actualSig) {
		return errInvalidToken
	}

	payloadJSON, err := base64.RawURLEncoding.DecodeString(payloadB64)
	if err != nil {
		return errInvalidToken
	}
	var claims jwtClaims
	if err := json.Unmarshal(payloadJSON, &claims); err != nil {
		return errInvalidToken
	}
	if claims.Exp != 0 && time.Now().Unix() > claims.Exp {
		return errExpiredToken
	}

	return nil
}

// JWTAuth returns an http middleware that validates a Bearer token against
// secret using HS256, matching the existing Genoa Golf Club backend's auth
// service. On failure it responds 401 with the standard error envelope.
func JWTAuth(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, models.NewError(models.CodeUnauthorized, "missing bearer token"))
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		if err := verifyHS256JWT(tokenStr, secret); err != nil {
			writeJSON(w, http.StatusUnauthorized, models.NewError(models.CodeUnauthorized, "invalid or expired token"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
