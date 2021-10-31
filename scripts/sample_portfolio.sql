
-- users table and portfolio table must exist.
INSERT INTO portfolio_positions (user_id, symbol, total_quantity, equity_value, price_per_share)
    VALUES (1, 'AAPL', 20, 2500, 125);

INSERT INTO portfolio_positions (user_id, symbol, total_quantity, equity_value, price_per_share)
    VALUES (1, 'TSLA', 5.0, 3000, 600);

INSERT INTO portfolio_positions (user_id, symbol, total_quantity, equity_value, price_per_share)
     VALUES (1, 'AMZN', 1.38461538, 4500, 150);