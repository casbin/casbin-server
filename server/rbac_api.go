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
	"errors"

	pb "github.com/casbin/casbin-server/proto"
)

// GetDomains gets the domains that a user has.
func (s *Server) GetDomains(ctx context.Context, in *pb.UserRoleRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}

	rm := e.GetModel()["g"]["g"].RM
	if rm == nil {
		return nil, errors.New("RoleManager is nil")
	}

	res, _ := rm.GetDomains(in.User)

	return &pb.ArrayReply{Array: res}, nil
}

// GetRolesForUser gets the roles that a user has.
func (s *Server) GetRolesForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}

	rm := e.GetModel()["g"]["g"].RM
	if rm == nil {
		return nil, errors.New("RoleManager is nil")
	}

	res, _ := rm.GetRoles(in.User, in.Domain...)

	return &pb.ArrayReply{Array: res}, nil
}

// GetImplicitRolesForUser gets implicit roles that a user has.
func (s *Server) GetImplicitRolesForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}
	res, err := e.GetImplicitRolesForUser(in.User)
	return &pb.ArrayReply{Array: res}, err
}

// GetUsersForRole gets the users that have a role.
func (s *Server) GetUsersForRole(ctx context.Context, in *pb.UserRoleRequest) (*pb.ArrayReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.ArrayReply{}, err
	}

	rm := e.GetModel()["g"]["g"].RM
	if rm == nil {
		return nil, errors.New("RoleManager is nil")
	}

	res, _ := rm.GetUsers(in.Role)

	return &pb.ArrayReply{Array: res}, nil
}

// HasRoleForUser determines whether a user has a role.
func (s *Server) HasRoleForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	roles, err := e.GetRolesForUser(in.User)
	if err != nil {
		return &pb.BoolReply{}, err
	}

	for _, r := range roles {
		if r == in.Role {
			return &pb.BoolReply{Res: true}, nil
		}
	}

	return &pb.BoolReply{}, nil
}

// AddRoleForUser adds a role for a user.
// Returns false if the user already has the role (aka not affected).
func (s *Server) AddRoleForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	ruleAdded, err := e.AddGroupingPolicy(in.User, in.Role)
	return &pb.BoolReply{Res: ruleAdded}, err
}

// DeleteRoleForUser deletes a role for a user.
// Returns false if the user does not have the role (aka not affected).
func (s *Server) DeleteRoleForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	ruleRemoved, err := e.RemoveGroupingPolicy(in.User, in.Role)
	return &pb.BoolReply{Res: ruleRemoved}, err
}

// DeleteRolesForUser deletes all roles for a user.
// Returns false if the user does not have any roles (aka not affected).
func (s *Server) DeleteRolesForUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	ruleRemoved, err := e.RemoveFilteredGroupingPolicy(0, in.User)
	return &pb.BoolReply{Res: ruleRemoved}, err
}

// DeleteUser deletes a user.
// Returns false if the user does not exist (aka not affected).
func (s *Server) DeleteUser(ctx context.Context, in *pb.UserRoleRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	ruleRemoved, err := e.RemoveFilteredGroupingPolicy(0, in.User)
	return &pb.BoolReply{Res: ruleRemoved}, err
}

// DeleteRole deletes a role.
func (s *Server) DeleteRole(ctx context.Context, in *pb.UserRoleRequest) (*pb.EmptyReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.EmptyReply{}, err
	}

	_, err = e.DeleteRole(in.Role)
	return &pb.EmptyReply{}, err
}

// DeletePermission deletes a permission.
// Returns false if the permission does not exist (aka not affected).
func (s *Server) DeletePermission(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	ruleRemoved, err := e.RemoveFilteredPolicy(1, in.Permissions...)
	return &pb.BoolReply{Res: ruleRemoved}, err
}

// AddPermissionForUser adds a permission for a user or role.
// Returns false if the user or role already has the permission (aka not affected).
func (s *Server) AddPermissionForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	ruleAdded, err := e.AddPolicy(s.convertPermissions(in.User, in.Permissions...)...)
	return &pb.BoolReply{Res: ruleAdded}, err
}

// DeletePermissionForUser deletes a permission for a user or role.
// Returns false if the user or role does not have the permission (aka not affected).
func (s *Server) DeletePermissionForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	ruleRemoved, err := e.RemovePolicy(s.convertPermissions(in.User, in.Permissions...)...)
	return &pb.BoolReply{Res: ruleRemoved}, err
}

// DeletePermissionsForUser deletes permissions for a user or role.
// Returns false if the user or role does not have any permissions (aka not affected).
func (s *Server) DeletePermissionsForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	ruleRemoved, err := e.RemoveFilteredPolicy(0, in.User)
	return &pb.BoolReply{Res: ruleRemoved}, err
}

// GetPermissionsForUser gets permissions for a user or role.
func (s *Server) GetPermissionsForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}

	return s.wrapPlainPolicy(e.GetFilteredPolicy(0, in.User)), nil
}

// GetImplicitPermissionsForUser gets all permissions(including children) for a user or role.
func (s *Server) GetImplicitPermissionsForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.Array2DReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.Array2DReply{}, err
	}
	resp, err := e.GetImplicitPermissionsForUser(in.User, in.Domain...)
	return s.wrapPlainPolicy(resp), err
}

// HasPermissionForUser determines whether a user has a permission.
func (s *Server) HasPermissionForUser(ctx context.Context, in *pb.PermissionRequest) (*pb.BoolReply, error) {
	e, err := s.getEnforcer(int(in.EnforcerHandler))
	if err != nil {
		return &pb.BoolReply{}, err
	}

	return &pb.BoolReply{Res: e.HasPolicy(s.convertPermissions(in.User, in.Permissions...)...)}, nil
}

func (s *Server) convertPermissions(user string, permissions ...string) []interface{} {
	params := make([]interface{}, 0, len(permissions)+1)
	params = append(params, user)
	for _, perm := range permissions {
		params = append(params, perm)
	}

	return params
}
