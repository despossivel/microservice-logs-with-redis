package main

import (
    "fmt"
    "log"
	"net/http"
    "github.com/gorilla/mux"
    "github.com/gomodule/redigo/redis"
)

func main() {
    r := mux.NewRouter()

    // IMPORTANTE: você deve especificar um correspondente do método OPTIONS para o middleware para definir cabeçalhos CORS
    r.HandleFunc("/", handler).Methods(http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodOptions)
    r.Use(mux.CORSMethodMiddleware(r))
    
    http.ListenAndServe(":8083", r)
}



func redisConnect(conn){
    // Estabeleça uma conexão com o servidor Redis atendendo na porta
    // 6379 da máquina local. 6379 é a porta padrão, portanto, a menos que
    // você já alterou o arquivo de configuração Redis, isso deve
    // trabalhos.
    conn, err := redis.Dial("tcp", "localhost:6379")
    if err != nil {
        log.Fatal(err)
    }

    // Importante, use adiar para garantir que a conexão esteja sempre
    // fechado corretamente antes de sair da função main ().
	// defer conn.Close()


    // Envie nosso comando pela conexão. O primeiro parâmetro para
    // Do () é sempre o nome do comando Redis (neste exemplo
    // HMSET), opcionalmente seguido por quaisquer argumentos necessários (neste
    // exemplo da chave, seguida pelos vários campos e valores de hash).
    _, err = conn.Do("HMSET", "album:2", "title", "Electric Ladyland", "artist", "Jimi Hendrix", "price", 4.95, "likes", 8)
        if err != nil {
            log.Fatal(err)
        }


}



func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    if r.Method == http.MethodOptions {
        return
    }

    w.Write([]byte("Golang"))
}