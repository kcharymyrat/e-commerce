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
