package server

import (
	"balanceService/pkg/model"
	"balanceService/pkg/repository"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	dbRepo repository.DBRepository
}

func NewHandler(dbRepo repository.DBRepository) *Handler {
	return &Handler{dbRepo: dbRepo}
}

func (handler *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.Atoi(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	if err != nil || userId < 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "wrong id"})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	var account model.Account
	account, err = handler.dbRepo.GetUserBalance(uint32(userId))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "user with given id was not found"})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		log.Fatal(err)
	}
}

func (handler *Handler) Credit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var crediting model.Crediting
	err := json.NewDecoder(r.Body).Decode(&crediting)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "failed to decode json"})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	err = handler.dbRepo.Credit(crediting.ID, crediting.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := fmt.Sprint("transaction failed: ", err)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: errMsg})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "successfully"})
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func (handler *Handler) Reserve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var reserving model.Reserve
	err := json.NewDecoder(r.Body).Decode(&reserving)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "failed to decode json"})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	err = handler.dbRepo.Reserve(
		reserving.UserID,
		reserving.ServiceID,
		reserving.OrderID,
		reserving.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := fmt.Sprint("transaction failed: ", err)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: errMsg})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "successfully"})
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func (handler *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var withdrawing model.Reserve
	err := json.NewDecoder(r.Body).Decode(&withdrawing)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "failed to decode json"})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	err = handler.dbRepo.Withdraw(
		withdrawing.UserID,
		withdrawing.ServiceID,
		withdrawing.OrderID,
		withdrawing.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := fmt.Sprint("transaction failed: ", err)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: errMsg})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "successfully"})
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func (handler *Handler) Refund(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var refund model.Reserve
	err := json.NewDecoder(r.Body).Decode(&refund)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "failed to decode json"})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	err = handler.dbRepo.Refund(
		refund.UserID,
		refund.ServiceID,
		refund.OrderID,
		refund.Value)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := fmt.Sprint("transaction failed: ", err)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: errMsg})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "successfully"})
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func (handler *Handler) GetReport(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	curDate := time.Now()
	vars := mux.Vars(r)
	var year, month int
	var err error
	year, err = strconv.Atoi(vars["year"])
	if err != nil || year < 2007 || year > curDate.Year() {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "wrong year number"})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	month, err = strconv.Atoi(vars["month"])
	if err != nil || month < 1 || month > 12 {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "wrong month number"})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	var filename string
	filename, err = handler.dbRepo.GetReport(uint32(year), uint32(month))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := fmt.Sprint("failed: ", err)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: errMsg})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "report file name: " + filename})
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}

func (handler *Handler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	fmt.Println(vars)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "wrong id"})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	var page int
	if vars["page"] == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(vars["page"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: "wrong page number"})
			if jsonErr != nil {
				log.Fatal(jsonErr)
			}
			return
		}
		if page < 1 {
			page = 1
		}
	}
	var sort model.SortField
	switch vars["sort"] {
	case "":
		sort = model.None
	case "value":
		sort = model.ByValue
	case "date":
		sort = model.ByDate
	}

	var order model.Order
	if vars["order"] == "" || vars["order"] == "asc" {
		order = model.ASC
	} else {
		order = model.DESC
	}
	transactions, err := handler.dbRepo.GetTransactionHistory(uint32(id), uint32(page), sort, order)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errMsg := fmt.Sprint("failed: ", err)
		jsonErr := json.NewEncoder(w).Encode(model.JSONMessage{Message: errMsg})
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	jsonErr := json.NewEncoder(w).Encode(transactions)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
}
