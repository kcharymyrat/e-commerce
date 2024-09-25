CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    phone VARCHAR(15) NOT NULL UNIQUE CHECK (phone ~ '^\+[1-9][0-9]{7,14}$'),
    password_hash bytea NOT NULL,

    first_name varchar(50),
    last_name varchar(50),
    patronymic varchar(50),
    dob date CHECK (dob BETWEEN '1900-01-01' AND CURRENT_DATE),
    email citext UNIQUE,

    is_active boolean NOT NULL DEFAULT FALSE,
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
    )STORED,
    bonus_points decimal(10, 2) NOT NULL DEFAULT 0.00 CHECK (bonus_points >= 0.00),

    is_staff boolean NOT NULL DEFAULT FALSE,
    is_admin boolean NOT NULL DEFAULT FALSE,
    is_superuser boolean NOT NULL DEFAULT FALSE,

    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    created_by_id uuid,
    updated_by_id uuid,

    version integer NOT NULL DEFAULT 1,

    CHECK (_dynamic_discount_percent >= 0.00 AND _dynamic_discount_percent <= 100.00),
    CHECK (dyn_disc_percent >= 0.00 AND dyn_disc_percent <= 10.00),
    CHECK (updated_at >= created_at)
);


CREATE TABLE IF NOT EXISTS user_referrals (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL UNIQUE,
    code varchar(50) NOT NULL UNIQUE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1,

    CHECK (updated_at > created_at)
);


CREATE TABLE IF NOT EXISTS user_product_referrals (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    product_id uuid NOT NULL,
    code varchar(50) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), 
    version integer NOT NULL DEFAULT 1,

    UNIQUE (user_id, product_id),
    CHECK (updated_at >= created_at)
);


CREATE TABLE IF NOT EXISTS user_bought_products (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    product_id uuid NOT NULL,
    quantity integer NOT NULL DEFAULT 1,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1,

    UNIQUE (user_id, product_id),
    CHECK (updated_at >= created_at)
);


CREATE TABLE IF NOT EXISTS languages (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    code varchar(10) NOT NULL UNIQUE,
    name varchar(50) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,
    version integer NOT NULL DEFAULT 1,

    CHECK (updated_at >= created_at)
);


CREATE TABLE IF NOT EXISTS translations (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    language_code varchar(10) NOT NULL,
    entity_id uuid NOT NULL,
    table_name varchar(50) NOT NULL,
    field_name varchar(50) NOT NULL,
    translated_value text NOT NULL,

    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL, 
    updated_by_id uuid NOT NULL,

    version integer NOT NULL DEFAULT 1,

    CHECK (updated_at >= created_at),
    UNIQUE (entity_id, language_code, table_name, field_name)
);


CREATE TABLE IF NOT EXISTS brands (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    logo_url text NOT NULL UNIQUE,
    title varchar(50) NOT NULL UNIQUE,
    slug varchar(50) NOT NULL UNIQUE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,
    version integer NOT NULL DEFAULT 1,

    CHECK (updated_at >= created_at)
);


CREATE TABLE IF NOT EXISTS categories (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    parent_id uuid,
    name varchar(50) NOT NULL,
    slug varchar(50) NOT NULL,
    description text,
    image_url text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,
    version integer NOT NULL DEFAULT 1,

    CHECK (updated_at >= created_at)
);


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
            WHEN stock_amount > 0 THEN TRUE
            ELSE FALSE
        END
    )STORED,
    price decimal(10, 2) NOT NULL DEFAULT 0.00 CHECK (price >= 0.00),
    sale_percent decimal(5, 2) NOT NULL DEFAULT 0.00 CHECK (price >= 0.00 AND price <= 100.00),
    sale_price decimal(10, 2) NOT NULL GENERATED ALWAYS AS (
        price * (1::decimal - sale_percent / 100::decimal)
    )STORED,
    image text NOT NULL,
    thumbnail text NOT NULL,
    video text NOT NULL,
    average_rating decimal(3, 2) NOT NULL DEFAULT 0.00 CHECK (average_rating BETWEEN 0.00 AND 5.00),
    number_of_reviews integer DEFAULT 0 CHECK (number_of_reviews >= 0),
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,
    version integer NOT NULL DEFAULT 1,


    CHECK (updated_at >= created_at),
    CHECK (sale_price <= price)
);


CREATE TABLE IF NOT EXISTS products_brands (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id uuid NOT NULL,
    brand_id uuid NOT NULL,

    UNIQUE (product_id, brand_id)
);


CREATE TABLE IF NOT EXISTS products_categories (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id uuid NOT NULL,
    category_id uuid NOT NULL,

    UNIQUE (product_id, category_id)
);


CREATE TABLE IF NOT EXISTS product_images (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id uuid NOT NULL,
    image_url text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,
    version integer NOT NULL DEFAULT 1,

    UNIQUE (product_id, image_url),
    CHECK (updated_at >= created_at)
);


CREATE TABLE IF NOT EXISTS product_reviews ( 
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id uuid NOT NULL,
    user_id uuid NOT NULL, 
    rating decimal(3, 2) NOT NULL CHECK (rating BETWEEN 0.00 AND 5.00),
    review_text text,
    image_url text,
    video_url text,
    approved_by_id uuid, 
    is_approved boolean NOT NULL DEFAULT FALSE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1,

    CHECK (updated_at >= created_at),
    UNIQUE (product_id, user_id)
);


CREATE TABLE IF NOT EXISTS attributes (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    name varchar(50) NOT NULL UNIQUE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,
    version integer NOT NULL DEFAULT 1,

    CHECK (updated_at >= created_at)
);


CREATE TABLE IF NOT EXISTS attribute_values (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id uuid NOT NULL,
    attribute_id uuid NOT NULL,
    value varchar(255) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,
    version integer NOT NULL DEFAULT 1,

    CHECK (updated_at >= created_at),
    UNIQUE (product_id, attribute_id)
);


CREATE TABLE IF NOT EXISTS catalog_managers (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL UNIQUE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,
    version integer NOT NULL DEFAULT 1,

    CHECK (updated_at >= created_at)
);


CREATE TABLE IF NOT EXISTS customers (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL UNIQUE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    created_by_id uuid NOT NULL,
    updated_by_id uuid NOT NULL,
    version integer NOT NULL DEFAULT 1,

    CHECK (updated_at >= created_at)
);
