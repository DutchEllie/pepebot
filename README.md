Substitute the environment variables in the docker-compose.yml file for your own

When changes have happened first use `docker-compose build` to update the pepe_server container's image. Then execute `docker-compose up -d` to start the containers.

When running for the first time, the database has to first start up properly, that can take some time. Also the database `badwords` (which is automatically created) won't have any tables or content in it yet. Similarly this happens when you have no volume yet to store your data. Create these yourself.
