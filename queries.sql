-- name: CreatePunishUsersTable :exec
CREATE TABLE IF NOT EXISTS punished_users (
    id SERIAL PRIMARY KEY,
    user_uuid uuid NOT NULL,
    reason text NOT NULL,
    done_by varchar(16) NOT NULL,
    punish_type punishtype NOT NULL,
    time_ends timestamptz NOT NULL,
    createdat timestamptz NOT NULL
);

-- name: CreateLookupUserTable :exec
CREATE TABLE IF NOT EXISTS lookup_users (
    user_uuid uuid PRIMARY KEY NOT NULL,
    user_name varchar(16) NOT NULL
);

-- name: CreatePunishType :exec
DO $$ BEGIN
    IF to_regtype('punishtype') IS NULL THEN
        CREATE TYPE punishtype AS enum (
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
    user_uuid, reason, done_by, punish_type, time_ends, createdat
) VALUES (
    $1, $2, $3, $4, $5, NOW()
);

-- name: GetPlayerUUID :one
SELECT user_uuid FROM lookup_users
WHERE LOWER(user_name) = LOWER($1);

-- name: IsPunishedPlayer :one
SELECT * FROM punished_users
WHERE user_uuid = $1 AND punish_type = $2 AND time_ends > NOW()
ORDER BY time_ends DESC LIMIT 1;

-- name: UnpunishPlayer :exec
UPDATE punished_users
SET time_ends = TIMESTAMP '2000-01-01 00:00:00' AT TIME ZONE 'UTC'
WHERE user_uuid = $1 AND punish_type = $2 AND time_ends > NOW();