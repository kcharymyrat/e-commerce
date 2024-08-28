-- Set time zone of database to Ashgabat
ALTER DATABASE ecommerce SET TIME ZONE 'Asia/Ashgabat';

-- Create uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create case-insensitive character string type
CREATE EXTENSION IF NOT EXISTS citext;

-- users table for authentiction
CREATE TABLE IF NOT EXISTS users(
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

-- indexes for performance
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_first_name ON users(first_name);
CREATE INDEX idx_users_last_name ON users(last_name);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_dob ON users(dob);
CREATE INDEX idx_users_is_staff ON users(is_staff);

-- Trigger to handel created_at and updated_at
CREATE OR REPLACE FUNCTION set_timestamps() RETURNS TRIGGER AS $$
BEGIN
    -- On INSERT, set the created_at and created_by columns only once
    IF (TG_OP = 'INSERT') THEN
        NEW.created_at := NOW();
        NEW.updated_at := NOW();
        RETURN NEW;
    END IF;

    -- On UPDATE, modify updated_at and updated_by, and ensure updated_at is not before the previous value
    IF (TG_OP = 'UPDATE') THEN
        IF (NEW.updated_at < OLD.updated_at) THEN
            RAISE EXCEPTION 'updated_at cannot be before the current updated_at value';
        END IF;
        NEW.updated_at := NOW();
        RETURN NEW;
    END IF;

    RETURN NULL; -- Result is ignored since this is a trigger
END;
$$ LANGUAGE plpgsql;

-- Apply the above trigger for users table
CREATE TRIGGER users_set_timestamps
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();


CREATE TABLE IF NOT EXISTS user_referrals (
    id bigserial PRIMARY KEY,
    user_id uuid NOT NULL UNIQUE,
    code varchar(50) NOT NULL UNIQUE,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    CHECK (updated_at > created_at),
    CONSTRAINT user_refferal_user_id_fk FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX idx_user_user_id ON user_referrals(user_id);
CREATE INDEX idx_user_referal_code ON user_referrals(code);


CREATE TRIGGER user_refferals_set_timeststamps()
BEFORE INSERT OR UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

