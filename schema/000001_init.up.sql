-- create database price_bot;

create table users
(
    id                     serial primary key,
    chat_id                bigint  not null,
    user_name              varchar,
    first_name             varchar,
    last_name              varchar,
    description            varchar,
    chat_type              varchar not null default 'private',
    photo                  varchar,
    title                  varchar,
    all_members_are_admins boolean,
    invite_link            varchar,
    created_at             timestamp with time zone,
    updated_at             timestamp with time zone
);


create table price_alerts
(
    id           serial primary key,
    user_id      int     not null,
    ticker       varchar not null,
    create_price decimal not null,
    alert_price  decimal not null,
    created_at   timestamp with time zone,

    constraint fk_price_alerts__user_id_on_users__id
        foreign key (user_id)
            references users
);
