package harego

import (
	"github.com/streadway/amqp"
)

// rabbitWrapper is defined to make it easy for passing a mocked connection.
type rabbitWrapper struct {
	*amqp.Connection
}

// Channel returns the underlying channel.
func (r *rabbitWrapper) Channel() (Channel, error) {
	return r.Connection.Channel()
}

// ConfigFunc is a function for setting up the Client. A config function returns
// an error if the client is already started.
type ConfigFunc func(*Client) error

// Connection sets the RabbitMQ connection.
func Connection(r RabbitMQ) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.conn = r
		return nil
	}
}

// AMQP uses the connection for the broker.
func AMQP(r *amqp.Connection) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.conn = &rabbitWrapper{r}
		return nil
	}
}

// QueueName sets the queue name.
func QueueName(name string) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.queueName = name
		return nil
	}
}

// RoutingKey sets the routing key of the queue.
func RoutingKey(key string) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.routingKey = key
		return nil
	}
}

// Workers sets the worker count for cusuming messages.
func Workers(n int) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.workers = n
		return nil
	}
}

// WithDeliveryMode sets the default delivery mode of messages.
func WithDeliveryMode(mode DeliveryMode) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.deliveryMode = mode
		return nil
	}
}

// PrefetchCount sets how many items should be prefetched for consumption. With
// a prefetch count greater than zero, the server will deliver that many
// messages to consumers before acknowledgments are received. The server ignores
// this option when consumers are started with noAck because no acknowledgments
// are expected or sent.
func PrefetchCount(i int) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.prefetchCount = i
		return nil
	}
}

// PrefetchSize sets the prefetch size of the Qos. If it is greater than zero,
// the server will try to keep at least that many bytes of deliveries flushed to
// the network before receiving acknowledgments from the consumers.
func PrefetchSize(i int) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.prefetchSize = i
		return nil
	}
}

// WithExchangeType sets the exchange type. The default is ExchangeTypeTopic.
func WithExchangeType(t ExchangeType) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.exchType = t
		return nil
	}
}

// ExchangeName sets the exchange name. For each worker, and additional string
// will be appended for the worker number.
func ExchangeName(name string) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.exchName = name
		return nil
	}
}

// ConsumerName sets the consumer name of the consuming queue.
func ConsumerName(name string) ConfigFunc {
	return func(c *Client) error {
		if c.started {
			return ErrAlreadyStarted
		}
		c.consumerName = name
		return nil
	}
}

// NotDurable marks the exchange and the queue not to be durable. Default is
// durable.
func NotDurable(c *Client) error {
	if c.started {
		return ErrAlreadyStarted
	}
	c.durable = false
	return nil
}

// AutoDelete marks the exchange and queues with autoDelete property which
// causes the messages to be automatically removed from the queue when
// consumed.
func AutoDelete(c *Client) error {
	if c.started {
		return ErrAlreadyStarted
	}
	c.autoDelete = true
	return nil
}

// Internal sets the exchange to be internal.
func Internal(c *Client) error {
	if c.started {
		return ErrAlreadyStarted
	}
	c.internal = true
	return nil
}

// NoWait marks the exchange as noWait. When noWait is true, declare
// without waiting for a confirmation from the server. The channel may be closed
// as a result of an error.
func NoWait(c *Client) error {
	if c.started {
		return ErrAlreadyStarted
	}
	c.noWait = true
	return nil
}
