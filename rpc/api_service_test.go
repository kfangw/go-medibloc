package rpc_test

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gavv/httpexpect"
	"github.com/medibloc/go-medibloc/consensus/dpos"
	"github.com/medibloc/go-medibloc/core"
	"github.com/medibloc/go-medibloc/rpc"
	"github.com/medibloc/go-medibloc/util/byteutils"
	"github.com/medibloc/go-medibloc/util/testutil"
	"github.com/medibloc/go-medibloc/util/testutil/blockutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPIService_GetAccount(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	bb := blockutil.New(t, testutil.DynastySize).AddKeyPairs(seed.Config.TokenDist).Block(seed.GenesisBlock()).ChildWithTimestamp(dpos.
		NextMintSlot2(time.Now().Unix()))
	payer := seed.Config.TokenDist[testutil.DynastySize]

	tb := bb.Tx()
	tx1 := tb.Nonce(1).StakeTx(payer, 400000000).Build()
	b := bb.ExecuteTx(tx1).SignProposer().Build()

	err := seed.Med.BlockManager().PushBlockDataSync(b.BlockData)
	assert.NoError(t, err)

	e := httpexpect.New(t, testutil.IP2Local(seed.Config.Config.Rpc.HttpListen[0]))

	e.GET("/v1/account").
		WithQuery("type", "tail").
		WithQuery("address", payer.Addr).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("address", payer.Address()).
		ValueEqual("balance", "0").
		ValueEqual("nonce", "1").
		ValueEqual("vesting", "400000000000000000000").
		ValueEqual("voted", []string{}).
		ValueNotEqual("bandwidth", "0").
		ValueEqual("unstaking", "0")

	e.GET("/v1/account").
		WithQuery("address", payer.Addr).
		WithQuery("height", 2).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("address", payer.Address()).
		ValueEqual("balance", "0").
		ValueEqual("nonce", "1").
		ValueEqual("vesting", "400000000000000000000").
		ValueEqual("voted", []string{}).
		ValueNotEqual("bandwidth", "0").
		ValueEqual("unstaking", "0")

	e.GET("/v1/account").
		WithQuery("address", payer.Addr).
		WithQuery("height", 3).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgInvalidBlockHeight)

	e.GET("/v1/account").
		WithQuery("address", payer.Addr).
		WithQuery("type", "genesis").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("vesting", "0")

	e.GET("/v1/account").
		WithQuery("address", payer.Addr).
		WithQuery("type", "confirmed").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("vesting", "0")

	e.GET("/v1/account").
		WithQuery("address", payer.Addr).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgInvalidRequest)

	e.GET("/v1/account").
		WithQuery("address", payer.Addr).
		WithQuery("alias", "Hello World").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgInvalidRequest)

	e.GET("/v1/account").
		WithQuery("type", "confirmed").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgInvalidRequest)
	// TODO @jiseob alias test
}

func TestAPIService_GetBlock(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	bb := blockutil.New(t, 3).AddKeyPairs(seed.Config.TokenDist)
	b := bb.Block(seed.GenesisBlock()).
		ChildWithTimestamp(dpos.NextMintSlot2(time.Now().Unix())).
		Stake().Tx().RandomTx().Execute().SignProposer().Build()

	err := seed.Med.BlockManager().PushBlockDataSync(b.BlockData)
	assert.NoError(t, err)

	e := httpexpect.New(t, testutil.IP2Local(seed.Config.Config.Rpc.HttpListen[0]))

	// The block response should be the same one for designated hash, type, height parameter
	e.GET("/v1/block").
		WithQuery("hash", "0123456789012345678901234567890123456789012345678901234567890123").
		Expect().
		Status(http.StatusInternalServerError).
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgBlockNotFound)

	e.GET("/v1/block").
		WithQuery("hash", byteutils.Bytes2Hex(b.Hash())).
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("height", "2")

	e.GET("/v1/block").
		WithQuery("type", "tail").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("height", "2")

	e.GET("/v1/block").
		WithQuery("height", "2").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("height", "2")

	// The block response should be the genesis one
	e.GET("/v1/block").
		WithQuery("type", "genesis").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("height", "1")

	// Check block response parameters
	e.GET("/v1/block").
		WithQuery("type", "tail").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("height").
		ContainsKey("hash").
		ContainsKey("parent_hash").
		ContainsKey("coinbase").
		ContainsKey("reward").
		ContainsKey("supply").
		ContainsKey("timestamp").
		ContainsKey("chain_id").
		ContainsKey("crypto_alg").
		ContainsKey("hash_alg").
		ContainsKey("sign").
		ContainsKey("accs_root").
		ContainsKey("txs_root").
		ContainsKey("dpos_root").
		ContainsKey("transactions")

	e.GET("/v1/block").
		WithQuery("type", "confirmed").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ValueEqual("height", "1")

	e.GET("/v1/block").
		WithQuery("height", "5").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgInvalidBlockHeight)

	e.GET("/v1/block").
		WithQuery("type", "genesis").
		WithQuery("height", "1").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgInvalidRequest)
}

func TestAPIService_GetBlocks(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	bb := blockutil.New(t, 3).AddKeyPairs(seed.Config.TokenDist)
	b := bb.Block(seed.GenesisBlock()).
		ChildWithTimestamp(dpos.NextMintSlot2(time.Now().Unix())).SignProposer().Build()

	err := seed.Med.BlockManager().PushBlockDataSync(b.BlockData)
	assert.NoError(t, err)

	e := httpexpect.New(t, testutil.IP2Local(seed.Config.Config.Rpc.HttpListen[0]))

	schema := `
	{
		"type":"object",
		"properties":{
			"blocks":{
				"type":"array",
				"items":{
					"type":"object",
					"properties":{
						"height": {
							"type":"string"
						}
					}
				},
				"minItems":2
			}
		}
	}`

	e.GET("/v1/blocks").
		WithQuery("from", "1").
		WithQuery("to", "5").
		Expect().JSON().Schema(schema)

	e.GET("/v1/blocks").
		WithQuery("from", "2").
		WithQuery("to", "1").
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgInvalidRequest)

	e.GET("/v1/blocks").
		WithQuery("from", "1").
		WithQuery("to", rpc.MaxBlocksCount+2).
		Expect().
		Status(http.StatusBadRequest).
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgTooManyBlocksRequest)
}

func TestAPIService_GetCandidates(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	e := httpexpect.New(t, testutil.IP2Local(seed.Config.Config.Rpc.HttpListen[0]))

	e.GET("/v1/candidates").
		Expect().JSON().
		Path("$.candidates").
		Array().Length().Equal(3)
}

func TestAPIService_GetDynasty(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	b := blockutil.New(t, testutil.DynastySize).AddKeyPairs(seed.Config.Dynasties).
		Block(seed.Tail()).Child().SignProposer().Build()

	err := seed.Med.BlockManager().PushBlockDataSync(b.BlockData)
	require.NoError(t, err)

	e := httpexpect.New(t, testutil.IP2Local(seed.Config.Config.Rpc.HttpListen[0]))

	addrs := e.GET("/v1/dynasty").
		Expect().JSON().
		Path("$.addresses").
		Array()
	addrs.Length().Equal(3)
	for _, addr := range addrs.Iter() {
		addr.String().Length().Equal(66)
	}
}

func TestAPIService_GetMedState(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	genesis := seed.Tail()

	bb := blockutil.New(t, 3).AddKeyPairs(seed.Config.TokenDist).Block(seed.GenesisBlock()).ChildWithTimestamp(dpos.
		NextMintSlot2(time.Now().Unix()))
	b := bb.SignProposer().Build()

	err := seed.Med.BlockManager().PushBlockDataSync(b.BlockData)
	assert.NoError(t, err)

	e := httpexpect.New(t, testutil.IP2Local(seed.Config.Config.Rpc.HttpListen[0]))

	e.GET("/v1/node/medstate").
		Expect().JSON().Object().
		ValueEqual("height", "2").
		ValueEqual("LIB", byteutils.Bytes2Hex(genesis.Hash())).
		ValueEqual("tail", byteutils.Bytes2Hex(b.Hash()))
}

func TestAPIService_GetPendingTransactions(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	bb := blockutil.New(t, 3).AddKeyPairs(seed.Config.TokenDist).Block(seed.GenesisBlock()).ChildWithTimestamp(dpos.
		NextMintSlot2(time.Now().Unix()))

	e := httpexpect.New(t, testutil.IP2Local(seed.Config.Config.Rpc.HttpListen[0]))

	e.GET("/v1/transactions/pending").
		Expect().JSON().
		Path("$.transactions").
		Array().Length().Equal(0)

	tb := bb.Tx()

	for i := 1; i <= 10; i++ {
		tx := tb.RandomTx().Build()
		seed.Med.TransactionManager().Push(tx)
		assert.Equal(t, tx, seed.Med.TransactionManager().Get(tx.Hash()))
	}

	e.GET("/v1/transactions/pending").
		Expect().JSON().
		Path("$.transactions").
		Array().Length().Equal(10)
}

func TestAPIService_GetTransaction(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	bb := blockutil.New(t, 3).AddKeyPairs(seed.Config.TokenDist).Block(seed.GenesisBlock()).ChildWithTimestamp(dpos.
		NextMintSlot2(time.Now().Unix())).Stake()
	tx := bb.Tx().RandomTx().Build()
	b := bb.ExecuteTx(tx).SignProposer().Build()

	err := seed.Med.BlockManager().PushBlockDataSync(b.BlockData)
	assert.NoError(t, err)

	e := httpexpect.New(t, testutil.IP2Local(seed.Config.Config.Rpc.HttpListen[0]))

	e.GET("/v1/transaction").
		WithQuery("hash", byteutils.Bytes2Hex(tx.Hash())).
		Expect().
		JSON().Object().
		ValueEqual("hash", byteutils.Bytes2Hex(tx.Hash())).
		ValueEqual("executed", true)

	e.GET("/v1/transaction").
		WithQuery("hash", "0123456789").
		Expect().
		Status(http.StatusNotFound).
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgInvalidTxHash)

	tx = bb.Tx().RandomTx().Build()
	seed.Med.TransactionManager().Push(tx)

	e.GET("/v1/transaction").
		WithQuery("hash", byteutils.Bytes2Hex(tx.Hash())).
		Expect().
		JSON().Object().
		ValueEqual("hash", byteutils.Bytes2Hex(tx.Hash())).
		ValueEqual("executed", false)

	e.GET("/v1/transaction").
		WithQuery("hash", "0123456789012345678901234567890123456789012345678901234567890123").
		Expect().
		JSON().Object().
		ValueEqual("error", rpc.ErrMsgTransactionNotFound)
}

func TestAPIService_HealthCheck(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	e := httpexpect.New(t, testutil.IP2Local(seed.Config.Config.Rpc.HttpListen[0]))

	e.GET("/v1/healthcheck").
		Expect().
		JSON().Object().
		ValueEqual("ok", true)
}

func TestAPIService_SendTransaction(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	bb := blockutil.New(t, 3).AddKeyPairs(seed.Config.TokenDist).Block(seed.GenesisBlock()).ChildWithTimestamp(dpos.NextMintSlot2(time.Now().Unix())).Stake()
	b := bb.SignProposer().Build()

	err := seed.Med.BlockManager().PushBlockDataSync(b.BlockData)
	assert.NoError(t, err)

	payer := seed.Config.TokenDist[3]
	receiver := seed.Config.TokenDist[4]
	tx := bb.Tx().Type(core.TxOpTransfer).To(receiver.Addr).Value(1).Nonce(3).SignPair(payer).Build()

	_, err = rpc.CoreTx2rpcTx(tx, false)
	assert.NoError(t, err)

	e := httpexpect.New(t, testutil.IP2Local(seed.Config.Config.Rpc.HttpListen[0]))

	TX, _ := rpc.CoreTx2rpcTx(tx, false)

	e.POST("/v1/transaction").
		WithJSON(TX).
		Expect().
		JSON().Object().
		ValueEqual("hash", TX.Hash)

	assert.Equal(t, seed.Med.TransactionManager().Get(tx.Hash()).Hash(), tx.Hash())

	TX.Sign = "123"
	e.POST("/v1/transaction").
		WithJSON(TX).
		Expect().
		JSON().Object().ValueEqual("error", rpc.ErrMsgInvalidTransaction)

	TX.Payload = "WRONG PAYLOAD"
	e.POST("/v1/transaction").
		WithJSON(TX).
		Expect().
		JSON().Object().ValueEqual("error", rpc.ErrMsgBuildTransactionFail)

	TX.Value = "abc"
	e.POST("/v1/transaction").
		WithJSON(TX).
		Expect().
		JSON().Object().ValueEqual("error", rpc.ErrMsgInvalidTxValue)
}

type Result struct {
	Topic string
	Hash  string
}

type Data struct {
	Result *Result
}

func TestAPIService_Subscribe(t *testing.T) {
	network := testutil.NewNetwork(t, testutil.DynastySize)
	defer network.Cleanup()
	network.SetLogTestHook()

	seed := network.NewSeedNode()
	seed.Start()
	network.WaitForEstablished()

	bb := blockutil.New(t, 3).AddKeyPairs(seed.Config.TokenDist).Block(seed.GenesisBlock()).ChildWithTimestamp(dpos.NextMintSlot2(time.Now().Unix())).Stake()
	b := bb.SignProposer().Build()

	err := seed.Med.BlockManager().PushBlockDataSync(b.BlockData)
	assert.NoError(t, err)

	tx := make([]*core.Transaction, testutil.DynastySize)
	payer := seed.Config.TokenDist[testutil.DynastySize]
	for i := 0; i < testutil.DynastySize; i++ {
		tx[i] = bb.Tx().Type(core.TxOpTransfer).To(payer.Addr).Value(1).Nonce(uint64(i + 2)).SignPair(payer).Build()
	}

	go func() {
		Client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/subscribe?topics=%s&topics=%s&topics=%s",
			testutil.IP2Local(seed.Config.Config.Rpc.
				HttpListen[0]), core.TopicPendingTransaction, core.TopicTransactionExecutionResult, core.TopicNewTailBlock), nil)
		assert.NoError(t, err)
		req.Header.Set("Accept", "text/event-stream")

		res, err := Client.Do(req)
		assert.NoError(t, err)
		br := bufio.NewReader(res.Body)
		defer res.Body.Close()

		i := 0
		for {
			bs, err := br.ReadBytes('\n')
			if err == io.EOF || i > 6 {
				break
			}
			assert.NoError(t, err)

			data := &Data{
				Result: &Result{},
			}

			err = json.Unmarshal(bs, data)
			assert.NoError(t, err)

			switch data.Result.Topic {
			case core.TopicPendingTransaction:
				assert.Equal(t, data.Result.Hash, byteutils.Bytes2Hex(tx[i%3].Hash()))
			case core.TopicTransactionExecutionResult:
				assert.Equal(t, data.Result.Hash, byteutils.Bytes2Hex(tx[i%3].Hash()))
			case core.TopicNewTailBlock:
				assert.Equal(t, data.Result.Hash, byteutils.Bytes2Hex(b.Hash()))
			}
			i = i + 1
		}
	}()

	for i := 0; i < testutil.DynastySize; i++ {
		// At least 3 seconds for next block
		time.Sleep(1000 * time.Millisecond)
		seed.Med.TransactionManager().Push(tx[i])
	}

	bb = bb.ChildWithTimestamp(dpos.NextMintSlot2(time.Now().Unix()))
	for i := 0; i < testutil.DynastySize; i++ {
		time.Sleep(500 * time.Millisecond)
		bb.ExecuteTx(tx[i])
	}
	b = bb.SignProposer().Build()
	err = seed.Med.BlockManager().PushBlockDataSync(b.BlockData)
	assert.NoError(t, err)
}
