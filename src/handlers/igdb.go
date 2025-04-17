package handlers

import (
	"exhibtion-proxy/caching"
	"fmt"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/igdb"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/proxy_models"
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
			ctx.JSON(http.StatusBadRequest, proxy_models.Error{
				ErrorMessage: "No search query",
				StatusCode:   http.StatusBadRequest,
			})
			return
		}
		metadataList := handleManager.CachingManager.GetMetadataListFromDBbyName(name)
		// If the entry is not in db fetch from igdb
		if metadataList == nil {
			newEntries, err := apiManager.GetGames(name)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, proxy_models.Error{
					ErrorMessage: "Rate limit exceeded",
					StatusCode:   http.StatusBadRequest,
				})
				fmt.Println(err)
				return
			}
			if len(newEntries) > 0 {
				for _, metadata := range newEntries {
					handleManager.CachingManager.AddMetadataToDB(&metadata)
					fmt.Println("Adding game to local db:", metadata.Name)
				}
				metadataList = handleManager.CachingManager.GetMetadataListFromDBbyName(name)
				ctx.JSON(http.StatusOK, metadataList)
				return
			}
			ctx.JSON(http.StatusBadRequest, proxy_models.Error{
				StatusCode:   http.StatusBadRequest,
				ErrorMessage: "No games found with by this name",
			})
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
			ctx.JSON(http.StatusBadRequest, proxy_models.Error{
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
			fmt.Println("Adding game to local db:", metadata.Name)
			handleManager.CachingManager.AddMetadataToDB(metadata)
		}
		ctx.JSON(http.StatusOK, metadata)
	}
}
