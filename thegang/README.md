# thegang

This folder contains files that are helpful for the vendor development and **will not be deployed to production
containers**.

The folder can be deleted after transition to the client and the vendor do not need to maintain the application.

### HOW TO COMPOSE?

Move the docker-compose to the root and change it to

```
web:
    build:
      context: .
      dockerfile: ./thegang/Dockerfile
      target: dev
```