create table campaign_images
(
    id          varchar(255) not null,
    campaign_id varchar(255) not null,
    file_name   varchar(255) not null,
    is_primary  tinyint(1) not null,
    created_at  datetime not null,
    updated_at  datetime,
    primary key (id),
    foreign key fk_campaign_images_campaign_id (campaign_id) references campaigns (id) on delete cascade on update restrict
) engine=InnoDB;