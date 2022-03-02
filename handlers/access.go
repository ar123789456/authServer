package handler

import (
	"auth/models"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const SECRET_KEY = "ffasdfu324i53t2fo43j"

func Access(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("GET Method Not Allowed"))
		return
	}
	var userinfo models.UserInfo
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&userinfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = user.Get(userinfo.Name)
	if err != nil || CheckPasswordHash(userinfo.Password, user.Password) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	tokens, err := generateTokenPair(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hahsRef := base64.StdEncoding.EncodeToString([]byte(tokens["refresh_token"]))
	err = user.UpdateAddNewToken(hahsRef)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(tokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func generateTokenPair(user models.User) (map[string]string, error) {

	var token map[string]string
	token = map[string]string{}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":  time.Now().Add(15 * time.Minute).Unix(),
		"name": user.Name,
		"id":   user.ID,
	})

	accessTokenString, err := accessToken.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return token, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		//    "exp": time.Now().Add(time.Second * 24 ).Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString([]byte("secret"))

	if err != nil {
		return token, err
	}

	token["access_token"] = accessTokenString
	token["refresh_token"] = refreshTokenString

	return token, err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
