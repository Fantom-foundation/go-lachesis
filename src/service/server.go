package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Fantom-foundation/go-lachesis/src/peers"
	"github.com/Fantom-foundation/go-lachesis/src/poset"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Service is a service required for handlers.
type Service interface {
	Stats() map[string]string
	Participants() (*peers.Peers, error)
	EventBlock(event poset.EventHash) (poset.Event, error)
	LastEventFrom(participant string) (poset.EventHash, bool, error)
	KnownEvents() map[uint64]int64
	ConsensusEvents() poset.EventHashes
	Round(roundIndex int64) (poset.RoundCreated, error)
	LastRound() int64
	RoundClothos(roundIndex int64) poset.EventHashes
	RoundEvents(roundIndex int64) int
	Root(rootIndex string) (poset.Root, error)
	Block(blockIndex int64) (poset.Block, error)
}

// Server http API service struct
type Server struct {
	bindAddress string
	service     Service
	logger      *logrus.Logger
}

// NewServer creates a new http API service
func NewServer(bindAddress string, service Service, logger *logrus.Logger) *Server {
	s := Server{
		bindAddress: bindAddress,
		service:     service,
		logger:      logger,
	}

	return &s
}

// Serve serves the API
func (s *Server) Serve() error {
	s.logger.WithField("bind_address", s.bindAddress).Debug("Server serving")
	mux := http.NewServeMux()

	mux.Handle("/stats", corsHandler(s.Stats))
	mux.Handle("/participants/", corsHandler(s.Participants))
	mux.Handle("/event/", corsHandler(s.EventBlock))
	mux.Handle("/lasteventfrom/", corsHandler(s.LastEventFrom))
	mux.Handle("/events/", corsHandler(s.KnownEvents))
	mux.Handle("/consensusevents/", corsHandler(s.ConsensusEvents))
	mux.Handle("/round/", corsHandler(s.Round))
	mux.Handle("/lastround/", corsHandler(s.LastRound))
	mux.Handle("/roundclothos/", corsHandler(s.RoundClothos))
	mux.Handle("/roundevents/", corsHandler(s.RoundEvents))
	mux.Handle("/root/", corsHandler(s.Root))
	mux.Handle("/block/", corsHandler(s.Block))

	if err := http.ListenAndServe(s.bindAddress, mux); err != nil {
		return errors.Wrap(err, "listen and serve")
	}

	return nil
}

func corsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		if r.Method == "OPTIONS" {
			/*w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")*/
		} else {
			/*w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")*/
			h.ServeHTTP(w, r)
		}
	}
}

// Stats returns all the node processing stats
func (s *Server) Stats(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("Stats")

	stats := s.service.Stats()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		s.logger.Debug(err)
	}
}

// Participants returns all the known participants
func (s *Server) Participants(w http.ResponseWriter, r *http.Request) {
	participants, err := s.service.Participants()
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing participants parameter")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(participants); err != nil {
		s.logger.Debug(err)
	}
}

// EventBlock returns a specific event block by id
func (s *Server) EventBlock(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/event/"):]

	var hash poset.EventHash
	err := hash.Parse(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing event hash %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event, err := s.service.EventBlock(hash)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving event %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		s.logger.Debug(err)
	}
}

// LastEventFrom returns the last event for a specific participant
func (s *Server) LastEventFrom(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/lasteventfrom/"):]
	event, _, err := s.service.LastEventFrom(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving event %s", event)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(event); err != nil {
		s.logger.Debug(err)
	}
}

// KnownEvents returns all known events by ID
func (s *Server) KnownEvents(w http.ResponseWriter, r *http.Request) {
	knownEvents := s.service.KnownEvents()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(knownEvents); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode known events: %v", knownEvents)
	}
}

// ConsensusEvents returns all the events that have reached consensus
func (s *Server) ConsensusEvents(w http.ResponseWriter, r *http.Request) {
	consensusEvents := s.service.ConsensusEvents()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(consensusEvents); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode consensus events: %v", consensusEvents)
	}
}

// Round returns a round for the given index
func (s *Server) Round(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/round/"):]
	roundIndex, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing roundIndex parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	round, err := s.service.Round(roundIndex)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving round %d", roundIndex)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(round); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode round: %v", round)
	}
}

// LastRound returns the last known round
func (s *Server) LastRound(w http.ResponseWriter, r *http.Request) {
	lastRound := s.service.LastRound()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lastRound); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode last round: %d", lastRound)
	}
}

// RoundClothos returns all clotho for a round
func (s *Server) RoundClothos(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/roundclothos/"):]
	roundClothosIndex, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing roundClothosIndex parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roundClothos := s.service.RoundClothos(roundClothosIndex)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(roundClothos); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode round clothos: %v", roundClothos)
	}
}

// RoundEvents returns all the events for a given round
func (s *Server) RoundEvents(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/roundevents/"):]
	roundEventsIndex, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing roundEventsIndex parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roundEvent := s.service.RoundEvents(roundEventsIndex)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(roundEvent); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode round event: %d", roundEvent)
	}
}

// Root returns the root for a given frame
func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/root/"):]
	root, err := s.service.Root(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving root %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(root); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode root: %v", root)
	}
}

// Block returns a specific block based on index
func (s *Server) Block(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/block/"):]
	blockIndex, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing block_index parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	block, err := s.service.Block(blockIndex)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving block %d", blockIndex)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(block); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode block: %v", block)
	}
}
