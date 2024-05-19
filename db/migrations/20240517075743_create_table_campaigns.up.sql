create table campaigns
(
    id                  varchar(255) not null,
    user_id             varchar(255) not null,
    name                varchar(255) not null,
    short_description   varchar(255) not null,
    description         text not null,
    goal_amount         int not null,
    current_amount      int not null,
    perks               text not null,
    backer_count        int not null,
    slug                varchar(255) not null,
    created_at          datetime not null,
    updated_at          datetime,
    primary key (id),
    foreign key fk_campaigns_user_id (user_id) references users (id)
) engine=InnoDB;