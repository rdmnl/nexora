// Nexora: Lightweight Server Monitoring & Resource Tracker
// Developed by Erdem Unal
// © 2024 All rights reserved.
//
// This code is licensed under the MIT License.
// GitHub: https://github.com/rdmnl/nexora

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rdmnl/nexora/config"
	"github.com/rdmnl/nexora/server"
	"github.com/rdmnl/nexora/version"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print the version and exit")
	configPath := flag.String("config-path", "config.yaml", "Path to the configuration file")

	port := flag.String("port", "12321", "Port for the HTTP server (default: 12321)")
	flag.Parse()

	if *versionFlag {
		version.PrintVersion()
		return
	}

	cfg := config.LoadConfig(*configPath)
	server.Init(cfg)

	address := fmt.Sprintf(":%s", *port)
	srv := server.StartHTTPServer(address)

	go server.StartBroadcast()

	fmt.Printf("Server running at http://localhost:%s\n", *port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
