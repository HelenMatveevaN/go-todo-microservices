-- +goose Down
-- +goose StatementBegin
drop table if exists tasks;
-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
create table tasks (
	id 			SERIAL	primary key,
	title 		text	not null,
	content		text,
	is_done		boolean	default false,
	created_at	timestamp	default now()
);
-- +goose StatementEnd
