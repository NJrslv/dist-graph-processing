package svc

type Serviceable interface {
	serve()
}

type Storage struct {
}

func (s *Storage) serve() {

}
