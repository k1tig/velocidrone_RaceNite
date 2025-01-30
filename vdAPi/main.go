package main

import (
	"fmt"
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

/*type bracketStatus struct {
	rev int
}*/

var Brackets []Bracket

func main() {
	router := gin.Default()
	router.GET("/brackets", getBrackets)
	router.POST("/brackets", initBracket)
	router.POST("/brackets/:id", editBracket)

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

func editBracket(c *gin.Context) {
	var bracket Bracket
	if err := c.BindJSON(&bracket); err != nil {
		return
	}
	id := c.Param("id")
	for x, b := range Brackets {
		if b.ID == id {
			for _, oringalRacer := range b.Racers {
				for _, editRacer := range bracket.Racers {
					if editRacer.Name == oringalRacer.Name {
						fmt.Println("Do something with object at", x) // not really. Need to add "IDs to racers to allow for editing all fields"
					}
				}
			}
		}
	}
}
