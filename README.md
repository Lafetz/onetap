The system is built using a module-based architecture to enable scalability and maintain a clear separation of concerns. Key components include:

## Introduction

This document explains the design, architecture, and decisions behind the implementation of a Tier Management and Cashback System. It outlines the technical solutions for managing customer tiers, cashback rewards, and real-time notifications.

## Design Overview

<img width="1028" alt="image" src="https://raw.githubusercontent.com/Lafetz/onetap/main/imgs/1.png">

The system is built using a A modular monolith architecture to enable scalability and maintain a clear separation of concerns. Key components include:

- **Loyalty Management System**: A modular service designed to manage customer loyalty programs.
- **Real-Time Notification Service**: Sends real-time updates to customers via Socket.IO when they are promoted to a new tier or receive cashback.

## Loyalty Management

The Loyalty Management System is a modular service designed to handle both tier management and cashback rewards. This system is constructed to be easily extendable, allowing for the addition of new modules or features as needed.

For programs with independent procedures—such as tier management or cashback rewards—a dedicated queue is assigned to handle specific tasks like tier promotions or cashback calculations. This design ensures that each program can process its own operations asynchronously, without interfering with others, while maintaining a clean separation of data and logic. This modular approach enhances scalability and simplifies future program expansions.

For each loyalty program, I decided to create a new schema and a separate customer instance to ensure independent operation with its own rules, tiers, and cashback configurations. Initially, MongoDB was chosen for its flexibility and schema-less design, which suited the rapid development of this test project. However, I anticipated that handling complex scenarios, such as demoting customers when a tier is removed and managing cascading updates, would require efficient joins. Given the complexity of these relational operations, I switched to an SQL database to simplify queries and improve performance, ensuring better management of relational data. Additionally, for merchants with many customers, a Redis cache could be added in the future to reduce database load.

The cashback management system was designed similarly to the tier management system, with the primary difference being the schema used to store and manage cashback-related data. While both systems operate independently and follow their own sets of rules and configurations, the schema for cashback focuses on tracking rewards, percentages, and expiration periods, distinct from the tier management schema that deals with tier levels and customer promotions.
## Notification
The notification service listens to messages from Redis Pub/Sub and checks if the user is online. If the user is available, it forwards the message to them. This ensures real-time delivery of notifications to active users.
## Running the Project
To run the project, first use Docker Compose to start the required services. This will spin up three services:

    mock_order: Sends a POST request to /order, simulating a customer order and publishing the data to Redis Pub/Sub.
    HTTP Service: Consumes the message, handling loyalty tier management and cashback management. If the user is new, it creates customer models for both services.
    Notification Service: Sends notifications if a new customer is detected.
    Merchant ID: 4f981b0b-accf-4eb7-8018-7cd651c7e907.(mock id)

Testing steps:

    Sign in using the mock sign-in route.
    Create loyalty tier and cashback entries.
    Trigger the /order route in the mock_order service to simulate an order.
    To test notifications, first connect to the service, then send a message with id and event parameters, using 4f981b0b-accf-4eb7-8018-7cd651c7e922 for the mock customer ID.
    Additionally, there are unit and integration tests available for the HTTP and Loyalty services to ensure functionality.as well as a Postman collection for testing the endpoints.