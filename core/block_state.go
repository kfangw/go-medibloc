package core

import (
	"sort"

	"github.com/gogo/protobuf/proto"
	"github.com/medibloc/go-medibloc/common"
	"github.com/medibloc/go-medibloc/common/trie"
	"github.com/medibloc/go-medibloc/core/pb"
	"github.com/medibloc/go-medibloc/storage"
	"github.com/medibloc/go-medibloc/util"
	"github.com/medibloc/go-medibloc/util/byteutils"
	"github.com/medibloc/go-medibloc/util/logging"
	"github.com/sirupsen/logrus"
)

type states struct {
	accState       *AccountStateBatch
	txsState       *TrieBatch
	usageState     *TrieBatch
	recordsState   *TrieBatch
	consensusState *TrieBatch
	candidacyState *TrieBatch

	reservationQueue *ReservationQueue

	storage storage.Storage
}

func newStates(stor storage.Storage) (*states, error) {
	accState, err := NewAccountStateBatch(nil, stor)
	if err != nil {
		return nil, err
	}

	txsState, err := NewTrieBatch(nil, stor)
	if err != nil {
		return nil, err
	}

	usageState, err := NewTrieBatch(nil, stor)
	if err != nil {
		return nil, err
	}

	recordsState, err := NewTrieBatch(nil, stor)
	if err != nil {
		return nil, err
	}

	consensusState, err := NewTrieBatch(nil, stor)
	if err != nil {
		return nil, err
	}

	candidacyState, err := NewTrieBatch(nil, stor)
	if err != nil {
		return nil, err
	}

	reservationQueue := NewEmptyReservationQueue(stor)

	return &states{
		accState:         accState,
		txsState:         txsState,
		usageState:       usageState,
		recordsState:     recordsState,
		consensusState:   consensusState,
		candidacyState:   candidacyState,
		reservationQueue: reservationQueue,
		storage:          stor,
	}, nil
}

func (st *states) Clone() (*states, error) {
	accState, err := NewAccountStateBatch(st.accState.RootHash(), st.storage)
	if err != nil {
		return nil, err
	}

	txsState, err := NewTrieBatch(st.txsState.RootHash(), st.storage)
	if err != nil {
		return nil, err
	}

	usageState, err := NewTrieBatch(st.usageState.RootHash(), st.storage)
	if err != nil {
		return nil, err
	}

	recordsState, err := NewTrieBatch(st.usageState.RootHash(), st.storage)
	if err != nil {
		return nil, err
	}

	consensusState, err := NewTrieBatch(st.consensusState.RootHash(), st.storage)
	if err != nil {
		return nil, err
	}

	candidacyState, err := NewTrieBatch(st.candidacyState.RootHash(), st.storage)
	if err != nil {
		return nil, err
	}

	reservationQueue, err := LoadReservationQueue(st.storage, st.reservationQueue.Hash())
	if err != nil {
		return nil, err
	}

	return &states{
		accState:         accState,
		txsState:         txsState,
		usageState:       usageState,
		recordsState:     recordsState,
		consensusState:   consensusState,
		candidacyState:   candidacyState,
		reservationQueue: reservationQueue,
		storage:          st.storage,
	}, nil
}

func (st *states) BeginBatch() error {
	if err := st.accState.BeginBatch(); err != nil {
		return err
	}
	if err := st.txsState.BeginBatch(); err != nil {
		return err
	}
	if err := st.usageState.BeginBatch(); err != nil {
		return err
	}
	if err := st.recordsState.BeginBatch(); err != nil {
		return err
	}
	if err := st.consensusState.BeginBatch(); err != nil {
		return err
	}
	if err := st.candidacyState.BeginBatch(); err != nil {
		return err
	}
	return st.reservationQueue.BeginBatch()
}

func (st *states) RollBack() error {
	if err := st.accState.RollBack(); err != nil {
		return err
	}
	if err := st.txsState.RollBack(); err != nil {
		return err
	}
	if err := st.usageState.RollBack(); err != nil {
		return err
	}
	if err := st.consensusState.RollBack(); err != nil {
		return err
	}
	if err := st.recordsState.RollBack(); err != nil {
		return err
	}
	if err := st.candidacyState.RollBack(); err != nil {
		return err
	}
	return st.reservationQueue.RollBack()
}

func (st *states) Commit() error {
	if err := st.accState.Commit(); err != nil {
		return err
	}
	if err := st.txsState.Commit(); err != nil {
		return err
	}
	if err := st.usageState.Commit(); err != nil {
		return err
	}
	if err := st.consensusState.Commit(); err != nil {
		return err
	}
	if err := st.recordsState.Commit(); err != nil {
		return err
	}
	if err := st.candidacyState.Commit(); err != nil {
		return err
	}
	return st.reservationQueue.Commit()
}

func (st *states) AccountsRoot() common.Hash {
	return common.BytesToHash(st.accState.RootHash())
}

func (st *states) TransactionsRoot() common.Hash {
	return common.BytesToHash(st.txsState.RootHash())
}

func (st *states) UsageRoot() common.Hash {
	return common.BytesToHash(st.usageState.RootHash())
}

func (st *states) RecordsRoot() common.Hash {
	return common.BytesToHash(st.recordsState.RootHash())
}

func (st *states) ConsensusRoot() common.Hash {
	return common.BytesToHash(st.consensusState.RootHash())
}

func (st *states) CandidacyRoot() common.Hash {
	return common.BytesToHash(st.candidacyState.RootHash())
}

func (st *states) ReservationQueueHash() common.Hash {
	return st.reservationQueue.Hash()
}

func (st *states) LoadAccountsRoot(rootHash common.Hash) error {
	accState, err := NewAccountStateBatch(rootHash.Bytes(), st.storage)
	if err != nil {
		return err
	}
	st.accState = accState
	return nil
}

func (st *states) LoadTransactionsRoot(rootHash common.Hash) error {
	txsState, err := NewTrieBatch(rootHash.Bytes(), st.storage)
	if err != nil {
		return err
	}
	st.txsState = txsState
	return nil
}

func (st *states) LoadUsageRoot(rootHash common.Hash) error {
	usageState, err := NewTrieBatch(rootHash.Bytes(), st.storage)
	if err != nil {
		return err
	}
	st.usageState = usageState
	return nil
}

func (st *states) LoadRecordsRoot(rootHash common.Hash) error {
	recordsState, err := NewTrieBatch(rootHash.Bytes(), st.storage)
	if err != nil {
		return err
	}
	st.recordsState = recordsState
	return nil
}

func (st *states) LoadConsensusRoot(rootHash common.Hash) error {
	consensusState, err := NewTrieBatch(rootHash.Bytes(), st.storage)
	if err != nil {
		return err
	}
	st.consensusState = consensusState
	return nil
}

func (st *states) LoadCandidacyRoot(rootHash common.Hash) error {
	candidacyState, err := NewTrieBatch(rootHash.Bytes(), st.storage)
	if err != nil {
		return err
	}
	st.candidacyState = candidacyState
	return nil
}

func (st *states) LoadReservationQueue(hash common.Hash) error {
	rq, err := LoadReservationQueue(st.storage, hash)
	if err != nil {
		return err
	}
	st.reservationQueue = rq
	return nil
}

func (st *states) GetAccount(address common.Address) (Account, error) {
	return st.accState.GetAccount(address.Bytes())
}

func (st *states) AddBalance(address common.Address, amount *util.Uint128) error {
	return st.accState.AddBalance(address.Bytes(), amount)
}

func (st *states) AddWriter(address common.Address, writer common.Address) error {
	return st.accState.AddWriter(address.Bytes(), writer.Bytes())
}

func (st *states) RemoveWriter(address common.Address, writer common.Address) error {
	return st.accState.RemoveWriter(address.Bytes(), writer.Bytes())
}

func (st *states) SubBalance(address common.Address, amount *util.Uint128) error {
	return st.accState.SubBalance(address.Bytes(), amount)
}

func (st *states) AddRecord(tx *Transaction, hash common.Hash, storage string,
	encKey []byte, seed []byte,
	owner common.Address, writer common.Address) error {
	record := &corepb.Record{
		Hash:      hash.Bytes(),
		Storage:   storage,
		Owner:     tx.from.Bytes(),
		Timestamp: tx.Timestamp(),
		Readers: []*corepb.Reader{
			{
				Address: tx.from.Bytes(),
				EncKey:  encKey,
				Seed:    seed,
			},
		},
	}
	recordBytes, err := proto.Marshal(record)
	if err != nil {
		return err
	}

	if err := st.recordsState.Put(hash.Bytes(), recordBytes); err != nil {
		return err
	}

	return st.accState.AddRecord(tx.from.Bytes(), hash.Bytes())
}

func (st *states) AddRecordReader(tx *Transaction, dataHash common.Hash, reader common.Address, encKey []byte, seed []byte) error {
	recordBytes, err := st.recordsState.Get(dataHash.Bytes())
	if err != nil {
		return err
	}
	pbRecord := new(corepb.Record)
	if err := proto.Unmarshal(recordBytes, pbRecord); err != nil {
		return err
	}
	if byteutils.Equal(pbRecord.Owner, tx.from.Bytes()) == false {
		return ErrTxIsNotFromRecordOwner
	}
	for _, r := range pbRecord.Readers {
		if byteutils.Equal(reader.Bytes(), r.Address) {
			return ErrRecordReaderAlreadyAdded
		}
	}
	pbRecord.Readers = append(pbRecord.Readers, &corepb.Reader{
		Address: reader.Bytes(),
		EncKey:  encKey,
		Seed:    seed,
	})
	b, err := proto.Marshal(pbRecord)
	if err != nil {
		return err
	}
	return st.recordsState.Put(dataHash.Bytes(), b)
}

func (st *states) GetRecord(hash common.Hash) (*corepb.Record, error) {
	recordBytes, err := st.recordsState.Get(hash.Bytes())
	if err != nil {
		return nil, err
	}
	pbRecord := new(corepb.Record)
	if err := proto.Unmarshal(recordBytes, pbRecord); err != nil {
		return nil, err
	}
	return pbRecord, nil
}

func (st *states) incrementNonce(address common.Address) error {
	return st.accState.IncrementNonce(address.Bytes())
}

func (st *states) GetTx(txHash common.Hash) ([]byte, error) {
	return st.txsState.Get(txHash.Bytes())
}

func (st *states) PutTx(txHash common.Hash, txBytes []byte) error {
	return st.txsState.Put(txHash.Bytes(), txBytes)
}

func (st *states) updateUsage(tx *Transaction, blockTime int64) error {
	weekSec := int64(604800)

	if tx.Timestamp() < blockTime-weekSec {
		return ErrTooOldTransaction
	}

	usageBytes, err := st.usageState.Get(tx.from.Bytes())
	switch err {
	case nil:
	case ErrNotFound:
		usage := &corepb.Usage{
			Timestamps: []*corepb.TxTimestamp{
				{
					Hash:      tx.Hash().Bytes(),
					Timestamp: tx.Timestamp(),
				},
			},
		}
		usageBytes, err = proto.Marshal(usage)
		if err != nil {
			return err
		}
		return st.usageState.Put(tx.from.Bytes(), usageBytes)
	default:
		return err
	}

	pbUsage := new(corepb.Usage)
	if err := proto.Unmarshal(usageBytes, pbUsage); err != nil {
		return err
	}

	var idx int
	for idx = range pbUsage.Timestamps {
		if blockTime-weekSec < tx.Timestamp() {
			break
		}
	}
	pbUsage.Timestamps = append(pbUsage.Timestamps[idx:], &corepb.TxTimestamp{Hash: tx.Hash().Bytes(), Timestamp: tx.Timestamp()})
	sort.Slice(pbUsage.Timestamps, func(i, j int) bool {
		return pbUsage.Timestamps[i].Timestamp < pbUsage.Timestamps[j].Timestamp
	})

	pbBytes, err := proto.Marshal(pbUsage)
	if err != nil {
		return err
	}

	return st.usageState.Put(tx.from.Bytes(), pbBytes)
}

func (st *states) GetUsage(addr common.Address) ([]*corepb.TxTimestamp, error) {
	usageBytes, err := st.usageState.Get(addr.Bytes())
	switch err {
	case nil:
	case ErrNotFound:
		return []*corepb.TxTimestamp{}, nil
	default:
		return nil, err
	}

	pbUsage := new(corepb.Usage)
	if err := proto.Unmarshal(usageBytes, pbUsage); err != nil {
		return nil, err
	}
	return pbUsage.Timestamps, nil
}

// Dynasty returns members belonging to the current dynasty.
func (st *states) Dynasty() (members []*common.Address, err error) {
	iter, err := st.consensusState.Iterator(nil)
	if err != nil && err != trie.ErrNotFound {
		logging.WithFields(logrus.Fields{
			"err": err,
		}).Error("Failed to get trie iterator.")
		return nil, err
	}
	if err != nil {
		return members, nil
	}

	exist, err := iter.Next()
	if err != nil {
		return nil, err
	}
	for exist {
		member := common.BytesToAddress(iter.Value())
		members = append(members, &member)
		exist, err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	return members, nil
}

// SetDynasty sets dynasty members.
func (st *states) SetDynasty(miners []*common.Address) error {
	for _, miner := range miners {
		err := st.consensusState.Put(miner.Bytes(), miner.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}

func (st *states) GetCandidate(address common.Address) (*corepb.Candidate, error) {
	candidateBytes, err := st.candidacyState.Get(address.Bytes())
	if err != nil {
		return nil, err
	}
	pbCandidate := new(corepb.Candidate)
	if err := proto.Unmarshal(candidateBytes, pbCandidate); err != nil {
		return nil, err
	}
	return pbCandidate, nil
}

func (st *states) AddCandidate(address common.Address, collateral *util.Uint128) error {
	_, err := st.GetCandidate(address)
	if err != nil && err != ErrNotFound {
		return err
	}
	if err == nil {
		return ErrAlreadyInCandidacy
	}
	if err := st.SubBalance(address, collateral); err != nil {
		return err
	}
	collateralBytes, err := collateral.ToFixedSizeByteSlice()
	if err != nil {
		return err
	}
	pbCandidate := &corepb.Candidate{
		Address:    address.Bytes(),
		Collateral: collateralBytes,
	}
	candidateBytes, err := proto.Marshal(pbCandidate)
	if err != nil {
		return err
	}
	if err := st.candidacyState.Put(address.Bytes(), candidateBytes); err != nil {
		return err
	}
	return nil
}

func (st *states) QuitCandidacy(address common.Address) error {
	candidate, err := st.GetCandidate(address)
	if err != nil {
		return err
	}
	collateral, err := util.NewUint128FromFixedSizeByteSlice(candidate.Collateral)
	if err != nil {
		return err
	}
	if err := st.AddBalance(address, collateral); err != nil {
		return err
	}
	return st.candidacyState.Delete(address.Bytes())
}

// GetReservedTasks returns reserved tasks in reservation queue
func (st *states) GetReservedTasks() []*ReservedTask {
	return st.reservationQueue.Tasks()
}

// AddReservedTask adds a reserved task in reservation queue
func (st *states) AddReservedTask(task *ReservedTask) error {
	return st.reservationQueue.AddTask(task)
}

// PopReservedTask pops reserved tasks which should be processed before 'before'
func (st *states) PopReservedTasks(before int64) []*ReservedTask {
	return st.reservationQueue.PopTasksBefore(before)
}

func (st *states) PeekHeadReservedTask() *ReservedTask {
	return st.reservationQueue.Peek()
}

// WithdrawVesting makes multiple reserved tasks for withdraw a certain amount of vesting
func (st *states) WithdrawVesting(address common.Address, amount *util.Uint128, blockTime int64) error {
	acc, err := st.GetAccount(address)
	if err != nil {
		return err
	}
	if amount.Cmp(acc.Vesting()) > 0 {
		return ErrVestingNotEnough
	}
	splitAmount, err := amount.Div(util.NewUint128FromUint(RtWithdrawNum))
	if err != nil {
		return err
	}
	amountLeft := amount.DeepCopy()
	payload := new(RtWithdraw)
	for i := 0; i < RtWithdrawNum; i++ {
		if amountLeft.Cmp(splitAmount) <= 0 {
			payload, err = NewRtWithdraw(amountLeft)
			if err != nil {
				return err
			}
		} else {
			payload, err = NewRtWithdraw(splitAmount)
			if err != nil {
				return err
			}
		}
		task := NewReservedTask(RtWithdrawType, address, payload, blockTime+int64(i+1)*RtWithdrawInterval)
		if err := st.AddReservedTask(task); err != nil {
			return err
		}
		amountLeft, _ = amountLeft.Sub(splitAmount)
	}
	return nil
}

func (st *states) Vest(address common.Address, amount *util.Uint128) error {
	if err := st.accState.SubBalance(address.Bytes(), amount); err != nil {
		return err
	}
	return st.accState.AddVesting(address.Bytes(), amount)
}

func (st *states) SubVesting(address common.Address, amount *util.Uint128) error {
	return st.accState.SubVesting(address.Bytes(), amount)
}

// BlockState possesses every states a block should have
type BlockState struct {
	*states
	snapshot *states
}

// NewBlockState creates a new block state
func NewBlockState(stor storage.Storage) (*BlockState, error) {
	states, err := newStates(stor)
	if err != nil {
		return nil, err
	}
	return &BlockState{
		states:   states,
		snapshot: nil,
	}, nil
}

// Clone clones block state
func (bs *BlockState) Clone() (*BlockState, error) {
	states, err := bs.states.Clone()
	if err != nil {
		return nil, err
	}
	return &BlockState{
		states:   states,
		snapshot: nil,
	}, nil
}

// BeginBatch begins batch
func (bs *BlockState) BeginBatch() error {
	snapshot, err := bs.states.Clone()
	if err != nil {
		return err
	}
	if err := bs.states.BeginBatch(); err != nil {
		return err
	}
	bs.snapshot = snapshot
	return nil
}

// RollBack rolls back batch
func (bs *BlockState) RollBack() error {
	bs.states = bs.snapshot
	bs.snapshot = nil
	return nil
}

// Commit saves batch updates
func (bs *BlockState) Commit() error {
	if err := bs.states.Commit(); err != nil {
		return err
	}
	bs.snapshot = nil
	return nil
}

// ExecuteTx and update internal states
func (bs *BlockState) ExecuteTx(tx *Transaction) error {
	return tx.ExecuteOnState(bs)
}

// AcceptTransaction and update internal txsStates
func (bs *BlockState) AcceptTransaction(tx *Transaction, blockTime int64) error {
	pbTx, err := tx.ToProto()
	if err != nil {
		return err
	}

	txBytes, err := proto.Marshal(pbTx)
	if err != nil {
		return err
	}

	if err := bs.PutTx(tx.hash, txBytes); err != nil {
		return err
	}

	if err := bs.updateUsage(tx, blockTime); err != nil {
		return err
	}

	return bs.incrementNonce(tx.from)
}

func (bs *BlockState) checkNonce(tx *Transaction) error {
	fromAcc, err := bs.GetAccount(tx.from)
	if err != nil {
		return err
	}

	expectedNonce := fromAcc.Nonce() + 1
	if tx.nonce > expectedNonce {
		return ErrLargeTransactionNonce
	} else if tx.nonce < expectedNonce {
		return ErrSmallTransactionNonce
	}
	return nil
}
