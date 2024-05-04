-- name: CreatePunishUsersTable :exec
CREATE TABLE IF NOT EXISTS punished_users (
    id SERIAL PRIMARY KEY,
    user_uuid uuid NOT NULL,
    punish_type punishtype NOT NULL,
    time_ends time NOT NULL
);

-- name: CreateLookupUserTable :exec
CREATE TABLE IF NOT EXISTS lookup_users (
    user_uuid uuid PRIMARY KEY NOT NULL,
    user_name varchar(16) NOT NULL
);

-- name: CreatePunishType :exec
DO $$ BEGIN
    IF to_regtype('punishtype') IS NULL THEN
        CREATE TYPE friendstatus AS enum (
            'BAN',
            'MUTE'
        );
    END IF;
END $$;

-- name: LogIntoLookupTable :exec
INSERT INTO lookup_users (
    user_uuid, user_name
) VALUES (
    $1, $2
)
ON CONFLICT(user_uuid)
DO UPDATE SET
user_name = $2;

-- name: PunishPlayer :exec
INSERT INTO punished_users (
    user_uuid, punish_type, time_ends
) VALUES (
    $1, $2, $3
);
