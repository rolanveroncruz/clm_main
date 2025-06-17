

 DROP TABLE IF EXISTS certificate;

CREATE TABLE IF NOT EXISTS certificate(
    pk integer primary key,
    sql_time_stamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_email text null,
    subject_common_name text,
    subject_organization text null,
    subject_organizational_unit text null,
    subject_country text null,
    subject_locality text null,
    subject_province text null,

    issuer_common_name text,
    issuer_organization text,
    issuer_organizational_unit text,
    issuer_country text,

    serial_number text,
    not_before text,
    not_after text,
    requested_server text null, -- this is the name of the server requested.
    FOREIGN KEY (user_email) REFERENCES user(email) ON DELETE CASCADE
);

create index idx_certificate on certificate(subject_common_name);
create index idx_issuer on certificate(issuer_common_name);
create index idx_certificate_user_email on certificate(user_email);
/*

 */
