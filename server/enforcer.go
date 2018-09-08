// Copyright 2018 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"errors"

	"github.com/casbin/casbin"
	pb "github.com/casbin/casbin-server/proto"
	"github.com/casbin/casbin/persist"
	"golang.org/x/net/context"
)

// Server is used to implement proto.CasbinServer.
type Server struct {
	enforcerMap map[int]*casbin.Enforcer
	adapterMap  map[int]persist.Adapter
}

func NewServer() *Server {
	s := Server{}

	s.enforcerMap = map[int]*casbin.Enforcer{}
	s.adapterMap = map[int]persist.Adapter{}

	return &s
}

func (s *Server) getEnforcer(handle int) (*casbin.Enforcer, error) {
	if _, ok := s.enforcerMap[handle]; ok {
		return s.enforcerMap[handle], nil
	} else {
		return nil, errors.New("enforcer not found")
	}
}

func (s *Server) getAdapter(handle int) (persist.Adapter, error) {
	if _, ok := s.adapterMap[handle]; ok {
		return s.adapterMap[handle], nil
	} else {
		return nil, errors.New("adapter not found")
	}
}

func (s *Server) addEnforcer(e *casbin.Enforcer) int {
	cnt := len(s.enforcerMap)
	s.enforcerMap[cnt] = e
	return cnt
}

func (s *Server) addAdapter(a persist.Adapter) int {
	cnt := len(s.adapterMap)
	s.adapterMap[cnt] = a
	return cnt
}

func (s *Server) NewEnforcer(ctx context.Context, in *pb.NewEnforcerRequest) (*pb.NewEnforcerReply, error) {
	a, err := s.getAdapter(int(in.AdapterHandle))
	if err != nil {
		return &pb.NewEnforcerReply{Handler: 0}, err
	}

	e := casbin.NewEnforcer(casbin.NewModel(in.ModelText), a)
	h := s.addEnforcer(e)

	return &pb.NewEnforcerReply{Handler: int32(h)}, nil
}

func (s *Server) NewAdapter(ctx context.Context, in *pb.NewAdapterRequest) (*pb.NewAdapterReply, error) {
	a, err := newAdapter(in)
	if err != nil {
		return nil, err
	}

	h := s.addAdapter(a)

	return &pb.NewAdapterReply{Handler: int32(h)}, nil
}

func (s *Server) Enforce(ctx context.Context, in *pb.EnforceRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{Res: false}, err
	}

	res := e.Enforce(in.Sub, in.Obj, in.Act)

	return &pb.BoolReply{Res: res}, nil
}

func (s *Server) LoadPolicy(ctx context.Context, in *pb.EmptyRequest) (*pb.EmptyReply, error) {
	e, err := s.getEnforcer(int(in.Handler))
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	err = e.LoadPolicy()

	return &pb.EmptyReply{}, err
}

func (s *Server) SavePolicy(ctx context.Context, in *pb.EmptyRequest) (*pb.EmptyReply, error) {
	e, err := s.getEnforcer(int(in.Handler))
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	err = e.SavePolicy()

	return &pb.EmptyReply{}, err
}

func (s *Server) AddPolicy(ctx context.Context, in *pb.PolicyRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	res := e.AddPolicy(in.Sub, in.Obj, in.Act)

	return &pb.BoolReply{Res: res}, err
}

func (s *Server) RemovePolicy(ctx context.Context, in *pb.PolicyRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	res := e.RemovePolicy(in.Sub, in.Obj, in.Act)

	return &pb.BoolReply{Res: res}, err
}

func (s *Server) AddNamedPolicy(ctx context.Context, in *pb.NamedPolicyRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	res := e.AddNamedPolicy(in.Ptype, in.Sub, in.Obj, in.Act)

	return &pb.BoolReply{Res: res}, err
}

func (s *Server) RemoveNamedPolicy(ctx context.Context, in *pb.NamedPolicyRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	res := e.RemoveNamedPolicy(in.Ptype, in.Sub, in.Obj, in.Act)

	return &pb.BoolReply{Res: res}, err
}

// GetAllSubjects gets the list of subjects that show up in the current policy.
func (s *Server) GetAllSubjects(ctx context.Context, in *pb.EmptyRequest) (*pb.ArrayReply, error) {
	return s.GetAllNamedSubjects(ctx, &pb.SimpleGetRequest{EnforcerHandler: in.Handler, PType: "p"})
}

// GetAllNamedSubjects gets the list of subjects that show up in the current named policy.
func (s *Server) GetAllNamedSubjects(ctx context.Context, in *pb.SimpleGetRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}

	return &pb.ArrayReply{Array: e.GetModel().GetValuesForFieldInPolicy("p", in.PType, 0)}, nil
}

// GetAllObjects gets the list of objects that show up in the current policy.
func (s *Server) GetAllObjects(ctx context.Context, in *pb.EmptyRequest) (*pb.ArrayReply, error) {
	return s.GetAllNamedObjects(ctx, &pb.SimpleGetRequest{EnforcerHandler: in.Handler, PType: "p"})
}

// GetAllNamedObjects gets the list of objects that show up in the current named policy.
func (s *Server) GetAllNamedObjects(ctx context.Context, in *pb.SimpleGetRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}

	return &pb.ArrayReply{Array: e.GetModel().GetValuesForFieldInPolicy("p", in.PType, 1)}, nil
}

// GetAllActions gets the list of actions that show up in the current policy.
func (s *Server) GetAllActions(ctx context.Context, in *pb.EmptyRequest) (*pb.ArrayReply, error) {
	return s.GetAllNamedActions(ctx, &pb.SimpleGetRequest{EnforcerHandler: in.Handler, PType: "p"})
}

// GetAllNamedActions gets the list of actions that show up in the current named policy.
func (s *Server) GetAllNamedActions(ctx context.Context, in *pb.SimpleGetRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}

	return &pb.ArrayReply{Array: e.GetModel().GetValuesForFieldInPolicy("p", in.PType, 2)}, nil
}

// GetAllRoles gets the list of roles that show up in the current policy.
func (s *Server) GetAllRoles(ctx context.Context, in *pb.EmptyRequest) (*pb.ArrayReply, error) {
	return s.GetAllNamedRoles(ctx, &pb.SimpleGetRequest{EnforcerHandler: in.Handler, PType: "g"})
}

// GetAllNamedRoles gets the list of roles that show up in the current named policy.
func (s *Server) GetAllNamedRoles(ctx context.Context, in *pb.SimpleGetRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}

	return &pb.ArrayReply{Array: e.GetModel().GetValuesForFieldInPolicy("g", in.PType, 1)}, nil
}

// GetPolicy gets all the authorization rules in the policy.
func (s *Server) GetPolicy(ctx context.Context, in *pb.EmptyRequest) (*pb.Array2DReply, error) {
	return s.GetNamedPolicy(ctx, &pb.GetPolicyRequest{EnforcerHandler: in.Handler, PType: "p"})
}

// GetNamedPolicy gets all the authorization rules in the named policy.
func (s *Server) GetNamedPolicy(ctx context.Context, in *pb.GetPolicyRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetModel().GetPolicy("p", in.PType)), nil
}

// GetFilteredPolicy gets all the authorization rules in the policy, field filters can be specified.
func (s *Server) GetFilteredPolicy(ctx context.Context, in *pb.GetPolicyRequest) (*pb.Array2DReply, error) {
	in.PType = "p"

	return s.GetFilteredNamedPolicy(ctx, in)
}

// GetFilteredNamedPolicy gets all the authorization rules in the named policy, field filters can be specified.
func (s *Server) GetFilteredNamedPolicy(ctx context.Context, in *pb.GetPolicyRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetModel().GetFilteredPolicy("p", in.PType, int(in.FieldIndex), in.FieldValues...)), nil
}

// GetGroupingPolicy gets all the role inheritance rules in the policy.
func (s *Server) GetGroupingPolicy(ctx context.Context, in *pb.EmptyRequest) (*pb.Array2DReply, error) {
	return s.GetNamedGroupingPolicy(ctx, &pb.GetPolicyRequest{EnforcerHandler: in.Handler, PType: "g"})
}

// GetNamedGroupingPolicy gets all the role inheritance rules in the policy.
func (s *Server) GetNamedGroupingPolicy(ctx context.Context, in *pb.GetPolicyRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetModel().GetPolicy("g", in.PType)), nil
}

// GetFilteredGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.
func (s *Server) GetFilteredGroupingPolicy(ctx context.Context, in *pb.GetPolicyRequest) (*pb.Array2DReply, error) {
	in.PType = "g"

	return s.GetFilteredNamedGroupingPolicy(ctx, in)
}

// GetFilteredNamedGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.
func (s *Server) GetFilteredNamedGroupingPolicy(ctx context.Context, in *pb.GetPolicyRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetModel().GetFilteredPolicy("g", in.PType, int(in.FieldIndex), in.FieldValues...)), nil
}

// HasPolicy determines whether an authorization rule exists.
func (s *Server) HasPolicy(ctx context.Context, in *pb.ExistRequest) (*pb.BoolReply, error) {
	return s.HasNamedPolicy(ctx, in)
}

// HasNamedPolicy determines whether a named authorization rule exists.
func (s *Server) HasNamedPolicy(ctx context.Context, in *pb.ExistRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.GetModel().HasPolicy("p", in.Ptype, in.Policy)}, nil
}

// HasGroupingPolicy determines whether a role inheritance rule exists.
func (s *Server) HasGroupingPolicy(ctx context.Context, in *pb.ExistRequest) (*pb.BoolReply, error) {
	in.Ptype = "g"

	return s.HasNamedGroupingPolicy(ctx, in)
}

// HasNamedGroupingPolicy determines whether a named role inheritance rule exists.
func (s *Server) HasNamedGroupingPolicy(ctx context.Context, in *pb.ExistRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.GetModel().HasPolicy("g", in.Ptype, in.Policy)}, nil
}

func (s *Server) wrapPlainPolicy(policy [][]string) *pb.Array2DReply {
	if len(policy) == 0 {
		return &pb.Array2DReply{}
	}

	policyReply := &pb.Array2DReply{}
	policyReply.D2 = make([]*pb.Array2DReplyD, len(policy))
	for e := range policy {
		policyReply.D2[e] = &pb.Array2DReplyD{D1: policy[e]}
	}

	return policyReply
}
