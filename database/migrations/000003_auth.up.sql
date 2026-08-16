CREATE TABLE IF NOT EXISTS auth_identities (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 guardian_id uuid REFERENCES guardians(id) ON DELETE CASCADE,
 school_id uuid REFERENCES schools(id) ON DELETE CASCADE,
 role varchar(40) NOT NULL CHECK(role IN('PARENT','SCHOOL_ADMIN','CASHIER','REGISTRAR','SUPER_ADMIN')),
 phone varchar(32),
 email varchar(255),
 status varchar(20) NOT NULL DEFAULT 'ACTIVE',
 created_at timestamptz NOT NULL DEFAULT now(),
 updated_at timestamptz NOT NULL DEFAULT now(),
 CHECK(guardian_id IS NOT NULL OR school_id IS NOT NULL OR role='SUPER_ADMIN')
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_auth_identity_phone ON auth_identities(phone) WHERE phone IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS ux_auth_identity_email ON auth_identities(lower(email)) WHERE email IS NOT NULL;
CREATE TABLE IF NOT EXISTS auth_otps (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 identity_id uuid NOT NULL REFERENCES auth_identities(id) ON DELETE CASCADE,
 code_hash varchar(64) NOT NULL,
 purpose varchar(30) NOT NULL DEFAULT 'LOGIN',
 expires_at timestamptz NOT NULL,
 consumed_at timestamptz,
 attempts int NOT NULL DEFAULT 0,
 created_at timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS auth_sessions (
 id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
 identity_id uuid NOT NULL REFERENCES auth_identities(id) ON DELETE CASCADE,
 token_hash varchar(64) NOT NULL UNIQUE,
 expires_at timestamptz NOT NULL,
 revoked_at timestamptz,
 created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS ix_auth_sessions_identity ON auth_sessions(identity_id);
