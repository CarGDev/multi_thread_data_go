# multi_thread_data_go

## Getting started (no Go installed)

### 1. Requirements
- **Go**, the version in the `go` line of `go.mod` or newer.
- **Git** (optional; you can download the project as a ZIP instead).

The project only uses Go's standard library, so there is nothing else to install.

### 2. Install Go

Download the installer for your system from <https://go.dev/dl/> and follow it, or use a package manager:

| System | Command |
|---|---|
| macOS (Homebrew) | `brew install go` |
| Windows (winget) | `winget install GoLang.Go` |
| Ubuntu / Debian | `sudo apt install golang-go` (may be an older version than `go.mod` requires; if so use the installer from go.dev) |
| Linux (manual) | extract the `.tar.gz` to `/usr/local` and add `export PATH=$PATH:/usr/local/go/bin` to your shell profile |

Open a **new terminal** and check it works:

```
go version
```

### 3. Get the code

```
git clone <repository-url>
cd multi_thread_data_go
```

(or download the ZIP, extract it and open a terminal in that folder).

### 4. Run it

From the project root, the folder that contains `go.mod`:

```
go run ./cmd/ridesharing
```

The CSV files are embedded in the program, so no paths or extra files are needed.

### 5. What you should see

- Console: `loaded ... drivers`, `loaded ... riders / ride requests`, a stream of worker/task log lines, and a final line like
  `done in 8.1s: 500 total, 500 successful, 0 failed (results in results.csv)`.
- New files in the project folder: `results.csv` (one row per task, see columns below) and `ridesharing.log` (the same log lines).

### 6. Optional checks

```
go vet ./...                          # static checks
go run -race ./cmd/ridesharing        # detects data races (slower)
go build -o ridesharing ./cmd/ridesharing && ./ridesharing   # build a binary (ridesharing.exe on Windows)
```

### 7. Troubleshooting

| Problem | Fix |
|---|---|
| `go: command not found` | Go is not on your `PATH`; reinstall or add it, then open a new terminal |
| `go.mod requires go >= X` | Your Go is too old; install a newer version from go.dev |
| `directory not found` / `no required module` | Run the command from the folder that contains `go.mod` |
| `pattern Drivers.csv: no matching files found` | `Drivers.csv` and `Riders.csv` must be inside `lib/simulation/` |

## Class Diagram

```markdown
classDiagram
    direction TB

    class RideSharingSystem {
        -TaskQueue taskQueue
        -ResultStore resultStore
        -Logger logger
        -Worker[] workers
        -bool running
        +initialize(workerCount) RideSharingSystem
        +submitTask(task)
        +submitTasks(tasks)
        +startWorkers()
        +shutdown()
        +awaitCompletion()
        +getResults() Result[]
        +writeResults(path)
        +stats() Stats
        +run(tasks, outPath)
        +close()
    }

    class Stats {
        +int total
        +int success
        +int failed
    }

    class Task {
        <<interface>>
        +process() Result
        +getId() int
        +getStatus() TaskStatus
    }

    class Base {
        -int taskId
        -TaskStatus status
        -DateTime createdAt
        -Mutex lock
        +init(taskId)
        +getId() int
        +getStatus() TaskStatus
        +setStatus(status)
        +getCreatedAt() DateTime
    }

    class RideRequest {
        -Rider rider
        -Location pickup
        -Location destination
        +process() Result
        +findDriver() Driver
        +calculateFare() double
        +tripDistanceKm() double
        +getRider() Rider
        +getPickup() Location
        +getDestination() Location
    }

    class Simulation {
        <<package>>
        +loadDrivers()
        +loadRides() Task[]
    }

    class TaskQueue {
        -chan Task tasks
        -RWMutex lock
        -bool closed
        +enqueue(task)
        +dequeue() Task
        +tryDequeue() Task, bool
        +size() int
        +isEmpty() bool
        +close()
        +isClosed() bool
    }

    class Worker {
        -int workerId
        -TaskQueue taskQueue
        -ResultStore resultStore
        -Logger logger
        -bool running
        +start()
        +run()
        +processTask(task)
        +stop()
        +getId() int
        +isRunning() bool
    }

    class Result {
        -int taskId
        -bool success
        -string message
        -DateTime completedAt
        +getTaskId() int
        +isSuccess() bool
        +getMessage() string
        +getCompletedAt() DateTime
        +setWorker(workerId)
        +getWorkerId() int
        +setStartedAt(time)
        +getDuration() Duration
        +setRide(info)
        +getRide() RideInfo
    }

    class RideInfo {
        +int riderId
        +string riderName
        +Location pickup
        +Location destination
        +int driverId
        +string driverName
        +Location driverStart
        +Location driverEnd
        +double driverToPickupKm
        +double tripDistanceKm
        +double fare
    }

    class ResultStore {
        -List~Result~ results
        -Mutex lock
        +addResult(result)
        +getResults() List~Result~
        +count() int
        +countSuccess() int
        +countFailed() int
        +writeToCSV(path)
    }

    class Logger {
        -Mutex lock
        -File file
        +newFileLogger(path) Logger
        +close()
        +info(message)
        +warn(message)
        +error(message)
        +logWorkerStart(workerId)
        +logWorkerComplete(workerId)
        +logTaskStart(workerId, taskId)
        +logTaskComplete(workerId, taskId)
        +logTaskError(workerId, taskId, message)
        +logException(workerId, exception)
    }

    class Rider {
        -int riderId
        -string name
        -Location currentLocation
        +updateLocation(location)
        +getID() int
        +getName() string
        +getLocation() Location
    }

    class Driver {
        -int driverId
        -string name
        -Location currentLocation
        -DriverStatus status
        +acquireNearest(pickup) Driver
        +getAll() Driver[]
        +getAllAvailable() Driver[]
        +acceptRide() bool
        +completeRide()
        +completeRideAt(destination)
        +goOffline() bool
        +goOnline()
        +getID() int
        +getName() string
        +getLocation() Location
        +getStatus() DriverStatus
        +updateLocation(location)
    }

    class Location {
        -double latitude
        -double longitude
        -string address
        +distanceTo(location) double
        +equals(location) bool
        +toString() string
    }

    class TaskStatus {
        <<enumeration>>
        PENDING
        PROCESSING
        COMPLETED
        FAILED
    }

    class DriverStatus {
        <<enumeration>>
        AVAILABLE
        BUSY
        OFFLINE
    }

    class ProcessingException {
        -string message
        -int taskId
        +getMessage() string
    }

    class QueueException {
        -string message
        +getMessage() string
    }

    class FileIOException {
        -string message
        +getMessage() string
    }

    class ErrorHelpers {
        <<functions>>
        +isProcessingError(err) bool
        +isQueueError(err) bool
        +isFileIOError(err) bool
    }

    RideSharingSystem *-- TaskQueue : owns
    RideSharingSystem *-- ResultStore : owns
    RideSharingSystem *-- Logger : owns
    RideSharingSystem *-- Worker : manages

    RideSharingSystem ..> Stats : returns
    Task <|.. RideRequest : implements
    Base <|-- RideRequest : embedded in
    Base --> TaskStatus : has

    TaskQueue o-- Task : stores
    Worker --> TaskQueue : retrieves tasks
    Worker --> ResultStore : writes results
    Worker --> Logger : logs activity
    Worker --> Task : processes

    ResultStore o-- Result : stores
    Task --> Result : produces

    RideRequest --> Rider : requested by
    RideRequest --> Driver : assigned to
    RideRequest --> Location : pickup/destination

    Rider --> Location : located at
    Driver --> Location : located at
    Driver --> DriverStatus : has

    Worker ..> ProcessingException : handles
    Worker ..> QueueException : handles
    ResultStore ..> FileIOException : handles
    Logger ..> ProcessingException : logs
```

## Implementation (Go)

### Run

See [Getting started](#getting-started-no-go-installed); the command is `go run ./cmd/ridesharing`.

Outputs: console summary, `results.csv` (one row per task) and `ridesharing.log` (worker/task start, completion and errors).

### Data

`lib/simulation/Drivers.csv` and `lib/simulation/Riders.csv` are embedded in the binary (`go:embed`).
Each row of `Riders.csv` is a rider **and** its ride request: the pickup is the rider's own
`latitude/longitude/address`, the destination comes from the `destination_*` columns.
Coordinates are worldwide, so distances and fares are large; this is only demo data.

### Flow

1. `simulation.LoadDrivers()` registers every driver; `simulation.LoadRides()` registers every rider and returns one `RideRequest` per rider.
2. `system.Initialize(workers)` creates the queue, result store, logger and workers.
3. `Run` starts the workers (goroutines), enqueues all tasks, closes the queue and waits for the workers to finish.
4. Each worker dequeues a task, processes it (simulated 50-200 ms delay, reserves a free driver, computes the fare with the haversine distance), stores the `Result` and logs it.
5. Results are written to `results.csv`.

### Concurrency design

| Requirement | Where |
|---|---|
| Shared queue | `lib/queue/task_queue.go`: buffered channel; each task goes to exactly one worker |
| Worker threads | `lib/worker/worker.go`: one goroutine per worker, started with a `sync.WaitGroup` |
| Simulated work | `lib/task/ride_request.go`: `Process()` sleeps a random 50-200 ms |
| No races on shared data | `sync.Mutex` in `ResultStore`, `Logger`, the driver registry and the location list; `sync.RWMutex` on the queue's closed flag; `atomic.Bool` for the worker running flag |
| No deadlock / safe termination | queue is closed after submitting; `Dequeue` returns a `QueueError` once it is closed and drained, so every worker exits and `WaitGroup.Wait()` returns |
| No lost or duplicated tasks | a channel delivers each task once; a worker that panics still stores a failed result |
| Error handling | `lib/errors`: `ProcessingError`, `QueueError`, `FileIOError`; functions return errors, `defer` closes files and recovers panics |
| Logging | `lib/logger/logger.go`: worker start/complete, task start/complete/error, exceptions; to console and `ridesharing.log` |

### `results.csv` columns

| Column | Meaning |
|---|---|
| `task_id`, `worker_id` | the task and the worker (goroutine) that processed it |
| `rider_id`, `rider_name` | the rider; `rider_id` is the `id` in `Riders.csv` |
| `pickup_lat/lon/address` | the rider's initial location (pickup) |
| `dest_lat/lon/address` | the ride destination |
| `driver_id`, `driver_name` | the assigned driver; `driver_id` is the `id` in `Drivers.csv` (empty if none was available) |
| `driver_start_lat/lon/address` | where the driver was when assigned |
| `driver_end_lat/lon/address` | where the driver ended the ride (the destination) |
| `driver_to_pickup_km`, `trip_distance_km`, `fare` | straight-line distances (haversine) and the fare |
| `success`, `message` | outcome and a short description |
| `started_at`, `completed_at`, `duration_ms` | timing of the task |

Rows are in completion order, so they also show how the workers interleave. A driver that never appears in `driver_id` was never picked.
