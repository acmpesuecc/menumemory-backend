-- name: GetRestaurantsLike :many
SELECT id, name, area, address, mapslink, mapsratingoutof5 FROM Restaurant WHERE name LIKE ?;

-- name: GetRestaurantHistory :many
SELECT id, date, time FROM Visit
    WHERE userId = ? AND restaurantId = ?;

-- name: GetOrdersForVisit :many
SELECT d.name, o.rating, o.reviewText FROM
    Orders o JOIN Dish d ON o.dishId = d.id
    WHERE o.visitId = ?;

-- name: CreateVisit :exec
INSERT INTO Visit(date, time, userId, restaurantId)
VALUES (?, ?, ?, ?);

-- name: CreateOrder :exec
INSERT INTO Orders(visitId, dishId, rating, reviewText)
VALUES (?, ?, ?, ?);

-- name: GetVisitById :one
SELECT id, date, time, userId, restaurantId FROM Visit WHERE id = ?;

-- name: UpdateVisit :exec
UPDATE Visit
SET date = ?, time = ?, restaurantId = ?
WHERE id = ? AND userId = ?;
