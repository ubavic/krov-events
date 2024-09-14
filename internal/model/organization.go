package model

import (
	"decentrala.org/events/internal/types"
	"github.com/google/uuid"
)

func (model *Model) GetOrganizations(city types.CityCode) ([]types.Organization, error) {
	organizations := []types.Organization{}

	parameters := []any{}
	whereClause := ""
	if city != "" {
		whereClause = " WHERE city_code=$1"
		parameters = append(parameters, city)
	}

	query := `
		SELECT
			org_code,
			org_name,
			org_description,
			website,
			email,
			org_address,
			city_code
		FROM organization
	` + whereClause + ` ORDER BY org_code`

	rows, err := model.db.Query(query, parameters...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var organization types.Organization
		err := rows.Scan(
			&organization.Code,
			&organization.Name,
			&organization.Description,
			&organization.Website,
			&organization.Email,
			&organization.Address,
			&organization.CityCode,
		)

		if err != nil {
			return nil, err
		}

		organization.CityName = types.GetCityName(organization.CityCode)

		organizations = append(organizations, organization)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func (model *Model) GetOrganization(code types.OrganizationCode) (types.Organization, error) {
	organization := types.Organization{}

	row := model.db.QueryRow(
		`SELECT
			org_code,
			org_name,
			org_description,
			website,
			email,
			token,
			api_allowed,
			is_admin,
			org_address,
			city_code,
			osm_url,
			created_at,
			modified_at
		FROM organization
		WHERE org_code = $1`,
		code,
	)

	var token string
	var admin, apiAllowed bool

	err := row.Scan(
		&organization.Code,
		&organization.Name,
		&organization.Description,
		&organization.Website,
		&organization.Email,
		&token,
		&apiAllowed,
		&admin,
		&organization.Address,
		&organization.CityCode,
		&organization.OsmUrl,
		&organization.CreatedAt,
		&organization.ModifiedAt,
	)

	organization.SetToken(token)
	organization.SetApiAllowed(apiAllowed)
	organization.SetAdmin(admin)

	if err != nil {
		return types.Organization{}, err
	}

	organization.CityName = types.GetCityName(organization.CityCode)

	return organization, nil
}

func (model *Model) CreateOrganization(organization types.Organization) error {
	token := uuid.New()

	_, err := model.db.Exec(
		`INSERT INTO
			organization (
				org_code,
				org_name,
				org_description,
				website,
				email,
				token,
				api_allowed,
				org_address,
				city_code,
				osm_url,
				created_at )
			VALUES
				($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())`,
		organization.Code,
		organization.Name,
		organization.Description,
		organization.Website,
		organization.Email,
		token,
		organization.IsApiAllowed(),
		organization.Address,
		string(organization.CityCode),
		organization.OsmUrl,
	)

	if err != nil {
		return err
	}

	return nil
}

func (model *Model) ModifyOrganization(organization types.Organization) error {
	_, err := model.db.Exec(
		`UPDATE organization 
			SET
				org_name = $1,
				org_description = $2,
				website = $3,
				email = $4,
				api_allowed = $5,
				org_address = $6,
				city_code = $7,
				osm_url = $8,
				modified_at = NOW()
			WHERE
				org_code = $9`,
		organization.Name,
		organization.Description,
		organization.Website,
		organization.Email,
		organization.IsApiAllowed(),
		organization.Address,
		string(organization.CityCode),
		organization.OsmUrl,
		organization.Code,
	)

	if err != nil {
		return err
	}

	return nil
}
