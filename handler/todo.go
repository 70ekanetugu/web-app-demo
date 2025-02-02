package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/70ekanetugu/webdemo/repository"
	"gorm.io/gorm"
)

// Todo一覧をidが若い順に5件取得する。クエリーでstart, endが指定されている場合はその範囲で取得する。
func GetTodoList(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	startId := uint(1)
	endId := uint(5)
	startIdInt, startErr := strconv.Atoi(params.Get("start"))
	if startErr == nil {
		startId = uint(startIdInt)
	}
	endIdInt, endErr := strconv.Atoi(params.Get("end"))
	if endErr == nil {
		endId = uint(endIdInt)
	}
	if endId <= startId {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("End must be greater than start"))
		return
	}

	todos, err := repository.GetTodos(startId, endId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("A server error has occured."))
		return
	}
	if todos == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not exists todos"))
		return
	}

	json, err := json.Marshal(todos)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("A server error has occured."))
		return
	}
	w.Write(json)
}

// パス中で指定されたidのTodoを取得する。
func GetTodoById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id required in last path"))
		return
	}
	todo, err := repository.GetTodoById(uint(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("A server error has occured."))
		return
	}
	if todo == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not exists todo"))
		return
	}

	json, err := json.Marshal(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("A server error has occured."))
		return
	}
	w.Write(json)
}

// Todoを保存する。パス中で有効なidが指定されている場合は更新し、idが無いまたは不正なidの場合新規作成する。
func SaveTodo(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only support application/json"))
		return
	}

	todo := &repository.Todo{}
	if err := json.NewDecoder(r.Body).Decode(todo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid Todo data"))
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	isNew := err != nil || id < 1
	if !isNew {
		todo.ID = uint(id)
	}

	if err = repository.SaveTodo(todo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("A server error has occured."))
		return
	}

	if isNew {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// 指定idに一致するTodoのcompletedを反転させる。
func SaveStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.PathValue("id"))
	if id < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id required greater than 1"))
		return
	}

	var todo *repository.Todo
	txErr := repository.Db.Transaction(func(tx *gorm.DB) (err error) {
		if todo, err = repository.GetTodoById(uint(id)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("A server error has occured."))
			return err
		}
		if todo == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not exists todo"))
			return err
		}

		todo.Completed = !todo.Completed

		if err = repository.SaveTodo(todo); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("A server error has occured."))
			return err
		}
		return nil
	})
	if txErr != nil {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
