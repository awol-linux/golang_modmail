CREATE TABLE tickets (
   id SERIAL,
   requester BIGINT NOT NULL,
   channel_id BIGINT,
   is_open BOOLEAN DEFAULT TRUE,
   PRIMARY KEY (id)
);
CREATE TABLE messages (
   sender BIGINT NOT NULL,
   ticket_id BIGINT NOT NULL,
   message_text VARCHAR(2000) NOT NULL,
   message_id BIGINT NOT NULL,
   channel_id BIGINT NOT NULL,
   PRIMARY KEY (message_id),
   CONSTRAINT ticket 
       FOREIGN KEY (ticket_id) 
           REFERENCES tickets(id)
);
