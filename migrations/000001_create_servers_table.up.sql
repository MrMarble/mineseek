CREATE TABLE IF NOT EXISTS servers(
   host VARCHAR(254) NOT NULL,
   port integer NOT NULL,
   version VARCHAR(254) NOT NULL,
   favicon text,
   motd text,
   max_players integer NOT NULL,
   PRIMARY KEY(host, port)
);