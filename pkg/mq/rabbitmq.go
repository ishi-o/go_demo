package mq

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ishi-o/go_demo/pkg/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var ErrClientClosed = errors.New("client is closed")

type RabbitMQClient struct {
	conn      *amqp.Connection
	pool      *ChannelPool
	url       string
	mu        sync.RWMutex
	reconnect chan struct{}
	closed    bool
}

type ChannelPool struct {
	conn     *amqp.Connection
	mu       sync.Mutex
	channels chan *amqp.Channel
	maxSize  int
}

func New(cfg *config.RabbitMQConfig) (*RabbitMQClient, error) {
	if cfg.PoolSize <= 0 {
		cfg.PoolSize = 10
	}

	client := &RabbitMQClient{
		url:       cfg.URL,
		reconnect: make(chan struct{}, 1),
	}

	if err := client.connect(cfg.PoolSize); err != nil {
		return nil, err
	}

	go client.handleReconnect(cfg.PoolSize)
	return client, nil
}

func (c *RabbitMQClient) connect(poolsize int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return ErrClientClosed
	}

	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}

	pool, err := NewChannelPool(conn, poolsize)
	if err != nil {
		conn.Close()
		return err
	}

	c.conn = conn
	c.pool = pool

	return nil
}

func (c *RabbitMQClient) handleReconnect(poolsize int) {
	for {
		<-c.reconnect

		c.mu.RLock()
		if c.closed {
			c.mu.RUnlock()
			return
		}
		c.mu.RUnlock()

		for i := 0; i < 5; i++ {
			if err := c.connect(poolsize); err == nil {
				break
			}
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}
}

func (c *RabbitMQClient) GetChannel() (*amqp.Channel, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, ErrClientClosed
	}

	return c.pool.Get()
}

func (c *RabbitMQClient) PutChannel(ch *amqp.Channel) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		ch.Close()
		return
	}

	c.pool.Put(ch)
}

// Close closes the client and releases all resources
func (c *RabbitMQClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return nil
	}

	c.closed = true
	close(c.reconnect)

	var errs []error

	if c.pool != nil {
		c.pool.Close() // pool.Close() doesn't return error
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close connection: %w", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("rabbitmq client close errors: %v", errs)
	}
	return nil
}

func NewChannelPool(conn *amqp.Connection, maxSize int) (*ChannelPool, error) {
	if conn == nil {
		return nil, ErrClientClosed
	}

	pool := &ChannelPool{
		conn:     conn,
		channels: make(chan *amqp.Channel, maxSize),
		maxSize:  maxSize,
	}

	for i := 0; i < maxSize/2; i++ {
		ch, err := conn.Channel()
		if err != nil {
			pool.Close()
			return nil, err
		}
		pool.channels <- ch
	}

	return pool, nil
}

func (p *ChannelPool) Get() (*amqp.Channel, error) {
	select {
	case ch := <-p.channels:
		if ch.IsClosed() {
			return p.createChannel()
		}
		return ch, nil
	default:
		return p.createChannel()
	}
}

func (p *ChannelPool) Put(ch *amqp.Channel) {
	if ch == nil {
		return
	}

	if ch.IsClosed() {
		return
	}

	select {
	case p.channels <- ch:
	default:
		ch.Close()
	}
}

func (p *ChannelPool) createChannel() (*amqp.Channel, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.conn.Channel()
}

func (p *ChannelPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.channels)
	for ch := range p.channels {
		ch.Close()
	}
}

func (c *RabbitMQClient) Publish(exchange, routingKey string, body []byte) error {
	ch, err := c.GetChannel()
	if err != nil {
		return err
	}
	defer c.PutChannel(ch)

	return ch.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		},
	)
}

func (c *RabbitMQClient) Consume(queue, consumer string, autoAck, exclusive bool) (<-chan amqp.Delivery, error) {
	ch, err := c.GetChannel()
	if err != nil {
		return nil, err
	}

	return ch.Consume(
		queue,
		consumer,
		autoAck,
		exclusive,
		false, // no-local
		false, // no-wait
		nil,   // args
	)
}

func (c *RabbitMQClient) DeclareQueue(name string, durable, autoDelete, exclusive bool) (*amqp.Queue, error) {
	ch, err := c.GetChannel()
	if err != nil {
		return nil, err
	}
	defer c.PutChannel(ch)

	q, err := ch.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		false, // no-wait
		nil,   // args
	)

	return &q, err
}

func (c *RabbitMQClient) DeclareExchange(name, kind string, durable, autoDelete bool) error {
	ch, err := c.GetChannel()
	if err != nil {
		return err
	}
	defer c.PutChannel(ch)

	return ch.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		false, // internal
		false, // no-wait
		nil,   // args
	)
}
