# Price Alerts Telegram Bot

Price Alerts Telegram Bot is an open-source project for creating a Telegram bot that tracks prices on exchanges and notifies users when threshold values are reached. You are welcome to contribute to the project in any way you see fit.

## Project Structure

```
.
├── cmd
│   ├── price_alerts_bot
│   │   └── main.go               # Launches the bot that communicates with Telegram.
│   └── worker
│       └── price_monitoring.go   # A worker that tracks exchange prices and compares them to the values stored in the database.
└── internal
    ├── app                       # Contains everything needed to launch services.
    ├── config
    │   └── config.go             # A configuration file for parsing the .env file.
    ├── models                    # Contains the database models.
    ├── providers                 # Creates various providers for tracking prices. Currently, only the Kucoin exchange provider is available.
    ├── telegram                  # A package for working with Telegram.
    └── price_alerts              # A package for working with the price_alerts model. This table stores all the prices tracked by users.
└── schema                        # Contains all migrations for the migrate package.
```

The provider selection works through the following function:

```go
func GetProvider(name string) (IProvider, error) {
	switch name {
	case "kucoin":
		p := &KucoinProvider{}
		p.SetParams()
		return p, nil
    //case "someProvider":
	//    p := &SomeProvider{}
    //    p.SetParams()
    //    return p, nil
	default:
		return nil, errors.New("the provider is not found")
	}
}
```

There are currently two migrations that create two tables:

```sql
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

alter table price_alerts
    add observable_price numeric not null default 0;

alter table price_alerts
    add observable_percent numeric not null default 0;
```

When a user joins the bot, they are recorded in the users table with their chat ID, which will be used to send push notifications.

The price_alerts table stores all requests for price tracking. Below is a description of each field:

```
id: Unique identifier for the price alert.
user_id: Foreign key referencing the user who created the alert.
ticker: Ticker symbol of the asset being tracked.
create_price: Price at which the alert was created.
alert_price: Price at which the user will be notified.
created_at: Timestamp when the alert was created.
observable_price: The observed price of the asset.
observable_percent: The observed percentage change of the asset's price.
```

## How get `TG_TOKEN` for .env
To obtain a TG_TOKEN for your Telegram bot, follow these steps:
1. Open the Telegram app and search for the BotFather bot. You can find it by searching for @BotFather in the search bar.
2. Start a chat with the BotFather by clicking on the "Start" button.
3. Create a new bot by sending the `/newbot` command to the BotFather. You will receive a message asking you to choose a name for your bot.
4. Enter a name for your bot `(e.g., "Price Alerts Bot")`. This is a display name and can contain spaces and special characters. The BotFather will then ask you to choose a username for your bot.
5. Choose a unique username for your bot that ends in bot `(e.g., "price_alerts_bot")`. Usernames must be unique and can only contain lowercase letters, numbers, and underscores. If the username is available, the BotFather will create your bot and provide you with a token.
6. Copy the token `(e.g., 123456789:ABCdefGHIJKLMNOPQrstUVWxyz)`. This is your TG_TOKEN, and you should paste it into the .env file, as described in the previous steps:
```dotenv
TG_TOKEN=your_token_from_botfather
```

## Run Project
To launch this project, follow the steps below:

1. Install the migrate tool if you haven't already. You can find the installation instructions [here](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).
2. Copy the provided .env-example file to a new .env file in the root directory of your project:
```shell
$ cp .env-example .env
```
Fill in the appropriate configuration settings in the .env file, such as the database connection string, Telegram bot token, and Kucoin API credentials:
```dotenv
TG_TOKEN=your_token_from_botfather

KUCOIN_API_URL=https://api.kucoin.com
KUCOIN_API_KEY=your_api_key_for_kucoin
KUCOIN_API_SECRET=your_api_secret_for_kucoin

DATABASE_URL=url_for_db
```
3. Run the migrations using the migrate tool. Replace path/to/schema with the actual path to the schema folder in your project, and your_database_connection_string with your database connection string from the .env file:
```shell
migrate -source file:schema -database 'your_database_connection_string' up
```
4. Use the provided Makefile to build and run the project:
   - To build the project, run: 
        ```shell
        $ make build
        ```
   - To run the project in development mode, run:
        ```shell
        $ make dev
        ```
   - To run the worker in development mode, run:
        ```shell
        $ make worker-dev
        ```
   - To run the project in production mode, first build the project using make build, and then run:
        ```shell
        $ make prod
        ```
That's it! The project should be up and running, with the database migrations applied.




## Contributing

Feel free to contribute to the project by opening issues, submitting pull requests, or suggesting new features and improvements. Your contributions are greatly appreciated.
