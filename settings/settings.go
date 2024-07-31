package settings

import "github.com/NikoMalik/GoTrack/data"

type accountSettings struct {
	MaxTrackings int
	Webhooks     bool
}

var Account = map[data.Plan]accountSettings{
	data.PlanStarter: {
		MaxTrackings: 5,
		Webhooks:     false,
	},
	data.PlanBusiness: {
		MaxTrackings: 10,
		Webhooks:     true,
	},
	data.PlanEnterprise: {
		MaxTrackings: 100,
		Webhooks:     true,
	},
}