-- users table fk constraints
ALTER TABLE users
ADD CONSTRAINT users_invited_by_id_fk FOREIGN KEY (invited_by_id) 
REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE users
ADD CONSTRAINT users_inv_ref_id_fk FOREIGN KEY (inv_ref_id) 
REFERENCES user_referrals(id) ON DELETE SET NULL;

ALTER TABLE users
ADD CONSTRAINT users_prod_ref_id_fk FOREIGN KEY (inv_prod_ref_id) 
REFERENCES user_product_referrals(id) ON DELETE SET NULL;

ALTER TABLE users
ADD CONSTRAINT users_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE users
ADD CONSTRAINT users_updated_by_id_fk 
FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE SET NULL;


-- user_referrals fk contraints
ALTER TABLE user_referrals
ADD CONSTRAINT user_refferal_user_id_fk FOREIGN KEY (user_id) 
REFERENCES users(id) ON DELETE RESTRICT;


-- user_product_referrals fk constraints
ALTER TABLE user_product_referrals
ADD CONSTRAINT user_prod_refs_user_id_fk FOREIGN KEY (user_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE user_product_referrals
ADD CONSTRAINT user_prod_refs_product_id_fk FOREIGN KEY (product_id) 
REFERENCES products(id) ON DELETE CASCADE;


-- languages table fk constraints
ALTER TABLE languages
ADD CONSTRAINT langs_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE languages
ADD CONSTRAINT langs_updated_by_id_fk FOREIGN KEY (updated_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;


-- translations table fk constraints
ALTER TABLE translations
ADD CONSTRAINT translations_language_code_fk FOREIGN KEY (language_code) 
REFERENCES languages(code) ON DELETE RESTRICT;

ALTER TABLE translations
ADD CONSTRAINT translations_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE translations
ADD CONSTRAINT translations_updated_by_id_fk FOREIGN KEY (updated_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;


-- brands tables fk constraints
ALTER TABLE brands
ADD CONSTRAINT brands_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE brands
ADD CONSTRAINT brands_updated_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;


-- categories table fk constraints
ALTER TABLE categories
ADD CONSTRAINT categories_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE categories
ADD CONSTRAINT categories_updated_by_id_fk FOREIGN KEY (updated_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;


-- products table fk constraints
ALTER TABLE products
ADD CONSTRAINT products_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE products
ADD CONSTRAINT products_updated_by_id_fk FOREIGN KEY (updated_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;


-- products_brands table fk constraints
ALTER TABLE products_brands
ADD CONSTRAINT products_brands_product_id_fk FOREIGN KEY (product_id) 
REFERENCES products(id) ON DELETE CASCADE;

ALTER TABLE products_brands
ADD CONSTRAINT products_brands_brand_id_fk FOREIGN KEY (brand_id) 
REFERENCES brands(id) ON DELETE CASCADE;


-- products_categories table fk constaints
ALTER TABLE products_categories
ADD CONSTRAINT products_categories_product_id_fk FOREIGN KEY (product_id) 
REFERENCES products(id) ON DELETE CASCADE;

ALTER TABLE products_categories
ADD CONSTRAINT products_categories_category_id_fk FOREIGN KEY (category_id) 
REFERENCES categories(id) ON DELETE CASCADE;


-- product_images table fk constaints
ALTER TABLE product_images
ADD CONSTRAINT prod_imgs_product_id_fk FOREIGN KEY (product_id) 
REFERENCES products(id) ON DELETE CASCADE;

ALTER TABLE product_images
ADD CONSTRAINT prod_imgs_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE product_images
ADD CONSTRAINT prod_imgs_updated_by_id_fk FOREIGN KEY (updated_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;


-- product_reviews table fk constaints
ALTER TABLE product_reviews
ADD CONSTRAINT prod_revs_product_id_fk FOREIGN KEY (product_id) 
REFERENCES products(id) ON DELETE CASCADE;

ALTER TABLE product_reviews
ADD CONSTRAINT prod_revs_user_id_fk FOREIGN KEY (user_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE product_reviews
ADD CONSTRAINT prod_revs_approved_by_id_fk FOREIGN KEY (approved_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;


-- attributes table fk constraints
ALTER TABLE attributes
ADD CONSTRAINT attributes_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE attributes
ADD CONSTRAINT attributes_updated_by_id_fk FOREIGN KEY (updated_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

-- attribute_values table fk constraints
ALTER TABLE attribute_values
ADD CONSTRAINT attr_vals_product_id_fk FOREIGN KEY (product_id) 
REFERENCES products(id) ON DELETE CASCADE;

ALTER TABLE attribute_values
ADD CONSTRAINT attr_vals_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE attribute_values
ADD CONSTRAINT atts_vals_updated_by_id_fk FOREIGN KEY (updated_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;


-- catalog_managers table fk constraints
ALTER TABLE catalog_managers
ADD CONSTRAINT cm_user_id_fk FOREIGN KEY (user_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE catalog_managers
ADD CONSTRAINT cm_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE catalog_managers
ADD CONSTRAINT cm_updated_by_id_fk FOREIGN KEY (updated_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;


-- customers table fk constraints
ALTER TABLE customers
ADD CONSTRAINT customers_user_id_fk FOREIGN KEY (user_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE customers
ADD CONSTRAINT customers_created_by_id_fk FOREIGN KEY (created_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE customers
ADD CONSTRAINT customers_updated_by_id_fk FOREIGN KEY (updated_by_id) 
REFERENCES users(id) ON DELETE RESTRICT;
