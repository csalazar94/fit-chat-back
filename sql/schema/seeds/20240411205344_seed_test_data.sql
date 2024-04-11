-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, full_name, email, encoded_hash, created_at, updated_at)
VALUES ('2a4fbca3-fb52-4e0f-9c31-83eb28f1c651', 'Anjela Riglesford', 'a.riglesford@example.com', '$argon2id$v=19$m=65536,t=1,p=32$Ff0YFbFnIpX/gj2V0HMLYHfGDUqboWI11AEGzd0q6Uw$A4HJf4OSGGF7psh8OX/ZS26jw2t9h/EyOoWUenzTvvKDJXzplNLeYZ84McaVU0q3s/zPvP9lugJykrkCNQIB3L3hYl6FysQuQXBs4Z0bslKKSlmIKPC2go7W9ZTxBiGKt43ZpNGqDBQB6Z0UPPS301fpiBPsrNtoypcqR9RrJagpcb5ssMyKnrt9CwvrurWmhhrdnJAqfOa8SxeXuEXb/aLeaR0VvEgYhomkNBLFeRmyYyKa1CFeTVIMm6r/6QgMrUHXQjfHJpxBli1ewvhmGyKiXrAgQXzCnS6nFU1SjyjFRplwQbeWbhR9/P3lAw43dUqoMYUSU8IKcXksS0/t4g', timezone('utc', now()), timezone('utc', now()));

INSERT INTO chats (id, user_id, title, created_at, updated_at)
VALUES ('ced2a910-cdf3-4763-9eca-61337bad8ee2', '2a4fbca3-fb52-4e0f-9c31-83eb28f1c651', 'Chat', timezone('utc', now()), timezone('utc', now()));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users
WHERE id = '2a4fbca3-fb52-4e0f-9c31-83eb28f1c651';

DELETE FROM chats
WHERE id = 'ced2a910-cdf3-4763-9eca-61337bad8ee2';
-- +goose StatementEnd
