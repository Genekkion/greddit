# Postgres

This folder contains Postgres-specific code.

## SQL files

In the `sql` folder, you can find SQL files which include the schema
for all tables required. Note that it uses a simple migration system
based on the numerical prefix of the file name. Should more migrations
be required, they should be added in the correct order, using the
same format.