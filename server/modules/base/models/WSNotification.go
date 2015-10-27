package baseModels
import "github.com/CodeCollaborate/CodeCollaborate/server/modules/base/requests"

type WSNotification struct {
	Type     string                 // Notification
	Action   string                 // Add, Update, Remove
	Resource string                 // Project vs file
	ResId    string                 // Id of resource
	Data     map[string]interface{} // Any other data
}

func NewNotification(baseRequest baseRequests.BaseRequest, data map[string]interface{}) *WSNotification {

	notification := new(WSNotification)
	notification.Type = "Notification"
	notification.Action = baseRequest.Action
	notification.Resource = baseRequest.Resource
	notification.ResId = baseRequest.ResId
	notification.Data = data

	return notification
}