package service

import (
	"context"
	"fmt"
	"housework/v1"
	"sync"
)

type HouseworkServer struct {
	mu     sync.Mutex
	chores []*housework.Chore
	housework.UnimplementedHouseworkServer
}

func (h *HouseworkServer) Add(_ context.Context, input *housework.AddInput) (*housework.AddPayload, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, c := range input.Chores {
		h.chores = append(h.chores, &housework.Chore{
			Complete:    c.Complete,
			Description: c.Description,
		})
	}
	return &housework.AddPayload{Message: "ok"}, nil
}

func (h *HouseworkServer) Get(_ context.Context, req *housework.GetChoreInput) (
	*housework.Chore, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, c := range h.chores {
		if c.Id == req.Id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("Chore %d not found", req.Id)
}

func (h *HouseworkServer) List(_ context.Context, _ *housework.Empty) (
	*housework.Chores, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.chores == nil {
		h.chores = make([]*housework.Chore, 0)
	}
	return &housework.Chores{
		Chores: h.chores,
	}, nil
}

func (h *HouseworkServer) Complete(_ context.Context, req *housework.CompleteInput) (
	*housework.CompletePayload, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.chores == nil || req.ChoreNumber < 1 || int(req.ChoreNumber) > len(h.chores) {
		return nil, fmt.Errorf("chore %d not found", req.ChoreNumber)
	}
	h.chores[req.ChoreNumber-1].Complete = true

	return &housework.CompletePayload{Message: "ok"}, nil
}

func Service() housework.HouseworkServer {
	return &HouseworkServer{
		chores: make([]*housework.Chore, 0),
	}
}
