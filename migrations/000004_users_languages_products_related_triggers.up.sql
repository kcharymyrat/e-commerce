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


-- users table triggers
CREATE TRIGGER users_prevent_created_at_update
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER users_set_timestamps
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER users_prevent_delete
BEFORE DELETE ON users
FOR EACH ROW
EXECUTE FUNCTION prevent_user_deletion();

CREATE TRIGGER enforce_banned_constraints_trigger
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION enforce_banned_user_constraints();

CREATE TRIGGER superuser_is_admin_and_is_staff
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW 
EXECUTE FUNCTION ensure_superuser_is_admin_staff();

CREATE TRIGGER admin_is_staff
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION ensure_admin_is_staff();


-- user_referrals table triggers
CREATE TRIGGER user_referrals_prevent_created_at_update
BEFORE UPDATE ON user_referrals
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER user_refferals_set_timeststamps
BEFORE INSERT OR UPDATE ON user_referrals
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER user_referrals_prevent_user_id_change
BEFORE UPDATE ON user_referrals
FOR EACH ROW
EXECUTE FUNCTION prevent_user_id_change();


-- user_product_referrals table triggers
CREATE TRIGGER user_prod_refs_prevent_created_at_update
BEFORE UPDATE ON user_product_referrals
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER user_prod_refs_set_timestamps
BEFORE INSERT OR UPDATE ON user_product_referrals
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER user_prod_refs_prevent_user_id_change
BEFORE UPDATE ON user_product_referrals
FOR EACH ROW
EXECUTE FUNCTION prevent_user_id_change();


-- user_bough_products table triggers
CREATE TRIGGER user_bght_prods_prevent_created_at_update
BEFORE UPDATE ON user_bought_products
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER user_bght_prods_set_timestamps
BEFORE INSERT OR UPDATE ON user_bought_products
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();


-- languages table triggers
CREATE TRIGGER langs_set_timestamps 
BEFORE INSERT OR UPDATE ON languages
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER langs_prevent_created_at_update
BEFORE UPDATE ON languages
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

-- translations table triggers
CREATE TRIGGER translations_set_timestamps
BEFORE INSERT OR UPDATE ON translations
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER translations_prevent_created_at_update
BEFORE UPDATE ON translations
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER user_referrals_prevent_created_by_id_change
BEFORE UPDATE ON translations
FOR EACH ROW
EXECUTE FUNCTION prevent_created_by_id_change();


-- brands table triggers
CREATE TRIGGER brands_set_timestamps
BEFORE INSERT OR UPDATE ON brands
FOR EACH ROW 
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER brands_prevent_created_at_update
BEFORE UPDATE ON brands
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER brands_prevent_created_by_id_change
BEFORE UPDATE ON brands
FOR EACH ROW
EXECUTE FUNCTION prevent_created_by_id_change();


-- categories table triggers
CREATE TRIGGER categories_set_timestamps
BEFORE INSERT OR UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER categories_prevent_created_at_update
BEFORE UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER categories_prevent_created_by_id_change
BEFORE UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION prevent_created_by_id_change();


-- products table triggers
CREATE TRIGGER products_set_timestamps
BEFORE INSERT OR UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER products_prevent_created_at_update
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER products_prevent_created_by_id_change
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION prevent_created_by_id_change();


-- product_images table triggers
CREATE TRIGGER prod_imgs_set_timestamps
BEFORE INSERT OR UPDATE ON product_images
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER prod_imgs_prevent_created_at_update
BEFORE UPDATE ON product_images
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER prod_imgs_prevent_created_by_id_change
BEFORE UPDATE ON product_images
FOR EACH ROW
EXECUTE FUNCTION prevent_created_by_id_change();


-- product_reviews table triggers
CREATE TRIGGER prod_revs_check_user_bought_product
BEFORE INSERT OR UPDATE ON product_reviews
FOR EACH ROW
EXECUTE FUNCTION check_user_bought_product();

CREATE TRIGGER prod_revs_set_timestamps
BEFORE INSERT OR UPDATE ON product_reviews
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER prod_revs_prevent_user_id_change
BEFORE UPDATE ON product_reviews
FOR EACH ROW
EXECUTE FUNCTION prevent_user_id_change();

CREATE TRIGGER prod_revs_validate_approved_by
BEFORE INSERT OR UPDATE ON product_reviews
FOR EACH ROW
EXECUTE FUNCTION validate_approved_by();

CREATE TRIGGER prod_revs_ensure_approved_by_set
BEFORE INSERT OR UPDATE ON product_reviews
FOR EACH ROW
EXECUTE FUNCTION ensure_approved_by_set();

CREATE TRIGGER prod_revs_prevent_created_at_update
BEFORE UPDATE ON product_reviews
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER prod_revs_update_product_rating
BEFORE INSERT OR UPDATE ON product_reviews
FOR EACH ROW
EXECUTE FUNCTION update_product_rating();


-- attributes table triggers
CREATE TRIGGER attributes_set_timestamps
BEFORE INSERT OR UPDATE ON attributes
FOR EACH ROW 
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER attributes_prevent_created_at_update
BEFORE UPDATE ON attributes
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER attributes_prevent_created_by_id_change
BEFORE UPDATE ON attributes
FOR EACH ROW
EXECUTE FUNCTION prevent_created_by_id_change();


-- attribute_values table triggers
CREATE TRIGGER attr_vals_set_timestamps
BEFORE INSERT OR UPDATE ON attribute_values
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER attr_vals_prevent_created_at_update
BEFORE UPDATE ON attribute_values
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER attr_vals_prevent_created_by_id_change
BEFORE UPDATE ON attribute_values
FOR EACH ROW
EXECUTE FUNCTION prevent_created_by_id_change();


-- catalog_managers table triggers
CREATE TRIGGER cm_set_timestamps
BEFORE INSERT OR UPDATE ON catalog_managers
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE INDEX cm_prevent_created_at_update
BEFORE UPDATE ON catalog_managers
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER cm_prevent_user_id_change
BEFORE UPDATE ON catalog_managers
FOR EACH ROW
EXECUTE FUNCTION prevent_user_id_change();

CREATE TRIGGER cm_prevent_created_by_id_change
BEFORE UPDATE ON catalog_managers
FOR EACH ROW
EXECUTE FUNCTION prevent_created_by_id_change();


-- customers table triggers
CREATE TRIGGER customers_set_timestamps
BEFORE INSERT OR UPDATE ON customers
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE INDEX customers_prevent_created_at_update
BEFORE UPDATE ON customers
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER customers_prevent_user_id_change
BEFORE UPDATE ON customers
FOR EACH ROW
EXECUTE FUNCTION prevent_user_id_change();

CREATE TRIGGER customers_prevent_created_by_id_change
BEFORE UPDATE ON customers
FOR EACH ROW
EXECUTE FUNCTION prevent_created_by_id_change();

