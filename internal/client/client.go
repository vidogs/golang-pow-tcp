package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protodelim"
	"google.golang.org/protobuf/proto"
	"net"
	"pow/internal/proof_of_work"
	"pow/protocol"
	"time"
)

type Client struct {
	ctx         context.Context
	log         *logrus.Entry
	address     string
	proofOfWork *proof_of_work.ProofOfWork
}

func NewClient(ctx context.Context, log *logrus.Entry, address string, proofOfWork *proof_of_work.ProofOfWork) *Client {
	return &Client{
		ctx:         ctx,
		log:         log,
		address:     address,
		proofOfWork: proofOfWork,
	}
}

func (c *Client) Run() error {
	conn, err := net.Dial("tcp", c.address)

	if err != nil {
		return fmt.Errorf("error connecting to server at %s: %w", c.address, err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			c.log.WithError(err).Errorf("error closing connection")
		}
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		var message = &protocol.ServerMessage{}

		if err := protodelim.UnmarshalFrom(reader, message); err != nil {
			return fmt.Errorf("error decoding message: %w", err)
		}

		switch data := message.Data.(type) {
		case *protocol.ServerMessage_Challenge_:
			nonce := c.proofOfWork.Solve(data.Challenge.Challenge, int(data.Challenge.Difficulty))

			if err := c.sendChallengeSolved(writer, nonce); err != nil {
				c.log.WithError(err).Errorf("error sending solved challenge")

				return err
			}
		case *protocol.ServerMessage_ChallengeFailed_:
			c.log.Errorf("challenge failed")

			return fmt.Errorf("challenge failed")
		case *protocol.ServerMessage_WordOfWisdom_:
			c.log.Infof("word of wisdom: %s\n", data.WordOfWisdom.Quote)

			return nil
		default:
			return fmt.Errorf("unexpected message from server: %s", message)
		}
	}
}

func (c *Client) sendChallengeSolved(writer *bufio.Writer, nonce []byte) error {
	message := &protocol.ClientMessage{
		Timestamp: time.Now().UnixNano(),
		Data: &protocol.ClientMessage_ChallengeSolved_{
			ChallengeSolved: &protocol.ClientMessage_ChallengeSolved{
				Nonce: nonce,
			},
		},
	}

	return c.sendMessage(writer, message)
}

func (c *Client) sendMessage(writer *bufio.Writer, message proto.Message) error {
	if _, err := protodelim.MarshalTo(writer, message); err != nil {
		return fmt.Errorf("error encoding message: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}
