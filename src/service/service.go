package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Fantom-foundation/go-lachesis/src/node"
	"github.com/sirupsen/logrus"
)

// Service http API service struct
type Service struct {
	bindAddress string
	node        *node.Node
	graph       *node.Graph
	logger      *logrus.Logger
}

// NewService creates a new http API service
func NewService(bindAddress string, n *node.Node, logger *logrus.Logger) *Service {
	service := Service{
		bindAddress: bindAddress,
		node:        n,
		graph:       node.NewGraph(n),
		logger:      logger,
	}

	return &service
}

// Serve serves the API
func (s *Service) Serve() {
	s.logger.WithField("bind_address", s.bindAddress).Debug("Service serving")
	mux := http.NewServeMux()
	mux.Handle("/stats", corsHandler(s.GetStats))
	mux.Handle("/participants/", corsHandler(s.GetParticipants))
	mux.Handle("/event/", corsHandler(s.GetEventBlock))
	mux.Handle("/lasteventfrom/", corsHandler(s.GetLastEventFrom))
	mux.Handle("/events/", corsHandler(s.GetKnownEvents))
	mux.Handle("/consensusevents/", corsHandler(s.GetConsensusEvents))
	mux.Handle("/round/", corsHandler(s.GetRound))
	mux.Handle("/lastround/", corsHandler(s.GetLastRound))
	mux.Handle("/roundclothos/", corsHandler(s.GetRoundClothos))
	mux.Handle("/roundevents/", corsHandler(s.GetRoundEvents))
	mux.Handle("/root/", corsHandler(s.GetRoot))
	mux.Handle("/block/", corsHandler(s.GetBlock))
	err := http.ListenAndServe(s.bindAddress, mux)
	if err != nil {
		s.logger.WithField("error", err).Error("Service failed")
	}
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

// GetStats returns all the node processing stats
func (s *Service) GetStats(w http.ResponseWriter, r *http.Request) {
	s.logger.Debug("Stats")

	stats := s.node.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// GetParticipants returns all the known participants
func (s *Service) GetParticipants(w http.ResponseWriter, r *http.Request) {
	participants, err := s.node.GetParticipants()
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing participants parameter")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participants)
}

// GetEventBlock returns a specific event block by id
func (s *Service) GetEventBlock(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/event/"):]
	event, err := s.node.GetEventBlock(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving event %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// GetLastEventFrom returns the last event for a specific participant
func (s *Service) GetLastEventFrom(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/lasteventfrom/"):]
	event, _, err := s.node.GetLastEventFrom(param)
	if err != nil {
		s.logger.WithError(err).Errorf("Retrieving event %s", event)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// GetKnownEvents returns all known events by ID
func (s *Service) GetKnownEvents(w http.ResponseWriter, r *http.Request) {
	knownEvents := s.node.GetKnownEvents()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(knownEvents); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode known events: %v", knownEvents)
	}
}

// GetConsensusEvents returns all the events that have reached consensus
func (s *Service) GetConsensusEvents(w http.ResponseWriter, r *http.Request) {
	consensusEvents := s.node.GetConsensusEvents()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(consensusEvents); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode consensus events: %v", consensusEvents)
	}
}

// GetRound returns a round for the given index
func (s *Service) GetRound(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/round/"):]
	roundIndex, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing roundIndex parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	round, err := s.node.GetRound(roundIndex)
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

// GetLastRound returns the last known round
func (s *Service) GetLastRound(w http.ResponseWriter, r *http.Request) {
	lastRound := s.node.GetLastRound()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(lastRound); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode last round: %d", lastRound)
	}
}

// GetRoundClothos returns all clotho for a round
func (s *Service) GetRoundClothos(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/roundclothos/"):]
	roundClothosIndex, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing roundClothosIndex parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roundClothos := s.node.GetRoundClothos(roundClothosIndex)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(roundClothos); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode round clothos: %v", roundClothos)
	}
}

// GetRoundEvents returns all the events for a given round
func (s *Service) GetRoundEvents(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/roundevents/"):]
	roundEventsIndex, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing roundEventsIndex parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	roundEvent := s.node.GetRoundEvents(roundEventsIndex)

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(roundEvent); err != nil {
		s.logger.WithError(err).Errorf("Failed to encode round event: %d", roundEvent)
	}
}

// GetRoot returns the root for a given frame
func (s *Service) GetRoot(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/root/"):]
	root, err := s.node.GetRoot(param)
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

// GetBlock returns a specific block based on index
func (s *Service) GetBlock(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Path[len("/block/"):]
	blockIndex, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		s.logger.WithError(err).Errorf("Parsing block_index parameter %s", param)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	block, err := s.node.GetBlock(blockIndex)
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
