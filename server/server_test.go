package main

import (
    "context"
    "fmt"
    "testing"
    "time"
    pb "github.com/liwei2001/go-grpc/organization"
)

func TestCreateOrganization(t *testing.T) {

  request1 := pb.CreateOrganizationRequest {
    Name: "Michael's Test Org",
    Description: "Michael's Test Org description",
  }

  request2 := pb.CreateOrganizationRequest{
      Name: "",
      Description: "Michael's Test Org description",
    }

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  organizationServiceServer := newServer()
   response1, error1 := organizationServiceServer.CreateOrganization(ctx, &request1)

   if error1 == nil {
        fmt.Println("Organization created successfully with organizationId: ", response1.Id)
    } else {
        t.Errorf("Organization not created with Name: %s, Description: %s", request1.Name, request1.Description)
    }

    _, error2 := organizationServiceServer.CreateOrganization(ctx, &request2)

   if error2 != nil {
        fmt.Println("Organization with empty name cannot be created")
    } else {
        t.Errorf("Organization with empty name should not be created")
    }
}

func TestFetchOrganizationList(t *testing.T) {

  orgCreationRequest1 := pb.CreateOrganizationRequest {
    Name: "Michael's Test Org I",
    Description: "Michael's Test Org description",
  }

  orgCreationRequest2 := pb.CreateOrganizationRequest {
      Name: "Michael's Test Org II",
      Description: "Michael's Test Org description",
  }

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  organizationServiceServer := newServer()

    organizationServiceServer.CreateOrganization(ctx, &orgCreationRequest1)
    organizationServiceServer.CreateOrganization(ctx, &orgCreationRequest2)

   response, error := organizationServiceServer.FetchOrganizationList(ctx, &pb.Empty{})

   if error == nil {
        fmt.Println("Number of organizations is: ", len(response.Organizations))
    } else {
        t.Errorf("Error retrieving organizations")
    }
}

func TestCreateUser(t *testing.T) {

  orgCreationRequest := pb.CreateOrganizationRequest {
    Name: "Michael's Test Org",
    Description: "Michael's Test Org description",
  }

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  organizationServiceServer := newServer()
   response, _ := organizationServiceServer.CreateOrganization(ctx, &orgCreationRequest)

   organizationId := response.Id

  request1 := pb.CreateUserRequest {
    OrganizationId: organizationId,
    Name: "Michael's Test User",
  }

  request2 := pb.CreateUserRequest{
    OrganizationId: "No-existing OrgId",
     Name: "Michael's Test User",
  }

  request3 := pb.CreateUserRequest{
    OrganizationId: "Michael's Test Org",
    Name: "",
  }

    response1, error1 := organizationServiceServer.CreateUser(ctx, &request1)

   if error1 == nil {
        fmt.Println("User created successfully with userId: ", response1.Id)
    } else {
        t.Errorf("User not created with Name: %s, Description: %s", request1.OrganizationId, request1.Name)
    }

    _, error2 := organizationServiceServer.CreateUser(ctx, &request2)

    if error2 != nil {
          fmt.Println("User with non-existing orgId cannot be created")
     } else {
          t.Errorf("User with non-existing orgId should not be created")
    }

    _, error3 := organizationServiceServer.CreateUser(ctx, &request3)

   if error3 != nil {
        fmt.Println("User with empty name cannot be created")
    } else {
        t.Errorf("User with empty name should not be created")
    }
}

func TestFetchUserList(t *testing.T) {

  orgCreationRequest1 := pb.CreateOrganizationRequest {
    Name: "Michael's Test Org I",
    Description: "Michael's Test Org description",
  }

  orgCreationRequest2 := pb.CreateOrganizationRequest {
      Name: "Michael's Test Org II",
      Description: "Michael's Test Org description",
  }

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  organizationServiceServer := newServer()

   orgResponse1, _ := organizationServiceServer.CreateOrganization(ctx, &orgCreationRequest1)
   organizationId1 := orgResponse1.Id

   orgResponse2, _ := organizationServiceServer.CreateOrganization(ctx, &orgCreationRequest2)
   organizationId2 := orgResponse2.Id

  userCreationRequest1 := pb.CreateUserRequest {
      OrganizationId: organizationId1,
      Name: "Michael's Test User I",
    }

  userCreationRequest2 := pb.CreateUserRequest{
      OrganizationId: organizationId2,
       Name: "Michael's Test User II",
    }

  userCreationRequest3 := pb.CreateUserRequest{
      OrganizationId: organizationId2,
       Name: "Michael's Test User III",
    }

  organizationServiceServer.CreateUser(ctx, &userCreationRequest1)
  organizationServiceServer.CreateUser(ctx, &userCreationRequest2)
  organizationServiceServer.CreateUser(ctx, &userCreationRequest3)

   response, error := organizationServiceServer.FetchUserList(ctx, &pb.Empty{})

   if error == nil {
        fmt.Println("Number of users retrieved from all organizations is", len(response.Users))
    } else {
        t.Errorf("Error retrieving users from all organizations")
    }
}

func TestFetchUserListByOrganization(t *testing.T) {

  orgCreationRequest := pb.CreateOrganizationRequest {
    Name: "Michael's Test Org",
    Description: "Michael's Test Org description",
  }

  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()

  organizationServiceServer := newServer()
   response, _ := organizationServiceServer.CreateOrganization(ctx, &orgCreationRequest)

   organizationId := response.Id

  userCreationRequest1 := pb.CreateUserRequest {
      OrganizationId: organizationId,
      Name: "Michael's Test User I",
    }

  userCreationRequest2 := pb.CreateUserRequest{
      OrganizationId: organizationId,
       Name: "Michael's Test User II",
    }

  organizationServiceServer.CreateUser(ctx, &userCreationRequest1)
  organizationServiceServer.CreateUser(ctx, &userCreationRequest2)

  request1 := pb.ByOrganizationRequest {OrganizationId: organizationId}

  request2 := pb.ByOrganizationRequest{OrganizationId: "No-existing OrgId"}

   response1, error1 := organizationServiceServer.FetchUserListByOrganization(ctx, &request1)

   if error1 == nil {
        fmt.Println("Number of users retrieved from organizationId: " + request1.OrganizationId + " is", len(response1.Users))
    } else {
        t.Errorf("Error retrieving users from the organization")
    }

    _, error2 := organizationServiceServer.FetchUserListByOrganization(ctx, &request2)

    if error2 != nil {
          fmt.Println("Cannot retrieve users from non-existing orgId")
     } else {
          t.Errorf("No users should be retrieved from non-existing orgId")
    }

}