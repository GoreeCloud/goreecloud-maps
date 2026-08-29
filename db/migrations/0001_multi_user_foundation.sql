BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE SCHEMA IF NOT EXISTS maps;

CREATE TABLE maps.users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    identity_subject text NOT NULL UNIQUE,
    display_name text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE maps.preferences (
    user_id uuid PRIMARY KEY REFERENCES maps.users(id) ON DELETE CASCADE,
    map_style text NOT NULL DEFAULT 'standard',
    units text NOT NULL DEFAULT 'system' CHECK (units IN ('system', 'metric', 'imperial')),
    avoid_tolls boolean NOT NULL DEFAULT false,
    avoid_highways boolean NOT NULL DEFAULT false,
    avoid_ferries boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE maps.saved_places (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_user_id uuid NOT NULL REFERENCES maps.users(id) ON DELETE CASCADE,
    provider text NOT NULL,
    provider_place_id text,
    name text NOT NULL,
    address text,
    position geography(Point, 4326) NOT NULL,
    note text,
    revision bigint NOT NULL DEFAULT 1 CHECK (revision > 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (owner_user_id, provider, provider_place_id)
);

CREATE INDEX saved_places_owner_idx ON maps.saved_places(owner_user_id);
CREATE INDEX saved_places_position_gix ON maps.saved_places USING gist(position);

CREATE TABLE maps.collections (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_user_id uuid NOT NULL REFERENCES maps.users(id) ON DELETE CASCADE,
    name text NOT NULL CHECK (char_length(name) BETWEEN 1 AND 160),
    description text,
    sharing_mode text NOT NULL DEFAULT 'private' CHECK (sharing_mode IN ('private', 'members', 'link')),
    share_token_hash bytea,
    revision bigint NOT NULL DEFAULT 1 CHECK (revision > 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CHECK ((sharing_mode = 'link' AND share_token_hash IS NOT NULL) OR sharing_mode <> 'link')
);

CREATE INDEX collections_owner_idx ON maps.collections(owner_user_id);

CREATE TABLE maps.collection_members (
    collection_id uuid NOT NULL REFERENCES maps.collections(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES maps.users(id) ON DELETE CASCADE,
    role text NOT NULL CHECK (role IN ('editor', 'viewer')),
    invited_by_user_id uuid NOT NULL REFERENCES maps.users(id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (collection_id, user_id)
);

CREATE INDEX collection_members_user_idx ON maps.collection_members(user_id);

CREATE TABLE maps.collection_items (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    collection_id uuid NOT NULL REFERENCES maps.collections(id) ON DELETE CASCADE,
    created_by_user_id uuid NOT NULL REFERENCES maps.users(id) ON DELETE RESTRICT,
    provider text NOT NULL,
    provider_place_id text,
    name text NOT NULL,
    address text,
    position geography(Point, 4326) NOT NULL,
    note text,
    sort_key bigint NOT NULL DEFAULT 0,
    revision bigint NOT NULL DEFAULT 1 CHECK (revision > 0),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX collection_items_collection_idx ON maps.collection_items(collection_id, sort_key, created_at);
CREATE INDEX collection_items_position_gix ON maps.collection_items USING gist(position);

CREATE TABLE maps.audit_events (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    actor_user_id uuid REFERENCES maps.users(id) ON DELETE SET NULL,
    collection_id uuid REFERENCES maps.collections(id) ON DELETE SET NULL,
    event_type text NOT NULL,
    event_data jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX audit_events_collection_idx ON maps.audit_events(collection_id, created_at DESC);
CREATE INDEX audit_events_actor_idx ON maps.audit_events(actor_user_id, created_at DESC);

CREATE OR REPLACE FUNCTION maps.current_user_id()
RETURNS uuid
LANGUAGE sql
STABLE
AS $$
    SELECT NULLIF(current_setting('goreecloud.maps_user_id', true), '')::uuid
$$;

CREATE OR REPLACE FUNCTION maps.collection_access_role(target_collection_id uuid)
RETURNS text
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = maps, pg_temp
AS $$
    SELECT CASE
        WHEN c.owner_user_id = maps.current_user_id() THEN 'owner'
        ELSE cm.role
    END
    FROM maps.collections AS c
    LEFT JOIN maps.collection_members AS cm
        ON cm.collection_id = c.id
       AND cm.user_id = maps.current_user_id()
    WHERE c.id = target_collection_id
$$;

REVOKE ALL ON FUNCTION maps.collection_access_role(uuid) FROM PUBLIC;

ALTER TABLE maps.preferences ENABLE ROW LEVEL SECURITY;
ALTER TABLE maps.saved_places ENABLE ROW LEVEL SECURITY;
ALTER TABLE maps.collections ENABLE ROW LEVEL SECURITY;
ALTER TABLE maps.collection_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE maps.collection_items ENABLE ROW LEVEL SECURITY;
ALTER TABLE maps.audit_events ENABLE ROW LEVEL SECURITY;

CREATE POLICY preferences_owner_all
ON maps.preferences
FOR ALL
USING (user_id = maps.current_user_id())
WITH CHECK (user_id = maps.current_user_id());

CREATE POLICY saved_places_owner_all
ON maps.saved_places
FOR ALL
USING (owner_user_id = maps.current_user_id())
WITH CHECK (owner_user_id = maps.current_user_id());

CREATE POLICY collections_select_authorized
ON maps.collections
FOR SELECT
USING (
    owner_user_id = maps.current_user_id()
    OR EXISTS (
        SELECT 1
        FROM maps.collection_members AS cm
        WHERE cm.collection_id = collections.id
          AND cm.user_id = maps.current_user_id()
    )
);

CREATE POLICY collections_insert_owner
ON maps.collections
FOR INSERT
WITH CHECK (owner_user_id = maps.current_user_id());

CREATE POLICY collections_update_owner_or_editor
ON maps.collections
FOR UPDATE
USING (maps.collection_access_role(id) IN ('owner', 'editor'))
WITH CHECK (owner_user_id = (
    SELECT c.owner_user_id
    FROM maps.collections AS c
    WHERE c.id = collections.id
));

CREATE POLICY collections_delete_owner
ON maps.collections
FOR DELETE
USING (owner_user_id = maps.current_user_id());

CREATE POLICY collection_members_select_authorized
ON maps.collection_members
FOR SELECT
USING (maps.collection_access_role(collection_id) IS NOT NULL);

CREATE POLICY collection_members_insert_owner
ON maps.collection_members
FOR INSERT
WITH CHECK (
    maps.collection_access_role(collection_id) = 'owner'
    AND invited_by_user_id = maps.current_user_id()
    AND user_id <> maps.current_user_id()
);

CREATE POLICY collection_members_update_owner
ON maps.collection_members
FOR UPDATE
USING (maps.collection_access_role(collection_id) = 'owner')
WITH CHECK (maps.collection_access_role(collection_id) = 'owner');

CREATE POLICY collection_members_delete_owner_or_self
ON maps.collection_members
FOR DELETE
USING (
    maps.collection_access_role(collection_id) = 'owner'
    OR user_id = maps.current_user_id()
);

CREATE POLICY collection_items_select_authorized
ON maps.collection_items
FOR SELECT
USING (maps.collection_access_role(collection_id) IS NOT NULL);

CREATE POLICY collection_items_insert_editor
ON maps.collection_items
FOR INSERT
WITH CHECK (
    maps.collection_access_role(collection_id) IN ('owner', 'editor')
    AND created_by_user_id = maps.current_user_id()
);

CREATE POLICY collection_items_update_editor
ON maps.collection_items
FOR UPDATE
USING (maps.collection_access_role(collection_id) IN ('owner', 'editor'))
WITH CHECK (maps.collection_access_role(collection_id) IN ('owner', 'editor'));

CREATE POLICY collection_items_delete_editor
ON maps.collection_items
FOR DELETE
USING (maps.collection_access_role(collection_id) IN ('owner', 'editor'));

CREATE POLICY audit_events_select_owner
ON maps.audit_events
FOR SELECT
USING (
    actor_user_id = maps.current_user_id()
    OR (collection_id IS NOT NULL AND maps.collection_access_role(collection_id) = 'owner')
);

CREATE POLICY audit_events_insert_actor
ON maps.audit_events
FOR INSERT
WITH CHECK (actor_user_id = maps.current_user_id());

COMMIT;
