package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

type Todo struct {
	// json変換時にkey名は小文字・スネークケースにしたいのでgorm.Modelは使用しない。
	ID        uint           `gorm:"primarykey" json:"id"`
	Text      string         `gorm:"not null" json:"text"`
	Completed bool           `gorm:"default:false" json:"completed"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Todo一覧を startId~endId の範囲で取得する。
func GetTodos(startId, endId uint) (todos []Todo, err error) {
	todos = make([]Todo, endId-startId+1)
	if err = Db.Where("id >= ? AND id <= ?", startId, endId).Find(&todos).Error; err == nil {
		return todos, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else {
		return nil, err
	}
}

// 指定idのTodoを取得する。
func GetTodoById(id uint) (todo *Todo, err error) {
	todo = new(Todo)
	if err = Db.First(todo, id).Error; err == nil {
		return todo, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else {
		return nil, err
	}
}

// IDがセットされていなければINSERT, セットされていればUPDATEする。
func SaveTodo(todo *Todo) (err error) {
	return Db.Save(todo).Error
}

// DBの初期化処理。サーバー起動前に呼び出すこと。
func Initialize(dbHost, dbPort, dbUser, dbPass string) (err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/todo_demo?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
	)

	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Cannot open connection", slog.String("err", err.Error()))
		return err
	}
	if err = Db.AutoMigrate(&Todo{}); err != nil {
		slog.Error("Failed migration", slog.String("err", err.Error()))
		return err
	}

	return nil
}
