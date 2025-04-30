package ethnode

import (
	"blockchain-gateway/config"
	"context"
	"fmt"
	"github.com/knstch/subtrack-kafka/producer"
)

type MempoolScanner struct {
	addressesToCheck map[string]interface{}
	producer         *producer.Producer
	node             *Node
}

func MakeBscMempoolScanner(ctx context.Context, cfg config.Config) (*MempoolScanner, error) {
	node, err := NewNode(ctx, cfg.Blockchains.BscWsAddr)
	if err != nil {
		return nil, fmt.Errorf("ethnode.NewNode: %w", err)
	}

	kafkaProducer := producer.NewProducer(cfg.KafkaAddr)

	addressesToCheck := make(map[string]interface{}, len(cfg.Blockchains.BscWsAddr))
	for _, addr := range cfg.Blockchains.BscDexAddresses {
		addressesToCheck[addr] = nil
	}

	return &MempoolScanner{
		addressesToCheck: addressesToCheck,
		producer:         kafkaProducer,
		node:             node,
	}, nil
}
