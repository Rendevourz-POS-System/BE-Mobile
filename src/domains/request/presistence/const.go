package presistence

type (
	Type   string
	Status string
)

const (
	Adoption   Type = "adoption"
	Donation   Type = "donation"
	Rescue     Type = "rescue"
	Monitoring Type = "monitoring"
	Publish    Type = "publish"
	Surrender  Type = "surrender"
)

const (
	New       Status = "new"
	Ongoing   Status = "ongoing"
	Completed Status = "completed"
	Rejected  Status = "rejected"
	Cancelled Status = "cancelled"
	Failed    Status = "failed"
)
