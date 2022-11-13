package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"balanceService/pkg/repository"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Server struct {
	server *http.Server
}

func (s *Server) Start() error {
	dbRepo := repository.DBRepository{}
	err := dbRepo.SetupDB()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer dbRepo.Shutdown()

	handler := NewHandler(dbRepo)
	router := mux.NewRouter()
	router.HandleFunc("/balance/{id:[0-9]+}", handler.GetUserBalance).Methods("GET", "OPTIONS")
	router.HandleFunc("/credit", handler.Credit).Methods("POST", "OPTIONS")
	router.HandleFunc("/reserve", handler.Reserve).Methods("POST", "OPTIONS")
	router.HandleFunc("/withdraw", handler.Withdraw).Methods("POST", "OPTIONS")
	router.HandleFunc("/refund", handler.Refund).Methods("POST", "OPTIONS")
	router.HandleFunc("/report/{year:[0-9]+}/{month:[0-9]+}", handler.GetReport).Methods("GET", "OPTIONS")
	router.HandleFunc("/history/{id:[0-9]+}/{page:[0-9]+}/{sort:|value|date}/{order:|asc|desc}", handler.GetTransactionHistory).Methods("GET", "OPTIONS")
	router.HandleFunc("/history/{id:[0-9]+}/{page:[0-9]+}/{sort:|value|date}", handler.GetTransactionHistory).Methods("GET", "OPTIONS")
	router.HandleFunc("/history/{id:[0-9]+}/{page:|[0-9]+}", handler.GetTransactionHistory).Methods("GET", "OPTIONS")

	viper.SetConfigFile("config/config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	s.server = &http.Server{
		Addr:         host + ":" + port,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		Handler:      router,
	}
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(c context.Context) error {
	return s.server.Shutdown(c)
}
