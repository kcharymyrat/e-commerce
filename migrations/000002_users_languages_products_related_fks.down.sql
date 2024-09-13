-- users table fk constraints
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_invited_by_id_fk;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_inv_ref_id_fk;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_prod_ref_id_fk;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_created_by_id_fk;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_updated_by_id_fk;

-- user_referrals table fk constraints
ALTER TABLE user_referrals DROP CONSTRAINT IF EXISTS user_refferal_user_id_fk;

-- user_product_referrals table fk constraints
ALTER TABLE user_product_referrals DROP CONSTRAINT IF EXISTS user_prod_refs_user_id_fk;
ALTER TABLE user_product_referrals DROP CONSTRAINT IF EXISTS user_prod_refs_product_id_fk;

-- languages table fk constraints
ALTER TABLE languages DROP CONSTRAINT IF EXISTS langs_created_by_id_fk;
ALTER TABLE languages DROP CONSTRAINT IF EXISTS langs_updated_by_id_fk;

-- translations table fk constraints
ALTER TABLE translations DROP CONSTRAINT IF EXISTS translations_language_code_fk;
ALTER TABLE translations DROP CONSTRAINT IF EXISTS translations_created_by_id_fk;
ALTER TABLE translations DROP CONSTRAINT IF EXISTS translations_updated_by_id_fk;

-- brands table fk constraints
ALTER TABLE brands DROP CONSTRAINT IF EXISTS brands_created_by_id_fk;
ALTER TABLE brands DROP CONSTRAINT IF EXISTS brands_updated_by_id_fk;

-- categories table fk constraints
ALTER TABLE categories DROP CONSTRAINT IF EXISTS categories_created_by_id_fk;
ALTER TABLE categories DROP CONSTRAINT IF EXISTS categories_updated_by_id_fk;

-- products table fk constraints
ALTER TABLE products DROP CONSTRAINT IF EXISTS products_created_by_id_fk;
ALTER TABLE products DROP CONSTRAINT IF EXISTS products_updated_by_id_fk;

-- products_brands table fk constraints
ALTER TABLE products_brands DROP CONSTRAINT IF EXISTS products_brands_product_id_fk;
ALTER TABLE products_brands DROP CONSTRAINT IF EXISTS products_brands_brand_id_fk;

-- products_categories table fk constraints
ALTER TABLE products_categories DROP CONSTRAINT IF EXISTS products_categories_product_id_fk;
ALTER TABLE products_categories DROP CONSTRAINT IF EXISTS products_categories_category_id_fk;

-- product_images table fk constraints
ALTER TABLE product_images DROP CONSTRAINT IF EXISTS prod_imgs_product_id_fk;
ALTER TABLE product_images DROP CONSTRAINT IF EXISTS prod_imgs_created_by_id_fk;
ALTER TABLE product_images DROP CONSTRAINT IF EXISTS prod_imgs_updated_by_id_fk;

-- product_reviews table fk constraints
ALTER TABLE product_reviews DROP CONSTRAINT IF EXISTS prod_revs_product_id_fk;
ALTER TABLE product_reviews DROP CONSTRAINT IF EXISTS prod_revs_user_id_fk;
ALTER TABLE product_reviews DROP CONSTRAINT IF EXISTS prod_revs_approved_by_id_fk;

-- attributes table fk constraints
ALTER TABLE attributes DROP CONSTRAINT IF EXISTS attributes_created_by_id_fk;
ALTER TABLE attributes DROP CONSTRAINT IF EXISTS attributes_updated_by_id_fk;

-- attribute_values table fk constraints
ALTER TABLE attribute_values DROP CONSTRAINT IF EXISTS attr_vals_product_id_fk;
ALTER TABLE attribute_values DROP CONSTRAINT IF EXISTS attr_vals_created_by_id_fk;
ALTER TABLE attribute_values DROP CONSTRAINT IF EXISTS atts_vals_updated_by_id_fk;

-- catalog_managers table fk constraints
ALTER TABLE catalog_managers DROP CONSTRAINT IF EXISTS cm_user_id_fk;
ALTER TABLE catalog_managers DROP CONSTRAINT IF EXISTS cm_created_by_id_fk;
ALTER TABLE catalog_managers DROP CONSTRAINT IF EXISTS cm_updated_by_id_fk;

-- customers table fk constraints
ALTER TABLE customers DROP CONSTRAINT IF EXISTS customers_user_id_fk;
ALTER TABLE customers DROP CONSTRAINT IF EXISTS customers_created_by_id_fk;
ALTER TABLE customers DROP CONSTRAINT IF EXISTS customers_updated_by_id_fk;
