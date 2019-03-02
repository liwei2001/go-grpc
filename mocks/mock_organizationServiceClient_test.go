package mocks_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	orgService "scytale/organization/organization"
	mocks "scytale/organization/mocks"
)


func TestCreateOrganization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrganizationServiceClient := mocks.NewMockOrganizationServiceClient(ctrl)
	req := &orgService.CreateOrganizationRequest{Name: "unit_test_org_name", Description: "unit_test_org_description"}
	mockOrganizationServiceClient.EXPECT().CreateOrganization(
		gomock.Any(),
		req,
	).Return(&orgService.OrganizationResponse{Id: "unit_test_org_id", Name: "unit_test_org_name", Description: "unit_test_org_description"}, nil)
	testCreateOrganization(t, mockOrganizationServiceClient)
}

func testCreateOrganization(t *testing.T, client orgService.OrganizationServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.CreateOrganization(ctx, &orgService.CreateOrganizationRequest{Name: "unit_test_org_name", Description: "unit_test_org_description"})
	if err != nil || r.Id != "unit_test_org_id" {
		t.Errorf("mocking failed")
	}
	t.Log("New organization created with organizationId : ", r.Id)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockOrganizationServiceClient := mocks.NewMockOrganizationServiceClient(ctrl)
	req := &orgService.CreateUserRequest{OrganizationId: "unit_test_org_id", Name: "unit_test_user_name"}
	mockOrganizationServiceClient.EXPECT().CreateUser(
		gomock.Any(),
		req,
	).Return(&orgService.UserResponse{Id: "unit_test_user_id", OrganizationId: "unit_test_org_id", Name: "unit_test_user_name"}, nil)
	testCreateUser(t, mockOrganizationServiceClient)
}

func testCreateUser(t *testing.T, client orgService.OrganizationServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.CreateUser(ctx, &orgService.CreateUserRequest{OrganizationId: "unit_test_org_id", Name: "unit_test_user_name"})
	if err != nil || r.Id != "unit_test_user_id" {
		t.Errorf("mocking failed")
	}
	t.Log("New user created with userId : ", r.Id)
}