-- file: pre_go_acc_user_two_factor.sql

-- EnableTwoFactor
-- name: EnableTwoFactorTypeEmail :exec
INSERT INTO pre_go_acc_user_two_factor_9999 (user_id, two_factor_auth_type, two_factor_email, two_factor_auth_secret, two_factor_is_active, two_factor_created_at, two_factor_updated_at)
VALUES (?, ?, ?, "OTP", FALSE, NOW(), NOW());

-- DisableTwoFactor
-- name: DisableTwoFactor :exec
UPDATE pre_go_acc_user_two_factor_9999
SET two_factor_is_active = FALSE, 
    two_factor_updated_at = NOW()
WHERE user_id = ? AND two_factor_auth_type = ?;

-- UpdateTwoFactorStatusVerification
-- name: UpdateTwoFactorStatus :exec
UPDATE pre_go_acc_user_two_factor_9999
SET two_factor_is_active = TRUE, two_factor_updated_at = NOW()
WHERE user_id = ? AND two_factor_auth_type = ? AND two_factor_is_active = FALSE;

-- VerifyTwoFactor
-- name: VerifyTwoFactor :one
SELECT COUNT(*)
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_auth_type = ? AND two_factor_is_active = TRUE;

-- GetTwoFactorStatus
-- name: GetTwoFactorStatus :one
SELECT two_factor_is_active
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_auth_type = ?;

-- IsTwoFactorEnabled
-- name: IsTwoFactorEnabled :one
SELECT COUNT(*)
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_is_active = TRUE;

-- AddOrUpdatePhoneNumber
-- name: AddOrUpdatePhoneNumber :exec
INSERT INTO pre_go_acc_user_two_factor_9999 (user_id, two_factor_phone, two_factor_is_active)
VALUES (?, ?, TRUE)
ON DUPLICATE KEY UPDATE 
    two_factor_phone = ?, 
    two_factor_updated_at = NOW();

-- AddOrUpdateEmail
-- name: AddOrUpdateEmail :exec
INSERT INTO pre_go_acc_user_two_factor_9999 (user_id, two_factor_email, two_factor_is_active)
VALUES (?, ?, TRUE)
ON DUPLICATE KEY UPDATE 
    two_factor_email = ?, 
    two_factor_updated_at = NOW();

-- GetUserTwoFactorMethods
-- name: GetUserTwoFactorMethods :many
SELECT two_factor_id, user_id, two_factor_auth_type, two_factor_auth_secret, 
       two_factor_phone, two_factor_email, 
       two_factor_is_active, two_factor_created_at, two_factor_updated_at
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ?;

-- ReactivateTwoFactor
-- name: ReactivateTwoFactor :exec
UPDATE pre_go_acc_user_two_factor_9999
SET two_factor_is_active = TRUE, 
    two_factor_updated_at = NOW()
WHERE user_id = ? AND two_factor_auth_type = ?;

-- RemoveTwoFactor
-- name: RemoveTwoFactor :exec
DELETE FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_auth_type = ?;

-- CountActiveTwoFactorMethods
-- name: CountActiveTwoFactorMethods :one
SELECT COUNT(*)
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_is_active = TRUE;

-- GetTwoFactorMethodByID
-- name: GetTwoFactorMethodByID :one
SELECT two_factor_id, user_id, two_factor_auth_type, two_factor_auth_secret, 
       two_factor_phone, two_factor_email, 
       two_factor_is_active, two_factor_created_at, two_factor_updated_at
FROM pre_go_acc_user_two_factor_9999
WHERE two_factor_id = ?;

-- GetTwoFactorMethodByIDAndType: select lay email de sen otp
-- name: GetTwoFactorMethodByIDAndType :one
SELECT two_factor_id, user_id, two_factor_auth_type, two_factor_auth_secret, 
       two_factor_phone, two_factor_email, 
       two_factor_is_active, two_factor_created_at, two_factor_updated_at
FROM pre_go_acc_user_two_factor_9999
WHERE user_id = ? AND two_factor_auth_type = ?;
