create table users 
(
    id          varchar(255) not null,
    name        varchar(255) not null,
    occupation varchar(255) not null,
    email varchar(255) not null,
    password_hash varchar(255) not null,
    avatar_file_name varchar(255),
    role varchar(255) not null,
    token varchar(255),
    created_at datetime not null,
    updated_at datetime,
    primary key (id)
) engine=InnoDB;