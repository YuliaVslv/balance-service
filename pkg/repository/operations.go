package repository

import (
	"balanceService/pkg/model"
	"balanceService/pkg/report"
	"fmt"
)

func (dbRepo *DBRepository) GetUserBalance(userId uint32) (model.Account, error) {
	row := dbRepo.db.QueryRow(
		"SELECT id, balance "+
			"FROM accounts "+
			"WHERE id = $1",
		userId)
	var accountData model.Account
	err := row.Scan(&accountData.ID, &accountData.Balance)
	return accountData, err
}

func (dbRepo *DBRepository) Credit(userId uint32, value float32) error {
	tx, err := dbRepo.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO accounts AS a (id, balance) "+
			"VALUES ($1, $2) "+
			"ON CONFLICT (id) "+
			"DO UPDATE "+
			"SET balance = a.balance + $2 "+
			"WHERE a.id = $1",
		userId,
		value)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO transactions (type, user_id, value, description) "+
			"VALUES ('credit', $1, $2, $3)",
		userId,
		value,
		fmt.Sprintf("credited to the account"))
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (dbRepo *DBRepository) Reserve(userId uint32, serviceId uint32, orderId uint32, price float32) error {
	tx, err := dbRepo.db.Begin()
	if err != nil {
		return err
	}
	var accountData model.Account
	accountData, err = dbRepo.GetUserBalance(userId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if accountData.Balance < price {
		err = fmt.Errorf("insufficient funds")
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		"UPDATE accounts "+
			"SET balance = balance - $2 "+
			"WHERE id = $1",
		userId,
		price)

	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO reserve (user_id, service_id, order_id, value) "+
			"VALUES ($1, $2, $3, $4)",
		userId,
		serviceId,
		orderId,
		price)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO transactions (type, user_id, value, description) "+
			"VALUES ('reserve', $1, $2, $3)",
		userId,
		price,
		fmt.Sprintf("reservation of funds (model %d, order %d)", serviceId, orderId))
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (dbRepo *DBRepository) Withdraw(userId uint32, serviceId uint32, orderId uint32, price float32) error {
	tx, err := dbRepo.db.Begin()
	if err != nil {
		return err
	}
	row := dbRepo.db.QueryRow(
		"SELECT value "+
			"FROM reserve "+
			"WHERE user_id = $1 and service_id = $2 and order_id = $3",
		userId,
		serviceId,
		orderId)

	var reserveData model.Reserve
	err = row.Scan(&reserveData.Value)
	if err != nil {
		tx.Rollback()
		return err
	}

	if reserveData.Value != price {
		err = fmt.Errorf("the price in the reserve does not match the asking price")
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		"DELETE FROM reserve "+
			"WHERE user_id = $1 and service_id = $2 and order_id = $3",
		userId,
		serviceId,
		orderId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO transactions (type, user_id, value, description) "+
			"VALUES ('debit', $1, $2, $3)",
		userId,
		price,
		fmt.Sprintf("withdrawal of funds from the reserve account (model %d, order %d)", serviceId, orderId))
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO reports AS r (year, month, service_id, revenue) "+
			"VALUES (EXTRACT(YEAR FROM CURRENT_DATE), EXTRACT(MONTH FROM CURRENT_DATE), $1, $2) "+
			"ON CONFLICT (year, month, service_id) "+
			"DO UPDATE "+
			"SET revenue = r.revenue + $2 "+
			"WHERE r.year = EXTRACT(YEAR FROM CURRENT_DATE) "+
			"and r.month = EXTRACT(MONTH FROM CURRENT_DATE) "+
			"and r.service_id = $1",
		serviceId,
		price)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (dbRepo *DBRepository) Refund(userId uint32, serviceId uint32, orderId uint32, price float32) error {
	tx, err := dbRepo.db.Begin()
	if err != nil {
		return err
	}
	row := dbRepo.db.QueryRow(
		"SELECT value "+
			"FROM reserve "+
			"WHERE user_id = $1 and service_id = $2 and order_id = $3",
		userId,
		serviceId,
		orderId)

	var reserveData model.Reserve
	err = row.Scan(&reserveData.Value)
	if err != nil {
		tx.Rollback()
		return err
	}

	if reserveData.Value != price {
		err = fmt.Errorf("the price in the reserve does not match the asking price")
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		"DELETE FROM reserve "+
			"WHERE user_id = $1 and service_id = $2 and order_id = $3",
		userId,
		serviceId,
		orderId)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		"UPDATE accounts "+
			"SET balance = balance + $2 "+
			"WHERE id = $1",
		userId,
		price)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO transactions (type, user_id, value, description) "+
			"VALUES ('refund', $1, $2, $3)",
		userId,
		price,
		fmt.Sprintf("return of funds to the account (model %d, order %d)", serviceId, orderId))
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (dbRepo *DBRepository) GetReport(year uint32, month uint32) (string, error) {
	rows, err := dbRepo.db.Query(
		"SELECT service_id, revenue "+
			"FROM reports "+
			"WHERE year = $1 and month = $2",
		year,
		month)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	data := make([][]string, 1)
	headers := []string{"service_id", "revenue"}
	data[0] = headers

	for rows.Next() {
		row := make([]string, 2)
		err = rows.Scan(&row[0], &row[1])
		if err != nil {
			return "", err
		}
		data = append(data, row)
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}
	return report.CreateReportFile(data, year, month)
}

func (dbRepo *DBRepository) GetTransactionHistory(
	userId uint32,
	page uint32,
	sort model.SortField,
	order model.Order) ([]model.Transaction, error) {

	var numTrOnPage uint32
	numTrOnPage = 10
	limit := numTrOnPage
	offset := numTrOnPage * (page - 1)
	rows, err := dbRepo.db.Query(""+
		"SELECT date, value, description "+
		"FROM transactions "+
		"WHERE user_id = $1 "+
		"ORDER BY "+string(sort)+" "+string(order)+" "+
		"LIMIT $2 "+
		"OFFSET $3",
		userId,
		limit,
		offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := make([]model.Transaction, 0)
	for rows.Next() {
		var tr model.Transaction
		err = rows.Scan(&tr.Date, &tr.Value, &tr.Description)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tr)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
