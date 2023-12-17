package helper

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GenPassword() string {
	//random
	random := RandomString(10)

	return random
}

func RandomString(n int) string {
	//var pool = "0123456789ABCEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	specials := "@#"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 10
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	buf[2] = upper[rand.Intn(len(upper))]
	for i := 3; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf)
	return string(str)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
