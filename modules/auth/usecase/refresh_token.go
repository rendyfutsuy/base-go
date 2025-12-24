package usecase

// RefreshToken generates a new access token based on the provided access token.
// If the old token is revoked (not found in Redis), it returns an error.
//
// Parameters:
// - ctx: The context for managing request lifecycle and cancellation.
// - accessToken: The existing access token to be refreshed.
//
// Returns:
// - string: The new access token.
// - error: An error if the token refresh fails or if the token is revoked.
// func (u *authUsecase) RefreshToken(ctx context.Context, accessToken string) (string, error) {
// 	// Check if token exists in Redis (not revoked)
// 	user, err := u.authRepo.GetUserByAccessToken(ctx, accessToken)
// 	if err != nil {
// 		if err == redis.Nil {
// 			return "", constants.ErrTokenRevoked
// 		}
// 		return "", err
// 	}

// 	// If user not found, token is revoked
// 	if user.ID == uuid.Nil {
// 		return "", constants.ErrTokenRevoked
// 	}

// 	// Destroy old token
// 	if err := u.authRepo.DestroyToken(ctx, accessToken); err != nil {
// 		utils.Logger.Error("failed to destroy old token", zap.Error(err))
// 	}

// 	// Generate new JWT token with same logic as Authenticate
// 	// Get Redis TTL from config (for session storage)
// 	ttlSeconds := utils.ConfigVars.Int("auth.redis_ttl_seconds")
// 	if ttlSeconds <= 0 {
// 		ttlSeconds = int((48 * time.Hour).Seconds()) // default 2 days
// 	}

// 	// Get JWT expiration time from config (for token expiration)
// 	jwtExpiresAtSeconds := utils.ConfigVars.Int("auth.access_token_ttl_seconds")
// 	if jwtExpiresAtSeconds <= 0 {
// 		jwtExpiresAtSeconds = 30 * 60 // Default 30 minutes
// 	}

// 	now := time.Now().UTC()

// 	// Generate new JWT token
// 	// JWT expires based on access_token_ttl_seconds (30 minutes)
// 	// Redis session TTL is based on redis_ttl_seconds (2 days)
// 	claims := AuthClaims{
// 		UserID: user.ID.String(),
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(jwtExpiresAtSeconds) * time.Second)),
// 			IssuedAt:  jwt.NewNumericDate(now),
// 			ID:        uuid.NewString(),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	newAccessToken, err := token.SignedString(u.signingKey)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to sign token: %w", err) // If fail to sign, return error.
// 	}

// 	// Record new access token to Redis
// 	if err := u.authRepo.AddUserAccessToken(ctx, newAccessToken, user.ID); err != nil {
// 		return "", fmt.Errorf("failed to store new access token: %w", err) // If fail to record, return error.
// 	}

// 	return newAccessToken, nil
// }
