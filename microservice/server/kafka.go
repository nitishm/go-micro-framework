package server

import (
	"encoding/json"
	"go-micro-framework/microservice/configuration"
	"go-micro-framework/microservice/endpoint"
	"runtime"
	"strconv"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	file, line := runtime.FuncForPC(pc).FileLine(pc)
	return runtime.FuncForPC(pc).Name() + "::" + file + "::" + strconv.Itoa(line)
}

// KafkaConsumer encapsulates all members used for the kafka server.
type KafkaConsumer struct {
	cfg    *configuration.KafkaConfig
	config *sarama.Config
	master sarama.Consumer
	topics map[string]Topic
	quit   chan bool
}

// Topic encapsulates a kafka channel name and a map of endpoints used on that channel
type Topic struct {
	Name               string
	ServiceEndpointMap ServiceEndpointMap
}

// MyTopic returns a new Topic instance
func MyTopic(name string) (topic Topic) {
	return Topic{
		Name:               name,
		ServiceEndpointMap: make(ServiceEndpointMap),
	}
}

// MyKafkaConsumer takes in configuration object and returns a new KafkaConsumer instance.
func MyKafkaConsumer(cfg *configuration.KafkaConfig) (kc *KafkaConsumer, err error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer(cfg.Brokers, config)
	if err != nil {
		return
	}

	kc = &KafkaConsumer{
		cfg,
		config,
		master,
		make(map[string]Topic),
		make(chan bool),
	}

	return
}

// Start is an implementation of the Server Start() method.
func (kc *KafkaConsumer) Start() (err error) {
	for _, topic := range kc.topics {
		go func(topic Topic) {
			consumer(kc.master, kc.quit, topic, ByteHandler, kc.cfg)
			select {
			case <-kc.quit:
				return
			}
		}(topic)
	}
	return
}

// Stop is an implementation of the Server Stop() method
func (kc *KafkaConsumer) Stop() (err error) {
	close(kc.quit)
	kc.master.Close()
	return
}

// RegisterNamespace is an implementation of the Server RegisterNamespace method.
// In kafka a namespace is the same as the channel name.
func (kc *KafkaConsumer) RegisterNamespace(name string) {
	kc.topics[name] = MyTopic(name)
	return
}

// RegisterService is an implementation of the server RegisterService method.
// In kafka this is used to demux the messages coming in on a single channel/topic.
func (kc *KafkaConsumer) RegisterService(namespace string, service Service, ep endpoint.Endpoint) {
	kc.topics[namespace].ServiceEndpointMap[service] = ep
	return
}

func consumer(master sarama.Consumer, quit chan bool, topic Topic, consume Handler, cfg *configuration.KafkaConfig) {
	log.WithFields(log.Fields{
		"FunctionName": funcName(),
		"Topic":        topic.Name,
	}).Info("Starting Consumer")

	consumer, err := master.ConsumePartition(topic.Name, cfg.Partition, cfg.Offset)
	if err != nil {
		log.WithFields(log.Fields{
			"FunctionName": funcName(),
			"Topic":        topic.Name,
			"Error":        err.Error(),
		}).Error("Failed ConsumerPartition")
	}

	defer func() {
		if err = consumer.Close(); err != nil {
			log.WithFields(log.Fields{
				"FunctionName": funcName(),
				"Error":        err.Error(),
			}).Error("Failed to close Producer")
		}

		if err := master.Close(); err != nil {
			log.WithFields(log.Fields{
				"FunctionName": funcName(),
				"Error":        err.Error(),
			}).Error("Failed to close Producer")
		}
	}()

	for {
		select {
		case <-quit:
			return
		case err := <-consumer.Errors():
			log.WithFields(log.Fields{
				"FunctionName": funcName(),
				"Error":        err.Error(),
			}).Debug("Received error on consumer channel.")
		case msg := <-consumer.Messages():
			log.WithFields(log.Fields{
				"FunctionName": funcName(),
				"Topic":        topic.Name,
			}).Debugf("CONSUMER - Message consumed %#v\n", string(msg.Value))
			_, err := consume(msg.Value, topic.ServiceEndpointMap)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Debug("Received error on consumer channel.")
			}
		}
	}
}

// KafkaProducer encapsulates all the members used by a kafka producer
type KafkaProducer struct {
	cfg     *configuration.KafkaConfig
	config  *sarama.Config
	master  sarama.SyncProducer
	msgChan chan sarama.ProducerMessage
	quit    chan bool
}

// MyKafkaProducer takes in a configuration object and returns a KafkaProducer instance.
func MyKafkaProducer(cfg *configuration.KafkaConfig) (kp *KafkaProducer, err error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	master, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return
	}

	kp = &KafkaProducer{
		cfg,
		config,
		master,
		make(chan sarama.ProducerMessage),
		make(chan bool),
	}

	return
}

// Start is used to start a go routine for the kafka producer.
// All messages produced are read by the msgChan passed to the producer() method.
func (kp *KafkaProducer) Start() (err error) {
	go func() {
		producer(kp.master, kp.msgChan, kp.quit, kp.cfg)
		select {
		case <-kp.quit:
			return
		}
	}()
	return
}

// Stop gracefully stops the kafka producer by writing to the quit channel.
func (kp *KafkaProducer) Stop() (err error) {
	close(kp.quit)
	kp.master.Close()
	return
}

// Produce is used by the client to send a message through kafka producer.
func (kp *KafkaProducer) Produce(topic string, msg interface{}) (err error) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"FunctionName": funcName(),
			"Value":        msg,
			"Error":        err,
		}).Error("Failed Marshal")
		return err
	}

	kp.msgChan <- sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(b),
	}
	return
}

func producer(master sarama.SyncProducer, msgChan chan sarama.ProducerMessage, quit chan bool, cfg *configuration.KafkaConfig) {
	log.WithFields(log.Fields{
		"FunctionName": funcName(),
	}).Info("Starting Producer")

	defer func() {
		if err := master.Close(); err != nil {
			log.WithFields(log.Fields{
				"FunctionName": funcName(),
				"Error":        err,
			}).Error("Failed to close Producer")
		}
	}()

	for {
		select {
		case <-quit:
			return
		case msg := <-msgChan:
			partition, offset, err := master.SendMessage(&msg)
			if err != nil {
				log.WithFields(log.Fields{
					"FunctionName": funcName(),
					"Error":        err,
				}).Error("Failed SendMessage")
				continue
			}

			log.WithFields(log.Fields{
				"FunctionName": funcName(),
				"Topic":        msg.Topic,
				"Partition":    partition,
				"Offset":       offset,
			}).Debug("PRODUCER - Message Sent")
		}
	}
}
