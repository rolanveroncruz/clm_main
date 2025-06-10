
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

