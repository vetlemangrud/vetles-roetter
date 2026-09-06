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

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := os.Getenv("CARROT_WRITE_KEY")
		res := subtle.ConstantTimeCompare([]byte(key), []byte(r.Header.Get("Authorization")))

		switch res {
		case 1:
			next(w, r)
		case 0:
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
	}
}

func (h CarrotHandler) VetleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := h.VetleTemplate.Execute(w, nil); err != nil {
		log.Printf("template: %v", err)
	}
}

func (h CarrotHandler) VetlePost(w http.ResponseWriter, r *http.Request) {

	key := os.Getenv("CARROT_WRITE_KEY")
	res := subtle.ConstantTimeCompare([]byte(key), []byte(r.PostFormValue("key")))

	switch res {
	case 1:
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
	case 0:
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

	data := struct{ CarrotAmount int }{CarrotAmount: count}
	if err := h.HomeTemplate.Execute(w, data); err != nil {
		log.Printf("template: %v", err)
	}
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
