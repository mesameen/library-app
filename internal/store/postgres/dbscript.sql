create table books (
	id SERIAL PRIMARY KEY,
	title VARCHAR(255) NOT NULL UNIQUE,
	available_copies INT NOT NULL CHECK (available_copies >= 0)
)

INSERT INTO books (title, available_copies) VALUES ('Book_1', 3);
INSERT INTO books (title, available_copies) VALUES ('Book_2', 4);
INSERT INTO books (title, available_copies) VALUES ('Book_3', 7);
INSERT INTO books (title, available_copies) VALUES ('Book_4', 5);
INSERT INTO books (title, available_copies) VALUES ('Book_5', 10);

select * from  books;

create table loans (
	id SERIAL PRIMARY KEY,
	title VARCHAR(256) NOT NULL,
	name_of_borrower VARCHAR(256) NOT NULL,
	loan_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	return_date TIMESTAMP NOT NULL,
	status VARCHAR(100) NOT NULL
)

select * from loans;
