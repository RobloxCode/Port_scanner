package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	err := os.MkdirAll("ScanResults", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	fileName := fmt.Sprintf("%v.txt", GetCurrentDateAndHour())
	file, err := CreateFile("ScanResults/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ip := "192.168.1.75"
	protocol := "tcp"
	var wg sync.WaitGroup

	results, err := ScanRangePortsConcurrently(&wg, protocol, ip, 0, 101)
	if err != nil {
		log.Fatal(err)
	}

	for port, open := range results {
		_, err := file.WriteString(fmt.Sprintf("port: %v, open: %v\n", port, open))
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Scan results saved to file: ", fileName)
	DisplayClosedPorts(ip, results)
	DisplayOpenPorts(ip, results)
}

func DisplayOpenPorts(ip string, ports map[string]bool) {
	fmt.Printf("The ip: %v has the following open ports:\n", ip)
	for port, open := range ports {
		if open == true {
			fmt.Printf("port: %v, open: %v\n", port, open)
		}
	}
}

func DisplayClosedPorts(ip string, ports map[string]bool) {
	fmt.Printf("The ip: %v has the following close ports:\n", ip)
	for port, open := range ports {
		if open == false {
			fmt.Printf("port: %v, open: %v\n", port, open)
		}
	}
}

func ScanRangePorts(protocol, ip string, bottomPort, topPort int) (map[string]bool, error) {
	results := make(map[string]bool)

	if bottomPort > topPort {
		return results, fmt.Errorf("last port: %v must be greater than start port: %v\n", topPort, bottomPort)
	}

	for i := bottomPort; i < topPort; i++ {
		port := strconv.Itoa(i)
		conn, err := IPConnectedToPort(protocol, ip, port)

		if err != nil {
			results[port] = false
		} else {
			results[port] = true
			conn.Close()
		}
	}
	return results, nil
}

func IPConnectedToPort(protocol, ip, port string) (net.Conn, error) {
	address := ip + ":" + port
	return net.DialTimeout(protocol, address, 1*time.Second)
}

func CreateFile(name string) (*os.File, error) {
	file, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetCurrentDateAndHour() string {
	now := time.Now()
	return now.Format("2006-01-02_15h04m05s")
}

func ScanRangePortsConcurrently(wg *sync.WaitGroup, protocol, ip string, bottomPort, topPort int) (map[string]bool, error) {
	results := make(map[string]bool)
	var mu sync.Mutex

	if bottomPort > topPort {
		return results, fmt.Errorf("last port: %v must be greater than start port: %v\n", topPort, bottomPort)
	}

	for i := bottomPort; i < topPort; i++ {
		wg.Add(1)
		port := strconv.Itoa(i)

		go func(port string) {
			defer wg.Done()
			conn, err := IPConnectedToPort(protocol, ip, port)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				results[port] = false
			} else {
				results[port] = true
				conn.Close()
			}
		}(port)
	}

	wg.Wait()
	return results, nil
}
