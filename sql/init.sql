CREATE TABLE accounts (account VARCHAR PRIMARY KEY);

CREATE TABLE currs (account VARCHAR REFERENCES accounts(account), curr VARCHAR);

CREATE TABLE balances (account VARCHAR REFERENCES accounts(account), balance INTEGER);

CREATE TABLE transfers (dt timestamp DEFAULT now(), account_from VARCHAR REFERENCES accounts(account), account_to VARCHAR REFERENCES accounts(account), amount INTEGER);

CREATE OR REPLACE FUNCTION curr_cmp()
RETURNS trigger AS $BODY$

DECLARE
	bal INTEGER;
BEGIN
IF EXISTS (SELECT 1 FROM currs WHERE account IN (NEW.account_from, NEW.account_to) GROUP BY curr HAVING count(*) = 2)
THEN
	UPDATE balances SET balance = balance - NEW.amount WHERE account = NEW.account_from RETURNING balance INTO bal;
	IF bal < 0
	THEN
		RAISE EXCEPTION 'Need more minerals';
	END IF;
	UPDATE balances SET balance = balance + NEW.amount WHERE account = NEW.account_to;
	RETURN NEW;
END IF;
	RAISE EXCEPTION 'Currencies mismatch';
END;
$BODY$
LANGUAGE 'plpgsql';

CREATE TRIGGER curr_cmp BEFORE INSERT ON transfers 
FOR EACH ROW EXECUTE PROCEDURE curr_cmp();