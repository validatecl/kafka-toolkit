package kafka_toolkit

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"github.com/Shopify/sarama"
)

var (
	InvalidUsernamePassword = "Usuario o Contraseña viene sin información"
)

type KafkaSASLSecurity struct {
	Username      string
	Password      string
	CaFile        string
	MechanismSASL string
}

// SaramaTLSConfig retorna TLS config
func SaramaTLSConfig(caFile string) (*tls.Config, error) {
	ca, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(ca)
	tlsConfig := new(tls.Config)
	tlsConfig.RootCAs = certPool
	tlsConfig.InsecureSkipVerify = true
	return tlsConfig, nil
}

func SaramaSASLConfig(input KafkaSASLSecurity) (*sarama.Config, error) {
	var err error
	saramaConf := sarama.NewConfig()

	if input.Username == "" || input.Password == "" {
		return nil, errors.New(InvalidUsernamePassword)
	}

	if input.MechanismSASL == "" {
		input.MechanismSASL = sarama.SASLTypePlaintext
	}

	saramaConf.Net.SASL.Enable = true
	saramaConf.Net.SASL.Mechanism = sarama.SASLMechanism(input.MechanismSASL)

	switch input.MechanismSASL {
	case sarama.SASLTypeSCRAMSHA512:
		saramaConf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA512} }
	case sarama.SASLTypeSCRAMSHA256:
		saramaConf.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &XDGSCRAMClient{HashGeneratorFcn: SHA256} }
	default:
	}

	saramaConf.Net.SASL.User = input.Username
	saramaConf.Net.SASL.Password = input.Password
	saramaConf.Net.TLS.Enable = true
	saramaConf.Net.TLS.Config, err = SaramaTLSConfig(input.CaFile)
	if err != nil {
		return nil, err
	}
	return saramaConf, nil
}
