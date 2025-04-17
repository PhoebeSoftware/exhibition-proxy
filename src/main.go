package main

import (
	"exhibtion-proxy/caching"
	"exhibtion-proxy/handlers"
	"fmt"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library"
	"github.com/PhoebeSoftware/exhibition-proxy-library/exhibition-proxy-library/igdb"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	//if err := godotenv.Load(); err != nil {
	//	fmt.Println(err)
	//	return
	//}

	dataPath := filepath.Join(".", "data")
	if err := os.MkdirAll(filepath.Dir(dataPath), 0666); err != nil {
		log.Fatal(err)
		return
	}

	proxy := exhibition_proxy_library.Proxy{
		SettingsPath: filepath.Join(dataPath, "proxy-settings.json"),
	}
	proxy.Init()

	if proxy.Settings.IgdbSettings.IgdbClient == "fill-in-pls" ||
		proxy.Settings.IgdbSettings.IgdbSecret == "fill-in-pls" {
		fmt.Println("Failed to launch: Please fill in the IGDB client and secret")
		fmt.Println("Config file path: " + proxy.SettingsPath)
		return
	}

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
