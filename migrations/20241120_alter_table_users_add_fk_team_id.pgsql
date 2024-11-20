ALTER TABLE users
ADD COLUMN IF NOT EXISTS team_id INT;

ALTER TABLE users
ADD CONSTRAINT fk_team_id
FOREIGN KEY (team_id)
REFERENCES teams (id)
ON DELETE SET NULL;

CREATE OR REPLACE FUNCTION set_team_id_to_null()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.team_id = 0 THEN
        NEW.team_id := NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_or_update_team_id
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_team_id_to_null();