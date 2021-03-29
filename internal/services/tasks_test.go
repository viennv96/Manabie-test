/*
Add by VienNV
*/

package services

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sqllite "github.com/manabie-com/togo/internal/storages/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

//Note: update valid Token before run test
const validToken string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MTcxMDQ4ODYsInVzZXJfaWQiOiJmaXJzdFVzZXIifQ.XdMAfvd1KD1t_0HmrTG75CcCmFNfpSh1Zu9jUgQgDfs"

var db, _ = sql.Open("sqlite3", "./../../data.db")

var s = &ToDoService{
	JWTKey: "wqGyEBBfPK9w3Lxw",
	Store: &sqllite.LiteDB{
		DB: db,
	},
}

var dbErr, _ = sql.Open("sqlite3", "./../../data1.db")

var sErr = &ToDoService{
	JWTKey: "wqGyEBBfPK9w3Lxw",
	Store: &sqllite.LiteDB{
		DB: dbErr,
	},
}

func TestMethodOptions(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodOptions, "/", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestLoginSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/login?user_id=firstUser&password=example", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestLoginWrongUsername(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/login?user_id=first&password=example", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 401 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestLoginInvalidUsername(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/login?user_id=&password=example", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 401 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/login?user_id=firstUser&password=exp", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 401 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/login?user_id=firstUser&password=", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 401 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestGetTaskListSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/tasks?created_date=2020-06-29", nil)
	request.Header.Add("Authorization", validToken)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestCreateTaskInvalidToken(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)

	writer := httptest.NewRecorder()
	requestBody := strings.NewReader("{\n    \"content\": \"another content some thing test\"\n}")
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Add("Authorization", "")
	mux.ServeHTTP(writer, request)

	if writer.Code != 401 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestCreateTaskDbError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", sErr.ServeHTTP)

	writer := httptest.NewRecorder()
	requestBody := strings.NewReader("{\n    \"content\": \"another content some thing test\"\n}")
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Add("Authorization", validToken)
	mux.ServeHTTP(writer, request)

	if writer.Code != 500 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestCreateTaskRequestError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)

	writer := httptest.NewRecorder()
	requestBody := strings.NewReader("{;fff}")
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Add("Authorization", validToken)
	mux.ServeHTTP(writer, request)

	if writer.Code != 500 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestCreateTaskSuccess(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.ServeHTTP)

	writer := httptest.NewRecorder()
	requestBody := strings.NewReader("{\n    \"content\": \"another content some thing test\"\n}")
	request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
	request.Header.Add("Authorization", validToken)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

func TestCreateTaskLimitReached(t *testing.T) {
	for ;;{
		mux := http.NewServeMux()
		mux.HandleFunc("/", s.ServeHTTP)

		writer := httptest.NewRecorder()
		requestBody := strings.NewReader("{\n    \"content\": \"another content some thing test\"\n}")
		request, _ := http.NewRequest(http.MethodPost, "/tasks", requestBody)
		request.Header.Add("Authorization", validToken)
		mux.ServeHTTP(writer, request)

		if writer.Code != 429 && writer.Code != 200 {
			t.Errorf("Response code is %v", writer.Code)
			return
		} else if writer.Code == 429 {
			return
		}
	}
}

func TestGetTaskListDbError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", sErr.ServeHTTP)

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/tasks?created_date=2020-06-29", nil)
	request.Header.Add("Authorization", validToken)
	mux.ServeHTTP(writer, request)

	if writer.Code != 500 {
		t.Errorf("Response code is %v", writer.Code)
	}
}

