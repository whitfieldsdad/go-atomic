# go-atomic

A Go client for the Atomic Red Team framework

⚠️ This project is under active development ⚠️

## Usage

### Command line interface

#### Search for and run tests by ATT&CK technique ID

```bash
go run main.go tasks --atomics-path=~/src/atomic-red-team/atomics/ list --attack-technique-id=T1082 --auto | jq
```

```json
...
{
  "id": "a8cf05d6-6e5f-5e1e-a9b1-3ed8961ba4c0",
  "template_id": "327cc050-9e99-4c8e-99b5-1d15f2fb6b96",
  "name": "T1082: Show System Integrity Protection status (MacOS)",
  "description": "Read and Display System Intergrety Protection status. csrutil is commonly used by malware and post-exploitation tools to determine whether certain files and directories on the system are writable or not.\n",
  "platforms": [
    "darwin"
  ],
  "elevation_required": false,
  "attack_technique_ids": [
    "T1082"
  ],
  "tags": [
    "T1082"
  ],
  "steps": [
    {
      "id": "8e542973-150a-4d33-9d22-6bdb5c0515d9",
      "step_type": "exec",
      "data": {
        "command_template": "csrutil status\n",
        "command_type": "sh"
      }
    }
  ]
}
...
```

To execute task `a8cf05d6-6e5f-5e1e-a9b1-3ed8961ba4c0`:`

```bash
go run main.go tasks --atomics-path=~/src/atomic-red-team/atomics/ exec --task-id a8cf05d6-6e5f-5e1e-a9b1-3ed8961ba4c0
```

```json
{
  "id": "9e07d0de-bc4f-4063-a121-7e55c629ebc8",
  "task_id": "a8cf05d6-6e5f-5e1e-a9b1-3ed8961ba4c0",
  "host_id": "475b147c-fe8d-5484-b74c-182255e488fe",
  "user_id": "f16fd24c-b117-5ad0-b1a6-41c854b07b6f",
  "steps": [
    {
      "id": "b4ec72a7-75ca-4be8-a16c-94a3fd0c1591",
      "step_id": "73fc5909-50bd-458f-a1c8-d6513560f64c",
      "step_type": "exec",
      "start_time": "2024-01-30T11:09:51.767313-05:00",
      "end_time": "2024-01-30T11:09:51.780706-05:00",
      "duration": 0.013393417,
      "data": {
        "subprocess": {
          "id": "d57ba2af-1193-425a-aa91-b1555e730160",
          "time": "2024-01-30T11:09:51.780652-05:00",
          "pid": 81773,
          "ppid": 81772,
          "name": "sh",
          "cwd": "/Users/tyler/go/src/github.com/go-atomic",
          "create_time": "2024-01-30T11:09:51.767-05:00",
          "exit_code": 0,
          "stdout": "System Integrity Protection status: enabled.\n",
          "user": {
            "id": "f16fd24c-b117-5ad0-b1a6-41c854b07b6f",
            "user_id": "501",
            "name": "Tyler Fisher",
            "username": "tyler",
            "primary_group_id": "20",
            "group_ids": [
              "20",
              "12",
              "61",
              "79",
              "80",
              "81",
              "98",
              "33",
              "100",
              "204",
              "250",
              "395",
              "398",
              "399",
              "400",
              "701"
            ],
            "home_dir": "/Users/tyler"
          },
          "user_ids": [
            501
          ],
          "executable": {
            "id": "2ceeaf36-41c0-47dc-b99e-0eefbc98f7ad",
            "time": "2024-01-30T11:09:51.768851-05:00",
            "path": "/bin/sh",
            "filename": "sh",
            "hashes": {
              "md5": "68a37d17986d5af3dc693748d56e9248",
              "sha1": "f001efb6072783430686cff41d07a6c5d4e4972b",
              "sha256": "192b7e21b34ce0de5abaf684347a0ed304a7819e74cd1017ec3a83d00f9969c6"
            }
          }
        }
      }
    }
  ],
  "start_time": "2024-01-30T11:09:51.767313-05:00",
  "end_time": "2024-01-30T11:09:51.780708-05:00",
  "duration": 0.013395166,
  "process": {
    "id": "6a4097a4-e00c-4fa8-bb5e-87ea7a367b02",
    "time": "2024-01-30T11:09:51.798808-05:00",
    "pid": 81772,
    "ppid": 81744,
    "name": "main",
    "cwd": "/Users/tyler/go/src/github.com/go-atomic",
    "argv": [
      "/var/folders/px/d9nhhv0144g2n9yjs36k6wgw0000gn/T/go-build2849268053/b001/exe/main",
      "tasks",
      "--atomics-path=~/src/atomic-red-team/atomics/",
      "exec",
      "--task-id",
      "a8cf05d6-6e5f-5e1e-a9b1-3ed8961ba4c0"
    ],
    "argc": 6,
    "create_time": "2024-01-30T11:09:51.496-05:00",
    "user": {
      "id": "f16fd24c-b117-5ad0-b1a6-41c854b07b6f",
      "user_id": "501",
      "name": "Tyler Fisher",
      "username": "tyler",
      "primary_group_id": "20",
      "group_ids": [
        "20",
        "12",
        "61",
        "79",
        "80",
        "81",
        "98",
        "33",
        "100",
        "204",
        "250",
        "395",
        "398",
        "399",
        "400",
        "701"
      ],
      "home_dir": "/Users/tyler"
    },
    "user_ids": [
      501
    ],
    "executable": {
      "id": "ac3c20c3-a03c-47b3-9671-b679ac5ab1db",
      "time": "2024-01-30T11:09:51.780994-05:00",
      "path": "/private/var/folders/px/d9nhhv0144g2n9yjs36k6wgw0000gn/T/go-build2849268053/b001/exe/main",
      "filename": "main",
      "hashes": {
        "md5": "9e309a71c9ac3c309b50d5c0ac35ee44",
        "sha1": "ef05419998a4ab5444b75bfbb83ea715ea18141d",
        "sha256": "0c143f43f7b44e6bf18f87e0a30d464cbd6ad31f2a5184e3e07c0c2346fed493"
      }
    }
  }
}
```

##### Search for and run tests by ATT&CK tactic ID

To search for tests that can be run on the current system by the current user and are related to specific ATT&CK tactics (e.g. [TA0007](https://attack.mitre.org/tactics/TA0007/)):

```bash
go run main.go tasks --atomics-path=~/src/atomic-red-team/atomics/ --attack-tactics-to-techniques-path=data/attack-tactics-to-techniques.json count --attack-tactic-id TA0007 --auto
```

```text
2024/01/30 11:11:23 INFO Selected 46 techniques related to TA0007: T1007, T1010, T1012, T1016, T1016.001, T1016.002, T1018, T1033, T1040, T1046, T1049, T1057, T1069, T1069.001, T1069.002, T1069.003, T1082, T1083, T1087, T1087.001, T1087.002, T1087.003, T1087.004, T1120, T1124, T1135, T1201, T1217, T1482, T1497, T1497.001, T1497.002, T1497.003, T1518, T1518.001, T1526, T1538, T1580, T1613, T1614, T1614.001, T1615, T1619, T1622, T1652, T1654
...
29
```

We can list the names of the discovered tests as follows using `jq`:

```bash
go run main.go tasks --atomics-path=~/src/atomic-red-team/atomics/ --attack-tactics-to-techniques-path=data/attack-tactics-to-techniques.json list --attack-tactic-id TA0007 --auto | jq -r '.display_name'
```

```text
...
T1016: System Network Configuration Discovery
T1018: Remote System Discovery - arp nix
T1018: Remote System Discovery - sweep
T1033: System Owner/User Discovery
T1046: Port Scan
T1049: System Network Connections Discovery FreeBSD, Linux & MacOS
T1057: Process Discovery - ps
T1069.001: Permission Groups Discovery (Local)
T1082: System Information Discovery
T1082: List OS Information
T1082: Hostname Discovery
T1082: Environment variables discovery on freebsd, macos and linux
T1082: Show System Integrity Protection status (MacOS)
T1083: Nix File and Directory Discovery
T1083: Nix File and Directory Discovery 2
T1087.001: View accounts with UID 0
T1087.001: List opened files by user
T1087.001: Enumerate users and groups
T1087.001: Enumerate users and groups
T1124: System Time Discovery in FreeBSD/macOS
T1135: Network Share Discovery
T1201: Examine password policy - macOS
T1217: List Mozilla Firefox Bookmark Database Files on macOS
T1217: List Google Chrome Bookmark JSON Files on macOS
T1217: List Safari Bookmarks on MacOS
T1497.001: Detect Virtualization Environment (MacOS)
T1518: Find and Display Safari Browser Version
T1518.001: Security Software Discovery - ps (macOS)
T1580: AWS - EC2 Enumeration from Cloud Instance
...
```

To run all 46 tests:

```bash
go run main.go tasks --atomics-path=~/src/atomic-red-team/atomics/ --attack-tactics-to-techniques-path=data/attack-tactics-to-techniques.json exec --attack-tactic-id TA0007 --auto > TA0007.jsonl
```
