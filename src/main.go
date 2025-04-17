package main

import (
	"exhibtion-proxy/caching"
	"exhibtion-proxy/handlers"
	"fmt"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/igdb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	godotenv.Load()

	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = filepath.Join(".", "data")
	}

	if err := os.MkdirAll(dataPath, 0777); err != nil {
		fmt.Println(err)
		fmt.Println("Could not create path: " + dataPath)
		return
	}

	proxy := exhibition_proxy_library.Proxy{
		SettingsPath: filepath.Join(dataPath, "proxy-settings.json"),
	}
	proxy.Init()

	cachingManager := caching.CachingManager{
		CacheDBPath: filepath.Join(dataPath, "cache.db"),
	}
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
