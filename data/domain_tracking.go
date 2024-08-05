package data

import (
	"context"
	"fmt"

	"github.com/NikoMalik/GoTrack/db"
	"github.com/NikoMalik/GoTrack/logEvent"
	"github.com/NikoMalik/GoTrack/util"
	"github.com/nedpals/supabase-go"
	"github.com/uptrace/bun"
)

func CountUserDomainTrackings(user *supabase.User) (int, error) {
	// #1: Counting domain trackings for a specific user.
	return db.Bun.NewSelect().Model(&DomainTracking{}).Where("user_id = ?", user.ID).Count(context.Background())
}

func GetDomainTrackings(filter map[string]any, limit, page int) ([]DomainTracking, error) {
	if limit == 0 {
		limit = defaultLimit
	}

	var trackings []DomainTracking
	builder := db.Bun.NewSelect().Model(&trackings).Limit(limit)

	// #2: Adding filter conditions to the query.
	for k, v := range filter {
		if vStr, ok := v.(string); ok && vStr != "" {
			builder.Where("? = ?", bun.Ident(k), vStr)
		}
	}

	// #3: Adding pagination offset.
	offset := (page - 1) * limit
	builder.Offset(offset)

	// #4: Executing the query.
	err := builder.Scan(context.Background(), &trackings)
	return trackings, err
}

func GetDomainTracking(filter map[string]any) (DomainTracking, error) {
	var tracking DomainTracking
	builder := db.Bun.NewSelect().Model(&tracking)

	// #5: Adding filter conditions to the query.
	for k, v := range filter {
		if vStr, ok := v.(string); ok && vStr != "" {
			builder.Where("? = ?", bun.Ident(k), vStr)
		}
	}

	// #6: Executing the query.
	err := builder.Scan(context.Background(), &tracking)
	return tracking, err
}

func DeleteDomainTracking(filter map[string]any) error {
	tracking, err := GetDomainTracking(filter)
	if err != nil {
		return err
	}

	// #7: Deleting the domain tracking.
	_, err = db.Bun.NewDelete().Model(&tracking).Where("id = ?", tracking.ID).Exec(context.Background())
	return err
}

func InsertDomain(tracking *DomainTracking) error {
	// #8: Inserting a new domain tracking.
	_, err := db.Bun.NewInsert().Model(tracking).Exec(context.Background())
	return err
}

func NewDomainTracking(trackings []*DomainTracking) error {
	tx, err := db.Bun.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			logEvent.Log("panic", fmt.Sprintf("%v", p))
		} else if err != nil {
			tx.Rollback()
			logEvent.Log("error", err.Error())
		} else {
			err = tx.Commit()
		}
	}()

	for _, tracking := range trackings {
		query := map[string]any{
			"domain_name": tracking.DomainName,
			"user_id":     tracking.User.ID,
		}
		existingTracking, err := GetDomainTracking(query)

		if err != nil {
			if util.IsErrNoRecords(err) {

				if err := InsertDomain(tracking); err != nil {
					return err
				}
			} else {
				logEvent.Log("error", err.Error())
			}
		} else if existingTracking.ID != 0 {
			// #10: Optional: Handle case where tracking already exists.
			if err := DeleteDomainTracking(query); err != nil {
				return err
			}
			// Optionally log or handle the existing tracking case
		}
	}
	return err
}

func UpdateAllTrackings(trackings []DomainTracking) error {
	_, err := db.Bun.NewUpdate().
		Model(&trackings).
		Column(
			"issuer",
			"expires",
			"signature_algo",
			"public_key_algo",
			"dns_names",
			"last_poll_at",
			"latency",
			"error",
			"status",
			"signature",
			"public_key",
			"key_usage",
			"ext_key_usages",
			"encoded_pem",
			"server_ip",
		).
		Bulk().
		Exec(context.Background())
	return err
}
