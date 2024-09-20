package dto

import (
	pb "github.com/casbin/casbin-server/proto"
)

type EnforceRequest struct {
	pb.EnforceRequest
}

type EnforceResponse struct {
	pb.BoolReply
}
