package connections

import (
	"context"
	"errors"
)

type ConnectionService struct {
	repo *ConnectionRepository
}

func NewConnectionService(repo *ConnectionRepository) *ConnectionService {
	return &ConnectionService{repo: repo}
}

var (
	ErrConnectionRequestNotFound = errors.New("Connection Request not found")
	ErrConnectionNotFound        = errors.New("Connection not found")
)

func (s *ConnectionService) CreateNewConnectionRequestService(ctx context.Context, senderID, receiverID string) (*Connection_Request, error) {

	connection_req_id, err := s.repo.CreateNewConnectionRequest(ctx, senderID, receiverID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetConnectionRequestByID(ctx, connection_req_id)
}

func (s *ConnectionService) GetConnectionRequestService(ctx context.Context, senderID, receiverID string) (*Connection_Request, error) {
	return s.repo.GetConnectionRequest(ctx, senderID, receiverID)
}

func (s *ConnectionService) GetConnectionRequest(ctx context.Context, senderID, receiverID string) (*Connection_Request, error) {
	return s.repo.GetConnectionRequest(ctx, senderID, receiverID)
}

func (s *ConnectionService) GetReceivedRequestsService(ctx context.Context, receiverID string) (*[]Connection_Request, error) {
	return s.repo.GetReceivedRequests(ctx, receiverID)
}

func (s *ConnectionService) GetSentRequestsService(ctx context.Context, senderID string) (*[]Connection_Request, error) {
	return s.repo.GetSentRequests(ctx, senderID)
}

func (s *ConnectionService) DeleteConnectionRequestService(ctx context.Context, senderID, receiverID string) error {
	return s.repo.DeleteConnectionRequest(ctx, senderID, receiverID)
}

func (s *ConnectionService) DeleteConnectionService(ctx context.Context, senderID, receiverID string) error {
	return s.repo.DeleteConnection(ctx, senderID, receiverID)
}

func (s *ConnectionService) AcceptConnectionRequestService(ctx context.Context, senderID, receiverID string) (*Connection, error) {

	connectionrequest, err := s.repo.GetConnectionRequest(ctx, senderID, receiverID)

	if err != nil {
		return nil, err
	}

	connection_id, cerr := s.repo.CreateNewConnection(ctx, senderID, receiverID)
	if cerr != nil {
		return nil, cerr
	}

	connection, cerr := s.repo.GetConnectionByID(ctx, connection_id)
	if cerr != nil {
		return nil, cerr
	}

	crerr := s.repo.DeleteConnectionRequestByID(ctx, connectionrequest.ConnectionRequestID)
	if crerr != nil {
		return nil, crerr
	}

	return connection, nil
}

func (s *ConnectionService) RejectConnectionRequestService(ctx context.Context, senderID, receiverID string) error {

	connectionrequest, err := s.repo.GetConnectionRequest(ctx, senderID, receiverID)

	if err != nil {
		return err
	}

	crerr := s.repo.DeleteConnectionRequestByID(ctx, connectionrequest.ConnectionRequestID)
	if crerr != nil {
		return crerr
	}

	return nil
}

func (s *ConnectionService) GetAllConnectionsService(ctx context.Context, userID string) (*[]ConnectionResponse, error) {
	return s.repo.GetAllConnections(ctx, userID)
}
