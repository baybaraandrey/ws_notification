package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

func CheckDjangoPassword(password string, encodedPassword string) bool {
	passArr := strings.Split(encodedPassword, "$")
	//method := passArr[0]
	iter, _ := strconv.Atoi(passArr[1])
	salt := passArr[2]
	shapass := passArr[3]
	dk := pbkdf2.Key([]byte(password), []byte(salt), iter, 32, sha256.New)
	pass := base64.StdEncoding.EncodeToString(dk)
	return pass == shapass
}
