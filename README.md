# Gin Auth

Implementation of authentication and authorization by middleware functions in Gin server

``` go
e.PUT("/post/:id",
    handle.JwtAuthenticationRequiredMw(jwtService),
    handle.UpdatePost(postRepo),
)

e.PUT("/post/force/:id",
    handle.JwtAuthenticationRequiredMw(jwtService),
    handle.JwtAuthorizationHasAnyRoleMv(auth.RoleAdmin, auth.RoleManager),
    handle.UpdatePostForcibly(postRepo),
)
```
