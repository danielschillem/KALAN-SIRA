CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE schools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(32) NOT NULL UNIQUE,
    name VARCHAR(200) NOT NULL,
    short_name VARCHAR(80),
    school_type VARCHAR(32) NOT NULL DEFAULT 'MIXED',
    status VARCHAR(24) NOT NULL DEFAULT 'ACTIVE',
    country_code CHAR(2) NOT NULL DEFAULT 'BF',
    city VARCHAR(120),
    address TEXT,
    phone VARCHAR(32),
    email VARCHAR(200),
    logo_url TEXT,
    currency CHAR(3) NOT NULL DEFAULT 'XOF',
    timezone VARCHAR(64) NOT NULL DEFAULT 'Africa/Ouagadougou',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE school_years (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id),
    name VARCHAR(20) NOT NULL,
    starts_at DATE NOT NULL,
    ends_at DATE NOT NULL,
    enrollment_open_at TIMESTAMPTZ,
    enrollment_close_at TIMESTAMPTZ,
    status VARCHAR(24) NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (school_id, name),
    CHECK (ends_at > starts_at)
);

CREATE TABLE levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id),
    code VARCHAR(32) NOT NULL,
    name VARCHAR(100) NOT NULL,
    cycle VARCHAR(40),
    display_order INT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (school_id, code)
);

CREATE TABLE classes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(32) NOT NULL UNIQUE,
    school_id UUID NOT NULL REFERENCES schools(id),
    school_year_id UUID NOT NULL REFERENCES school_years(id),
    level_id UUID NOT NULL REFERENCES levels(id),
    name VARCHAR(100) NOT NULL,
    capacity INT,
    status VARCHAR(24) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (school_id, school_year_id, name),
    CHECK (capacity IS NULL OR capacity > 0)
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(32),
    email VARCHAR(200),
    password_hash TEXT,
    status VARCHAR(24) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX users_email_unique ON users (lower(email)) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX users_phone_unique ON users (phone) WHERE phone IS NOT NULL;

CREATE TABLE user_school_roles (
    user_id UUID NOT NULL REFERENCES users(id),
    school_id UUID NOT NULL REFERENCES schools(id),
    role VARCHAR(32) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, school_id, role)
);

CREATE TABLE students (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(32) NOT NULL UNIQUE,
    school_id UUID NOT NULL REFERENCES schools(id),
    school_student_no VARCHAR(64),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    gender VARCHAR(16),
    birth_date DATE,
    birth_place VARCHAR(150),
    status VARCHAR(24) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX students_school_number_unique ON students(school_id, school_student_no) WHERE school_student_no IS NOT NULL;

CREATE TABLE guardians (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(32) NOT NULL UNIQUE,
    user_id UUID REFERENCES users(id),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    email VARCHAR(200),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX guardians_phone_idx ON guardians(phone);

CREATE TABLE student_guardians (
    student_id UUID NOT NULL REFERENCES students(id),
    guardian_id UUID NOT NULL REFERENCES guardians(id),
    relationship VARCHAR(32) NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    can_pay BOOLEAN NOT NULL DEFAULT TRUE,
    can_receive_notifications BOOLEAN NOT NULL DEFAULT TRUE,
    can_access_student BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY(student_id, guardian_id)
);

CREATE TABLE enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(40) NOT NULL UNIQUE,
    school_id UUID NOT NULL REFERENCES schools(id),
    school_year_id UUID NOT NULL REFERENCES school_years(id),
    student_id UUID NOT NULL REFERENCES students(id),
    class_id UUID NOT NULL REFERENCES classes(id),
    enrollment_type VARCHAR(24) NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'DRAFT',
    enrolled_at TIMESTAMPTZ,
    validated_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(student_id, school_year_id)
);

CREATE TABLE fee_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id),
    school_year_id UUID NOT NULL REFERENCES school_years(id),
    level_id UUID REFERENCES levels(id),
    class_id UUID REFERENCES classes(id),
    name VARCHAR(160) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'XOF',
    status VARCHAR(24) NOT NULL DEFAULT 'DRAFT',
    effective_from DATE,
    effective_to DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (level_id IS NOT NULL OR class_id IS NOT NULL)
);

CREATE TABLE fee_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fee_schedule_id UUID NOT NULL REFERENCES fee_schedules(id),
    code VARCHAR(40) NOT NULL,
    label VARCHAR(160) NOT NULL,
    amount BIGINT NOT NULL CHECK(amount >= 0),
    mandatory BOOLEAN NOT NULL DEFAULT TRUE,
    refundable BOOLEAN NOT NULL DEFAULT FALSE,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(fee_schedule_id, code)
);

CREATE TABLE installment_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fee_schedule_id UUID NOT NULL REFERENCES fee_schedules(id),
    name VARCHAR(160) NOT NULL,
    allow_partial_payment BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE installments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    installment_plan_id UUID NOT NULL REFERENCES installment_plans(id),
    sequence INT NOT NULL CHECK(sequence > 0),
    label VARCHAR(160) NOT NULL,
    amount BIGINT NOT NULL CHECK(amount >= 0),
    due_date DATE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(installment_plan_id, sequence)
);

CREATE TABLE student_charges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id),
    student_id UUID NOT NULL REFERENCES students(id),
    enrollment_id UUID NOT NULL REFERENCES enrollments(id),
    fee_item_id UUID REFERENCES fee_items(id),
    installment_id UUID REFERENCES installments(id),
    label VARCHAR(160) NOT NULL,
    original_amount BIGINT NOT NULL CHECK(original_amount >= 0),
    adjustment_amount BIGINT NOT NULL DEFAULT 0,
    net_amount BIGINT NOT NULL CHECK(net_amount >= 0),
    amount_paid BIGINT NOT NULL DEFAULT 0 CHECK(amount_paid >= 0),
    balance BIGINT NOT NULL CHECK(balance >= 0),
    due_date DATE,
    status VARCHAR(24) NOT NULL DEFAULT 'UPCOMING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CHECK (amount_paid <= net_amount),
    CHECK (balance = net_amount - amount_paid)
);
CREATE INDEX student_charges_enrollment_idx ON student_charges(enrollment_id, due_date);

CREATE TABLE adjustments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id),
    student_charge_id UUID NOT NULL REFERENCES student_charges(id),
    type VARCHAR(24) NOT NULL,
    amount BIGINT NOT NULL,
    reason TEXT NOT NULL,
    approved_by UUID REFERENCES users(id),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payment_intents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(48) NOT NULL UNIQUE,
    school_id UUID NOT NULL REFERENCES schools(id),
    student_id UUID NOT NULL REFERENCES students(id),
    enrollment_id UUID NOT NULL REFERENCES enrollments(id),
    guardian_id UUID REFERENCES guardians(id),
    amount BIGINT NOT NULL CHECK(amount > 0),
    currency CHAR(3) NOT NULL DEFAULT 'XOF',
    provider VARCHAR(32) NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'CREATED',
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(48) NOT NULL UNIQUE,
    school_id UUID NOT NULL REFERENCES schools(id),
    student_id UUID NOT NULL REFERENCES students(id),
    enrollment_id UUID NOT NULL REFERENCES enrollments(id),
    payment_intent_id UUID REFERENCES payment_intents(id),
    amount BIGINT NOT NULL CHECK(amount > 0),
    currency CHAR(3) NOT NULL DEFAULT 'XOF',
    payment_method VARCHAR(32) NOT NULL,
    status VARCHAR(24) NOT NULL,
    paid_at TIMESTAMPTZ,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payment_allocations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL REFERENCES payments(id),
    student_charge_id UUID NOT NULL REFERENCES student_charges(id),
    amount BIGINT NOT NULL CHECK(amount > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(payment_id, student_charge_id)
);

CREATE TABLE provider_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_intent_id UUID NOT NULL REFERENCES payment_intents(id),
    payment_id UUID REFERENCES payments(id),
    provider VARCHAR(32) NOT NULL,
    provider_reference VARCHAR(160),
    provider_status VARCHAR(64),
    request_payload JSONB,
    response_payload JSONB,
    callback_payload JSONB,
    requested_at TIMESTAMPTZ,
    confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX provider_reference_unique ON provider_transactions(provider, provider_reference) WHERE provider_reference IS NOT NULL;

CREATE TABLE receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(48) NOT NULL UNIQUE,
    school_id UUID NOT NULL REFERENCES schools(id),
    payment_id UUID NOT NULL UNIQUE REFERENCES payments(id),
    receipt_number VARCHAR(64) NOT NULL,
    pdf_url TEXT,
    verification_token_hash TEXT NOT NULL,
    issued_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(school_id, receipt_number)
);

CREATE TABLE payment_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    public_id VARCHAR(48) NOT NULL UNIQUE,
    school_id UUID NOT NULL REFERENCES schools(id),
    student_id UUID NOT NULL REFERENCES students(id),
    guardian_id UUID REFERENCES guardians(id),
    token_hash TEXT NOT NULL UNIQUE,
    amount BIGINT NOT NULL CHECK(amount > 0),
    expires_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    used_at TIMESTAMPTZ
);

CREATE TABLE notification_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id),
    event VARCHAR(48) NOT NULL,
    offset_days INT NOT NULL DEFAULT 0,
    channel VARCHAR(24) NOT NULL,
    template_key VARCHAR(80) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL REFERENCES schools(id),
    student_id UUID REFERENCES students(id),
    guardian_id UUID REFERENCES guardians(id),
    type VARCHAR(48) NOT NULL,
    channel VARCHAR(24) NOT NULL,
    template_key VARCHAR(80),
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    status VARCHAR(24) NOT NULL DEFAULT 'PENDING',
    scheduled_at TIMESTAMPTZ,
    sent_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    failed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID REFERENCES schools(id),
    user_id UUID REFERENCES users(id),
    action VARCHAR(80) NOT NULL,
    entity_type VARCHAR(80) NOT NULL,
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX audit_logs_school_created_idx ON audit_logs(school_id, created_at DESC);
