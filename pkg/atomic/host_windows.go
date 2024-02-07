package atomic

func getOSVersion() string {
	key, _ := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	defer key.Close()
	productName, _, _ := key.GetStringValue("ProductName")
	return productName
}
