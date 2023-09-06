package kindoperator

import (
	"testing"
)

func TestGetSubnetAndGateway(t *testing.T) {
	subnet := "192.168.0.0/24"
	expectedSubnet := "192.168.0.0/24"
	expectedGateway := "192.168.0.1"

	actualSubnet, actualGateway, err := getSubnetAndGateway(subnet)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualSubnet != expectedSubnet {
		t.Errorf("Expected subnet %s, but got %s", expectedSubnet, actualSubnet)
	}

	if actualGateway != expectedGateway {
		t.Errorf("Expected gateway %s, but got %s", expectedGateway, actualGateway)
	}
}

func TestGetSubnetAndGatewayInvalidSubnet(t *testing.T) {
	subnet := "invalid_subnet"
	expectedError := "invalid CIDR address: invalid_subnet"

	_, _, err := getSubnetAndGateway(subnet)
	if err == nil {
		t.Errorf("Expected error, but got nil")
	}

	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', but got '%s'", expectedError, err.Error())
	}
}

func TestCreateNetwork(t *testing.T) {
	name := "test_network"
	subnet := "192.168.0.0/24"

	err := createNetwork(name, subnet)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// TODO: Add assertions for network creation
}
