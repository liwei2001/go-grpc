This is a simple go-grpc implementation.

Here are the steps to explore this project:

1. Clone the project:
    git clone https://github.com/liwei2001/go-grpc.git

2. At project root directory, start the server in docker:
    ./run_server.sh

3. The client is not running in docker, build the client at project root directory
    ./build_client.sh

4. Test the above APIs, for example:

    a. Demo usage:
    ./test_run

    b. CreateOrganization action:
    ./test_run CreateOrganization "Test Org" "Test Org Description"

    List of Actions:
    1. CreateOrganization {name} {description}
    2. FetchOrganizationList
    3. CreateUser {organization_id} {name}
    4. FetchUserList
    5. FetchUserListByOrganization {organization_id}
