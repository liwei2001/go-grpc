//go:generate protoc -I ../organization --go_out=plugins=grpc:../organization ../organization/organization.proto

// It implements the organization service whose definition can be found in organization/organization.proto.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
    "github.com/rs/xid"
    "google.golang.org/grpc"
	pb "github.com/liwei2001/go-grpc/organization"
)

var (
	port       = flag.Int("port", 10000, "The server port")
)

type organization struct {
  id string
  name string
  description string
}

type userInfo struct {
  id string
  name string
}

type organizationServiceServer struct {
    organizations []organization
    users map[string][]userInfo  //organization_id as key
}


func (o *organizationServiceServer) CreateOrganization(ctx context.Context, request *pb.CreateOrganizationRequest) (*pb.OrganizationResponse, error) {
    newOrgId := xid.New().String()

    newOrganization := organization {
        id: newOrgId,
        name: request.Name,
        description: request.Description,
    }

    o.organizations = append(o.organizations, newOrganization)

    return &pb.OrganizationResponse {
        Id: newOrgId,
        Name: request.Name,
        Description: request.Description,
    }, nil
}

func (o *organizationServiceServer) FetchOrganizationList(ctx context.Context, empty *pb.Empty) (*pb.OrganizationListResponse, error) {

    organizations := make([]*pb.OrganizationResponse, 0)

    for _, org := range o.organizations {
        orgRes := &pb.OrganizationResponse{
            Id:         org.id,
            Name:     org.name,
            Description: org.description,
        }

        organizations = append(organizations, orgRes)
    }

    return &pb.OrganizationListResponse{
        Organizations: organizations,
    }, nil
}

func (o *organizationServiceServer) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.UserResponse, error) {
    newUserId := xid.New().String()

    newUserInfo := userInfo {
        id: newUserId,
        name: request.Name,
    }

    userInfoList := o.users[request.OrganizationId]

    userInfoList = append(userInfoList, newUserInfo)

    o.users[request.OrganizationId] = userInfoList

    return &pb.UserResponse {
        Id: newUserId,
        OrganizationId: request.OrganizationId,
        Name: request.Name,
    }, nil
}

func (o *organizationServiceServer) FetchUserList(ctx context.Context, empty *pb.Empty) (*pb.UserListResponse, error) {

    users := make([]*pb.UserResponse, 0)

    for key, userInfoList := range o.users {
        for i := 0; i < len(userInfoList); i++ {
            userRes := &pb.UserResponse {
                Id:         userInfoList[i].id,
                OrganizationId: key,
                Name:     userInfoList[i].name,
            }

            users = append(users, userRes)
        }
    }

    return &pb.UserListResponse{
        Users: users,
    }, nil
}

func (o *organizationServiceServer) FetchUserListByOrganization(ctx context.Context, request *pb.ByOrganizationRequest) (*pb.UserListResponse, error) {
    var users []*pb.UserResponse

    userInfoList := o.users[request.OrganizationId]

    for i := 0; i < len(userInfoList); i++ {
        userRes := &pb.UserResponse{
            Id:         userInfoList[i].id,
            OrganizationId: request.OrganizationId,
            Name:     userInfoList[i].name,
        }

        users = append(users, userRes)
    }

    return &pb.UserListResponse{
        Users: users,
    }, nil
}


func newServer() *organizationServiceServer {
	s := &organizationServiceServer {
	    organizations: make([]organization, 0),
	    users: make(map[string][]userInfo),
	}

	return s
}


func main() {

    flag.Parse()
    lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    log.Println("Listening on ", *port)
    var opts []grpc.ServerOption

    server := grpc.NewServer(opts...)

    pb.RegisterOrganizationServiceServer(server, newServer())

    if err := server.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
