# API

/login
/logout
/users/:userId
/user (current user, authenticated)
/user/jobs
/user/:orgId/jobs

/organizations?q=<search>
/organizations/:orgId
/organizations/:orgId/employees
/organizations/:orgId/employees/:emplId
/organizations/:orgId/roles

/organizations/:orgId/products
Authorization:
- userId from token (Authroization: "Bearer ...")
- orgID from api
- GetUserJob(userID, orgID)
```
{
  "employee": {
    "user_id": userID,
    "organization_id": orgID,
    "role": {
      "name": roleName,
      "permissions": [
        {
          "module": moduleName,
          "crud": {
            "create": c,
            "read": r,
            "update": u,
            "delete": d,
          },
        }
      ],
    }
  }
}
```

# Services and interfaces
GetUserJobs(userID string)
GetUserJob(userID string, orgID string)
