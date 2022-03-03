package handler

import (
	"auth/info"
	"auth/models"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func Access(w http.ResponseWriter, r *http.Request) {
	//Проверка http метода
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("GET Method Not Allowed"))
		return
	}
	//Достаем отправленные данные
	var userinfo models.UserInfo
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&userinfo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//Ищем юзера и аудентифицируем
	user.Name = userinfo.Name
	err = user.Get()
	if err != nil || CheckPasswordHash(userinfo.Password, user.Password) {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Генерируем пару токенов
	tokens, err := generateTokenPair(user, time.Now().Add(15*time.Minute).Unix(), time.Now().Add(30*24*time.Hour).Unix())
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//Сохраняем зашиврованный Refresh токен
	hahsRef := base64.StdEncoding.EncodeToString([]byte(tokens["refresh_token"]))
	err = user.UpdateAddNewToken(hahsRef)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Маршалим и отправляем юзеру
	res, err := json.Marshal(tokens)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

// generateTokenPair Генериратор пары токенов
func generateTokenPair(user models.User, timeAccess, timeRefrash int64) (map[string]string, error) {

	token := map[string]string{}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":  timeAccess,
		"name": user.Name,
		"id":   user.ID,
	})

	accessTokenString, err := accessToken.SignedString([]byte(info.SECRET_KEY))
	if err != nil {
		return token, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":  user.ID,
		"exp": timeRefrash,
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(info.SECRET_KEY))

	if err != nil {
		return token, err
	}

	token["access_token"] = accessTokenString
	token["refresh_token"] = refreshTokenString

	return token, err
}

//HashPassword шифратор паролей
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash сверка паролей
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
