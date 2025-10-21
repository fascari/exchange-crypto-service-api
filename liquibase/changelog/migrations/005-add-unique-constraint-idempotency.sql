--liquibase formatted sql

--changeset felipe-ascari:add-unique-constraint-idempotency-key
CREATE UNIQUE INDEX idx_transactions_account_idempotency_unique ON transactions(account_id, idempotency_key) WHERE idempotency_key IS NOT NULL AND deleted_at IS NULL;
--rollback DROP INDEX idx_transactions_account_idempotency_unique;

