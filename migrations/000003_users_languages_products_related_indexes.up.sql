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
CREATE INDEX IF NOT EXISTS idx_brands_title ON brands(title);
CREATE INDEX IF NOT EXISTS idx_brands_slug ON brands(slug);

-- categories table indexes
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
CREATE INDEX IF NOT EXISTS idx_categories_slug on categories(slug);
CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);

-- products table indexes
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
CREATE INDEX IF NOT EXISTS idx_products_code ON products(code);
CREATE INDEX IF NOT EXISTS idx_products_is_new ON products(is_new);
CREATE INDEX IF NOT EXISTS idx_products_sale_percent ON products(sale_percent); 

-- products_brands table indexes
CREATE INDEX IF NOT EXISTS idx_products_brands_product_id ON products_brands(product_id);
CREATE INDEX IF NOT EXISTS idx_products_brands_brand_id ON products_brands(brand_id);
CREATE INDEX IF NOT EXISTS idx_products_brands_together ON products_brands(product_id, brand_id);

-- products_categories table indexes
CREATE INDEX IF NOT EXISTS idx_products_categories_product_id ON products_categories(product_id);
CREATE INDEX IF NOT EXISTS idx_products_categories_category_id ON products_categories(category_id);
CREATE INDEX IF NOT EXISTS idx_products_categories_together ON products_categories(product_id, category_id);

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
