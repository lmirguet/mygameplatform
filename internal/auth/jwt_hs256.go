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

var ErrInvalidToken = errors.New("invalid_token")
var ErrTokenExpired = errors.New("token_expired")

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

func (s *HS256Signer) Verify(token string, now time.Time) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", ErrInvalidToken
	}
	if parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return "", ErrInvalidToken
	}

	enc := base64.RawURLEncoding
	headerJSON, err := enc.DecodeString(parts[0])
	if err != nil {
		return "", ErrInvalidToken
	}
	payloadJSON, err := enc.DecodeString(parts[1])
	if err != nil {
		return "", ErrInvalidToken
	}

	var header map[string]any
	if err := json.Unmarshal(headerJSON, &header); err != nil {
		return "", ErrInvalidToken
	}
	if alg, _ := header["alg"].(string); alg != "HS256" {
		return "", ErrInvalidToken
	}

	unsigned := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, s.secret)
	if _, err := mac.Write([]byte(unsigned)); err != nil {
		return "", ErrInvalidToken
	}
	expected := mac.Sum(nil)
	got, err := enc.DecodeString(parts[2])
	if err != nil {
		return "", ErrInvalidToken
	}
	if !hmac.Equal(expected, got) {
		return "", ErrInvalidToken
	}

	var payload map[string]any
	if err := json.Unmarshal(payloadJSON, &payload); err != nil {
		return "", ErrInvalidToken
	}
	sub, _ := payload["sub"].(string)
	if sub == "" {
		return "", ErrInvalidToken
	}
	if expRaw, ok := payload["exp"]; ok {
		switch v := expRaw.(type) {
		case float64:
			if now.Unix() > int64(v) {
				return "", ErrTokenExpired
			}
		case json.Number:
			exp, err := v.Int64()
			if err != nil {
				return "", ErrInvalidToken
			}
			if now.Unix() > exp {
				return "", ErrTokenExpired
			}
		default:
			return "", ErrInvalidToken
		}
	}
	return sub, nil
}

func LooksLikeJWT(token string) bool {
	parts := strings.Split(token, ".")
	return len(parts) == 3 && parts[0] != "" && parts[1] != "" && parts[2] != ""
}
