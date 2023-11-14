package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var version = "1.0.0"

func main() {
	var targetIP, port, password, tlsVersion, certFile, certKey string
	var useTLS, showHelp, showVersion, setCommand bool
	var numConnections int

	flag.StringVar(&targetIP, "ip", "your-redis-host", "Redis server IP address")
	flag.StringVar(&port, "port", "6379", "Redis server port")
	flag.StringVar(&password, "password", "", "Redis server password")
	flag.StringVar(&tlsVersion, "tlsVersion", "1.2", "TLS version (1.2 or 1.3)")
	flag.StringVar(&certFile, "certFile", "", "Path to client certificate")
	flag.StringVar(&certKey, "certKey", "", "Path to client private key")
	flag.BoolVar(&useTLS, "tls", false, "Use TLS for connection")
	flag.BoolVar(&setCommand, "setCommand", false, "Send additional SET command for every connection")
	flag.BoolVar(&showHelp, "help", false, "Display usage")
	flag.BoolVar(&showVersion, "version", false, "Display version")
	flag.IntVar(&numConnections, "numConnections", 100, "Number of connections to establish")

	flag.Parse()

	if showHelp {
		flag.Usage()
		return
	}

	if showVersion {
		fmt.Printf("Version: %s\n", version)
		return
	}

	redisAddress := fmt.Sprintf("%s:%s", targetIP, port)

	if flag.NFlag() == 0 {
		fmt.Println("Please provide at least one flag")
		flag.PrintDefaults()
		return
	}

	if flag.Lookup("tls") != nil && flag.Lookup("tls").DefValue != flag.Lookup("tls").Value.String() {
		useTLS = true
	}

	if useTLS {
		// Use TLS connection
		fmt.Println("Using TLS for connection")

		// Create a TLS configuration
		tlsConfig := &tls.Config{}

		// Load the client certificate and key if provided
		if certFile != "" && certKey != "" {
			cert, err := tls.LoadX509KeyPair(certFile, certKey)
			if err != nil {
				fmt.Println("Error loading client certificate:", err)
				return
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		// Choose TLS version
		switch tlsVersion {
		case "1.2":
			tlsConfig.MaxVersion = tls.VersionTLS12
		case "1.3":
			tlsConfig.MaxVersion = tls.VersionTLS13
		default:
			fmt.Println("Invalid TLS version specified")
			return
		}

		// Test dial the Redis server with TLS
		conn, err := redis.Dial("tcp", redisAddress, redis.DialPassword(password), redis.DialTLSConfig(tlsConfig))
		if err != nil {
			fmt.Println("Failed to connect to Redis with TLS:", err)
			return
		}
		defer conn.Close()

		// Measure connection rate for encrypted connection
		startTime := time.Now()

		var totalConnectionTime int64

		for i := 0; i < numConnections; i++ {
			connStartTime := time.Now()

			conn, err := redis.Dial("tcp", redisAddress, redis.DialPassword(password), redis.DialTLSConfig(tlsConfig))
			if err != nil {
				fmt.Println("Failed to connect to Redis with TLS:", err)
				return
			}

			if setCommand {
				_, err := conn.Do("SET", "test_key", "test_value")
				if err != nil {
					fmt.Println("Failed to execute SET command:", err)
					conn.Close()  // Close the connection immediately if SET command fails
					return
				}
			}
		
			defer conn.Close()

			connElapsedTime := time.Since(connStartTime)
			totalConnectionTime += connElapsedTime.Microseconds()
		}

		elapsedTime := time.Since(startTime)
		avgConnectionTime := float64(totalConnectionTime) / float64(numConnections) / 1000
		fmt.Printf("Established %d TLS connections in %v\nAverage connection time: %.3f milliseconds\n", numConnections, elapsedTime, avgConnectionTime)
	} else {
		// Use unencrypted connection
		fmt.Println("Using unencrypted connection")

		// Test dial the Redis server without TLS
		conn, err := redis.Dial("tcp", redisAddress, redis.DialPassword(password))
		if err != nil {
			fmt.Println("Failed to connect to Redis without TLS:", err)
			return
		}
		
		if setCommand {
			_, err := conn.Do("SET", "test_key", "test_value")
			if err != nil {
				fmt.Println("Failed to execute SET command:", err)
				conn.Close()  // Close the connection immediately if SET command fails
				return
			}
		}
	
		defer conn.Close()

		// Measure connection rate for unencrypted connection
		startTime := time.Now()

		var totalConnectionTime int64

		for i := 0; i < numConnections; i++ {
			connStartTime := time.Now()

			conn, err := redis.Dial("tcp", redisAddress, redis.DialPassword(password))
			if err != nil {
				fmt.Println("Failed to connect to Redis without TLS:", err)
				return
			}
			conn.Close()

			connElapsedTime := time.Since(connStartTime)
			totalConnectionTime += connElapsedTime.Microseconds()
		}

		elapsedTime := time.Since(startTime)
		avgConnectionTime := float64(totalConnectionTime) / float64(numConnections) / 1000
		fmt.Printf("Established %d unencrypted connections in %v\nAverage connection time: %.3f milliseconds\n", numConnections, elapsedTime, avgConnectionTime)
	}
}
