--liquibase formatted sql

--liquibase formatted sql

--changeset felipe-ascari:seed-exchanges
INSERT INTO exchanges (name, minimum_age, maximum_transfer_amount)
    VALUES ('Binance', 18, 1000.00),
           ('Coinbase', 18, 500.00),
           ('Kraken', 21, 7500.00),
           ('Bitfinex', 18, 2000.00),
           ('Huobi', 18, 800.00);
--rollback DELETE FROM exchanges WHERE name IN ('Binance', 'Coinbase', 'Kraken', 'Bitfinex', 'Huobi');


