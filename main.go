/*
*------------------------------------------------------------------------------
* Система храниея и обработки библиотеки в формате TXT
* Michael S. Merzlyakov AFKA predator_pc@09112018
*------------------------------------------------------------------------------
*/

package main

import (
	"sync"
	"net/url"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"
	"net/http"
	"strconv"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/foolin/gin-template"
	"github.com/jinzhu/gorm"
    "gopkg.in/gcfg.v1"
    "os"
	_"github.com/jinzhu/gorm/dialects/mysql"	
)

// по умолчанию отключаем
var DevDebug = false
var cfg Config // Конфиг инстанс
var db *gorm.DB  //ДБ Инстанс

type FileInfoData struct {
	sync.Mutex
	sizeData int
	splittedData int
	uniqueData int
}

// Получаем мини данные для базы
func getFileInfoData(fn string) FileInfoData {
	var f FileInfoData
	sourceData, _ := ioutil.ReadFile(fn)
	sizeData := len(sourceData)
	stringData := string(sourceData)
    splittedData := Splitter(stringData, " ,.:;-")
	uniqueData := Unique(splittedData)      

	f.sizeData = sizeData
	f.splittedData = len(splittedData)
	f.uniqueData = len(uniqueData)
	return f
}

// Тип для хранения конфига
type Config struct {

    	General struct {
			Name string
			Port int
			Db string
  		}
    	Section struct {
			Name string
			DevDebug bool
    	}
}
// Тип для БД таблицы
type (	
        booksModel struct {
            gorm.Model		
    		Title string `json:"title"`
            Wcount int  `json:"wcount"`
            Ucount int `json:"ucount"`       
            SizeBytes int `json: "sizebytes"` 
	    }
)

/*
*
* Вспомогательные функции
*
*/

func checkerr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// мульити разбивка

func Splitter(s string, splits string) []string {
	m := make(map[rune]int)
	for _, r := range splits {
		m[r] = 1
	}
	splitter := func(r rune) bool { return m[r] == 1 }
	return strings.FieldsFunc(s, splitter)
}

// уникальные значения

func Unique(Slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range Slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// подсчет вхождений

func oneWordStat(Slice []string, Word string) int {
	counter:=0
	for _, entry := range Slice {
		if DevDebug { fmt.Println(entry," = ",Word) }
			if entry == Word { 
				if DevDebug { fmt.Println(entry," = ",Word) }
				counter++ 
			}		
	}
	return counter
}

// Канал убийцы не существующих сущностей 
// уборка мусора и т.п.

func cleanTrash() <-chan string {
    c := make(chan string)
	var books []booksModel //много книг
	var notFoundFiles []string //лист трэша	
	var flagIdentity bool = false // есть нет?

	go func() { //фигачим
			notFoundFiles = nil // сбросим наш лист		
			db.Find(&books) // спросим базу
			files, err := filepath.Glob("lib/*.txt") // спросим файлуху
            checkerr(err)
			// поищем теперь
            for _, b := range books {
                for _,f := range files {
                    if b.Title == f {
                        flagIdentity = true
					}										                 
				}
				if flagIdentity!=true {
					notFoundFiles = append(notFoundFiles, b.Title)
				}
				flagIdentity = false
			}
			// если нашли файлы в базе которых нет которых на файлухе
            if len(notFoundFiles)>0 {
				if DevDebug { fmt.Println("Trash was founded:") }
				for _, b := range books {
					for _,f := range notFoundFiles {				
						if b.Title == f {
							if DevDebug { fmt.Println(b.Title," removing...") }
							db.Unscoped().Delete(&b); //грохаем из базы безвозвратно
						}
					}
				}
            }		
    }()
    return c
}

// Канал основной прямой обработки данных

func myChan() <-chan string {
    c := make(chan string) 
	var books []booksModel //много книг
	var newFiles []string // много файлов
	var flagIdentity bool = false //есть нет?
	var finfo FileInfoData

	go func() {
		for {          //будем крутить до посинения
			var cnt int //счетчик итераций
			db.Find(&books) //прочитаем зановго все книжки 
			db.Find(&books).Count(&cnt) // и сразу кол-во запросим
			newFiles = nil // сбросим лист

			finfo.Lock()
			defer finfo.Unlock()

			//если мы хотим посмотреть все файлы в базе
			if DevDebug { 
				for _, b := range books {
			 		fmt.Println(b.Title," in db library list")
			 	}		
			}

			// прочитаем нашу директорию
			files, err := filepath.Glob("lib/*.txt")
			checkerr(err)

			if DevDebug { fmt.Println("OS file system list = ", files) }
			
			if cnt > 0 {
				//дебаг всякий
				if DevDebug { 
					fmt.Println("[STATUS] scaning for changes...") 
					fmt.Println("[STATUS] Files in DB = ", cnt) 
					fmt.Println("[STATUS] Files at FS = ", len(files)) 
				}

				for _, f := range files {
					for _, b:= range books {

						if f==b.Title {

							if DevDebug { fmt.Printf("[WATCHDOG] %s - %s\n",b.Title,f) }

							finfo:=getFileInfoData(f)

							 //sourceData, _ := ioutil.ReadFile(f)
							 //sizeData := len(sourceData)

							 //sizeData
							 if b.SizeBytes != finfo.sizeData{
							 	if DevDebug {  fmt.Printf("[CHANGED] = %s(%d) - %s(%d)\n",b.Title,b.SizeBytes,f,finfo.sizeData) }

							 //	stringData := string(sourceData)
                        	 	//splittedData := Splitter(stringData, " ,.:;-")
                        	 	//uniqueData := Unique(splittedData)                         

								b.Title = f
								b.Wcount = finfo.splittedData //len(splittedData)
								b.Ucount = finfo.uniqueData  //len(uniqueData)
								b.SizeBytes = finfo.sizeData
							
                        	 	db.Save(&b)   			 							 					
							 }
							flagIdentity = true
						} 
					}
					if flagIdentity!=true {
						newFiles = append(newFiles, f)							
					}
					flagIdentity = false
				}
				flagIdentity = false
				finfo = FileInfoData{}
			} else {
				if DevDebug {  fmt.Println("[STATUS] first run filling db") }
				  for _, f := range files {
					  finfo:=getFileInfoData(f)
					//	sourceData, _ := ioutil.ReadFile(f)
					//	sizeData := len(sourceData)
					//	stringData := string(sourceData)
                    //  splittedData := Splitter(stringData, " ,.:;-")
					//	uniqueData := Unique(splittedData)                         
						
						book := booksModel{Title: f, Wcount:finfo.splittedData, 
										Ucount:  finfo.uniqueData, SizeBytes: finfo.sizeData}
                        db.Save(&book)   
				  }				
				  finfo = FileInfoData{}
			}
			if DevDebug { fmt.Println("[STATUS] new files found ", len(newFiles) ) }
			if len(newFiles) > 0 {
				if DevDebug { fmt.Println("[STATUS] adding new files to db ... ",newFiles) }
				for _, f := range newFiles {
						finfo:=getFileInfoData(f)
					// sourceData, _ := ioutil.ReadFile(f)
					// sizeData := len(sourceData)
					// stringData := string(sourceData)
					// splittedData := Splitter(stringData, " ,.:;-")
					// uniqueData := Unique(splittedData)                         
					// book := booksModel{Title: f, Wcount:len(splittedData), 
					
						book := booksModel{Title: f, Wcount:finfo.splittedData, 
										Ucount:  finfo.uniqueData, SizeBytes: finfo.sizeData}
					db.Save(&book)   
				}							
				finfo = FileInfoData{}
			}
            go cleanTrash() //теперь почистим мусор
			time.Sleep(10 * time.Second) // поспим чуть чуть см. ТЗ
		}
	}()
	return c
}

// Загрузка конфига и обработка параметров

func loadParams() {
	var actionArg string
	_ = actionArg

	if len(os.Args) > 1 {
	    actionArg = os.Args[1]
            fmt.Println("[ Current command ]",actionArg)
	}

	err := gcfg.FatalOnly(gcfg.ReadFileInto(&cfg, "settings.ini"));
	if err!=nil {
		fmt.Println("[ERROR] - %s",err);
                os.Exit(3) // exit anyway
	} else {

		switch actionArg { 
		case "config": {
			if cfg.General.Name!="" || cfg.General.Port!=0 {
		            fmt.Println("[ General ]")
		            if cfg.General.Name != "" {                                 
        		        fmt.Println("[ -- Name ]",cfg.General.Name)             
	        	    }                                                           
		            if cfg.General.Db != "" {                                 
	        	        fmt.Println("[ -- Db ]",cfg.General.Db)             
		            }                                                           
		            if cfg.General.Port != 0 {                                  
	        	        fmt.Println("[ -- Port ]",cfg.General.Port)             
		            } else{                                                     
		                fmt.Println("[ -- Empty ]")                             
		            }
			} 
			if cfg.Section.Name !="" || cfg.Section.DevDebug!=false {
				fmt.Println("[ Section ]");	
				if cfg.Section.Name != "" {
					fmt.Println("[ -- Name ]",cfg.Section.Name);			
				}
				if cfg.Section.DevDebug != false {
					fmt.Println("[ -- Reboot ]",cfg.Section.DevDebug);
				} else{
					fmt.Println("[ -- Empty ]")
				}
			} else {
					fmt.Println("[ Empty configuration ]")
				}	
                os.Exit(3) // exit anyway
		}
		case "run": break;
		default: {
				fmt.Println("Usage: [this-file] command options")		
				fmt.Println("Commands: --run - start API serverice")		
				fmt.Println("          --config - show usable ini file settings")		
				fmt.Println("          --help /none - show this message")		
                os.Exit(3)
			}
		}
	}
} 

// инициализация

func init() {
	loadParams() // грузим конф

	// подключаемся
	if cfg.General.Db !="" {		
		var err error
		db, err = gorm.Open("mysql", cfg.General.Db)
		if err != nil {
			fmt.Println("[Error] description: %s",err);
			panic("[RESULT] Failed to connect database")
		}		
		// подтягиваем модели
		db.AutoMigrate(&booksModel{})

	} else { //базы нет не работаем
		fmt.Println("[ db ]",cfg.General.Db);
		panic("[RESULT] Failed to obtain database connection from settings.ini")
	} 

	if cfg.Section.DevDebug!=false {
		DevDebug = true
	}
}

func main() {
    myChan()   // запускаем наш сканер файлов грузим базу и т.д.

	// запускаем роутер 
	router := gin.Default()
	router.HTMLRender = gintemplate.Default()

	v1 := router.Group("/api/books/")
	{
		v1.GET("/", fetchBooks)
		v1.GET("/:id", fetchBookStat)
	}

	v2 := router.Group("/api/stat/")
	{
		v2.GET("/*info", fetchBookStatByWord)
	}

	if cfg.General.Port!=0 {
		router.Run(":"+strconv.Itoa(cfg.General.Port))
	} else {
		panic("[ERROR] Failed to obtain server port from settings.ini")
	}
}

// отдаем все книги

func fetchBooks(c *gin.Context) {
	var allBooks []booksModel

	db.Select("title, size_bytes, created_at").Find(&allBooks);
	data, _ := json.Marshal(allBooks)

	if DevDebug { fmt.Println(string(data)); }

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": string(data)})	
}

// отдаем статистику по конкретной книге

func fetchBookStat(c *gin.Context){
	var book booksModel

	path:=fmt.Sprintf("lib\\%s",c.Param("id"))

	if DevDebug {  fmt.Println("path = ",path) }

	db.Where("title = ?", path).First(&book)
	data, _ := json.Marshal(book)

	if DevDebug {  fmt.Println(string(data)); }

	if book.Title!="" {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": string(data)})	
	} else {
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "error": "not found"})	
	}
}

// отдаем статистику по слову

func fetchBookStatByWord(c *gin.Context){

	//for Ajax
	c.Writer.Header().Set("Content-Type", "application/json")
    c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	var book booksModel
	counter:=0
	params:=strings.Split(c.Request.RequestURI,"/")

	if DevDebug { fmt.Println("PARAMS",params) }

	if len(params) > 4 {

		www, _ := url.PathUnescape(params[4])

		if DevDebug {  
			fmt.Println("PARAMS --> ",www)
			fmt.Println("[PARAMS]",params)
		}

		path:=fmt.Sprintf("lib\\%s",params[3])

		if DevDebug { fmt.Println("path = ",path) }

		db.Where("title = ?", path).First(&book)		

		if book==(booksModel{}) {
			if DevDebug {  fmt.Println("[NO DATA SELECTED]"); }
		} else  { 
			if DevDebug {  fmt.Println("[COUNTING]"); }
			sourceData, _ := ioutil.ReadFile(path)
			stringData := string(sourceData)
			splittedData := Splitter(stringData, " ,.:;-")

			counter = oneWordStat(splittedData, www)
		} 

		if DevDebug { fmt.Println("[COUNTER] = ",counter) }
		
		if counter > 0 {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "counter": counter, "wcount": book.Wcount})	
		} else {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "error": "not found"})	
		}
	}
}

