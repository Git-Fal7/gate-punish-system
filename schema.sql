CREATE TYPE punishtype AS enum (
  'BAN',
  'MUTE'
);

CREATE TABLE punished_users (
    id SERIAL PRIMARY KEY,
    user_uuid uuid NOT NULL,
    punish_type punishtype NOT NULL,
    time_ends time NOT NULL
);

CREATE TABLE lookup_users (
  user_uuid uuid PRIMARY KEY NOT NULL,
  user_name varchar(16) NOT NULL
);