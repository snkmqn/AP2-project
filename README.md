# Electronics store microservices platform

This project is an Electronics Store application that allows users to browse, search, and purchase a wide range of electronic products such as smartphones, laptops, accessories, and home appliances. The core idea is to demonstrate how e-commerce platforms can be structured using microservices, event-driven communication with NATS, and centralized access through an API Gateway.

Service           | Description                                                                 |
|------------------|-----------------------------------------------------------------------------|
|**User Service**  | Handles user registration, authentication, and profile retrieval.           |
|**Product Service** | Manages product CRUD operations and stock levels.                         
|**Order Service** | Creates orders and emits `OrderCreated` events via NATS.                   |
|**Consumer Service** | Listens to events (e.g., `OrderCreated`) and updates stock accordingly.     |
|**Producer Service** | Publishes mock events to NATS for testing and debugging purposes.           |
|**API Gateway**   | Exposes a RESTful interface (HTTP) and proxies requests to gRPC services.  |

<br>

## Technologies used

- **Golang** – Main programming language used for building all microservices
- **MongoDB** – NoSQL database for storing product, user, and order data
- **NATS** – Lightweight, high-performance messaging system for event-driven communication
- **Redis** – In-memory data store used for caching and session management
- **gRPC** – High-performance, open-source RPC framework for communication between internal services
- **Prometheus** – Monitoring system and time series database for collecting metrics
- **Grafana** – Data visualization tool for monitoring dashboards and alerting based on Prometheus metrics


## How to run locally

### Prerequisites

- Go 1.20+
- MongoDB (default: `mongodb://localhost:27017`)
- NATS server (default: `nats://localhost:4222`)
- Protocol buffer tools (`protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`)
- Docker

### Steps to run

### 1. Clone the Repository
- git clone https://github.com/arxshi/advncdfnl2
- cd electronics-store

### 2. Set Environment Variables
- look env.example

### 3. Run Docker-compose

``` docker-compose up --build```
<br>

This will: 
- Build and start all services
- Launch NATS, Redis, Prometheus, and Grafana
- Expose relevant ports

### 3. Run Services (each in a separate terminal)



## How to run tests

### Testing tools
- **REST APIs**: Use Postman, Insomnia, or curl against the API Gateway.

- **gRPC**: Use grpcurl or gRPC client libraries for manual calls.

- **NATS events**: Use producer-service to simulate event publishing (e.g. OrderCreated).
- **Mock Unit Tests**: User microservice includes unit tests that mock external dependencies such as databases or Redis. <br>```go test -v ./internal/infrastructure/repositories```
- **Mock Integration Tests**: Can be run locally to test service interactions using mocked external services. <br> ``` go test -v -tags integration ./internal/infrastructure/repositories```


## GRPC endpoints

### ProductService

- **CreateProduct(CreateProductRequest) → ProductResponse**  
  Creates a new product in the catalog.

- **GetProductByID(GetProductByIDRequest) → ProductResponse**  
  Retrieves product information by its ID.

- **ListProducts(ListProductsRequest) → ListProductsResponse**  
  Returns a list of products with support for pagination and filtering.

- **UpdateProduct(UpdateProductRequest) → ProductResponse**  
  Updates the details of an existing product.

- **DeleteProduct(DeleteProductRequest) → DeleteProductResponse**  
  Deletes a product from the catalog.

- **CheckStock(CheckStockRequest) → CheckStockResponse**  
  Checks the stock availability of a product.

- **DecreaseStock(DecreaseStockRequest) → DecreaseStockResponse**  
  Decreases the stock quantity after an order.

### OrderService

- **CreateOrder(CreateOrderRequest) → CreateOrderResponse**  
  Creates a new order for a user.

- **GetOrderByID(GetOrderRequest) → GetOrderResponse**  
  Retrieves order details by order ID.

- **UpdateOrder(UpdateOrderRequest) → UpdateOrderResponse**  
  Updates order information (e.g., status).

- **GetOrdersByUserID(GetOrdersByUserRequest) → GetOrdersByUserResponse**  
  Retrieves a list of all orders for a user with pagination support.

### UserService

- **RegisterUser(RegisterRequest) → RegisterResponse**  
Registers a new user.

- **LoginUser(LoginRequest) → LoginResponse**  
  Authenticates a user and returns an access token.

- **RetrieveProfile(RetrieveProfileRequest) → RetrieveProfileResponse**  
  Retrieves user profile information.

- **DeleteUser(DeleteUserRequest) → DeleteUserResponse**  
  Deletes a user from the system.

## List of implemented features

- **User Management**
  - User registration and authentication (JWT-based)

- **Product Management**
  - List products with pagination and filtering
  - Stock management: check and decrease stock levels
  
- **Microservices Architecture**
  - 3 core microservices: User, Inventory, Order
  - API Gateway for routing and authentication
  - Event-driven communication via NATS (publisher and subscriber services)
- **Communication Protocols**
  - gRPC for internal microservice communication
  - REST API via API Gateway for external clients
- **Transactions & Caching**
  - MongoDB for consecutive transactions
  - Redis for caching and session management
- **Observability**
  - Metrics collection using Prometheus
  - Monitoring dashboards and alerts via Grafana
- **Testing**
  - Unit tests and integration tests
  - Mock tests and mock integration tests for reliable testing
- **Containerization**
  - Docker support for local development and deployment
  - Includes containers for Redis, NATS, and telemetry tools