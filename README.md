# multi_thread_data_go

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
        +getRider() Rider
        +getPickup() Location
        +getDestination() Location
    }

    class Simulation {
        <<package>>
        generator.go (to implement)
    }

    class TaskQueue {
        -Queue~Task~ tasks
        -Mutex lock
        -Condition taskAvailable
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
    }

    class ResultStore {
        -List~Result~ results
        -Mutex lock
        +addResult(result)
        +getResults() List~Result~
        +count() int
        +countSuccess() int
        +countFailed() int
        +writeToFile(path)
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
        +acquireAvailable() Driver
        +getAll() Driver[]
        +acceptRide() bool
        +completeRide()
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
