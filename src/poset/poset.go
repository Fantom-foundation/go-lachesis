package poset

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/hashicorp/golang-lru"

	"github.com/Fantom-foundation/go-lachesis/src/common"
	"github.com/Fantom-foundation/go-lachesis/src/log"
	"github.com/Fantom-foundation/go-lachesis/src/peers"
)

// Core is an interface for interacting with a core.
type Core interface {
	Head() string
	HexID() string
}

//Poset is a DAG of Events. It also contains methods to extract a consensus
//order of Events and map them onto a blockchain.
type Poset struct {
	Participants            *peers.Peers     //[public key] => id
	Store                   Store            //store of Events, Rounds, and Blocks
	UndeterminedEvents      []string         //[index] => hash . FIFO queue of Events whose consensus order is not yet determined
	PendingRounds           []*pendingRound  //FIFO queue of Rounds which have not attained consensus yet
	LastConsensusRound      *int64           //index of last consensus round
	FirstConsensusRound     *int64           //index of first consensus round (only used in tests)
	AnchorBlock             *int64           //index of last block with enough signatures
	LastCommitedRoundEvents int              //number of events in round before LastConsensusRound
	SigPool                 []BlockSignature //Pool of Block signatures that need to be processed
	consensusTransactions   uint64           //number of consensus transactions
	pendingLoadedEvents     int64            //number of loaded events that are not yet committed
	commitCh                chan Block       //channel for committing Blocks
	topologicalIndex        int64            //counter used to order events in topological order (only local)
	superMajority           int
	trustCount              int
	core                    Core

	ancestorCache     *lru.Cache
	selfAncestorCache *lru.Cache
	stronglySeeCache  *lru.Cache
	roundCache        *lru.Cache
	timestampCache    *lru.Cache

	logger *logrus.Entry

	undeterminedEventsLocker      sync.RWMutex
	pendingLoadedEventsLocker     sync.RWMutex
	firstLastConsensusRoundLocker sync.RWMutex
	consensusTransactionsLocker   sync.RWMutex
}

//NewPoset instantiates a Poset from a list of participants, underlying
//data store and commit channel
func NewPoset(participants *peers.Peers, store Store, commitCh chan Block, logger *logrus.Entry) *Poset {
	if logger == nil {
		log := logrus.New()
		log.Level = logrus.DebugLevel
		lachesis_log.NewLocal(log, log.Level.String())
		logger = logrus.NewEntry(log)
	}

	superMajority := 2*participants.Len()/3 + 1
	trustCount := int(math.Ceil(float64(participants.Len()) / float64(3)))

	cacheSize := store.CacheSize()
	ancestorCache, err := lru.New(cacheSize)
	if err != nil {
		logger.Fatal("Unable to init Poset.ancestorCache")
	}
	selfAncestorCache, err := lru.New(cacheSize)
	if err != nil {
		logger.Fatal("Unable to init Poset.selfAncestorCache")
	}
	stronglySeeCache, err :=  lru.New(cacheSize)
	if err != nil {
		logger.Fatal("Unable to init Poset.stronglySeeCache")
	}
	roundCache, err :=        lru.New(cacheSize)
	if err != nil {
		logger.Fatal("Unable to init Poset.roundCache")
	}
	timestampCache, err :=    lru.New(cacheSize)
	if err != nil {
		logger.Fatal("Unable to init Poset.timestampCache")
	}
	poset := Poset{
		Participants:      participants,
		Store:             store,
		commitCh:          commitCh,
		ancestorCache:     ancestorCache,
		selfAncestorCache: selfAncestorCache,
		stronglySeeCache:  stronglySeeCache,
		roundCache:        roundCache,
		timestampCache:    timestampCache,
		logger:            logger,
		superMajority:     superMajority,
		trustCount:        trustCount,
	}

	participants.OnNewPeer(func(peer *peers.Peer) {
		poset.superMajority = 2*participants.Len()/3 + 1
		poset.trustCount = int(math.Ceil(float64(participants.Len()) / float64(3)))
	})

	return &poset
}

// SetCore sets a core for poset.
func (p *Poset) SetCore(core Core) {
	p.core = core
}

/*******************************************************************************
Private Methods
*******************************************************************************/

//true if y is an ancestor of x
func (p *Poset) ancestor(x, y string) (bool, error) {
	if c, ok := p.ancestorCache.Get(Key{x, y}); ok {
		return c.(bool), nil
	}

	if len(x) == 0 || len(y) == 0 {
		return false, nil
	}

	a, err := p.ancestor2(x, y)
	if err != nil {
		return false, err
	}
	p.ancestorCache.Add(Key{x, y}, a)
	return a, nil
}

func (p *Poset) ancestor2(x, y string) (bool, error) {
	if x == y {
		return true, nil
	}

	ex, err := p.Store.GetEvent(x)
	if err != nil {
		roots, err2 := p.Store.RootsBySelfParent()
		if err2 != nil {
			return false, err2
		}
		for _, root := range roots {
			if other, ok := root.Others[y]; ok {
				return other.Hash == x, nil
			}
		}
		return false, nil
	}
	if lamportDiff, err := p.lamportTimestampDiff(x, y); err != nil || lamportDiff > 0 {
		return false, err
	}

	ey, err := p.Store.GetEvent(y)
	if err != nil {
		// check y roots
		roots, err2 := p.Store.RootsBySelfParent()
		if err2 != nil {
			return false, err2
		}
		if root, ok := roots[y]; ok {
			peer, ok := p.Participants.GetById(root.SelfParent.CreatorID)
			if !ok {
				return false, fmt.Errorf("ancestor2: participant not found")
			}
			yCreator := peer.PubKeyHex
			if ex.Creator() == yCreator {
				return ex.Index() >= root.SelfParent.Index, nil
			}
		} else {
			return false, nil
		}
	} else {
		// check if creators are equals and check indexes
		if ex.Creator() == ey.Creator() {
			return ex.Index() >= ey.Index(), nil
		}
	}

	res, err := p.ancestor(ex.SelfParent(), y)
	if err != nil {
		return false, err
	}

	if res {
		return true, nil
	}

	return p.ancestor(ex.OtherParent(), y)
}

//true if y is a self-ancestor of x
func (p *Poset) selfAncestor(x, y string) (bool, error) {
	if c, ok := p.selfAncestorCache.Get(Key{x, y}); ok {
		return c.(bool), nil
	}
	if len(x) == 0 || len(y) == 0 {
		return false, nil
	}
	a, err := p.selfAncestor2(x, y)
	if err != nil {
		return false, err
	}
	p.selfAncestorCache.Add(Key{x, y}, a)
	return a, nil
}

func (p *Poset) selfAncestor2(x, y string) (bool, error) {
	if x == y {
		return true, nil
	}
	ex, err := p.Store.GetEvent(x)
	if err != nil {
		roots, err := p.Store.RootsBySelfParent()
		if err != nil {
			return false, err
		}
		if root, ok := roots[x]; ok {
			if root.SelfParent.Hash == y {
				return true, nil
			}
		}
		return false, err
	}

	ey, err := p.Store.GetEvent(y)
	if err != nil {
		roots, err2 := p.Store.RootsBySelfParent()
		if err2 != nil {
			return false, err2
		}
		if root, ok := roots[y]; ok {
			peer, ok := p.Participants.GetById(root.SelfParent.CreatorID)
			if !ok {
				return false, fmt.Errorf("selfAncestor2: participant not found")
			}
			yCreator := peer.PubKeyHex
			if ex.Creator() == yCreator {
				return ex.Index() >= root.SelfParent.Index, nil
			}
		}
	} else {
		if ex.Creator() == ey.Creator() {
			return ex.Index() >= ey.Index(), nil
		}
	}

	return false, nil
}

//true if x sees y
func (p *Poset) see(x, y string) (bool, error) {
	return p.ancestor(x, y)
	//it is not necessary to detect forks because we assume that the InsertEvent
	//function makes it impossible to insert two Events at the same height for
	//the same participant.
}

//true if x strongly sees y
func (p *Poset) stronglySee(x, y string) (bool, error) {
	if len(x) == 0 || len(y) == 0 {
		return false, nil
	}

	if c, ok := p.stronglySeeCache.Get(Key{x, y}); ok {
		return c.(bool), nil
	}
	ss, err := p.stronglySee2(x, y)
	if err != nil {
		return false, err
	}
	p.stronglySeeCache.Add(Key{x, y}, ss)
	return ss, nil
}

// Possible improvement: Populate the cache for upper and downer events
// that also stronglySee y
func (p *Poset) stronglySee2(x, y string) (bool, error) {
	sentinels := make(map[string]bool)

	if err := p.MapSentinels(x, y, sentinels); err != nil {
		return false, err
	}

	return len(sentinels) >= p.superMajority, nil
}

// participants in x's ancestry that see y
func (p *Poset) MapSentinels(x, y string, sentinels map[string]bool) error {
	if x == "" {
		return nil
	}

	if see, err := p.see(x, y); err != nil || !see {
		return err
	}

	ex, err := p.Store.GetEvent(x)

	if err != nil {
		roots, err2 := p.Store.RootsBySelfParent()

		if err2 != nil {
			return err2
		}

		if root, ok := roots[x]; ok {
			creator, ok := p.Participants.GetById(root.SelfParent.CreatorID)
			if !ok {
				return fmt.Errorf("MapSentinels: participant not found")
			}

			sentinels[creator.PubKeyHex] = true

			return nil
		}

		return err
	}

	creator, ok := p.Participants.GetById(ex.CreatorID())
	if !ok {
		return fmt.Errorf("MapSentinels: participant not found, two")
	}
	sentinels[creator.PubKeyHex] = true

	if x == y {
		return nil
	}

	if err := p.MapSentinels(ex.OtherParent(), y, sentinels); err != nil {
		return err
	}

	return p.MapSentinels(ex.SelfParent(), y, sentinels)
}

func (p *Poset) round(x string) (int64, error) {
	if c, ok := p.roundCache.Get(x); ok {
		return c.(int64), nil
	}
	r, err := p.round2(x)
	if err != nil {
		return -1, err
	}
	p.roundCache.Add(x, r)
	return r, nil
}

func (p *Poset) round2(x string) (int64, error) {

	/*
		x is the Root
		Use Root.SelfParent.Round
	*/
	rootsBySelfParent, _ := p.Store.RootsBySelfParent()
	if r, ok := rootsBySelfParent[x]; ok {
		return r.SelfParent.Round, nil
	}

	ex, err := p.Store.GetEvent(x)
	if err != nil {
		return math.MinInt64, err
	}

	root, err := p.Store.GetRoot(ex.Creator())
	if err != nil {
		return math.MinInt64, err
	}

	/*
		The Event is directly attached to the Root.
	*/
	if ex.SelfParent() == root.SelfParent.Hash {
		//Root is authoritative EXCEPT if other-parent is not in the root
		if other, ok := root.Others[ex.Hex()]; (ex.OtherParent() == "") ||
			(ok && other.Hash == ex.OtherParent()) {

			return root.NextRound, nil
		}
	}

	/*
		The Event's parents are "normal" Events.
		Use the whitepaper formula: parentRound + roundInc
	*/
	spRound, err := p.round(ex.SelfParent())
	if err != nil {
		return math.MinInt64, err
	}
	var parentRound = spRound
	var opRound int64

	if ex.OtherParent() != "" {
		//XXX
		if other, ok := root.Others[ex.Hex()]; ok && other.Hash == ex.OtherParent() {
			opRound = root.NextRound
		} else {
			opRound, err = p.round(ex.OtherParent())
			if err != nil {
				return math.MinInt64, err
			}
		}

		if opRound > parentRound {
			var (
				found           bool
				seeOpRoundRoots int64
			)

			// if in a flag table there are witnesses of the current round, then
			// current round is other parent round.
			ws := p.Store.RoundWitnesses(opRound)
			ft, _ := ex.GetFlagTable()
			for k := range ft {
				for _, w := range ws {
					if w == k && w != ex.Hex() {
						see, err := p.see(ex.Hex(), w)
						if err != nil {
							return math.MinInt32, err
						}

						if see {
							if !found {
								found = true
							}
							seeOpRoundRoots++
						}
					}
				}
			}

			if seeOpRoundRoots >= int64(p.superMajority) {
				return opRound + 1, nil
			}

			if found {
				return opRound, nil
			}

			parentRound = opRound
		}
	}

	ws := p.Store.RoundWitnesses(parentRound)

	isSee := func(poset *Poset, root string, witnesses []string) bool {
		for _, w := range ws {
			if w == root && w != ex.Hex() {
				see, err := poset.see(ex.Hex(), w)
				if err != nil {
					return false
				}
				if see {
					return true
				}
			}
		}
		return false
	}

	// check wp
	if len(ex.Message.WitnessProof) >= p.superMajority {
		count := 0

		for _, root := range ex.Message.WitnessProof {
			if isSee(p, root, ws) {
				count++
			}
		}

		if count >= p.superMajority {
			return parentRound + 1, err
		}
	}

	// check ft
	ft, _ := ex.GetFlagTable()
	if len(ft) >= p.superMajority {
		count := 0

		for root := range ft {
			if isSee(p, root, ws) {
				count++
			}
		}

		if count >= p.superMajority {
			return parentRound + 1, err
		}
	}

	return parentRound, nil
}

// witness if is true then x is a witness (first event of a round for the owner)
func (p *Poset) witness(x string) (bool, error) {
	ex, err := p.Store.GetEvent(x)
	if err != nil {
		return false, err
	}
	if ex.OtherParent() == "" {
		return true, nil
	}

	xRound, err := p.round(x)
	if err != nil {
		return false, err
	}
	spRound, err := p.round(ex.SelfParent())
	if err != nil {
		return false, err
	}
	return xRound > spRound, nil
}

func (p *Poset) roundReceived(x string) (int64, error) {

	ex, err := p.Store.GetEvent(x)
	if err != nil {
		return -1, err
	}

	return ex.Message.RoundReceived, nil
}

func (p *Poset) lamportTimestamp(x string) (int64, error) {
	if c, ok := p.timestampCache.Get(x); ok {
		return c.(int64), nil
	}
	r, err := p.lamportTimestamp2(x)
	if err != nil {
		return -1, err
	}
	p.timestampCache.Add(x, r)
	return r, nil
}

func (p *Poset) lamportTimestamp2(x string) (int64, error) {
	/*
		x is the Root
		User Root.SelfParent.LamportTimestamp
	*/
	rootsBySelfParent, _ := p.Store.RootsBySelfParent()
	if r, ok := rootsBySelfParent[x]; ok {
		return r.SelfParent.LamportTimestamp, nil
	}

	ex, err := p.Store.GetEvent(x)
	if err != nil {
		return math.MinInt64, err
	}

	//We are going to need the Root later
	root, err := p.Store.GetRoot(ex.Creator())
	if err != nil {
		return math.MinInt64, err
	}

	plt := int64(math.MinInt64)
	//If it is the creator's first Event, use the corresponding Root
	if ex.SelfParent() == root.SelfParent.Hash {
		plt = root.SelfParent.LamportTimestamp
	} else {
		t, err := p.lamportTimestamp(ex.SelfParent())
		if err != nil {
			return math.MinInt64, err
		}
		plt = t
	}

	if ex.OtherParent() != "" {
		opLT := int64(math.MinInt64)
		if _, err := p.Store.GetEvent(ex.OtherParent()); err == nil {
			//if we know the other-parent, fetch its Round directly
			t, err := p.lamportTimestamp(ex.OtherParent())
			if err != nil {
				return math.MinInt64, err
			}
			opLT = t
		} else if other, ok := root.Others[x]; ok && other.Hash == ex.OtherParent() {
			//we do not know the other-parent but it is referenced  in Root.Others
			//we use the Root's LamportTimestamp
			opLT = other.LamportTimestamp
		}

		if opLT > plt {
			plt = opLT
		}
	}

	return plt + 1, nil
}

// lamport(y) - lamport(x)
func (p *Poset) lamportTimestampDiff(x, y string) (int64, error) {
	xlt, err := p.lamportTimestamp(x)
	if err != nil {
		return 0, err
	}
	ylt, err := p.lamportTimestamp(y)
	if err != nil {
		return 0, err
	}
	return ylt - xlt, nil
}

//round(x) - round(y)
func (p *Poset) roundDiff(x, y string) (int64, error) {

	xRound, err := p.round(x)
	if err != nil {
		return math.MinInt64, fmt.Errorf("event %s has negative round", x)
	}

	yRound, err := p.round(y)
	if err != nil {
		return math.MinInt64, fmt.Errorf("event %s has negative round", y)
	}

	return xRound - yRound, nil
}

//Check the SelfParent is the Creator's last known Event
func (p *Poset) checkSelfParent(event Event) error {
	selfParent := event.SelfParent()
	creator := event.Creator()

	creatorLastKnown, _, err := p.Store.LastEventFrom(creator)

	p.logger.WithFields(logrus.Fields{
		"selfParent":       selfParent,
		"creator":          creator,
		"creatorLastKnown": creatorLastKnown,
		"event":            event.Hex(),
	}).Debugf("checkSelfParent")

	if err != nil {
		return err
	}

	selfParentLegit := selfParent == creatorLastKnown

	if !selfParentLegit {
		return fmt.Errorf("self-parent not last known event by creator")
	}

	return nil
}

//Check if we know the OtherParent
func (p *Poset) checkOtherParent(event Event) error {
	otherParent := event.OtherParent()
	if otherParent != "" {
		//Check if we have it
		_, err := p.Store.GetEvent(otherParent)
		if err != nil {
			//it might still be in the Root
			root, err := p.Store.GetRoot(event.Creator())
			if err != nil {
				return err
			}
			other, ok := root.Others[event.Hex()]
			if ok && other.Hash == event.OtherParent() {
				return nil
			}
			return fmt.Errorf("other-parent not known")
		}
	}
	return nil
}

func (p *Poset) createSelfParentRootEvent(ev Event) (RootEvent, error) {
	sp := ev.SelfParent()
	spLT, err := p.lamportTimestamp(sp)
	if err != nil {
		return RootEvent{}, err
	}
	spRound, err := p.round(sp)
	if err != nil {
		return RootEvent{}, err
	}
	peer, ok := p.Participants.GetByPubKey(ev.Creator())
	if !ok {
		return RootEvent{}, fmt.Errorf("createSelfParentRootEvent: participant not found")
	}
	selfParentRootEvent := RootEvent{
		Hash:             sp,
		CreatorID:        peer.ID,
		Index:            ev.Index() - 1,
		LamportTimestamp: spLT,
		Round:            spRound,
		//FlagTable:ev.FlagTable,
		//flags:ev.flags,
	}
	return selfParentRootEvent, nil
}

func (p *Poset) createOtherParentRootEvent(ev Event) (RootEvent, error) {

	op := ev.OtherParent()

	//it might still be in the Root
	root, err := p.Store.GetRoot(ev.Creator())
	if err != nil {
		return RootEvent{}, err
	}
	if other, ok := root.Others[ev.Hex()]; ok && other.Hash == op {
		return *other, nil
	}

	otherParent, err := p.Store.GetEvent(op)
	if err != nil {
		return RootEvent{}, err
	}
	opLT, err := p.lamportTimestamp(op)
	if err != nil {
		return RootEvent{}, err
	}
	opRound, err := p.round(op)
	if err != nil {
		return RootEvent{}, err
	}
	peer, ok := p.Participants.GetByPubKey(otherParent.Creator())
	if !ok {
		return RootEvent{}, fmt.Errorf("createOtherParentRootEvent: participant not found")
	}
	otherParentRootEvent := RootEvent{
		Hash:             op,
		CreatorID:        peer.ID,
		Index:            otherParent.Index(),
		LamportTimestamp: opLT,
		Round:            opRound,
	}
	return otherParentRootEvent, nil

}

func (p *Poset) createRoot(ev Event) (Root, error) {

	evRound, err := p.round(ev.Hex())
	if err != nil {
		return Root{}, err
	}

	/*
		SelfParent
	*/
	selfParentRootEvent, err := p.createSelfParentRootEvent(ev)
	if err != nil {
		return Root{}, err
	}

	/*
		OtherParent
	*/
	var otherParentRootEvent *RootEvent
	if ev.OtherParent() != "" {
		opre, err := p.createOtherParentRootEvent(ev)
		if err != nil {
			return Root{}, err
		}
		otherParentRootEvent = &opre
	}

	root := Root{
		NextRound:  evRound,
		SelfParent: &selfParentRootEvent,
		Others:     map[string]*RootEvent{},
	}

	if otherParentRootEvent != nil {
		root.Others[ev.Hex()] = otherParentRootEvent
	}

	return root, nil
}

func (p *Poset) SetWireInfo(event *Event) error {
	return p.setWireInfo(event)
}
func (p *Poset) SetWireInfoAndSign(event *Event, privKey *ecdsa.PrivateKey) error {
	if err := p.setWireInfo(event); err != nil {
		return err
	}
	return event.Sign(privKey)
}

func (p *Poset) setWireInfo(event *Event) error {
	selfParentIndex := int64(-1)
	otherParentCreatorID := int64(-1)
	otherParentIndex := int64(-1)

	//could be the first Event inserted for this creator. In this case, use Root
	if lf, isRoot, _ := p.Store.LastEventFrom(event.Creator()); isRoot && lf == event.SelfParent() {
		root, err := p.Store.GetRoot(event.Creator())
		if err != nil {
			return err
		}
		selfParentIndex = root.SelfParent.Index
	} else {
		selfParent, err := p.Store.GetEvent(event.SelfParent())
		if err != nil {
			return err
		}
		selfParentIndex = selfParent.Index()
	}

	if event.OtherParent() != "" {
		//Check Root then regular Events
		root, err := p.Store.GetRoot(event.Creator())
		if err != nil {
			return err
		}
		if other, ok := root.Others[event.Hex()]; ok && other.Hash == event.OtherParent() {
			otherParentCreatorID = other.CreatorID
			otherParentIndex = other.Index
		} else {
			otherParent, err := p.Store.GetEvent(event.OtherParent())
			if err != nil {
				return err
			}
			peer, ok := p.Participants.GetByPubKey(otherParent.Creator())
			if !ok {
				return fmt.Errorf("setWireInfo: participant not found")
			}
			otherParentCreatorID = peer.ID
			otherParentIndex = otherParent.Index()
		}
	}

	peer, ok := p.Participants.GetByPubKey(event.Creator())
	if !ok {
		return fmt.Errorf("setWireInfo: participant not found, two")
	}
	event.SetWireInfo(selfParentIndex,
		otherParentCreatorID,
		otherParentIndex,
		peer.ID)

	return nil
}

func (p *Poset) updatePendingRounds(decidedRounds map[int64]int64) {
	for _, ur := range p.PendingRounds {
		if _, ok := decidedRounds[ur.Index]; ok {
			ur.Decided = true
		}
	}
}

//Remove processed Signatures from SigPool
func (p *Poset) removeProcessedSignatures(processedSignatures map[int64]bool) {
	var newSigPool []BlockSignature
	for _, bs := range p.SigPool {
		if _, ok := processedSignatures[bs.Index]; !ok {
			newSigPool = append(newSigPool, bs)
		}
	}
	p.SigPool = newSigPool
}

/*******************************************************************************
Public Methods
*******************************************************************************/

//InsertEvent attempts to insert an Event in the DAG. It verifies the signature,
//checks the ancestors are known, and prevents the introduction of forks.
func (p *Poset) InsertEvent(event Event, setWireInfo bool) error {
	//verify signature
	if ok, err := event.Verify(); !ok {
		if err != nil {
			return err
		}

		p.logger.WithFields(logrus.Fields{
			"event":      event,
			"creator":    event.Creator(),
			"selfParent": event.SelfParent(),
			"index":      event.Index(),
			"hex":        event.Hex(),
		}).Debugf("Invalid Event signature")

		return fmt.Errorf("invalid Event signature")
	}

	if err := p.checkSelfParent(event); err != nil {
		return fmt.Errorf("CheckSelfParent: %s", err)
	}

	if err := p.checkOtherParent(event); err != nil {
		return fmt.Errorf("CheckOtherParent: %s", err)
	}

	event.Message.TopologicalIndex = p.topologicalIndex
	p.topologicalIndex++

	if setWireInfo {
		if err := p.setWireInfo(&event); err != nil {
			return fmt.Errorf("SetWireInfo: %s", err)
		}
	}

	if err := p.Store.SetEvent(event); err != nil {
		return fmt.Errorf("SetEvent: %s", err)
	}

	p.undeterminedEventsLocker.Lock()
	p.UndeterminedEvents = append(p.UndeterminedEvents, event.Hex())
	p.undeterminedEventsLocker.Unlock()

	if event.IsLoaded() {
		p.pendingLoadedEventsLocker.Lock()
		p.pendingLoadedEvents++
		p.pendingLoadedEventsLocker.Unlock()
	}

	blockSignatures := make([]BlockSignature, len(event.BlockSignatures()))
	for i, v := range event.BlockSignatures() {
		blockSignatures[i] = *v
	}
	p.SigPool = append(p.SigPool, blockSignatures...)

	return nil
}

/*
DivideRounds assigns a Round and LamportTimestamp to Events, and flags them as
witnesses if necessary. Pushes Rounds in the PendingRounds queue if necessary.
*/
func (p *Poset) DivideRounds() error {

	p.undeterminedEventsLocker.RLock()
	defer p.undeterminedEventsLocker.RUnlock()

	for _, hash := range p.UndeterminedEvents {

		ev, err := p.Store.GetEvent(hash)
		if err != nil {
			return err
		}

		updateEvent := false

		/*
		   Compute Event's round, update the corresponding Round object, and
		   add it to the PendingRounds queue if necessary.
		*/
		if ev.Message.Round == RoundNIL {

			roundNumber, err := p.round(hash)
			if err != nil {
				return err
			}

			ev.SetRound(roundNumber)
			updateEvent = true

			roundInfo, err := p.Store.GetRound(roundNumber)
			if err != nil && !common.Is(err, common.KeyNotFound) {
				return err
			}

			/*
				Why the lower bound?
				Normally, once a Round has attained consensus, it is impossible for
				new Events from a previous Round to be inserted; the lower bound
				appears redundant. This is the case when the poset grows
				linearly, without jumps, which is what we intend by 'Normally'.
				But the Reset function introduces a discontinuity  by jumping
				straight to a specific place in the poset. This technique relies
				on a base layer of Events (the corresponding Frame's Events) for
				other Events to be added on top, but the base layer must not be
				reprocessed.
			*/
			if !roundInfo.Message.Queued && roundNumber >= p.GetLastConsensusRound() {

				p.PendingRounds = append(p.PendingRounds, &pendingRound{roundNumber, false})
				roundInfo.Message.Queued = true
			}

			witness, err := p.witness(hash)
			if err != nil {
				return err
			}
			roundInfo.AddEvent(hash, witness)

			err = p.Store.SetRound(roundNumber, roundInfo)
			if err != nil {
				return err
			}

			if witness {
				// if event is self head
				if p.core != nil && ev.Hex() == p.core.Head() &&
					ev.Creator() == p.core.HexID() {

					replaceFlagTable := func(event *Event, round int64) {
						ft := make(map[string]int64)
						ws := p.Store.RoundWitnesses(round)
						for _, v := range ws {
							ft[v] = 1
						}
						event.ReplaceFlagTable(ft)
					}

					// special case
					if ev.GetRound() == 0 {
						replaceFlagTable(&ev, 0)
						root, err := p.Store.GetRoot(ev.Creator())
						if err != nil {
							return err
						}
						ev.Message.WitnessProof = []string{root.SelfParent.Hash}
					} else {
						replaceFlagTable(&ev, ev.GetRound())
						roots := p.Store.RoundWitnesses(ev.GetRound() - 1)
						ev.Message.WitnessProof = roots
					}
				}
			}
		}

		/*
			Compute the Event's LamportTimestamp
		*/
		if ev.Message.LamportTimestamp == LamportTimestampNIL {

			lamportTimestamp, err := p.lamportTimestamp(hash)
			if err != nil {
				return err
			}

			ev.SetLamportTimestamp(lamportTimestamp)
			updateEvent = true
		}

		if updateEvent {
			if ev.CreatorID() == 0 {
				p.setWireInfo(&ev)
			}
			p.Store.SetEvent(ev)
		}
	}

	return nil
}

//DecideFame decides if witnesses are famous
func (p *Poset) DecideFame() error {

	//Initialize the vote map
	votes := make(map[string]map[string]bool) //[x][y]=>vote(x,y)
	setVote := func(votes map[string]map[string]bool, x, y string, vote bool) {
		if votes[x] == nil {
			votes[x] = make(map[string]bool)
		}
		votes[x][y] = vote
	}

	decidedRounds := map[int64]int64{} // [round number] => index in p.PendingRounds
	c := 3

	for pos, r := range p.PendingRounds {
		roundIndex := r.Index
		roundInfo, err := p.Store.GetRound(roundIndex)
		if err != nil {
			return err
		}
		for _, x := range roundInfo.Witnesses() {
			if roundInfo.IsDecided(x) {
				continue
			}
		VOTE_LOOP:
			for j := roundIndex + 1; j <= p.Store.LastRound(); j++ {
				for _, y := range p.Store.RoundWitnesses(j) {
					diff := j - roundIndex
					if diff == 1 {
						ycx, err := p.see(y, x)
						if err != nil {
							return err
						}
						setVote(votes, y, x, ycx)
					} else {
						//count votes
						var ssWitnesses []string
						for _, w := range p.Store.RoundWitnesses(j - 1) {
							ss, err := p.stronglySee(y, w)
							if err != nil {
								return err
							}
							if ss {
								ssWitnesses = append(ssWitnesses, w)
							}
						}
						yays := 0
						nays := 0
						for _, w := range ssWitnesses {
							if votes[w][x] {
								yays++
							} else {
								nays++
							}
						}
						v := false
						t := nays
						if yays >= nays {
							v = true
							t = yays
						}

						//normal round
						if math.Mod(float64(diff), float64(c)) > 0 {
							if t >= p.superMajority {
								roundInfo.SetFame(x, v)
								setVote(votes, y, x, v)
								break VOTE_LOOP //break out of j loop
							} else {
								setVote(votes, y, x, v)
							}
						} else { //coin round
							if t >= p.superMajority {
								setVote(votes, y, x, v)
							} else {
								setVote(votes, y, x, middleBit(y)) //middle bit of y's hash
							}
						}
					}
				}
			}
		}

		err = p.Store.SetRound(roundIndex, roundInfo)
		if err != nil {
			return err
		}

		if roundInfo.WitnessesDecided() {
			decidedRounds[roundIndex] = int64(pos)
		}

	}

	p.updatePendingRounds(decidedRounds)
	return nil
}

//DecideRoundReceived assigns a RoundReceived to undetermined events when they
//reach consensus
func (p *Poset) DecideRoundReceived() error {

	p.undeterminedEventsLocker.Lock()
	defer p.undeterminedEventsLocker.Unlock()

	var newUndeterminedEvents []string

	/* From whitepaper - 18/03/18
	   "[...] An event is said to be “received” in the first round where all the
	   unique famous witnesses have received it, if all earlier rounds have the
	   fame of all witnesses decided"
	*/
	for _, x := range p.UndeterminedEvents {

		received := false
		r, err := p.round(x)
		if err != nil {
			return err
		}

		for i := r + 1; i <= p.Store.LastRound(); i++ {

			tr, err := p.Store.GetRound(i)
			if err != nil {
				//Can happen after a Reset/FastSync
				if r < p.GetLastConsensusRound() {
					received = true
					break
				}
				return err
			}

			//We are looping from earlier to later rounds; so if we encounter
			//one round with undecided witnesses, we are sure that this event
			//is not "received". Break out of i loop
			if !(tr.WitnessesDecided()) {
				break
			}

			fws := tr.FamousWitnesses()
			//set of famous witnesses that see x
			var s []string
			for _, w := range fws {
				see, err := p.see(w, x)
				if err != nil {
					return err
				}
				if see {
					s = append(s, w)
				}
			}

			if len(s) == len(fws) && len(s) > 0 {

				received = true

				ex, err := p.Store.GetEvent(x)
				if err != nil {
					return err
				}
				ex.SetRoundReceived(i)

				err = p.Store.SetEvent(ex)
				if err != nil {
					return err
				}

				tr.SetConsensusEvent(x)
				err = p.Store.SetRound(i, tr)
				if err != nil {
					return err
				}

				//break out of i loop
				break
			}

		}

		if !received {
			newUndeterminedEvents = append(newUndeterminedEvents, x)
		}
	}

	p.UndeterminedEvents = newUndeterminedEvents

	return nil
}

//ProcessDecidedRounds takes Rounds whose witnesses are decided, computes the
//corresponding Frames, maps them into Blocks, and commits the Blocks via the
//commit channel
func (p *Poset) ProcessDecidedRounds() error {

	//Defer removing processed Rounds from the PendingRounds Queue
	processedIndex := 0
	defer func() {
		p.PendingRounds = p.PendingRounds[processedIndex:]
	}()

	for _, r := range p.PendingRounds {

		//Although it is possible for a Round to be 'decided' before a previous
		//round, we should NEVER process a decided round before all the previous
		//rounds are processed.
		if !r.Decided {
			break
		}

		//This is similar to the lower bound introduced in DivideRounds; it is
		//redundant in normal operations, but becomes necessary after a Reset.
		//Indeed, after a Reset, LastConsensusRound is added to PendingRounds,
		//but its ConsensusEvents (which are necessarily 'under' this Round) are
		//already deemed committed. Hence, skip this Round after a Reset.
		if r.Index == p.GetLastConsensusRound() {
			continue
		}

		frame, err := p.GetFrame(r.Index)
		if err != nil {
			return fmt.Errorf("getting Frame %d: %v", r.Index, err)
		}

		round, err := p.Store.GetRound(r.Index)
		if err != nil {
			return err
		}
		p.logger.WithFields(logrus.Fields{
			"round_received": r.Index,
			"witnesses":      round.FamousWitnesses(),
			"events":         len(frame.Events),
			"roots":          frame.Roots,
		}).Debugf("Processing Decided Round")

		if len(frame.Events) > 0 {

			for _, e := range frame.Events {
				ev := e.ToEvent()
				err := p.Store.AddConsensusEvent(ev)
				if err != nil {
					return err
				}
				p.consensusTransactionsLocker.Lock()
				p.consensusTransactions += uint64(len(ev.Transactions()))
				p.consensusTransactionsLocker.Unlock()
				if ev.IsLoaded() {
					p.pendingLoadedEventsLocker.Lock()
					p.pendingLoadedEvents--
					p.pendingLoadedEventsLocker.Unlock()
				}
			}

			lastBlockIndex := p.Store.LastBlockIndex()
			block, err := NewBlockFromFrame(lastBlockIndex+1, frame)
			if err != nil {
				return err
			}
			if len(block.Transactions()) > 0 {
				if err := p.Store.SetBlock(block); err != nil {
					return err
				}

				if p.commitCh != nil {
					p.commitCh <- block
				}
			}

		} else {
			p.logger.Debugf("No Events to commit for ConsensusRound %d", r.Index)
		}

		processedIndex++

		if p.GetLastConsensusRound() < 0 || r.Index > p.GetLastConsensusRound() {
			p.setLastConsensusRound(r.Index)
		}

	}

	return nil
}

//GetFrame computes the Frame corresponding to a RoundReceived.
func (p *Poset) GetFrame(roundReceived int64) (Frame, error) {

	//Try to get it from the Store first
	frame, err := p.Store.GetFrame(roundReceived)
	if err == nil || !common.Is(err, common.KeyNotFound) {
		return frame, err
	}

	//Get the Round and corresponding consensus Events
	round, err := p.Store.GetRound(roundReceived)
	if err != nil {
		return Frame{}, err
	}

	var events []Event
	for _, eh := range round.ConsensusEvents() {
		e, err := p.Store.GetEvent(eh)
		if err != nil {
			return Frame{}, err
		}
		events = append(events, e)
	}

	sort.Stable(ByLamportTimestamp(events))

	// Get/Create Roots
	roots := make(map[string]Root)
	//The events are in topological order. Each time we run into the first Event
	//of a participant, we create a Root for it.
	for _, ev := range events {
		c := ev.Creator()
		if _, ok := roots[c]; !ok {
			root, err := p.createRoot(ev)
			if err != nil {
				return Frame{}, err
			}
			roots[ev.Creator()] = root
		}
	}

	//Every participant needs a Root in the Frame. For the participants that
	//have no Events in this Frame, we create a Root from their last consensus
	//Event, or their last known Root
	for _, peer := range p.Participants.ToPubKeySlice() {
		if _, ok := roots[peer]; !ok {
			var root Root
			lastConsensusEventHash, isRoot, err := p.Store.LastConsensusEventFrom(peer)
			if err != nil {
				return Frame{}, err
			}
			if isRoot {
				root, _ = p.Store.GetRoot(peer)
			} else {
				lastConsensusEvent, err := p.Store.GetEvent(lastConsensusEventHash)
				if err != nil {
					return Frame{}, err
				}
				root, err = p.createRoot(lastConsensusEvent)
				if err != nil {
					return Frame{}, err
				}
			}
			roots[peer] = root
		}
	}

	//Some Events in the Frame might have other-parents that are outside of the
	//Frame (cf root.go ex 2)
	//When inserting these Events in a newly reset poset, the CheckOtherParent
	//method would return an error because the other-parent would not be found.
	//So we make it possible to also look for other-parents in the creator's Root.
	treated := map[string]bool{}
	eventMessages := make([]*EventMessage, len(events))
	for i, ev := range events {
		treated[ev.Hex()] = true
		otherParent := ev.OtherParent()
		if otherParent != "" {
			opt, ok := treated[otherParent]
			if !opt || !ok {
				if ev.SelfParent() != roots[ev.Creator()].SelfParent.Hash {
					other, err := p.createOtherParentRootEvent(ev)
					if err != nil {
						return Frame{}, err
					}
					roots[ev.Creator()].Others[ev.Hex()] = &other
				}
			}
		}
		eventMessages[i] = new(EventMessage)
		*eventMessages[i] = ev.Message
	}

	//order roots
	orderedRoots := make([]*Root, p.Participants.Len())
	for i, peer := range p.Participants.ToPeerSlice() {
		root := roots[peer.PubKeyHex]
		orderedRoots[i] = new(Root)
		*orderedRoots[i] = root
	}

	res := Frame{
		Round:  roundReceived,
		Roots:  orderedRoots,
		Events: eventMessages,
	}

	if err := p.Store.SetFrame(res); err != nil {
		return Frame{}, err
	}

	return res, nil
}

//ProcessSigPool runs through the SignaturePool and tries to map a Signature to
//a known Block. If a Signature is found to be valid for a known Block, it is
//appended to the block and removed from the SignaturePool
func (p *Poset) ProcessSigPool() error {
	processedSignatures := map[int64]bool{} //index in SigPool => Processed?
	defer p.removeProcessedSignatures(processedSignatures)

	for i, bs := range p.SigPool {
		//check if validator belongs to list of participants
		validatorHex := fmt.Sprintf("0x%X", bs.Validator)
		if _, ok := p.Participants.GetByPubKey(validatorHex); !ok {
			p.logger.WithFields(logrus.Fields{
				"index":     bs.Index,
				"validator": validatorHex,
			}).Warning("Verifying Block signature. Unknown validator")
			continue
		}
		//only check if bs is greater than AnchorBlock, otherwise simply remove
		if p.AnchorBlock == nil ||
			bs.Index > *p.AnchorBlock {
			block, err := p.Store.GetBlock(bs.Index)
			if err != nil {
				p.logger.WithFields(logrus.Fields{
					"index": bs.Index,
					"msg":   err,
				}).Warning("Verifying Block signature. Could not fetch Block")
				continue
			}
			valid, err := block.Verify(bs)
			if err != nil {
				p.logger.WithFields(logrus.Fields{
					"index": bs.Index,
					"msg":   err,
				}).Error("Verifying Block signature")
				return err
			}
			if !valid {
				p.logger.WithFields(logrus.Fields{
					"index":     bs.Index,
					"validator": func (args ... interface{})interface{} {
						return args[0]
					}(p.Participants.GetByPubKey(validatorHex)),
					"block":     block,
				}).Warning("Verifying Block signature. Invalid signature")
				continue
			}

			block.SetSignature(bs)

			if err := p.Store.SetBlock(block); err != nil {
				p.logger.WithFields(logrus.Fields{
					"index": bs.Index,
					"msg":   err,
				}).Warning("Saving Block")
			}

			if len(block.Signatures) > p.trustCount &&
				(p.AnchorBlock == nil ||
					block.Index() > *p.AnchorBlock) {
				p.setAnchorBlock(block.Index())
				p.logger.WithFields(logrus.Fields{
					"block_index": block.Index(),
					"signatures":  len(block.Signatures),
					"trustCount":  p.trustCount,
				}).Debug("Setting AnchorBlock")
			}
		}

		processedSignatures[int64(i)] = true
	}

	return nil
}

//GetAnchorBlockWithFrame returns the AnchorBlock and the corresponding Frame.
//This can be used as a base to Reset a Poset
func (p *Poset) GetAnchorBlockWithFrame() (Block, Frame, error) {

	if p.AnchorBlock == nil {
		return Block{}, Frame{}, fmt.Errorf("no Anchor Block")
	}

	block, err := p.Store.GetBlock(*p.AnchorBlock)
	if err != nil {
		return Block{}, Frame{}, err
	}

	frame, err := p.GetFrame(block.RoundReceived())
	if err != nil {
		return Block{}, Frame{}, err
	}

	return block, frame, nil
}

//Reset clears the Poset and resets it from a new base.
func (p *Poset) Reset(block Block, frame Frame) error {

	//Clear all state
	p.firstLastConsensusRoundLocker.Lock()
	p.LastConsensusRound = nil
	p.FirstConsensusRound = nil
	p.firstLastConsensusRoundLocker.Unlock()
	p.AnchorBlock = nil

	p.undeterminedEventsLocker.Lock()
	p.UndeterminedEvents = []string{}
	p.undeterminedEventsLocker.Unlock()
	p.PendingRounds = []*pendingRound{}
	p.pendingLoadedEventsLocker.Lock()
	p.pendingLoadedEvents = 0
	p.pendingLoadedEventsLocker.Unlock()
	p.topologicalIndex = 0

	cacheSize := p.Store.CacheSize()
	ancestorCache, err := lru.New(cacheSize)
	if err != nil {
		p.logger.Fatal("Unable to reset Poset.ancestorCache")
	}
	selfAncestorCache, err := lru.New(cacheSize)
	if err != nil {
		p.logger.Fatal("Unable to reset Poset.selfAncestorCache")
	}
	stronglySeeCache, err := lru.New(cacheSize)
	if err != nil {
		p.logger.Fatal("Unable to reset Poset.stronglySeeCache")
	}
	roundCache, err := lru.New(cacheSize)
	if err != nil {
		p.logger.Fatal("Unable to reset Poset.roundCache")
	}
	p.ancestorCache = ancestorCache
	p.selfAncestorCache = selfAncestorCache
	p.stronglySeeCache = stronglySeeCache
	p.roundCache = roundCache

	participants := p.Participants.ToPeerSlice()

	//Initialize new Roots
	rootMap := map[string]Root{}
	for id, root := range frame.Roots {
		p := participants[id]
		rootMap[p.PubKeyHex] = *root
	}
	if err := p.Store.Reset(rootMap); err != nil {
		return err
	}

	//Insert Block
	if err := p.Store.SetBlock(block); err != nil {
		return err
	}

	p.setLastConsensusRound(block.RoundReceived())

	//Insert Frame Events
	for _, ev := range frame.Events {
		if err := p.InsertEvent(ev.ToEvent(), false); err != nil {
			return err
		}
	}

	return nil
}

//Bootstrap loads all Events from the Store's DB (if there is one) and feeds
//them to the Poset (in topological order) for consensus ordering. After this
//method call, the Poset should be in a state coherent with the 'tip' of the
//Poset
func (p *Poset) Bootstrap() error {
	if badgerStore, ok := p.Store.(*BadgerStore); ok {
		//Retreive the Events from the underlying DB. They come out in topological
		//order
		topologicalEvents, err := badgerStore.dbTopologicalEvents()
		if err != nil {
			return err
		}

		//Insert the Events in the Poset
		for _, e := range topologicalEvents {
			if err := p.InsertEvent(e, true); err != nil {
				return err
			}
		}

		//Compute the consensus order of Events
		if err := p.DivideRounds(); err != nil {
			return err
		}
		if err := p.DecideFame(); err != nil {
			return err
		}
		if err := p.DecideRoundReceived(); err != nil {
			return err
		}
		if err := p.ProcessDecidedRounds(); err != nil {
			return err
		}
		if err := p.ProcessSigPool(); err != nil {
			return err
		}
	}

	return nil
}

//ReadWireInfo converts a WireEvent to an Event by replacing int IDs with the
//corresponding public keys.
func (p *Poset) ReadWireInfo(wevent WireEvent) (*Event, error) {
	selfParent := rootSelfParent(wevent.Body.CreatorID)
	otherParent := ""
	var err error

	creator, ok := p.Participants.GetById(wevent.Body.CreatorID)
	// FIXIT: creator can be nil when wevent.Body.CreatorID == 0
	if !ok {
		return nil, fmt.Errorf("unknown wevent.Body.CreatorID=%v", wevent.Body.CreatorID)
	}
	creatorBytes, err := hex.DecodeString(creator.PubKeyHex[2:])
	if err != nil {
		return nil, err
	}

	if wevent.Body.SelfParentIndex >= 0 {
		selfParent, err = p.Store.ParticipantEvent(creator.PubKeyHex, wevent.Body.SelfParentIndex)
		if err != nil {
			return nil, err
		}
	}
	if wevent.Body.OtherParentIndex >= 0 {
		otherParentCreator, ok := p.Participants.GetById(wevent.Body.OtherParentCreatorID)
		if ok {
			otherParent, err = p.Store.ParticipantEvent(otherParentCreator.PubKeyHex, wevent.Body.OtherParentIndex)
			if err != nil {
				//PROBLEM Check if other parent can be found in the root
				//problem, we do not known the WireEvent's EventHash, and
				//we do not know the creators of the roots RootEvents
				root, err := p.Store.GetRoot(creator.PubKeyHex)
				if err != nil {
					return nil, err
				}
				//loop through others
				found := false
				for _, re := range root.Others {
					if re.CreatorID == wevent.Body.OtherParentCreatorID &&
						re.Index == wevent.Body.OtherParentIndex {
						otherParent = re.Hash
						found = true
						break
					}
				}

				if !found {
					return nil, fmt.Errorf("OtherParent not found")
				}
			}
		} else {
			// unknown participant
			// TODO: we should handle this nicely
			return nil, fmt.Errorf("unknown participant ID=%v", wevent.Body.OtherParentCreatorID)
		}
	}

	if len(wevent.FlagTable) == 0 {
		return nil, fmt.Errorf("flag table is null")
	}

	transactions := make([]*InternalTransaction, len(wevent.Body.InternalTransactions))
	for i, v := range wevent.Body.InternalTransactions {
		transactions[i] = new(InternalTransaction)
		*transactions[i] = v
	}
	signatureValues := wevent.BlockSignatures(creatorBytes)
	blockSignatures := make([]*BlockSignature, len(signatureValues))
	for i, v := range signatureValues {
		blockSignatures[i] = new(BlockSignature)
		*blockSignatures[i] = v
	}
	body := EventBody{
		Transactions:         wevent.Body.Transactions,
		InternalTransactions: transactions,
		Parents:              []string{selfParent, otherParent},
		Creator:              creatorBytes,
		Index:                wevent.Body.Index,
		BlockSignatures:      blockSignatures,
	}

	event := &Event{
		Message: EventMessage{
			Body:         &body,
			Signature:    wevent.Signature,
			FlagTable:    wevent.FlagTable,
			WitnessProof: wevent.WitnessProof,
			SelfParentIndex:      wevent.Body.SelfParentIndex,
			OtherParentCreatorID: wevent.Body.OtherParentCreatorID,
			OtherParentIndex:     wevent.Body.OtherParentIndex,
			CreatorID:            wevent.Body.CreatorID,
			LamportTimestamp:     LamportTimestampNIL,
			Round:                RoundNIL,
			RoundReceived:        RoundNIL,
		},
	}

	p.logger.WithFields(logrus.Fields{
		"event.Signature":  event.Message.Signature,
		"wevent.Signature": wevent.Signature,
	}).Debug("Return Event from ReadFromWire")

	return event, nil
}

//CheckBlock returns an error if the Block does not contain valid signatures
//from MORE than 1/3 of participants
func (p *Poset) CheckBlock(block Block) error {
	validSignatures := 0
	for _, s := range block.GetBlockSignatures() {
		ok, _ := block.Verify(s)
		if ok {
			validSignatures++
		}
	}
	if validSignatures <= p.trustCount {
		return fmt.Errorf("not enough valid signatures: got %d, need %d", validSignatures, p.trustCount+1)
	}

	p.logger.WithField("valid_signatures", validSignatures).Debug("CheckBlock")
	return nil
}

/*******************************************************************************
Setters
*******************************************************************************/

func (p *Poset) setLastConsensusRound(i int64) {
	p.firstLastConsensusRoundLocker.Lock()
	defer p.firstLastConsensusRoundLocker.Unlock()
	if p.LastConsensusRound == nil {
		p.LastConsensusRound = new(int64)
	}
	*p.LastConsensusRound = i

	if p.FirstConsensusRound == nil {
		p.FirstConsensusRound = new(int64)
		*p.FirstConsensusRound = i
	}
}

func (p *Poset) setAnchorBlock(i int64) {
	if p.AnchorBlock == nil {
		p.AnchorBlock = new(int64)
	}
	*p.AnchorBlock = i
}

/*******************************************************************************
Getters
*******************************************************************************/

func (p *Poset) GetFlagTableOfRandomUndeterminedEvent() (result map[string]int64, err error) {
	p.undeterminedEventsLocker.RLock()
	defer p.undeterminedEventsLocker.RUnlock()

	perm := rand.Perm(len(p.UndeterminedEvents))
	for i := 0; i < len(perm); i++ {
		hash := p.UndeterminedEvents[perm[i]]
		ev, err := p.Store.GetEvent(hash)
		if err != nil {
			continue
		}
		ft, err := ev.GetFlagTable()
		if err != nil {
			continue
		}
		if len(ft) >= p.Participants.Len() {
			continue
		}
		return ft, nil
	}
	return nil, err
}

func (p *Poset) GetUndeterminedEvents() []string {
	p.undeterminedEventsLocker.RLock()
	defer p.undeterminedEventsLocker.RUnlock()
	return p.UndeterminedEvents
}

func (p *Poset)  GetPendingLoadedEvents() int64 {
	p.pendingLoadedEventsLocker.RLock()
	defer p.pendingLoadedEventsLocker.RUnlock()
	return p.pendingLoadedEvents
}

func (p *Poset) GetLastConsensusRound() int64 {
	p.firstLastConsensusRoundLocker.RLock()
	defer p.firstLastConsensusRoundLocker.RUnlock()
	if p.LastConsensusRound == nil {
		// -2 is less that undefined round index, -1
		return -2
	}
	return *p.LastConsensusRound
}

func (p *Poset) GetConsensusTransactionsCount() uint64 {
	p.consensusTransactionsLocker.RLock()
	defer p.consensusTransactionsLocker.RUnlock()
	return p.consensusTransactions
}

/*******************************************************************************
   Helpers
*******************************************************************************/

func middleBit(ehex string) bool {
	hash, err := hex.DecodeString(ehex[2:])
	if err != nil {
		fmt.Printf("ERROR decoding hex string: %s\n", err)
	}
	if len(hash) > 0 && hash[len(hash)/2] == 0 {
		return false
	}
	return true
}
