package kafka_toolkit

import (
	"errors"

	"github.com/Shopify/sarama"
)

const (
	//RoundRobin estrategia round robin
	RoundRobin = "roundrobin"
	//Range estrategia range
	Range = "range"
)

//BalanceStrategyResolver resuelve estrategia de balance
type BalanceStrategyResolver interface {
	Resolve(balanceStrategy string) (sarama.BalanceStrategy, error)
}

type balanceStrategyResolver struct{}

//NewBalanceStrategyResolver constructo de BalanceStrategyResolver
func NewBalanceStrategyResolver() BalanceStrategyResolver {
	return &balanceStrategyResolver{}
}

//Resolve resuelve a estrategia de sarama, si estrategia es invalida retorna error
func (r *balanceStrategyResolver) Resolve(balanceStrategy string) (sarama.BalanceStrategy, error) {
	switch balanceStrategy {
	case RoundRobin:
		return sarama.BalanceStrategyRoundRobin, nil
	case Range:
		return sarama.BalanceStrategyRange, nil
	default:
		return nil, errors.New(InvalidBalanceStrategyKind)
	}
}
