CREATE TABLE organization (
    org_code varchar(20) PRIMARY KEY,
    org_name varchar(100) NOT NULL UNIQUE,
    org_description varchar(2000) NOT NULL DEFAULT '',
    website varchar(200) NOT NULL DEFAULT '',
    email varchar(50) NOT NULL UNIQUE,
    token uuid NOT NULL,
    api_allowed boolean NOT NULL DEFAULT false,
    is_admin boolean NOT NULL DEFAULT false,
    org_address varchar(100) NOT NULL DEFAULT '',
    city_code varchar(8) NOT NULL DEFAULT '',
    osm_url varchar(100) NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL,
    modified_at timestamptz,
    CONSTRAINT code_format CHECK (org_code ~ '^[a-z]+$')
);

CREATE INDEX idx_org_city ON organization (city_code);

CREATE TABLE event (
    id serial PRIMARY KEY,
    organization varchar(20) NOT NULL,
    event_name varchar(100) NOT NULL,
    event_description varchar(2000) NOT NULL DEFAULT '',
    website varchar(200) NOT NULL DEFAULT '',
    starts_at timestamptz NOT NULL,
    ends_at timestamptz,
    city_code varchar(8) NOT NULL DEFAULT 0,
    event_address varchar(100) NOT NULL DEFAULT '',
    entry_price smallint NOT NULL DEFAULT 0,
    event_language smallint NOT NULL DEFAULT 0,
    event_type smallint NOT NULL DEFAULT 0,
    canceled boolean NOT NULL DEFAULT false,
    canceled_at timestamptz,
    created_at timestamptz NOT NULL,
    modified_at timestamptz,
    deleted_at timestamptz
    CONSTRAINT positive_entry CHECK (entry_price >= 0),
    CONSTRAINT fk_organization FOREIGN KEY (organization)
    REFERENCES organization (org_code) ON DELETE CASCADE
);

CREATE INDEX idx_event_start ON event (starts_at);
CREATE INDEX idx_event_org ON event (organization, starts_at);

CREATE TABLE migration (
    id uuid PRIMARY KEY,
    migration_file varchar(10),
    migration_start timestamptz NOT NULL,
    migration_end timestamptz NOT NULL,
    success boolean NOT NULL,
    error text NOT NULL DEFAULT ''
);
