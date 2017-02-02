package main

import (
	"log"
	"math/rand"
	"os"
	"time"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
	"github.com/golang/protobuf/proto"
	pb "github.com/fercamp09/elastic-transcoder/tasks"
)


func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func connectToRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	return conn, err
}

func imageRPC(i string, o string, p int) (resp *pb.Response, err error) {
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	//failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()
	conn, err := connectToRabbitMQ()
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	args := make(amqp.Table)
	args["x-max-priority"] = int32(2)

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when usused
		true,  // exclusive
		false, // noWait
		args,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	corrId := randomString(32)
	log.Printf("ID: ", corrId)
 
        // Encode with protocol buffer
	//t := &pb.Task{
        //	Filename:  i,
       // 	NewName:  o,
       // 	Priority:  int32(p),
	//	FileId: corrId,
	//}
	//out, err := proto.Marshal(t)
       	//failOnError(err, "Failed to encode task:")
	
	// Publish task
	//err = ch.Publish(
	//	"",          // exchange
	//	"rpc_queue2",// routing key
	//	false,       // mandatory
	//	false,       // immediate
	//	amqp.Publishing{
	//		ContentType:   "application/octet-stream",
	//		CorrelationId: corrId,
	//		ReplyTo:       q.Name,
	//		Priority:      uint8(p),
	//		Body:          out,
	//	})
	//failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			resp = &pb.Response{}
			err := proto.Unmarshal(d.Body, resp)
			log.Printf(resp.FileLocation)
			failOnError(err, "Failed to convert body to string")
			break
		}
	}

	return
}

func cancelTask(id string){
	conn, err := connectToRabbitMQ()

        ch, err := conn.Channel()
        failOnError(err, "Failed to open a channel")
        defer ch.Close()

        args := make(amqp.Table)
        args["x-max-priority"] = int32(2)

        q, err := ch.QueueDeclare(
                "cancel_queue",    // name
                false, // durable
                false, // delete when usused
                true,  // exclusive
                false, // noWait
                args,   // arguments
        )
	failOnError(err, "Failed to declare a queue")

        //corrId := randomString(32)
	corrId := id
        
	// Encode with protocol buffer
        t := &pb.Cancel{
		FileId:  corrId,
        }
        out, err := proto.Marshal(t)
        failOnError(err, "Failed to encode task:")

        // Publish task
        err = ch.Publish(
                "",          // exchange
                q.Name,      // routing key
                false,       // mandatory
                false,       // immediate
                amqp.Publishing{
                        ContentType:   "application/octet-stream",
                        Body:          out,
                })
        failOnError(err, "Failed to publish a message")

        log.Printf("Task %s cancelled", corrId)

        return
}

func readTask(id string){
	conn, err := connectToRabbitMQ()
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	args := make(amqp.Table)
        args["x-max-priority"] = int32(2)
	
	q, err := ch.QueueDeclare(
		"rpc_queue2", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		args,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range msgs {
		task := &pb.Task{}
                err := proto.Unmarshal(d.Body, task)
                failOnError(err, "Failed to parse task message")
                if id == task.FileId {
                      log.Printf("FileId: %s, Filename:%s, New Filename: %s, New Format: %s", task.FileId, task.Filename, task.NewName, task.Priority, task.Format)
                      break
		}
	}

	log.Printf(" [*] Searching for message. To exit press CTRL+C")
	
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if strings.Compare(os.Args[1],"create") == 0 {
		input, output, filetype, priority := bodyFrom(os.Args)
		log.Printf(" [x] Requesting image(%s)", input)
		res, err := imageRPC(input, output, priority)
		failOnError(err, "Failed to handle RPC request")
		log.Printf(" [.] Image processed found in %s, %s", res.FileLocation, filetype)
	} else if strings.Compare(os.Args[1], "read") == 0 {
		id := os.Args[2]
		readTask(id)
	} else if os.Args[1] == "cancel" {
		id := os.Args[2]
		cancelTask(id)
	} else {
		log.Printf("Wrong arguments, valid options: read, cancel, create")
	}
	
}

func bodyFrom(args []string) (string, string, string, int) {
	var i, o, s string
	var p int
	var err error

	if len(os.Args) < 5 {
		log.Print("Please supply an input and output filename e.g. go run rpc_client.go input.jpg output.jpg jpg 1")
		os.Exit(3)
		} else if len(os.Args) == 5 {
 		p = 0
		i = os.Args[2]
		o = os.Args[3]
		s = os.Args[4]
		log.Print("Added to Non-Priority Queue")
	} else {
		i = os.Args[2]
		o = os.Args[3]
		s = os.Args[4]
		p, err = strconv.Atoi(os.Args[5])
	}
	failOnError(err, "Wrong arguments")
	return i, o, s, p
}
