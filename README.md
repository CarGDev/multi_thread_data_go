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
        +initialize(workerCount)
        +submitTask(task)
        +startWorkers()
        +shutdown()
        +awaitCompletion()
    }

    class Task {
        <<abstract>>
        -int taskId
        -TaskStatus status
        -DateTime createdAt
        +process() Result
        +getId() int
        +getStatus() TaskStatus
    }

    class RideRequest {
        -Rider rider
        -Location pickup
        -Location destination
        +process() Result
        +findDriver() Driver
        +calculateFare() double
    }

    class TaskQueue {
        -Queue~Task~ tasks
        -Mutex lock
        -Condition taskAvailable
        -bool closed
        +enqueue(task)
        +dequeue() Task
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
    }

    class Result {
        -int taskId
        -bool success
        -string message
        -DateTime completedAt
        +getTaskId() int
        +isSuccess() bool
    }

    class ResultStore {
        -List~Result~ results
        -Mutex lock
        +addResult(result)
        +getResults() List~Result~
        +writeToFile(path)
    }

    class Logger {
        -Mutex lock
        +info(message)
        +error(message)
        +logWorkerStart(workerId)
        +logWorkerComplete(workerId)
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
        +acceptRide() bool
        +completeRide()
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

    RideSharingSystem *-- TaskQueue : owns
    RideSharingSystem *-- ResultStore : owns
    RideSharingSystem *-- Logger : owns
    RideSharingSystem *-- Worker : manages

    Task <|-- RideRequest
    Task --> TaskStatus : has

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
