CREATE OR REPLACE FUNCTION enforce_null_team_id()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.team_id = 0 THEN
        NEW.team_id := NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_null_on_zero_team_id
BEFORE INSERT OR UPDATE ON projects
FOR EACH ROW
EXECUTE FUNCTION enforce_null_team_id();
