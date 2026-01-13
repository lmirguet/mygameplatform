package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type HS256Signer struct {
	secret []byte
	ttl    *time.Duration
}

func NewHS256Signer(secret []byte, ttl *time.Duration) (*HS256Signer, error) {
	if len(secret) == 0 {
		return nil, errors.New("jwt secret is required")
	}
	return &HS256Signer{secret: secret, ttl: ttl}, nil
}

func (s *HS256Signer) Sign(sub string, now time.Time) (string, error) {
	hdr := map[string]any{"alg": "HS256", "typ": "JWT"}
	payload := map[string]any{
		"sub": sub,
		"iat": now.Unix(),
	}
	if s.ttl != nil {
		payload["exp"] = now.Add(*s.ttl).Unix()
	}

	headerJSON, err := json.Marshal(hdr)
	if err != nil {
		return "", err
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	enc := base64.RawURLEncoding
	h := enc.EncodeToString(headerJSON)
	p := enc.EncodeToString(payloadJSON)

	unsigned := h + "." + p
	mac := hmac.New(sha256.New, s.secret)
	if _, err := mac.Write([]byte(unsigned)); err != nil {
		return "", err
	}
	sig := enc.EncodeToString(mac.Sum(nil))
	return unsigned + "." + sig, nil
}

func LooksLikeJWT(token string) bool {
	parts := strings.Split(token, ".")
	return len(parts) == 3 && parts[0] != "" && parts[1] != "" && parts[2] != ""
}

