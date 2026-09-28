package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Task struct {
	ID    int
	Title string
	Done  bool
}

func (t *Task) Complete() {
	t.Done = true
}

func (t *Task) Rename(name string) {
	t.Title = name
}
func addTask(reader *bufio.Reader, tasks map[int]Task, nextID *int) map[int]Task { // reader — ввод; tasks — задачи; nextID — адрес следующего ID; результат — обновлённая map

	fmt.Println("Добавить задачу") // выводим подсказку пользователю

	task, err := reader.ReadString('\n') // task — введённый текст; err — ошибка чтения; '\n' — Enter

	if err != nil { // nil — ошибки нет; err != nil — чтение завершилось с ошибкой
		fmt.Println("Ошибка:", err) // показываем пользователю причину ошибки
		return tasks                // возвращаем исходные задачи без изменений
	}

	task = strings.TrimSpace(task) // убираем Enter и пробелы по краям строки

	newTask := Task{ // создаём новую структуру задачи
		ID:    *nextID, // *nextID — получаем число по переданному адресу
		Title: task,    // сохраняем очищенный пользовательский текст
		Done:  false,   // новая задача изначально не выполнена
	}

	tasks[*nextID] = newTask // сохраняем задачу в map под ключом текущего ID

	(*nextID)++ // увеличиваем настоящее значение nextID для следующей задачи

	return tasks // возвращаем map с добавленной задачей
}

func showTasks(tasks map[int]Task) {
	fmt.Println("Показать задачи")
	if len(tasks) == 0 {
		fmt.Println("Список пуст")
	} else {
		for ID, task := range tasks { // ID = ключ map, task = текущая задача типа Task
			if task.Done {
				fmt.Println(ID, task.Title, "[✅]")
			} else {
				fmt.Println(ID, task.Title, "[❌]")
			}
		}
	}
}

func deleteTask(tasks map[int]Task) map[int]Task {
	var task_ID int

	fmt.Println("Удалить задачу")
	if len(tasks) == 0 {
		fmt.Println("Список задач пуст")
		return tasks
	}
	fmt.Println("Введите ID задачи:")

	fmt.Scan(&task_ID)

	_, ok := tasks[task_ID]

	if ok {
		delete(tasks, task_ID)
		fmt.Println("Ключ удален")
		return tasks
	}
	fmt.Println("Ключа нету")
	return tasks
}

func deleteTaskByID(tasks map[int]Task, taskID int) map[int]Task {
	_, ok := tasks[taskID]
	// проверяем, существует ли задача с таким ID

	if ok {
		delete(tasks, taskID)
		// удаляем задачу, только если ключ существует
	}

	return tasks
	// возвращаем map:
	// при ok == true задача уже удалена
	// при ok == false map осталась без изменений
}

func changeTask(reader *bufio.Reader, tasks map[int]Task) map[int]Task {
	var task_change int
	fmt.Println("Заменить задачу")
	if len(tasks) == 0 {
		fmt.Println("Список задач пуст")
		return tasks
	}
	fmt.Println("Введите номер задачи для изменения")
	fmt.Scan(&task_change)

	currentTask, ok := tasks[task_change]
	// tasks — map со всеми задачами
	// task_change — ID задачи, которую пользователь хочет переименовать
	// tasks[task_change] — пытаемся найти задачу по этому ID
	// currentTask — копия найденной задачи
	// ok — результат поиска:
	//      true  → задача с таким ID существует
	//      false → задачи с таким ID нет
	if ok {
		fmt.Println("Введите замену.")
		task, err := reader.ReadString('\n') // читаем текст до нажатия Enter; '\n' — символ переноса строки

		if err != nil { // проверяем, не произошла ли ошибка при чтении текста
			fmt.Println("Ошибка:", err) // выводим причину ошибки
			return tasks                // возвращаем map без изменений
		}

		task = strings.TrimSpace(task) // убираем '\n' и лишние пробелы по краям
		currentTask.Rename(task)
		// меняем Title у найденной копии

		tasks[task_change] = currentTask
		// сохраняем изменённую копию обратно в map

		return tasks
		// возвращаем обновлённую map и завершаем функцию
	}

	fmt.Println("Нету такого")
	// эта строка выполнится только тогда, когда ok == false

	return tasks
	// возвращаем прежнюю map без изменений
}

func completeTaskByID(tasks map[int]Task, taskID int) map[int]Task { // tasks — задачи; taskID — ID нужной задачи; map[int]Task — результат после изменения
	currentTask, ok := tasks[taskID] // currentTask — копия задачи; ok — найден ли ключ

	if ok { // работаем только с найденной задачей
		currentTask.Complete()      // у копии задачи меняем Done на true
		tasks[taskID] = currentTask // сохраняем изменённую копию обратно в map
	}

	return tasks // возвращаем изменённую либо исходную map
}

func completeTask(tasks map[int]Task) map[int]Task {
	fmt.Println("Проверка наличия ")
	if len(tasks) == 0 {
		fmt.Println("Список пуст")
		return tasks
	}
	fmt.Println("Введите номер задачи")
	var task_number int
	fmt.Scan(&task_number)
	_, ok := tasks[task_number] // _ — задачу не сохраняем; ok — найден ли ключ task_number
	// tasks — map со всеми задачами
	// task_number — ID задачи, которую пользователь хочет переименовать
	// tasks[task_number] — пытаемся найти задачу по этому ID
	// ok — результат поиска:
	//
	//	true  → задача с таким ID существует
	//	false → задачи с таким ID нет
	//
	// := создаёт сразу две переменные: currentTask и ok
	if ok { // выполняем блок, только если задача существует
		tasks = completeTaskByID(tasks, task_number) // tasks — где меняем; task_number — какую задачу

		fmt.Println("Задача сделана")
		return tasks // возвращаем обновлённую map
	}
	fmt.Println("Нету такого")
	// выполняется только при ok == false

	return tasks
	// возвращаем исходную map без изменений

}
func containsID(ids []int, targetID int) bool { // ids — список ID; targetID — искомый ID; bool — результат поиска
	for _, id := range ids { // _ — пропущенный индекс; id — текущее значение
		if id == targetID { // сравниваем текущий и искомый ID
			return true // ID найден
		}
	}

	return false // ID не найден
}

func findMax(num []int) int { // num — срез чисел; int — найденный максимум
	max := num[0] // max — текущий максимум; сначала это первый элемент

	for _, number := range num { // _ — ненужный индекс; number — текущее число
		if number > max { // текущее число больше сохранённого максимума
			max = number // обновляем текущий максимум
		}
	}

	return max // возвращаем максимум после проверки всего среза
}

func main() {
	nextID := 1
	reader := bufio.NewReader(os.Stdin)
	tasks := make(map[int]Task)
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1. Добавить задачу")
		fmt.Println("2. Показать задачи")
		fmt.Println("3. Удалить задачу")
		fmt.Println("4. Заменить задачу")
		fmt.Println("5. Отметить задачу выполненой")
		fmt.Println("6. Выход")

		var ch int
		fmt.Print("Введите номер:")
		fmt.Scan(&ch)

		switch ch {
		case 1:
			tasks = addTask(reader, tasks, &nextID)
		case 2:
			showTasks(tasks)
		case 3:
			tasks = deleteTask(tasks)
		case 4:
			tasks = changeTask(reader, tasks)
		case 5:
			tasks = completeTask(tasks)
		case 6:
			fmt.Println("Выход")
			return
		default:
			fmt.Println("Не туда тыкнул сарделькой")
		}
	}
}
