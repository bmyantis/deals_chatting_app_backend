# Dealls Chatting App
System for a Dating App 

## Directory Structure
internal/
├── config/
├── constant/
├── data/
├── middleware/
├── model/
├── repository/
├── router/
└── service/

## Components
1. Config: Contains configuration files such as environment variables or application settings.
2. Constant: Houses constant values used throughout the service.
3. Data: Holds response and Request data structures and types used internally within the service.
4. Middleware: Includes middleware functions for request processing, such as authentication, logging, etc.
5. Model: Defines the data models used by the service, including database schemas.
6. Repository: Implements data access logic and interacts with the database.
7. Router: Handles HTTP request routing and request handling using a web framework such as Gin or Echo.
8. Service: Implements the core business logic of the service, including use cases and application-specific logic.
9. Controller: Handles incoming requests, processing them, and returning an appropriate response.

## Technologies Used
1. Keycloak for Authentication: Keycloak is used to securely manage user logins and permissions. It's reliable and makes it easy to add authentication features like login, signup, and user management to the app.

2. PostgreSQL: PostgreSQL is our database. It stores all the app's data securely, ensuring it's organized and easy to access. It's known for being powerful, reliable, and scalable.

3. Gomock for Unit Testing: Gomock helps us test our code. It creates pretend versions of other parts of the code so we can check if everything works as expected. This ensures our app runs smoothly without any surprises.

4. Docker for Containerization: Docker makes it easy to run our app on different computers without any setup hassles. It wraps everything our app needs into a neat package called a container, making it easy to deploy and manage. This simplifies our development and deployment processes.

## How to run it
1. Install Docker: Ensure you have Docker installed on your system. You can download and install Docker Desktop from the official Docker website.

2. Clone the Repository: Clone this repository https://github.com/bmyantis/deals_chatting_app_backend from GitHub to your local machine.

3. Set Up Keycloak Container: Run this `docker run -p 8080:8080 -e KEYCLOAK_ADMIN=admin -e KEYCLOAK_ADMIN_PASSWORD=admin quay.io/keycloak/keycloak:25.0.0 start-dev` to run keycloak and create admin account in local

4. Access Keycloak Admin Console: Open a web browser and go to http://localhost:8080/auth. This is the Keycloak admin console URL. Input username and password(admin, admin)

5. Set Up Keycloak Realm: Log in to the Keycloak admin console using the admin credentials you created. Create a new realm for your service. You can name it something relevant to your project.

6. Set Up Client in Keycloak: Within your Keycloak realm, create a new client for your service. Configure the client settings according to your service requirements. Make sure to note down the client ID and client secret.

7. Configure Service Environment Variables: In `config.go`, set up environment variables to store the Keycloak realm URL, client ID, client secret, and other relevant configurations.

8. Build and run Docker Image: 
docker-compose up --build

9. Access Your Service: You can now access your service by going to http://localhost:8090 in your web browser.


## Testing
- run this script `go test -v -coverprofile=coverage.txt ./...`
- import postman collections from this file `deals.postman_collection.json`