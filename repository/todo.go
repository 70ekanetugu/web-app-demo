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
var todoRepository TodoRepository

type Todo struct {
	// json変換時にkey名は小文字・スネークケースにしたいのでgorm.Modelは使用しない。
	ID        uint           `gorm:"primarykey" json:"id"`
	Text      string         `gorm:"not null" json:"text"`
	Completed bool           `gorm:"default:false" json:"completed"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type TodoRepository interface {
	getTodos(startId, endId uint) ([]Todo, error)
	getTodoById(id uint) (*Todo, error)
	saveTodo(todo *Todo) error
}

type dbRepository struct{}
type testRepository struct{}

// Todo一覧を startId~endId の範囲で取得する。
func GetTodos(startId, endId uint) ([]Todo, error) {
	return todoRepository.getTodos(startId, endId)
}

// 指定idのTodoを取得する。
func GetTodoById(id uint) (*Todo, error) {
	return todoRepository.getTodoById(id)
}

// IDがセットされていなければINSERT, セットされていればUPDATEする。
func SaveTodo(todo *Todo) error {
	return todoRepository.saveTodo(todo)
}

// dbRepositoryの実装
func (d *dbRepository) getTodos(startId, endId uint) (todos []Todo, err error) {
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

func (d *dbRepository) getTodoById(id uint) (todo *Todo, err error) {
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

func (d *dbRepository) saveTodo(todo *Todo) (err error) {
	return Db.Save(todo).Error
}

// fsReepositoryの実装
func (f *testRepository) getTodos(startId, endId uint) (todos []Todo, err error) {
	todos = []Todo{
		{ID: 1, Text: "Todo1", Completed: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 2, Text: "Todo2", Completed: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 3, Text: "Todo3", Completed: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 4, Text: "Todo4", Completed: true, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: 5, Text: "Todo5", Completed: false, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	return todos, nil
}

func (f *testRepository) getTodoById(id uint) (todo *Todo, err error) {
	todo = &Todo{ID: id, Text: "Todo" + fmt.Sprint(id), Completed: false, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return todo, nil
}

func (f *testRepository) saveTodo(todo *Todo) (err error) {
	return nil
}

// DBの初期化処理。サーバー起動前に呼び出すこと。
func Initialize(dbHost, dbPort, dbUser, dbPass string) (err error) {
	if dbHost == "" || dbPort == "" || dbUser == "" || dbPass == "" {
		todoRepository = newTodoRepository(false)
		slog.Info("DB connection info is not set. Use test repository.")
		return nil
	}

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

	slog.Info("DB connection established")
	todoRepository = newTodoRepository(true)
	return nil
}

// TodoRepositoryの実装を取得する。
func newTodoRepository(useDb bool) TodoRepository {
	if useDb {
		return &dbRepository{}
	} else {
		return &testRepository{}
	}
}
