# Tenant Model 

A **Tenant** represents a logical customer or organization using the application. All users, projects, and resources are scoped under a tenant, ensuring strict **data isolation** between different customers.

---

## What is a Tenant?

A tenant is a high-level unit of ownership in the system. Each tenant:

- Has its own users, projects, and settings
- Is fully isolated from other tenants
- Acts as the root context for access control and ownership
- Can be managed independently (created, suspended, deleted)

---

## Tenant Isolation

Multi-tenancy is enforced by:
- Tenant ID scoping: All entities include a tenant_id field
- Authentication context: On login, user's tenant_id is stored in session
- Middleware: Injects tenant into request context automatically
- Usecases: Always perform logic within the current tenant scope

This ensures no user or process can access data outside their assigned tenant.

---

## Testing Considerations

In unit and integration tests, simulate multiple tenants to ensure:
- Data is not leaked between tenants
- Usecases reject invalid cross-tenant access
- Deletion or suspension of one tenant doesn't affect others

---

## Relationships

| Entity    | Relationship                |
| --------- | --------------------------- |
| `Tenant`  | Has many `Users`            |
| `Tenant`  | Has many `Projects`         |
| `User`    | Belongs to one `Tenant`     |
| `Session` | Scoped by user's `TenantID` |

---

## Summary

| Feature            | Value                                     |
| ------------------ | ----------------------------------------- |
| Role               | Root context for access and data grouping |
| Storage Strategy   | Shared DB tables with `tenant_id` FK      |
| Access Enforcement | Middleware + usecase layer                |
| Authentication     | Session-based, tied to a specific tenant  |
| Testing            | Requires multi-tenant test coverage       |

--- 