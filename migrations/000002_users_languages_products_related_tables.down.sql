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
    inv_ref_id uuid,
    inv_prod_ref_id uuid,

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

    CONSTRAINT users_invited_by_id_fk FOREIGN KEY (invited_by_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT users_inv_ref_id_fk FOREIGN KEY (inv_ref_id) REFERENCES user_referrals(id) ON DELETE SET NULL,
    CONSTRAINT users_prod_ref_id_fk FOREIGN KEY (inv_prod_ref_id) REFERENCES user_product_referrals(id) ON DELETE SET NULL,
    CONSTRAINT users_created_by_id_fk FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT users_updated_by_id_fk FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE SET NULL,

    CHECK valid__dynamic_discount (_dynamic_discount_percent >= 0.00 AND _dynamic_discount_percent <= 100.00),
    CHECK valid_dyn_disc (dyn_disc_percent >= 0.00 AND dyn_disc_percent <= 10.00),
    CHECK valid_updated_at (updated_at >= created_at)
);

CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_first_name ON users(first_name);
CREATE INDEX IF NOT EXISTS idx_users_last_name ON users(last_name);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_dob ON users(dob);
CREATE INDEX IF NOT EXISTS idx_users_is_staff ON users(is_staff);

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
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL UNIQUE,
    code varchar(50) NOT NULL UNIQUE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    CONSTRAINT user_refferal_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CHECK (updated_at > created_at),
);

CREATE INDEX IF NOT EXISTS idx_user_referrals_user_id ON user_referrals(user_id);
CREATE INDEX IF NOT EXISTS idx_user_referrals_code ON user_referrals(code);

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
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
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

CREATE INDEX IF NOT EXISTS idx_user_prod_refs_user_id ON user_product_referrals(user_id);
CREATE INDEX IF NOT EXISTS idx_user_prod_refs_product_id ON user_product_referrals(product_id);
CREATE INDEX IF NOT EXISTS idx_user_prod_refs_code ON user_product_referrals(code);

CREATE TRIGGER user_prod_refs_prevent_created_at_update
BEFORE UPDATE ON user_product_referrals
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER user_prod_refs_set_timestamps
BEFORE INSERT OR UPDATE ON user_product_referrals
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();


-- user_bough_products table
CREATE TABLE IF NOT EXISTS user_bought_products (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    product_id uuid NOT NULL,
    quantity integer NOT NULL DEFAULT 1;
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW();
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW();

    UNIQUE (user_id, product_id),
    CHECK (updated_at >= created_at)
);

CREATE INDEX IF NOT EXISTS idx_user_bght_prods_user_id ON user_bought_products(user_id);
CREATE INDEX IF NOT EXISTS idx_user_bght_prods_prod_id ON user_bought_products(product_id);

CREATE TRIGGER user_bght_prods_prevent_created_at_update
BEFORE UPDATE ON user_bought_products
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();

CREATE TRIGGER user_bght_prods_set_timestamps
BEFORE INSERT OR UPDATE ON user_bought_products
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();


-- languages table
CREATE TABLE IF NOT EXISTS languages (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    code varchar(10) NOT NULL UNIQUE,
    name varchar(50) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,

    CHECK (updated_at >= created_at),
    CONSTRAINT langs_created_by_id_fk FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT langs_updated_by_id_fk FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_languages_code ON languages(code);

CREATE TRIGGER langs_set_timestamps 
BEFORE INSERT OR UPDATE ON languages
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER langs_prevent_created_at_update
BEFORE UPDATE ON languages
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();


-- translations table
CREATE TABLE IF NOT EXISTS translations (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    language_id uuid NOT NULL, 
    entity_id uuid NOT NULL,
    table_name text NOT NULL,
    field_name text NOT NULL,
    translation text NOT NULL,

    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL, 
    updated_by_id uuid NOT NULL,

    CHECK (updated_at >= created_at),
    CONSTRAINT translations_lang_id_fk FOREIGN KEY (language_id) REFERENCES languages(id) ON DELETE RESTRICT,
    CONSTRAINT translations_created_by_id_fk FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT translations_updated_by_id_fk FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE RESTRICT,
);

CREATE INDEX IF NOT EXISTS idx_translations_lang_id ON translations(language_id);
CREATE INDEX IF NOT EXISTS idx_translations_table_name ON translations(table_name);
CREATE INDEX IF NOT EXISTS idx_translations_entity_id ON translations(entity_id);


CREATE TRIGGER translations_set_timestamps
BEFORE INSERT OR UPDATE ON translations
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER translations_prevent_created_at_update
BEFORE UPDATE ON translations
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();


-- brand table
CREATE TABLE IF NOT EXISTS brands (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    logo_url text NOT NULL UNIQUE,
    title varchar(50) NOT NULL UNIQUE,
    slug varchar(50) NOT NULL UNIQUE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,

    CHECK (updated_at >= created_at),
    CONSTRAINT brands_created_by_id_fk FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT brands_updated_by_id_fk FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_brands_title ON brands(title);
CREATE INDEX IF NOT EXISTS idx_brands_slug ON brands(slug);

CREATE TRIGGER brands_set_timestamps
BEFORE INSERT OR UPDATE ON brands
FOR EACH ROW 
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER brands_prevent_created_at_update
BEFORE UPDATE ON brands
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();


-- categories table
CREATE TABLE IF NOT EXISTS categories (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent uuid,
    name varchar(50),
    slug varchar(50),
    description text,
    image_url text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,

    CHECK (updated_at >= created_at),
    CONSTRAINT categories_created_by_id_fk FOREIGN KEY (created_by_id) ON DELETE RESTRICT,
    CONSTRAINT categories_created_by_id_fk FOREIGN KEY (updated_by_id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name);
CREATE INDEX IF NOT EXISTS idx_categories_slug on categories(slug);


CREATE TRIGGER categories_set_timestamps
BEFORE INSERT OR UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER categories_prevent_created_at_update
BEFORE UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION prevent_created_at_update();


-- products table
CREATE TABLE IF NOT EXISTS products (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(50) NOT NULL UNIQUE,
    slug varchar(50) NOT NULL UNIQUE,
    description text,
    code varchar(32) NOT NULL UNIQUE,
    weight_kg decimal(5, 2) NOT NULL DEFAULT 0.00 CHECK (weight_kg >= 0.00),
    stock_amount integer NOT NULL DEFAULT 0 CHECK (stock_amount >= 0),
    is_adult boolean NOT NULL DEFAULT FALSE,
    is_new boolean NOT NULL DEFAULT TRUE,
    is_active boolean NOT NULL DEFAULT TRUE,
    in_stock boolean NOT NULL GENERATED ALWAYS AS (
        CASE 
            WHEN stock_amount >= 0 THEN TRUE
            ELSE FALSE
        END
    )STORED,
    price decimal(10, 2) NOT NULL DEFAULT 0.00 CHECK (price >= 0.00),
    sale_percent decimal(5, 2) NOT NULL DEFAULT 0.00 CHECK (price >= 0.00 AND price <= 100.00),
    sale_price decimal(10, 2) NOT NULL GENERATED ALWAYS AS (
        price * (1::decimal - sale_percent / 100::decimal)
    )STORED
    image text NOT NULL,
    thumbnail text NOT NULL,
    video text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,


    CHECK (updated_at >= created_at),
    CHECK (sale_price <= price),
    CONSTRAINT products_created_by_id_fk FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE RESTRICT,
    CONSTRAINT products_updated_by_id_fk FOREIGN KEY  (updated_by_id) REFERENCES users(id) ON DELETE RESTRICT
);


CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
CREATE INDEX IF NOT EXISTS idx_products_code ON products(code);
CREATE INDEX IF NOT EXISTS idx_products_is_new ON products(is_new);
CREATE INDEX IF NOT EXISTS idx_products_sale_percent ON products(sale_percent); 