package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	port := getEnv("PORT", "80")
	staticDir := getEnv("STATIC_DIR", "./www")
	staticPath := getEnv("STATIC_PATH", "/")
	fsEnabled := getEnv("FS_ENABLED", "no")
	fsDir := getEnv("FS_DIR", "./files")
	fsPath := getEnv("FS_PATH", "/files")
	debug := getEnv("DEBUG", "false")
	metrics := getEnv("METRICS", "true")
	metricsPort := getEnv("METRICS_PORT", "9696")

	gin.SetMode(gin.ReleaseMode)

	if debug == "true" {
		gin.SetMode(gin.DebugMode)
	}

	r := gin.New()

	logger, _ := zap.NewProduction()

	if debug == "true" {
		logger, _ = zap.NewDevelopment()
	}

	r.Use(ginzap.Ginzap(logger, time.RFC3339, true))

	if fsEnabled == "yes" {
		r.StaticFS(fsPath, http.Dir(fsDir))
	}

	r.Static(staticPath, staticDir)

	// Prometheus Metrics
	if metrics == "true" {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			log.Fatal(http.ListenAndServe(":"+metricsPort, nil))
			r.GET("/metrics", gin.WrapH(promhttp.Handler()))
		}()
	}

	// Gin web server
	err := r.Run(":" + port)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

// getEnv gets an environment variable or sets a default if
// one does not exist.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}
