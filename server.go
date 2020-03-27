package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
)

//criando o pool de conexçoes do redis
var pool *redis.Pool

func initPool() {
	pool = &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				log.Printf("ERROR: Fail init redis pool: %s", err.Error())
				os.Exit(1)
			}
			return conn, err
		},
	}
}

func ping(conn redis.Conn) {
	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		log.Printf("ERROR: fail ping redis conn: %s", err.Error())
		os.Exit(1)
	}
}

//Redis PUT
func set(key string, application string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, application)
	if err != nil {
		log.Printf("Error: fail key %s, val %s", key, err.Error())
		return err
	}
	return nil
}

//Redis GET
func get(key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", key))

	if err != nil {
		log.Printf("ERROR: fail get key %s, error %s", key, err.Error())
		return "", err
	}

	return s, nil
}

//Redis SADD
func sadd(key string, val string) error {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SADD", key, val)
	if err != nil {
		log.Printf("ERROR: fail add val %s to set %s, error %s", val, key, err.Error())
		return err
	}

	return nil
}

//Redis SMEMBERS
func smembers(key string) ([]string, error) {
	// get conn and put back when exit from method
	conn := pool.Get()
	defer conn.Close()

	s, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		log.Printf("ERROR: fail get set %s , error %s", key, err.Error())
		return nil, err
	}

	return s, nil
}

// func initStore() {
// 	// get conn and put back when exit from method
// 	conn := pool.Get()
// 	defer conn.Close()

// 	macs := []string{"00000C  Cisco", "00000D  FIBRONICS", "00000E  Fujitsu",
// 		"00000F  Next", "000010  Hughes"}

// 	for _, mac := range macs {
// 		pair := strings.Split(mac, "  ")
// 		set(pair[0], pair[1])
// 	}
// }

func main() {
	initPool()

	// initStore()

	// obter valor que existe
	// log.Printf(get("00000E"))

	// // obter valor que não existe
	// log.Printf(get("0000E"))

	// // adicionar membros para definir
	// sadd("mystiko", "0000E")
	// sadd("mystiko", "0000D")

	// // obter memebers do conjunto
	// s, _ := smembers("mystiko")
	// log.Printf("%v", s)

	r := mux.NewRouter()

	// IMPORTANTE: você deve especificar um correspondente do método OPTIONS para o middleware para definir cabeçalhos CORS
	r.HandleFunc("/{key}", handler).Methods(http.MethodGet) //, http.MethodPut, http.MethodPatch, http.MethodOptions

	r.HandleFunc("/publish", handlerPost).Methods(http.MethodPost)

	r.Use(mux.CORSMethodMiddleware(r))

	http.ListenAndServe(":8083", r)
}

type Registre struct {
	Key         string `json:"key"`
	Application string `json:"application"`
	Body        string `json:"body"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json") //definindo heade json

	params := mux.Vars(r) //obter params da rota

	data, err := get(params["key"])
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(data)

	// w.Write([]byte(params["key"]))
}

func handlerPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var registre Registre
	_ = json.NewDecoder(r.Body).Decode(&registre)

	err := set(registre.Key, registre.Application)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode("sucesso")

}
