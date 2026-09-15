-- Script d'initialisation de la base de données Roundcube
-- Ce script crée les tables nécessaires pour Roundcube

-- Connecter à la base de données office1789
\c office1789

-- Créer les tables Roundcube
CREATE TABLE IF NOT EXISTS roundcube_users (
    user_id serial PRIMARY KEY,
    username varchar(128) NOT NULL,
    mail_host varchar(128) NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    last_login timestamp with time zone DEFAULT NULL,
    failed_login timestamp with time zone DEFAULT NULL,
    failed_login_counter integer DEFAULT 0,
    language varchar(16),
    preferences text DEFAULT ''::text,
    CONSTRAINT roundcube_users_username_key UNIQUE (username, mail_host)
);

CREATE TABLE IF NOT EXISTS roundcube_session (
    sess_id varchar(128) NOT NULL PRIMARY KEY,
    changed timestamp with time zone DEFAULT now() NOT NULL,
    ip varchar(40) NOT NULL,
    vars text NOT NULL
);

CREATE TABLE IF NOT EXISTS roundcube_identities (
    identity_id serial PRIMARY KEY,
    user_id integer NOT NULL REFERENCES roundcube_users(user_id) ON DELETE CASCADE,
    changed timestamp with time zone DEFAULT now() NOT NULL,
    del smallint DEFAULT 0 NOT NULL,
    standard smallint DEFAULT 0 NOT NULL,
    name varchar(128),
    organization varchar(128),
    email varchar(128) NOT NULL,
    reply_to varchar(128),
    bcc varchar(128),
    signature text,
    html_signature smallint DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS roundcube_contacts (
    contact_id serial PRIMARY KEY,
    changed timestamp with time zone DEFAULT now() NOT NULL,
    del smallint DEFAULT 0 NOT NULL,
    name varchar(128) DEFAULT ''::varchar NOT NULL,
    email text DEFAULT ''::text NOT NULL,
    firstname varchar(128) DEFAULT ''::varchar NOT NULL,
    surname varchar(128) DEFAULT ''::varchar NOT NULL,
    vcard text,
    words text,
    user_id integer NOT NULL REFERENCES roundcube_users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS roundcube_contactgroups (
    contactgroup_id serial PRIMARY KEY,
    user_id integer NOT NULL REFERENCES roundcube_users(user_id) ON DELETE CASCADE,
    changed timestamp with time zone DEFAULT now() NOT NULL,
    del smallint DEFAULT 0 NOT NULL,
    name varchar(128) NOT NULL
);

CREATE TABLE IF NOT EXISTS roundcube_cache (
    user_id integer NOT NULL REFERENCES roundcube_users(user_id) ON DELETE CASCADE,
    cache_key varchar(128) NOT NULL,
    created timestamp with time zone DEFAULT now() NOT NULL,
    expires timestamp with time zone DEFAULT NULL,
    data text NOT NULL,
    PRIMARY KEY (user_id, cache_key)
);

CREATE TABLE IF NOT EXISTS roundcube_cache_index (
    user_id integer NOT NULL REFERENCES roundcube_users(user_id) ON DELETE CASCADE,
    mailbox varchar(255) NOT NULL,
    changed timestamp with time zone DEFAULT now() NOT NULL,
    valid smallint DEFAULT 0 NOT NULL,
    data text NOT NULL,
    PRIMARY KEY (user_id, mailbox)
);

CREATE TABLE IF NOT EXISTS roundcube_cache_messages (
    user_id integer NOT NULL REFERENCES roundcube_users(user_id) ON DELETE CASCADE,
    mailbox varchar(255) NOT NULL,
    uid integer NOT NULL,
    changed timestamp with time zone DEFAULT now() NOT NULL,
    data text NOT NULL,
    flags integer DEFAULT 0 NOT NULL,
    PRIMARY KEY (user_id, mailbox, uid)
);

-- Accorder les permissions
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO robespierre;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO robespierre;

-- Afficher un message de succès
\echo 'Tables Roundcube créées avec succès!'
