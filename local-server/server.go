package local_server

type Server struct {
	HttpAddr string
}

func (s *Server) Start() {
	if s.HttpAddr != "" {
		startHttp(s.HttpAddr)
	}
}

func (s *Server) Stop() error {
	if s.HttpAddr != "" {
		err := stopHttp()
		return err
	}
	return nil
}

func New(httpAddr string) Server {
	return Server{httpAddr}
}
