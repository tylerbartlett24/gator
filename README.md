# Gator

A multi-user CLI tool for aggrgating RSS feeds and saving posts to a local database. The posts can then be browsed at the user's leisure.

## Installation

Make sure you have the [Go toolchain]((https://golang.org/dl/)) installed. Then ensure you have a local Postgres database. If both requirements are met, you can install gator with:

```bash
go install ...
```

Goose should also be installed for database management.

## Configuration

In your home directory create a .gatorconfig.json file. The file should have 2 fields, "db_url", which is a url to the local Postgres database, and "current_user_name", which will be the username of whoever is presently logged in (this can be left blank initiallly). Change into the sql.schema directory and run a goose upmigration to initialize the database.

## Commmands

Gator supports the following commands:

- login <username>: makes username the current user
- register <username>: adds username to the list of users in the database.
- reset: removes all registered users from the database, along with all their foolowed feeds and saved posts.
- users: lists all registered users
- agg <duration between requests>: causes gator to continually scrape the current users followed feeds starting from the least recently updated, with one request per span of time given by the argument. Be careful not to choose to brief a duration, for fear of DOSing the target. Scraped posts are saved to the database. Use ctrl+c to exit.
- addfeed <feed name> <feed url>: adds the given feed to the database and automatically causes the current user to follow it.
- feeds: lists all feeds in the database.
- follow <feed url>: causes the current user to follow the given feed.
- unfollow <feed url>: causes the current user to unfollow the given feed.
- following: lists all feeds followed by the current user
- browse <optional limit>: retrieves the most recent posts from the feeds followed by the current user, up to the given limit (default 2 posts).

