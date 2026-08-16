DROP INDEX IF EXISTS ux_provider_transactions_provider_reference;
DROP INDEX IF EXISTS ux_provider_transactions_provider_event;
ALTER TABLE provider_transactions DROP COLUMN IF EXISTS event_id;
