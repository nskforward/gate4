package finam

import "errors"

func (c *Client) refreshToken() error {
	return errors.New("not implemented")
}

/*
func (c *Client) watchToken() {

	attempt := 0
	sleep := time.Second

	for {
		attempt++
		err := c.waitExpiration()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}

			slog.Error("finam auth token expiration waiting failed", "account", c.accountID, "attempt", attempt, "msg", err.Error())

			if attempt > 10 {
				slog.Error("max attempts reached during wait the finam auth token expiration", "account", c.accountID)
				c.Close()
				return
			}

			select {
			case <-c.ctx.Done():
				return
			case <-time.After(sleep):
				sleep = sleep * 2
			}

			continue
		}

		attempt = 0
		sleep = time.Second

		for {
			attempt++
			err = c.refreshToken()
			if err == nil {
				attempt = 0
				sleep = time.Second
				break
			}

			if errors.Is(err, context.Canceled) {
				return
			}

			slog.Error("finam auth token refresh failed", "account", c.accountID, "attempt", attempt, "msg", err.Error())

			if attempt > 10 {
				slog.Error("max attempts reached during the finam auth token refresh", "account", c.accountID)
				c.Close()
				return
			}

			select {
			case <-c.ctx.Done():
				return
			case <-time.After(sleep):
				sleep = sleep * 2
			}
		}
	}
}

func (c *Client) getToken() string {
	return c.tokenStore.Get()
}

func (c *Client) refreshToken() error {
	ctx, cancel := context.WithTimeout(c.ctx, 30*time.Second)
	defer cancel()
	resp, err := c.service.auth.Auth(ctx, &auth.AuthRequest{
		Secret: c.secret,
	})
	if err != nil {
		if tools.IsGRPCCancelled(err) {
			return context.Canceled
		}
		return err
	}
	c.tokenStore.Set(resp.GetToken())
	slog.Debug("refreshed finam auth token", "account", c.accountID, "token", TokenSuffix(resp.GetToken()))
	return nil
}

func (c *Client) getTokenExpitation() (time.Time, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 30*time.Second)
	defer cancel()
	resp, err := c.service.auth.TokenDetails(ctx, &auth.TokenDetailsRequest{
		Token: c.getToken(),
	})
	if err != nil {
		if tools.IsGRPCCancelled(err) {
			return time.Time{}, context.Canceled
		}
		return time.Time{}, err
	}
	return resp.ExpiresAt.AsTime(), nil
}

func (c *Client) waitExpiration() error {

	attempt := 0
	sleep := time.Second

	for {
		attempt++

		expires, err := c.getTokenExpitation()
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return context.Canceled
			}

~			slog.Error("finam auth token getting expiration safailed", "account", c.accountID, "attempt", attempt, "msg", err.Error())

			if attempt > 10 {
				slog.Error("max attempts reached during wait the finam auth token expiration", "account", c.accountID)
				c.Close()
				return
			}

			select {
			case <-c.ctx.Done():
				return
			case <-time.After(sleep):
				sleep = sleep * 2
			}

			continue
		}
	}

	slog.Debug("recognized the finam auth token expiration", "date", expires.Format("2006-01-02 15:04"))

	select {
	case <-c.ctx.Done():
		return context.Canceled

	case <-time.After(time.Until(expires) - 10*time.Minute):
		return nil
	}
}
*/
