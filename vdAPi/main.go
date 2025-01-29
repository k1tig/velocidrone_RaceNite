package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// add grey group for catch-ups or late qualify

type Bracket struct {
	ID     string `json:"id"`
	Rev    int    `json:"rev"`
	Racers []struct {
		Name        string  `json:"name"`
		Qualifytime float32 `json:"qualifytime"`
	} `json:"racers "`
}

type Racer struct {
	Name        string  `json:"name"`
	Qualifytime float32 `json:"qualifytime"`
}

type Bracket struct {
	ID     string  `json:"id"`
	Rev    int     `json:"rev"`
	Racers []Racer `json:"racers"`
}

/*type bracketStatus struct {
	rev int
}*/

var Brackets []Bracket

func main() {
	router := gin.Default()
	router.GET("/brackets", getBrackets)
	router.POST("/brackets", initBracket)

	router.Run("localhost:8080")
}

func initBracket(c *gin.Context) {
	var initBracket Bracket

	if err := c.BindJSON(&initBracket); err != nil {
		return
	}

	Brackets = append(Brackets, initBracket)
	c.IndentedJSON(http.StatusCreated, initBracket)

}

func getBrackets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, Brackets)
}
