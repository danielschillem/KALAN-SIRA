ALTER TABLE provider_transactions ADD COLUMN IF NOT EXISTS event_id varchar(160);
CREATE UNIQUE INDEX IF NOT EXISTS ux_provider_transactions_provider_event ON provider_transactions(provider,event_id) WHERE event_id IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS ux_provider_transactions_provider_reference ON provider_transactions(provider,provider_reference) WHERE provider_reference IS NOT NULL;
