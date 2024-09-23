CREATE OR REPLACE FUNCTION prevent_created_at_update() 
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.created_at <> OLD.created_at THEN
        RAISE EXCEPTION 'Updating the value for created_at is not allowed';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION prevent_created_by_id_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.created_by_id IS DISTINCT FROM NEW.created_by_id THEN
        RAISE EXCEPTION 'created_by_id cannot be changed from % to %', OLD.created_by_id, NEW.created_by_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION set_timestamps() 
RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        NEW.created_at := CURRENT_TIMESTAMP;
        NEW.updated_at := CURRENT_TIMESTAMP;
        RETURN NEW;
    END IF;

    IF (TG_OP = 'UPDATE') THEN
        IF (NEW.updated_at < OLD.updated_at) THEN
            RAISE EXCEPTION 'updated_at cannot be before the current updated_at value';
        END IF;
        NEW.updated_at := CURRENT_TIMESTAMP;
        RETURN NEW;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION prevent_user_deletion() 
RETURNS TRIGGER AS $$
BEGIN
    -- Instead of deleting, mark the user as inactive
    UPDATE users
    SET is_active = FALSE
    WHERE id = OLD.id;

    -- Log that the user was marked as inactive
    RAISE NOTICE 'User % marked as inactive instead of being deleted', OLD.id;

    -- Prevent the actual deletion
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION enforce_banned_user_constraints()
RETURNS TRIGGER AS $$
BEGIN
    -- If the user is banned, make sure is_active and is_trusted are set to FALSE
    IF NEW.is_banned = TRUE THEN
        NEW.is_active := FALSE;
        NEW.is_trusted := FALSE;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION ensure_superuser_is_admin_staff()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_superuser THEN
        NEW.is_admin = TRUE;
        NEW.is_staff = TRUE;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION ensure_admin_is_staff()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_admin THEN
        NEW.is_staff = TRUE;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION update_product_rating() 
RETURNS TRIGGER AS $$
DECLARE
    review_stats RECORD;
BEGIN
    -- Fetch number of reviews and average rating in one query
    SELECT COUNT(*) AS num_reviews, COALESCE(AVG(rating), 0) as avg_rating
    INTO review_stats
    FROM product_reviews
    WHERE product_id = NEW.product_id;

    -- Update the product's number_of_reviews and average_rating
    UPDATE products
    SET 
        number_of_reviews = review_stats.num_reviews,
        average_rating = review_stats.avg_rating
    WHERE id = NEW.product_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION check_user_bought_product()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the user has purchased the product
    IF NOT EXISTS (
        SELECT 1
        FROM user_bought_products
        WHERE user_id = NEW.user_id
          AND product_id = NEW.product_id
    ) THEN
        RAISE EXCEPTION 'User % has not bought the product %', NEW.user_id, NEW.product_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION prevent_user_id_change()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.user_id IS DISTINCT FROM NEW.user_id THEN
        RAISE EXCEPTION 'user_id cannot be changed from % to %', OLD.user_id, NEW.user_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION validate_approved_by()
RETURNS TRIGGER AS $$
BEGIN
    -- Ensure that approved_by_id exists and is a catalog manager
    IF NEW.is_approved THEN
        IF NOT EXISTS (
            SELECT 1
            FROM catalog_managers
            WHERE user_id = NEW.approved_by_id
        ) THEN
            RAISE EXCEPTION 'approved_by_id % must be a valid catalog manager', NEW.approved_by_id;
        END IF;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION ensure_approved_by_set()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.is_approved AND NEW.approved_by_id IS NULL THEN
        RAISE EXCEPTION 'approved_by_id must be set when is_approved is TRUE';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
