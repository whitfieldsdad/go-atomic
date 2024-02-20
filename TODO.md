# TODO

- [ ] Create UI (Python: Polars, NetworkX, FastAPI, HTMX, Tailwind CSS)
- [ ] Create agent (Go: DuckDB, Afero, NATS)
- [ ] Create C2 service (Python: Polars, NetworkX, NATS)
- [ ] Create event correlation service (Python: Polars, NetworkX, NATS)

---

## Implementation notes

- Host IDs are unique to each host and persist across reboots.
- User IDs are unique to each user, are derived from a combination of usernames and host IDs, and are persist across reboots.
- Agent IDs are ephemeral and are unique to each process running an instance of the agent.
- Tasks can be scheduled to run on a specific host, user, or agent.

---

## Message formats

- [ ] Tenants
- [ ] Agents
- [ ] Hosts
    - [ ] Processes
    - [ ] Files
    - [ ] Users
- [ ] Tasks
- [ ] Task invocation requests
- [ ] Task invocations
- [ ] Task invocation statuses
- [ ] Task invocation responses

## Tenants

```json
{
    "id": "73edb347-49a9-4979-9df5-5a322baaaa63",
    "name": "Tyler's Homelab"
}
```

### Agents

```json
{
    "id": "788a3f10-f11b-4a72-9b70-f73c9af2967d"
}
```

### Hosts

```json
{
    "id": "84e12fd9-d863-4400-b0c3-fa94bfc019d8",
    "hostname": "tylers-mac-mini.lan",
    "os": {
        "type": "darwin",
        "arch": "arm64",
        "version": "MacOS 13.5"
    }
}
```

#### Processes

```json
{
    "id": "eb837afc-b065-4d8b-9b19-43bac314b233",
    "pid": 30222,
    "ppid": 10751,
    "argv": ["/Library/Developer/CommandLineTools/Library/Frameworks/Python3.framework/Versions/3.9/Resources/Python.app/Contents/MacOS/Python"],
    "argc": 1,
    "user": {
        "id": "9f24bc57-e8ae-4d67-9562-7a16316cec3e",
        "user_id": 501,
        "username": "tyler",
        "name": "Tyler",
        "homedir": "/Users/tyler"
    },
    "executable": {
        "id": "0354f0f5-f8a9-4b08-8063-74068f4dbbd9",
        "path": "/Library/Developer/CommandLineTools/Library/Frameworks/Python3.framework/Versions/3.9/Resources/Python.app/Contents/MacOS/Python",
        "filename": "Python",
        "hashes": {
            "md5": "90b306cfe439583b92032e20c1620153",
            "sha1": "740c717798755c81a85f06b3b2a6897fdeda93b7",
            "sha256": "ff57195e2dbd1c21e617de91d22c7020c185e5547c2a4c4710171a22aec17fbf"
        }
    }
}
```

#### Files

```json
{
    "id": "0354f0f5-f8a9-4b08-8063-74068f4dbbd9",
    "path": "/Library/Developer/CommandLineTools/Library/Frameworks/Python3.framework/Versions/3.9/Resources/Python.app/Contents/MacOS/Python",
    "filename": "Python",
    "hashes": {
        "md5": "90b306cfe439583b92032e20c1620153",
        "sha1": "740c717798755c81a85f06b3b2a6897fdeda93b7",
        "sha256": "ff57195e2dbd1c21e617de91d22c7020c185e5547c2a4c4710171a22aec17fbf"
    }
}
```

#### Users

```json
{
    "id": "9f24bc57-e8ae-4d67-9562-7a16316cec3e",
    "user_id": 501,
    "username": "tyler",
    "name": "Tyler",
    "homedir": "/Users/tyler"
}
```

### Task templates

To execute a command with input arguments:

```json
{
    "id": "ab5c3d09-29da-426f-bcfa-e608cc346d70",
    "input_arguments": [
        {
            "name": "process_to_enumerate",
            "type": "string",
            "value": "lsass",
            "description": "Process name string to search for."
        }
    ],
    "steps": [
        {
            "id": "0f59b092-00a8-49e5-8cfd-bf92ab5052d7",
            "type": "execute-command",
            "data": {
                "command": "tasklist | findstr #{process_to_enumerate}",
                "command_type": "command_prompt"
            }
        }
    ]
}
```

To execute one command followed by another:

```json
{
    "id": "904dd6fb-4f6d-4fee-bb21-5e937b74c179",
    "input_arguments": [
        {
            "name": "process_to_enumerate",
            "type": "string",
            "value": "lsass",
            "description": "Process name string to search for."
        }
    ],
    "steps": [
        {
            "id": "fdd9ff2b-662e-4c51-8a25-c44a6fc025fc",
            "type": "execute-command",
            "data": {
                "command": "tasklist | findstr #{process_to_enumerate}",
                "command_type": "command_prompt"
            }
        },
        {
            "id": "749f0b1c-a57d-4120-a256-5a6e9f1db875",
            "type": "execute-command",
            "data": {
                "command": "Get-WmiObject -Class Win32_Process -Filter \"Name='#{process_to_enumerate}'\"",
                "command_type": "powershell"
            }
        }
    ],
    "flows": [
        {
            "s": "fdd9ff2b-662e-4c51-8a25-c44a6fc025fc",
            "p": "on-success",
            "o": "749f0b1c-a57d-4120-a256-5a6e9f1db875"
        }
    ]
}
```

### Task invocations

A task invocation consist of a task and a set of input arguments.

```json
{
    "id": "ff8b423b-d456-47a6-8f71-b55889e35e7e",
    "agent_id": "e9578d59-bdc0-445b-9307-7f5785bc863c",
    "host_id": "696eeb39-4567-47bb-a3e0-657e3fbc2e29",
    "user_id": "9f24bc57-e8ae-4d67-9562-7a16316cec3e",
    "task_id": "904dd6fb-4f6d-4fee-bb21-5e937b74c179",
    "input_arguments": [
        {
            "name": "process_to_enumerate",
            "type": "string",
            "value": "lsass",
            "description": "Process name string to search for."
        }
    ],
}
```


### Task invocation statuses

Task invocation statuses can have one of the following status types:

- `created`: the task invocation has been created server-side
- `queued`: the task invocation has been queued by the agent for execution
- `running`: the task invocation is being executed by the agent
- `success`: the task invocation has been successfully executed by the agent
- `failed`: the agent failed to process the task invocation request (e.g. because the agent could not understand the request)

For example, the following message indicates that a task invocation has been created:

```json
{
    "id": "f5458ad1-8f3d-48ae-b5f4-1b6179768d83",
    "time": "2024-02-19T18:01:41.904487",
    "task_id": "904dd6fb-4f6d-4fee-bb21-5e937b74c179",
    "task_invocation_id": "96016224-772f-4104-a42f-c5ddda251ce2",
    "host_id": "696eeb39-4567-47bb-a3e0-657e3fbc2e29",
    "user_id": "9f24bc57-e8ae-4d67-9562-7a16316cec3e",
    "agent_id": "e9578d59-bdc0-445b-9307-7f5785bc863c",
    "status": "created"
}
```

### Task invocation results

```json
{
    "id": "85f08a59-8f75-4d01-9d91-e5f7bd3a7340",
    "time": "2024-02-19T18:05:39.721319",
    "task_id": "904dd6fb-4f6d-4fee-bb21-5e937b74c179",
    "task_invocation_id": "96016224-772f-4104-a42f-c5ddda251ce2",
    "host_id": "696eeb39-4567-47bb-a3e0-657e3fbc2e29",
    "user_id": "9f24bc57-e8ae-4d67-9562-7a16316cec3e",
    "agent_id": "e9578d59-bdc0-445b-9307-7f5785bc863c",
    "process": {
        "id": "b5d44075-9274-44e0-ae8d-7e440d62422a",
        "pid": 123,
        "ppid": 456
    },
    "user": {
        "id": "9f24bc57-e8ae-4d67-9562-7a16316cec3e",
        "username": "tyler"
    },
    "steps": [
        {
            "id": "66427cff-577d-4188-b596-8c5ca8f83ad2",
            "step_id": "fdd9ff2b-662e-4c51-8a25-c44a6fc025fc",
            "data": {
                "command": "tasklist | findstr lsass",
                "command_type": "command_prompt",
                "process": {
                    "id": "e6fe51fb-88b6-4259-aebe-94f6a60267f7",
                    "pid": 125,
                    "ppid": 123
                }
            }
        },
        {
            "id": "68d1ea1f-d468-4f02-9b97-6dd06de24d73",
            "step_id": "749f0b1c-a57d-4120-a256-5a6e9f1db875",
            "data": {
                "command": "Get-WmiObject -Class Win32_Process -Filter \"Name='#{process_to_enumerate}'\"",
                "command_type": "powershell",
                "process": {
                    "id": "91a00432-00e1-4862-90e6-f3854b33d37b",
                    "pid": 135,
                    "ppid": 123
                }
            }
        }
    ]
}
```

---

## Agent

### Tasks

- [ ] Read tasks from a directory
- [ ] Write tasks to a directory

### Task invocation requests

- [ ] Read task invocation requests from a directory
- [ ] Write task invocation requests to a directory

### Task invocations

- [ ] Read task invocations from a directory
- [ ] Write task invocations to a directory

### Task invocation statuses

- [ ] Read task invocation statuses from a directory
- [ ] Write task invocation statuses to a directory

### Task invocation results

- [ ] Read task invocation results from a directory
- [ ] Write task invocation results to a directory
