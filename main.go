package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

// ToDo（タスク）を管理するタスクマネージャ
// タスクマネージャはタスクが保存されるJSONファイルを読み込み、
// taskスライス、次に割り振るタスクIDを管理します
type taskManager struct {
	tasks []task        // タスク状態管理
	targetPath string   // タスク保存先パス
}

// ToDo（タスク）を表すエンティティ
type task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// タスクマネージャを初期化するファクトリ関数
// targetPathからJSONファイルを読み込み、taskスライスにデコードします
// JSONファイルがない場合は初回起動とみなして空のスライスを生成します
func newTaskManager(targetPath string) (*taskManager, error) {
	// 空のタスクマネージャを生成
    taskManager := &taskManager{
		tasks: []task{},
		targetPath: targetPath,
	}
	// JSONファイルの読み込み
	if err := taskManager.loadTasks(); err != nil {
		// JSONファイルが存在しないケースは初回起動（正常）とみなす
		if errors.Is(err, os.ErrNotExist) {
        	return taskManager, nil
   		}
		return nil, err
	}
	return taskManager, nil
}

// 最大ID取得
func(m *taskManager) getMaxID() int{
	var maxID int
	for _, task := range m.tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID
}

func(m *taskManager) issueID() int{
	return m.getMaxID() + 1
}

// タスク新規追加
func(m *taskManager) addTask(title string) (task) {
	newTask := task{
		ID:        m.issueID(),
		Title:     title,
		Completed: false,
	}
	m.tasks = append(m.tasks, newTask)
	return newTask
}

// タスク完了
func(m *taskManager) completeTask(done int) error {
	for i, task := range m.tasks {
		if task.ID == done {
			m.tasks[i].Completed = true
			return nil
		}
	}
	return fmt.Errorf("task id %d does not exist\n", done)
}

// Struct->JSONへの変換
func(m *taskManager) marshalTasks() ([]byte, error) {
	b, err := json.Marshal(m.tasks)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// JSON->Structへの変換
func(m *taskManager) unmarshalTasks(data []byte) error {
	if err := json.Unmarshal(data, &m.tasks); err != nil {
		return err
	}
	return nil
}

// ファイル保存
func(m *taskManager) save() error{
	byte, err := m.marshalTasks()
	if err != nil {
		return err
	}
	if err := os.WriteFile(m.targetPath, byte, 0644); err != nil {
		return err
	}
	return  nil
}

// タスク読み込み処理
func(m *taskManager) loadTasks() error {
	// JSONファイルの読込
	data, err := os.ReadFile(m.targetPath)
	if err != nil {
		return err
	}
	// JSON->Struct変換
	if err := m.unmarshalTasks(data); err != nil {
		return err
	}
	return nil
}

// UI関数
func viewTasks(tasks []task) {
	fmt.Printf("%-6s%-8s%s\n", "ID", "STATUS", "TASK")
	for _, task := range tasks {
		if task.Completed {
			fmt.Printf("%-7d%-7s%s\n", task.ID, "[x]", task.Title)
			continue
		}
		fmt.Printf("%-7d%-7s%s\n", task.ID, "[ ]", task.Title)
	}
}

func main() {

	// ToDoリストの保管先ファイルパス
	var targetPath = "./tasks.json"

	// タスクマネージャ初期化
	taskManager, err := newTaskManager(targetPath)
	if err != nil {
		fmt.Printf("ファイルの読み込み中にエラーが発生しました: %v\n", err)
		fmt.Println("Todo管理プログラムを終了します")
		return
	}

	// フラグの定義
	add := flag.String("add", "", "add a task")
	list := flag.Bool("list", false, "show registered tasks")
	done := flag.Int("done", 0, "set task completed")

	// 実行コマンドのフラグ情報を解析
	flag.Parse()

	// フラグごとの処理
	if *add != "" {

		newTask := taskManager.addTask(*add)

		if err := taskManager.save(); err != nil {
			fmt.Printf("ファイルへの書き込みに失敗しました: %v\n", err)
			return
		}
		fmt.Printf("タスクを登録しました: %s（ID: %d）\n", newTask.Title, newTask.ID)
	}
	if *list {
		viewTasks(taskManager.tasks)
	}
	if *done != 0 {
		if err := taskManager.completeTask(*done); err != nil {
			fmt.Println("該当のタスクIDは存在しません")
			return
		}
		if err := taskManager.save(); err != nil {
			fmt.Printf("ファイルへの書き込みに失敗しました: %v\n", err)
			return
		}
		fmt.Printf("タスクID %dを完了しました\n", *done)
	}
}
