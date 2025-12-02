---
apply: always
---

# Pull Request Message Rule

## Purpose
This rule ensures that all pull request descriptions follow the project's standardized template and include all required information for proper review and documentation.

## Required Template Structure

All pull requests messages MUST follow this structure:

### 1. Summary Section
- **Required**: Yes
- **Description**: A clear and concise description of the changes in the PR
- **Guidelines**:
  - Explain WHAT changed
  - Explain WHY the change was necessary
  - Reference the diff.log output if available to understand the scope of changes
  - Use present tense (e.g., "Add feature" not "Added feature")
  - Be specific and avoid vague descriptions

### 2. Type of Change
- **Required**: Yes
- **Description**: Exactly ONE checkbox must be selected
- **Valid Options**:
  - `[ ] Feature` - New functionality or enhancement
  - `[ ] Fix` - Bug fix or correction
  - `[ ] Refactor` - Code restructuring without functional changes
  - `[ ] Documentation` - Documentation updates only
  - `[ ] Other` - Any change that doesn't fit the above categories

**Format**: Use `[x]` to mark the selected option

### 3. Additional Notes
- **Required**: No (optional)
- **Description**: Extra context, implementation details, screenshots, or related information
- **Suggestions**:
  - Breaking changes and migration notes
  - Performance implications
  - Testing approach
  - Related issues or PRs
  - Screenshots for UI changes
  - Deployment considerations

## Validation Checklist

When reviewing or creating a PR message, verify:

- [ ] Summary section is present and non-empty
- [ ] Summary clearly describes the changes
- [ ] Exactly one "Type of Change" checkbox is marked with `[x]`
- [ ] No placeholder text remains (e.g., "TODO", "Update this")
- [ ] PR title is concise and descriptive
- [ ] Changes align with the stated type
- [ ] Breaking changes are documented in Additional Notes (if applicable)

## Examples

### ✅ Good PR Message

```markdown
# Pull Request

## Summary

Add user authentication middleware to protect admin endpoints. This change implements JWT-based authentication that validates tokens before allowing access to /admin/* routes. Previously, these endpoints were accessible without authentication, posing a security risk.

---

## Type of Change

Please select one:

- [x] Feature
- [ ] Fix
- [ ] Refactor
- [ ] Documentation
- [ ] Other

---

## Additional Notes (optional)

- Requires new environment variable: JWT_SECRET
- All admin endpoints now return 401 for unauthenticated requests
- Added unit tests for middleware
- Updated API documentation
```
