package kafka_toolkit

import (
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var FlushFrequencyByDefault = time.Millisecond * time.Duration(200)

// BaseProducerConfigInput entrada de configuracion de producer
type BaseProducerConfigInput struct {
	Brokers          string
	Topic            string
	Ack              int16
	Retries          int
	Security         bool
	Username         string
	Password         string
	Mechanism        string
	CaFile           string
	Version          string
	ClientID         string
	FlushFrequencyMs int64
	TimeoutMs        int64
}

// BaseProducerConfig configuracion base de producer
type BaseProducerConfig struct {
	Brokers      []string
	SaramaConfig *sarama.Config
}

// BaseProducerConfigurer Convierte input en configuracion de Producer
type BaseProducerConfigurer interface {
	GenerateConfig(confInput BaseProducerConfigInput) (*BaseProducerConfig, error)
}

type baseProducerConfigurer struct {
}

// NewBaseProducerConfigurer crea un nuevo producer configurer
func NewBaseProducerConfigurer() BaseProducerConfigurer {
	return &baseProducerConfigurer{}
}

func (b *baseProducerConfigurer) GenerateConfig(confInput BaseProducerConfigInput) (*BaseProducerConfig, error) {
	brokers := strings.Split(confInput.Brokers, ",")

	config := sarama.NewConfig()

	if err := configVersion(confInput, config); err != nil {
		return nil, err
	}

	if confInput.ClientID != "" {
		config.ClientID = confInput.ClientID
	}
	config.Producer.RequiredAcks = sarama.RequiredAcks(confInput.Ack)
	config.Producer.Retry.Max = confInput.Retries
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	config.Producer.Retry.Backoff = 250 * time.Millisecond

	if confInput.FlushFrequencyMs > 0 {
		config.Producer.Flush.Frequency = time.Millisecond * time.Duration(confInput.FlushFrequencyMs)
	} else {
		config.Producer.Flush.Frequency = FlushFrequencyByDefault
	}

	if confInput.Security {
		saslConfig, err := SaramaSASLConfig(KafkaSASLSecurity{
			Username:      confInput.Username,
			Password:      confInput.Password,
			CaFile:        confInput.CaFile,
			MechanismSASL: confInput.Mechanism,
		})
		if err != nil {
			return nil, err
		}

		config.Net = saslConfig.Net
	}

	return &BaseProducerConfig{brokers, config}, nil
}

func configVersion(confInput BaseProducerConfigInput, config *sarama.Config) error {
	if len(confInput.Version) > 0 {
		version, err := sarama.ParseKafkaVersion(confInput.Version)

		if err != nil {
			return err
		}

		config.Version = version
	}

	return nil
}
