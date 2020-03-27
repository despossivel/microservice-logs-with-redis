# Micro-serviço de logs escrito em GO utilizando Redis

![golang-redis](https://www.restapiexample.com/wp-content/uploads/2018/06/golang-redis-databse-example.png)


#### Acesse o diretório redis contendo os arquivos do docker para o redis

Para efetue o build do container
```bash
docker-compose build
```
 após suba-o

```bash
docker-compose up -d
```

#### Redis comands

Para acessar o container
```bash
docker exec -it <nome do container> /bin/bash
```
Acesse a CLI do redis com
```bash
redis-cli
```
Para listar todas as keys
```bash
KEYS *
```
ou para listar uma KEY expecifica 
```bash
KEYS <nome da key>
```


### Golang

Na pasta raiz do projeto faça o download das dependências

```bash
go mod init
```

Inicie o server
```bash
go run server.go
```

Para fazer o build execute
```bash
go build server.go
```





[Tutorial redis em golang](https://www.alexedwards.net/blog/working-with-redis)

[Tutorial gorilla mux](https://medium.com/@hugo.bjarred/rest-api-with-golang-and-mux-e934f581b8b5)