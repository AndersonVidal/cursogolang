package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Quotation struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

const APIUrl = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func main() {
	db, err := sql.Open("sqlite3", "./database/applicationDB.db")
	if err != nil {
		log.Panic("[CONN] Error in database connection: " + err.Error())
	}

	defer db.Close()

	err = createQuotationTable(db)
	if err != nil {
		log.Panic("[CONN] Database creation error: " + err.Error())
	}

	http.HandleFunc("/cotacao", GetBid)

	log.Println("Init server on port :8080")
	http.ListenAndServe(":8080", nil)
}

func GetBid(w http.ResponseWriter, r *http.Request) {
	log.Println("Receive GET request on '/cotacao'")

	quotation, err := GetQuotation()
	if err != nil {
		http.Error(w, "Internal Server Error: error to get quotation", http.StatusInternalServerError)
		return
	}

	err = registerQuotation(quotation)
	if err != nil {
		http.Error(w, "Internal Server Error: error in register database quotation", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(quotation.Bid)
	log.Println("Finish process for GET request")
}

func GetQuotation() (*Quotation, error) {
	log.Println("Requesting for " + APIUrl + "...")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", APIUrl, nil)
	if err != nil {
		log.Println("[GET-QUOTATION] Error in build request for `" + APIUrl + "` : " + err.Error())
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("[GET-QUOTATION] Error doing request to `" + APIUrl + "` : " + err.Error())
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		log.Printf("[GET-QUOTATION] Unexpected status code: %d", res.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("[GET-QUOTATION] Error in read response of `" + APIUrl + "` : " + err.Error())
		return nil, err
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		log.Println("[GET-QUOTATION] Error unmarshalling response: " + err.Error())
		return nil, err
	}

	usdBrl, ok := jsonResponse["USDBRL"].(map[string]interface{})
	if !ok {
		log.Println("[GET-QUOTATION] Error to get USDBRL tag from json response")
	}

	usdBrlJson, err := json.Marshal(usdBrl)
	if err != nil {
		log.Println("[GET-QUOTATION] Error to get USDBRL tag from json: " + err.Error())
	}

	var q Quotation
	err = json.Unmarshal(usdBrlJson, &q)
	if err != nil {
		log.Println("[GET-QUOTATION] Error unmarshalling quotation response: " + err.Error())
		return nil, err
	}

	return &q, nil
}

// DATABASE

func createQuotationTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS quotations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		code_id TEXT,
		name TEXT,
		high REAL,
		low REAL,
		var_bid REAL,
		pct_change REAL,
		bid REAL,
		ask REAL,
		timestamp TEXT,
		create_date TEXT
	);`

	_, err := db.Exec(query)

	if err != nil {
		log.Println("[CREATE-QUOTATION-TABLE] Create quotations table failed: " + err.Error())
		return err
	}

	return nil
}

func registerQuotation(q *Quotation) error {
	insert := `INSERT INTO quotations (
		code, code_id, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	db, err := sql.Open("sqlite3", "./database/applicationDB.db")
	if err != nil {
		log.Println("[REGISTER-QUOTATION] Database connection error: " + err.Error())
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err = db.ExecContext(
		ctx,
		insert,
		q.Code,
		q.Codein,
		q.Name,
		q.High,
		q.Low,
		q.VarBid,
		q.PctChange,
		q.Bid,
		q.Ask,
		q.Timestamp,
		q.CreateDate,
	)

	if err != nil {
		log.Println("[REGISTER-QUOTATION] Error during insert quotation in database: " + err.Error())
		return err
	}

	return nil
}
