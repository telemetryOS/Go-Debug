package debug_test

import (
	"fmt"
	"time"

	"github.com/telemetryos/Go-Debug/debug"
)

// Initialize different debug scopes for different parts of the application
var (
	dbDebug   = debug.Bind("app:db")
	httpDebug = debug.Bind("app:http")
	mainDebug = debug.Bind("app:main")
)

func simulateDBQuery() {
	dbDebug.Trace("Opening database connection")
	time.Sleep(100 * time.Millisecond)
	dbDebug.Tracef("Executing query with timeout: %d ms", 500)
	time.Sleep(200 * time.Millisecond)
	dbDebug.Trace("Closing database connection")
}

func simulateHTTPRequest() {
	httpDebug.Trace("Received HTTP request")
	time.Sleep(50 * time.Millisecond)
	httpDebug.Tracef("Processing request to path: %s", "/api/users")
	time.Sleep(150 * time.Millisecond)
	httpDebug.Trace("Sending HTTP response")
}

func Example() {
	mainDebug.Trace("Application starting")

	fmt.Println("Running application...")
	fmt.Println("To see debug traces, run with DEBUG environment variable set:")
	fmt.Println("DEBUG=app:main,app:http,app:db go run main.go  # Show all traces")
	fmt.Println("DEBUG=app:db go run main.go                    # Show only database traces")
	fmt.Println("DEBUG=* go run main.go                         # Show all traces with wildcard")

	mainDebug.Trace("Initializing services")

	simulateHTTPRequest()
	simulateDBQuery()

	mainDebug.Trace("Application shutting down")
	fmt.Println("Application finished")

	// Output:
	// Running application...
	// To see debug traces, run with DEBUG environment variable set:
	// DEBUG=app:main,app:http,app:db go run main.go  # Show all traces
	// DEBUG=app:db go run main.go                    # Show only database traces
	// DEBUG=* go run main.go                         # Show all traces with wildcard
	// Application finished
}
