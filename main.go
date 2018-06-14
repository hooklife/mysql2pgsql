package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/xwb1989/sqlparser"
	"os"

)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	sql := "select * from `hj_banner` where `status` = '1' and `location` = '1'"
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		print(err.Error())
	}

	log.SetFormatter(&log.TextFormatter{})
	// Otherwise do something with stmt
	switch stmt.(type) {
	case *sqlparser.Select:
		log.Debug("AAAA");
	case *sqlparser.Insert:
		print(22222)
	}
	//print(stmt)
	os.Exit(1)

	conn, err := amqp.Dial("amqp://rabbitmq:rabbitmq@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
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
			log.Printf("Received a message: %s", d.Body)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
