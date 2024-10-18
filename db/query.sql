-- name: GetRestaurantsLike :many
SELECT * FROM Restaurant where Name like ?;

-- name: GetRestaurantHistory :many
SELECT id, Date, Time from Visit
    where UserId = ? and RestaurantId = ?;

-- name: GetOrdersForVisit :many
SELECT d.Name, o.Rating, o.ReviewText from
    Orders o join Dish d on o.DishId = d.id
    where o.VisitId = ?;

-- name: CreateVisit :exec
INSERT INTO Visit(Date, Time, UserId, RestaurantId)
VALUES (?, ?, ?, ?);

-- name: CreateOrder :exec
INSERT INTO Orders(VisitId, DishId, Rating, ReviewText)
VALUES (?, ?, ?, ?);

-- name: GetVisitOwnerByID :one
SELECT UserId FROM Visit WHERE id = ?;

-- name: DeleteVisitByID :exec
DELETE FROM Visit WHERE id = ?;
