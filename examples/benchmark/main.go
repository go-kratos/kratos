package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	http2 "net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
)

var (
	port              = 8080
	sleepTime         = 0
	cpuBound          bool
	target            = 15
	sleepTimeDuration time.Duration
	message           = []byte("hello world")
	samplingPoint     = 20 // seconds
)

func main() {
	args := os.Args
	argsLen := len(args)
	webFramework := "kratos"
	if argsLen > 1 {
		webFramework = args[1]
	}
	if argsLen > 2 {
		sleepTime, _ = strconv.Atoi(args[2])
		if sleepTime == -1 {
			cpuBound = true
			sleepTime = 0
		}
	}
	if argsLen > 3 {
		port, _ = strconv.Atoi(args[3])
	}
	if argsLen > 4 {
		samplingPoint, _ = strconv.Atoi(args[4])
	}
	sleepTimeDuration = time.Duration(sleepTime) * time.Millisecond
	samplingPointDuration := time.Duration(samplingPoint) * time.Second

	go func() {
		time.Sleep(samplingPointDuration)
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		var u uint64 = 1024 * 1024
		fmt.Printf("TotalAlloc: %d\n", mem.TotalAlloc/u)
		fmt.Printf("Alloc: %d\n", mem.Alloc/u)
		fmt.Printf("HeapAlloc: %d\n", mem.HeapAlloc/u)
		fmt.Printf("HeapSys: %d\n", mem.HeapSys/u)
	}()

	switch webFramework {
	case "kratos":
		KratosServer()
	case "gin":
		GinServer()
	case "echo":
		EchoServer()
	case "mux":
		MuxServer()
	default:
		fmt.Println("--------------------------------------------------------------------")
		fmt.Println("------------- Unknown framework given!!! Check libs.sh -------------")
		fmt.Println("------------- Unknown framework given!!! Check libs.sh -------------")
		fmt.Println("------------- Unknown framework given!!! Check libs.sh -------------")
		fmt.Println("--------------------------------------------------------------------")
	}
}

func KratosServer() {
	httpSrv := http.NewServer(
		http.Address(":" + strconv.Itoa(port)),
	)
	httpSrv.HandleFunc("/", kratosHandle)
	app := kratos.New(
		kratos.Name("benchmark"),
		kratos.Server(
			httpSrv,
		),
	)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func kratosHandle(w http2.ResponseWriter, q *http2.Request) {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	w.Write(message)
}

func GinServer() {
	gin.SetMode(gin.ReleaseMode)
	mux := gin.New()
	mux.GET("/", ginHandler)
	mux.Run(":" + strconv.Itoa(port))
}

func ginHandler(c *gin.Context) {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.Writer.Write(message)
}

func EchoServer() {
	e := echo.New()
	e.GET("/", echoHandler)
	e.Start(":" + strconv.Itoa(port))
}

func echoHandler(c echo.Context) error {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	c.Response().Write(message)
	return nil
}

func MuxServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", muxHandle)
	http2.Handle("/", r)
	http2.ListenAndServe(":"+strconv.Itoa(port), r)
}

func muxHandle(w http2.ResponseWriter, q *http2.Request) {
	if cpuBound {
		pow(target)
	} else {
		if sleepTime > 0 {
			time.Sleep(sleepTimeDuration)
		} else {
			runtime.Gosched()
		}
	}
	w.Write(message)
}

func pow(targetBits int) {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for {
		data := "hello world " + strconv.Itoa(nonce)
		hash = sha256.Sum256([]byte(data))
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}

		if nonce%100 == 0 {
			runtime.Gosched()
		}
	}
}
