package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// add grey group for catch-ups or late qualify

type Bracket struct {
	BracketID string `json:"bracketid"`
	Rev       int    `json:"rev"`
	Racers    []struct {
		RaceID      int     `json:"raceid"`
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
	router.PUT("/brackets/:id", editBracket)

	router.Run("localhost:8080")
}

func initBracket(c *gin.Context) {
	var newBracket Bracket
	newBracketOK := true

	if err := c.BindJSON(&newBracket); err != nil {
		return
	}
	for _, i := range Brackets {
		if i.BracketID == newBracket.BracketID {
			c.IndentedJSON(http.StatusOK, gin.H{"message": " error, id already exists"}) // not correct status
			newBracketOK = false
			break
		}
	}
	if newBracketOK {
		Brackets = append(Brackets, newBracket)
		c.IndentedJSON(http.StatusCreated, Brackets)
	}

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
		if b.BracketID == id {
			bracketUpdated := false
			for i, oringalRacer := range b.Racers {
				for _, editRacer := range bracket.Racers {
					if editRacer.RaceID == oringalRacer.RaceID {
						Brackets[x].Racers[i] = editRacer
						if !bracketUpdated {
							bracketUpdated = true
						}
					}
				}
			}
			if bracketUpdated {
				c.IndentedJSON(http.StatusOK, Brackets)
			} else {
				c.IndentedJSON(http.StatusOK, gin.H{"message": "no update to brackets"})
			}
			break
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bracket not found"})
	}

}
