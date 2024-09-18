package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"decentrala.org/events/internal/types"
)

func (model *Model) GetEvents(city types.CityCode, organization types.OrganizationCode, from, to *time.Time) ([]types.Event, error) {
	events := []types.Event{}

	var cityClause, orgClause, fromClause, toClause string
	var params []any

	paramNo := 1

	if city != "" {
		cityClause = fmt.Sprintf("AND event.city_code=$%d", paramNo)
		paramNo += 1
		params = append(params, string(city))
	}

	if organization != "" {
		orgClause = fmt.Sprintf("AND event.organization=$%d", paramNo)
		paramNo += 1
		params = append(params, string(organization))
	}

	if from != nil {
		fromClause = fmt.Sprintf("AND event.starts_at>=$%d", paramNo)
		paramNo += 1
		params = append(params, *from)
	}

	if to != nil {
		toClause = fmt.Sprintf("AND event.ends_at>=$%d", paramNo)
		paramNo += 1
		params = append(params, *to)
	}

	query := `
		SELECT
			event.id,
			event.organization,
			event.event_name,
			event.starts_at,
			event.ends_at,
			event.city_code,
			event.event_type,
			event.canceled,
			organization.org_name
		FROM event
		LEFT JOIN organization ON event.organization = organization.org_code
		WHERE event.deleted_at IS NULL ` +
		strings.Join([]string{cityClause, orgClause, fromClause, toClause}, " ") +
		" ORDER BY event.starts_at"

	rows, err := model.db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event types.Event
		err := rows.Scan(
			&event.Id,
			&event.OrganizationCode,
			&event.Name,
			&event.StartsAt,
			&event.EndsAt,
			&event.CityCode,
			&event.EventType,
			&event.Canceled,
			&event.OrganizationName,
		)

		if err != nil {
			return nil, err
		}

		event.CityName = types.GetCityName(event.CityCode)
		event.EventTypeName = types.GetEventTypeName(event.EventType)
		event.FormatDates()

		events = append(events, event)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (model *Model) GetEvent(id types.EventId) (types.Event, error) {
	event := types.Event{}

	row := model.db.QueryRow(
		`SELECT
			event.id,
			event.organization,
			event.event_name,
			event.event_description,
			event.website,
			event.starts_at,
			event.ends_at,
			event.city_code,
			event.event_address,
			event.entry_price,
			event.event_language,
			event.event_type,
			event.canceled,
			event.canceled_at,
			event.created_at,
			event.modified_at,
			event.deleted_at,
			organization.org_name
		FROM event
		LEFT JOIN organization ON event.organization = organization.org_code
		WHERE event.id = $1`, id,
	)

	err := row.Scan(
		&event.Id,
		&event.OrganizationCode,
		&event.Name,
		&event.Description,
		&event.Website,
		&event.StartsAt,
		&event.EndsAt,
		&event.CityCode,
		&event.Address,
		&event.EntryPrice,
		&event.Language,
		&event.EventType,
		&event.Canceled,
		&event.CanceledAt,
		&event.CreatedAt,
		&event.ModifiedAt,
		&event.DeletedAt,
		&event.OrganizationName,
	)

	if err != nil {
		return types.Event{}, err
	}

	event.CityName = types.GetCityName(event.CityCode)
	event.EventTypeName = types.GetEventTypeName(event.EventType)
	event.FormatDates()

	return event, nil
}

func (model *Model) CreateEvent(event types.Event) (types.EventId, error) {
	var id int
	err := model.db.QueryRow(
		`INSERT INTO event (
			organization,
			event_name,
			event_description,
			website,
			starts_at,
			ends_at,
			city_code,
			event_address,
			entry_price,
			event_language,
			event_type,
			canceled,
			canceled_at,
			created_at )
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW())
		RETURNING id`,
		event.OrganizationCode,
		event.Name,
		event.Description,
		event.Website,
		event.StartsAt,
		event.EndsAt,
		event.CityCode,
		event.Address,
		event.EntryPrice,
		event.Language,
		event.EventType,
		event.Canceled,
		event.CanceledAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return types.EventId(id), err
}

func (model *Model) DeleteEvent(id string) error {
	_, err := model.db.Exec(
		`UPDATE event
			SET deleted_at = NOW()
		WHERE
			id = $1`, id,
	)

	return err
}

func (model *Model) ChangeEvent(id types.EventId, event types.Event) error {
	return nil
}

func (model *Model) GetOrganizationNext10Events(organization types.OrganizationCode) ([]types.Event, error) {
	events := []types.Event{}

	if organization == "" {
		return nil, errors.New("organization must be set")
	}

	query := `
		SELECT
			event.id,
			event.organization,
			event.event_name,
			event.starts_at,
			event.ends_at,
			event.city_code,
			event.event_type,
			event.canceled,
			organization.org_name
		FROM event
		LEFT JOIN organization ON event.organization = organization.org_code
		WHERE
			deleted_at IS NULL
			AND event.organization = $1
			AND starts_at >= $2
		ORDER BY starts_at
		LIMIT 10`

	rows, err := model.db.Query(query, organization, time.Now().Add(-2*time.Hour))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event types.Event
		err := rows.Scan(
			&event.Id,
			&event.OrganizationCode,
			&event.Name,
			&event.StartsAt,
			&event.EndsAt,
			&event.CityCode,
			&event.EventType,
			&event.Canceled,
			&event.OrganizationName,
		)

		if err != nil {
			return nil, err
		}

		event.CityName = types.GetCityName(event.CityCode)
		event.EventTypeName = types.GetEventTypeName(event.EventType)
		event.FormatDates()

		events = append(events, event)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return events, nil
}
