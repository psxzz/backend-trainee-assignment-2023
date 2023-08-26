package service

type Storage interface {
}

type Service struct {
	storage Storage
}

func New(storage Storage) (*Service, error) {
	return &Service{
		storage: storage,
	}, nil
}

func (svc *Service) CreateSegment() {
}

func (svc *Service) DeleteSegment() {

}

func (svc *Service) AddUser() {

}

func (svc *Service) ListUser() {

}
