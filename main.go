package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

// タスク構造体
type task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// ID発行機
func counter(init int) func() int {
	var counter = init
	return func() int {
		counter++
		return counter
	}
}

// タスク新規追加
func addTask(tasks []task, id int, title string) ([]task, task) {
	newTask := task{
		ID:        id,
		Title:     title,
		Completed: false,
	}
	tasks = append(tasks, newTask)
	return tasks, newTask
}

// Struct->JSONへの変換
func toJson(tasks []task) ([]byte, error) {
	b, err := json.Marshal(tasks)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// JSON->Structへの変換
func toStruct(data []byte, tasks *[]task) error {
	if err := json.Unmarshal(data, tasks); err != nil {
		return err
	}
	return nil
}

// 最大ID取得
func getMaxID(tasks []task) int {
	var max int
	for _, task := range tasks {
		if task.ID > max {
			max = task.ID
		}
	}
	return max
}

// 初期化処理
func loadFile(targetPath string, tasks *[]task) error {
	// JSONファイルの読込
	byte, err := os.ReadFile(targetPath)
	if err != nil {
		return err
	}
	// JSON->Struct変換
	if err := toStruct(byte, tasks); err != nil {
		return err
	}
	return nil
}

func main() {

	// ToDoリスト
	var tasks []task

	// ToDoリストの保管先ファイルパス
	var targetPath = "./tasks.json"

	// JSONファイルの読み込み
	if err := loadFile(targetPath, &tasks); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("%sが存在しませんでした。ToDoリストは新規ファイルに保存されます。: %v\n", targetPath, err)
		} else {
			fmt.Printf("内部エラーが発生しました: %v\n", err)
			return
		}
	}

	// IDの最新値を取得
	maxID := getMaxID(tasks)

	// ID発行機クロージャ生成
	idCounter := counter(maxID)

	// フラグの定義
	add := flag.String("add", "", "add a task")
	list := flag.Bool("list", false, "show tasks")
	done := flag.Int("done", 0, "choose task id that is completed")

	// 実行コマンドのフラグ情報を解析
	flag.Parse()

	// フラグごとの処理
	if *add != "" {

		tasks, newTask := addTask(tasks, idCounter(), *add)

		byte, err := toJson(tasks)
		if err != nil {
			fmt.Printf("内部エラーが発生しました: %v\n", err)
		}

		if err := os.WriteFile(targetPath, byte, 0644); err != nil {
			fmt.Printf("ファイルへの書き込みに失敗しました: %v\n", err)
		}

		fmt.Printf("タスクを登録しました: %v\n", newTask)

	}
	if *list != false {

		fmt.Printf("%-6s%-8s%s\n", "ID", "STATUS", "TASK")
		for _, task := range tasks {
			if task.Completed == true {
				fmt.Printf("%-7d%-7s%s\n", task.ID, "[x]", task.Title)
				continue
			}
			fmt.Printf("%-7d%-7s%s\n", task.ID, "[ ]", task.Title)
		}
	}
	if *done != 0 {

		for i, task := range tasks {
			if task.ID == *done {
				tasks[i].Completed = true
			}
		}

		byte, err := toJson(tasks)
		if err != nil {
			fmt.Printf("内部エラーが発生しました: %v\n", err)
		}

		if err := os.WriteFile(targetPath, byte, 0644); err != nil {
			fmt.Printf("ファイルへの書き込みに失敗しました: %v\n", err)
		}
		fmt.Printf("タスクID %dを完了しました\n", *done)
	}
}
