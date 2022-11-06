package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

var servers []Server

type Server struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Host string   `json:"host"`
	Tags []string `json:"tags"`
}

func NewServerHandler(ctx *gin.Context) {
	var server Server
	if err := ctx.ShouldBindJSON(&server); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	server.ID = xid.New().String()
	servers = append(servers, server)
	ctx.JSON(http.StatusOK, server)

}
func ListServersHander(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, servers)
}
func UpdateServerHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	var server Server
	if err := ctx.ShouldBindJSON(&server); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	index := -1

	for i := 0; i < len(servers); i++ {
		if servers[i].ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}

	servers[index] = server
	ctx.JSON(http.StatusOK, server)

}
func DeleteServerHandler(ctx *gin.Context) {
	id := ctx.Param("id")

	index := -1

	for i := 0; i < len(servers); i++ {
		if servers[i].ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}
	servers = append(servers[:index], servers[index+1:]...)
	ctx.JSON(http.StatusOK, gin.H{"message": "Recipe has been deleted"})

}
func SearchServersHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfServers := make([]Server, 0)
	for i := 0; i < len(servers); i++ {
		found := false
		for _, t := range servers[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfServers = append(listOfServers, servers[i])
		}
	}
	c.JSON(http.StatusOK, listOfServers)
}

func GetServerHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	for i := 0; i < len(servers); i++ {
		if servers[i].ID == id {
			ctx.JSON(http.StatusOK, servers[i])
			return
		}
	}

	ctx.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})

}
func init() {
	fmt.Printf("ini dipanggil dulu\n")
	servers = make([]Server, 0)
	file, _ := ioutil.ReadFile("servers.json")
	_ = json.Unmarshal([]byte(file), &servers)

}

func main() {
	r := gin.Default()
	r.POST("/servers", NewServerHandler)
	r.GET("/servers", ListServersHander)
	r.GET("/servers/search", SearchServersHandler)
	r.PUT("/servers/:id", UpdateServerHandler)
	r.DELETE("/servers/:id", DeleteServerHandler)
	r.GET("/servers/:id", GetServerHandler)
	r.Run()
}
