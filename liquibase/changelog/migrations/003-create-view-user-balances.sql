--liquibase formatted sql

--changeset felipe-ascari:create-user-balances-view
CREATE VIEW user_exchange_balances AS
SELECT a.user_id,
       u.username,
       e.name AS exchange_name,
       a.balance,
       a.created_at,
       a.updated_at
    FROM accounts a
         INNER JOIN exchanges e ON e.id = a.exchange_id
         INNER JOIN users u ON u.id = a.user_id
    WHERE a.balance > 0
      AND a.deleted_at IS NULL
      AND u.deleted_at IS NULL;