package ent

// DialectName returns the active SQL dialect for the client.
func (c *Client) DialectName() string {
	if c == nil || c.driver == nil {
		return ""
	}

	if d, ok := c.driver.(interface{ Dialect() string }); ok {
		return d.Dialect()
	}

	return ""
}
