alter table price_alerts
    add observable_price numeric not null default 0;

alter table price_alerts
    add observable_percent numeric not null default 0;