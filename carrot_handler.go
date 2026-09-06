package main

import (
	"crypto/subtle"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

type CarrotHandler struct {
	Repository    CarrotRepository
	HomeTemplate  *template.Template
	VetleTemplate *template.Template
}

func keyIsCorrect(key string) bool{
	correct := os.Getenv("CARROT_WRITE_KEY")
	res := subtle.ConstantTimeCompare([]byte(correct), []byte(key))
	return res == 1
}


func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if IsVetle(r) {
			next(w, r)
		} else {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func IsVetle(r *http.Request) bool {
	// Check Auth header
	if keyIsCorrect(r.Header.Get("Authorization")) { return true }

	// Check vetle_key cookie
	key, err := r.Cookie("vetle_key")
	if(err == nil && keyIsCorrect(key.Value)) { return true }

	return false
}

func (h CarrotHandler) VetleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.VetleTemplate.Execute(w, nil); err != nil {
		log.Printf("template: %v", err)
	}
}

func (h CarrotHandler) VetlePost(w http.ResponseWriter, r *http.Request) {
	key := r.PostFormValue("key")
	if keyIsCorrect(key) {
		http.SetCookie(w, &http.Cookie{
			Name:     "vetle_key",
			Value:    key,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   7 * 24 * 3600,
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

}

func (h CarrotHandler) HomeGet(w http.ResponseWriter, r *http.Request) {
	carrots, err := h.Repository.findCarrots()
	if err != nil {
		http.Error(w, "Failed to get carrots from DB :(", http.StatusInternalServerError)
		return
	}
	count := len(carrots)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	data := struct{ CarrotAmount int; IsVetle bool }{CarrotAmount: count, IsVetle: IsVetle(r)}
	if err := h.HomeTemplate.Execute(w, data); err != nil {
		log.Printf("template: %v", err)
	}
}

func (h CarrotHandler) HomePost(w http.ResponseWriter, r *http.Request) {
	if _, err := h.Repository.addCarrot(); err != nil {
		http.Error(w, "Failed to add carrot :(", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h CarrotHandler) APIGet(w http.ResponseWriter, r *http.Request) {
	carrots, err := h.Repository.findCarrots()
	if err != nil {
		http.Error(w, "Failed to get carrots from DB :(", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(carrots); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h CarrotHandler) APIPost(w http.ResponseWriter, r *http.Request) {
	carrot, err := h.Repository.addCarrot()
	if err != nil {
		http.Error(w, "Failed to post carrot :(", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(carrot); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
