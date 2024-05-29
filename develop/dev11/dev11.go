package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

/*
Реализовать HTTP-сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP-библиотекой.


В рамках задания необходимо:
Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
Реализовать middleware для логирования запросов


Методы API:
POST /create_event
POST /update_event
POST /delete_event
GET /events_for_day
GET /events_for_week
GET /events_for_month


Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09). В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON-документ содержащий либо {"result": "..."} в случае успешного выполнения метода, либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
Реализовать все методы.
Бизнес логика НЕ должна зависеть от кода HTTP сервера.
В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400.
В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
*/

type event struct {
	UserId int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Info   string    `json:"info"`
}

type response struct {
	Result []event `json:"result"`
}

type responseErr struct {
	Error string `json:"error"`
}

var events map[int][]event

func createEventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var userIDStr, dateStr, info string

	err = validateInput(body, &userIDStr, &dateStr, &info)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Bad request", 400)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)
	date, _ := time.Parse("2006-01-02", dateStr)
	ev := event{
		UserId: userID,
		Date:   date,
		Info:   info,
	}
	events[userID] = append(events[userID], ev)

	resp := response{Result: []event{ev}}
	respJson, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Fprintf(w, string(respJson))
}

func updateEventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var userIDStr, dateStr, info string

	err = validateInput(body, &userIDStr, &dateStr, &info)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Bad request", 400)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)
	date, _ := time.Parse("2006-01-02", dateStr)
	ev := event{
		UserId: userID,
		Date:   date,
		Info:   info,
	}
	_, ok := events[userID]
	if !ok {
		respErr := responseErr{Error: "user does not exist"}
		reply, _ := json.MarshalIndent(respErr, "", "\t")
		slog.Error("user does not exist")
		fmt.Fprintf(w, string(reply))
		return
	}

	var resp response

	for id, evt := range events[userID] {
		if ev.Date == evt.Date {
			events[userID][id].Info = ev.Info
			resp.Result = append(resp.Result, events[userID][id])
		}
	}

	respJson, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Fprintf(w, string(respJson))
}

func deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	var userIDStr, dateStr, info string

	err = validateInput(body, &userIDStr, &dateStr, &info)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Bad request", 400)
		return
	}

	userID, _ := strconv.Atoi(userIDStr)
	date, _ := time.Parse("2006-01-02", dateStr)
	ev := event{
		UserId: userID,
		Date:   date,
		Info:   info,
	}
	_, ok := events[userID]
	if !ok {
		respErr := responseErr{Error: "user does not exist"}
		reply, _ := json.MarshalIndent(respErr, "", "\t")
		slog.Error("user does not exist")
		fmt.Fprintf(w, string(reply))
		return
	}

	if len(events[userID]) < 1 {
		slog.Error("nothing to delete")
		http.Error(w, "Bad request", 400)
		return
	}

	for id, evt := range events[userID] {
		if ev.Date == evt.Date {
			events[userID] = slices.Delete(events[userID], id, id+1)
		}
	}

	resp := response{Result: []event{ev}}
	respJson, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Fprintf(w, string(respJson))
}

func eventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	err := validate(w, r)
	if err != nil {
		return
	}

	date, _ := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))

	var resp response

	y, m, d := date.Date()
	for _, ev := range events[userID] {
		y2, m2, d2 := ev.Date.Date()

		if y == y2 && m == m2 && d == d2 {
			resp.Result = append(resp.Result, ev)
		}
	}
	respJson, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal server Error", 500)
		return
	}
	fmt.Fprintf(w, string(respJson))
}

func eventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	err := validate(w, r)
	if err != nil {
		return
	}

	date, _ := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))

	var resp response

	y, week := date.ISOWeek()
	for _, ev := range events[userID] {
		y2, week2 := ev.Date.ISOWeek()

		if y == y2 && week == week2 {
			resp.Result = append(resp.Result, ev)
		}
	}
	respJson, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal server Error", 500)
		return
	}
	fmt.Fprintf(w, string(respJson))
}

func eventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	err := validate(w, r)
	if err != nil {
		return
	}

	date, _ := time.Parse("2006-01-02", r.URL.Query().Get("date"))
	userID, _ := strconv.Atoi(r.URL.Query().Get("user_id"))

	var resp response

	y, m, _ := date.Date()
	for _, ev := range events[userID] {
		y2, m2, _ := ev.Date.Date()

		if y == y2 && m == m2 {
			resp.Result = append(resp.Result, ev)
		}
	}
	respJson, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Internal server Error", 500)
		return
	}
	fmt.Fprintf(w, string(respJson))
}

func validateInput(body []byte, id *string, d *string, i *string) error {
	query := strings.Split(string(body), "&")
	if len(query) != 3 {
		slog.Error("incorrect input")
		return errors.New("incorrect input")
	}

	for _, item := range query {
		box := strings.Split(item, "=")
		switch box[0] {
		case "user_id":
			*id = box[1]
		case "date":
			*d = box[1]
		case "info":
			*i = box[1]
		}
	}
	_, err := strconv.Atoi(*id)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	_, err = time.Parse("2006-01-02", *d)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	return nil
}

func validate(w http.ResponseWriter, r *http.Request) error {
	dateStr := r.URL.Query().Get("date")
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Bad request", 400)
		return err
	}

	_, err = time.Parse("2006-01-02", dateStr)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Bad request", 400)
		return err
	}

	_, ok := events[userID]
	if !ok {
		respErr := responseErr{Error: "user does not exist"}
		reply, _ := json.MarshalIndent(respErr, "", "\t")
		slog.Error("user does not exist")
		fmt.Fprintf(w, string(reply))
		return errors.New("user does not exist")
	}
	return nil
}

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	l.handler.ServeHTTP(w, r)
	log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}

func main() {
	events = make(map[int][]event)
	box := event{
		UserId: 3,
		Date:   time.Now(),
		Info:   "test 1",
	}
	box2 := event{
		UserId: 3,
		Date:   time.Now().AddDate(0, 0, 1),
		Info:   "test 2",
	}

	events[box.UserId] = append(events[box.UserId], box)
	events[box2.UserId] = append(events[box2.UserId], box2)

	http.HandleFunc("POST /create_event", createEventHandler)
	http.HandleFunc("POST /update_event", updateEventHandler)
	http.HandleFunc("POST /delete_event", deleteEventHandler)
	http.HandleFunc("GET /events_for_day", eventsForDayHandler)
	http.HandleFunc("GET /events_for_week", eventsForWeekHandler)
	http.HandleFunc("GET /events_for_month", eventsForMonthHandler)

	wrappedMux := NewLogger(http.DefaultServeMux)

	fmt.Println("Server is running at http://localhost:8082")
	log.Fatal(http.ListenAndServe(":8082", wrappedMux))
}
