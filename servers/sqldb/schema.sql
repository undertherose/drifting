create table if not exists users (
    id int primary key not null auto_increment,
    email  varchar(254) unique,
    passHash binary(60),
    userName varchar(32) unique,
    firstName varchar(255),
    lastName varchar(255),
    photoURL varchar(64)
);