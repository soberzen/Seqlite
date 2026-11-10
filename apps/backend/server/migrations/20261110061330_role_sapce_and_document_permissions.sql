-- +goose Up
INSERT INTO role_space_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN space_permissions p
WHERE
  (r.code = 'owner' AND p.code IN (
    'space:read','space:update','space:delete',
    'space:member:invite','space:member:remove','space:member:role',
    'space:document:create','space:document:delete'
  ))
OR
  (r.code = 'admin' AND p.code IN (
    'space:read','space:update',
    'space:member:invite','space:member:remove','space:member:role',
    'space:document:create','space:document:delete'
  ))
OR (r.code = 'editor' AND p.code IN ('space:read','space:document:create'))
OR (r.code IN ('viewer','commenter') AND p.code = 'space:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;


INSERT INTO role_document_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
CROSS JOIN document_permissions p
WHERE
  (r.code = 'owner'  AND p.code IN (
    'document:read','document:update','document:delete',
    'document:share','document:comment','document:comment:delete','document:manage'
  ))
OR
  (r.code = 'admin'  AND p.code IN (
    'document:read','document:update','document:delete',
    'document:share','document:comment','document:comment:delete'
  ))
OR
  (r.code = 'editor' AND p.code IN (
    'document:read','document:update','document:share','document:comment'
  ))
OR (r.code = 'viewer'    AND p.code = 'document:read')
OR (r.code = 'commenter' AND p.code IN ('document:read','document:comment'))
ON CONFLICT (role_id, permission_id) DO NOTHING;


-- +goose Down
DELETE FROM role_space_permissions
WHERE role_id IN
( SELECT id FROM roles
  WHERE code IN ('owner','admin','editor','viewer','commenter')
)
AND permission_id IN (SELECT id FROM space_permissions WHERE code LIKE 'space:%');

DELETE FROM role_document_permissions
WHERE role_id IN
( SELECT id FROM roles
  WHERE code IN ('owner','admin','editor','viewer','commenter')
)
AND permission_id IN (SELECT id FROM document_permissions WHERE code LIKE 'document:%');
