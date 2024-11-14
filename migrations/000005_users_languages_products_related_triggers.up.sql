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


-- countries table triggers
CREATE TRIGGER countries_set_timestamps 
BEFORE INSERT OR UPDATE ON countries
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER countries_prevent_created_at_update
BEFORE UPDATE ON countries
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();



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


-- promotions table triggers
CREATE TRIGGER promotions_set_timestamps
BEFORE INSERT OR UPDATE ON promotions
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER promotions_prevent_created_at_update
BEFORE UPDATE ON promotions
FOR EACH ROW 
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER promotions_prevent_created_by_id_change
BEFORE UPDATE ON promotions
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


-- product_price_histories table triggers
CREATE TRIGGER prod_price_hist_prevent_created_at_update
BEFORE UPDATE ON product_price_histories
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();


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

CREATE TRIGGER cm_prevent_created_at_update
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

CREATE TRIGGER customers_prevent_created_at_update
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

