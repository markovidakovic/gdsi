package players

import "context"

type service struct {
	store *store
}

func (s *service) getAllPlayers(ctx context.Context) (string, error) {
	result, err := s.store.queryPlayers(ctx)
	if err != nil {
		return "", err
	}
	return result, nil
}

func newService(store *store) *service {
	var s = &service{
		store,
	}
	return s
}
