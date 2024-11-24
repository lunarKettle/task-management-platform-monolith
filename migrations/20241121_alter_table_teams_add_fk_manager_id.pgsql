ALTER TABLE teams
ADD COLUMN IF NOT EXISTS manager_id INT;

ALTER TABLE teams
ADD CONSTRAINT fk_manager_id
FOREIGN KEY (manager_id)
REFERENCES users (id)
ON DELETE SET NULL;

CREATE OR REPLACE FUNCTION set_manager_id_to_null()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.manager_id = 0 THEN
        NEW.manager_id := NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_or_update_manager_id
BEFORE INSERT OR UPDATE ON teams
FOR EACH ROW
EXECUTE FUNCTION set_manager_id_to_null();