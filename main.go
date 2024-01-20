package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var version = "1.0.2"

func main() {
	var targetIP, port, password, tlsVersion, certFile, certKey, caCertFile string
	var useTLS, tlsInsecureSkipVerify, showHelp, showVersion, setCommand, parallel, hello, outputFile bool
	var numConnections int

	flag.StringVar(&targetIP, "ip", "localhost", "Redis server IP address")
	flag.StringVar(&port, "port", "6379", "Redis server port")
	flag.StringVar(&password, "password", "", "Redis server password")
	flag.StringVar(&tlsVersion, "tlsVersion", "1.2", "TLS version (1.2 or 1.3)")
	flag.StringVar(&certFile, "certFile", "", "Path to client certificate")
	flag.StringVar(&certKey, "certKey", "", "Path to client private key")
	flag.StringVar(&caCertFile, "caCertFile", "", "Path a file containing CA certificates")
	flag.BoolVar(&useTLS, "tls", false, "Use TLS for connection")
	flag.BoolVar(&tlsInsecureSkipVerify, "tlsSkipVerify", false, "Skip verification of server certificate")
	flag.BoolVar(&setCommand, "setCommand", false, "Send additional SET command for every connection")
	flag.BoolVar(&showHelp, "help", false, "Display usage")
	flag.BoolVar(&parallel, "parallel", false, "Run connections in parallel")
	flag.BoolVar(&showVersion, "version", false, "Display version")
	flag.BoolVar(&hello, "hello", false, "Send Hello command after connection")
	flag.BoolVar(&outputFile, "outputFile", false, "Send per connection timing to a file, provide a file name")
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

	if flag.Lookup("tls") != nil && flag.Lookup("tls").DefValue != flag.Lookup("tls").Value.String() {
		useTLS = true
	}

	if useTLS {
		establishTLSConnections(redisAddress, password, tlsVersion, certFile, certKey, caCertFile, tlsInsecureSkipVerify, numConnections, parallel, hello, outputFile)
	} else {
		establishUnencryptedConnections(redisAddress, password, numConnections, parallel, hello, outputFile)
	}
}

func establishTLSConnections(redisAddress, password, tlsVersion, certFile, certKey, caCertFile string, InsecureSkipVerify bool, numConnections int, parallel, hello, outputFile bool) {
	fmt.Println("Using TLS for connection")

	tlsConfig, tlsConfigError := createTLSConfig(tlsVersion, certFile, certKey, caCertFile, InsecureSkipVerify)

	if tlsConfigError != nil {
		fmt.Println("error in TLS config: %v", tlsConfigError)
	}

	if parallel {
		testAndMeasureConnectionsParallel(redisAddress, password, tlsConfig, numConnections, hello, outputFile)
	} else {
		testAndMeasureConnections(redisAddress, password, tlsConfig, numConnections, hello, outputFile)
	}
}

func establishUnencryptedConnections(redisAddress, password string, numConnections int, parallel, hello, outputFile bool) {
	fmt.Println("Using unencrypted connection")

	if parallel {
		testAndMeasureConnectionsParallel(redisAddress, password, nil, numConnections, hello, outputFile)
	} else {
		testAndMeasureConnections(redisAddress, password, nil, numConnections, hello, outputFile)
	}
}

func createTLSConfig(tlsVersion, certFile, certKey, caCertFile string, tlsInsecureSkipVerify bool) (*tls.Config, error) {
	tlsConfig := &tls.Config{}

	// Load client certificate and private key if provided
	if certFile != "" && certKey != "" {
		cert, err := tls.LoadX509KeyPair(certFile, certKey)
		if err != nil {
			return nil, fmt.Errorf("error loading client certificate: %v", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// Load CA certificates if provided
	if caCertFile != "" {
		caCert, err := os.ReadFile(caCertFile)
		if err != nil {
			return nil, fmt.Errorf("error reading CA certificate file: %v", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	if tlsVersion != "" {
		// Set TLS version based on the specified parameter
		switch tlsVersion {
		case "1.2":
			tlsConfig.MaxVersion = tls.VersionTLS12
		case "1.3":
			tlsConfig.MaxVersion = tls.VersionTLS13
		default:
			return nil, fmt.Errorf("invalid TLS version specified")
		}
	}

	// InsecureSkipVerify controls whether a client verifies the
	// server's certificate chain and host name.
	// If InsecureSkipVerify is true, TLS accepts any certificate
	// presented by the server and any host name in that certificate.
	// In this mode, TLS is susceptible to man-in-the-middle attacks.
	// This should be used only for testing.
	tlsConfig.InsecureSkipVerify = tlsInsecureSkipVerify

	return tlsConfig, nil
}

func testAndMeasureConnections(redisAddress, password string, tlsConfig *tls.Config, numConnections int, hello, outputFile bool) {
	startTime := time.Now()
	var totalConnectionTime float32
	var file *os.File

	for i := 0; i < numConnections; i++ {
		connStartTime := time.Now()

		conn, err := redis.Dial("tcp", redisAddress, redis.DialPassword(password))

		if tlsConfig != nil {
			conn, err = redis.Dial(
				"tcp",
				redisAddress,
				redis.DialPassword(password),
				redis.DialTLSConfig(tlsConfig),
				redis.DialUseTLS(true),
				redis.DialTLSSkipVerify(tlsConfig.InsecureSkipVerify),
				redis.DialKeepAlive(time.Duration(0)))
		}

		if err != nil {
			fmt.Printf("Failed to connect to Redis %s:%s: %v\n", redisAddress, password, err)
			return
		}
		connElapsedTime := time.Since(connStartTime)
		if outputFile {
			file, err = os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			_, err = fmt.Fprintf(file, "%.3f\n", float32(connElapsedTime.Milliseconds()))
			if err != nil {
				fmt.Println("Error writing to the file:", err)
				return
			}
			defer file.Close()
		}

		if hello {
			_, err = conn.Do("HELLO")
			duration := 10 * time.Second
			time.Sleep(duration)
			if err != nil {
				fmt.Println("Failed to execute HELLO command:", err)
				conn.Close()
				return
			}
		}

		defer conn.Close()

		totalConnectionTime += float32(connElapsedTime.Milliseconds())
	}

	elapsedTime := time.Since(startTime)
	avgConnectionTime := totalConnectionTime / float32(numConnections)

	if tlsConfig != nil {
		fmt.Printf("Established %d TLS connections in %v\nAverage connection time: %.3fms\n", numConnections, elapsedTime, avgConnectionTime)
	} else {
		fmt.Printf("Established %d unencrypted connections in %v\nAverage connection time:  %.3fms\n", numConnections, elapsedTime, avgConnectionTime)
	}
}

func testAndMeasureConnectionsParallel(redisAddress, password string, tlsConfig *tls.Config, numConnections int, hello, outputFile bool) {
	var wg sync.WaitGroup
	startTime := time.Now()
	var totalConnectionTime float32
	var file *os.File

	for i := 0; i < numConnections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			connStartTime := time.Now()

			conn, err := redis.Dial("tcp", redisAddress, redis.DialPassword(password))

			if tlsConfig != nil {
				conn, err = redis.Dial("tcp", redisAddress, redis.DialPassword(password), redis.DialTLSConfig(tlsConfig), redis.DialKeepAlive(time.Duration(0)))
			}

			if err != nil {
				fmt.Printf("Failed to connect to Redis %s:%s: %v\n", redisAddress, password, err)
				return
			}
			connElapsedTime := time.Since(connStartTime)
			if outputFile {
				file, err = os.OpenFile("output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

				if err != nil {
					fmt.Println("Error creating file:", err)
					return
				}
				_, err = fmt.Fprintf(file, "%.3f\n", float32(connElapsedTime.Milliseconds()))
				if err != nil {
					fmt.Println("Error writing to the file:", err)
					return
				}
				defer file.Close()
			}

			if hello {
				_, err = conn.Do("HELLO")
				duration := 10 * time.Second
				time.Sleep(duration)
				if err != nil {
					fmt.Println("Failed to execute HELLO command:", err)
					conn.Close()
					return
				}
			}

			defer conn.Close()
			totalConnectionTime += float32(connElapsedTime.Milliseconds())
		}()
	}

	wg.Wait()

	elapsedTime := time.Since(startTime)
	avgConnectionTime := totalConnectionTime / float32(numConnections)

	if tlsConfig != nil {
		fmt.Printf("Established %d TLS connections in %v\nAverage connection time: %.3fms\n", numConnections, elapsedTime, avgConnectionTime)
	} else {
		fmt.Printf("Established %d unencrypted connections in %v\nAverage connection time: %.3fms\n", numConnections, elapsedTime, avgConnectionTime)
	}
}
