-- name: create_citext_extension
CREATE
EXTENSION citext;

-- public.orders definition

-- Drop table

-- DROP TABLE public.orders;

-- name: create_orders_table
CREATE TABLE public.orders
(
    id               serial    NOT NULL,
    recurly_uid      varchar NULL,
    product_id       int4 NULL,
    email            varchar NULL,
    created_at       timestamp NOT NULL,
    updated_at       timestamp NOT NULL,
    subscription_uid varchar NULL,
    state            varchar NULL DEFAULT 'active':: character varying,
    CONSTRAINT orders_pkey PRIMARY KEY (id)

);

-- name: create_email_index
CREATE
INDEX index_orders_on_email ON public.orders USING btree (email);

-- name: create_orders_data
INSERT INTO orders (email, product_id, created_at, updated_at)
VALUES ('foo@example.com', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
       ('bar@example.com', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: create_users_table
CREATE TABLE public.users
(
    id                           serial  NOT NULL,
    provider                     varchar NOT NULL,
    uid                          varchar NOT NULL DEFAULT '':: character varying,
    encrypted_password           varchar NOT NULL DEFAULT '':: character varying,
    reset_password_token         varchar NULL,
    reset_password_sent_at       timestamp NULL,
    remember_created_at          timestamp NULL,
    sign_in_count                int4    NOT NULL DEFAULT 0,
    current_sign_in_at           timestamp NULL,
    last_sign_in_at              timestamp NULL,
    current_sign_in_ip           varchar NULL,
    last_sign_in_ip              varchar NULL,
    confirmation_token           varchar NULL,
    confirmed_at                 timestamp NULL,
    confirmation_sent_at         timestamp NULL,
    unconfirmed_email            varchar NULL,
    "name"                       varchar NULL,
    nickname                     varchar NULL,
    image                        varchar NULL,
    email                        citext NULL,
    tokens                       text NULL,
    created_at                   timestamp NULL,
    updated_at                   timestamp NULL,
    subscription_expired         bool NULL DEFAULT true,
    first_name                   varchar NULL,
    last_name                    varchar NULL,
    occupation                   varchar NULL,
    company                      varchar NULL,
    country                      varchar NULL,
    state                        varchar NULL,
    city                         varchar NULL,
    phone                        varchar NULL,
    address_line_1               varchar NULL,
    address_line_2               varchar NULL,
    zip_code                     varchar NULL,
    infusionsoft_contact_id      varchar NULL,
    avatar                       oid NULL,
    invoice_sent_at              timestamp NULL,
    twitter_screenname           varchar NULL,
    nyse_token                   varchar NULL,
    is_dev                       bool NULL DEFAULT false,
    signed_agreements            bool NULL DEFAULT false,
    ent_nasdaq                   bool NULL DEFAULT false,
    ent_date_nasdaq              timestamp NULL,
    ent_nyse                     bool NULL DEFAULT false,
    ent_date_nyse                timestamp NULL,
    ent_otc                      bool NULL DEFAULT false,
    ent_date_otc                 timestamp NULL,
    last_event                   varchar NULL,
    tos_agreement                bool NULL DEFAULT false,
    tos_agreement_date           timestamp NULL,
    is_pro                       bool NULL DEFAULT false,
    pro_start_date               timestamp NULL,
    pro_order_id                 varchar NULL,
    recurly_account_code         varchar NULL,
    compliance_checked           bool    NOT NULL DEFAULT false,
    compliance_check_date        timestamp NULL,
    is_compliant                 bool NULL,
    is_suspended                 bool    NOT NULL DEFAULT false,
    suspended_date               timestamp NULL,
    "comment"                    text NULL,
    affiliate_plan_id            int4 NULL,
    data_plan                    varchar NULL,
    has_grace_period             bool    NOT NULL DEFAULT false,
    grace_period_end_date        timestamp NULL,
    level2_subscriptions_active  bool    NOT NULL DEFAULT false,
    created_by_admin             bool    NOT NULL DEFAULT false,
    read_message_ids             _int4 NULL       DEFAULT '{}':: integer [],
    recurly_status               jsonb NULL DEFAULT '{}'::jsonb,
    canceled_addons              _jsonb NULL      DEFAULT '{}'::jsonb[],
    setting_file_name            varchar NULL,
    setting_content_type         varchar NULL,
    setting_file_size            int4 NULL,
    setting_updated_at           timestamp NULL,
    profile_picture_file_name    varchar NULL,
    profile_picture_content_type varchar NULL,
    profile_picture_file_size    int4 NULL,
    profile_picture_updated_at   timestamp NULL,
    trial_agreement_accepted     bool NULL DEFAULT false,
    trial_agreement_date         timestamp NULL,
    log_file_name                varchar NULL,
    log_content_type             varchar NULL,
    log_file_size                int4 NULL,
    log_updated_at               timestamp NULL,
    backup_setting_file_name     varchar NULL,
    backup_setting_content_type  varchar NULL,
    backup_setting_file_size     int4 NULL,
    backup_setting_updated_at    timestamp NULL,
    wp_id                        int4 NULL,
    infusionsoft_order_id        jsonb NULL DEFAULT '{}'::jsonb,
    utm_source                   text NULL,
    utm_content                  text NULL,
    utm_medium                   text NULL,
    utm_term                     text NULL,
    utm_campaign                 text NULL,
    map_code                     text NULL,
    sms_number                   varchar NULL,
    sms_phone_number             varchar NULL,
    recurly_addons               _varchar NULL    DEFAULT '{}':: character varying [],
    created_by_api               bool    NOT NULL DEFAULT false,
    external_account             varchar NULL,
    offer_mapping_id             int4 NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

--name: create_users_data
INSERT INTO users (provider, email)
VALUES ('email', 'foo@example.com'),
       ('email', 'bar@example.com'),
       ('email', 'fizz@example.com');

--name: create_entitlement_logs_table
CREATE TABLE public.entitlement_logs
(
    id              serial    NOT NULL,
    ent_nasdaq      bool NULL DEFAULT false,
    ent_date_nasdaq timestamp NULL,
    ent_nyse        bool NULL DEFAULT false,
    ent_date_nyse   timestamp NULL,
    ent_otc         bool NULL DEFAULT false,
    ent_date_otc    timestamp NULL,
    created_at      timestamp NOT NULL,
    updated_at      timestamp NOT NULL,
    user_id         int4 NULL,
    "event"         varchar NULL,
    CONSTRAINT entitlement_logs_pkey PRIMARY KEY (id)
);

--name: create_nyse_entries_table
CREATE TABLE public.nyse_entries
(
    id                          serial    NOT NULL,
    name_of_vendor              varchar NULL,
    section1_agreed             bool NULL,
    section2_12_a               bool NULL,
    section2_12_b               bool NULL,
    section2_12_c               bool NULL,
    section2_12_d               bool NULL,
    section2_12_e               bool NULL,
    section2_12_f               bool NULL,
    section2_12_g               bool NULL,
    section2_12_h               bool NULL,
    section2_12_i               bool NULL,
    section2_12_j               bool NULL,
    section2_12_k               bool NULL,
    section2_12_certificated    bool NULL,
    section2_agreed             bool NULL,
    completed                   bool NULL,
    user_id                     int4 NULL,
    created_at                  timestamp NOT NULL,
    updated_at                  timestamp NOT NULL,
    nasdaq_agreed               bool NULL,
    is_non_pro                  bool NULL,
    personal_first_name         text NULL,
    personal_last_name          text NULL,
    personal_middle_name        text NULL,
    personal_address            text NULL,
    personal_city               text NULL,
    personal_state              text NULL,
    personal_country            text NULL,
    personal_zip                text NULL,
    personal_email              text NULL,
    employer_occupations        text NULL,
    employer_retiree_name       text NULL,
    employer_retiree_address    text NULL,
    employer_retiree_city       text NULL,
    employer_retiree_state_prov text NULL,
    employer_retiree_country    text NULL,
    employer_retiree_zip        text NULL,
    employer_title              text NULL,
    employer_functions          text NULL,
    nyse_agreed                 bool NULL,
    otc_agreed                  bool NULL,
    signature_agreed            bool NULL,
    completed_at                timestamp NULL,
    void_agreement              bool NULL DEFAULT false,
    void_agreement_date         timestamp NULL,
    CONSTRAINT nyse_entries_pkey PRIMARY KEY (id)
);

-- name: create_nyse_entries_data
INSERT INTO nyse_entries (user_id, created_at, updated_at)
VALUES ((SELECT id FROM users WHERE email like 'bar@example.com'), CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

--name: create_referrals_table
CREATE TABLE public.referrals
(
    id                      serial    NOT NULL,
    user_id                 int4 NULL,
    affiliate_id            int4 NULL,
    subscription_price      float8 NULL,
    subscription_event_date timestamp NULL,
    subscription_start_date timestamp NULL,
    trial_period            int4 NULL,
    is_cancelled            bool NULL,
    is_expired              bool NULL,
    affiliate_code          varchar NULL,
    commission_value        float8 NULL,
    commission_type         varchar NULL,
    commission_recurrence   varchar NULL,
    commission_total        float8 NULL,
    created_at              timestamp NOT NULL,
    updated_at              timestamp NOT NULL,
    date_cancelled          timestamp NULL,
    date_expired            timestamp NULL,
    subscription_id         varchar NULL,
    plan_code               varchar NULL,
    one_time                bool      NOT NULL DEFAULT false,
    CONSTRAINT referrals_pkey PRIMARY KEY (id)
);

--name: create_referrals_data
INSERT INTO referrals (user_id, is_cancelled, is_expired, created_at, updated_at)
VALUES ((SELECT id FROM users WHERE email LIKE 'fizz@example.com'), FALSE, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

--name: create_recurly_subs_table
CREATE TABLE public.recurly_subs
(
    id                      serial NOT NULL,
    user_id                 int4 NULL,
    recurly_subscription_id varchar NULL,
    utm_source              varchar NULL,
    utm_content             varchar NULL,
    utm_medium              varchar NULL,
    utm_term                varchar NULL,
    utm_campaign            varchar NULL,
    CONSTRAINT recurly_subs_pkey PRIMARY KEY (id)
);

--name: create_recurly_subs_data
INSERT INTO public.recurly_subs (user_id, recurly_subscription_id, utm_source, utm_content, utm_medium, utm_term,
                                 utm_campaign)
VALUES (3805, '4f5d7e1db40b80b1aa319a418a8fe9de', NULL, NULL, NULL, NULL, NULL);