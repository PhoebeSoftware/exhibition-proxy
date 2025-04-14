package handlers

import (
	"exhibtion-proxy/caching"
	"fmt"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/igdb"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type HandleManager struct {
	CachingManager *caching.CachingManager
}


func (handleManager *HandleManager) HandleSearchByName(apiManager *igdb.APIManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Query("name")
		if name == "" {
			ctx.JSON(http.StatusBadRequest, models.Error{
				ErrorMessage: "No search query",
				StatusCode:   http.StatusBadRequest,
			})
			return
		}

		gameDataList, err := apiManager.GetGames(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		ctx.JSON(http.StatusOK, gameDataList)
	}
}
func (handleManager *HandleManager) HandleSearchByID(apiManager *igdb.APIManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idString := ctx.Param("igdbid")
		id, err := strconv.Atoi(idString)
		if err != nil {
			fmt.Println(err)
			return
		}
		gameData, err := apiManager.GetGameData(id)
		if err != nil {
			fmt.Println(err)
			return
		}
		handleManager.CachingManager.AddMetadataToDatabase(gameData)
		ctx.JSON(http.StatusOK, gameData)
	}
}
