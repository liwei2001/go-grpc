
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
	"strconv"
    "math/rand"
    "google.golang.org/grpc"
	pb "github.com/liwei2001/go-grpc/organization"
)

var (
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
)

func printOrganization(client pb.OrganizationServiceClient, createOrganizationRequest *pb.CreateOrganizationRequest) {

    log.Printf("Create a new organization and print it")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	organizationResponse, err := client.CreateOrganization(ctx, createOrganizationRequest)
	if err != nil {
		log.Fatalf("%v.CreateOrganization(_, _) = _, %v: ", client, err)
	}
	log.Println(organizationResponse)
}

func getOrganizationList(client pb.OrganizationServiceClient, empty *pb.Empty) (*pb.OrganizationListResponse, error) {

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    return client.FetchOrganizationList(ctx, empty)
}

func getRandomOrganizationId (organizationListResponse *pb.OrganizationListResponse) string {
    randomOrgIndex := rand.Intn(len(organizationListResponse.GetOrganizations()))
    //log.Printf("generated random org index is " + strconv.Itoa(randomOrgIndex))

    organizationId := organizationListResponse.GetOrganizations()[randomOrgIndex].Id
    //log.Printf("generated random organizationId for user is " + organizationId)

    return organizationId
}

func printOrganizationList(organizationListResponse *pb.OrganizationListResponse) {

    log.Printf("fetch organization list and print all")

    organizations := organizationListResponse.GetOrganizations()
    fmt.Println("The total number of organization is " + strconv.Itoa(len(organizationListResponse.GetOrganizations())))

    if len(organizations) > 0 {
        fmt.Printf("Id\t\t\tName\t\t\tDescription\n")
    } else {
        fmt.Println("No organizations found")
    }

    for _, org := range organizations {
        fmt.Printf("%s\t\t\t%s\t\t\t%s\n",
            org.GetId(),
            org.GetName(),
            org.GetDescription())
    }
}

func printUser(client pb.OrganizationServiceClient, createUserRequest *pb.CreateUserRequest) {

    log.Printf("Create a new user and print it")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	userResponse, err := client.CreateUser(ctx, createUserRequest)
	if err != nil {
		log.Fatalf("%v.CreateUser(_, _) = _, %v: ", client, err)
	}
	log.Println(userResponse)
}

func printUserList(client pb.OrganizationServiceClient, empty *pb.Empty) {

    log.Printf("fetch user list and print all")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    userListResponse, err := client.FetchUserList(ctx, empty)
    if err != nil {
        log.Fatalf("%v.FetchUserList(_, _) = _, %v: ", client, err)
    }

    users := userListResponse.GetUsers()
    if len(users) > 0 {
        fmt.Printf("Id\t\t\tOrganizationId\t\t\tName\n")
    } else {
        fmt.Println("No users found")
    }

    for _, user := range users {
        fmt.Printf("%s\t\t\t%s\t\t\t%s\n",
            user.GetId(),
            user.GetOrganizationId(),
            user.GetName())
    }
}

func printUserListByOrganization(client pb.OrganizationServiceClient, byOrganizationRequest *pb.ByOrganizationRequest) {

    log.Printf("fetch user list by organization and print all")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    userListResponse, err := client.FetchUserListByOrganization(ctx, byOrganizationRequest)
    if err != nil {
        log.Fatalf("%v.FetchUserList(_, _) = _, %v: ", client, err)
    }

    users := userListResponse.GetUsers()
    if len(users) > 0 {
        fmt.Printf("Id\t\t\tOrganizationId\t\t\tName\n")
    } else {
        fmt.Println("No users found")
    }

    for _, user := range users {
        fmt.Printf("%s\t\t\t%s\t\t\t%s\n",
            user.GetId(),
            user.GetOrganizationId(),
            user.GetName())
    }
}

func main() {
    flag.Parse()

    conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
    if err != nil {
        log.Fatalf("grpc.Dial err: %v", err)
    }

    client := pb.NewOrganizationServiceClient(conn)

    for i := 0; i < 5; i++ {
        printOrganization(client, &pb.CreateOrganizationRequest{Name: "Michael's Org " + strconv.Itoa(i+1), Description: "Michael's Testing Org " + strconv.Itoa(i+1)})
    }

    organizationListResponse, err := getOrganizationList(client, &pb.Empty{})
    printOrganizationList(organizationListResponse)

    for i := 0; i < 20; i++ {
        printUser(client, &pb.CreateUserRequest{OrganizationId: getRandomOrganizationId(organizationListResponse), Name: "Random User " + strconv.Itoa(i+1)})
    }

    printUserList(client, &pb.Empty{})

    randomOrgId := getRandomOrganizationId(organizationListResponse)
    log.Printf("generated random organizationId for user is " + randomOrgId)

    printUserListByOrganization(client, &pb.ByOrganizationRequest{OrganizationId : randomOrgId})
}
