package main

import (
	"context"
	"log"
	"os"
	containerinfo "pinger/containerInfo"
	"strings"
	"time"

	"github.com/docker/docker/client"
)

var (
	logger = log.New(os.Stdout, "docker-pinger: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
)

func parseEnv() containerinfo.Env {
	var env containerinfo.Env
	networksEnv := os.Getenv("DOCKER_NETWORKS")
	if networksEnv == "" {
		logger.Fatal("DOCKER_NETWORKS environment variable is required")
		os.Exit(1)
	}
	backendUrl := os.Getenv("BACKEND_URL")
	if backendUrl == "" {
		logger.Fatal("BACKEND_URL environment variable is required")
		os.Exit(1)
	}
	networkList := strings.Split(networksEnv, ",")
	logger.Printf("Monitoring networks: %v\n", networkList)
	env = containerinfo.Env{
		Networks: networkList,
		BackURL:  backendUrl,
	}
	return env
}

func main() {
	env := parseEnv()
	for {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			logger.Fatal("Docker client init error:", err)
		}

		if _, err := cli.Ping(context.Background()); err != nil {
			logger.Fatalf("Docker API connection error: %v", err)
		}
		logger.Println("Successfully connected to Docker API")
		containerinfo.CheckContainers(cli, env)
		time.Sleep(5 * time.Second)
		cli.Close()
	}
}