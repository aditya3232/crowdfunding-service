create table transactions
(
    id          varchar(255) not null,
    campaign_id varchar(255) not null,
    user_id     varchar(255) not null,
    amount      int not null,
    status      varchar(255) not null,
    code       varchar(255) not null,
    payment_url varchar(255) not null,
    created_at  datetime not null,
    updated_at  datetime,
    primary key (id),
    foreign key fk_transactions_campaign_id (campaign_id) references campaigns (id),
    foreign key fk_transactions_user_id (user_id) references users (id)
) engine=InnoDB;