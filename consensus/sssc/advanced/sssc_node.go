package sssc

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/KOSASIH/pi-nexus-autonomous-banking-network/consensus/sssc/pb"
	"github.com/dgraph-io/badger"
	"google.golang.org/grpc"
)

type SSSCNode struct {
	pb.UnimplementedSSSCNodeServer
	mu      sync.Mutex
	shardID string
	nodes   map[string]*grpc.ClientConn
	db      *badger.DB
	votes   map[string]map[string]bool
}

func (n *SSSCNode) ProposeBlock(ctx context.Context, req *pb.ProposeBlockRequest) (*pb.ProposeBlockResponse, error) {
	if req == nil || req.Block == nil || req.Block.Hash == "" {
		return nil, fmt.Errorf("invalid propose block request")
	}

	// Store proposed block in database
	err := n.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("proposed_block_"+req.Block.Hash), []byte(req.Block.Data))
	})
	if err != nil {
		return nil, err
	}

	// Broadcast proposed block to other nodes
	for _, node := range n.nodes {
		client := pb.NewSSSCNodeClient(node)
		_, err := client.ProposeBlock(ctx, req)
		if err != nil {
			log.Printf("Error broadcasting proposed block to node %s: %v", node.Target(), err)
		}
	}

	return &pb.ProposeBlockResponse{Result: "success"}, nil
}

func (n *SSSCNode) VoteBlock(ctx context.Context, req *pb.VoteBlockRequest) (*pb.VoteBlockResponse, error) {
	if req == nil || req.Block == nil || req.Block.Hash == "" || req.Vote == "" {
		return nil, fmt.Errorf("invalid vote request")
	}

	// Store vote in database
	err := n.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("vote_"+req.Block.Hash+"_"+req.Vote), []byte("true"))
	})
	if err != nil {
		return nil, err
	}

	// Update vote count
	n.mu.Lock()
	if _, ok := n.votes[req.Block.Hash]; !ok {
		n.votes[req.Block.Hash] = make(map[string]bool)
	}
	n.votes[req.Block.Hash][req.Vote] = true
	committed := n.isCommitted(req.Block.Hash)
	n.mu.Unlock()

	// Check if block is committed
	if committed {
		// Commit block to blockchain
		err := n.commitBlock(req.Block)
		if err != nil {
			return nil, err
		}
	}

	return &pb.VoteBlockResponse{Result: "success"}, nil
}

func (n *SSSCNode) isCommitted(blockHash string) bool {
	// Check if block has been committed
	votes, ok := n.votes[blockHash]
	if !ok {
		return false
	}
	if len(votes) >= len(n.nodes)/2 {
		return true
	}
	return false
}

func (n *SSSCNode) commitBlock(block *pb.Block) error {
	// Commit block to blockchain
	err := n.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("block_"+block.Hash), []byte(block.Data))
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db, err := badger.Open("sssc.db")
	if err != nil {
		log.Fatalf("failed to open badger db: %v", err)
	}
	defer db.Close()

	srv := grpc.NewServer()
	pb.RegisterSSSCNodeServer(srv, &SSSCNode{
		shardID: "shard-1",
		nodes:   make(map[string]*grpc.ClientConn),
		db:      db,
		votes:   make(map[string]map[string]bool),
	})

	log.Println("SSSC node listening on port 50053")
	srv.Serve(lis)
}
