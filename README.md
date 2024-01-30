# go-atomic

A Go client for the Atomic Red Team framework

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

To execute task `a8cf05d6-6e5f-5e1e-a9b1-3ed8961ba4c0` and write the output to `a8cf05d6-6e5f-5e1e-a9b1-3ed8961ba4c0.json`:

```bash
go run main.go tasks --atomics-path=~/src/atomic-red-team/atomics/ exec --task-id a8cf05d6-6e5f-5e1e-a9b1-3ed8961ba4c0 | jq '.' > 'a8cf05d6-6e5f-5e1e-a9b1-3ed8961ba4c0.json' 
```

```json

```

##### Search for and run tests by ATT&CK ID

To search for tests that can be run on the current system by the current user and are related to specific ATT&CK tactics (e.g. TA0007, TA0008):

```bash
go run main.go tasks --atomics-path=~/src/atomic-red-team/atomics/ --attack-tactics-to-techniques-path=data/attack-tactics-to-techniques.json count --attack-tactic-id TA0007 --attack-tactic-id TA0008 --auto
```

```text
29
```

```bash
go run main.go tasks --atomics-path=~/src/atomic-red-team/atomics/ --attack-tactics-to-techniques-path=data/attack-tactics-to-techniques.json list --attack-tactic-id TA0006 --attack-tactic-id TA0007 --
attack-tactic-id TA0009 --auto | jq -r '.name'
```

```text
...
2024/01/19 11:01:21 DEBUG Selected 64 techniques related to TA0006: T1003, T1003.001, T1003.002, T1003.003, T1003.004, T1003.005, T1003.006, T1003.007, T1003.008, T1040, T1056, T1056.001, T1056.002, T1056.003, T1056.004, T1110, T1110.001, T1110.002, T1110.003, T1110.004, T1111, T1187, T1212, T1528, T1539, T1552, T1552.001, T1552.002, T1552.003, T1552.004, T1552.005, T1552.006, T1552.007, T1552.008, T1555, T1555.001, T1555.002, T1555.003, T1555.004, T1555.005, T1555.006, T1556, T1556.001, T1556.002, T1556.003, T1556.004, T1556.005, T1556.006, T1556.007, T1556.008, T1557, T1557.001, T1557.002, T1557.003, T1558, T1558.001, T1558.002, T1558.003, T1558.004, T1606, T1606.001, T1606.002, T1621, T1649
2024/01/19 11:01:21 DEBUG Selected 46 techniques related to TA0007: T1007, T1010, T1012, T1016, T1016.001, T1016.002, T1018, T1033, T1040, T1046, T1049, T1057, T1069, T1069.001, T1069.002, T1069.003, T1082, T1083, T1087, T1087.001, T1087.002, T1087.003, T1087.004, T1120, T1124, T1135, T1201, T1217, T1482, T1497, T1497.001, T1497.002, T1497.003, T1518, T1518.001, T1526, T1538, T1580, T1613, T1614, T1614.001, T1615, T1619, T1622, T1652, T1654
2024/01/19 11:01:21 DEBUG Selected 37 techniques related to TA0009: T1005, T1025, T1039, T1056, T1056.001, T1056.002, T1056.003, T1056.004, T1074, T1074.001, T1074.002, T1113, T1114, T1114.001, T1114.002, T1114.003, T1115, T1119, T1123, T1125, T1185, T1213, T1213.001, T1213.002, T1213.003, T1530, T1557, T1557.001, T1557.002, T1557.003, T1560, T1560.001, T1560.002, T1560.003, T1602, T1602.001, T1602.002
...
T1016: System Network Configuration Discovery
T1018: Remote System Discovery - arp nix
T1018: Remote System Discovery - sweep
T1033: System Owner/User Discovery
T1046: Port Scan
T1049: System Network Connections Discovery FreeBSD, Linux & MacOS
T1056.001: MacOS Swift Keylogger
T1056.002: AppleScript - Prompt User for Password
T1056.002: AppleScript - Spoofing a credential prompt using osascript
T1057: Process Discovery - ps
T1069.001: Permission Groups Discovery (Local)
T1074.001: Stage data from Discovery.sh
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
T1110.004: SSH Credential Stuffing From MacOS
T1113: Screencapture
T1113: Screencapture (silent)
T1115: Execute commands from clipboard
T1123: using Quicktime Player
T1124: System Time Discovery in FreeBSD/macOS
T1135: Network Share Discovery
T1201: Examine password policy - macOS
T1217: List Mozilla Firefox Bookmark Database Files on macOS
T1217: List Google Chrome Bookmark JSON Files on macOS
T1217: List Safari Bookmarks on MacOS
T1497.001: Detect Virtualization Environment (MacOS)
T1518: Find and Display Safari Browser Version
T1518.001: Security Software Discovery - ps (macOS)
T1539: Steal Chrome Cookies via Remote Debugging (Mac)
T1552: AWS - Retrieve EC2 Password Data using stratus
T1552.001: Find AWS credentials
T1552.001: Extract passwords with grep
T1552.001: Find and Access Github Credentials
T1552.003: Search Through Bash History
T1552.004: Discover Private SSH Keys
T1552.004: Copy Private SSH Keys with rsync
T1552.004: Copy the users GnuPG directory with rsync
T1555.001: Export Certificate Item(s)
T1555.001: Import Certificate Item(s) into Keychain
T1555.003: Search macOS Safari Cookies
T1555.003: Simulating Access to Chrome Login Data - MacOS
T1560.001: Data Compressed - nix - zip
T1560.001: Data Compressed - nix - gzip Single File
T1560.001: Data Compressed - nix - tar Folder or File
T1560.001: Data Encrypted with zip and gpg symmetric
T1560.001: Encrypts collected data with AES-256 and Base64
...
```
