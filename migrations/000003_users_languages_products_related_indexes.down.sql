-- users table indexes
DROP INDEX IF EXISTS idx_users_phone;
DROP INDEX IF EXISTS idx_users_first_name;
DROP INDEX IF EXISTS idx_users_last_name;
DROP INDEX IF EXISTS idx_users_full_name;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_is_staff;

-- user_referrals table indexes
DROP INDEX IF EXISTS idx_user_referrals_user_id;
DROP INDEX IF EXISTS idx_user_referrals_code;

-- user_product_referrals table indexes
DROP INDEX IF EXISTS idx_user_prod_refs_user_id;
DROP INDEX IF EXISTS idx_user_prod_refs_product_id;
DROP INDEX IF EXISTS idx_user_prod_refs_code;

-- user_bough_products table indexes
DROP INDEX IF EXISTS idx_user_bght_prods_user_id;
DROP INDEX IF EXISTS idx_user_bght_prods_prod_id;

-- languages table indexes
DROP INDEX IF EXISTS idx_languages_code;
DROP INDEX IF EXISTS idx_languages_name;

-- translations table indexes
DROP INDEX IF EXISTS idx_translations_lang_code;
DROP INDEX IF EXISTS idx_translations_table_name;
DROP INDEX IF EXISTS idx_translations_field_name;
DROP INDEX IF EXISTS idx_translations_entity_id;
DROP INDEX IF EXISTS idx_translations_entity_language;

-- brands table indexes
DROP INDEX IF EXISTS idx_brands_title;
DROP INDEX IF EXISTS idx_brands_slug;

-- categories table indexes
DROP INDEX IF EXISTS idx_categories_name;
DROP INDEX IF EXISTS idx_categories_slug;
DROP INDEX IF EXISTS idx_categories_parent_id;

-- promotions table indexes
DROP INDEX IF EXISTS idx_promotions_type;
DROP INDEX IF EXISTS idx_promotions_sale_percent;
DROP INDEX IF EXISTS idx_promotions_start_date;
DROP INDEX IF EXISTS idx_promotions_end_date;
DROP INDEX IF EXISTS idx_promotions_is_active;

-- products table indexes
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_slug;
DROP INDEX IF EXISTS idx_products_code;
DROP INDEX IF EXISTS idx_products_is_new;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_price;

-- products_brands table indexes
DROP INDEX IF EXISTS idx_products_brands_product_id;
DROP INDEX IF EXISTS idx_products_brands_brand_id;
DROP INDEX IF EXISTS idx_products_brands_together;

-- products_categories table indexes
DROP INDEX IF EXISTS idx_products_categories_product_id;
DROP INDEX IF EXISTS idx_products_categories_category_id;
DROP INDEX IF EXISTS idx_products_categories_together;

-- products_promotions table indexes
DROP INDEX IF EXISTS idx_products_promotions_product_id;
DROP INDEX IF EXISTS idx_products_promotions_promotion_id;
DROP INDEX IF EXISTS idx_products_promotions_together;

-- product_price_histories table indexes
DROP INDEX IF EXISTS idx_prod_price_hist_product_id;
DROP INDEX IF EXISTS idx_prod_price_hist_created_at;

-- product_images table indexes
DROP INDEX IF EXISTS idx_prod_imgs_product_id;

-- product_reviews table indexes
DROP INDEX IF EXISTS idx_prod_revs_product_id;
DROP INDEX IF EXISTS idx_prod_revs_user_id;
DROP INDEX IF EXISTS idx_prod_revs_rating;

-- attributes table indexes
DROP INDEX IF EXISTS idx_attributes_name;

-- attribute_values table indexes
DROP INDEX IF EXISTS idx_attr_vals_product_id;
DROP INDEX IF EXISTS idx_attr_vals_attribute_id;
DROP INDEX IF EXISTS idx_attr_vals_attribute_product_ids;

-- catalog_managers table indexes
DROP INDEX IF EXISTS idx_cm_user_id;

-- customers table indexes
DROP INDEX IF EXISTS idx_customers_user_id;
