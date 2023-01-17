package network

import (
	"bytes"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/rachit77/Eigen-Chain/core"
	"github.com/rachit77/Eigen-Chain/crypto"
	"github.com/rachit77/Eigen-Chain/types"
)

var defaultBlockTime = 5 * time.Second

//server options
type ServerOpts struct {
	ID            string     //to uniquely  identify the server
	Logger        log.Logger //safe for concurrent use
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration //for validator to know when to consume the mempool and produce a block
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	memPool     *TxPool
	chain       *core.Blockchain
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}

	s := &Server{
		ServerOpts:  opts,
		chain:       chain,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil, //potential validator will need a private key
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	//if RPC processor is absent from server options than server will be used as default processor
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

func (s *Server) Start() {
	s.initTransports()

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
			}

			//call processor to process the message
			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				s.Logger.Log("error", err)
			}

		case <-s.quitCh:
			break free

		}
	}

	s.Logger.Log("msg", "Server is shutting down")
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)
	s.Logger.Log("msg", "Starting validator loop", "blockTime", s.BlockTime)

	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	}
	return nil
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {

		}
	}
	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {

	hash := tx.Hash(core.TxHasher{})

	//check if transaction is already in mempool
	if s.memPool.Has(hash) {
		return nil
	}

	//verify the transaction
	if err := tx.Verify(); err != nil {
		return err
	}

	tx.SetFirstSeen(time.Now().UnixNano())

	//add transaction to the mempool
	s.Logger.Log(
		"msg", "adding new tx to mempool",
		"hash", hash,
		"mempoolLength", s.memPool.Len(),
	)

	//broadcast transaction to the peers
	// if err := s.broadcastTx(tx); err != nil {
	// 	logrus.Error(err)
	// }
	go s.broadcastTx(tx)

	return s.memPool.Add(tx)

}

func (s *Server) broadcastBlock(b *core.Block) error {
	return nil
}

//TODO: if not using return value of this function than log the error in this function instead of returning a value
func (s *Server) broadcastTx(tx *core.Transaction) error {
	//encode this transaction
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	//all communication between are of message type
	msg := NewMessage(MessageTypeTx, buf.Bytes())

	//broadcast the message to peers
	return s.broadcast(msg.Bytes())
}

func (s *Server) initTransports() {
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	//TODO: decide parameter and complexity for size of block
	txx := s.memPool.Transactions()
	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	s.memPool.Flush()

	return nil
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: uint64(000000),
	}

	b, _ := core.NewBlock(header, nil)
	return b
}
