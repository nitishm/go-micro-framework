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

func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	file, line := runtime.FuncForPC(pc).FileLine(pc)
	return runtime.FuncForPC(pc).Name() + "::" + file + "::" + strconv.Itoa(line)
}

type KafkaConsumer struct {
	cfg    *configuration.KafkaConfig
	config *sarama.Config
	master sarama.Consumer
	topics map[string]Topic
	quit   chan bool
}
type Topic struct {
	Name               string
	ServiceEndpointMap ServiceEndpointMap
}

func MyTopic(name string) (topic Topic) {
	return Topic{
		Name:               name,
		ServiceEndpointMap: make(ServiceEndpointMap),
	}
}

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

func (kc *KafkaConsumer) Stop() (err error) {
	close(kc.quit)
	kc.master.Close()
	return
}

func (kc *KafkaConsumer) RegisterNamespace(name string) {
	kc.topics[name] = MyTopic(name)
	return
}

func (kc *KafkaConsumer) RegisterService(namespace string, service Service, ep endpoint.Endpoint) {
	kc.topics[namespace].ServiceEndpointMap[service] = ep
	return
}

func consumer(master sarama.Consumer, quit chan bool, topic Topic, consume Handler, cfg *configuration.KafkaConfig) {
	log.WithFields(log.Fields{
		"FunctionName": FuncName(),
		"Topic":        topic.Name,
	}).Info("Starting Consumer")

	consumer, err := master.ConsumePartition(topic.Name, cfg.Partition, cfg.Offset)
	if err != nil {
		log.WithFields(log.Fields{
			"FunctionName": FuncName(),
			"Topic":        topic.Name,
			"Error":        err.Error(),
		}).Error("Failed ConsumerPartition")
	}

	defer func() {
		if err = consumer.Close(); err != nil {
			log.WithFields(log.Fields{
				"FunctionName": FuncName(),
				"Error":        err.Error(),
			}).Error("Failed to close Producer")
		}

		if err := master.Close(); err != nil {
			log.WithFields(log.Fields{
				"FunctionName": FuncName(),
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
				"FunctionName": FuncName(),
				"Error":        err.Error(),
			}).Debug("Received error on consumer channel.")
		case msg := <-consumer.Messages():
			log.WithFields(log.Fields{
				"FunctionName": FuncName(),
				"Topic":        topic.Name,
			}).Debugf("CONSUMER - Message consumed %#v\n", string(msg.Value))
			err := consume(msg.Value, topic.ServiceEndpointMap)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Debug("Received error on consumer channel.")
			}
		}
	}
}

type KafkaProducer struct {
	cfg     *configuration.KafkaConfig
	config  *sarama.Config
	master  sarama.SyncProducer
	msgChan chan sarama.ProducerMessage
	quit    chan bool
}

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

func (kp *KafkaProducer) Stop() (err error) {
	close(kp.quit)
	kp.master.Close()
	return
}

func (kp *KafkaProducer) Produce(topic string, msg interface{}) (err error) {
	b, err := json.Marshal(msg)
	if err != nil {
		log.WithFields(log.Fields{
			"FunctionName": FuncName(),
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
		"FunctionName": FuncName(),
	}).Info("Starting Producer")

	defer func() {
		if err := master.Close(); err != nil {
			log.WithFields(log.Fields{
				"FunctionName": FuncName(),
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
					"FunctionName": FuncName(),
					"Error":        err,
				}).Error("Failed SendMessage")
				continue
			}

			log.WithFields(log.Fields{
				"FunctionName": FuncName(),
				"Topic":        msg.Topic,
				"Partition":    partition,
				"Offset":       offset,
			}).Debug("PRODUCER - Message Sent")
		}
	}
}
