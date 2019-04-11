package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
	case "POST":
	case "HEAD":
	default:
		// Unsupportted method PUT etc..
		JSONOptions(w, r)
		return
	}

	testhost := r.FormValue("host")
	testurl := r.FormValue("url")

	w.Header().Set("Host-Var-Length", strconv.Itoa(len(testhost)))
	w.Header().Set("Url-Var-Length", strconv.Itoa(len(testurl)))

	switch {
	case (len(testhost) > 0):
		_, err := url.Parse("http://" + testhost + "/")
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed host varialble: %s", testhost))
			return
		}
	case (len(testurl) > 0):
		parsedURL, err := url.Parse(testurl)
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed url varialble: %s", testurl))
			return
		}
		testhost = parsedURL.Hostname()
	default:
		JSONError(w, r, fmt.Errorf("host or url varialbles are required"))
		return
	}
	w.Header().Set("X-Test-Domain", testhost)

	statement, err := database.Prepare("INSERT or REPLACE INTO blacklist (domain) VALUES (?)")
	if err != nil {
		JSONError(w, r, err)
		return

	}
	_, err = statement.Exec(testhost)
	if err != nil {
		JSONError(w, r, err)
		return

	}

	newmap := make(map[string]interface{})
	newmap["Insert-Domain"] = testhost

	json.NewEncoder(w).Encode(newmap)

}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
	case "POST":
	case "HEAD":
	default:
		// Unsupportted method PUT etc..
		JSONOptions(w, r)
		return
	}
	testhost := r.FormValue("host")
	testurl := r.FormValue("url")
	w.Header().Set("Host-Var-Length", strconv.Itoa(len(testhost)))
	w.Header().Set("Url-Var-Length", strconv.Itoa(len(testurl)))

	switch {
	case (len(testhost) > 0):
		_, err := url.Parse("http://" + testhost + "/")
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed host varialble: %s", testhost))
			return
		}
	case (len(testurl) > 0):
		parsedURL, err := url.Parse(testurl)
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed url varialble: %s", testurl))
			return
		}
		testhost = parsedURL.Hostname()
	default:
		JSONError(w, r, fmt.Errorf("host or url varialbles are required"))
		return
	}
	w.Header().Set("X-Test-Domain", testhost)

	statement, err := database.Prepare("DELETE FROM blacklist WHERE domain = ?")
	if err != nil {
		JSONError(w, r, err)
		return

	}
	_, err = statement.Exec(testhost)
	if err != nil {
		JSONError(w, r, err)
		return

	}

	newmap := make(map[string]interface{})
	newmap["Delete-Domain"] = testhost

	json.NewEncoder(w).Encode(newmap)
}

func demoIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
	case "POST":
	case "HEAD":
	default:
		// Unsupportted method PUT etc..
		JSONOptions(w, r)
		return
	}
	testhost := r.FormValue("host")
	testurl := r.FormValue("url")
	w.Header().Set("Host-Var-Length", strconv.Itoa(len(testhost)))
	w.Header().Set("Url-Var-Length", strconv.Itoa(len(testurl)))

	switch {
	case (len(testhost) > 0):
		_, err := url.Parse("http://" + testhost + "/")
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed host varialble: %s", testhost))
			return
		}
	case (len(testurl) > 0):
		parsedURL, err := url.Parse(testurl)
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed url varialble: %s", testurl))
			return
		}
		testhost = parsedURL.Hostname()
	default:
		JSONError(w, r, fmt.Errorf("host or url varialbles are required"))
		return
	}
	w.Header().Set("X-Test-Domain", testhost)
}

func search(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
	case "POST":
	case "HEAD":
	default:
		// Unsupportted method PUT etc..
		JSONOptions(w, r)
		return
	}
	testhost := r.FormValue("host")
	testurl := r.FormValue("url")
	w.Header().Set("Host-Var-Length", strconv.Itoa(len(testhost)))
	w.Header().Set("Url-Var-Length", strconv.Itoa(len(testurl)))

	switch {
	case (len(testhost) > 0):
		_, err := url.Parse("http://" + testhost + "/")
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed host varialble: %s", testhost))
			return
		}
	case (len(testurl) > 0):
		parsedURL, err := url.Parse(testurl)
		if err != nil {
			JSONError(w, r, fmt.Errorf("malformed url varialble: %s", testurl))
			return
		}
		testhost = parsedURL.Hostname()
	default:
		JSONError(w, r, fmt.Errorf("host or url varialbles are required"))
		return
	}

	response := make(map[string]interface{})
	response["OK"] = "1"

	sqlStmt := `SELECT domain FROM blacklist WHERE domain = ?`
	var resVar string
	err := database.QueryRow(sqlStmt, testhost).Scan(&resVar)
	switch {
	case err == sql.ErrNoRows:
	case err != nil:
		JSONError(w, r, fmt.Errorf(err.Error()))
		return
	default:
		// There is a row and there is no error aka blacklisted
		response["blacklisted"] = "true"
		response["OK"] = "0"
		w.Header().Set("X-Vote", "BLOCK")
	}

	json.NewEncoder(w).Encode(response)
}
