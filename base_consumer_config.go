package kafka_toolkit

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

//ConsumerGroupInput Represents a consumer group config input
type ConsumerGroupInput struct {
	Brokers         string
	Topic           string
	Group           string
	ClientID        string
	BalanceStrategy string
	// Deprecado (A ser removido en la proxima minor version)
	Earliest               bool
	Latest                 bool
	Version                string
	Security               bool
	Username               string
	Password               string
	Mechanism              string
	CaFile                 string
	SessionDurationSeconds int64
}

//ConsumerGroupConfig represents a consumer group config
type ConsumerGroupConfig struct {
	Topic        string
	Brokers      []string
	Group        string
	SaramaConfig *sarama.Config
}

//SaramaConsumerConfigurer generates Sarama Consumer config
type SaramaConsumerConfigurer interface {
	GenerateConfig(input ConsumerGroupInput) (*ConsumerGroupConfig, error)
}

type saramaConsumerConfigurer struct {
	balanceStrategyResolver BalanceStrategyResolver
}

//NewSaramaConsumerConfigurer constructor
func NewSaramaConsumerConfigurer(b BalanceStrategyResolver) SaramaConsumerConfigurer {
	return &saramaConsumerConfigurer{balanceStrategyResolver: b}
}

func (s *saramaConsumerConfigurer) GenerateConfig(input ConsumerGroupInput) (*ConsumerGroupConfig, error) {
	consumerConfig := new(ConsumerGroupConfig)

	consumerConfig.Topic = input.Topic
	consumerConfig.Brokers = strings.Split(input.Brokers, ",")
	consumerConfig.Group = input.Group

	saramaConf, err := s.parseSaramaConsumerConfig(input)

	if err != nil {
		return nil, err
	}

	consumerConfig.SaramaConfig = saramaConf

	return consumerConfig, nil
}

func (s *saramaConsumerConfigurer) parseSaramaConsumerConfig(input ConsumerGroupInput) (*sarama.Config, error) {
	saramaConf := sarama.NewConfig()

	if input.ClientID != "" {
		clientIDWithPID := fmt.Sprintf("%s_%d", input.ClientID, os.Getpid())

		saramaConf.ClientID = clientIDWithPID
	}

	version, err := sarama.ParseKafkaVersion(input.Version)

	if err != nil {
		return nil, err
	}

	saramaConf.Version = version

	strategy, err := s.balanceStrategyResolver.Resolve(input.BalanceStrategy)

	if err != nil {
		return nil, err
	}

	saramaConf.Consumer.Group.Rebalance.Strategy = strategy

	// TODO: Remover earliest en proxima version minor
	if input.Latest || input.Earliest {
		saramaConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		saramaConf.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	saramaConf.Consumer.Offsets.AutoCommit.Interval = 250 * time.Millisecond

	if !input.Security {
		return saramaConf, nil
	}

	saslConfig, err := SaramaSASLConfig(KafkaSASLSecurity{
		Username:      input.Username,
		Password:      input.Password,
		CaFile:        input.CaFile,
		MechanismSASL: input.Mechanism,
	})

	if input.SessionDurationSeconds > 0 {
		saramaConf.Consumer.Group.Session.Timeout = time.Duration(input.SessionDurationSeconds) * time.Second
	}

	if err != nil {
		return nil, err
	}

	saramaConf.Net = saslConfig.Net

	return saramaConf, nil

}
