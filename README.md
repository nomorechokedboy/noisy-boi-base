# Notification Server

This is a notification server implemented in Go, leveraging Server-Sent Events (SSE) to push real-time notifications to clients. You can find out more details in this [post](https://dev.to/hadius/go-server-sent-events-2ng)

## Features

Real-time notifications: Clients can receive notifications in real time without the need for manual polling.
Simple and lightweight: The server is built using Go, which provides a performant and efficient runtime.
Scalable: The server is designed to handle a large number of concurrent connections, making it suitable for high-traffic applications.

## Prerequisites

Go (version 1.18 or above)
NodeJs (14+)

## Getting Started

1. Clone the repository:

```sh
git clone https://github.com/nomorechokedboy/noisy-boi-base.git
```

2. Navigate to the project directory:

```bash
cd noisy-boi-base/sse-server
```

3. Build the project:

```bash
go build -o api .
```

4. Run the server:

```bash
./api
```

The server should now be running on the default port (e.g., http://localhost:3500).
