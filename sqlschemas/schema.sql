CREATE TABLE tickets (
   id SERIAL,
   requester BIGINT NOT NULL,
   ticket_channel_id BIGINT,
   is_open BOOLEAN DEFAULT TRUE,
   PRIMARY KEY (id)
);

CREATE TABLE forwarded (
   sendto_message_id BIGINT NOT NULL,
   sendto_channel_id BIGINT NOT NULL,
   PRIMARY KEY (sendto_message_id)
);

CREATE TABLE messages (
   sender BIGINT NOT NULL,
   ticket_id BIGINT NOT NULL REFERENCES tickets (id),
   message_text VARCHAR(2000) NOT NULL,
   message_id BIGINT NOT NULL,
   channel_id BIGINT NOT NULL,
   forwarded BIGINT REFERENCES forwarded (sendto_message_id),
   PRIMARY KEY (message_id)
);
