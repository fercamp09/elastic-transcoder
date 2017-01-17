package main

import (
	"log"
	"math/rand"
	"os"
	"time"

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

func imageRPC(i string, o string, p int) (resp *pb.Response, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

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

        // Encode with protocol buffer
	t := &pb.Task{
        	Filename:  i,
        	NewName:  o,
        	Priority:  int32(p),
	}
	out, err := proto.Marshal(t)
       	failOnError(err, "Failed to encode task:")
	
	// Publish task
	err = ch.Publish(
		"",          // exchange
		"rpc_queue2",// routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "application/octet-stream",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Priority:      uint8(p),
			Body:          out,
		})
	failOnError(err, "Failed to publish a message")

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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	input, output, priority := bodyFrom(os.Args)

	log.Printf(" [x] Requesting image(%s)", input)
	res, err := imageRPC(input, output, priority)
	failOnError(err, "Failed to handle RPC request")

	log.Printf(" [.] Image processed found in %s", res.FileLocation)
}

func bodyFrom(args []string) (string, string, int) {
	//var i, o, s string
	//var p int
	//if len(os.Args) < 3 {
	//	log.Print("Please supply an input and output filename e.g. go run rpc_client.go input.jpg output.jpg")	
	//	p = 0
	//} else if len(os.Args) == 3 {
 	//	p = 0
	//} else {
	//	i = os.Args[1]
	//	o = os.Args[2]
	//	s = os.Args[3]
	//	p = 0
	//}
	i := "bvb.png"
	o := "bvb.jpg"
	p := 0  
	return i, o, p
}
