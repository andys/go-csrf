# go-csrf
Golang helper library for generating and verifying anti-CSRF tokens

Example init:

```go
import (
  csrf "github.com/andys/go-csrf"
)

   // Initialize CSRF protection
   csrfSecret := make([]byte, 32)
   if _, err := rand.Read(csrfSecret); err != nil {
           log.Fatal().Err(err).Msg("Failed to generate CSRF secret")
   }
   csrf.Setup(csrfSecret)

```

Example Fiber middleware:

```go
   // CSRF protection middleware - verify token on POST requests
   app.Use(func(c *fiber.Ctx) error {
           if c.Method() == "POST" {
                   token := c.FormValue("csrf_token")
                   if !csrf.IsTokenValid(token) {
                           return c.Status(fiber.StatusForbidden).SendString("Invalid CSRF token")
                   }
           }
           return c.Next()
   })
```

Example use:

```html
<input type="hidden" name="csrf_token" value="{%s csrf.GenerateToken() %}">
```
