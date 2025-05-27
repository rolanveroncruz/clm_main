drop table user;
CREATE TABLE if not exists user(
        id integer primary key,
        email text,
        name text,
        password text

);

create  index idx_user_email on user(email);
