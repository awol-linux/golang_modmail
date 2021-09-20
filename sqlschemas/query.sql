-- name: GetOpenTicket :one
SELECT channel_id,
    requester,
    id
FROM tickets
WHERE requester = $1
    AND is_open = TRUE;
-- name: GetAllTickets :many
SELECT tickets.channel_id,
    tickets.requester,
    tickets.id
FROM messages,
    tickets
WHERE tickets.requester = $1
    AND messages.ticket_id = tickets.id;
-- name: AddTicket :exec
INSERT INTO tickets (requester, is_open)
VALUES ($1, TRUE);
-- name: InsertChannel :exec
UPDATE tickets
SET channel_id = $1
WHERE id = $2;
-- name: GetMessages :many
SELECT messages.sender,
    messages.ticket_id,
    messages.message_text,
    messages.message_id,
    messages.channel_id,
    tickets.channel_id,
    tickets.requester,
    tickets.id
FROM messages,
    tickets
WHERE messages.ticket_id = tickets.id
    AND ticket.id = $1;
-- name: AddMessage :exec
INSERT INTO messages (
        sender,
        ticket_id,
        message_text,
        message_id,
        channel_id
    )
VALUES(
        $1,
        $2,
        $3,
        $4,
        $5
    );
-- name: testAddIntoDb :exec
INSERT INTO messages (
        sender,
        requester,
        ticket_id,
        message_text,
        message_id,
        channel_id
    )
VALUES (
        1234,
        12345,
        12346,
        '1234567',
        12347,
        12347
    );
-- name: DropMessages :exec
DROP TABLE messages,
tickets;