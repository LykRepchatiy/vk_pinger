package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

var (
	logger = log.New(os.Stdout, "docker-pinger: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
)

type DBContainer struct {
	ContainerID string            `json:"containerID"`
	IP          map[string]string `json:"ip"`
	Status      string            `json:"status"`
	Timestamp   time.Time         `json:"timestamp"`
	Datestamp   time.Time         `json:"datestamp"`
}

type Env struct {
	Networks []string
	BackURL  string
}

var env Env

func sendToBack(req []DBContainer) {
	json, err := json.Marshal(req)
	if err != nil {
		logger.Println(err)
	}
	http.Post(env.BackURL, "application/json", bytes.NewBuffer(json))
}

func ParseEnv() {
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
	env = Env{networkList, backendUrl}
}

func checkContainers(cli *client.Client) {
	allContainers := make(map[string]types.Container)

	for _, network := range env.Networks {
		containers, err := getNetworkContainers(cli, network)
		if err != nil {
			logger.Printf("Error getting containers in network %s: %v\n", network, err)
			continue
		}

		for _, c := range containers {
			if _, exists := allContainers[c.ID]; !exists {
				allContainers[c.ID] = c
			}
		}
	}

	req := []DBContainer{}
	for _, c := range allContainers {
		status, err := getContainerStatus(cli, c.ID)
		if err != nil {
			logger.Printf("Container %s status error: %v\n", c.Names[0], err)
			continue
		}
		
		containerNetworks := getContainerNetworks(c, env.Networks)
		ips := getContainerIPs(c, env.Networks)
		pingTime := time.Now()
		req = append(req, DBContainer{
			ContainerID: c.ID,
			IP:          ips,
			Status:      status,
			Timestamp:   pingTime,
			Datestamp:   pingTime,
		})
		go sendToBack(req)
		logger.Printf(
			"Container: %-20s Status: %-10s Networks: %-30s IPs: %v\n",
			c.Names[0],
			status,
			strings.Join(containerNetworks, ", "),
			ips,
		)
	}
}

func getContainerIPs(c types.Container, networkList []string) map[string]string {
	ips := make(map[string]string)
	for _, network := range networkList {
		if net, exists := c.NetworkSettings.Networks[network]; exists {
			ips[network] = net.IPAddress
		}
	}
	return ips
}

func getContainerNetworks(c types.Container, networkList []string) []string {
	var networks []string
	for _, network := range networkList {
		if _, exists := c.NetworkSettings.Networks[network]; exists {
			networks = append(networks, network)
		}
	}
	return networks
}

func getNetworkContainers(cli *client.Client, network string) ([]types.Container, error) {
	filter := filters.NewArgs()
	filter.Add("network", network)

	return cli.ContainerList(
		context.Background(),
		container.ListOptions{
			All:     true,
			Filters: filter,
		},
	)
}

func getContainerStatus(cli *client.Client, containerID string) (string, error) {
	info, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return "", err
	}
	return info.State.Status, nil
}

func main() {
	ParseEnv()

	for {
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			logger.Fatal("Docker client init error:", err)
		}

		if _, err := cli.Ping(context.Background()); err != nil {
			logger.Fatalf("Docker API connection error: %v", err)
		}
		logger.Println("Successfully connected to Docker API")
		checkContainers(cli)
		time.Sleep(10 * time.Second)
		cli.Close()
	}
}
