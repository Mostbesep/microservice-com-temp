# microservice-com-temp
# Microservice Commercial Template (microservice-com-temp)
### A commercial-grade template for building a scalable and maintainable microservices architecture, using Golang, GraphQL, gRPC, Elasticsearch, and PostgreSQL.
## Overview
The microservice-com-temp template provides a foundation for building a robust and efficient microservices-based system. This project includes:
A set of Go (Golang) services that communicate with each other using gRPC
A GraphQL API layer for exposing the service functionality to clients
Elasticsearch integration for storing and querying data
PostgreSQL database for persisting business logic
### Features
### Service Architecture
The template consists of multiple microservices, each responsible for a specific domain or capability. These services are designed to be loosely coupled, allowing for independent development, deployment, and scaling.
### GraphQL API
A GraphQL API layer is provided on top of the service architecture, enabling clients to query and mutate data using a simple, intuitive interface.
### gRPC Communication
Services communicate with each other using gRPC, which provides efficient, high-performance communication over HTTP/2 or TCP.
### Elasticsearch Integration
Elasticsearch is used for storing and querying large amounts of data. This allows for fast search capabilities and scalable indexing.
### PostgreSQL Database
A PostgreSQL database is provided to store business logic and persist service state.
## Getting Started
To get started with the microservice-com-temp template, follow these steps:
Clone this repository: git clone https://github.com/Mostbesep/microservice-com-temp.git
Install dependencies using Go modules: go mod init && go mod tidy
Run the services using Docker (recommended): docker-compose up -d
Use a GraphQL client library to interact with the API


**GraphQL API Documentation**

### Overview

This GraphQL API provides endpoints for managing accounts, products, and orders. It supports CRUD operations on these entities.

### Scalars

* `Time`: A scalar type representing a timestamp in ISO format (e.g., "2022-01-01T12:00:00Z").

### Types

#### Account

| Field | Type | Description |
| --- | --- | --- |
| id | String! | Unique identifier for the account. |
| name | String! | Name of the account holder. |
| orders | [Order!]! | List of orders associated with this account. |

#### Product

| Field | Type | Description |
| --- | --- | --- |
| id | String! | Unique identifier for the product. |
| name | String! | Name of the product. |
| description | String | Brief description of the product. |
| price | Float! | Price of the product in decimal format (e.g., 19.99). |

#### Order

| Field | Type | Description |
| --- | --- | --- |
| id | String! | Unique identifier for the order. |
| createdAt | Time! | Timestamp when the order was created. |
| totalPrice | Float! | Total price of all products in this order. |
| products | [OrderedProduct!]! | List of ordered products associated with this order. |

#### OrderedProduct

| Field | Type | Description |
| --- | --- | --- |
| id | String! | Unique identifier for the product being ordered. |
| name | String! | Name of the product being ordered. |
| description | String | Brief description of the product being ordered. |
| price | Float! | Price of the product in decimal format (e.g., 19.99). |
| quantity | Int! | Number of units of this product being ordered. |

### Inputs

#### PaginationInput

| Field | Type | Description |
| --- | --- | --- |
| skip | Int! | Skip a specified number of records before returning results. |
| take | Int! | Return only the first N records matching the query criteria. |

#### AccountInput

| Field | Type | Description |
| --- | --- | --- |
| name | String! | Name of the account holder to create or update. |

#### ProductInput

| Field | Type | Description |
| --- | --- | --- |
| name | String! | Name of the product to create or update. |
| description | String | Brief description of the product to create or update. |
| price | Float! | Price of the product in decimal format (e.g., 19.99). |

#### OrderProductInput

| Field | Type | Description |
| --- | --- | --- |
| id | String! | Unique identifier for the product being ordered. |
| quantity | Int! | Number of units of this product to order. |

#### OrderInput

| Field | Type | Description |
| --- | --- | --- |
| AccountId | String! | ID of the account associated with this order. |
| Products | [OrderProductInput!]! | List of ordered products for this order. |

### Mutations

* `createAccount(account: AccountInput!)`: Creates a new account.
    + Input fields:
        - name (String!)
* `createProduct(product: ProductInput!)`: Creates a new product.
    + Input fields:
        - name (String!)
        - description (String)
        - price (Float!)
* `createOrder(order: OrderInput!)`: Creates a new order.
    + Input fields:
        - AccountId (String!)
        - Products ([OrderProductInput!]!)

### Queries

* `accounts(pagination: PaginationInput, id: String): [Account!]!` : Retrieves accounts matching the specified criteria.
    + Optional input fields:
        - pagination (PaginationInput)
        - id (String)
* `products(pagination: PaginationInput, query:String, id:String): [Product!]!`: Retrieves products matching the specified criteria.
    + Optional input fields:
        - pagination (PaginationInput)
        - query (String)
        - id (String)

### Example Queries

```graphql
query {
  accounts(id: "12345") {
    name
    orders {
      createdAt
      totalPrice
      products {
        name
        quantity
      }
    }
  }

  products(query: "phone", pagination: { skip: 0, take: 10 }) {
    id
    name
    description
    price
  }
}
```

```graphql
mutation {
  createAccount(account: { name: "John Doe" }) {
    id
    name
  }

  createProduct(product: { name: "iPhone", description: "Smartphone", price: 999.99 }) {
    id
    name
    description
    price
  }

  createOrder(order: { AccountId: "12345", Products: [{ id: "ABCDEF", quantity: 2 }] }) {
    id
    createdAt
    totalPrice
    products {
      name
      quantity
    }
  }
}
```

## Contributing
We welcome contributions and feedback! If you'd like to contribute or report an issue, please follow these steps:
Fork this repository: git fork https://github.com/Mostbesep/microservice-com-temp.git
Create a new branch for your changes
Make your changes and commit them
Open a pull request against the original repository
## License
The microservice-com-temp template is licensed under the MIT License.
Acknowledgments
This project was inspired by various open-source projects, including:
GraphQL-go: A Go implementation of GraphQL
gRPC-go: The official Go implementation of gRPC
Elasticsearch: A search and analytics engine
We appreciate their contributions to the development of these technologies!