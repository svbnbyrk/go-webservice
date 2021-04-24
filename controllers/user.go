package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/svbnbyrk/webservice/models"
)

type userController struct {
	userIDPattern *regexp.Regexp
}

func (uc userController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/users" {
		switch r.Method {
		case http.MethodGet:
			uc.getAll(w, r)
		case http.MethodPost:
			uc.post(w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	} else {
		matches := uc.userIDPattern.FindStringSubmatch(r.URL.Path)
		if len(matches) == 0 {
			w.WriteHeader(http.StatusNotFound)
		}

		id, err := strconv.Atoi(matches[1])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
		switch r.Method {
		case http.MethodGet:
			uc.getAll(w, r)
		case http.MethodPost:
			uc.post(w, r)
		case http.MethodDelete:
			uc.delete(id, w)
		case http.MethodPut:
			uc.put(id, w, r)
		default:
			w.WriteHeader(http.StatusNotImplemented)
		}
	}
}

func (uc *userController) parseRequest(r *http.Request) (models.User, error) {
	dec := json.NewDecoder(r.Body)
	var u models.User
	err := dec.Decode(&u)
	if err != nil {
		return models.User{}, err
	}

	return u, nil
}

func (uc *userController) getAll(w http.ResponseWriter, r *http.Request) {
	encodeResponseAsJSON(models.GetUsers(), w)
	w.WriteHeader(http.StatusOK)
}

func (uc *userController) get(id int, w http.ResponseWriter) {
	u, err := models.GetUserById(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	encodeResponseAsJSON(u, w)
	w.WriteHeader(http.StatusOK)
}

func (uc *userController) post(w http.ResponseWriter, r *http.Request) {
	v, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Coult not parse user object"))
	}

	u, err := models.AddUser(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	encodeResponseAsJSON(u, w)
	w.WriteHeader(http.StatusOK)
}

func (uc *userController) put(id int, w http.ResponseWriter, r *http.Request) {
	v, err := uc.parseRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not parse User object"))
		return
	}

	if id != v.ID {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ID of submitted user must match ID in URL"))
	}

	u, err := models.UpdateUser(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	encodeResponseAsJSON(u, w)
	w.WriteHeader(http.StatusOK)
}

func (uc *userController) delete(id int, w http.ResponseWriter) {
	err := models.RemoveUserByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.WriteHeader(http.StatusOK)
}

// yapısal(struct) veri tipleriyle çalışırken var üretmeden & adresini kullanabilirim, örneğin 21 gibi int veri tipi için adresde spesifik bir yer ayrılmazken struck için ayrılır
func newUserController() *userController {
	return &userController{
		userIDPattern: regexp.MustCompile(`^/users/(\d+)/?`),
	}
}
