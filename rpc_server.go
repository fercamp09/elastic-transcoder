package main

import (
	"log"
	"io/ioutil"

	"github.com/streadway/amqp"
	"github.com/quirkey/magick"
	"github.com/golang/protobuf/proto"
	pb "github.com/fercamp09/elastic-transcoder/tasks"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Setting two priority levels
	args := make(amqp.Table)
	args["x-max-priority"] = int32(2)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"rpc_queue2", // name
		false,       // durable
		false,       // delete when usused
		false,       // exclusive
		false,       // no-wait
		args,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Decode task
			task := &pb.Task{}
			err := proto.Unmarshal(d.Body, task)
			failOnError(err, "Failed to parse task")
			
			// Process image
			source, _ := ioutil.ReadFile(task.Filename)
			image, err := magick.NewFromBlob(source, "jpg")
			failOnError(err, "Error reading from file")
			response := task.NewName
			err = image.ToFile(response)
			failOnError(err, "Problem with writing") 
			log.Printf(" [.] image (%s)", response)
			image.Destroy()

			// Encode response
			res := &pb.Response{
				FileLocation: response,
			}
			log.Printf(res.FileLocation)
			body, err := proto.Marshal(res)
			failOnError(err, "Failed to encode response")
			
			err = ch.Publish(
				"",        // exchange
				d.ReplyTo, // routing key
				false,     // mandatory
				false,     // immediate
				amqp.Publishing{
					ContentType:   "application/octet-stream",
					CorrelationId: d.CorrelationId,
					Body:          body,
				})
			failOnError(err, "Failed to publish a message")

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}
