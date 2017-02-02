package main

import (
	"log"
	"io/ioutil"
	"fmt"
	"net/http"
	"io"
	"os"
	"mime/multipart"
	"bytes"
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

func Upload(url, file string) (err error) {
    // Prepare a form that you will submit to that URL.
    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    // Add your image file
    f, err := os.Open(file)
    if err != nil {
        return
    }
    defer f.Close()
    fw, err := w.CreateFormFile("image", file)
    if err != nil {
        return
    }
    if _, err = io.Copy(fw, f); err != nil {
        return
    }
    // Add the other fields
    if fw, err = w.CreateFormField("key"); err != nil {
        return
    }
    if _, err = fw.Write([]byte("KEY")); err != nil {
        return
    }
    // Don't forget to close the multipart writer.
    // If you don't close it, your request will be missing the terminating boundary.
    w.Close()

    // Now that you have a form, you can submit it to your handler.
    req, err := http.NewRequest("PUT", url, &b)
    if err != nil {
        return
    }
    // Don't forget to set the content type, this will contain the boundary.
    req.Header.Set("Content-Type", w.FormDataContentType())

    // Submit the request
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        return
    }
    defer res.Body.Close()
    // Check the response
    if res.StatusCode != http.StatusOK {
        err = fmt.Errorf("bad status: %s", res.Status)
    }
    return
}


func main() {
	var task_cancelled bool
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
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		args,         // arguments
	)

	cq, err := ch.QueueDeclare(
		"cancel_queue",  //name
		false,           //durable
		false,           //delete when unused
		false,           //exclusive
		false,           //no-wait
		args,            // arguments

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

	cmsgs, err := ch.Consume(
		cq.Name, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// Decode task
			task := &pb.Task{}
			err := proto.Unmarshal(d.Body, task)
			failOnError(err, "Failed to parse task")


			for t := range cmsgs {
				cancel := &pb.Cancel{}
				err := proto.Unmarshal(t.Body, cancel)
				failOnError(err, "Failed to parse cancel message")
				if cancel.FileId == task.FileId {
					task_cancelled = true
					t.Ack(false)
					break
				}else{
					t.Nack(false, true)
					break
				}
			}

			if task_cancelled {
				d.Ack(false)
				break
			}

			// Process image
			url := "http://localhost:3000/files/" + task.FileId //se ingresa el id con el que se descargara el archivo
			resp, err := http.Get(url)
		  defer resp.Body.Close()
		  out, err := os.Create(task.Filename)
		  if err != nil {
		    // panic?
		  }
		  defer out.Close()
		  io.Copy(out, resp.Body)



			source, _ := ioutil.ReadFile(task.Filename)
			image, err := magick.NewFromBlob(source, "jpg")
			failOnError(err, "Error reading from file")
			response := task.NewName
			err = image.ToFile(response)
			failOnError(err, "Problem with writing")
			log.Printf(" [.] image (%s)", response)
			image.Destroy()


			Upload(url, task.NewName) //aqui se sube el archivo que se convirtió

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
