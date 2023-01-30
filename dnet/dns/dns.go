// Copyright (c) 2018-2019 Aalborg University
// Use of this source code is governed by a GPLv3
// license that can be found in the LICENSE file.

package dns

import (
	"context"
	"github.com/aau-network-security/sandbox/virtual"

	//"github.com/aau-network-security/defatt/sandbox"
	//"github.com/aau-network-security/sandbox2/store"
	//"github.com/aau-network-security/defatt/virtual"

	"github.com/aau-network-security/sandbox/virtual/docker"
	"github.com/rs/zerolog/log"
)

type Server struct {
	cont     docker.Container
	IPanswer string
	IPbind   string
	bindPort string
	//corefile string
	//zonefile string
	//ipList   map[string]string
}

func New(ctx context.Context, bridge string) (*Server, error) {

	container := docker.NewContainer(docker.ContainerConfig{
		Image: "docker.io/rasmim/fakedns",
		Labels: map[string]string{
			"nap-sandbox": bridge,
		},
		Resources: &docker.Resources{
			MemoryMB: 100,
			CPU:      0.5,
		},
		Cmd: []string{"/fakedns/fakedns.py", "-a", "10.10.10.55"},
	})

	return &Server{
		cont: container,
		//corefile: Corefile,
		//zonefile: zonefile,
	}, nil
}

func (s *Server) Container() docker.Container {
	return s.cont
}

func (s *Server) Run(ctx context.Context) error {
	return s.cont.Run(ctx)
}

//func (s *Server) Close(ctx context.Context) error {
//	if err := os.Remove(s.corefile); err != nil {
//		log.Warn().Msgf("error while removing DNS configuration file: %s", err)
//	}
//
//	if err := s.cont.Close(); err != nil {
//		log.Warn().Msgf("error while closing DNS container: %s", err)
//	}
//
//	return nil
//}

func (s *Server) Close() error {
	//if err := os.Remove(s.cont); err != nil {
	//	log.Warn().Msgf("error while removing DNS configuration file: %s", err)
	//}

	if err := s.cont.Close(); err != nil {
		log.Warn().Msgf("error while closing DNS container: %s", err)
	}

	return nil
}

func (s *Server) Create(ctx context.Context) error {
	if err := s.cont.Create(ctx); err != nil {
		log.Error().Err(err).Msg("creating container DNS")
		return err
	}
	return nil
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.cont.Start(ctx); err != nil {
		log.Error().Err(err).Msg("starting container")
		return err
	}
	return nil
}

//func (s *Server) Execute(ctx context.Context, i []string, s2 string) error {
//	panic("implement me")
//}
//
//func (s *Server) Suspend(ctx context.Context) error {
//	panic("implement me")
//}
//
//func (s *Server) Info() virtual.InstanceInfo {
//	panic("implement me")
//}

func (s *Server) ID() string {
	return s.cont.ID()
}

func (s *Server) Stop() error {
	return s.cont.Stop()
}

func (s *Server) Execute(ctx context.Context, strings []string, s2 string) error {
	panic("implement me")
}

func (s *Server) Suspend(ctx context.Context) error {
	return s.cont.Suspend(ctx)
}

func (s *Server) Info() virtual.InstanceInfo {
	panic("implement me")
}
