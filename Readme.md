# Go API with Redis Caching

This application is a Go-based API that provides endpoints for managing users, pets, and adoptions. It uses Redis for caching responses to improve performance and reduce load on the primary database.

## Features

- **User Management**: APIs to get all users, retrieve a user by ID, login, create, update, and delete users.
- **Pet Management**: APIs to manage pets, including retrieving all pets and performing CRUD operations on individual pets.
- **Adoption Management**: APIs to handle adoptions, including retrieving adoption history and performing CRUD operations.
- **Caching with Redis**: The application caches responses in Redis to speed up frequent requests.


## Docker
	docker build -t reverse-proxy:lastest .


