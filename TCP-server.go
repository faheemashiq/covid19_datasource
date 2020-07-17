package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	//"bufio"
	//"os"
	//dataSource "data"
)

//DataStructure

type Data struct {
	Date                       string
	Cumulative_Test_positive   string
	Cumulative_tests_performed string
	Expired                    string
	Still_admitted             string
	Discharged                 string
	Region                     string
}

func Load(path string) []Data {
	table := make([]Data, 0)
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err.Error())
		}
		data := Data{
			Cumulative_Test_positive:   row[0],
			Cumulative_tests_performed: row[1],
			Date:                       row[2],
			Discharged:                 row[3],
			Expired:                    row[4],
			Region:                     row[5],
			Still_admitted:             row[6],
		}
		table = append(table, data)
	}
	return table
}

func Find(table []Data, filter string) []Data {
	result := make([]Data, 0)
	filter = strings.ToUpper(filter)
	for _, dat := range table {
		if strings.Contains(strings.ToUpper(dat.Region), filter) ||
			strings.Contains(strings.ToUpper(dat.Date), filter) {

			result = append(result, dat)
		}

	}
	return result
}

var myData = Load(`C:\Go\Workspace\data.csv`)

func main() {
	var addr string
	var network string
	flag.StringVar(&addr, "e", ":4040", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.Parse()

	// create a listener for provided network and host address
	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Fatal("failed to create listener:", err)
	}
	defer ln.Close()
	log.Println("**** Covid19 Data of Pakistan ***")
	log.Printf("Service started: (%s) %s\n", network, addr)

	// connection-loop - handle incoming requests
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			if err := conn.Close(); err != nil {
				log.Println("failed to close listener:", err)
			}
			continue
		}
		log.Println("Connected to", conn.RemoteAddr())

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	if _, err := conn.Write([]byte("Connected...\nUsage: GET <region, date>\n")); err != nil {
		log.Println("error writing:", err)
		return
	}

	// loop to stay connected with client until client breaks connection
	for {
		// buffer for client command
		cmdLine := make([]byte, (1024 * 4))
		n, err := conn.Read(cmdLine)
		if n == 0 || err != nil {
			log.Println("connection read error:", err)
			return
		}

		cmd, param := parseCommand(string(cmdLine[0:n]))
		//fmt.Println(cmd, param)
		if cmd == "" {
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
				log.Println("failed to write:", err)
				return
			}
			continue
		}
		//log.Println(param)
		//execute command
		switch strings.ToUpper(cmd) {
		case "GET":
			result := Find(myData, param)
			//fmt.Print(result[1])
			if len(result) == 0 {
				if _, err := conn.Write([]byte("Nothing found\n")); err != nil {
					log.Println("failed to write:", err)
				}
				continue
			}
			structArr := make([]Data, 0)
			for _, cur := range result {
				var dataObj Data
				dataObj.Date = cur.Date
				dataObj.Cumulative_Test_positive = cur.Cumulative_Test_positive
				dataObj.Cumulative_tests_performed = cur.Cumulative_tests_performed
				dataObj.Expired = cur.Expired
				dataObj.Still_admitted = cur.Still_admitted
				dataObj.Discharged = cur.Discharged
				dataObj.Region = cur.Region
				structArr = append(structArr, dataObj)
			}

			data, _ := json.Marshal(structArr)

			_, err := conn.Write([]byte(
				fmt.Sprintln(
					string(data)),
			))
			if err != nil {
				log.Println("failed to write response:", err)
				return
			}

		default:
			if _, err := conn.Write([]byte("Invalid command\n")); err != nil {
				log.Println("failed to write:", err)
				return
			}
		}
		fmt.Println()
	}
}

func parseCommand(cmdLine string) (cmd, param string) {

	parts := strings.Split(cmdLine, " ")
	if len(parts) != 2 {
		return "", ""
	}
	cmd = strings.TrimSpace(parts[0])
	param = strings.TrimSpace(parts[1])

	return
}
