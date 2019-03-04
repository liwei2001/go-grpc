//go:generate protoc -I ../organization --go_out=plugins=grpc:../organization ../organization/organization.proto

// It implements the organization service whose definition can be found in organization/organization.proto.

package main

import (
	"context"
	"flag"
	"log"
	"net"
    "fmt"
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

type argError struct {
   message string
}

func (e *argError) Error() string {
    return fmt.Sprintf("%s", e.message)
}

func (o *organizationServiceServer) CreateOrganization(ctx context.Context, request *pb.CreateOrganizationRequest) (*pb.OrganizationResponse, error) {

    if len(request.Name) == 0 {
        return nil, &argError{"Orgnaization name cannot be empty"}
    }

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
    if len(request.Name) == 0 {
            return nil, &argError{"User name cannot be empty"}
    }

    newUserId := xid.New().String()

    newUserInfo := userInfo {
        id: newUserId,
        name: request.Name,
    }

    organizationExists := false

    for _, org := range o.organizations {
        if org.id == request.OrganizationId {
            organizationExists = true
        }
    }

    if organizationExists {

        userInfoList := o.users[request.OrganizationId]

        userInfoList = append(userInfoList, newUserInfo)

        o.users[request.OrganizationId] = userInfoList

        return &pb.UserResponse {
            Id: newUserId,
            OrganizationId: request.OrganizationId,
            Name: request.Name,
        }, nil
    } else {
        return nil, &argError{"No corresponding organization exists"}
    }
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

    organizationExists := false

    for _, org := range o.organizations {
        if org.id == request.OrganizationId {
            organizationExists = true
        }
    }

    if organizationExists {

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
    } else {
        return &pb.UserListResponse{
            Users: users,
        }, &argError{"No corresponding organization exists"}
    }
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
    //lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
    lis, err := net.Listen("tcp", ":3000")

    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    //log.Println("Listening on ", *port)
    log.Println("Listening on port 3000")

    var opts []grpc.ServerOption

    server := grpc.NewServer(opts...)

    pb.RegisterOrganizationServiceServer(server, newServer())

    if err := server.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
