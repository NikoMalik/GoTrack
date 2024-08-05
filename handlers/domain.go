package handlers

import (
	"math"
	"net/url"
	"strconv"

	"github.com/NikoMalik/GoTrack/data"
	"github.com/NikoMalik/GoTrack/sb"
	"github.com/NikoMalik/GoTrack/views/layouts"
	"github.com/gofiber/fiber/v2"
)

var limitFilters = []int{
	0,
	10,
	25,
	50,
	100,
}

var statusFilters = []string{
	"all",
	data.StatusHealthy,
	data.StatusExpires,
	data.StatusExpired,
	data.StatusInvalid,
	data.StatusOffline,
	data.StatusUnresponsive,
}

type TrackingFilter struct {
	Limit  int
	Page   int
	Status string
	Sort   string
}

func HandleDomainList(c *fiber.Ctx) error {
	user, err := sb.Client.Auth.User(c.Context(), c.Cookies("access_Token"))
	count, err := data.CountUserDomainTrackings(user)
	if err != nil {
		return err

	}

	if count == 0 {
		return Render(c, layouts.Dashboard(c))

	}

	filter, err := buildTrackingFilter(c)

	if err != nil {
		return err
	}

	filterContext := buildFilterContext(filter)
	query := map[string]any{
		"user_id": user.ID,
	}
	if filter.Status != "all" {
		query["status"] = filter.Status
	}

	domainTrackings, err := data.GetDomainTrackings(query, filter.Limit, filter.Page)

	if err != nil {
		return err
	}

	data := map[string]any{
		"domainTrackings":  domainTrackings,
		"filters":          filterContext,
		"userHasTrackings": true,
		"pages":            buildPages(count, filter.Limit),
		"queryParams":      filter.encode(),
	}

	return Render(c, layouts.Dashboard(c, data))

}

func (f *TrackingFilter) encode() string {
	values := url.Values{}
	if f.Limit != 0 {
		values.Set("limit", strconv.Itoa(f.Limit))
	}
	if f.Page != 0 {
		values.Set("page", strconv.Itoa(f.Page))
	}
	values.Set("status", f.Status)
	return values.Encode()
}

func buildTrackingFilter(c *fiber.Ctx) (*TrackingFilter, error) {
	filter := new(TrackingFilter)
	if err := c.QueryParser(filter); err != nil {
		return nil, err
	}
	if filter.Limit == 0 {
		filter.Limit = 25
	}
	return filter, nil
}

func buildFilterContext(filter *TrackingFilter) map[string]any {
	return map[string]any{
		"statuses":       statusFilters,
		"limits":         limitFilters,
		"selectedStatus": filter.Status,
		"selectedLimit":  filter.Limit,
		"selectedPage":   filter.Page,
	}
}

func buildPages(results int, limit int) []int {
	lenPages := float64(results) / float64(limit)
	pages := make([]int, int(math.RoundToEven(lenPages)))
	for i := 0; i < len(pages); i++ {
		pages[i] = i + 1
	}
	return pages
}
