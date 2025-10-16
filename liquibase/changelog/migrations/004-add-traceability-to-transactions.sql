--liquibase formatted sql

--changeset felipe-ascari:add-financial-audit-fields
ALTER TABLE transactions ADD COLUMN previous_balance decimal(15, 2);
ALTER TABLE transactions ADD COLUMN new_balance decimal(15, 2);
ALTER TABLE transactions ADD COLUMN transaction_id varchar(255) UNIQUE;
ALTER TABLE transactions ADD COLUMN idempotency_key varchar(255);
--rollback ALTER TABLE transactions DROP COLUMN previous_balance;
--rollback ALTER TABLE transactions DROP COLUMN new_balance;
--rollback ALTER TABLE transactions DROP COLUMN transaction_id;
--rollback ALTER TABLE transactions DROP COLUMN idempotency_key;

--changeset felipe-ascari:create-transactions-audit-indexes
CREATE INDEX idx_transactions_transaction_id ON transactions (transaction_id);
CREATE INDEX idx_transactions_idempotency_key ON transactions (idempotency_key);
CREATE INDEX idx_transactions_account_created ON transactions (account_id, created_at DESC);
--rollback DROP INDEX idx_transactions_transaction_id;
--rollback DROP INDEX idx_transactions_idempotency_key;
--rollback DROP INDEX idx_transactions_account_created;
