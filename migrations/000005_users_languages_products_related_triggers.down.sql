-- users table triggers
DROP TRIGGER IF EXISTS users_prevent_created_at_update ON users;
DROP TRIGGER IF EXISTS users_set_timestamps ON users;
DROP TRIGGER IF EXISTS users_prevent_delete ON users;
DROP TRIGGER IF EXISTS enforce_banned_constraints_trigger ON users;
DROP TRIGGER IF EXISTS superuser_is_admin_and_is_staff ON users;
DROP TRIGGER IF EXISTS admin_is_staff ON users;

-- user_referrals table triggers
DROP TRIGGER IF EXISTS user_referrals_prevent_created_at_update ON user_referrals;
DROP TRIGGER IF EXISTS user_refferals_set_timeststamps ON user_referrals;
DROP TRIGGER IF EXISTS user_referrals_prevent_user_id_change ON user_referrals;

-- user_product_referrals table triggers
DROP TRIGGER IF EXISTS user_prod_refs_prevent_created_at_update ON user_product_referrals;
DROP TRIGGER IF EXISTS user_prod_refs_set_timestamps ON user_product_referrals;
DROP TRIGGER IF EXISTS user_prod_refs_prevent_user_id_change ON user_product_referrals;

-- user_bough_products table triggers
DROP TRIGGER IF EXISTS user_bght_prods_prevent_created_at_update ON user_bought_products;
DROP TRIGGER IF EXISTS user_bght_prods_set_timestamps ON user_bought_products;

-- languages table triggers
DROP TRIGGER IF EXISTS langs_set_timestamps ON languages;
DROP TRIGGER IF EXISTS langs_prevent_created_at_update ON languages;

-- translations table triggers
DROP TRIGGER IF EXISTS translations_set_timestamps ON translations;
DROP TRIGGER IF EXISTS translations_prevent_created_at_update ON translations;
DROP TRIGGER IF EXISTS user_referrals_prevent_created_by_id_change ON translations;

-- brands table triggers
DROP TRIGGER IF EXISTS brands_set_timestamps ON brands;
DROP TRIGGER IF EXISTS brands_prevent_created_at_update ON brands;
DROP TRIGGER IF EXISTS brands_prevent_created_by_id_change ON brands;

-- categories table triggers
DROP TRIGGER IF EXISTS categories_set_timestamps ON categories;
DROP TRIGGER IF EXISTS categories_prevent_created_at_update ON categories;
DROP TRIGGER IF EXISTS categories_prevent_created_by_id_change ON categories;

-- promotions table triggers
DROP TRIGGER IF EXISTS promotions_set_timestamps ON promotions;
DROP TRIGGER IF EXISTS promotions_prevent_created_at_update ON promotions;
DROP TRIGGER IF EXISTS promotions_prevent_created_by_id_change ON promotions;

-- products table triggers
DROP TRIGGER IF EXISTS products_set_timestamps ON products;
DROP TRIGGER IF EXISTS products_prevent_created_at_update ON products;
DROP TRIGGER IF EXISTS products_prevent_created_by_id_change ON products;

-- product_price_histories table triggers
DROP TRIGGER IF EXISTS prod_price_hist_prevent_created_at_update ON product_price_histories;

-- product_images table triggers
DROP TRIGGER IF EXISTS prod_imgs_set_timestamps ON product_images;
DROP TRIGGER IF EXISTS prod_imgs_prevent_created_at_update ON product_images;
DROP TRIGGER IF EXISTS prod_imgs_prevent_created_by_id_change ON product_images;

-- product_reviews table triggers
DROP TRIGGER IF EXISTS prod_revs_check_user_bought_product ON product_reviews;
DROP TRIGGER IF EXISTS prod_revs_set_timestamps ON product_reviews;
DROP TRIGGER IF EXISTS prod_revs_prevent_user_id_change ON product_reviews;
DROP TRIGGER IF EXISTS prod_revs_validate_approved_by ON product_reviews;
DROP TRIGGER IF EXISTS prod_revs_ensure_approved_by_set ON product_reviews;
DROP TRIGGER IF EXISTS prod_revs_prevent_created_at_update ON product_reviews;
DROP TRIGGER IF EXISTS prod_revs_update_product_rating ON product_reviews;

-- attributes table triggers
DROP TRIGGER IF EXISTS attributes_set_timestamps ON attributes;
DROP TRIGGER IF EXISTS attributes_prevent_created_at_update ON attributes;
DROP TRIGGER IF EXISTS attributes_prevent_created_by_id_change ON attributes;

-- attribute_values table triggers
DROP TRIGGER IF EXISTS attr_vals_set_timestamps ON attribute_values;
DROP TRIGGER IF EXISTS attr_vals_prevent_created_at_update ON attribute_values;
DROP TRIGGER IF EXISTS attr_vals_prevent_created_by_id_change ON attribute_values;

-- catalog_managers table triggers
DROP TRIGGER IF EXISTS cm_set_timestamps ON catalog_managers;
DROP TRIGGER IF EXISTS cm_prevent_created_at_update ON catalog_managers;
DROP TRIGGER IF EXISTS cm_prevent_user_id_change ON catalog_managers;
DROP TRIGGER IF EXISTS cm_prevent_created_by_id_change ON catalog_managers;

-- customers table triggers
DROP TRIGGER IF EXISTS customers_set_timestamps ON customers;
DROP TRIGGER IF EXISTS customers_prevent_created_at_update ON customers;
DROP TRIGGER IF EXISTS customers_prevent_user_id_change ON customers;
DROP TRIGGER IF EXISTS customers_prevent_created_by_id_change ON customers;

