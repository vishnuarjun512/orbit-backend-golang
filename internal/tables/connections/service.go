package connections

import "context"

type ConnectionService struct {
	repo *ConnectionRepository
}

func NewConnectionService(repo *ConnectionRepository) *ConnectionService {
	return &ConnectionService{repo: repo}
}

func (s *ConnectionService) CreateNewConnectionRequestService(ctx context.Context, senderID, receiverID string) (*Connection_Request, error) {

	connection_req_id, err := s.repo.CreateNewConnectionRequest(ctx, senderID, receiverID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetConnectionRequest(ctx, connection_req_id)
}
