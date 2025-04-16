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
		metadataList := handleManager.CachingManager.GetMetadataListFromDBbyName(name)
		if metadataList == nil {
			newEntries, err := apiManager.GetGames(name)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, models.Error{
					ErrorMessage: "Error finding games on IGDB side",
					StatusCode:   http.StatusBadRequest,
				})
				return
			}
			fmt.Println("Adding game to local db:", name)
			for _, metadata := range newEntries {
				handleManager.CachingManager.AddMetadataToDB(&metadata)
			}
			metadataList = handleManager.CachingManager.GetMetadataListFromDBbyName(name)
			ctx.JSON(http.StatusOK, metadataList)
			return
		}
		ctx.JSON(http.StatusOK, metadataList)
	}
}
func (handleManager *HandleManager) HandleSearchByID(apiManager *igdb.APIManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idString := ctx.Param("igdbid")
		id, err := strconv.Atoi(idString)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, models.Error{
				ErrorMessage: "Error parsing " + idString + " to int",
				StatusCode:   http.StatusBadRequest,
			})
			return
		}
		metadata := handleManager.CachingManager.GetMetadataFromDBbyID(id)
		// If the entry is not in db fetch from igdb
		if metadata == nil {
			metadata, err = apiManager.GetGameData(id)
			if err != nil {
				fmt.Println(err)
				return
			}
			handleManager.CachingManager.AddMetadataToDB(metadata)
			fmt.Println("not in db... fetching from igdb")
		}
		ctx.JSON(http.StatusOK, metadata)
	}
}
