package log

import (
	"log/slog"
	"net"
	"time"
)

type LogstashConnection struct {
	net.Conn
	url string
}

func (l *LogstashConnection) Write(data []byte) (int, error) {

	if l.Conn == nil {

		conn, err := newTcpConnection(l.url, 0)
		if err != nil {
			slog.Error("[Logstash] unable to connect to Logstash. Data will be dropped", slog.Any("error", err), slog.String("data", string(data)))
			return 0, err
		}

		l.Conn = conn
	}

	size, err := l.Conn.Write(data)
	if err != nil {
		slog.Error("[Logstash] write log data failed. Data will be dropped", slog.Any("error", err), slog.String("data", string(data)))
		l.Conn.Close()
		l.Conn = nil
		return size, err
	}

	slog.Info("[Logstash] write log data success")
	return size, nil
}

func (l *LogstashConnection) Close() error {

	if l.Conn == nil {
		return nil
	}

	return l.Conn.Close()
}

const (
	connectRetryCount    = 2
	retryCooldown        = 2 * time.Second
	newConnectionTimeout = 5 * time.Second
)

func NewLogstashConnection(url string) (*LogstashConnection, error) {

	newLogstashConn := &LogstashConnection{
		url: url,
	}

	conn, err := newTcpConnection(url, connectRetryCount)
	if err != nil {
		return newLogstashConn, err
	}

	newLogstashConn.Conn = conn

	return newLogstashConn, nil
}

func newTcpConnection(url string, retryCount int) (net.Conn, error) {

	var conn net.Conn
	var err error
	for i := 0; i <= retryCount; i++ {
		conn, err = net.DialTimeout("tcp", url, newConnectionTimeout)

		if err == nil {
			slog.Info("[Logstash] Connected to feed log", slog.String("url", url))
			return conn, nil
		}

		slog.Warn("[Logstash] Unable to connect to logstash",
			slog.String("url", url),
			slog.Int("attempt", i+1),
			slog.Any("error", err),
		)

		if i != retryCount {
			time.Sleep(retryCooldown)
		}
	}

	return nil, err
}
