package endpoint

type Service interface {
	CreateSegment()
	DeleteSegment()
	AddUser()
	ListUser()
}

type Endpoint struct {
	svc Service
}

func New(svc Service) *Endpoint {
	return &Endpoint{
		svc: svc,
	}
}

// func (e *Endpoint) CreateSegment(ctx echo.Context) error {

// }

// func (e *Endpoint) DeleteSegment(ctx echo.Context) error {

// }

// func (e *Endpoint) AddSegment(ctx echo.Context) error {

// }

// func (e *Endpoint) UserSegments(ctx echo.Context) error {

// }
