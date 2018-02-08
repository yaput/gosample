package page

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/yaput/gosample/redis"
)

type ServerConfig struct {
	Name string
}

type Config struct {
	Server ServerConfig
}

type PageModule struct {
	cfg *Config
}

type (
	Category struct {
		ID   int
		Name string
		URL  string
	}

	User struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		MSISDN       string `json:"msisdn"`
		Email        string `json:"email"`
		Birth_date   string `json:"birth_date"`
		Created_time string `json:"created_time"`
		Update_time  string `json:"update_time"`
		User_age     string `json:"user_age"`
	}

	Rows struct {
		Rows []User `json:"rows"`
	}
)

var dbHome *sql.DB
var err error

// func NewPageModule() *PageModule {

// 	var cfg Config

// 	ok := logging.ReadModuleConfig(&cfg, "config", "hello") || logging.ReadModuleConfig(&cfg, "files/etc/gosample", "hello")
// 	if !ok {
// 		// when the app is run with -e switch, this message will automatically be redirected to the log file specified
// 		log.Fatalln("failed to read config")
// 	}

// 	// this message only shows up if app is run with -debug option, so its great for debugging
// 	logging.Debug.Println("Page init called", cfg.Server.Name)

// 	return &PageModule{
// 		cfg: &cfg,
// 	}

// }

// Connection to Database
func connection() (*sql.DB, error) {
	return sql.Open("postgres", "postgres://st140804:apaajadeh@devel-postgre.tkpd/tokopedia-user?sslmode=disable")
}

// Index return index page and with visit counter
func Index(w http.ResponseWriter, r *http.Request) {
	// Set Redis
	_ = redis.NewRedisModule()
	var count int
	if c, err := redis.Get("counter_yaput"); c == "" || c == "0" {
		if err != nil {
			err = redis.Set("counter_yaput", 1)
			if err != nil {
				count = 1
			}
		}
		log.Println(err.Error())
	} else {
		_, err := redis.Incr("counter_yaput")
		if err != nil {
			tmp, _ := redis.Get("counter_yaput")
			count, _ = strconv.Atoi(tmp)
		} else {
			log.Println(err.Error())
		}
	}
	span, _ := opentracing.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()
	// Connection to Database
	dbHome, err = connection()
	if err != nil {
		log.Print(err)
	}
	defer dbHome.Close()
	// Execute HTML page
	t, _ := template.ParseFiles("page/view.html") //setp 1
	t.Execute(w, count)                           //step 2
}

// GetUsers return list of Users from User Table
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Connection to Database
	dbHome, err = connection()
	if err != nil {
		log.Print(err)
	}
	defer dbHome.Close()
	// Query Start
	stmt, err := dbHome.Prepare("select user_id, full_name, msisdn, user_email, birth_date, create_time, update_time from ws_user limit 1000")
	rows, err := stmt.Query()
	if err != nil {
		log.Print(err)
	}
	UserList := []User{}
	defer rows.Close()
	for rows.Next() {
		c := &User{}
		err := rows.Scan(&c.ID, &c.Name, &c.MSISDN, &c.Email, &c.Birth_date, &c.Created_time, &c.Update_time)
		if err != nil {
			log.Print("Error scan rows")
		}

		UserList = append(UserList, *c)
	}
	for _, el := range UserList {
		log.Print(el)
	}
	hasil, err := json.Marshal(UserList)
	if err != nil {
		log.Println("ERROR JSON")
	}

	// Return JSON Result
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hasil))
}
