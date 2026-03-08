package internal

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

const (
	WIKIMEDIA_URL = "https://stream.wikimedia.org/v2/stream/recentchange"
)

func WikimediaProduceKafka(producer sarama.SyncProducer, topic string) error {

	// Create a new HTTP request
	req, err := http.NewRequest("GET", WIKIMEDIA_URL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return err
	}

	// Set a custom User-Agent header
	req.Header.Set("User-Agent", "wikimedia-pet-project/1.0")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	timerChan := time.After(time.Minute * 1)

	fmt.Println("scanning producer")
Loop:
	for {
		select {
		case <-timerChan:
			fmt.Println("Timer ended, finishing loop")
			break Loop
		default:
			if scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "data:") {
					line = line[6:]
					PushMessageToQueue(producer, topic, []byte(line))
				}
			} else {
				break Loop
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading stream: %v", err)
		return err
	}

	return nil
}
