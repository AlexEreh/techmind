# API Integration Guide - DataRush Frontend

## ✅ Completed Setup

### Configuration
- **Base URL**: Configured in `/lib/api/config.ts`
- **Environment Variable**: `NEXT_PUBLIC_API_URL` (defaults to `http://localhost:8080`)
- **Axios Client**: Configured with auth token interceptor and 401 handler

### API Client Structure
```
/lib/api/
├── config.ts          # Axios client setup & BASE_URL
├── types.ts           # All TypeScript interfaces
├── auth.ts            # Authentication endpoints
├── company.ts         # Company & user management
├── folders.ts         # Folder management with lazy loading
├── documents.ts       # Document management
├── tags.ts            # Tag management
└── senders.ts         # Sender management
```

---

## 📋 API Methods Reference

### Authentication (`authApi`)

#### Login
```typescript
import { authApi } from '@/lib/api/auth';

const response = await authApi.login({
  email: 'user@example.com',
  password: 'password123'
});
// Returns: { token: string, expires_at: string }
```

#### Register
```typescript
const response = await authApi.register({
  email: 'user@example.com',
  password: 'password123',
  name: 'John Doe'
});
// Returns: { token: string, user: User }
```

#### Logout
```typescript
authApi.logout(); // Clears localStorage and redirects to /login
```

---

### Companies (`companyApi`)

#### Get My Companies
```typescript
import { companyApi } from '@/lib/api/company';

const { companies } = await companyApi.getUserCompanies();
// Returns: { companies: Company[] }
// ⚠️ Backend endpoint needed: GET /private/companies/my
```

#### Get Company Users (with JOIN)
```typescript
const { users, total } = await companyApi.getCompanyUsers(companyId);
// Returns: { users: CompanyUserWithDetails[], total: number }
// ⚠️ Backend endpoint needed: GET /private/companies/{id}/users
// Each user includes: { id, name, email, role, company_user_id }
```

#### Update User Role
```typescript
await companyApi.updateUserRole(companyUserId, newRole);
// ⚠️ Backend endpoint needed: PUT /private/company-users/{id}/role
```

#### Remove User from Company
```typescript
await companyApi.removeUser(companyUserId);
// ⚠️ Backend endpoint needed: DELETE /private/company-users/{id}
```

#### Invite User
```typescript
await companyApi.inviteUser(companyId, 'newuser@example.com', UserRole.MEMBER);
// ⚠️ Backend endpoint needed: POST /private/companies/{id}/invite
```

---

### Folders (`foldersApi`) - ✅ Lazy Loading Supported

#### Get Folders by Parent (Lazy Loading)
```typescript
import { foldersApi } from '@/lib/api/folders';

// Get root folders
const { folders, total } = await foldersApi.getByParent(companyId);

// Get subfolders
const { folders, total } = await foldersApi.getByParent(companyId, parentId);

// ✅ Backend endpoint exists: POST /private/folders/by-parent
```

#### Get All Company Folders
```typescript
const { folders, total } = await foldersApi.getAllCompanyFolders(companyId);
// ✅ Backend endpoint exists: GET /private/folders/company/{id}
```

#### Get Folder by ID
```typescript
const folder = await foldersApi.getById(folderId);
// ✅ Backend endpoint exists: GET /private/folders/{id}
```

#### Create Folder
```typescript
const folder = await foldersApi.create({
  company_id: companyId,
  name: 'New Folder',
  parent_id: parentId // optional
});
// ✅ Backend endpoint exists: POST /private/folders
```

#### Rename Folder
```typescript
const folder = await foldersApi.rename(folderId, 'New Name');
// ✅ Backend endpoint exists: PUT /private/folders/{id}/rename
```

#### Delete Folder
```typescript
await foldersApi.delete(folderId);
// ✅ Backend endpoint exists: DELETE /private/folders/{id}
```

---

### Documents (`documentsApi`)

#### Search Documents
```typescript
import { documentsApi } from '@/lib/api/documents';

const { documents, total } = await documentsApi.search({
  company_id: companyId,
  query: 'invoice',
  folder_id: folderId, // optional
  tag_ids: ['tag1', 'tag2'], // optional
  sender_id: senderId, // optional
  page: 1,
  page_size: 20
});
// ✅ Backend endpoint exists: POST /private/documents/search
```

#### Get Document by ID
```typescript
const document = await documentsApi.getById(documentId);
// ✅ Backend endpoint exists: GET /private/documents/{id}
```

#### Get Documents by Folder
```typescript
const { documents, total } = await documentsApi.getByFolder(folderId);
// ✅ Backend endpoint exists: GET /private/documents/folder/{id}
```

#### Get Documents by Company
```typescript
const { documents, total } = await documentsApi.getByCompany(companyId);
// ✅ Backend endpoint exists: GET /private/documents/company/{id}
```

#### Upload Document
```typescript
const document = await documentsApi.upload({
  company_id: companyId,
  name: 'document.pdf',
  file: fileObject,
  folder_id: folderId, // optional
  sender_id: senderId // optional
});
// ✅ Backend endpoint exists: POST /private/documents
```

#### Update Document
```typescript
const document = await documentsApi.update(documentId, {
  name: 'new-name.pdf',
  folder_id: newFolderId,
  sender_id: senderId
});
// ✅ Backend endpoint exists: PUT /private/documents/{id}
```

#### Delete Document
```typescript
await documentsApi.delete(documentId);
// ✅ Backend endpoint exists: DELETE /private/documents/{id}
```

#### Get Download URL
```typescript
const { url, expires_at } = await documentsApi.getDownloadUrl(documentId);
// ✅ Backend endpoint exists: GET /private/documents/{id}/download
```

#### Get Preview URL
```typescript
const { url, expires_at } = await documentsApi.getPreviewUrl(documentId);
// ✅ Backend endpoint exists: GET /private/documents/{id}/preview
```

---

### Tags (`tagsApi`)

#### Get Company Tags
```typescript
import { tagsApi } from '@/lib/api/tags';

const { tags, total } = await tagsApi.getByCompany(companyId);
// ✅ Backend endpoint exists: GET /private/document-tags/company/{id}
```

#### Get Document Tags
```typescript
const { tags, total } = await tagsApi.getByDocument(documentId);
// ✅ Backend endpoint exists: GET /private/document-tags/document/{id}
```

#### Get Tag by ID
```typescript
const tag = await tagsApi.getById(tagId);
// ✅ Backend endpoint exists: GET /private/document-tags/tags/{id}
```

#### Create Tag
```typescript
const tag = await tagsApi.create({
  company_id: companyId,
  name: 'Important'
});
// ✅ Backend endpoint exists: POST /private/document-tags/tags
```

#### Update Tag
```typescript
const tag = await tagsApi.update(tagId, { name: 'Very Important' });
// ✅ Backend endpoint exists: PUT /private/document-tags/tags/{id}
```

#### Delete Tag
```typescript
await tagsApi.delete(tagId);
// ✅ Backend endpoint exists: DELETE /private/document-tags/tags/{id}
```

#### Add Tag to Document
```typescript
await tagsApi.addToDocument(documentId, tagId);
// ✅ Backend endpoint exists: POST /private/document-tags/add
```

#### Remove Tag from Document
```typescript
await tagsApi.removeFromDocument(documentId, tagId);
// ✅ Backend endpoint exists: POST /private/document-tags/remove
```

---

### Senders (`sendersApi`)

#### Get Company Senders
```typescript
import { sendersApi } from '@/lib/api/senders';

const { senders, total } = await sendersApi.getByCompany(companyId);
// ⚠️ Backend endpoint needed: GET /private/senders/company/{id}
```

#### Get Sender by ID
```typescript
const sender = await sendersApi.getById(senderId);
// ⚠️ Backend endpoint needed: GET /private/senders/{id}
```

#### Create Sender
```typescript
const sender = await sendersApi.create({
  company_id: companyId,
  name: 'ACME Corp',
  email: 'contact@acme.com' // optional
});
// ⚠️ Backend endpoint needed: POST /private/senders
```

#### Update Sender
```typescript
const sender = await sendersApi.update(senderId, {
  name: 'ACME Corporation',
  email: 'info@acme.com'
});
// ⚠️ Backend endpoint needed: PUT /private/senders/{id}
```

#### Delete Sender
```typescript
await sendersApi.delete(senderId);
// ⚠️ Backend endpoint needed: DELETE /private/senders/{id}
```

---

## 🎨 Using AuthContext

```typescript
'use client';
import { useAuth } from '@/contexts/AuthContext';

function MyComponent() {
  const { 
    user,              // Current user
    companies,         // User's companies
    currentCompany,    // Selected company
    setCurrentCompany, // Switch company
    login,             // Login function
    register,          // Register function
    logout,            // Logout function
    isLoading          // Loading state
  } = useAuth();

  // Use in your component
  if (isLoading) return <div>Loading...</div>;
  if (!user) return <div>Not authenticated</div>;

  return <div>Hello, {user.name}!</div>;
}
```

---

## 🔑 User Roles

```typescript
import { UserRole } from '@/lib/api/types';

// Available roles:
UserRole.OWNER   // 0 - Full control
UserRole.ADMIN   // 1 - Manage users
UserRole.MEMBER  // 2 - Upload/edit documents
UserRole.VIEWER  // 3 - Read-only
```

---

## 📝 Type Definitions

### User
```typescript
interface User {
  id: string;
  email: string;
  name: string;
}
```

### Company
```typescript
interface Company {
  id: string;
  name: string;
}
```

### CompanyUserWithDetails (from JOIN)
```typescript
interface CompanyUserWithDetails {
  id: string;           // user id
  name: string;
  email: string;
  role: number;         // role in the company
  company_user_id: string; // id of company_users record
}
```

### Folder
```typescript
interface Folder {
  id: string;
  company_id: string;
  parent_folder_id?: string;
  name: string;
  size: number;   // total size in bytes
  count: number;  // number of documents
  children?: Folder[]; // for tree structure
}
```

### Document
```typescript
interface Document {
  id: string;
  company_id: string;
  folder_id?: string;
  name: string;
  file_path: string;
  preview_file_path?: string;
  preview_url?: string;
  download_url?: string;
  file_size: number;
  mime_type: string;
  checksum: string;
  sender_id?: string;
  created_at: string;
  tags?: Tag[];
}
```

### Tag
```typescript
interface Tag {
  id: string;
  company_id: string;
  name: string;
}
```

### Sender
```typescript
interface Sender {
  id: string;
  name: string;
  email?: string;
}
```

---

## ⚙️ Environment Setup

Create `.env.local` file in the frontend root:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## 🚀 Key Features Implemented

✅ **Lazy Loading for Folders**: Use `foldersApi.getByParent()` to load subfolders on demand  
✅ **User Role Management**: Complete type-safe role system  
✅ **Company User JOIN**: Get user details with roles via `getCompanyUsers()`  
✅ **Auth Context**: Ready-to-use `useAuth()` hook  
✅ **Token Interceptor**: Automatic auth token injection  
✅ **401 Handler**: Auto-redirect to login on auth failure  
✅ **Type Safety**: Full TypeScript support for all API calls

---

## ⚠️ Backend Endpoints Still Needed

See `api-endpoints-needed.md` for complete list of missing endpoints.

**High Priority**:
- GET /private/companies/my
- GET /private/companies/{id}/users
- POST /private/companies
- GET /private/senders/company/{id}
- POST /private/senders

---

## 📦 Next Steps

1. Create UI components for file manager (Obsidian-style)
2. Implement search page
3. Implement settings page with tag/user management
4. Add file upload drag & drop
5. Add document preview
6. Add folder tree component with lazy loading

Frontend is **100% ready** to consume the API once backend implements missing endpoints!

