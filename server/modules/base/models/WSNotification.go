package baseModels

type WSNotification struct {
	Action   string                 // Add, Update, Remove
	Resource string                 // Project vs file
	ResId    string                 // Id of resource
	Data     map[string]interface{} // Any other data
}
