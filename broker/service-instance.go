package broker

const (
	InProgressOperationSate        OperationState       = "in progress"
	InProgressOperationDescription OperationDescription = "creating service instance..."
	SucceedOperationState          OperationState       = "succeeded"
	SucceedOperationDescription    OperationDescription = "successfully created service instance"
	FailedOperationState           OperationState       = "failed"
	FailedOperationDescription     OperationDescription = "failed to create service instance"
)

type ServiceInstance struct {
	ID               string `json:"id"`
	DashboardURL     string `json:"dashboard_url"`
	InternalID       string `json:"internalId,omitempty"`
	ServiceID        string `json:"service_id"`
	PlanID           string `json:"plan_id"`
	OrganizationGUID string `json:"organization_guid"`
	SpaceGUID        string `json:"space_guid"`

	LastOperation *LastOperation `json:"last_operation,omitempty"`

	Parameters struct {
		Parameters
	} `json:"parameters,omitempty"`

	Context Context `json:"context"`
}

type LastOperation struct {
	State                    OperationState       `json:"state"`
	Description              OperationDescription `json:"description"`
	AsyncPollIntervalSeconds int                  `json:"async_poll_interval_seconds,omitempty"`
}

type CreateServiceInstanceResponse struct {
	DashboardURL  string         `json:"dashboard_url"`
	LastOperation *LastOperation `json:"last_operation,omitempty"`
}

type GetServiceInstanceLastOperationResponse struct {
	*LastOperation
}

type OperationState string

type OperationDescription string

type Parameters struct {
	ClusterName      string `json:"cluster_name"`
	Replicas         int32  `json:"replicas,omitempty"`
	ProxySQLReplicas int32  `json:"proxy_sql_replicas,omitempty"`
	TopologyKey      string `json:"topology_key,omitempty"`
	Size             string `json:"size,omitempty"`
}

type Context struct {
	ClusterID string `json:"clusterid,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Platform  string `json:"platform,omitempty"`
}