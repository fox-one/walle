package broker

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/fox-one/pkg/store/db"
	"github.com/fox-one/walle/core"
	"github.com/fox-one/walle/store/encrypt"
)

func init() {
	db.RegisterMigrate(func(db *db.DB) error {
		tx := db.Update().Model(core.Broker{})

		if err := tx.AutoMigrate(core.Broker{}).Error; err != nil {
			return err
		}

		if err := tx.AddUniqueIndex("idx_brokers_wallet", "wallet_id").Error; err != nil {
			return err
		}

		return nil
	})
}

func New(db *db.DB, enc encrypt.Encrypter) core.BrokerStore {
	return &brokerStore{
		db:  db,
		enc: enc,
	}
}

type brokerStore struct {
	db  *db.DB
	enc encrypt.Encrypter
}

func (s *brokerStore) Create(ctx context.Context, broker *core.Broker) error {
	pin := broker.Pin
	ciphertext, err := s.enc.Encrypt(broker.WalletID + pin)
	if err != nil {
		return err
	}
	broker.Pin = base64.StdEncoding.EncodeToString(ciphertext)

	if err := s.db.Update().Where("wallet_id = ?", broker.WalletID).FirstOrCreate(broker).Error; err != nil {
		return err
	}

	broker.Pin = pin
	return nil
}

func (s *brokerStore) Find(ctx context.Context, walletID string) (*core.Broker, error) {
	var broker core.Broker
	if err := s.db.View().Where("wallet_id = ?", walletID).Take(&broker).Error; err != nil {
		return nil, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(broker.Pin)
	if err != nil {
		return nil, err
	}

	plaintext, err := s.enc.Decrypt(ciphertext)
	if err != nil {
		return nil, err
	}

	broker.Pin = strings.TrimPrefix(plaintext, broker.WalletID)
	return &broker, nil
}
