package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/encoding/protodelim"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"pow/internal/dao"
	"pow/internal/proof_of_work"
	"pow/internal/server_config"
	"pow/protocol"
	"time"
)

type Server struct {
	ctx              context.Context
	log              *logrus.Entry
	solveTimeout     time.Duration
	challengeConfig  server_config.Challenge
	bind             string
	proofOfWork      *proof_of_work.ProofOfWork
	wordsOfWisdomDao *dao.WordsOfWisdom
}

func NewServer(
	ctx context.Context,
	log *logrus.Entry,
	bind string,
	solveTimeout time.Duration,
	challengeConfig server_config.Challenge,
	proofOfWork *proof_of_work.ProofOfWork,
	wordsOfWisdomDao *dao.WordsOfWisdom,
) *Server {
	return &Server{
		ctx:              ctx,
		log:              log,
		bind:             bind,
		solveTimeout:     solveTimeout,
		challengeConfig:  challengeConfig,
		proofOfWork:      proofOfWork,
		wordsOfWisdomDao: wordsOfWisdomDao,
	}
}

func (s *Server) Demo() error {

	return nil
}

func (s *Server) Run() error {
	s.log.Infof("starting server at %s", s.bind)

	listen, err := net.Listen("tcp", s.bind)

	if err != nil {
		return err
	}

	s.log.Infof("server started at %s", s.bind)

	defer func() {
		if err := listen.Close(); err != nil {
			s.log.WithError(err).Errorf("failed to close listener")
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			s.log.Warnf("context is done, shutting down")
			return nil
		default:
			conn, err := listen.Accept()

			if err != nil {
				return err
			}

			s.log.Infof("accepted new connection from %s", conn.RemoteAddr())

			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	log := s.log.WithField("service", "connection").WithField("remote_addr", conn.RemoteAddr())

	done := make(chan struct{})

	defer func() {
		log.Infof("closing connection")

		if err := conn.Close(); err != nil {
			log.WithError(err).Errorf("failed to close connection")
		}
	}()

	go func() {
		defer close(done)

		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)

		challenge, err := s.proofOfWork.GenerateChallenge(s.challengeConfig.Length)

		if err != nil {
			log.WithError(err).Errorf("failed to generate challenge")
			return
		}

		if err := s.sendChallenge(writer, challenge, s.challengeConfig.Difficulty); err != nil {
			log.WithError(err).Errorf("failed to send challenge")
			return
		}

		for {
			var message = &protocol.ClientMessage{}

			if err := protodelim.UnmarshalFrom(reader, message); err != nil {
				if errors.Is(err, io.EOF) {
					log.Debugf("EOF")
					return
				}

				log.WithError(err).Errorf("failed to read from connection")
				return
			}

			switch data := message.Data.(type) {
			case *protocol.ClientMessage_ChallengeSolved_:
				if !s.proofOfWork.Verify(challenge, data.ChallengeSolved.Nonce, s.challengeConfig.Difficulty) {
					log.WithError(err).Errorf("failed to verify challenge")

					if err := s.sendChallengeFailed(writer); err != nil {
						log.WithError(err).Errorf("failed to send challenge failed")
						return
					}
				} else {
					log.Debug("challenge solved")

					quote, err := s.wordsOfWisdomDao.GetRandomQuote()

					if err != nil {
						log.WithError(err).Errorf("failed to get quote")
						return
					}

					if err := s.sendWordOfWisdom(writer, quote); err != nil {
						log.WithError(err).Errorf("failed to send word of Wisdom")
						return
					}
				}
			default:
				log.Errorf("unexpected message from client: %s", message)
				return
			}
		}
	}()

	timeout := time.After(s.solveTimeout)

	select {
	case <-s.ctx.Done():
	case <-done:
	case <-timeout:
		log.Warnf("connection timed out")
	}
}

func (s *Server) sendChallenge(writer *bufio.Writer, challenge []byte, difficulty int) error {
	message := &protocol.ServerMessage{
		Timestamp: time.Now().UnixNano(),
		Data: &protocol.ServerMessage_Challenge_{
			Challenge: &protocol.ServerMessage_Challenge{
				Difficulty: int32(difficulty),
				Challenge:  challenge,
			},
		},
	}

	return s.sendMessage(writer, message)
}

func (s *Server) sendChallengeFailed(writer *bufio.Writer) error {
	message := &protocol.ServerMessage{
		Timestamp: time.Now().UnixNano(),
		Data: &protocol.ServerMessage_ChallengeFailed_{
			ChallengeFailed: &protocol.ServerMessage_ChallengeFailed{},
		},
	}

	return s.sendMessage(writer, message)
}

func (s *Server) sendWordOfWisdom(writer *bufio.Writer, quote string) error {
	message := &protocol.ServerMessage{
		Timestamp: time.Now().UnixNano(),
		Data: &protocol.ServerMessage_WordOfWisdom_{
			WordOfWisdom: &protocol.ServerMessage_WordOfWisdom{
				Quote: quote,
			},
		},
	}

	return s.sendMessage(writer, message)
}

func (s *Server) sendMessage(writer *bufio.Writer, message proto.Message) error {
	if _, err := protodelim.MarshalTo(writer, message); err != nil {
		return fmt.Errorf("error encoding message: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer: %w", err)
	}

	return nil
}
