-- +goose Up
INSERT INTO space_permissions (code, name, description) VALUES
('space:read',            '读取空间',     '查看空间信息'),
('space:update',          '编辑空间',     '修改空间设置'),
('space:delete',          '删除空间',     '移入回收站'),
('space:member:invite',   '邀请成员',     '添加空间成员'),
('space:member:remove',   '移除成员',     '移除空间成员'),
('space:member:role',     '修改成员角色',  '修改成员的空间角色'),
('space:document:create', '创建文档',  '在空间内新建文档'),
('space:document:delete', '删除文档',  '删除空间内的文档')
ON CONFLICT (code) DO NOTHING;

-- +goose Down
DELETE FROM space_permissions
  WHERE code IN (
    'space:read',
    'space:update',
    'space:delete',
    'space:member:invite',
    'space:member:remove',
    'space:member:role',
    'space:document:create',
    'space:document:delete'
  );
