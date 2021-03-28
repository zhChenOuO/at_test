-- +goose Up
INSERT INTO
    "public"."identity_accounts"(
        "id",
        "name",
        "email",
        "phone",
        "phone_area_code",
        "password",
        "phone_verify_status",
        "email_verify_status",
        "created_at",
        "updated_at"
    )
VALUES
    (
        10000,
        'test name',
        'test already exists email',
        '',
        '',
        'daHWp25ImHWZ2OMxwI0mCw==',
        0,
        0,
        '2020-10-24 10:31:20.231808+00',
        '2020-10-24 10:31:20.231808+00'
    ),
    (
        10001,
        'test name',
        '',
        'test already exists phone',
        'test already exists phone area code',
        'daHWp25ImHWZ2OMxwI0mCw==',
        0,
        0,
        '2020-10-24 10:31:20.231808+00',
        '2020-10-24 10:31:20.231808+00'
    );