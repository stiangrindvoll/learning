package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	cid "github.com/ipfs/go-cid"
	iaddr "github.com/ipfs/go-ipfs-addr"
	libp2p "github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	mh "github.com/multiformats/go-multihash"

	inet "github.com/libp2p/go-libp2p-net"
)

// IPFS bootstrap nodes. Used to find other peers in the network.
// var bootstrapPeers = []string{
// 	"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
// 	"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
// 	"/ip4/104.236.76.40/tcp/4001/ipfs/QmSoLV4Bbm51jM9C4gDYZQ9Cy3U6aXMJDAbzgu2fzaDs64",
// 	"/ip4/128.199.219.111/tcp/4001/ipfs/QmSoLSafTMBsPKadTEgaXctDQVcqN88CNLHXMkTNwMKPnu",
// 	"/ip4/178.62.158.247/tcp/4001/ipfs/QmSoLer265NRgSp2LA3dPaeykiS1J6DifTC88f5uVQKNAd",
// }
var (
	bootstrapPeers             = []string{"/ip4/81.166.72.135/tcp/4001/ipfs/QmYkQ9SxH71iT6AttBajMqrrkPx1rmm8eBnRuZuqWDLBxB"}
	pid            protocol.ID = "/espgame/0.1.0"
	myPeers        []pstore.PeerInfo
	kadDHT         *dht.IpfsDHT
	peers          = make(map[peer.ID]*bufio.ReadWriter)
)

func createNetwork(rendezvous string) {
	ctx := context.Background()

	host, err := libp2p.New(ctx, libp2p.NATPortMap())
	if err != nil {
		panic(err)
	}
	defer host.Close()

	// Set a function as stream handler.
	// This function is called when a peer initiate a connection and starts a stream with this peer.
	host.SetStreamHandler(pid, handleStream)

	kadDHT, err = dht.New(ctx, host)
	if err != nil {
		panic(err)
	}
	defer kadDHT.Close()

	// Let's connect to the bootstrap nodes first. They will tell us about the other nodes in the network.
	for _, peerAddr := range bootstrapPeers {
		addr, _ := iaddr.ParseString(peerAddr)
		peerinfo, _ := pstore.InfoFromP2pAddr(addr.Multiaddr())

		if err := host.Connect(ctx, *peerinfo); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Connection established with bootstrap node: ", *peerinfo)
		}
	}

	// We use a rendezvous point "meet me here" to announce our location.
	// This is like telling your friends to meet you at the Eiffel Tower.
	v1b := cid.V1Builder{Codec: cid.Raw, MhType: mh.SHA2_256}
	rendezvousPoint, _ := v1b.Sum([]byte(rendezvous))

	fmt.Println("announcing ourselves...")
	tctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if err := kadDHT.Provide(tctx, rendezvousPoint, true); err != nil {
		panic(err)
	}

	// Now, look for others who have announced
	// This is like your friend telling you the location to meet you.
	// 'FindProviders' will return 'PeerInfo' of all the peers which
	// have 'Provide' or announced themselves previously.
	fmt.Println("searching for other peers...")
	tctx, cancel = context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	myPeers, err = kadDHT.FindProviders(tctx, rendezvousPoint)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d peers!\n", len(myPeers))

	for _, p := range myPeers {
		fmt.Println("Peer: ", p)
	}

	for _, p := range myPeers {

		if p.ID == host.ID() || len(p.Addrs) == 0 {
			// No sense connecting to ourselves or if addrs are not available
			fmt.Println("Skipping:", p.ID)
			continue
		}

		stream, err := host.NewStream(ctx, p.ID, pid)
		if err != nil {
			fmt.Println("ERROR: ", err)
		} else {
			peers[p.ID] = bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
			// go writeData(rw)
			go readData(peers[p.ID])
			fmt.Println("Connected to: ", p.ID)
		}

	}

	select {}
}

func handleStream(stream inet.Stream) {
	peer := stream.Conn().RemotePeer()
	log.Println("Got a new stream from:", peer)

	peerInfo := kadDHT.FindLocal(peer)
	fmt.Println("IS in local ps?", peerInfo)

	// Create a buffer stream for non blocking read and write.
	peers[peerInfo.ID] = bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))

	go readData(peers[peerInfo.ID])

	// 'stream' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')

		if err != nil {
			panic(err)
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}
