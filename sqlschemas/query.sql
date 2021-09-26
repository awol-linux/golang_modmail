-- name: GetOpenTicket :one
SELECT ticket_channel_id,
    requester,
    id
FROM tickets
WHERE requester = $1
    AND is_open = TRUE;
-- name: GetAllTickets :many
SELECT tickets.ticket_channel_id,
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
SET ticket_channel_id = $1
WHERE id = $2;
-- name: CloseTicket :exec
UPDATE tickets
SET is_open = FALSE
WHERE ticket_channel_id = $1;
-- name: GetMessages :many
SELECT MESSAGES.SENDER,
    MESSAGES.TICKET_ID,
    MESSAGES.MESSAGE_TEXT,
    MESSAGES.MESSAGE_ID,
    MESSAGES.CHANNEL_ID,
    TICKETS.ticket_channel_id,
    TICKETS.REQUESTER,
    TICKETS.ID,
    FORWARDED.sendto_channel_id,
    FORWARDED.sendto_message_id
FROM MESSAGES,
    TICKETS,
    FORWARDED
WHERE MESSAGES.TICKET_ID = TICKETS.ID
    AND FORWARDED.sendto_message_id = MESSAGES.FORWARDED
    AND ticket.id = $1;
-- name: GetMessage :one
SELECT MESSAGES.SENDER,
    MESSAGES.TICKET_ID,
    MESSAGES.MESSAGE_TEXT,
    MESSAGES.MESSAGE_ID,
    MESSAGES.CHANNEL_ID,
    TICKETS.ticket_channel_id,
    TICKETS.REQUESTER,
    TICKETS.ID,
    FORWARDED.sendto_channel_id,
    FORWARDED.sendto_message_id
FROM MESSAGES,
    TICKETS,
    FORWARDED
WHERE MESSAGES.TICKET_ID = TICKETS.ID
    AND FORWARDED.sendto_message_id = MESSAGES.FORWARDED
    AND (
        messages.message_id = $1
        OR forwarded.sendto_message_id = $1
    );
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
-- name: InsertForward :exec
INSERT INTO forwarded (sendto_message_id, sendto_channel_id)
VALUES ($1, $2);
-- name: LinkForward :exec
UPDATE messages
SET forwarded = $1
WHERE message_id = $2;
-- name: GetForwarded :one
SELECT forwarded.sendto_channel_id,
    forwarded.sendto_message_id
FROM forwarded,
    messages
WHERE messages.message_id = $1;
-- name: DropMessages :exec
DROP TABLE messages,
tickets,
forwarded;