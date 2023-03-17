alter table price_alerts
    drop constraint fk_price_alerts__user_id_on_users__id;

drop table price_alerts;

drop table users;