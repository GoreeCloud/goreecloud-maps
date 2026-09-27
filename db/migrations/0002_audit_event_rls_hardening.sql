BEGIN;

DROP POLICY IF EXISTS audit_events_insert_actor ON maps.audit_events;

CREATE POLICY audit_events_insert_authorized_actor
ON maps.audit_events
FOR INSERT
WITH CHECK (
    actor_user_id = maps.current_user_id()
    AND (
        collection_id IS NULL
        OR maps.collection_access_role(collection_id) IS NOT NULL
    )
);

COMMIT;
