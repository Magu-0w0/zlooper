package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/martinlindhe/notify"
	_ "github.com/mattn/go-sqlite3"
)

// Импортируем штуту для быз данных(database/sql)
// Импоритуем логи
// Импортируем либу для обращения к файловой системе операционной системы
// либа для sqlite в go
// Либа обращающаяся ко времени в системе

func Exists(name string) bool { // Фунуция проверки файла
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

type Work struct {
	Id         int
	Text       string
	Complected bool
	Timestamp  time.Time
}

//	Содаём тело базы данных ( строки и типы строк )

func CreateWork(text string) { // Функция создания строки в базе данных
	db, err := sql.Open("sqlite3", "works.db") // Открываем базу данных
	if err != nil {                            // Проверяем, есть ли ошибки
		panic(err) // Паника !!
	}
	defer db.Close() // Говорим компилятору о том, что набо бы закрыть базу после завершения программы
	_, err = db.Exec("insert into works (text, complected) values ($1, $2)",
		text, false) // Добовляем строку в базу данных. где $1 это text, $2 это false
	if err != nil { // Панически проверяем, есть ли ошибка
		panic(err)
	}
}

func ListWorks() *[]Work { // Функция считываения базы в массив
	db, err := sql.Open("sqlite3", "works.db") // Открываем базу
	if err != nil {                            //Проверяем, не обосрались ли мы
		panic(err) //Паникуем если обосрались
	}
	defer db.Close()                             // В лифте родился что ли ?
	rows, err := db.Query("select * from works") // Достаем из таблицы works
	if err != nil {
		panic(err)
	}
	defer rows.Close() // В лифте родился что ли ?

	works := []Work{} // Создаём массив из структуры Work

	for rows.Next() { // Бегаем по строкам базы данных, добовляя их в массив структуры
		work := Work{}
		err := rows.Scan(&work.Id, &work.Text, &work.Complected, &work.Timestamp) // Читаем строку из массива и записывает её во временную структуру
		if err != nil {
			log.Println(err) // " что то явно сломалось "
			continue
		}
		works = append(works, work) // добовляем временную структуру в массив
	}

	return &works // Возврощаем массив с структурами
}

func ComplectedWork(id int) { // Функция говорящая нам о том, что мы завершаем оповещение ( задаем значение Complected 1 )
	db, err := sql.Open("sqlite3", "works.db") //
	if err != nil {
		panic(err) // вставай, ты обосрался
	}
	defer db.Close() // Закрываем за собой базу

	_, err = db.Exec("update works set complected = $1 where id = $2", true, id) // Обновляем строку с айдишником равным тому айдишнику, который мы ищем
	if err != nil {
		panic(err) // oh no, chringe
	}
}

func init() { // Инициализируем приложение
	if !Exists("works.db") { // Проверяем наличие базы данных
		log.Printf("[WARN] Database no found! Creating database.") // Здравствуйте, мы обосрались, досвидания
		db, err := sql.Open("sqlite3", "works.db")                 // Открываем базу данных
		if err != nil {                                            // Панически проверяем, не навернулосб ли все
			panic(err)
		}
		defer db.Close()                                                                                                                                               // Снова говорим компилятору о том, что набо бы закрыть базу после завершения программы
		_, err = db.Exec("CREATE TABLE works (id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE, text TEXT, complected BOOL, timestamp DATETIME DEFAULT CURRENT_TIMESTAMP)") // Cоздаем новую таблицу в базе
		if err != nil {                                                                                                                                                // Панически проверяем
			panic(err)
		}
	}
}

func MenuAddWork() { // Функция ввода в меню
	var scanner = bufio.NewScanner(os.Stdin) //
	fmt.Printf("> ")
	for scanner.Scan() {
		output := strings.TrimSpace(scanner.Text())
		if output != "" {
			CreateWork(output)
			break
		} else {
			fmt.Println("Error your a stupid")
			fmt.Printf("> ")
		}
	}
}

func MenuComplectedWork() {
	var scanner = bufio.NewScanner(os.Stdin)
	fmt.Printf("> ")
	for scanner.Scan() {
		output := strings.TrimSpace(scanner.Text())
		if output != "" {
			i, err := strconv.Atoi(output)
			if err != nil {
				fmt.Println("No working string to use")
			}
			ComplectedWork(i)
			break
		} else {
			fmt.Println("Error your a stupid")
			fmt.Printf("> ")
		}
	}
}
func NotifyRing() {
	// отображение уведомления
	notify.Notify("ZlooPer", "cum", "in my ass", "/Data/jesus.jpg") // "Название программы" "главный текст" "иконка"
}

func main() { // Функция интерфейса
	//NotifyRing()
	fmt.Println("Welcome to zlooper...") // Заголовок интерфейса
	fmt.Println("Menu:")
	fmt.Println("1) Add new work")
	fmt.Println("2) List works")
	fmt.Println("3) Complected work")
	fmt.Println("4) Exit")
	var scanner = bufio.NewScanner(os.Stdin) // Чтение ввода строки
	fmt.Printf("> ")                         // Вывод строки
	for scanner.Scan() {                     // Функция вывода интерфейса
		output := strings.TrimSpace(scanner.Text())
		if output != "" {
			switch output {
			case "1":
				MenuAddWork()
			case "2":
				works := ListWorks()
				for _, work := range *works {
					if work.Complected {
						fmt.Printf("%d|%s|%s\n", work.Id, work.Text, "Complected")
					} else {
						fmt.Printf("%d|%s|%s\n", work.Id, work.Text, "No complected")
					}
				}
			case "3":
				MenuComplectedWork()
			case "4":
				os.Exit(0)
			}
		}
		fmt.Printf("> ")
	}

	err := scanner.Err()
	if err != nil {
		panic(err)
	}
}
