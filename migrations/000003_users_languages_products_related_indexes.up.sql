-- users table indexes
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_first_name ON users(first_name);
CREATE INDEX IF NOT EXISTS idx_users_last_name ON users(last_name);
CREATE INDEX IF NOT EXISTS idx_users_full_name ON users(first_name, last_name);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_is_staff ON users(is_staff);

-- user_referrals table indexes
CREATE INDEX IF NOT EXISTS idx_user_referrals_user_id ON user_referrals(user_id);
CREATE INDEX IF NOT EXISTS idx_user_referrals_code ON user_referrals(code);

-- user_product_referrals table indexes
CREATE INDEX IF NOT EXISTS idx_user_prod_refs_user_id ON user_product_referrals(user_id);
CREATE INDEX IF NOT EXISTS idx_user_prod_refs_product_id ON user_product_referrals(product_id);
CREATE INDEX IF NOT EXISTS idx_user_prod_refs_code ON user_product_referrals(code);

-- user_bough_products table indexes
CREATE INDEX IF NOT EXISTS idx_user_bght_prods_user_id ON user_bought_products(user_id);
CREATE INDEX IF NOT EXISTS idx_user_bght_prods_prod_id ON user_bought_products(product_id);

-- languages table indexes
CREATE INDEX IF NOT EXISTS idx_languages_code ON languages(code);
CREATE INDEX IF NOT EXISTS idx_languages_name ON languages(name);

-- translations table indexes
CREATE INDEX IF NOT EXISTS idx_translations_lang_code ON translations(language_code);
CREATE INDEX IF NOT EXISTS idx_translations_table_name ON translations(table_name);
CREATE INDEX IF NOT EXISTS idx_translations_field_name ON translations(field_name);
CREATE INDEX IF NOT EXISTS idx_translations_entity_id ON translations(entity_id);
CREATE INDEX IF NOT EXISTS idx_translations_entity_language ON translations(entity_id, language_code);

-- brands table indexes
CREATE INDEX IF NOT EXISTS idx_brands_name ON brands(name);
CREATE INDEX IF NOT EXISTS idx_brands_slug ON brands(slug);

-- categories table indexes
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
CREATE INDEX IF NOT EXISTS idx_categories_slug on categories(slug);
CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);

-- promotions table indexes
CREATE INDEX IF NOT EXISTS idx_promotions_type ON promotions(type);
CREATE INDEX IF NOT EXISTS idx_promotions_sale_percent ON promotions(sale_percent);
CREATE INDEX IF NOT EXISTS idx_promotions_start_date ON promotions(start_date);
CREATE INDEX IF NOT EXISTS idx_promotions_end_date ON promotions(end_date);
CREATE INDEX IF NOT EXISTS idx_promotions_is_active ON promotions(is_active);

-- products table indexes
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
CREATE INDEX IF NOT EXISTS idx_products_code ON products(code);
CREATE INDEX IF NOT EXISTS idx_products_is_new ON products(is_new);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);
CREATE INDEX IF NOT EXISTS idx_products_price ON products(price); 

-- products_brands table indexes
CREATE INDEX IF NOT EXISTS idx_products_brands_product_id ON products_brands(product_id);
CREATE INDEX IF NOT EXISTS idx_products_brands_brand_id ON products_brands(brand_id);
CREATE INDEX IF NOT EXISTS idx_products_brands_together ON products_brands(product_id, brand_id);

-- products_categories table indexes
CREATE INDEX IF NOT EXISTS idx_products_categories_product_id ON products_categories(product_id);
CREATE INDEX IF NOT EXISTS idx_products_categories_category_id ON products_categories(category_id);
CREATE INDEX IF NOT EXISTS idx_products_categories_together ON products_categories(product_id, category_id);

-- products_promotions table indexes
CREATE INDEX IF NOT EXISTS idx_products_promotions_product_id ON products_promotions(product_id);
CREATE INDEX IF NOT EXISTS idx_products_promotions_promotion_id ON products_promotions(promotion_id);
CREATE INDEX IF NOT EXISTS idx_products_promotions_together ON products_promotions(product_id, promotion_id);

-- product_price_histories table indexes
CREATE INDEX IF NOT EXISTS idx_prod_price_hist_product_id ON product_price_histories(product_id);
CREATE INDEX IF NOT EXISTS idx_prod_price_hist_created_at ON product_price_histories(created_at);

-- product_images table indexes
CREATE INDEX IF NOT EXISTS idx_prod_imgs_product_id ON product_images(product_id);

-- product_reviews table indexes
CREATE INDEX IF NOT EXISTS idx_prod_revs_product_id ON product_reviews(product_id);
CREATE INDEX IF NOT EXISTS idx_prod_revs_user_id ON product_reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_prod_revs_rating ON product_reviews(rating);

-- attributes table indexes
CREATE INDEX IF NOT EXISTS idx_attributes_name ON attributes(name);

-- attribute_values table indexes
CREATE INDEX IF NOT EXISTS idx_attr_vals_product_id ON attribute_values(product_id);
CREATE INDEX IF NOT EXISTS idx_attr_vals_attribute_id ON attribute_values(attribute_id);
CREATE INDEX IF NOT EXISTS idx_attr_vals_attribute_product_ids ON attribute_values(attribute_id, product_id);

-- catalog_managers table indexes
CREATE INDEX IF NOT EXISTS idx_cm_user_id ON catalog_managers(user_id);

-- customers table indexes
CREATE INDEX IF NOT EXISTS idx_customers_user_id ON customers(user_id);
