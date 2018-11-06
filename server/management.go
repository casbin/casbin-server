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
	"context"
	pb "github.com/casbin/casbin-server/proto"
)

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
	return s.GetNamedPolicy(ctx, &pb.PolicyRequest{EnforcerHandler: in.Handler, PType: "p"})
}

// GetNamedPolicy gets all the authorization rules in the named policy.
func (s *Server) GetNamedPolicy(ctx context.Context, in *pb.PolicyRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetModel().GetPolicy("p", in.PType)), nil
}

// GetFilteredPolicy gets all the authorization rules in the policy, field filters can be specified.
func (s *Server) GetFilteredPolicy(ctx context.Context, in *pb.FilteredPolicyRequest) (*pb.Array2DReply, error) {
	in.PType = "p"

	return s.GetFilteredNamedPolicy(ctx, in)
}

// GetFilteredNamedPolicy gets all the authorization rules in the named policy, field filters can be specified.
func (s *Server) GetFilteredNamedPolicy(ctx context.Context, in *pb.FilteredPolicyRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetModel().GetFilteredPolicy("p", in.PType, int(in.FieldIndex), in.FieldValues...)), nil
}

// GetGroupingPolicy gets all the role inheritance rules in the policy.
func (s *Server) GetGroupingPolicy(ctx context.Context, in *pb.EmptyRequest) (*pb.Array2DReply, error) {
	return s.GetNamedGroupingPolicy(ctx, &pb.PolicyRequest{EnforcerHandler: in.Handler, PType: "g"})
}

// GetNamedGroupingPolicy gets all the role inheritance rules in the policy.
func (s *Server) GetNamedGroupingPolicy(ctx context.Context, in *pb.PolicyRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetModel().GetPolicy("g", in.PType)), nil
}

// GetFilteredGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.
func (s *Server) GetFilteredGroupingPolicy(ctx context.Context, in *pb.FilteredPolicyRequest) (*pb.Array2DReply, error) {
	in.PType = "g"

	return s.GetFilteredNamedGroupingPolicy(ctx, in)
}

// GetFilteredNamedGroupingPolicy gets all the role inheritance rules in the policy, field filters can be specified.
func (s *Server) GetFilteredNamedGroupingPolicy(ctx context.Context, in *pb.FilteredPolicyRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetModel().GetFilteredPolicy("g", in.PType, int(in.FieldIndex), in.FieldValues...)), nil
}

// HasPolicy determines whether an authorization rule exists.
func (s *Server) HasPolicy(ctx context.Context, in *pb.PolicyRequest) (*pb.BoolReply, error) {
	return s.HasNamedPolicy(ctx, in)
}

// HasNamedPolicy determines whether a named authorization rule exists.
func (s *Server) HasNamedPolicy(ctx context.Context, in *pb.PolicyRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.GetModel().HasPolicy("p", in.PType, in.Params)}, nil
}

// HasGroupingPolicy determines whether a role inheritance rule exists.
func (s *Server) HasGroupingPolicy(ctx context.Context, in *pb.PolicyRequest) (*pb.BoolReply, error) {
	in.PType = "g"

	return s.HasNamedGroupingPolicy(ctx, in)
}

// HasNamedGroupingPolicy determines whether a named role inheritance rule exists.
func (s *Server) HasNamedGroupingPolicy(ctx context.Context, in *pb.PolicyRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.GetModel().HasPolicy("g", in.PType, in.Params)}, nil
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
