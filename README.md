A Front End service, that just displays web pages;

An Authentication service, with a Postgres database;

A Logging service, with a MongoDB database;

A Listener service, which receives messages from RabbitMQ and acts upon them;

A Broker service, which is an optional single point of entry into the microservice cluster;

A Mail service, which takes a JSON payload, converts into a formatted email, and send it out.