package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
)

func MemAvail() string {
	fname := "/proc/meminfo"
	FileBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Printf("!Err-> can't read %s\n", fname)
		return "error"
	}
	bufr := bytes.NewBuffer(FileBytes)
	for {
		line, err := bufr.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Sprintf("!Err-> can't read %s\n", fname)
		}
		ndx := strings.Index(line, "MemFree:")
		if ndx >= 0 {
			line = strings.TrimSpace(line[9:])
			fmt.Printf("%q\n", line)
			line = line[:len(line)-3]
			fmt.Printf("%q\n", line)
			mem, err := strconv.ParseInt(line, 10, 64)
			if err == nil {
				return line
			}
			// some problem if parse failed
			fmt.Printf("line: %s\n", line)
			n, err := fmt.Sscan(line, "%d", &mem)
			if err != nil {
				fmt.Printf("!Err-> can't scan %s\n", line)
				return "error"
			}
			if n != 1 {
				fmt.Printf("!Err-> can't scan all %s\n", line)
				return "error"
			}
			return "error"
		}
	}
	fmt.Printf("didn't find MemFree in /proc/meminfo\n")
	return "error"
}

func FatalError(err error) {
	fmt.Printf("!Err-> %s\n", err)
	os.Exit(1)
}

func main() {
	fmt.Println("Go RabbitMQ")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Succesfully Connected to our RabbitMQ Instance")

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	// defer ch.Close()

	forever := make(chan bool)
	fmt.Println("Looping...")
	go func() {
		q, err := ch.QueueDeclare(
			"TestQueue",
			false,
			false,
			false,
			false,
			nil,
		)

		if err != nil {
			fmt.Println(err)
			panic(err)
		}

		fmt.Println(q)

		var mm = MemAvail()

		err = ch.Publish(
			"",
			"TestQueue",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(mm),
			},
		)

		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}()

	fmt.Println("Succesfully Published Message to Queue")
	<-forever

}
