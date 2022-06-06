package cache

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eko/gocache/v2/store"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

type Nats struct {
	Conn   *nats.Conn
	KV     nats.KeyValue
	logger *zerolog.Logger
}

type NatsConfig struct {
	CAPath   string
	CertPath string
	KeyPath  string
	Address  string
	KVConfig nats.KeyValueConfig
}

func NewNatsJSCache(logger *zerolog.Logger, cfg NatsConfig) (*Nats, error) {
	var opts []nats.Option
	if cfg.KVConfig.Bucket == "" {
		return nil, fmt.Errorf("NatsConfig:%v - Bucket is required", cfg.KVConfig)
	}
	if cfg.CAPath != "" {
		opts = append(opts, nats.RootCAs(cfg.CAPath))
	}

	if cfg.CertPath != "" {
		opts = append(opts, nats.ClientCert(cfg.CertPath, cfg.KeyPath))
	}

	nc, err := nats.Connect(cfg.Address, opts...)
	if err != nil {
		logger.Err(err).Msg("error connecting NATS client")
		return nil, err
	}
	jscontext, err := nc.JetStream()
	if err != nil {
		logger.Err(err).Msg("error initiating JetStream context")
		return nil, err
	}
	kv, err := jscontext.CreateKeyValue(&cfg.KVConfig)
	if err != nil {
		logger.Err(err).Msg("error creating nats key-value store")
		return nil, err
	}
	logger.Debug().Msg("NATS client connected")
	return &Nats{
		Conn:   nc,
		logger: logger,
		KV:     kv,
	}, nil
}

func (n *Nats) Get(ctx context.Context, key interface{}) (interface{}, error) {
	entry, err := n.KV.Get(fmt.Sprintf("%v", key))
	if err != nil {
		n.logger.Error().Msgf("NATS GET: %v - error %v", key, err)
		return nil, err
	}
	return entry.Value(), nil
}

func (n *Nats) Set(ctx context.Context, key, object interface{}, options *store.Options) error {
	switch object.(type) {
	case string:
		{
			revision, err := n.KV.PutString(fmt.Sprintf("%v", key), fmt.Sprintf("%s", object))
			n.logger.Trace().Msgf("NATS PUT String: %v : %v - revision %v", key, object, revision)
			return err
		}
	case []byte:
		{
			revision, err := n.KV.Put(fmt.Sprintf("%v", key), []byte(fmt.Sprintf("%s", object)))
			n.logger.Trace().Msgf("NATS PUT: %v : %v - revision %v", key, object, revision)
			return err
		}
	default:
		{
			value, err := json.Marshal(object)
			if err != nil {
				n.logger.Error().Msgf("Failed to marshal object: %v", object)
			}
			revision, err := n.KV.Put(fmt.Sprintf("%v", key), value)
			n.logger.Trace().Msgf("NATS PUT: %v : %v - revision %v", key, object, revision)
			return err
		}
	}
}
func (n *Nats) Delete(ctx context.Context, key interface{}) error {
	return n.KV.Delete(fmt.Sprintf("%v", key))
}
func (n *Nats) Invalidate(ctx context.Context, options store.InvalidateOptions) error {
	keys, err := n.KV.Keys()
	if err != nil {
		return err
	}
	for _, key := range keys {
		delerr := n.KV.Delete(key)
		if delerr != nil {
			n.logger.Error().Msgf("NATS Delete: %v - error %v", key, delerr)
		}
	}
	return nil
}
func (n *Nats) Clear(ctx context.Context) error {
	keys, err := n.KV.Keys()
	if err != nil {
		return err
	}
	for _, key := range keys {
		err = n.KV.Purge(key)
		if err != nil {
			n.logger.Error().Msgf("NATS Purge: %v - error %v", key, err)
		}
	}
	return n.KV.PurgeDeletes()
}
func (n *Nats) GetType() string {
	return "nats-jetstream"
}
