
DROP TABLE IF EXISTS user; /* Users. */
CREATE TABLE if not exists user(
        pk integer primary key,
        email text,
        name text,
        password text,
        uuid text

);

create  index idx_user_email on user(email);
create  index idx_user_uuid on user(uuid);
/*

 */
 DROP TABLE IF EXISTS certificate;

CREATE TABLE IF NOT EXISTS certificate(
    pk integer primary key,
    user_email text null,
    subject_common_name text,
    subject_organization text null,
    subject_organization_unit text null,
    subject_country text null,
    subject_locality text null,
    subject_province text null,

    issuer_id integer,

    serial_number text,
    not_before text,
    not_after text,
    FOREIGN KEY (user_email) REFERENCES user(email) ON DELETE CASCADE
);

create index idx_certificate on certificate(subject_common_name);
create index idx_certificate_user_email on certificate(user_email);
/*

 */
