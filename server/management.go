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
func (s *Server) GetAllSubjects(ctx context.Context, in *pb.GetSujectsRequest) (*pb.SubjectsReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.SubjectsReply{}, err
	}

	return &pb.SubjectsReply{Subjects: e.GetAllSubjects()}, nil
}

// GetAllNamedSubjects gets the list of subjects that show up in the current named policy.
func (s *Server) GetAllNamedSubjects(ctx context.Context, in *pb.GetSujectsRequest) (*pb.SubjectsReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.SubjectsReply{}, err
	}

	return &pb.SubjectsReply{Subjects: e.GetModel().GetValuesForFieldInPolicy("p", in.PType, 0)}, nil
}

// GetAllObjects gets the list of objects that show up in the current policy.
func (s *Server) GetAllObjects(ctx context.Context, in *pb.GetObjectsRequest) (*pb.ObjectsReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ObjectsReply{}, err
	}

	return &pb.ObjectsReply{Objects: e.GetAllNamedObjects("p")}, nil
}

// GetAllNamedObjects gets the list of objects that show up in the current named policy.
func (s *Server) GetAllNamedObjects(ctx context.Context, in *pb.GetObjectsRequest) (*pb.ObjectsReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ObjectsReply{}, err
	}

	return &pb.ObjectsReply{Objects: e.GetModel().GetValuesForFieldInPolicy("p", in.PType, 1)}, nil
}

// GetAllActions gets the list of actions that show up in the current policy.
func (s *Server) GetAllActions(ctx context.Context, in *pb.GetActionsRequest) (*pb.ActionsReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ActionsReply{}, err
	}

	return &pb.ActionsReply{Actions: e.GetAllNamedActions(in.PType)}, nil
}

// GetAllNamedActions gets the list of actions that show up in the current named policy.
func (s *Server) GetAllNamedActions(ctx context.Context, in *pb.GetActionsRequest) (*pb.ActionsReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ActionsReply{}, err
	}

	return &pb.ActionsReply{Actions: e.GetModel().GetValuesForFieldInPolicy("p", in.PType, 2)}, nil
}

// GetAllRoles gets the list of roles that show up in the current policy.
func (s *Server) GetAllRoles(ctx context.Context, in *pb.GetRolesRequest) (*pb.RolesReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.RolesReply{}, err
	}

	return &pb.RolesReply{Roles: e.GetAllNamedRoles("g")}, nil
}

// GetAllNamedRoles gets the list of roles that show up in the current named policy.
func (s *Server) GetAllNamedRoles(ctx context.Context, in *pb.GetRolesRequest) (*pb.RolesReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.RolesReply{}, err
	}

	return &pb.RolesReply{Roles: e.GetModel().GetValuesForFieldInPolicy("g", in.PType, 1)}, nil
}

// GetPolicy gets all the authorization rules in the policy.
func (s *Server) GetPolicy(ctx context.Context, in *pb.GetPolicyRequest) (*pb.PolicyReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.PolicyReply{}, err
	}

	policy := e.GetNamedPolicy("p")
	if len(policy) > 0 {
		return &pb.PolicyReply{}, nil
	}

	policyReply := &pb.PolicyReply{}
	policyReply.Policy = make([]*pb.PolicyReplyItem, len(policy))
	for e := range policy {
		policyReply.Policy[e] = &pb.PolicyReplyItem{Items: policy[e]}
	}

	return policyReply, nil
}
