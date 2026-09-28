package main // package — пакет файла; main — основной пакет приложения

import ( // import — подключение внешних пакетов; () — список подключений
	"bufio"   // работа с буферизированным вводом; нужен для bufio.NewReader
	"strings" // работа со строками; нужен для strings.NewReader
	"testing" // стандартное тестирование; нужен для testing.T
)

func TestComplete(t *testing.T) { // TestComplete — название теста; t — управление тестом; *testing.T — указатель на тест
	task := Task{ // task — переменная; Task — структура задачи; {} — создание структуры
		Done: false, // Done — статус выполнения; false — задача не выполнена
	}

	task.Complete() // . — вызов метода; Complete меняет Done на true

	if !task.Done { // if — условие; ! — «не»; task.Done — фактический статус
		t.Error("после Complete() ожидали Done = true, но получили false") // Error — провалить тест
	}
}

func TestCompleteTaskByID(t *testing.T) { // t позволяет сообщать об ошибках теста
	tasks := map[int]Task{ // создаём начальную map с задачами
		1: { // 1 — ключ задачи в map
			ID:   1,     // ID самой задачи
			Done: false, // до вызова функция ещё не выполнена
		},
	}

	tasks = completeTaskByID(tasks, 1) // выполняем задачу с ключом 1 и сохраняем результат

	foundTask, ok := tasks[1] // foundTask — найденная задача; ok — существует ли ключ 1

	if !ok { // если ключ не найден, функция почему-то потеряла задачу
		t.Error("ожидали задачу с ID 1, но она не найдена")
	}

	if !foundTask.Done { // если Done == false, функция не выполнила задачу
		t.Error("ожидали Done = true, но получили false")
	}
}

func TestCompleteTaskByIDMissing(t *testing.T) { // проверяем вызов с отсутствующим ID
	tasks := map[int]Task{ // int — тип ключа; Task — тип значения
		1: {
			ID:   1,     // ID существующей задачи
			Done: false, // до вызова задача не выполнена
		},
	}

	tasks = completeTaskByID(tasks, 2) // tasks — список задач; 2 — отсутствующий ID

	foundTask, ok := tasks[1] // foundTask — задача 1; ok — найден ли ключ 1

	if !ok { // ошибка, если существующая задача исчезла
		t.Error("ожидали задачу с ID 1, но она не найдена")
	}

	if foundTask.Done { // ошибка, если задача 1 случайно стала выполненной
		t.Error("ожидали Done = false, но получили true")
	}
}

func TestRename(t *testing.T) { // TestRename — тест Rename; t — управление тестом
	task := Task{ // Task — структура задачи; {} — создание значения
		Title: "Старое название", // Title — поле задачи; справа записано исходное значение
	}

	expectedTitle := "Новое название" // := — создать переменную; expectedTitle — ожидаемый результат
	task.Rename(expectedTitle)        // .Rename — вызвать метод и передать новое название

	if task.Title != expectedTitle { // != — «не равно»; сравниваем результат с ожиданием
		t.Errorf( // Errorf — провалить тест и подставить значения
			"ожидали Title = %q, получили %q", // %q — место для строки в кавычках
			expectedTitle, // значение для первого %q
			task.Title,    // значение для второго %q
		)
	}
}

func TestRenameSeveralTitles(t *testing.T) { // тест нескольких вариантов Rename
	tests := []struct { // tests — сценарии; [] — срез; struct — набор связанных полей
		name     string // name — название подтеста; string — строка
		oldTitle string // oldTitle — исходный Title
		newTitle string // newTitle — ожидаемый Title
	}{
		{
			name:     "обычное название", // имя первого сценария
			oldTitle: "Старое",           // исходное значение
			newTitle: "Новое",            // ожидаемое значение
		},
		{
			name:     "название из нескольких слов", // имя второго сценария
			oldTitle: "старый бычий пенис",          // исходное значение
			newTitle: "новый моржовый бубенец",      // ожидаемое значение
		},
	}

	for _, test := range tests { // range — перебор; _ — индекс не нужен; test — текущий сценарий
		t.Run(test.name, func(t *testing.T) { // t.Run — подтест; test.name — его название
			task := Task{ // новая задача для текущего сценария
				Title: test.oldTitle, // берём исходный Title из test
			}

			task.Rename(test.newTitle) // передаём ожидаемое название текущего сценария

			if task.Title != test.newTitle { // сравниваем фактический и ожидаемый Title
				t.Errorf( // проваливаем текущий подтест
					"ожидали %q, получили %q",
					test.newTitle, // ожидаемый результат
					task.Title,    // фактический результат
				)
			}
		})
	}
}

func TestDeleteTaskByIDSeveralCases(t *testing.T) { // тест разных вариантов удаления
	tests := []struct { // tests — срез структур с тестовыми сценариями
		name       string // название подтеста
		deleteID   int    // ID для удаления; int — целое число
		expectedOK bool   // ожидаемое наличие ключа; bool — true/false
	}{
		{
			name:       "существующий ключ удаляется", // название первого сценария
			deleteID:   1,                             // удаляем существующий ключ 1
			expectedOK: false,                         // после удаления ключа быть не должно
		},
		{
			name:       "несуществующий ключ", // название второго сценария
			deleteID:   2,                     // пытаемся удалить отсутствующий ключ 2
			expectedOK: true,                  // ключ 1 должен остаться
		},
	}

	for _, test := range tests { // test — текущий сценарий; _ — индекс пропускаем
		t.Run(test.name, func(t *testing.T) { // запускаем сценарий как отдельный подтест
			tasks := map[int]Task{ // tasks — map; int — ключ; Task — значение
				1: { // 1 — ключ map
					ID: 1, // ID — поле структуры Task
				},
			}

			tasks = deleteTaskByID(tasks, test.deleteID) // = — заменить значение; вызываем функцию из first.go
			_, ok := tasks[1]                            // _ — Task не нужен; ok — существует ли ключ 1

			if ok != test.expectedOK { // сравниваем фактический bool с ожидаемым
				t.Errorf( // проваливаем подтест при несовпадении
					"ожидали %t, получили %t", // %t — место для значения bool
					test.expectedOK,           // ожидаемый bool
					ok,                        // фактический bool
				)
			}
		})
	}
}

func TestAddTaskReadError(t *testing.T) { // проверяем addTask при ошибке чтения
	tasks := map[int]Task{}          // при ошибке должна остаться пустой
	nextID := 1                      // при ошибке не должен увеличиться
	input := strings.NewReader("")   // пустой источник; ReadString получит EOF
	reader := bufio.NewReader(input) // создаём тип, который принимает addTask

	tasks = addTask(reader, tasks, &nextID) // вызываем addTask из first.go

	if len(tasks) != 0 { // ошибка, если при неудачном чтении появилась задача
		t.Errorf("ожидали 0 задач, получили %d", len(tasks))
	}

	if nextID != 1 { // ошибка, если ID увеличился несмотря на ошибку
		t.Errorf("ожидали nextID = %d, получили %d", 1, nextID) // 1 — ожидаемое; nextID — фактическое
	}
}

func TestAddTask(t *testing.T) { // тест функции addTask; t — управление тестом
	tasks := map[int]Task{} // пустая map: int — ключ, Task — значение
	nextID := 1             // следующий ID; := создаёт переменную

	input := strings.NewReader("Купить хлеб\n") // NewReader — искусственный ввод; \n — Enter
	reader := bufio.NewReader(input)            // bufio.Reader — тип ввода для addTask

	tasks = addTask(reader, tasks, &nextID) // & — адрес nextID; результат сохраняем в tasks
	createdTask, ok := tasks[1]             // createdTask — найденная Task; ok — найден ли ключ

	if !ok { // !ok — ключ 1 не найден
		t.Error("ожидали задачу с ID 1, но она не найдена") // проваливаем тест
	}

	if createdTask.Title != "Купить хлеб" { // проверяем поле Title
		t.Errorf( // Errorf подставляет значения в сообщение
			"ожидали Title = %q, получили %q", // %q — строка в кавычках
			"Купить хлеб",                     // ожидаемый Title
			createdTask.Title,                 // фактический Title
		)
	}

	if nextID != 2 { // следующий ID должен увеличиться до 2
		t.Errorf(
			"ожидали nextID = %d, получили %d", // %d — целое число
			2,      // ожидаемое значение
			nextID, // фактическое значение
		)
	}

	if createdTask.ID != 1 { // проверяем ID внутри структуры
		t.Errorf(
			"ожидали ID = %d, получили %d",
			1,              // ожидаемый ID
			createdTask.ID, // фактический ID
		)
	}

	if createdTask.Done { // условие истинно, если Done == true
		t.Error("ожидали Done = false, но получили true") // новая задача не должна быть выполнена
	}
}

func TestContainsID(t *testing.T) { // тест containsID; t — управление тестом
	ids := []int{1, 2, 3, 4}     // ids — срез; []int — несколько целых чисел
	actual := containsID(ids, 2) // actual — результат; ищем ID 2 в ids

	if !actual { // !actual — функция вернула false
		t.Error("ожидали true, получили false") // проваливаем тест, потому что ID 2 существует
	}
	missingActual := containsID(ids, 5) // ищем отсутствующий ID 5
	if missingActual {                  // ошибка, если функция вернула true
		t.Error("ожидали false, получили true")
	}
}
func TestSeveralContainsID(t *testing.T) {
	tests := []struct {
		name           string // название сценария
		findID         int    // какой ID ищем
		expectedResult bool   // какой результат ожидаем
	}{
		{
			name:           "существующий ID найден",
			findID:         2,
			expectedResult: true,
		},
		{
			name:           "несуществующий ID не найден",
			findID:         10,
			expectedResult: false,
		},
	}

	ids := []int{1, 2, 3, 4} // общие исходные данные

	for _, test := range tests { // _ — индекс не нужен; test — текущий сценарий; tests — все сценарии
		t.Run(test.name, func(t *testing.T) { // test.name — название подтеста; t — управление этим подтестом
			actual := containsID(ids, test.findID) // вызываем функцию приложения из first.go

			if actual != test.expectedResult { // сравниваем фактический и ожидаемый результаты
				t.Errorf("ожидали %t, получили %t", test.expectedResult, actual)
			}
		})
	}
}

func TestFindMax(t *testing.T) { // t — управление тестом; *testing.T — стандартный тип теста
	numbers := []int{4, 9, 2, 7, 3, 1, 5} // numbers — исходные числа; ожидаемый максимум — 9

	actual := findMax(numbers) // actual — результат, фактически возвращённый findMax

	if actual != 9 { // ошибка, если фактический максимум не равен ожидаемому
		t.Errorf("фактический максимум %d не равен ожидаемому %d ", actual, 9)
	} // подставляется вместо первого %d 9,подставляется вместо второго %d
}
