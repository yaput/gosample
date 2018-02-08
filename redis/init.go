package redis

import (
	"log"
	"net/http"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"gopkg.in/tokopedia/logging.v1"
)

type RedisConfig struct {
	MaxIdle int
	Address string
}
type Config struct {
	Redis RedisConfig
}
type RedisModule struct {
	cfg *Config
}

var Pool *redigo.Pool

func New(c *RedisConfig) {
	Pool = &redigo.Pool{
		MaxIdle:     c.MaxIdle,
		IdleTimeout: 24 * time.Hour,
		Dial:        func() (redigo.Conn, error) { return redigo.Dial("tcp", c.Address) },
	}
}
func NewRedisModule() *RedisModule {
	var cfg Config
	ok := logging.ReadModuleConfig(&cfg, "config", "redis") || logging.ReadModuleConfig(&cfg, "files/etc/redis", "redis")
	if !ok {
		// when the app is run with -e switch, this message will automatically be redirected to the log file specified
		log.Fatalln("failed to read config")
	}
	// this message only shows up if app is run with -debug option, so its great for debugging
	logging.Debug.Println("redis init called")
	New(&cfg.Redis)
	return &RedisModule{
		cfg: &cfg,
	}
}
func (m *RedisModule) SetHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	value := r.FormValue("value")
	if err := Set(key, value); err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// w.Write([]byte(fmt.Sprintln("success: key:%s - value:%s", key, value)))
	if str, err := Get(key); err != nil {
		w.Write([]byte(err.Error()))
		return
	} else {
		w.Write([]byte(str))
		return
	}
	return
}
func Set(key string, value interface{}) error {
	con := Pool.Get()
	defer con.Close()
	_, err := con.Do("SET", key, value)
	return err
}
func Get(key string) (string, error) {
	con := Pool.Get()
	defer con.Close()
	return redigo.String(con.Do("GET", key))
}

func Incr(key string) (string, error) {
	con := Pool.Get()
	defer con.Close()
	return redigo.String(con.Do("INCR", key))
}
