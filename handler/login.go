package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type userForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var user userForm
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "json parse error: "+err.Error(), http.StatusBadRequest)
		return
	}

	userId, err := h.authenticator.Login(user.Login, user.Password)
	if err != nil {
		http.Error(w, "wrong login/password", http.StatusForbidden)
		return
	}

	session, err := h.authenticator.GetSession(*userId)
	if err != nil {
		http.Error(w, "can't gen a session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "user_id", Value: strconv.Itoa(*userId), Path: "/", Domain: "localhost", Expires: time.Now().AddDate(0, 0, 7)})
	http.SetCookie(w, &http.Cookie{Name: "token", Value: session, Path: "/", Domain: "localhost", Expires: time.Now().AddDate(0, 0, 7)})

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LoginForm(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserId(r)
	tp, err := template.New("").ParseFiles("templates/base.html", "templates/login.html")
	if err != nil {
		log.Fatal("template rendering error:", err)
	}

	err = tp.ExecuteTemplate(w, "base", baseTemplate{"Вход", user != nil, struct{}{}})
	if err != nil {
		http.Error(w, "template rendering error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user userForm
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "json parse error: "+err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.authenticator.NewUser(user.Login, user.Password)
	if err != nil {
		http.Error(w, "user signup error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SignUpForm(w http.ResponseWriter, r *http.Request) {
	user := h.GetUserId(r)

	tp, err := template.New("").ParseFiles("templates/base.html", "templates/signup.html")
	if err != nil {
		log.Fatal("template rendering error:", err)
	}

	err = tp.ExecuteTemplate(w, "base", baseTemplate{"Регистрация", user != nil, struct{}{}})
	if err != nil {
		http.Error(w, "template rendering error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "user_id", Value: "0", Path: "/", Domain: "localhost", Expires: time.Now().Add(time.Millisecond * 300)})
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "del", Path: "/", Domain: "localhost", Expires: time.Now().Add(time.Millisecond * 300)})

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("<script>setTimeout(() => document.location.href = '/', 500)</script>"))
}
