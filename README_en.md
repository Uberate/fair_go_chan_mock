# Go Channel Queue System

This project implements a multi-tenant queue system with two different queueing strategies: FIFO (First-In-First-Out) and Fair (tenant-specific round-robin).

## System Architecture

### Core Components

- **Queue Interface** (`type.go:5-9`): Defines the common interface for all queue implementations
  - `Next() (tenant string, message int)` - Retrieves the next message
  - `Put(tenant string, message int) error` - Adds a message to the queue
  - `IsEmpty() bool` - Checks if the queue is empty

- **Message Struct** (`type.go:13-16`): Represents a message with tenant identifier and sequential message index

- **Generator Function** (`type.go:18-29`): Factory function that generates tenant-specific sequential messages

### Queue Implementations

1. **FIFO Queue** (`fifo.go`): Single-channel implementation where all messages are processed in strict arrival order
2. **Fair Queue** (`fair.go`): Multi-channel implementation with round-robin scheduling per tenant

## Output Format Specification

### Input Pattern Display

The system displays input patterns using the following notation:

- **`+`** (Plus Sign): Indicates a write operation (message input) for a specific tenant in that step
- **`.`** (Dot): Indicates no operation (no message input) for a specific tenant in that step

**Example Input Pattern:**
```
t1: ++++++
t2: +...+.
t3: ..+..+
```

This means:
- Tenant `t1` had 6 consecutive write operations
- Tenant `t2` had writes in step 1, step 5 (with gaps in steps 2-4)
- Tenant `t3` had writes in step 3 and step 6

### Next Call Pattern Display

The system displays Next() call patterns using similar notation:

- **`+`** (Plus Sign): Indicates a read operation (message consumption) for a specific tenant in that step
- **`.`** (Dot): Indicates no read operation for a specific tenant in that step

**Example Next Pattern:**
```
t1: ++.+.+
t2: .+.+..
t3: ..+..+
```

This shows how messages are consumed from the queue over time.

## Step Concept and Progression Logic

### Step Definition
A **step** represents a discrete unit of time in the simulation where:
1. **Input Phase**: All messages for the current step are added to the queue
2. **Processing Phase**: One `Next()` call is made to consume a message from the queue

### Step Progression Logic

The system processes steps in the following sequence:

1. **Step Initialization**: For each step in `exampleSteps`, all messages are added to the queue
2. **Message Consumption**: After input, exactly one `Next()` call is made to retrieve a message
3. **Pattern Tracking**: Both input and consumption patterns are tracked per tenant per step
4. **Visual Alignment**: Output is formatted to align all tenants vertically for easy comparison

### Step Characteristics

- **Variable Input Size**: Different steps can have different numbers of messages (0 to many)
- **Multiple Writes per Step**: A single tenant can write multiple messages in one step
- **Empty Steps**: Some steps may contain no messages (testing queue draining behavior)
- **Continuous Processing**: The system continues processing until all messages are consumed

## Example Execution Flow

Given the `exampleSteps` configuration:

```go
var exampleSteps = [][]Message{
    {t1, t1, t2},     // Step 0: 2 messages from t1, 1 from t2
    {t1, t1},         // Step 1: 2 messages from t1
    {},               // Step 2: No messages (empty step)
    {t1, t2, t3},     // Step 3: 1 message each from t1, t2, t3
    // ... additional steps
}
```

The system will:
1. Process Step 0: Add 3 messages, consume 1 message
2. Process Step 1: Add 2 messages, consume 1 message
3. Process Step 2: Add 0 messages, consume 1 message (if queue not empty)
4. Process Step 3: Add 3 messages, consume 1 message
5. Continue until all messages are consumed

## Queue Behavior Differences

### FIFO Queue Behavior
- All messages processed in strict arrival order
- No tenant prioritization
- Single shared channel

### Fair Queue Behavior
- Round-robin scheduling between tenants
- Each tenant gets its own channel
- Fair distribution of processing time
- Prevents tenant starvation

## Usage

Run the simulation:
```bash
go run .
```

This will execute both FIFO and Fair queue implementations with the same input pattern, allowing comparison of their output behaviors.