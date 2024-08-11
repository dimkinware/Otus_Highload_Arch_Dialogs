package main

import (
	"HighArch-dialogs/api"
	"HighArch-dialogs/service"
	"HighArch-dialogs/storage"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	dialogService *service.DialogService
	authService   *service.AuthService
}

func NewServer(goCqlSession *gocql.Session) *Server {
	dialogStore := storage.NewDialogStore(goCqlSession)
	return &Server{
		dialogService: service.NewDialogService(dialogStore),
	}
}

func (s *Server) GetDialogListHandler(w http.ResponseWriter, req *http.Request) {
	currentUserId, err := getUserIdFromContext(req.Context())
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	peerId := mux.Vars(req)["user_id"]
	res, err := s.dialogService.GetDialog(currentUserId, peerId)
	if err != nil {
		respondError(w, err)
	} else {
		renderJSON(w, res)
	}
}

func (s *Server) GetSendMessageHandler(w http.ResponseWriter, req *http.Request) {
	currentUserId, err := getUserIdFromContext(req.Context())
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	var sendMessageModel api.DialogMessageSendApiModel
	err = parseJSON(req, &sendMessageModel)
	if err != nil { // validation error
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		peerId := mux.Vars(req)["user_id"]
		err := s.dialogService.AddDialogMessage(currentUserId, peerId, sendMessageModel.Text)
		if err != nil {
			respondError(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}
}

// Auth middleware methods

const userIdKey string = "user_id"

func (s *Server) GetAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		xRequestId := r.Header.Get(xRequestIdName)
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		userId, err := s.authService.Authenticate(tokenString, xRequestId)
		if err != nil || userId == nil || *userId == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		println("Checked auth for user: " + *userId)
		ctx := context.WithValue(r.Context(), userIdKey, *userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserIdFromContext(ctx context.Context) (string, error) {
	userId, ok := ctx.Value(userIdKey).(string)
	if !ok {
		return "", fmt.Errorf("user id not found in context")
	}
	return userId, nil
}

// X-Request-Id middleware methods

const xRequestIdName string = "X-Request-Id"

func (s *Server) GetXRequestIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestId := r.Header.Get(xRequestIdName)
		if xRequestId == "" {
			r.Header.Set(xRequestIdName, uuid.NewString())
		}
		next.ServeHTTP(w, r)
	})
}

// Utils methods

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js) // TODO: should handle error???
}

func parseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func respondError(w http.ResponseWriter, err error) {
	log.Println(err)
	if errors.Is(err, service.ErrorNotFound) {
		w.WriteHeader(http.StatusNotFound)
	} else if errors.Is(err, service.ErrorValidation) {
		w.WriteHeader(http.StatusBadRequest)
	} else if errors.Is(err, service.ErrorStoreError) {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
