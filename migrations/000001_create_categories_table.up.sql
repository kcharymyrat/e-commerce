-- Set time zone of database to Ashgabat
ALTER DATABASE ecommerce SET TIME ZONE 'Asia/Ashgabat';

-- Create uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create case-insensitive character string type
CREATE EXTENSION IF NOT EXISTS citext;

-- Triggers
CREATE OR REPLACE FUNCTION prevent_created_at_update() RETURN TRIGGER AS $$
BEGIN
    IF NEW.created_at <> OLD.created_at THEN
        RAISE EXCEPTION 'Updating the value for created_at is not allowed'
    ENF IF;
ENG;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION set_timestamps() RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        NEW.created_at := NOW();
        NEW.updated_at := NOW();
        RETURN NEW;
    END IF;

    IF (TG_OP = 'UPDATE') THEN
        IF (NEW.updated_at < OLD.updated_at) THEN
            RAISE EXCEPTION 'updated_at cannot be before the current updated_at value';
        END IF;
        NEW.updated_at := NOW();
        RETURN NEW;
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- users table for authentication
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone varchar(20) UNIQUE NOT NULL,
    password bytea NOT NULL,

    first_name varchar(50),
    last_name varchar(50),
    patronymic varchar(50),
    dob date CHECK (dob BETWEEN '1900-01-01' AND CURRENT_DATE),
    email citext UNIQUE,

    is_active boolean NOT NULL DEFAULT TRUE,
    is_banned boolean NOT NULL DEFAULT FALSE,
    is_trusted boolean NOT NULL DEFAULT FALSE,

    invited_by_id uuid,
    inv_ref_id bigint,
    inv_prod_ref_id bigint,

    ref_signups integer NOT NULL DEFAULT 0 CHECK (ref_signups >= 0),
    prod_ref_signups integer NOT NULL DEFAULT 0 CHECK (prod_ref_signups >= 0),
    prod_ref_bought integer NOT NULL DEFAULT 0 CHECK (prod_ref_bought >= 0),

    total_referrals integer NOT NULL DEFAULT 0 CHECK (total_referrals >= 0),
    _dynamic_discount_percent decimal(5, 2) NOT NULL DEFAULT 0.00,
    dyn_disc_percent decimal(5, 2) GENERATED ALWAYS AS ( 
        CASE 
            WHEN _dynamic_discount_percent >= 10.00 THEN 10.00
            ELSE _dynamic_discount_percent 
        END
    ) STORED
    bonus_points decimal(10, 2) NOT NULL DEFAULT 0.00 CHECK (bonus_points >= 0.00),

    is_staff boolean NOT NULL DEFAULT FALSE,
    is_admin boolean NOT NULL DEFAULT FALSE,
    is_superuser boolean NOT NULL DEFAULT FALSE,

    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    created_by_id uuid,
    updated_by_id uuid,

    CONSTRAINT user_invited_by_id_fk FOREIGN KEY (invited_by_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT user_inv_ref_id_fk FOREIGN KEY (inv_ref_id) REFERENCES user_referral(id) ON DELETE SET NULL,
    CONSTRAINT user_prod_ref_id_fk FOREIGN KEY (inv_prod_ref_id) REFERENCES user_product_referral(id) ON DELETE SET NULL,
    CONSTRAINT user_created_by_id_fk FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT user_updated_by_id_fk FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,

    CHECK valid__dynamic_discount (_dynamic_discount_percent >= 0.00 AND _dynamic_discount_percent <= 100.00),
    CHECK valid_dyn_disc (dyn_disc_percent >= 0.00 AND dyn_disc_percent <= 10.00),
    CHECK valid_updated_at (updated_at >= created_at)
);

CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_first_name ON users(first_name);
CREATE INDEX idx_users_last_name ON users(last_name);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_dob ON users(dob);
CREATE INDEX idx_users_is_staff ON users(is_staff);

CREATE TRIGGER users_prevent_created_at_update
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER users_set_timestamps
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

-- user_referrals table 
CREATE TABLE IF NOT EXISTS user_referrals (
    id bigserial PRIMARY KEY,
    user_id uuid NOT NULL UNIQUE,
    code varchar(50) NOT NULL UNIQUE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    CONSTRAINT user_refferal_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CHECK (updated_at > created_at),
);

CREATE INDEX idx_user_referrals_user_id ON user_referrals(user_id);
CREATE INDEX idx_user_referrals_code ON user_referrals(code);

CREATE TRIGGER user_referrals_prevent_created_at_update
BEFORE UPDATE ON user_referrals
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER user_refferals_set_timeststamps
BEFORE INSERT OR UPDATE ON user_referrals
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();


-- user_product_referrals table
CREATE TABLE IF NOT EXISTS user_product_referrals (
    id bigserial PRIMARY KEY,
    user_id uuid NOT NULL,
    product_id uuid NOT NULL,
    code varchar(50) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), 

    CONSTRAINT user_prod_refs_user_id_fk FOREIGN KEY users(id) ON DELETE CASCADE,
    CONSTRAINT user_prod_refs_product_id_fk FOREIGN kEY products(id) ON DELETE CASCADE,
    
    UNIQUE (user_id, product_id),
    CHECK (updated_at >= created_at)
);

CREATE INDEX idx_user_prod_refs_user_id ON user_product_referrals(user_id);
CREATE INDEX idx_user_prod_refs_product_id ON user_product_referrals(product_id);
CREATE INDEX idx_user_prod_refs_code ON user_product_referrals(code);

CREATE TRIGGER user_prod_refs_prevent_created_at_update
BEFORE UPDATE ON user_product_referrals
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER user_prod_refs_set_timestamps
BEFORE INSERT OR UPDATE ON user_product_referrals
FOR EACH ROW
EXECUTE set_timestamps();


-- user_bough_products table
CREATE TABLE IF NOT EXISTS user_bought_products (
    id bigserial PRIMARY KEY,
    user_id uuid NOT NULL,
    product_id uuid NOT NULL,
    quantity integer NOT NULL DEFAULT 1;
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW();
    update_at timestamp(0) with time zone NOT NULL DEFAULT NOW();

    UNIQUE (user_id, product_id),
    CHECK (updated_at >= created_at)
);

CREATE INDEX idx_user_bght_prods_user_id ON user_bought_products(user_id);
CREATE INDEX idx_user_bght_prods_prod_id ON user_bought_products(product_id);

CREATE TRIGGER user_bght_prods_prevent_created_at_update
BEFORE UPDATE ON user_bought_products
FOR EACH ROW
EXECUTE prevent_created_at_update();

CREATE TRIGGER user_bght_prods_set_timestamps
BEFORE INSERT OR UPDATE ON user_bought_products
FOR EACH ROW
EXECUTE set_timestamps();