package main

import (
	"exhibtion-proxy/caching"
	"exhibtion-proxy/handlers"
	"fmt"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/igdb"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strconv"
)

func main() {
	proxy := exhibition_proxy_library.Proxy{}
	proxy.Init()

	cachingManager := caching.CachingManager{
		CacheDBPath: filepath.Join(proxy.DataPath, "cache.db"),
		ProxySettings: proxy.Settings,
	}
	err := cachingManager.DBInit()
	if err != nil {
		fmt.Println(err)
		return
	}

	StartHttpServer(&proxy, &cachingManager)
}

func StartHttpServer(p *exhibition_proxy_library.Proxy, cachingManager *caching.CachingManager) {
	if !p.Settings.DebugMode {
		gin.SetMode(gin.ReleaseMode)
	}
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
