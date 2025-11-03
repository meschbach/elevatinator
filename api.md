## Web Service API

The web service in `cmd/webservice` listens on `http://localhost:8999` by default and exposes a small REST/WS surface so external clients can enumerate scenarios, configure custom simulations, and drive a session tick by tick. Unless otherwise noted, responses use JSON and failures return `400` for malformed inputs, `422` for validation errors, and `404` when the referenced resource is missing.

### Basics

- `GET /` — returns `webservice running` so you can sanity check the server.
- `GET /healthz` — returns `ok` with `200` and is intended for container probes.

### Typical client flow

The reference CLI in `cmd/webservice/client` walks through the API in the following order (see `client/main.go`):

1. `GET /scenarios` — pick a scenario by its `Name` (e.g. `single-up`).
2. `GET /controllers` — pick a controller by its `Name` (currently `queue`).
3. `POST /session` — pass both names in the body to start a run and capture the returned `sessionID`.
4. Loop `POST /session/{sessionID}/tick` until the reply’s `completed` flag becomes `true`.
5. `GET /session/{sessionID}/events` — dump the full simulator log for insight into elevator activity.

The CLI defaults to `http://localhost:8999` but exposes a `-baseURL` flag, so clients should treat the base URL as configurable and avoid hard-coding it elsewhere.

### Scenario Catalog

- `GET /scenarios` — lists built-in and dynamically created scenarios:
  ```json
  {
    "available": [
      { "Name": "single-up", "Description": "a single person to go up" }
    ]
  }
  ```
  The capitalized property names reflect the Go struct field names.
- `GET /scenario` — lists only dynamically created scenarios:
  ```json
  {
    "scenarios": [
      { "id": "c0c7f5...", "name": "demo", "description": "custom test" }
    ]
  }
  ```
- `POST /scenario` — creates a placeholder dynamic scenario. Body:
  ```json
  { "name": "demo", "description": "custom test" }
  ```
  Returns `{ "id": "<uuid>" }`. Newly created scenarios have zero floors/elevators until a full definition is supplied.
- `PUT /scenario/{id}` — replaces the full scenario definition. Body must be a `DynamicScenarioWire`:
  ```json
  {
    "id": "<uuid>",
    "name": "demo",
    "description": "custom test",
    "floors": 10,
    "elevators": 2,
    "actors": [
      {
        "name": "Alice",
        "starting-floor": 0,
        "starting-tick": 0,
        "goal-floor": 9
      }
    ]
  }
  ```
  `floors` and `elevators` must be ≥ 1 and the `id` in the payload must already exist. Success returns `202 Accepted`.
- `DELETE /scenario/{id}` — removes the dynamic scenario and replies with `202 Accepted` (no body).

### Controllers

- `GET /controllers` — enumerates available controllers. Response:
  ```json
  {
    "available": [{ "Name": "queue" }]
  }
  ```

### Sessions

- `POST /session` — creates a new session by referencing a scenario/controller pair:
  ```json
  {
    "scenario": "single-up",
    "controller": "queue"
  }
  ```
  Both properties are required and must match the exact strings returned by the discovery endpoints above (the Go client treats them as pointers and raises `422` if either is absent). Returns `{ "sessionID": "<uuid>" }`. The session stays in-memory for subsequent ticks.
- `POST /session/{sessionID}/tick` — advances the simulation one tick and returns `{ "completed": false }` until the scenario finishes.
- `GET /session/{sessionID}/events` — streams the accumulated simulator events for the session:
  ```json
  {
    "events": [
      { "eventType": "TickStart", "timestamp": 0 },
      { "eventType": "ElevatorCalled", "floor": 3 },
      { "eventType": "ElevatorArrived", "floor": 3, "elevator": 0 }
    ]
  }
  ```
  The client decodes these into Go pointer fields so each attribute is present only when the simulator emitted it (e.g. `timestamp` is a `*int64`). Possible `eventType` values include `TickStart`, `TickDone`, `InitStart`, `InitDone`, `InformElevator`, `InformFloor`, `ElevatorCalled`, `ElevatorArrived`, `ElevatorFloorRequest`, `ActorFinished`, and `ElevatorAtFloor`; unknown events fall back to the simulator’s string form.

### Real-time channel

- `GET /real-time` — upgrades to a WebSocket that currently logs every inbound message. The server advertises the `elevatinator/v1` sub-protocol and keeps the connection alive by responding to `Ping` frames; no structured payloads are defined yet.

All endpoints are unauthenticated and intended for local experimentation. Combine the scenario APIs to define problems, the controller list to understand available AIs, session APIs to drive simulations, and (optionally) the WebSocket for future real-time integrations. For a concrete example, run `go run ./cmd/webservice/client -baseURL http://localhost:8999` and mirror its request sequence in your own integrations.

## Lifecycle stage

This project is very young. Ideally the following enhancements can be made:

- gRPC to a controller
