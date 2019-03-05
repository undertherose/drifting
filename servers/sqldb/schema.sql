create table if not exists users (
    id int primary key not null auto_increment,
    userName varchar(32) unique,
    passHash binary(60),
    userType varchar(32),
    isSuspended bool not null default false
);
