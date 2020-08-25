package utility

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/aka-achu/watcher/logging"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strings"
	"time"
)

func Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func CreateToken(user string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = user
	atClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString([]byte(os.Getenv("ACCESS_SECRET")))
}

func VerifyToken(bearToken string) error {
	strArr := strings.Split(bearToken, " ")
	if len(strArr) != 2 {
		logging.Error.Printf(" [APP] Invalid auth token in the request header")
		return errors.New("invalid bearToken format")
	}
	if token, err := jwt.Parse(strArr[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	}); err != nil {
		logging.Error.Printf(" [APP] Failed to parse the JWT token. %v", err)
		return err
	} else {
		if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
			logging.Warn.Printf(" [APP] Unauthorized JWT token")
			return err
		} else {
			return nil
		}
	}
}
