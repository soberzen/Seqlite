-- +goose Up
INSERT INTO document_permissions (code, name, description) VALUES
('document:read',           '读取文档',     '查看文档内容'),
('document:update',         '编辑文档',     '修改文档内容'),
('document:delete',         '删除文档',     '移入回收站'),
('document:share',          '分享文档',     '分配他人文档角色'),
('document:comment',        '评论文档',     '发表评论'),
('document:comment:delete', '删除任意评论', '删除他人评论'),
('document:manage',         '管理文档权限', '修改文档级授权')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DELETE FROM document_permissions
  WHERE code IN (
    'document:read',
    'document:update',
    'document:delete',
    'document:share',
    'document:comment',
    'document:comment:delete',
    'document:manage'
  );
