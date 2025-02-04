package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

// pingIP выполняет ping для заданного IP и выводит результат.
func pingIP(ip string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		log.Printf("Ошибка создания пингера для %s: %v\n", ip, err)
		return
	}
	
	// Устанавливаем режим привилегированного пинга (может потребоваться root-права)
	pinger.SetPrivileged(true)
	pinger.Count = 3               // Количество ICMP-запросов
	pinger.Timeout = 3 * time.Second // Общий таймаут пинга
	
	
	err = pinger.Run() // Запуск пинга
	if err != nil {
		log.Printf("Ошибка при пинге %s: %v\n", ip, err)
		return
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		fmt.Printf("Timestamp: %s.\n IP %s доступен (%d/%d пакетов получено)\n", time.Now().Format(time.RFC1123), ip, stats.PacketsRecv, stats.PacketsSent)
	} else {
		fmt.Printf("Timestamp: %s.\n IP %s недоступен\n", time.Now().Format(time.RFC1123), ip)
	}
}

func main() {
	var Ips []string
	file, err := os.Open("ips.json")
	if err != nil {
		log.Fatal(err)
	}
	byteSl, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(byteSl, &Ips)
	if err != nil {
		log.Fatal(err)
}

	var wg sync.WaitGroup

	// Запускаем пинг каждой IP-адреса в отдельной горутине
	for _, ip := range Ips {
		wg.Add(1)
		go pingIP(ip, &wg)
	}

	wg.Wait()
}