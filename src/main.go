package main

import (
	"exhibtion-proxy/caching"
	"exhibtion-proxy/handlers"
	"fmt"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/igdb"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	//if err := godotenv.Load(); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	proxy := exhibition_proxy_library.Proxy{
		SettingsPath: "../proxy-settings.json",
	}
	proxy.Init()

	cachingManager := caching.CachingManager{}
	err := cachingManager.DBInit()
	if err != nil {
		fmt.Println(err)
		return
	}

	StartHttpServer(&proxy, &cachingManager)
}

func StartHttpServer(p *exhibition_proxy_library.Proxy, cachingManager *caching.CachingManager) {
	router := gin.Default()
	apiManager, err := igdb.NewAPI(p.Settings, p.SettingsManger)
	if err != nil {
		panic(err)
	}

	handleManager := handlers.HandleManager{
		CachingManager: cachingManager,
	}

	router.GET("/game/:igdbid", handleManager.HandleSearchByID(apiManager))
	router.GET("/game/", handleManager.HandleSearchByName(apiManager))

	err = router.Run(":" + strconv.Itoa(p.Settings.Port))
	if err != nil {
		panic(err)
	}
}
