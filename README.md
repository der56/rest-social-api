# Rest Social API - Go

This is a RESTful API built with Go, which provides basic user authentication and profile management functionalities. The API supports user registration, login, logout, profile management (view, update username/password), and user following features, all protected using JWT for authentication. It uses a PostgreSQL database to store user and follower data.

## Features

- **Register**: Create a new user with username, email, and password.
- **Login**: Authenticate a user and generate a JWT token.
- **Logout**: Invalidate the JWT token to log the user out.
- **Profile Management**:
  - `/profile`: View your own profile information.
  - `/update-username`: Update your username.
  - `/update-password`: Update your password.
  - `/profile/:id`: View the profile of another user by their unique UUID.
- **Follow System**:
  - `/followuser/:id`: Follow a user by their UUID.
  - `/unfollowuser/:id`: Unfollow a user by their UUID.
  - `/getfollowers/:id`: Get a list of followers for a user.

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/your-repo-name/rest-social-api-go.git
cd rest-social-api-go
```

### 2. Set up the environment

Create a `.env` file in the root directory and fill it with the following environment variables:

```
PORT=your-port-number
POSTGRES_URL=your-postgres-url
SECRET_JWT=your-jwt-secret-key
```

### 3. Set up the PostgreSQL database

Run the following SQL commands in your PostgreSQL database to create the necessary tables:

```sql
BEGIN;

-- Enable UUID extension if it doesn't exist
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create user_profile table
CREATE TABLE IF NOT EXISTS user_profile (
    ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    Username VARCHAR(20) NOT NULL UNIQUE,
    Email VARCHAR(150) NOT NULL UNIQUE,
    Password VARCHAR(45) NOT NULL,
    FirstName VARCHAR(18) NOT NULL,
    LastName VARCHAR(50) NOT NULL,
    DateOfEntry TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    Picture VARCHAR(255),
    Description VARCHAR(255),
    CONSTRAINT chk_username_min_length CHECK (CHAR_LENGTH(Username) >= 3),
    CONSTRAINT chk_password_min_length CHECK (CHAR_LENGTH(Password) >= 6)
);

-- Create followers table
CREATE TABLE IF NOT EXISTS followers (
    follower_id UUID NOT NULL,
    following_id UUID NOT NULL,
    followed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id),
    FOREIGN KEY (follower_id) REFERENCES user_profile(ID) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES user_profile(ID) ON DELETE CASCADE,
    CONSTRAINT chk_self_follow CHECK (follower_id != following_id)
);

COMMIT;
```

### 4. Install Go dependencies

Make sure to install the dependencies defined in `go.mod` and `go.sum` by running:

```bash
go mod tidy
```

### 5. Running the Application

You have two options to run the application in development:

- **Using `air` (Hot reload)**:
  - Go to the `/system/api/src` directory.
  - Run `air` for live reloading during development:

    ```bash
    cd /system/api/src
    air
    ```

- **Using `go run`**:
  - If you're on Linux, you can use the `Makefile` to start the app.
  - If you're on Windows, use `Makefile.ps1` to run the application.

### 6. Running Tests

If you want to run tests for the API, you can use the following command:
```bash
go test -v ./system/api/src/tests/
```

## Endpoints

### Authentication

- **POST /register**: Register a new user.

  **Request Body**:
  ```json
  {
    "Username": "yourusername", 
    "Password": "yourpass", 
    "Firstname": "yourfirstname", 
    "Lastname": "yourlastname", 
    "Email": "youremail@email.com"
  }
  ```

- **POST /login**: Login and get a JWT token.

  **Request Body**:
  ```json
  {
    "username": "yourusername", 
    "password": "yourpass"
  }
  ```
  Or:
  ```json
  {
    "email": "youremail@email.com", 
    "password": "yourpass"
  }
  ```

- **POST /logout**: Logout by invalidating the JWT token.

### User Profile

- **GET /profile**: Get your own profile.
- **POST /update-username**: Update your username.

  **Request Body**:
  ```json
  {
    "username": "yournewusername"
  }
  ```

- **POST /update-password**: Update your password.

  **Request Body**:
  ```json
  {
    "password": "yournewpassword"
  }
  ```

- **GET /profile/:id**: Get the profile of another user by UUID.

### Follow System

- **POST /followuser/:id**: Follow a user.
- **POST /unfollowuser/:id**: Unfollow a user.
- **GET /getfollowers/:id**: Get a list of followers of a user.
