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
)

func Refresh(w http.ResponseWriter, r *http.Request) {
	//Check method
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("GET Method Not Allowed"))
		return
	}

	//Получаем Токены
	var token = map[string]string{}

	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshTokenString := token["refresh_token"]

	//Достаем токен и проверяем
	refToken, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(info.SECRET_KEY), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !refToken.Valid && refToken.Claims.Valid() == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Получаем инфу для Refrash токена
	claims, ok := refToken.Claims.(jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Вытаскиваем юзера по токену
	var user models.User
	hahsRef := base64.StdEncoding.EncodeToString([]byte(refreshTokenString))
	err = user.GetByRefrashToken(hahsRef)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Генерируем токен с инфой из предыдушего
	tokens, err := generateRefrashTokenPair(user, time.Now().Add(15*time.Minute).Unix(), claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	//Генерируем добавляем новый
	hahsRef = base64.StdEncoding.EncodeToString([]byte(tokens["refresh_token"]))
	err = user.UpdateAddNewToken(hahsRef)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//Отправляем юзеру
	res, err := json.Marshal(tokens)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

//generateRefrashTokenPair генератор токенов с инфой из прошлого
func generateRefrashTokenPair(user models.User, timeAccess int64, claims jwt.MapClaims) (map[string]string, error) {

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

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	refreshTokenString, err := refreshToken.SignedString([]byte(info.SECRET_KEY))

	if err != nil {
		return token, err
	}

	token["access_token"] = accessTokenString
	token["refresh_token"] = refreshTokenString

	return token, err
}
