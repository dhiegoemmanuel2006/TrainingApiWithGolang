package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// album base para representar como seria um entidade no banco de dados
type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

// albums para simular um banco de dados onde seram feitas as operações
var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}
// main base de entrada para a aplicação
func main() {
    router := gin.Default()
    router.GET("/albums", getAlbums)
    router.GET("/albums/:id", getAlbumByID)
    router.POST("/albums", postAlbums)
    router.DELETE("/albums/:id", deleteAlbumByID)
    router.PUT("/albums/:id", updateAlbumByID)

    router.Run("localhost:8080")
}


// função que retorna toda a lista de albums
func getAlbums(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, albums)
    log.Println("Alguém consultou nosso repositório de albums")
}

// postAlbums função para receber um JSON com os dados para criar um novo album
func postAlbums(c *gin.Context) {
    var newAlbum album

    if err := c.BindJSON(&newAlbum); err != nil {
        log.Println("Uma requisição foi feita  com  um formato de  JSON incorreto para criar um album")
        c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "erro  no formato da requisição"})
        return
    }

    albums = append(albums, newAlbum)
    c.IndentedJSON(http.StatusCreated, newAlbum)
    log.Println("Um album foi criado com sucesso")
}

// getAlbumByID função para busca o album pelo seu id e retorna-lo
func getAlbumByID(c *gin.Context) {
    id := c.Param("id")
    
    for _, a := range albums {
        if a.ID == id {
            c.IndentedJSON(http.StatusOK, a)
            return
        }
    }
    c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// deleteAlbumByID função para buscar o album pelo seu id e depois remove-lo da lista
func deleteAlbumByID(c *gin.Context){
    id := c.Param("id")

    var index, err = findAlbumIndex(albums, id)

    if err != nil{
        c.IndentedJSON(http.StatusNotFound, gin.H{"mensagem": err.Error()})
        log.Println("Alguma requisição para deletar não consegui ser executada")
        return
    }

    albums = append(albums[:index], albums[index+1:]...)
    c.IndentedJSON(http.StatusOK, gin.H{"message": "Item deletado com sucesso"})
    log.Println("Algum item foi deletado com sucesso")
}


//updateAlbumByID função para coletar da request do id e depois recebe o body e atualiza o album
func updateAlbumByID(c *gin.Context){
    id := c.Param("id")

    var updatedAlbum album
    
    if err := c.ShouldBindJSON(&updatedAlbum); err != nil{
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "formato do json está incorreto"})
        log.Println("Algum erro ao tentar dar update em algum album")
        return
    }
    for i, a := range albums {
        if a.ID == id {
            updatedAlbum.ID = id // garantir que o ID não seja alterado
            albums[i] = updatedAlbum
            c.JSON(http.StatusOK, updatedAlbum)
            log.Println("Album atualizado com sucesso")
            return
        }
    }

    c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
    log.Println("Album para o método update não encontrado")
}

// findAlbumIndex função auxiliar para conseguir buscar o index de determinado Album
func findAlbumIndex(albumsList []album, itemIndex string) (int, error) {
    for i, a := range albumsList {
        if a.ID == itemIndex {
            return i, nil
        }
    }
    return -1, errors.New("nenhum item encontrado com este index")
}     