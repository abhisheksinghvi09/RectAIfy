# RectAify

RectAify is an AI-Powered Startup Idea Rectifier that enables users to refine their startup ideas with the help of AI. It evaluates metrics like market power and other relevant factors to guide users in improving their business concepts.

## Features

- AI-driven analysis of startup ideas.
- Metrics evaluation including market power and potential.
- Rate-limiting and caching for optimal performance.
- Secure authentication with bearer tokens.

## System Architecture

The following diagram illustrates the system architecture of RectAify:

```plaintext
+---------------+        +-----------------+        +------------------+
|               |        |                 |        |                  |
|    Client     +-------->  RectAify API   +-------->   OpenAI API     |
|   (Frontend)  |        |                 |        |                  |
+---------------+        +-----------------+        +------------------+
       |                         ^
       |                         |
       v                         |
+---------------+        +-----------------+        +------------------+
|               |        |                 |        |                  |
|    Database   <--------+  Rate Limiter   +-------->   Cache Layer    |
|   (Postgres)  |        |                 |        |                  |
+---------------+        +-----------------+        +------------------+
```

## Installation and Setup

### Prerequisites

- Go (latest version)
- PostgreSQL
- OpenAI API Key

### Clone the Repository

```bash
git clone https://github.com/abhisheksinghvi09/RectAIfy.git
cd RectAIfy
```

### Environment Configuration

Create a `.env` file in the root directory or use the provided `.env.example` file:

```dotenv
# Example Environment File
OPENAI_API_KEY=your-api-key-here
DB_DSN=postgres://your-user@localhost:5432/rectaify?sslmode=disable
HTTP_ADDR=:9444
LOG_LEVEL=info
```

### Database Setup

Start PostgreSQL and create a new database:

```bash
psql -U your-user -c "CREATE DATABASE rectaify;"
```

### Run the Application

Build and run the application:

```bash
go build -o rectaify
./rectaify
```

The server will start at `http://localhost:9444`.

### Testing

Run the tests to ensure everything is set up correctly:

```bash
go test ./...
```

### Contributing

Feel free to fork the repository and submit pull requests for new features, bug fixes, or documentation improvements.

## License

This project is currently unlicensed. Contact the repository owner for more information.
