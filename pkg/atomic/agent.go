package atomic

import "log"

func GetAgentId() string {
	hostId := GetHostId()
	userId, err := GetCurrentUser()
	if err != nil {
		log.Fatalf("Failed to lookup current user: %v", err)
	}
	m := map[string]interface{}{
		"host_id": hostId,
		"user_id": userId,
	}
	id, err := NewUUID5FromMap(appId, m)
	if err != nil {
		log.Fatalf("Failed to generate agent ID: %v", err)
	}
	return id
}
